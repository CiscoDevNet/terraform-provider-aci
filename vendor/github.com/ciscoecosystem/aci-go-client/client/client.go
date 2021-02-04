package client

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const authPayload = `{
	"aaaUser" : {
		"attributes" : {
			"name" : "%s",
			"pwd" : "%s"
		}
	}
}`

// Used authAppPayload to authenticate against the APIC using:
// AppName, App Certificate DN, Signed Request
const authAppPayload = `{
	"aaaAppToken" : {
		"attributes" : {
			"appName" : "%s"
		}
	}
}`

// Default timeout for NGINX in ACI is 90 Seconds.
// Allow the client to set a shorter or longer time depending on their
// environment
const DefaultReqTimeoutVal uint32 = 100
const DefaultMOURL = "/api/node/mo"

// Client is the main entry point
type Client struct {
	BaseURL            *url.URL
	MOURL              string
	httpClient         *http.Client
	AuthToken          *Auth
	l                  sync.Mutex
	username           string
	password           string
	privatekey         string
	adminCert          string
	insecure           bool
	reqTimeoutSet      bool
	reqTimeoutVal      uint32
	proxyUrl           string
	preserveBaseUrlRef bool
	skipLoggingPayload bool
	appUserName        string
	*ServiceManager
}

// singleton implementation of a client
var clientImpl *Client

type Option func(*Client)

func Insecure(insecure bool) Option {
	return func(client *Client) {
		client.insecure = insecure
	}
}

func MoURL(moURL string) Option {
	return func(sm *Client) {
		sm.MOURL = moURL
	}
}

func Password(password string) Option {
	return func(client *Client) {
		client.password = password
	}
}

func PrivateKey(privatekey string) Option {
	return func(client *Client) {
		client.privatekey = privatekey
	}
}

func AdminCert(adminCert string) Option {
	return func(client *Client) {
		client.adminCert = adminCert
	}
}

func AppUserName(appUserName string) Option {
	return func(client *Client) {
		client.appUserName = appUserName
	}
}

func ProxyUrl(pUrl string) Option {
	return func(client *Client) {
		client.proxyUrl = pUrl
	}
}

// HttpClient option: allows for caller to set 'httpClient' with 'Transport'.
// When this option is set 'client.proxyUrl' option is ignored.
func HttpClient(httpcl *http.Client) Option {
	return func(client *Client) {
		client.httpClient = httpcl
	}
}

func SkipLoggingPayload(skipLoggingPayload bool) Option {
	return func(client *Client) {
		client.skipLoggingPayload = skipLoggingPayload
	}
}

func PreserveBaseUrlRef(preserveBaseUrlRef bool) Option {
	return func(client *Client) {
		client.preserveBaseUrlRef = preserveBaseUrlRef
	}
}

func ReqTimeout(timeout uint32) Option {
	return func(client *Client) {
		client.reqTimeoutSet = true
		client.reqTimeoutVal = timeout
	}
}

func initClient(clientUrl, username string, options ...Option) *Client {
	var transport *http.Transport
	bUrl, err := url.Parse(clientUrl)
	if err != nil {
		// cannot move forward if url is undefined
		log.Fatal(err)
	}
	client := &Client{
		BaseURL:    bUrl,
		username:   username,
		MOURL:      DefaultMOURL,
	}

	for _, option := range options {
		option(client)
	}

	if client.httpClient == nil {
		transport = client.useInsecureHTTPClient(client.insecure)
		if client.proxyUrl != "" {
			transport = client.configProxy(transport)
		}
		client.httpClient = &http.Client{
			Transport: transport,
		}
	}

	var timeout time.Duration
	if client.reqTimeoutSet {
		timeout = time.Second * time.Duration(client.reqTimeoutVal)
	} else {
		timeout = time.Second * time.Duration(DefaultReqTimeoutVal)
	}

	client.httpClient.Timeout = timeout
	client.ServiceManager = NewServiceManager(client.MOURL, client)
	return client
}

// GetClient returns a singleton
func GetClient(clientUrl, username string, options ...Option) *Client {
	if clientImpl == nil {
		clientImpl = initClient(clientUrl, username, options...)
	} else {
		// making sure it is the same client
		bUrl, err := url.Parse(clientUrl)
		if err != nil {
			// cannot move forward if url is undefined
			log.Fatal(err)
		}
		if bUrl != clientImpl.BaseURL {
			clientImpl = initClient(clientUrl, username, options...)
		}
	}
	return clientImpl
}

// NewClient returns a new Instance of the client - allowing for simultaneous connections to the same APIC
func NewClient(clientUrl, username string, options ...Option) *Client {
	// making sure it is the same client
	_, err := url.Parse(clientUrl)
	if err != nil {
		// cannot move forward if url is undefined
		log.Fatal(err)
	}

	// initClient always returns a new struct, so always create a new pointer to allow for
	// multiple object instances
	newClientImpl := initClient(clientUrl, username, options...)

	return newClientImpl
}

func (c *Client) configProxy(transport *http.Transport) *http.Transport {
	pUrl, err := url.Parse(c.proxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	transport.Proxy = http.ProxyURL(pUrl)
	return transport

}

func (c *Client) useInsecureHTTPClient(insecure bool) *http.Transport {
	// proxyUrl, _ := url.Parse("http://10.0.1.167:3128")
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       insecure,
			MinVersion:               tls.VersionTLS11,
			MaxVersion:               tls.VersionTLS12,
		},
	}

	return transport

}

func (c *Client) MakeRestRequest(method string, rpath string, body *container.Container, authenticated bool) (*http.Request, error) {

	pathURL, err := url.Parse(rpath)
	if err != nil {
		return nil, err
	}

	fURL, err := url.Parse(c.BaseURL.String())
	if err != nil {
		return nil, err
	}

	if c.preserveBaseUrlRef {
		// Default is false for preserveBaseUrlRef - matching original behavior to strip out BaseURL
		fURLStr := fURL.String() + pathURL.String()
		fURL, err = url.Parse(fURLStr)
		if err != nil {
			return nil, err
		}
	} else {
		// Original behavior to strip down BaseURL
		fURL = fURL.ResolveReference(pathURL)
	}

	var req *http.Request
	log.Printf("[DEBUG] BaseURL: %s, pathURL: %s, fURL: %s", c.BaseURL.String(), pathURL.String(), fURL.String())
	if method == "GET" {
		req, err = http.NewRequest(method, fURL.String(), nil)
	} else {
		req, err = http.NewRequest(method, fURL.String(), bytes.NewBuffer((body.Bytes())))
	}
	if err != nil {
		return nil, err
	}

	if c.skipLoggingPayload {
		log.Printf("HTTP request %s %s", method, rpath)
	} else {
		log.Printf("HTTP request %s %s %v", method, rpath, req)
	}
	if authenticated {
		req, err = c.InjectAuthenticationHeader(req, rpath)
		if err != nil {
			return req, err
		}
	}

	if !c.skipLoggingPayload {
		log.Printf("HTTP request after injection %s %s %v", method, rpath, req)
	}

	return req, nil
}

// Authenticate is used to
func (c *Client) Authenticate() error {
	method := "POST"
	path := "/api/aaaLogin.json"
	authenticated := false

	// Adding the follwing replace allows support for (1) Login Domains, where login is in the format of: apic#LOCAL\admin2
	// (2) escapes out the password to support scenarios where the user password includes backslashes
	escUserName := strings.ReplaceAll(c.username, `\`, `\\`)
	escPwd := strings.ReplaceAll(c.password, `\`, `\\`)
	body, err := container.ParseJSON([]byte(fmt.Sprintf(authPayload, escUserName, escPwd)))
	if c.appUserName != "" {
		path = "/api/requestAppToken.json"
		body, err = container.ParseJSON([]byte(fmt.Sprintf(authAppPayload, c.appUserName)))
		authenticated = true
	}

	if err != nil {
		return err
	}

	req, err := c.MakeRestRequest(method, path, body, authenticated)
	obj, _, err := c.Do(req)

	if err != nil {
		return err
	}
	if obj == nil {
		return errors.New("Empty response")
	}
	err = CheckForErrors(obj, method, c.skipLoggingPayload)
	if err != nil {
		return err
	}

	token := obj.S("imdata").Index(0).S("aaaLogin", "attributes", "token").String()
	creationTimeStr := stripQuotes(obj.S("imdata").Index(0).S("aaaLogin", "attributes", "creationTime").String())
	refreshTimeStr := stripQuotes(obj.S("imdata").Index(0).S("aaaLogin", "attributes", "refreshTimeoutSeconds").String())

	creationTimeInt, err := StrtoInt(creationTimeStr, 10, 64)
	if err != nil {
		return err
	}
	refreshTimeInt, err := StrtoInt(refreshTimeStr, 10, 64)
	if err != nil {
		return err
	}
	if token == "" {
		return errors.New("Invalid Username or Password")
	}

	if c.AuthToken == nil {
		c.AuthToken = &Auth{}
	}
	c.AuthToken.Token = stripQuotes(token)
	c.AuthToken.apicCreatedAt = time.Unix(creationTimeInt, 0)
	c.AuthToken.realCreatedAt = time.Now()
	c.AuthToken.CalculateExpiry(refreshTimeInt)
	c.AuthToken.CaclulateOffset()

	return nil
}
func StrtoInt(s string, startIndex int, bitSize int) (int64, error) {
	return strconv.ParseInt(s, startIndex, bitSize)

}
func (c *Client) Do(req *http.Request) (*container.Container, *http.Response, error) {
	log.Printf("[DEBUG] Begining DO method %s", req.URL.String())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if !c.skipLoggingPayload {
		log.Printf("\n\n\n HTTP request: %v", req.Body)
	}
	log.Printf("\nHTTP Request: %s %s", req.Method, req.URL.String())
	if !c.skipLoggingPayload {
		log.Printf("\nHTTP Response: %d %s %v", resp.StatusCode, resp.Status, resp)
	} else {
		log.Printf("\nHTTP Response: %d %s", resp.StatusCode, resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)
	resp.Body.Close()
	if !c.skipLoggingPayload {
		log.Printf("\n HTTP response unique string %s %s %s", req.Method, req.URL.String(), bodyStr)
	}
	obj, err := container.ParseJSON(bodyBytes)

	if err != nil {

		log.Printf("Error occured while json parsing %+v", err)
		return nil, resp, err
	}
	log.Printf("[DEBUG] Exit from do method")
	return obj, resp, err

}

func (c *Client) DoRaw(req *http.Request) (*http.Response, error) {

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if !c.skipLoggingPayload {
		log.Printf("\n\n\n HTTP request: %v", req.Body)
	}
	log.Printf("\nHTTP Request: %s %s", req.Method, req.URL.String())
	if !c.skipLoggingPayload {
		log.Printf("\nHTTP Response: %d %s %v", resp.StatusCode, resp.Status, resp)
	} else {
		log.Printf("\nHTTP Response: %d %s", resp.StatusCode, resp.Status)
	}

	return resp, err
}

func stripQuotes(word string) string {
	if strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "\""), "\"")
	}
	return word
}

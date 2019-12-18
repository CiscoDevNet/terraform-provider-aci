package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

const DefaultMOURL = "/api/node/mo"

// Client is the main entry point
type Client struct {
	BaseURL    *url.URL
	MOURL      string
	httpClient *http.Client
	AuthToken  *Auth
	username   string
	password   string
	privatekey string
	adminCert  string
	insecure   bool
	proxyUrl   string
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

func ProxyUrl(pUrl string) Option {
	return func(client *Client) {
		client.proxyUrl = pUrl
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
		httpClient: http.DefaultClient,
		MOURL:      DefaultMOURL,
	}

	for _, option := range options {
		option(client)
	}

	if client.insecure {
		transport = client.useInsecureHTTPClient()
	}
	if client.proxyUrl != "" {
		transport = client.configProxy(transport)
	}
	client.httpClient = &http.Client{
		Transport: transport,
	}
	client.ServiceManager = NewServiceManager(client.MOURL, client)
	return client
}

// GetClient returns a singleton
func GetClient(clientUrl, username string, options ...Option) *Client {
	if clientImpl == nil {
		clientImpl = initClient(clientUrl, username, options...)
	}
	return clientImpl
}
func (c *Client) configProxy(transport *http.Transport) *http.Transport {
	pUrl, err := url.Parse(c.proxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	transport.Proxy = http.ProxyURL(pUrl)
	return transport

}
func (c *Client) useInsecureHTTPClient() *http.Transport {
	// proxyUrl, _ := url.Parse("http://10.0.1.167:3128")
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       true,
			MinVersion:               tls.VersionTLS11,
			MaxVersion:               tls.VersionTLS11,
		},
	}

	return transport

}

func (c *Client) MakeRestRequest(method string, path string, body *container.Container, authenticated bool) (*http.Request, error) {

	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	fURL := c.BaseURL.ResolveReference(url)
	var req *http.Request
	if method == "GET" {
		req, err = http.NewRequest(method, fURL.String(), nil)
	} else {
		req, err = http.NewRequest(method, fURL.String(), bytes.NewBuffer((body.Bytes())))
	}
	if err != nil {
		return nil, err
	}
	log.Printf("HTTP request %s %s %v", method, path, req)

	if authenticated {

		req, err = c.InjectAuthenticationHeader(req, path)
		if err != nil {
			return req, err
		}
	}

	return req, nil
}

// Authenticate is used to
func (c *Client) Authenticate() error {
	method := "POST"
	path := "/api/aaaLogin.json"
	body, err := container.ParseJSON([]byte(fmt.Sprintf(authPayload, c.username, c.password)))

	if err != nil {
		return err
	}

	fmt.Println(body.String())
	req, err := c.MakeRestRequest(method, path, body, false)
	obj, _, err := c.Do(req)

	if err != nil {
		return err
	}
	if obj == nil {
		return errors.New("Empty response")
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
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("\n\n\n HTTP request: %v", req.Body)
	log.Printf("\nHTTP Request: %s %s", req.Method, req.URL.String())
	log.Printf("nHTTP Response: %d %s %v", resp.StatusCode, resp.Status, resp)

	decoder := json.NewDecoder(resp.Body)
	obj, err := container.ParseJSONDecoder(decoder)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Error occurred.")
		return nil, resp, err
	}
	log.Printf("[DEBUG] Exit from do method")
	return obj, resp, err

}

func stripQuotes(word string) string {
	if strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "\""), "\"")
	}
	return word
}

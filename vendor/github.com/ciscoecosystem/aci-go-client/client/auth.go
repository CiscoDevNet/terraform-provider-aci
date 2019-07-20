package client

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Auth struct {
	Token         string
	Expiry        time.Time
	apicCreatedAt time.Time
	realCreatedAt time.Time
	offset        int64
}

func (au *Auth) IsValid() bool {
	if au.Token != "" && au.Expiry.Unix() > au.estimateExpireTime() {
		return true
	}
	return false
}

func (t *Auth) CalculateExpiry(willExpire int64) {
	t.Expiry = time.Unix((t.apicCreatedAt.Unix() + willExpire), 0)
}
func (t *Auth) CaclulateOffset() {
	t.offset = t.apicCreatedAt.Unix() - t.realCreatedAt.Unix()
}

func (t *Auth) estimateExpireTime() int64 {
	return time.Now().Unix() + t.offset
}

func (client *Client) InjectAuthenticationHeader(req *http.Request, path string) (*http.Request, error) {

	if client.password != "" {
		if client.AuthToken == nil || !client.AuthToken.IsValid() {
			fmt.Println(client)
			err := client.Authenticate()
			fmt.Println(client)
			if err != nil {
				return nil, err
			}

		}
		req.AddCookie(&http.Cookie{
			Name:  "APIC-Cookie",
			Value: client.AuthToken.Token,
		})
		return req, nil
	} else if client.privatekey != "" && client.adminCert != "" {
		buffer, _ := ioutil.ReadAll(req.Body)
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buffer))

		req.Body = rdr2
		bodyStr := string(buffer)
		contentStr := ""
		if bodyStr != "{}" {
			contentStr = fmt.Sprintf("%s%s%s", req.Method, path, bodyStr)
		} else {
			contentStr = fmt.Sprintf("%s%s", req.Method, path)

		}
		fmt.Println("Content " + contentStr)
		content := []byte(contentStr)

		signature, err := createSignature(content, client.privatekey)
		fmt.Println("sig" + signature)
		if err != nil {
			return req, err
		}
		req.AddCookie(&http.Cookie{
			Name:  "APIC-Request-Signature",
			Value: signature,
		})
		req.AddCookie(&http.Cookie{
			Name:  "APIC-Certificate-Algorithm",
			Value: "v1.0",
		})
		req.AddCookie(&http.Cookie{
			Name:  "APIC-Certificate-Fingerprint",
			Value: "fingerprint",
		})
		req.AddCookie(&http.Cookie{
			Name:  "APIC-Certificate-DN",
			Value: fmt.Sprintf("uni/userext/user-%s/usercert-%s", client.username, client.adminCert),
		})

		return req, nil

	} else {

		return req, fmt.Errorf("Anyone of password or privatekey/certificate name is must.")
	}

	return req, nil
}

func createSignature(content []byte, keypath string) (string, error) {
	hasher := sha256.New()
	hasher.Write(content)
	privkey, err := loadPrivateKey(keypath)
	if err != nil {
		return "", err
	}
	signedData, err := rsa.SignPKCS1v15(nil, privkey, crypto.SHA256, hasher.Sum(nil))
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(signedData)
	return signature, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parsePrivateKey(data)
}

func parsePrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return privkey, err
	case "PRIVATE KEY":
		parsedresult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		privkey := parsedresult.(*rsa.PrivateKey)
		return privkey, nil
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
}

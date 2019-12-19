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
	"log"
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
	log.Printf("[DEBUG] Begin Injection")
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
		var bodyStr string
		if req.Method != "GET" {
			buffer, _ := ioutil.ReadAll(req.Body)
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buffer))

			req.Body = rdr2
			bodyStr = string(buffer)
		}
		contentStr := ""
		if bodyStr != "{}" {
			contentStr = fmt.Sprintf("%s%s%s", req.Method, path, bodyStr)
		} else {
			contentStr = fmt.Sprintf("%s%s", req.Method, path)

		}
		log.Printf("Content %s", contentStr)
		content := []byte(contentStr)

		signature, err := createSignature(content, client.privatekey)
		log.Printf("signature %s" + signature)
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
		log.Printf("[DEBUG] finished signature creation")
		return req, nil

	} else {

		return req, fmt.Errorf("Anyone of password or privatekey/certificate name is must.")
	}

	return req, nil
}

func createSignature(content []byte, keypath string) (string, error) {
	log.Printf("[DEBUG] Begin Create signature")
	hasher := sha256.New()
	hasher.Write(content)
	log.Printf("[DEBUG] Begin Read private key inside createsignature")

	privkey, err := loadPrivateKey(keypath)
	log.Printf("[DEBUG] finish read private key inside Create signature")

	if err != nil {
		return "", err
	}
	log.Printf("[DEBUG] Begin signing signature")

	signedData, err := rsa.SignPKCS1v15(nil, privkey, crypto.SHA256, hasher.Sum(nil))
	log.Printf("[DEBUG] finish signing signature")

	if err != nil {
		return "", err
	}
	log.Printf("[DEBUG] Begin final encoding signature")

	signature := base64.StdEncoding.EncodeToString(signedData)
	log.Printf("[DEBUG] finish final signature")

	return signature, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	log.Printf("[DEBUG] Begin load private key inside loadPrivateKey")

	data, err := ioutil.ReadFile(path)
	log.Printf("[DEBUG] priavte key read finish  inside loadPrivateKey")

	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] finish load private key inside loadPrivateKey")

	return parsePrivateKey(data)
}

func parsePrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	log.Printf("[DEBUG] Begin parse private key inside parsePrivateKey")

	block, _ := pem.Decode(pemBytes)
	log.Printf("[DEBUG] pem decode finish parse private key inside parsePrivateKey")

	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		log.Printf("[DEBUG] x509 parsing  private key inside parsePrivateKey")

		if err != nil {
			return nil, err
		}
		return privkey, err
	case "PRIVATE KEY":
		parsedresult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		log.Printf("[DEBUG] x509 parsing private key inside parsePrivateKey")

		if err != nil {
			return nil, err
		}
		privkey := parsedresult.(*rsa.PrivateKey)
		log.Printf("[DEBUG] finish private parse private key inside parsePrivateKey")

		return privkey, nil
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
}

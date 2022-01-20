package client

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

type ServiceManager struct {
	MOURL  string
	client *Client
}

func NewServiceManager(moURL string, client *Client) *ServiceManager {

	sm := &ServiceManager{
		MOURL:  moURL,
		client: client,
	}
	return sm
}

func (sm *ServiceManager) Get(dn string) (*container.Container, error) {
	finalURL := fmt.Sprintf("%s/%s.json", sm.MOURL, dn)
	req, err := sm.client.MakeRestRequest("GET", finalURL, nil, true)

	if err != nil {
		return nil, err
	}

	obj, _, err := sm.client.Do(req)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, errors.New("Empty response body")
	}
	log.Printf("[DEBUG] Exit from GET %s", finalURL)
	return obj, CheckForErrors(obj, "GET", sm.client.skipLoggingPayload)
}

func createJsonPayload(payload map[string]string) (*container.Container, error) {
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
			}
		}
	}`, payload["classname"]))

	return container.ParseJSON(containerJSON)
}

func (sm *ServiceManager) Save(obj models.Model) error {

	jsonPayload, _, err := sm.PrepareModel(obj)

	if err != nil {
		return err
	}

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	cont, _, err := sm.client.Do(req)
	if err != nil {
		return err
	}

	return CheckForErrors(cont, "POST", sm.client.skipLoggingPayload)
}

// CheckForErrors parses the response and checks of there is an error attribute in the response
func CheckForErrors(cont *container.Container, method string, skipLoggingPayload bool) error {
	number, err := strconv.Atoi(models.G(cont, "totalCount"))
	if err != nil {
		if !skipLoggingPayload {
			log.Printf("[DEBUG] Exit from errors, Unable to parse error count from response %v", cont)
		} else {
			log.Printf("[DEBUG] Exit from errors %s", err.Error())
		}
		return err
	}
	imdata := cont.S("imdata").Index(0)
	if number > 0 {

		if imdata.Exists("error") {
			errorCode := models.StripQuotes(imdata.Path("error.attributes.code").String())
			// Ignore errors of type "Cannot create object"
			if errorCode == "103" {
				if !skipLoggingPayload {
					log.Printf("[DEBUG] Exit from error 103 %v", cont)
				}
				return nil
			} else if method == "DELETE" && (errorCode == "1" || errorCode == "107") { // Ignore errors of type "Cannot delete object"
				if !skipLoggingPayload {
					log.Printf("[DEBUG] Exit from error 1 or 107 %v", cont)
				}
				return nil
			} else {
				if models.StripQuotes(imdata.Path("error.attributes.text").String()) == "" && errorCode == "403" {
					if !skipLoggingPayload {
						log.Printf("[DEBUG] Exit from authentication error 403 %v", cont)
					}
					return errors.New("Unable to authenticate. Please check your credentials")
				}
				if !skipLoggingPayload {
					log.Printf("[DEBUG] Exit from errors %v", cont)
				}

				return errors.New(models.StripQuotes(imdata.Path("error.attributes.text").String()))
			}
		}

	}

	if imdata.String() == "{}" && method == "GET" {
		if !skipLoggingPayload {
			log.Printf("[DEBUG] Exit from error (Empty response) %v", cont)
		}

		return errors.New("Error retrieving Object: Object may not exists")
	}
	if !skipLoggingPayload {
		log.Printf("[DEBUG] Exit from errors %v", cont)
	}
	return nil
}

func (sm *ServiceManager) Delete(obj models.Model) error {

	jsonPayload, className, err := sm.PrepareModel(obj)

	if err != nil {
		return err
	}

	jsonPayload.Set("deleted", className, "attributes", "status")
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) PostViaURL(url string, obj models.Model) (*container.Container, error) {

	jsonPayload, _, err := sm.PrepareModel(obj)

	if err != nil {
		return nil, err
	}

	req, err := sm.client.MakeRestRequest("POST", url, jsonPayload, true)

	if err != nil {
		return nil, err
	}

	cont, _, err := sm.client.Do(req)
	if !sm.client.skipLoggingPayload {
		log.Printf("PostViaUrl %+v", obj)
	}
	if err != nil {
		return nil, err
	}

	if cont == nil {
		return nil, errors.New("Empty response body")
	}
	return cont, CheckForErrors(cont, "POST", sm.client.skipLoggingPayload)

}

func (sm *ServiceManager) GetViaURL(url string) (*container.Container, error) {
	req, err := sm.client.MakeRestRequest("GET", url, nil, true)

	if err != nil {
		return nil, err
	}

	obj, _, err := sm.client.Do(req)
	if !sm.client.skipLoggingPayload {
		log.Printf("Getvia url %+v", obj)
	}
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, errors.New("Empty response body")
	}
	return obj, CheckForErrors(obj, "GET", sm.client.skipLoggingPayload)

}

func (sm *ServiceManager) DeleteByDn(dn, className string) error {
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"status": "deleted"
			}
		}
	}`, className, dn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil

}
func (sm *ServiceManager) PrepareModel(obj models.Model) (*container.Container, string, error) {
	cont, err := obj.ToMap()
	if err != nil {
		return nil, "", err
	}
	jsonPayload, err := createJsonPayload(cont)
	if err != nil {
		return nil, "", err
	}
	className := cont["classname"]
	delete(cont, "classname")

	for key, value := range cont {
		jsonPayload.Set(value, className, "attributes", key)
	}
	return jsonPayload, className, nil
}

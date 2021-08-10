package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateConfigurationImportPolicy(name string, description string, configImportPattr models.ConfigurationImportPolicyAttributes) (*models.ConfigurationImportPolicy, error) {
	rn := fmt.Sprintf("fabric/configimp-%s", name)
	parentDn := fmt.Sprintf("uni")
	configImportP := models.NewConfigurationImportPolicy(rn, parentDn, description, configImportPattr)
	err := sm.Save(configImportP)
	return configImportP, err
}

func (sm *ServiceManager) ReadConfigurationImportPolicy(name string) (*models.ConfigurationImportPolicy, error) {
	dn := fmt.Sprintf("uni/fabric/configimp-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	configImportP := models.ConfigurationImportPolicyFromContainer(cont)
	return configImportP, nil
}

func (sm *ServiceManager) DeleteConfigurationImportPolicy(name string) error {
	dn := fmt.Sprintf("uni/fabric/configimp-%s", name)
	return sm.DeleteByDn(dn, models.ConfigimportpClassName)
}

func (sm *ServiceManager) UpdateConfigurationImportPolicy(name string, description string, configImportPattr models.ConfigurationImportPolicyAttributes) (*models.ConfigurationImportPolicy, error) {
	rn := fmt.Sprintf("fabric/configimp-%s", name)
	parentDn := fmt.Sprintf("uni")
	configImportP := models.NewConfigurationImportPolicy(rn, parentDn, description, configImportPattr)

	configImportP.Status = "modified"
	err := sm.Save(configImportP)
	return configImportP, err

}

func (sm *ServiceManager) ListConfigurationImportPolicy() ([]*models.ConfigurationImportPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/configImportP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ConfigurationImportPolicyListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationconfigRsImportSourceFromConfigurationImportPolicy(parentDn, tnFileRemotePathName string) error {
	dn := fmt.Sprintf("%s/rsImportSource", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFileRemotePathName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "configRsImportSource", dn, tnFileRemotePathName))

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

func (sm *ServiceManager) DeleteRelationconfigRsImportSourceFromConfigurationImportPolicy(parentDn string) error {
	dn := fmt.Sprintf("%s/rsImportSource", parentDn)
	return sm.DeleteByDn(dn, "configRsImportSource")
}

func (sm *ServiceManager) ReadRelationconfigRsImportSourceFromConfigurationImportPolicy(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "configRsImportSource")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "configRsImportSource")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationtrigRsTriggerableFromConfigurationImportPolicy(parentDn, tnTrigTriggerableName string) error {
	dn := fmt.Sprintf("%s/rsTriggerable", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "trigRsTriggerable", dn, tnTrigTriggerableName))

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

func (sm *ServiceManager) ReadRelationtrigRsTriggerableFromConfigurationImportPolicy(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "trigRsTriggerable")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "trigRsTriggerable")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationconfigRsRemotePathFromConfigurationImportPolicy(parentDn, tnFileRemotePathName string) error {
	dn := fmt.Sprintf("%s/rsRemotePath", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFileRemotePathName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "configRsRemotePath", dn, tnFileRemotePathName))

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

func (sm *ServiceManager) DeleteRelationconfigRsRemotePathFromConfigurationImportPolicy(parentDn string) error {
	dn := fmt.Sprintf("%s/rsRemotePath", parentDn)
	return sm.DeleteByDn(dn, "configRsRemotePath")
}

func (sm *ServiceManager) ReadRelationconfigRsRemotePathFromConfigurationImportPolicy(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "configRsRemotePath")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "configRsRemotePath")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}

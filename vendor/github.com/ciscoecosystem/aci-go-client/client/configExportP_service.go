package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateConfigurationExportPolicy(name string, description string, configExportPattr models.ConfigurationExportPolicyAttributes) (*models.ConfigurationExportPolicy, error) {
	rn := fmt.Sprintf("fabric/configexp-%s", name)
	parentDn := fmt.Sprintf("uni")
	configExportP := models.NewConfigurationExportPolicy(rn, parentDn, description, configExportPattr)
	err := sm.Save(configExportP)
	return configExportP, err
}

func (sm *ServiceManager) ReadConfigurationExportPolicy(name string) (*models.ConfigurationExportPolicy, error) {
	dn := fmt.Sprintf("uni/fabric/configexp-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	configExportP := models.ConfigurationExportPolicyFromContainer(cont)
	return configExportP, nil
}

func (sm *ServiceManager) DeleteConfigurationExportPolicy(name string) error {
	dn := fmt.Sprintf("uni/fabric/configexp-%s", name)
	return sm.DeleteByDn(dn, models.ConfigexportpClassName)
}

func (sm *ServiceManager) UpdateConfigurationExportPolicy(name string, description string, configExportPattr models.ConfigurationExportPolicyAttributes) (*models.ConfigurationExportPolicy, error) {
	rn := fmt.Sprintf("fabric/configexp-%s", name)
	parentDn := fmt.Sprintf("uni")
	configExportP := models.NewConfigurationExportPolicy(rn, parentDn, description, configExportPattr)

	configExportP.Status = "modified"
	err := sm.Save(configExportP)
	return configExportP, err

}

func (sm *ServiceManager) ListConfigurationExportPolicy() ([]*models.ConfigurationExportPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/configExportP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ConfigurationExportPolicyListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationconfigRsExportDestinationFromConfigurationExportPolicy(parentDn, tnFileRemotePathName string) error {
	dn := fmt.Sprintf("%s/rsExportDestination", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFileRemotePathName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "configRsExportDestination", dn, tnFileRemotePathName))

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

func (sm *ServiceManager) DeleteRelationconfigRsExportDestinationFromConfigurationExportPolicy(parentDn string) error {
	dn := fmt.Sprintf("%s/rsExportDestination", parentDn)
	return sm.DeleteByDn(dn, "configRsExportDestination")
}

func (sm *ServiceManager) ReadRelationconfigRsExportDestinationFromConfigurationExportPolicy(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "configRsExportDestination")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "configRsExportDestination")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationtrigRsTriggerableFromConfigurationExportPolicy(parentDn, tnTrigTriggerableName string) error {
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

func (sm *ServiceManager) ReadRelationtrigRsTriggerableFromConfigurationExportPolicy(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationconfigRsRemotePathFromConfigurationExportPolicy(parentDn, tnFileRemotePathName string) error {
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

func (sm *ServiceManager) DeleteRelationconfigRsRemotePathFromConfigurationExportPolicy(parentDn string) error {
	dn := fmt.Sprintf("%s/rsRemotePath", parentDn)
	return sm.DeleteByDn(dn, "configRsRemotePath")
}

func (sm *ServiceManager) ReadRelationconfigRsRemotePathFromConfigurationExportPolicy(parentDn string) (interface{}, error) {
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
func (sm *ServiceManager) CreateRelationconfigRsExportSchedulerFromConfigurationExportPolicy(parentDn, tnTrigSchedPName string) error {
	dn := fmt.Sprintf("%s/rsExportScheduler", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnTrigSchedPName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "configRsExportScheduler", dn, tnTrigSchedPName))

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

func (sm *ServiceManager) DeleteRelationconfigRsExportSchedulerFromConfigurationExportPolicy(parentDn string) error {
	dn := fmt.Sprintf("%s/rsExportScheduler", parentDn)
	return sm.DeleteByDn(dn, "configRsExportScheduler")
}

func (sm *ServiceManager) ReadRelationconfigRsExportSchedulerFromConfigurationExportPolicy(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "configRsExportScheduler")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "configRsExportScheduler")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}

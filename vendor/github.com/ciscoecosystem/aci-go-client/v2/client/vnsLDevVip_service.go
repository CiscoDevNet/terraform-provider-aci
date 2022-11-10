package client

import (
	"encoding/json"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateL4ToL7Devices(name string, tenant string, nameAlias string, vnsLDevVipAttr models.L4ToL7DevicesAttributes) (*models.L4ToL7Devices, error) {
	rn := fmt.Sprintf(models.RnvnsLDevVip, name)
	parentDn := fmt.Sprintf(models.ParentDnvnsLDevVip, tenant)
	vnsLDevVip := models.NewL4ToL7Devices(rn, parentDn, nameAlias, vnsLDevVipAttr)
	err := sm.Save(vnsLDevVip)
	return vnsLDevVip, err
}

func (sm *ServiceManager) ReadL4ToL7Devices(name string, tenant string) (*models.L4ToL7Devices, error) {
	dn := fmt.Sprintf(models.DnvnsLDevVip, tenant, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsLDevVip := models.L4ToL7DevicesFromContainer(cont)
	return vnsLDevVip, nil
}

func (sm *ServiceManager) DeleteL4ToL7Devices(name string, tenant string) error {
	dn := fmt.Sprintf(models.DnvnsLDevVip, tenant, name)
	return sm.DeleteByDn(dn, models.VnsldevvipClassName)
}

func (sm *ServiceManager) UpdateL4ToL7Devices(name string, tenant string, nameAlias string, vnsLDevVipAttr models.L4ToL7DevicesAttributes) (*models.L4ToL7Devices, error) {
	rn := fmt.Sprintf(models.RnvnsLDevVip, name)
	parentDn := fmt.Sprintf(models.ParentDnvnsLDevVip, tenant)
	vnsLDevVip := models.NewL4ToL7Devices(rn, parentDn, nameAlias, vnsLDevVipAttr)
	vnsLDevVip.Status = "modified"
	err := sm.Save(vnsLDevVip)
	return vnsLDevVip, err
}

func (sm *ServiceManager) ListL4ToL7Devices(tenant string) ([]*models.L4ToL7Devices, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/vnsLDevVip.json", models.BaseurlStr, tenant)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.L4ToL7DevicesListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationvnsRsALDevToDomP(parentDn, annotation, switchingMode string, tDn string) error {
	dn := fmt.Sprintf("%s/rsALDevToDomP", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vnsRsALDevToDomP", dn, annotation, tDn))

	attributes := map[string]interface{}{
		"switchingMode": switchingMode,
	}
	var output map[string]interface{}
	err_output := json.Unmarshal([]byte(containerJSON), &output)
	if err_output != nil {
		return err_output
	}
	for _, mo := range output {
		if mo_map, ok := mo.(map[string]interface{}); ok {
			for _, mo_attributes := range mo_map {
				if mo_attributes_map, ok := mo_attributes.(map[string]interface{}); ok {
					for key, value := range attributes {
						if value != "" {
							mo_attributes_map[key] = value
						}
					}

				}
			}
		}

	}
	input, out_err := json.Marshal(output)
	if out_err != nil {
		return out_err
	}
	jsonPayload, err := container.ParseJSON(input)
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
	fmt.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationvnsRsALDevToDomP(parentDn string) error {
	dn := fmt.Sprintf("%s/rsALDevToDomP", parentDn)
	return sm.DeleteByDn(dn, "vnsRsALDevToDomP")
}

func (sm *ServiceManager) ReadRelationvnsRsALDevToDomP(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vnsRsALDevToDomP")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vnsRsALDevToDomP")

	st := make([]map[string]string, 0, 1)
	for _, contItem := range contList {
		paramMap := make(map[string]string)
		paramMap["switchingMode"] = models.G(contItem, "switchingMode")
		paramMap["tDn"] = models.G(contItem, "tDn")
		st = append(st, paramMap)
	}
	return st, err
}

func (sm *ServiceManager) CreateRelationvnsRsALDevToPhysDomP(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsALDevToPhysDomP", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "vnsRsALDevToPhysDomP", dn, annotation, tDn))

	jsonPayload, err := container.ParseJSON(containerJSON)
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
	fmt.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationvnsRsALDevToPhysDomP(parentDn string) error {
	dn := fmt.Sprintf("%s/rsALDevToPhysDomP", parentDn)
	return sm.DeleteByDn(dn, "vnsRsALDevToPhysDomP")
}

func (sm *ServiceManager) ReadRelationvnsRsALDevToPhysDomP(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "vnsRsALDevToPhysDomP")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "vnsRsALDevToPhysDomP")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

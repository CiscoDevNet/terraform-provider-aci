package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateAccessPortSelector(access_port_selector_type string, name string, leaf_interface_profile string, description string, infraHPortSattr models.AccessPortSelectorAttributes) (*models.AccessPortSelector, error) {
	rn := fmt.Sprintf("hports-%s-typ-%s", name, access_port_selector_type)
	parentDn := fmt.Sprintf("uni/infra/accportprof-%s", leaf_interface_profile)
	infraHPortS := models.NewAccessPortSelector(rn, parentDn, description, infraHPortSattr)
	err := sm.Save(infraHPortS)
	return infraHPortS, err
}

func (sm *ServiceManager) ReadAccessPortSelector(access_port_selector_type string, name string, leaf_interface_profile string) (*models.AccessPortSelector, error) {
	dn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s", leaf_interface_profile, name, access_port_selector_type)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraHPortS := models.AccessPortSelectorFromContainer(cont)
	return infraHPortS, nil
}

func (sm *ServiceManager) DeleteAccessPortSelector(access_port_selector_type string, name string, leaf_interface_profile string) error {
	dn := fmt.Sprintf("uni/infra/accportprof-%s/hports-%s-typ-%s", leaf_interface_profile, name, access_port_selector_type)
	return sm.DeleteByDn(dn, models.InfrahportsClassName)
}

func (sm *ServiceManager) UpdateAccessPortSelector(access_port_selector_type string, name string, leaf_interface_profile string, description string, infraHPortSattr models.AccessPortSelectorAttributes) (*models.AccessPortSelector, error) {
	rn := fmt.Sprintf("hports-%s-typ-%s", name, access_port_selector_type)
	parentDn := fmt.Sprintf("uni/infra/accportprof-%s", leaf_interface_profile)
	infraHPortS := models.NewAccessPortSelector(rn, parentDn, description, infraHPortSattr)

	infraHPortS.Status = "modified"
	err := sm.Save(infraHPortS)
	return infraHPortS, err

}

func (sm *ServiceManager) ListAccessPortSelector(leaf_interface_profile string) ([]*models.AccessPortSelector, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/accportprof-%s/infraHPortS.json", baseurlStr, leaf_interface_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.AccessPortSelectorListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsAccBaseGrpFromAccessPortSelector(parentDn, tnInfraAccBaseGrpName string) error {
	dn := fmt.Sprintf("%s/rsaccBaseGrp", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsAccBaseGrp", dn, tnInfraAccBaseGrpName))

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

func (sm *ServiceManager) DeleteRelationinfraRsAccBaseGrpFromAccessPortSelector(parentDn string) error {
	dn := fmt.Sprintf("%s/rsaccBaseGrp", parentDn)
	return sm.DeleteByDn(dn, "infraRsAccBaseGrp")
}

func (sm *ServiceManager) ReadRelationinfraRsAccBaseGrpFromAccessPortSelector(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsAccBaseGrp")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsAccBaseGrp")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}

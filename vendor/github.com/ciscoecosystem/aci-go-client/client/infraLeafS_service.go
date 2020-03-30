package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSwitchAssociation(switch_association_type string, name string, leaf_profile string, description string, infraLeafSattr models.SwitchAssociationAttributes) (*models.SwitchAssociation, error) {
	rn := fmt.Sprintf("leaves-%s-typ-%s", name, switch_association_type)
	parentDn := fmt.Sprintf("uni/infra/nprof-%s", leaf_profile)
	infraLeafS := models.NewSwitchAssociation(rn, parentDn, description, infraLeafSattr)
	err := sm.Save(infraLeafS)
	return infraLeafS, err
}

func (sm *ServiceManager) ReadSwitchAssociation(switch_association_type string, name string, leaf_profile string) (*models.SwitchAssociation, error) {
	dn := fmt.Sprintf("uni/infra/nprof-%s/leaves-%s-typ-%s", leaf_profile, name, switch_association_type)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraLeafS := models.SwitchAssociationFromContainer(cont)
	return infraLeafS, nil
}

func (sm *ServiceManager) DeleteSwitchAssociation(switch_association_type string, name string, leaf_profile string) error {
	dn := fmt.Sprintf("uni/infra/nprof-%s/leaves-%s-typ-%s", leaf_profile, name, switch_association_type)
	return sm.DeleteByDn(dn, models.InfraleafsClassName)
}

func (sm *ServiceManager) UpdateSwitchAssociation(switch_association_type string, name string, leaf_profile string, description string, infraLeafSattr models.SwitchAssociationAttributes) (*models.SwitchAssociation, error) {
	rn := fmt.Sprintf("leaves-%s-typ-%s", name, switch_association_type)
	parentDn := fmt.Sprintf("uni/infra/nprof-%s", leaf_profile)
	infraLeafS := models.NewSwitchAssociation(rn, parentDn, description, infraLeafSattr)

	infraLeafS.Status = "modified"
	err := sm.Save(infraLeafS)
	return infraLeafS, err

}

func (sm *ServiceManager) ListSwitchAssociation(leaf_profile string) ([]*models.SwitchAssociation, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/nprof-%s/infraLeafS.json", baseurlStr, leaf_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.SwitchAssociationListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsAccNodePGrpFromSwitchAssociation(parentDn, tnInfraAccNodePGrpName string) error {
	dn := fmt.Sprintf("%s/rsaccNodePGrp", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnInfraAccNodePGrpName": "%s"
								
			}
		}
	}`, "infraRsAccNodePGrp", dn, tnInfraAccNodePGrpName))

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

func (sm *ServiceManager) DeleteRelationinfraRsAccNodePGrpFromSwitchAssociation(parentDn string) error {
	dn := fmt.Sprintf("%s/rsaccNodePGrp", parentDn)
	return sm.DeleteByDn(dn, "infraRsAccNodePGrp")
}

func (sm *ServiceManager) ReadRelationinfraRsAccNodePGrpFromSwitchAssociation(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsAccNodePGrp")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsAccNodePGrp")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnInfraAccNodePGrpName")
		return dat, err
	} else {
		return nil, err
	}

}

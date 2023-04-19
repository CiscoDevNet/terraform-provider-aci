package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateL3extVirtualLIfPLagPolicy(floating_path_attributes_tDn string, logical_interface_profile_encap string, logical_interface_profile_nodeDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, l3extVirtualLIfPLagPolAttAttr models.L3extVirtualLIfPLagPolicyAttributes) (*models.L3extVirtualLIfPLagPolicy, error) {
	rn := fmt.Sprintf(models.Rnl3extVirtualLIfPLagPolAtt)
	parentDn := fmt.Sprintf(models.ParentDnl3extVirtualLIfPLagPolAtt, tenant, l3_outside, logical_node_profile, logical_interface_profile, logical_interface_profile_nodeDn, logical_interface_profile_encap, floating_path_attributes_tDn)
	l3extVirtualLIfPLagPolAtt := models.NewL3extVirtualLIfPLagPolicy(rn, parentDn, description, l3extVirtualLIfPLagPolAttAttr)
	err := sm.Save(l3extVirtualLIfPLagPolAtt)
	return l3extVirtualLIfPLagPolAtt, err
}

func (sm *ServiceManager) ReadL3extVirtualLIfPLagPolicy(floating_path_attributes_tDn string, logical_interface_profile_encap string, logical_interface_profile_nodeDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.L3extVirtualLIfPLagPolicy, error) {
	dn := fmt.Sprintf(models.Dnl3extVirtualLIfPLagPolAtt, tenant, l3_outside, logical_node_profile, logical_interface_profile, logical_interface_profile_nodeDn, logical_interface_profile_encap, floating_path_attributes_tDn)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extVirtualLIfPLagPolAtt := models.L3extVirtualLIfPLagPolicyFromContainer(cont)
	return l3extVirtualLIfPLagPolAtt, nil
}

func (sm *ServiceManager) DeleteL3extVirtualLIfPLagPolicy(floating_path_attributes_tDn string, logical_interface_profile_encap string, logical_interface_profile_nodeDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf(models.Dnl3extVirtualLIfPLagPolAtt, tenant, l3_outside, logical_node_profile, logical_interface_profile, logical_interface_profile_nodeDn, logical_interface_profile_encap, floating_path_attributes_tDn)
	return sm.DeleteByDn(dn, models.L3extvirtuallifplagpolattClassName)
}

func (sm *ServiceManager) UpdateL3extVirtualLIfPLagPolicy(floating_path_attributes_tDn string, logical_interface_profile_encap string, logical_interface_profile_nodeDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, l3extVirtualLIfPLagPolAttAttr models.L3extVirtualLIfPLagPolicyAttributes) (*models.L3extVirtualLIfPLagPolicy, error) {
	rn := fmt.Sprintf(models.Rnl3extVirtualLIfPLagPolAtt)
	parentDn := fmt.Sprintf(models.ParentDnl3extVirtualLIfPLagPolAtt, tenant, l3_outside, logical_node_profile, logical_interface_profile, logical_interface_profile_nodeDn, logical_interface_profile_encap, floating_path_attributes_tDn)
	l3extVirtualLIfPLagPolAtt := models.NewL3extVirtualLIfPLagPolicy(rn, parentDn, description, l3extVirtualLIfPLagPolAttAttr)
	l3extVirtualLIfPLagPolAtt.Status = "modified"
	err := sm.Save(l3extVirtualLIfPLagPolAtt)
	return l3extVirtualLIfPLagPolAtt, err
}

func (sm *ServiceManager) ListL3extVirtualLIfPLagPolicy(floating_path_attributes_tDn string, logical_interface_profile_encap string, logical_interface_profile_nodeDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.L3extVirtualLIfPLagPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/lifp-%s/vlifp-[%s]-[%s]/rsdynPathAtt-[%s]/l3extVirtualLIfPLagPolAtt.json", models.BaseurlStr, tenant, l3_outside, logical_node_profile, logical_interface_profile, logical_interface_profile_nodeDn, logical_interface_profile_encap, floating_path_attributes_tDn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.L3extVirtualLIfPLagPolicyListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationl3extRsVSwitchEnhancedLagPol(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsvSwitchEnhancedLagPol-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
	    "l3extVirtualLIfPLagPolAtt": {
            "attributes": {
                "annotation": "orchestrator:terraform",
                "dn": "%s"
            },
            "children": [
                {
                    "l3extRsVSwitchEnhancedLagPol": {
                        "attributes": {
                            "dn": "%s",
                            "annotation": "orchestrator:terraform",
                            "tDn": "%s"
                        }
                    }
                }
            ]
        }
	}`, parentDn, dn, tDn))

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

func (sm *ServiceManager) DeleteRelationl3extRsVSwitchEnhancedLagPol(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsvSwitchEnhancedLagPol-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "l3extRsVSwitchEnhancedLagPol")
}

func (sm *ServiceManager) ReadRelationl3extRsVSwitchEnhancedLagPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "l3extRsVSwitchEnhancedLagPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "l3extRsVSwitchEnhancedLagPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateVirtualLogicalInterfaceProfile(encap string, nodeDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, l3extVirtualLIfPattr models.VirtualLogicalInterfaceProfileAttributes) (*models.VirtualLogicalInterfaceProfile, error) {
	rn := fmt.Sprintf("vlifp-[%s]-[%s]", nodeDn, encap)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	l3extVirtualLIfP := models.NewVirtualLogicalInterfaceProfile(rn, parentDn, description, l3extVirtualLIfPattr)
	err := sm.Save(l3extVirtualLIfP)
	return l3extVirtualLIfP, err
}

func (sm *ServiceManager) ReadVirtualLogicalInterfaceProfile(encap string, nodeDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) (*models.VirtualLogicalInterfaceProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/vlifp-[%s]-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, nodeDn, encap)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extVirtualLIfP := models.VirtualLogicalInterfaceProfileFromContainer(cont)
	return l3extVirtualLIfP, nil
}

func (sm *ServiceManager) DeleteVirtualLogicalInterfaceProfile(encap string, nodeDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s/vlifp-[%s]-[%s]", tenant, l3_outside, logical_node_profile, logical_interface_profile, nodeDn, encap)
	return sm.DeleteByDn(dn, models.L3extvirtuallifpClassName)
}

func (sm *ServiceManager) UpdateVirtualLogicalInterfaceProfile(encap string, nodeDn string, logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string, description string, l3extVirtualLIfPattr models.VirtualLogicalInterfaceProfileAttributes) (*models.VirtualLogicalInterfaceProfile, error) {
	rn := fmt.Sprintf("vlifp-[%s]-[%s]", nodeDn, encap)
	parentDn := fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", tenant, l3_outside, logical_node_profile, logical_interface_profile)
	l3extVirtualLIfP := models.NewVirtualLogicalInterfaceProfile(rn, parentDn, description, l3extVirtualLIfPattr)

	l3extVirtualLIfP.Status = "modified"
	err := sm.Save(l3extVirtualLIfP)
	return l3extVirtualLIfP, err

}

func (sm *ServiceManager) ListVirtualLogicalInterfaceProfile(logical_interface_profile string, logical_node_profile string, l3_outside string, tenant string) ([]*models.VirtualLogicalInterfaceProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/out-%s/lnodep-%s/lifp-%s/l3extVirtualLIfP.json", baseurlStr, tenant, l3_outside, logical_node_profile, logical_interface_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.VirtualLogicalInterfaceProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationl3extRsDynPathAttFromLogicalInterfaceProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsdynPathAtt-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s", "annotation":"orchestrator:terraform"				
			}
		}
	}`, "l3extRsDynPathAtt", dn))

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

func (sm *ServiceManager) DeleteRelationl3extRsDynPathAttFromLogicalInterfaceProfile(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsdynPathAtt-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "l3extRsDynPathAtt")
}

func (sm *ServiceManager) ReadRelationl3extRsDynPathAttFromLogicalInterfaceProfile(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l3extRsDynPathAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l3extRsDynPathAtt")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}

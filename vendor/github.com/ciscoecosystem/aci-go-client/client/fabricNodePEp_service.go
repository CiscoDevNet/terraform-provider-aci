package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateNodePolicyEndPoint(node_policy_end_point_id string, vpc_explicit_protection_group string, description string, fabricNodePEpattr models.NodePolicyEndPointAttributes) (*models.NodePolicyEndPoint, error) {
	rn := fmt.Sprintf("nodepep-%s", node_policy_end_point_id)
	parentDn := fmt.Sprintf("uni/fabric/protpol/expgep-%s", vpc_explicit_protection_group)
	fabricNodePEp := models.NewNodePolicyEndPoint(rn, parentDn, description, fabricNodePEpattr)
	err := sm.Save(fabricNodePEp)
	return fabricNodePEp, err
}

func (sm *ServiceManager) ReadNodePolicyEndPoint(node_policy_end_point_id string, vpc_explicit_protection_group string) (*models.NodePolicyEndPoint, error) {
	dn := fmt.Sprintf("uni/fabric/protpol/expgep-%s/nodepep-%s", vpc_explicit_protection_group, node_policy_end_point_id)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricNodePEp := models.NodePolicyEndPointFromContainer(cont)
	return fabricNodePEp, nil
}

func (sm *ServiceManager) DeleteNodePolicyEndPoint(node_policy_end_point_id string, vpc_explicit_protection_group string) error {
	dn := fmt.Sprintf("uni/fabric/protpol/expgep-%s/nodepep-%s", vpc_explicit_protection_group, node_policy_end_point_id)
	return sm.DeleteByDn(dn, models.FabricnodepepClassName)
}

func (sm *ServiceManager) UpdateNodePolicyEndPoint(node_policy_end_point_id string, vpc_explicit_protection_group string, description string, fabricNodePEpattr models.NodePolicyEndPointAttributes) (*models.NodePolicyEndPoint, error) {
	rn := fmt.Sprintf("nodepep-%s", node_policy_end_point_id)
	parentDn := fmt.Sprintf("uni/fabric/protpol/expgep-%s", vpc_explicit_protection_group)
	fabricNodePEp := models.NewNodePolicyEndPoint(rn, parentDn, description, fabricNodePEpattr)

	fabricNodePEp.Status = "modified"
	err := sm.Save(fabricNodePEp)
	return fabricNodePEp, err

}

func (sm *ServiceManager) ListNodePolicyEndPoint(vpc_explicit_protection_group string) ([]*models.NodePolicyEndPoint, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/fabric/protpol/expgep-%s/fabricNodePEp.json", baseurlStr, vpc_explicit_protection_group)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.NodePolicyEndPointListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationfabricRsToPeerNodeCfgFromNodePolicyEndPoint(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rstoPeerNodeCfg-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s"				
			}
		}
	}`, "fabricRsToPeerNodeCfg", dn))

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

func (sm *ServiceManager) ReadRelationfabricRsToPeerNodeCfgFromNodePolicyEndPoint(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fabricRsToPeerNodeCfg")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fabricRsToPeerNodeCfg")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}

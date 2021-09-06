package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFabricNodeOrg(fabric_node_id string, fabric_pod_id string, description string, fabricNodeattr models.OrgFabricNodeAttributes) (*models.OrgFabricNode, error) {
	rn := fmt.Sprintf("node-%s", fabric_node_id)
	parentDn := fmt.Sprintf("uni/topology/pod-%s", fabric_pod_id)
	fabricNode := models.OrgNewFabricNode(rn, parentDn, description, fabricNodeattr)
	err := sm.Save(fabricNode)
	return fabricNode, err
}

func (sm *ServiceManager) ReadFabricNodeOrg(fabric_node_id string, fabric_pod_id string) (*models.OrgFabricNode, error) {
	dn := fmt.Sprintf("uni/topology/pod-%s/node-%s", fabric_pod_id, fabric_node_id)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricNode := models.OrgFabricNodeFromContainer(cont)
	return fabricNode, nil
}

func (sm *ServiceManager) DeleteFabricNodeOrg(fabric_node_id string, fabric_pod_id string) error {
	dn := fmt.Sprintf("uni/topology/pod-%s/node-%s", fabric_pod_id, fabric_node_id)
	return sm.DeleteByDn(dn, models.FabricnodeClassName)
}

func (sm *ServiceManager) UpdateFabricNodeOrg(fabric_node_id string, fabric_pod_id string, description string, fabricNodeattr models.OrgFabricNodeAttributes) (*models.OrgFabricNode, error) {
	rn := fmt.Sprintf("node-%s", fabric_node_id)
	parentDn := fmt.Sprintf("uni/topology/pod-%s", fabric_pod_id)
	fabricNode := models.OrgNewFabricNode(rn, parentDn, description, fabricNodeattr)

	fabricNode.Status = "modified"
	err := sm.Save(fabricNode)
	return fabricNode, err

}

func (sm *ServiceManager) ListFabricNodeOrg(fabric_pod_id string) ([]*models.OrgFabricNode, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/topology/pod-%s/fabricNode.json", baseurlStr, fabric_pod_id)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.OrgFabricNodeListFromContainer(cont)

	return list, err
}

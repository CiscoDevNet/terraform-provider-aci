package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFabricPathEndpoint(name string, fabric_path_end_point_container_nodeId string, fabric_pod_id string, description string, fabricPathEpattr models.FabricPathEndpointAttributes) (*models.FabricPathEndpoint, error) {
	rn := fmt.Sprintf("pathep-[%s]", name)
	parentDn := fmt.Sprintf("uni/topology/pod-%s/paths-%s", fabric_pod_id, fabric_path_end_point_container_nodeId)
	fabricPathEp := models.NewFabricPathEndpoint(rn, parentDn, description, fabricPathEpattr)
	err := sm.Save(fabricPathEp)
	return fabricPathEp, err
}

func (sm *ServiceManager) ReadFabricPathEndpoint(name string, fabric_path_end_point_container_nodeId string, fabric_pod_id string) (*models.FabricPathEndpoint, error) {
	dn := fmt.Sprintf("uni/topology/pod-%s/paths-%s/pathep-[%s]", fabric_pod_id, fabric_path_end_point_container_nodeId, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricPathEp := models.FabricPathEndpointFromContainer(cont)
	return fabricPathEp, nil
}

func (sm *ServiceManager) DeleteFabricPathEndpoint(name string, fabric_path_end_point_container_nodeId string, fabric_pod_id string) error {
	dn := fmt.Sprintf("uni/topology/pod-%s/paths-%s/pathep-[%s]", fabric_pod_id, fabric_path_end_point_container_nodeId, name)
	return sm.DeleteByDn(dn, models.FabricpathepClassName)
}

func (sm *ServiceManager) UpdateFabricPathEndpoint(name string, fabric_path_end_point_container_nodeId string, fabric_pod_id string, description string, fabricPathEpattr models.FabricPathEndpointAttributes) (*models.FabricPathEndpoint, error) {
	rn := fmt.Sprintf("pathep-[%s]", name)
	parentDn := fmt.Sprintf("uni/topology/pod-%s/paths-%s", fabric_pod_id, fabric_path_end_point_container_nodeId)
	fabricPathEp := models.NewFabricPathEndpoint(rn, parentDn, description, fabricPathEpattr)

	fabricPathEp.Status = "modified"
	err := sm.Save(fabricPathEp)
	return fabricPathEp, err

}

func (sm *ServiceManager) ListFabricPathEndpoint(fabric_path_end_point_container_nodeId string, fabric_pod_id string) ([]*models.FabricPathEndpoint, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/topology/pod-%s/paths-%s/fabricPathEp.json", baseurlStr, fabric_pod_id, fabric_path_end_point_container_nodeId)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FabricPathEndpointListFromContainer(cont)

	return list, err
}

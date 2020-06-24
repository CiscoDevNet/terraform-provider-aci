package client

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) ReadTopologyFabricNode(pod int, node int) (*models.TopologyFabricNode, error) {
	dn := fmt.Sprintf("topology/pod-%d/node-%d", pod, node)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	fabricNode := models.TopologyFabricNodeFromContainer(cont)
	return fabricNode, nil
}

func (sm *ServiceManager) ListTopologyFabricNode() ([]*models.TopologyFabricNode, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/fabricNode.json", baseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.TopologyFabricNodeListFromContainer(cont)
	return list, err
}

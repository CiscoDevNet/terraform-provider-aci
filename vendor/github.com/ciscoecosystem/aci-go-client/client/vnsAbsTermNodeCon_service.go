package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateConsumerTerminalNode(name string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsTermNodeConattr models.ConsumerTerminalNodeAttributes) (*models.ConsumerTerminalNode, error) {
	rn := fmt.Sprintf("AbsTermNodeCon-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s", tenant, l4_l7_service_graph_template)
	vnsAbsTermNodeCon := models.NewConsumerTerminalNode(rn, parentDn, description, vnsAbsTermNodeConattr)
	err := sm.Save(vnsAbsTermNodeCon)
	return vnsAbsTermNodeCon, err
}

func (sm *ServiceManager) ReadConsumerTerminalNode(name string, l4_l7_service_graph_template string, tenant string) (*models.ConsumerTerminalNode, error) {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsTermNodeCon-%s", tenant, l4_l7_service_graph_template, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsTermNodeCon := models.ConsumerTerminalNodeFromContainer(cont)
	return vnsAbsTermNodeCon, nil
}

func (sm *ServiceManager) DeleteConsumerTerminalNode(name string, l4_l7_service_graph_template string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsTermNodeCon-%s", tenant, l4_l7_service_graph_template, name)
	return sm.DeleteByDn(dn, models.VnsabstermnodeconClassName)
}

func (sm *ServiceManager) UpdateConsumerTerminalNode(name string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsTermNodeConattr models.ConsumerTerminalNodeAttributes) (*models.ConsumerTerminalNode, error) {
	rn := fmt.Sprintf("AbsTermNodeCon-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s", tenant, l4_l7_service_graph_template)
	vnsAbsTermNodeCon := models.NewConsumerTerminalNode(rn, parentDn, description, vnsAbsTermNodeConattr)

	vnsAbsTermNodeCon.Status = "modified"
	err := sm.Save(vnsAbsTermNodeCon)
	return vnsAbsTermNodeCon, err

}

func (sm *ServiceManager) ListConsumerTerminalNode(l4_l7_service_graph_template string, tenant string) ([]*models.ConsumerTerminalNode, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/AbsGraph-%s/vnsAbsTermNodeCon.json", baseurlStr, tenant, l4_l7_service_graph_template)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ConsumerTerminalNodeListFromContainer(cont)

	return list, err
}

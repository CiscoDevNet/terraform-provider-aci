package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateProviderTerminalNode(name string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsTermNodeProvattr models.ProviderTerminalNodeAttributes) (*models.ProviderTerminalNode, error) {
	rn := fmt.Sprintf("AbsTermNodeProv-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s", tenant, l4_l7_service_graph_template)
	vnsAbsTermNodeProv := models.NewProviderTerminalNode(rn, parentDn, description, vnsAbsTermNodeProvattr)
	err := sm.Save(vnsAbsTermNodeProv)
	return vnsAbsTermNodeProv, err
}

func (sm *ServiceManager) ReadProviderTerminalNode(name string, l4_l7_service_graph_template string, tenant string) (*models.ProviderTerminalNode, error) {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsTermNodeProv-%s", tenant, l4_l7_service_graph_template, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsTermNodeProv := models.ProviderTerminalNodeFromContainer(cont)
	return vnsAbsTermNodeProv, nil
}

func (sm *ServiceManager) DeleteProviderTerminalNode(name string, l4_l7_service_graph_template string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsTermNodeProv-%s", tenant, l4_l7_service_graph_template, name)
	return sm.DeleteByDn(dn, models.VnsabstermnodeprovClassName)
}

func (sm *ServiceManager) UpdateProviderTerminalNode(name string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsTermNodeProvattr models.ProviderTerminalNodeAttributes) (*models.ProviderTerminalNode, error) {
	rn := fmt.Sprintf("AbsTermNodeProv-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s", tenant, l4_l7_service_graph_template)
	vnsAbsTermNodeProv := models.NewProviderTerminalNode(rn, parentDn, description, vnsAbsTermNodeProvattr)

	vnsAbsTermNodeProv.Status = "modified"
	err := sm.Save(vnsAbsTermNodeProv)
	return vnsAbsTermNodeProv, err

}

func (sm *ServiceManager) ListProviderTerminalNode(l4_l7_service_graph_template string, tenant string) ([]*models.ProviderTerminalNode, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/AbsGraph-%s/vnsAbsTermNodeProv.json", baseurlStr, tenant, l4_l7_service_graph_template)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ProviderTerminalNodeListFromContainer(cont)

	return list, err
}

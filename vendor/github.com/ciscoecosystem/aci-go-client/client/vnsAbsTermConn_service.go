package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateTerminalConnector(consumer_terminal_node string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsTermConnattr models.TerminalConnectorAttributes) (*models.TerminalConnector, error) {
	rn := fmt.Sprintf("AbsTConn")
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsTermNodeCon-%s", tenant, l4_l7_service_graph_template, consumer_terminal_node)
	vnsAbsTermConn := models.NewTerminalConnector(rn, parentDn, description, vnsAbsTermConnattr)
	err := sm.Save(vnsAbsTermConn)
	return vnsAbsTermConn, err
}

func (sm *ServiceManager) ReadTerminalConnector(consumer_terminal_node string, l4_l7_service_graph_template string, tenant string) (*models.TerminalConnector, error) {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsTermNodeCon-%s/AbsTConn", tenant, l4_l7_service_graph_template, consumer_terminal_node)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsTermConn := models.TerminalConnectorFromContainer(cont)
	return vnsAbsTermConn, nil
}

func (sm *ServiceManager) DeleteTerminalConnector(consumer_terminal_node string, l4_l7_service_graph_template string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsTermNodeCon-%s/AbsTConn", tenant, l4_l7_service_graph_template, consumer_terminal_node)
	return sm.DeleteByDn(dn, models.VnsabstermconnClassName)
}

func (sm *ServiceManager) UpdateTerminalConnector(consumer_terminal_node string, l4_l7_service_graph_template string, tenant string, description string, vnsAbsTermConnattr models.TerminalConnectorAttributes) (*models.TerminalConnector, error) {
	rn := fmt.Sprintf("AbsTConn")
	parentDn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s/AbsTermNodeCon-%s", tenant, l4_l7_service_graph_template, consumer_terminal_node)
	vnsAbsTermConn := models.NewTerminalConnector(rn, parentDn, description, vnsAbsTermConnattr)

	vnsAbsTermConn.Status = "modified"
	err := sm.Save(vnsAbsTermConn)
	return vnsAbsTermConn, err

}

func (sm *ServiceManager) ListTerminalConnector(consumer_terminal_node string, l4_l7_service_graph_template string, tenant string) ([]*models.TerminalConnector, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/AbsGraph-%s/AbsTermNodeCon-%s/vnsAbsTermConn.json", baseurlStr, tenant, l4_l7_service_graph_template, consumer_terminal_node)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.TerminalConnectorListFromContainer(cont)

	return list, err
}

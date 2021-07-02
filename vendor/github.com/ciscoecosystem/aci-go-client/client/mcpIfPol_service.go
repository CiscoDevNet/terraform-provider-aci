package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateMiscablingProtocolInterfacePolicy(name string, description string, mcpIfPolattr models.MiscablingProtocolInterfacePolicyAttributes) (*models.MiscablingProtocolInterfacePolicy, error) {
	rn := fmt.Sprintf("infra/mcpIfP-%s", name)
	parentDn := fmt.Sprintf("uni")
	mcpIfPol := models.NewMiscablingProtocolInterfacePolicy(rn, parentDn, description, mcpIfPolattr)
	err := sm.Save(mcpIfPol)
	return mcpIfPol, err
}

func (sm *ServiceManager) ReadMiscablingProtocolInterfacePolicy(name string) (*models.MiscablingProtocolInterfacePolicy, error) {
	dn := fmt.Sprintf("uni/infra/mcpIfP-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	mcpIfPol := models.MiscablingProtocolInterfacePolicyFromContainer(cont)
	return mcpIfPol, nil
}

func (sm *ServiceManager) DeleteMiscablingProtocolInterfacePolicy(name string) error {
	dn := fmt.Sprintf("uni/infra/mcpIfP-%s", name)
	return sm.DeleteByDn(dn, models.McpifpolClassName)
}

func (sm *ServiceManager) UpdateMiscablingProtocolInterfacePolicy(name string, description string, mcpIfPolattr models.MiscablingProtocolInterfacePolicyAttributes) (*models.MiscablingProtocolInterfacePolicy, error) {
	rn := fmt.Sprintf("infra/mcpIfP-%s", name)
	parentDn := fmt.Sprintf("uni")
	mcpIfPol := models.NewMiscablingProtocolInterfacePolicy(rn, parentDn, description, mcpIfPolattr)

	mcpIfPol.Status = "modified"
	err := sm.Save(mcpIfPol)
	return mcpIfPol, err

}

func (sm *ServiceManager) ListMiscablingProtocolInterfacePolicy() ([]*models.MiscablingProtocolInterfacePolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/mcpIfPol.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.MiscablingProtocolInterfacePolicyListFromContainer(cont)

	return list, err
}

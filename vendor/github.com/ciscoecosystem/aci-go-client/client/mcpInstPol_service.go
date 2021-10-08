package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateMiscablingProtocolInstancePolicy(name string, description string, nameAlias string, mcpInstPolAttr models.MiscablingProtocolInstancePolicyAttributes) (*models.MiscablingProtocolInstancePolicy, error) {
	rn := fmt.Sprintf(models.RnmcpInstPol, name)
	parentDn := fmt.Sprintf(models.ParentDnmcpInstPol)
	mcpInstPol := models.NewMiscablingProtocolInstancePolicy(rn, parentDn, description, nameAlias, mcpInstPolAttr)
	err := sm.Save(mcpInstPol)
	return mcpInstPol, err
}

func (sm *ServiceManager) ReadMiscablingProtocolInstancePolicy(name string) (*models.MiscablingProtocolInstancePolicy, error) {
	dn := fmt.Sprintf(models.DnmcpInstPol, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	mcpInstPol := models.MiscablingProtocolInstancePolicyFromContainer(cont)
	return mcpInstPol, nil
}

func (sm *ServiceManager) DeleteMiscablingProtocolInstancePolicy(name string) error {
	dn := fmt.Sprintf(models.DnmcpInstPol, name)
	return sm.DeleteByDn(dn, models.McpinstpolClassName)
}

func (sm *ServiceManager) UpdateMiscablingProtocolInstancePolicy(name string, description string, nameAlias string, mcpInstPolAttr models.MiscablingProtocolInstancePolicyAttributes) (*models.MiscablingProtocolInstancePolicy, error) {
	rn := fmt.Sprintf(models.RnmcpInstPol, name)
	parentDn := fmt.Sprintf(models.ParentDnmcpInstPol)
	mcpInstPol := models.NewMiscablingProtocolInstancePolicy(rn, parentDn, description, nameAlias, mcpInstPolAttr)
	mcpInstPol.Status = "modified"
	err := sm.Save(mcpInstPol)
	return mcpInstPol, err
}

func (sm *ServiceManager) ListMiscablingProtocolInstancePolicy() ([]*models.MiscablingProtocolInstancePolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/mcpInstPol.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MiscablingProtocolInstancePolicyListFromContainer(cont)
	return list, err
}

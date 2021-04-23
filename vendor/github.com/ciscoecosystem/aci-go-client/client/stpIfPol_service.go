package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSpanningTreeInterfacePolicy(name string, description string, nameAlias string, stpIfPolAttr models.SpanningTreeInterfacePolicyAttributes) (*models.SpanningTreeInterfacePolicy, error) {
	rn := fmt.Sprintf(models.RnstpIfPol, name)
	parentDn := fmt.Sprintf(models.ParentDnstpIfPol)
	stpIfPol := models.NewSpanningTreeInterfacePolicy(rn, parentDn, description, nameAlias, stpIfPolAttr)
	err := sm.Save(stpIfPol)
	return stpIfPol, err
}

func (sm *ServiceManager) ReadSpanningTreeInterfacePolicy(name string) (*models.SpanningTreeInterfacePolicy, error) {
	dn := fmt.Sprintf(models.DnstpIfPol, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	stpIfPol := models.SpanningTreeInterfacePolicyFromContainer(cont)
	return stpIfPol, nil
}

func (sm *ServiceManager) DeleteSpanningTreeInterfacePolicy(name string) error {
	dn := fmt.Sprintf(models.DnstpIfPol, name)
	return sm.DeleteByDn(dn, models.StpifpolClassName)
}

func (sm *ServiceManager) UpdateSpanningTreeInterfacePolicy(name string, description string, nameAlias string, stpIfPolAttr models.SpanningTreeInterfacePolicyAttributes) (*models.SpanningTreeInterfacePolicy, error) {
	rn := fmt.Sprintf(models.RnstpIfPol, name)
	parentDn := fmt.Sprintf(models.ParentDnstpIfPol)
	stpIfPol := models.NewSpanningTreeInterfacePolicy(rn, parentDn, description, nameAlias, stpIfPolAttr)
	stpIfPol.Status = "modified"
	err := sm.Save(stpIfPol)
	return stpIfPol, err
}

func (sm *ServiceManager) ListSpanningTreeInterfacePolicy() ([]*models.SpanningTreeInterfacePolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/stpIfPol.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SpanningTreeInterfacePolicyListFromContainer(cont)
	return list, err
}

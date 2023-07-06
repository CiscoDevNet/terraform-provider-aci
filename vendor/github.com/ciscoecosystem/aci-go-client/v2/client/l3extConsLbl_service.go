package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateL3ExtConsLbl(name string, l3_outside string, tenant string, description string, consLblAttr models.L3ExtConsLblAttributes) (*models.L3ExtConsLbl, error) {
	rn := fmt.Sprintf(models.RnL3extConsLbl, name)
	parentDn := fmt.Sprintf(models.ParentDnL3ExtConsLbl, tenant, l3_outside)
	consLbl := models.NewL3ExtConsLbl(rn, parentDn, description, consLblAttr)
	err := sm.Save(consLbl)
	return consLbl, err
}

func (sm *ServiceManager) ReadL3ExtConsLbl(name string, l3_outside string, tenant string) (*models.L3ExtConsLbl, error) {
	dn := fmt.Sprintf(models.DnL3extConsLbl, tenant, l3_outside, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	consLbl := models.L3ExtConsLblFromContainer(cont)
	return consLbl, nil
}

func (sm *ServiceManager) DeleteL3ExtConsLbl(name string, l3_outside string, tenant string) error {
	dn := fmt.Sprintf(models.DnL3extConsLbl, tenant, l3_outside, name)
	return sm.DeleteByDn(dn, models.L3extConsLblClassName)
}

func (sm *ServiceManager) UpdateL3ExtConsLbl(name string, l3_outside string, tenant string, description string, consLblAttr models.L3ExtConsLblAttributes) (*models.L3ExtConsLbl, error) {
	rn := fmt.Sprintf(models.RnL3extConsLbl, name)
	parentDn := fmt.Sprintf(models.ParentDnL3ExtConsLbl, tenant, l3_outside)
	consLbl := models.NewL3ExtConsLbl(rn, parentDn, description, consLblAttr)

	consLbl.Status = "modified"
	err := sm.Save(consLbl)
	return consLbl, err
}

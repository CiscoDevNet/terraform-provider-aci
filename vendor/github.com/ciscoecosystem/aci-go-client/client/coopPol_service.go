package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateCOOPGroupPolicy(name string, description string, nameAlias string, coopPolAttr models.COOPGroupPolicyAttributes) (*models.COOPGroupPolicy, error) {
	rn := fmt.Sprintf(models.RncoopPol, name)
	parentDn := fmt.Sprintf(models.ParentDncoopPol)
	coopPol := models.NewCOOPGroupPolicy(rn, parentDn, description, nameAlias, coopPolAttr)
	err := sm.Save(coopPol)
	return coopPol, err
}

func (sm *ServiceManager) ReadCOOPGroupPolicy(name string) (*models.COOPGroupPolicy, error) {
	dn := fmt.Sprintf(models.DncoopPol, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	coopPol := models.COOPGroupPolicyFromContainer(cont)
	return coopPol, nil
}

func (sm *ServiceManager) DeleteCOOPGroupPolicy(name string) error {
	dn := fmt.Sprintf(models.DncoopPol, name)
	return sm.DeleteByDn(dn, models.CooppolClassName)
}

func (sm *ServiceManager) UpdateCOOPGroupPolicy(name string, description string, nameAlias string, coopPolAttr models.COOPGroupPolicyAttributes) (*models.COOPGroupPolicy, error) {
	rn := fmt.Sprintf(models.RncoopPol, name)
	parentDn := fmt.Sprintf(models.ParentDncoopPol)
	coopPol := models.NewCOOPGroupPolicy(rn, parentDn, description, nameAlias, coopPolAttr)
	coopPol.Status = "modified"
	err := sm.Save(coopPol)
	return coopPol, err
}

func (sm *ServiceManager) ListCOOPGroupPolicy() ([]*models.COOPGroupPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/coopPol.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.COOPGroupPolicyListFromContainer(cont)
	return list, err
}

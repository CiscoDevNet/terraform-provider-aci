package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateErrorDisabledRecoveryPolicy(name string, description string, nameAlias string, edrErrDisRecoverPolAttr models.ErrorDisabledRecoveryPolicyAttributes) (*models.ErrorDisabledRecoveryPolicy, error) {
	rn := fmt.Sprintf(models.RnedrErrDisRecoverPol, name)
	parentDn := fmt.Sprintf(models.ParentDnedrErrDisRecoverPol)
	edrErrDisRecoverPol := models.NewErrorDisabledRecoveryPolicy(rn, parentDn, description, nameAlias, edrErrDisRecoverPolAttr)
	err := sm.Save(edrErrDisRecoverPol)
	return edrErrDisRecoverPol, err
}

func (sm *ServiceManager) ReadErrorDisabledRecoveryPolicy(name string) (*models.ErrorDisabledRecoveryPolicy, error) {
	dn := fmt.Sprintf(models.DnedrErrDisRecoverPol, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	edrErrDisRecoverPol := models.ErrorDisabledRecoveryPolicyFromContainer(cont)
	return edrErrDisRecoverPol, nil
}

func (sm *ServiceManager) DeleteErrorDisabledRecoveryPolicy(name string) error {
	dn := fmt.Sprintf(models.DnedrErrDisRecoverPol, name)
	return sm.DeleteByDn(dn, models.EdrerrdisrecoverpolClassName)
}

func (sm *ServiceManager) UpdateErrorDisabledRecoveryPolicy(name string, description string, nameAlias string, edrErrDisRecoverPolAttr models.ErrorDisabledRecoveryPolicyAttributes) (*models.ErrorDisabledRecoveryPolicy, error) {
	rn := fmt.Sprintf(models.RnedrErrDisRecoverPol, name)
	parentDn := fmt.Sprintf(models.ParentDnedrErrDisRecoverPol)
	edrErrDisRecoverPol := models.NewErrorDisabledRecoveryPolicy(rn, parentDn, description, nameAlias, edrErrDisRecoverPolAttr)
	edrErrDisRecoverPol.Status = "modified"
	err := sm.Save(edrErrDisRecoverPol)
	return edrErrDisRecoverPol, err
}

func (sm *ServiceManager) ListErrorDisabledRecoveryPolicy() ([]*models.ErrorDisabledRecoveryPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/edrErrDisRecoverPol.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.ErrorDisabledRecoveryPolicyListFromContainer(cont)
	return list, err
}

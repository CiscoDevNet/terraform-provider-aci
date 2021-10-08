package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateErrorDisabledRecoveryEvent(event string, error_disabled_recovery_policy string, description string, nameAlias string, edrEventPAttr models.ErrorDisabledRecoveryEventAttributes) (*models.ErrorDisabledRecoveryEvent, error) {
	rn := fmt.Sprintf(models.RnedrEventP, event)
	parentDn := fmt.Sprintf(models.ParentDnedrEventP, error_disabled_recovery_policy)
	edrEventP := models.NewErrorDisabledRecoveryEvent(rn, parentDn, description, nameAlias, edrEventPAttr)
	err := sm.Save(edrEventP)
	return edrEventP, err
}

func (sm *ServiceManager) ReadErrorDisabledRecoveryEvent(event string, error_disabled_recovery_policy string) (*models.ErrorDisabledRecoveryEvent, error) {
	dn := fmt.Sprintf(models.DnedrEventP, error_disabled_recovery_policy, event)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	edrEventP := models.ErrorDisabledRecoveryEventFromContainer(cont)
	return edrEventP, nil
}

func (sm *ServiceManager) DeleteErrorDisabledRecoveryEvent(event string, error_disabled_recovery_policy string) error {
	dn := fmt.Sprintf(models.DnedrEventP, error_disabled_recovery_policy, event)
	return sm.DeleteByDn(dn, models.EdreventpClassName)
}

func (sm *ServiceManager) UpdateErrorDisabledRecoveryEvent(event string, error_disabled_recovery_policy string, description string, nameAlias string, edrEventPAttr models.ErrorDisabledRecoveryEventAttributes) (*models.ErrorDisabledRecoveryEvent, error) {
	rn := fmt.Sprintf(models.RnedrEventP, event)
	parentDn := fmt.Sprintf(models.ParentDnedrEventP, error_disabled_recovery_policy)
	edrEventP := models.NewErrorDisabledRecoveryEvent(rn, parentDn, description, nameAlias, edrEventPAttr)
	edrEventP.Status = "modified"
	err := sm.Save(edrEventP)
	return edrEventP, err
}

func (sm *ServiceManager) ListErrorDisabledRecoveryEvent(error_disabled_recovery_policy string) ([]*models.ErrorDisabledRecoveryEvent, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/edrErrDisRecoverPol-%s/edrEventP.json", models.BaseurlStr, error_disabled_recovery_policy)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.ErrorDisabledRecoveryEventListFromContainer(cont)
	return list, err
}

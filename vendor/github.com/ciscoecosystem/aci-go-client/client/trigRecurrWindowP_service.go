package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRecurringWindow(name string, scheduler string, nameAlias string, trigRecurrWindowPAttr models.RecurringWindowAttributes) (*models.RecurringWindow, error) {
	rn := fmt.Sprintf(models.RntrigRecurrWindowP, name)
	parentDn := fmt.Sprintf(models.ParentDntrigRecurrWindowP, scheduler)
	trigRecurrWindowP := models.NewRecurringWindow(rn, parentDn, nameAlias, trigRecurrWindowPAttr)
	err := sm.Save(trigRecurrWindowP)
	return trigRecurrWindowP, err
}

func (sm *ServiceManager) ReadRecurringWindow(name string, scheduler string) (*models.RecurringWindow, error) {
	dn := fmt.Sprintf(models.DntrigRecurrWindowP, scheduler, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	trigRecurrWindowP := models.RecurringWindowFromContainer(cont)
	return trigRecurrWindowP, nil
}

func (sm *ServiceManager) DeleteRecurringWindow(name string, scheduler string) error {
	dn := fmt.Sprintf(models.DntrigRecurrWindowP, scheduler, name)
	return sm.DeleteByDn(dn, models.TrigrecurrwindowpClassName)
}

func (sm *ServiceManager) UpdateRecurringWindow(name string, scheduler string, nameAlias string, trigRecurrWindowPAttr models.RecurringWindowAttributes) (*models.RecurringWindow, error) {
	rn := fmt.Sprintf(models.RntrigRecurrWindowP, name)
	parentDn := fmt.Sprintf(models.ParentDntrigRecurrWindowP, scheduler)
	trigRecurrWindowP := models.NewRecurringWindow(rn, parentDn, nameAlias, trigRecurrWindowPAttr)
	trigRecurrWindowP.Status = "modified"
	err := sm.Save(trigRecurrWindowP)
	return trigRecurrWindowP, err
}

func (sm *ServiceManager) ListRecurringWindow(scheduler string) ([]*models.RecurringWindow, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/schedp-%s/trigRecurrWindowP.json", models.BaseurlStr, scheduler)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.RecurringWindowListFromContainer(cont)
	return list, err
}

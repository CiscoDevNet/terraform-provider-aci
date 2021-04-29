package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFilterEntry(name string, filter string, tenant string, description string, vzEntryattr models.FilterEntryAttributes) (*models.FilterEntry, error) {
	rn := fmt.Sprintf("e-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/flt-%s", tenant, filter)
	vzEntry := models.NewFilterEntry(rn, parentDn, description, vzEntryattr)
	err := sm.Save(vzEntry)
	return vzEntry, err
}

func (sm *ServiceManager) ReadFilterEntry(name string, filter string, tenant string) (*models.FilterEntry, error) {
	dn := fmt.Sprintf("uni/tn-%s/flt-%s/e-%s", tenant, filter, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vzEntry := models.FilterEntryFromContainer(cont)
	return vzEntry, nil
}

func (sm *ServiceManager) DeleteFilterEntry(name string, filter string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/flt-%s/e-%s", tenant, filter, name)
	return sm.DeleteByDn(dn, models.VzentryClassName)
}

func (sm *ServiceManager) UpdateFilterEntry(name string, filter string, tenant string, description string, vzEntryattr models.FilterEntryAttributes) (*models.FilterEntry, error) {
	rn := fmt.Sprintf("e-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/flt-%s", tenant, filter)
	vzEntry := models.NewFilterEntry(rn, parentDn, description, vzEntryattr)

	vzEntry.Status = "modified"
	err := sm.Save(vzEntry)
	return vzEntry, err

}

func (sm *ServiceManager) ListFilterEntry(filter string, tenant string) ([]*models.FilterEntry, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/flt-%s/vzEntry.json", baseurlStr, tenant, filter)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FilterEntryListFromContainer(cont)

	return list, err
}

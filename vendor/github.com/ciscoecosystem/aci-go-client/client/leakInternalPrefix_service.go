package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLeakInternalPrefix(ip string, vrf string, tenant string, description string, nameAlias string, leakInternalPrefixAttr models.LeakInternalPrefixAttributes) (*models.LeakInternalPrefix, error) {
	rn := fmt.Sprintf(models.RnleakInternalPrefix, ip)
	parentDn := fmt.Sprintf(models.ParentDnleakInternalPrefix, tenant, vrf)
	leakInternalPrefix := models.NewLeakInternalPrefix(rn, parentDn, description, nameAlias, leakInternalPrefixAttr)
	err := sm.Save(leakInternalPrefix)
	return leakInternalPrefix, err
}

func (sm *ServiceManager) ReadLeakInternalPrefix(ip string, vrf string, tenant string) (*models.LeakInternalPrefix, error) {
	dn := fmt.Sprintf(models.DnleakInternalPrefix, tenant, vrf, ip)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	leakInternalPrefix := models.LeakInternalPrefixFromContainer(cont)
	return leakInternalPrefix, nil
}

func (sm *ServiceManager) DeleteLeakInternalPrefix(ip string, vrf string, tenant string) error {
	dn := fmt.Sprintf(models.DnleakInternalPrefix, tenant, vrf, ip)
	return sm.DeleteByDn(dn, models.LeakInternalPrefixClassName)
}

func (sm *ServiceManager) UpdateLeakInternalPrefix(ip string, vrf string, tenant string, description string, nameAlias string, leakInternalPrefixAttr models.LeakInternalPrefixAttributes) (*models.LeakInternalPrefix, error) {
	rn := fmt.Sprintf(models.RnleakInternalPrefix, ip)
	parentDn := fmt.Sprintf(models.ParentDnleakInternalPrefix, tenant, vrf)
	leakInternalPrefix := models.NewLeakInternalPrefix(rn, parentDn, description, nameAlias, leakInternalPrefixAttr)
	leakInternalPrefix.Status = "modified"
	err := sm.Save(leakInternalPrefix)
	return leakInternalPrefix, err
}

func (sm *ServiceManager) ListLeakInternalPrefix(vrf string, tenant string) ([]*models.LeakInternalPrefix, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ctx-%s/leakroutes/leakInternalPrefix.json", models.BaseurlStr, tenant, vrf)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.LeakInternalPrefixListFromContainer(cont)
	return list, err
}

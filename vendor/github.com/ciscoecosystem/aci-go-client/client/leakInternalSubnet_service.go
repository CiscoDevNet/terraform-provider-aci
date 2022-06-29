package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLeakInternalSubnet(ip string, vrf string, tenant string, description string, nameAlias string, leakInternalSubnetAttr models.LeakInternalSubnetAttributes) (*models.LeakInternalSubnet, error) {
	rn := fmt.Sprintf(models.RnleakInternalSubnet, ip)
	parentDn := fmt.Sprintf(models.ParentDnleakInternalSubnet, tenant, vrf)
	leakInternalSubnet := models.NewLeakInternalSubnet(rn, parentDn, description, nameAlias, leakInternalSubnetAttr)
	err := sm.Save(leakInternalSubnet)
	return leakInternalSubnet, err
}

func (sm *ServiceManager) ReadLeakInternalSubnet(ip string, vrf string, tenant string) (*models.LeakInternalSubnet, error) {
	dn := fmt.Sprintf(models.DnleakInternalSubnet, tenant, vrf, ip)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	leakInternalSubnet := models.LeakInternalSubnetFromContainer(cont)
	return leakInternalSubnet, nil
}

func (sm *ServiceManager) DeleteLeakInternalSubnet(ip string, vrf string, tenant string) error {
	dn := fmt.Sprintf(models.DnleakInternalSubnet, tenant, vrf, ip)
	return sm.DeleteByDn(dn, models.LeakinternalsubnetClassName)
}

func (sm *ServiceManager) UpdateLeakInternalSubnet(ip string, vrf string, tenant string, description string, nameAlias string, leakInternalSubnetAttr models.LeakInternalSubnetAttributes) (*models.LeakInternalSubnet, error) {
	rn := fmt.Sprintf(models.RnleakInternalSubnet, ip)
	parentDn := fmt.Sprintf(models.ParentDnleakInternalSubnet, tenant, vrf)
	leakInternalSubnet := models.NewLeakInternalSubnet(rn, parentDn, description, nameAlias, leakInternalSubnetAttr)
	leakInternalSubnet.Status = "modified"
	err := sm.Save(leakInternalSubnet)
	return leakInternalSubnet, err
}

func (sm *ServiceManager) ListLeakInternalSubnet(vrf string, tenant string) ([]*models.LeakInternalSubnet, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ctx-%s/leakroutes/leakInternalSubnet.json", models.BaseurlStr, tenant, vrf)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.LeakInternalSubnetListFromContainer(cont)
	return list, err
}

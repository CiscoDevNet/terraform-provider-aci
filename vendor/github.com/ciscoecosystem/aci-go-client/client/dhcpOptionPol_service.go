package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateDHCPOptionPolicy(name string, tenant string, description string, dhcpOptionPolattr models.DHCPOptionPolicyAttributes) (*models.DHCPOptionPolicy, error) {
	rn := fmt.Sprintf("dhcpoptpol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	dhcpOptionPol := models.NewDHCPOptionPolicy(rn, parentDn, description, dhcpOptionPolattr)
	err := sm.Save(dhcpOptionPol)
	return dhcpOptionPol, err
}

func (sm *ServiceManager) ReadDHCPOptionPolicy(name string, tenant string) (*models.DHCPOptionPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/dhcpoptpol-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpOptionPol := models.DHCPOptionPolicyFromContainer(cont)
	return dhcpOptionPol, nil
}

func (sm *ServiceManager) DeleteDHCPOptionPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/dhcpoptpol-%s", tenant, name)
	return sm.DeleteByDn(dn, models.DhcpoptionpolClassName)
}

func (sm *ServiceManager) UpdateDHCPOptionPolicy(name string, tenant string, description string, dhcpOptionPolattr models.DHCPOptionPolicyAttributes) (*models.DHCPOptionPolicy, error) {
	rn := fmt.Sprintf("dhcpoptpol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	dhcpOptionPol := models.NewDHCPOptionPolicy(rn, parentDn, description, dhcpOptionPolattr)

	dhcpOptionPol.Status = "modified"
	err := sm.Save(dhcpOptionPol)
	return dhcpOptionPol, err

}

func (sm *ServiceManager) ListDHCPOptionPolicy(tenant string) ([]*models.DHCPOptionPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/dhcpOptionPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.DHCPOptionPolicyListFromContainer(cont)

	return list, err
}

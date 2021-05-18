package client

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateDHCPOption(name string, dhcp_option_policy string, tenant string, dhcpOptionattr models.DHCPOptionAttributes) (*models.DHCPOption, error) {
	rn := fmt.Sprintf("opt-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/dhcpoptpol-%s", tenant, dhcp_option_policy)
	dhcpOption := models.NewDHCPOption(rn, parentDn, dhcpOptionattr)
	err := sm.Save(dhcpOption)
	return dhcpOption, err
}

func (sm *ServiceManager) ReadDHCPOption(name string, dhcp_option_policy string, tenant string) (*models.DHCPOption, error) {
	dn := fmt.Sprintf("uni/tn-%s/dhcpoptpol-%s/opt-%s", tenant, dhcp_option_policy, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpOption := models.DHCPOptionFromContainer(cont)
	return dhcpOption, nil
}

func (sm *ServiceManager) DeleteDHCPOption(name string, dhcp_option_policy string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/dhcpoptpol-%s/opt-%s", tenant, dhcp_option_policy, name)
	return sm.DeleteByDn(dn, models.DhcpoptionClassName)
}

func (sm *ServiceManager) UpdateDHCPOption(name string, dhcp_option_policy string, tenant string, dhcpOptionattr models.DHCPOptionAttributes) (*models.DHCPOption, error) {
	rn := fmt.Sprintf("opt-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/dhcpoptpol-%s", tenant, dhcp_option_policy)
	dhcpOption := models.NewDHCPOption(rn, parentDn, dhcpOptionattr)

	dhcpOption.Status = "modified"
	err := sm.Save(dhcpOption)
	return dhcpOption, err

}

func (sm *ServiceManager) ListDHCPOption(dhcp_option_policy string, tenant string) ([]*models.DHCPOption, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/dhcpoptpol-%s/dhcpOption.json", baseurlStr, tenant, dhcp_option_policy)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.DHCPOptionListFromContainer(cont)

	return list, err
}

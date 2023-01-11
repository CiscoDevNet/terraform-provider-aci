package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateDhcpRelayGwExtIp(parentDn string, description string, dhcpRelayGwExtIpAttr models.DhcpRelayGwExtIpAttributes) (*models.DhcpRelayGwExtIp, error) {
	dhcpRelayGwExtIp := models.NewDhcpRelayGwExtIp(models.RndhcpRelayGwExtIp, parentDn, description, dhcpRelayGwExtIpAttr)
	err := sm.Save(dhcpRelayGwExtIp)
	return dhcpRelayGwExtIp, err
}

func (sm *ServiceManager) ReadDhcpRelayGwExtIp(parentDn string) (*models.DhcpRelayGwExtIp, error) {
	dn := fmt.Sprintf("%s/%s", parentDn, models.RndhcpRelayGwExtIp)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpRelayGwExtIp := models.DhcpRelayGwExtIpFromContainer(cont)
	return dhcpRelayGwExtIp, nil
}

func (sm *ServiceManager) DeleteDhcpRelayGwExtIp(parentDn string) error {
	dn := fmt.Sprintf("%s/%s", parentDn, models.RndhcpRelayGwExtIp)
	return sm.DeleteByDn(dn, models.DhcprelaygwextipClassName)
}

func (sm *ServiceManager) UpdateDhcpRelayGwExtIp(parentDn string, description string, dhcpRelayGwExtIpAttr models.DhcpRelayGwExtIpAttributes) (*models.DhcpRelayGwExtIp, error) {
	dhcpRelayGwExtIp := models.NewDhcpRelayGwExtIp(models.RndhcpRelayGwExtIp, parentDn, description, dhcpRelayGwExtIpAttr)
	dhcpRelayGwExtIp.Status = "modified"
	err := sm.Save(dhcpRelayGwExtIp)
	return dhcpRelayGwExtIp, err
}

func (sm *ServiceManager) ListDhcpRelayGwExtIp(parentDn string) ([]*models.DhcpRelayGwExtIp, error) {
	dnUrl := fmt.Sprintf("%s/%s/dhcpRelayGwExtIp.json", models.BaseurlStr, parentDn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.DhcpRelayGwExtIpListFromContainer(cont)
	return list, err
}

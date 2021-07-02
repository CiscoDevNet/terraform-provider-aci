package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreatePortSecurityPolicy(name string, description string, l2PortSecurityPolattr models.PortSecurityPolicyAttributes) (*models.PortSecurityPolicy, error) {
	rn := fmt.Sprintf("infra/portsecurityP-%s", name)
	parentDn := fmt.Sprintf("uni")
	l2PortSecurityPol := models.NewPortSecurityPolicy(rn, parentDn, description, l2PortSecurityPolattr)
	err := sm.Save(l2PortSecurityPol)
	return l2PortSecurityPol, err
}

func (sm *ServiceManager) ReadPortSecurityPolicy(name string) (*models.PortSecurityPolicy, error) {
	dn := fmt.Sprintf("uni/infra/portsecurityP-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l2PortSecurityPol := models.PortSecurityPolicyFromContainer(cont)
	return l2PortSecurityPol, nil
}

func (sm *ServiceManager) DeletePortSecurityPolicy(name string) error {
	dn := fmt.Sprintf("uni/infra/portsecurityP-%s", name)
	return sm.DeleteByDn(dn, models.L2portsecuritypolClassName)
}

func (sm *ServiceManager) UpdatePortSecurityPolicy(name string, description string, l2PortSecurityPolattr models.PortSecurityPolicyAttributes) (*models.PortSecurityPolicy, error) {
	rn := fmt.Sprintf("infra/portsecurityP-%s", name)
	parentDn := fmt.Sprintf("uni")
	l2PortSecurityPol := models.NewPortSecurityPolicy(rn, parentDn, description, l2PortSecurityPolattr)

	l2PortSecurityPol.Status = "modified"
	err := sm.Save(l2PortSecurityPol)
	return l2PortSecurityPol, err

}

func (sm *ServiceManager) ListPortSecurityPolicy() ([]*models.PortSecurityPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/l2PortSecurityPol.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.PortSecurityPolicyListFromContainer(cont)

	return list, err
}

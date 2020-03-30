package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateMonitoringPolicy(name string, tenant string, description string, monEPGPolattr models.MonitoringPolicyAttributes) (*models.MonitoringPolicy, error) {
	rn := fmt.Sprintf("monepg-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	monEPGPol := models.NewMonitoringPolicy(rn, parentDn, description, monEPGPolattr)
	err := sm.Save(monEPGPol)
	return monEPGPol, err
}

func (sm *ServiceManager) ReadMonitoringPolicy(name string, tenant string) (*models.MonitoringPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/monepg-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	monEPGPol := models.MonitoringPolicyFromContainer(cont)
	return monEPGPol, nil
}

func (sm *ServiceManager) DeleteMonitoringPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/monepg-%s", tenant, name)
	return sm.DeleteByDn(dn, models.MonepgpolClassName)
}

func (sm *ServiceManager) UpdateMonitoringPolicy(name string, tenant string, description string, monEPGPolattr models.MonitoringPolicyAttributes) (*models.MonitoringPolicy, error) {
	rn := fmt.Sprintf("monepg-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	monEPGPol := models.NewMonitoringPolicy(rn, parentDn, description, monEPGPolattr)

	monEPGPol.Status = "modified"
	err := sm.Save(monEPGPol)
	return monEPGPol, err

}

func (sm *ServiceManager) ListMonitoringPolicy(tenant string) ([]*models.MonitoringPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/monEPGPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.MonitoringPolicyListFromContainer(cont)

	return list, err
}

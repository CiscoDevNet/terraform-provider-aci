package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateIPSLAMonitoringPolicy(name string, tenant string, description string, nameAlias string, fvIPSLAMonitoringPolAttr models.IPSLAMonitoringPolicyAttributes) (*models.IPSLAMonitoringPolicy, error) {
	rn := fmt.Sprintf(models.RnfvIPSLAMonitoringPol, name)
	parentDn := fmt.Sprintf(models.ParentDnfvIPSLAMonitoringPol, tenant)
	fvIPSLAMonitoringPol := models.NewIPSLAMonitoringPolicy(rn, parentDn, description, nameAlias, fvIPSLAMonitoringPolAttr)
	err := sm.Save(fvIPSLAMonitoringPol)
	return fvIPSLAMonitoringPol, err
}

func (sm *ServiceManager) ReadIPSLAMonitoringPolicy(name string, tenant string) (*models.IPSLAMonitoringPolicy, error) {
	dn := fmt.Sprintf(models.DnfvIPSLAMonitoringPol, tenant, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvIPSLAMonitoringPol := models.IPSLAMonitoringPolicyFromContainer(cont)
	return fvIPSLAMonitoringPol, nil
}

func (sm *ServiceManager) DeleteIPSLAMonitoringPolicy(name string, tenant string) error {
	dn := fmt.Sprintf(models.DnfvIPSLAMonitoringPol, tenant, name)
	return sm.DeleteByDn(dn, models.FvipslamonitoringpolClassName)
}

func (sm *ServiceManager) UpdateIPSLAMonitoringPolicy(name string, tenant string, description string, nameAlias string, fvIPSLAMonitoringPolAttr models.IPSLAMonitoringPolicyAttributes) (*models.IPSLAMonitoringPolicy, error) {
	rn := fmt.Sprintf(models.RnfvIPSLAMonitoringPol, name)
	parentDn := fmt.Sprintf(models.ParentDnfvIPSLAMonitoringPol, tenant)
	fvIPSLAMonitoringPol := models.NewIPSLAMonitoringPolicy(rn, parentDn, description, nameAlias, fvIPSLAMonitoringPolAttr)
	fvIPSLAMonitoringPol.Status = "modified"
	err := sm.Save(fvIPSLAMonitoringPol)
	return fvIPSLAMonitoringPol, err
}

func (sm *ServiceManager) ListIPSLAMonitoringPolicy(tenant string) ([]*models.IPSLAMonitoringPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/fvIPSLAMonitoringPol.json", models.BaseurlStr, tenant)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.IPSLAMonitoringPolicyListFromContainer(cont)
	return list, err
}

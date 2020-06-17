package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateVMMSecurityPolicy(domain_tDn string, application_epg string, application_profile string, tenant string, description string, vmmSecPattr models.VMMSecurityPolicyAttributes) (*models.VMMSecurityPolicy, error) {
	rn := fmt.Sprintf("sec")
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rsdomAtt-[%s]", tenant, application_profile, application_epg, domain_tDn)
	vmmSecP := models.NewVMMSecurityPolicy(rn, parentDn, description, vmmSecPattr)
	err := sm.Save(vmmSecP)
	return vmmSecP, err
}

func (sm *ServiceManager) ReadVMMSecurityPolicy(domain_tDn string, application_epg string, application_profile string, tenant string) (*models.VMMSecurityPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rsdomAtt-[%s]/sec", tenant, application_profile, application_epg, domain_tDn)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vmmSecP := models.VMMSecurityPolicyFromContainer(cont)
	return vmmSecP, nil
}

func (sm *ServiceManager) DeleteVMMSecurityPolicy(domain_tDn string, application_epg string, application_profile string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rsdomAtt-[%s]/sec", tenant, application_profile, application_epg, domain_tDn)
	return sm.DeleteByDn(dn, models.VmmsecpClassName)
}

func (sm *ServiceManager) UpdateVMMSecurityPolicy(domain_tDn string, application_epg string, application_profile string, tenant string, description string, vmmSecPattr models.VMMSecurityPolicyAttributes) (*models.VMMSecurityPolicy, error) {
	rn := fmt.Sprintf("sec")
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rsdomAtt-[%s]", tenant, application_profile, application_epg, domain_tDn)
	vmmSecP := models.NewVMMSecurityPolicy(rn, parentDn, description, vmmSecPattr)

	vmmSecP.Status = "modified"
	err := sm.Save(vmmSecP)
	return vmmSecP, err

}

func (sm *ServiceManager) ListVMMSecurityPolicy(domain_tDn string, application_epg string, application_profile string, tenant string) ([]*models.VMMSecurityPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/epg-%s/rsdomAtt-[%s]/vmmSecP.json", baseurlStr, tenant, application_profile, application_epg, domain_tDn)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.VMMSecurityPolicyListFromContainer(cont)

	return list, err
}

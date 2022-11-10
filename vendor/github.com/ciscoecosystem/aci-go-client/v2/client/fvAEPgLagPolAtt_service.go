package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateApplicationEPGLagPolicy(domain_dn string, application_epg string, application_profile string, tenant string, fvAEPgLagPolAttAttr models.ApplicationEPGLagPolicyAttributes) (*models.ApplicationEPGLagPolicy, error) {
	rn := fmt.Sprintf(models.RnfvAEPgLagPolAtt)
	parentDn := fmt.Sprintf(models.ParentDnfvAEPgLagPolAtt, tenant, application_profile, application_epg, domain_dn)
	fvAEPgLagPolAtt := models.NewApplicationEPGLagPolicy(rn, parentDn, fvAEPgLagPolAttAttr)
	err := sm.Save(fvAEPgLagPolAtt)
	return fvAEPgLagPolAtt, err
}

func (sm *ServiceManager) ReadApplicationEPGLagPolicy(domain_dn string, application_epg string, application_profile string, tenant string) (*models.ApplicationEPGLagPolicy, error) {
	dn := fmt.Sprintf(models.DnfvAEPgLagPolAtt, tenant, application_profile, application_epg, domain_dn)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvAEPgLagPolAtt := models.ApplicationEPGLagPolicyFromContainer(cont)
	return fvAEPgLagPolAtt, nil
}

func (sm *ServiceManager) DeleteApplicationEPGLagPolicy(domain_dn string, application_epg string, application_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnfvAEPgLagPolAtt, tenant, application_profile, application_epg, domain_dn)
	return sm.DeleteByDn(dn, models.FvAEPgLagPolAttClassName)
}

func (sm *ServiceManager) UpdateApplicationEPGLagPolicy(domain_dn string, application_epg string, application_profile string, tenant string, fvAEPgLagPolAttAttr models.ApplicationEPGLagPolicyAttributes) (*models.ApplicationEPGLagPolicy, error) {
	rn := fmt.Sprintf(models.RnfvAEPgLagPolAtt)
	parentDn := fmt.Sprintf(models.ParentDnfvAEPgLagPolAtt, tenant, application_profile, application_epg, domain_dn)
	fvAEPgLagPolAtt := models.NewApplicationEPGLagPolicy(rn, parentDn, fvAEPgLagPolAttAttr)
	fvAEPgLagPolAtt.Status = "modified"
	err := sm.Save(fvAEPgLagPolAtt)
	return fvAEPgLagPolAtt, err
}

func (sm *ServiceManager) ListApplicationEPGLagPolicy(domain_dn string, application_epg string, application_profile string, tenant string) ([]*models.ApplicationEPGLagPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/epg-%s/rsdomAtt-[%s]/fvAEPgLagPolAtt.json", models.BaseurlStr, tenant, application_profile, application_epg, domain_dn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.ApplicationEPGLagPolicyListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationfvRsVmmVSwitchEnhancedLagPol(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsvmmVSwitchEnhancedLagPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "fvRsVmmVSwitchEnhancedLagPol", dn, annotation, tDn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationfvRsVmmVSwitchEnhancedLagPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsvmmVSwitchEnhancedLagPol", parentDn)
	return sm.DeleteByDn(dn, "fvRsVmmVSwitchEnhancedLagPol")
}

func (sm *ServiceManager) ReadRelationfvRsVmmVSwitchEnhancedLagPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsVmmVSwitchEnhancedLagPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsVmmVSwitchEnhancedLagPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

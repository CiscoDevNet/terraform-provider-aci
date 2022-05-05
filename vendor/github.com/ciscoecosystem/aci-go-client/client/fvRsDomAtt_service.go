package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFVDomain(tDn string, application_epg string, application_profile string, tenant string, description string, fvRsDomAttattr models.FVDomainAttributes) (*models.FVDomain, error) {
	rn := fmt.Sprintf("rsdomAtt-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant, application_profile, application_epg)
	fvRsDomAtt := models.NewFVDomain(rn, parentDn, description, fvRsDomAttattr)
	err := sm.Save(fvRsDomAtt)
	return fvRsDomAtt, err
}

func (sm *ServiceManager) ReadFVDomain(tDn string, application_epg string, application_profile string, tenant string) (*models.FVDomain, error) {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rsdomAtt-[%s]", tenant, application_profile, application_epg, tDn)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvRsDomAtt := models.FVDomainFromContainer(cont)
	return fvRsDomAtt, nil
}

func (sm *ServiceManager) DeleteFVDomain(tDn string, application_epg string, application_profile string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rsdomAtt-[%s]", tenant, application_profile, application_epg, tDn)
	return sm.DeleteByDn(dn, models.FvrsdomattClassName)
}

func (sm *ServiceManager) UpdateFVDomain(tDn string, application_epg string, application_profile string, tenant string, description string, fvRsDomAttattr models.FVDomainAttributes) (*models.FVDomain, error) {
	rn := fmt.Sprintf("rsdomAtt-[%s]", tDn)
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant, application_profile, application_epg)
	fvRsDomAtt := models.NewFVDomain(rn, parentDn, description, fvRsDomAttattr)

	fvRsDomAtt.Status = "modified"
	err := sm.Save(fvRsDomAtt)
	return fvRsDomAtt, err

}

func (sm *ServiceManager) ListFVDomain(application_epg string, application_profile string, tenant string) ([]*models.FVDomain, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/epg-%s/fvRsDomAtt.json", baseurlStr, tenant, application_profile, application_epg)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FVDomainListFromContainer(cont)

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
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn+"/epglagpolatt", "fvRsVmmVSwitchEnhancedLagPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsVmmVSwitchEnhancedLagPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

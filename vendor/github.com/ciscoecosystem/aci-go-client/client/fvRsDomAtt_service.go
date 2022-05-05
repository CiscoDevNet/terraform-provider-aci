package client

import (
	"fmt"

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

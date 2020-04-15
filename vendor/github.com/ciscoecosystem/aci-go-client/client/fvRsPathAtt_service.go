package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"



	


)









func (sm *ServiceManager) CreateStaticPath(tDn string ,application_epg string ,application_profile string ,tenant string , description string, fvRsPathAttattr models.StaticPathAttributes) (*models.StaticPath, error) {	
	rn := fmt.Sprintf("rspathAtt-[%s]",tDn)
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant ,application_profile ,application_epg )
	fvRsPathAtt := models.NewStaticPath(rn, parentDn, description, fvRsPathAttattr)
	err := sm.Save(fvRsPathAtt)
	return fvRsPathAtt, err
}

func (sm *ServiceManager) ReadStaticPath(tDn string ,application_epg string ,application_profile string ,tenant string ) (*models.StaticPath, error) {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rspathAtt-[%s]", tenant ,application_profile ,application_epg ,tDn )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvRsPathAtt := models.StaticPathFromContainer(cont)
	return fvRsPathAtt, nil
}

func (sm *ServiceManager) DeleteStaticPath(tDn string ,application_epg string ,application_profile string ,tenant string ) error {
	dn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rspathAtt-[%s]", tenant ,application_profile ,application_epg ,tDn )
	return sm.DeleteByDn(dn, models.FvrspathattClassName)
}

func (sm *ServiceManager) UpdateStaticPath(tDn string ,application_epg string ,application_profile string ,tenant string  ,description string, fvRsPathAttattr models.StaticPathAttributes) (*models.StaticPath, error) {
	rn := fmt.Sprintf("rspathAtt-[%s]",tDn)
	parentDn := fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", tenant ,application_profile ,application_epg )
	fvRsPathAtt := models.NewStaticPath(rn, parentDn, description, fvRsPathAttattr)

    fvRsPathAtt.Status = "modified"
	err := sm.Save(fvRsPathAtt)
	return fvRsPathAtt, err

}

func (sm *ServiceManager) ListStaticPath(application_epg string ,application_profile string ,tenant string ) ([]*models.StaticPath, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/epg-%s/fvRsPathAtt.json", baseurlStr , tenant ,application_profile ,application_epg )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.StaticPathListFromContainer(cont)

	return list, err
}



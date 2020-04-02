package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/ciscoecosystem/aci-go-client/container"



	


)









func (sm *ServiceManager) CreateSPANSourceGroup(name string ,tenant string , description string, spanSrcGrpattr models.SPANSourceGroupAttributes) (*models.SPANSourceGroup, error) {	
	rn := fmt.Sprintf("srcgrp-%s",name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant )
	spanSrcGrp := models.NewSPANSourceGroup(rn, parentDn, description, spanSrcGrpattr)
	err := sm.Save(spanSrcGrp)
	return spanSrcGrp, err
}

func (sm *ServiceManager) ReadSPANSourceGroup(name string ,tenant string ) (*models.SPANSourceGroup, error) {
	dn := fmt.Sprintf("uni/tn-%s/srcgrp-%s", tenant ,name )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	spanSrcGrp := models.SPANSourceGroupFromContainer(cont)
	return spanSrcGrp, nil
}

func (sm *ServiceManager) DeleteSPANSourceGroup(name string ,tenant string ) error {
	dn := fmt.Sprintf("uni/tn-%s/srcgrp-%s", tenant ,name )
	return sm.DeleteByDn(dn, models.SpansrcgrpClassName)
}

func (sm *ServiceManager) UpdateSPANSourceGroup(name string ,tenant string  ,description string, spanSrcGrpattr models.SPANSourceGroupAttributes) (*models.SPANSourceGroup, error) {
	rn := fmt.Sprintf("srcgrp-%s",name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant )
	spanSrcGrp := models.NewSPANSourceGroup(rn, parentDn, description, spanSrcGrpattr)

    spanSrcGrp.Status = "modified"
	err := sm.Save(spanSrcGrp)
	return spanSrcGrp, err

}

func (sm *ServiceManager) ListSPANSourceGroup(tenant string ) ([]*models.SPANSourceGroup, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/spanSrcGrp.json", baseurlStr , tenant )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.SPANSourceGroupListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationspanRsSrcGrpToFilterGrpFromSPANSourceGroup( parentDn, tnSpanFilterGrpName string) error {
	dn := fmt.Sprintf("%s/rssrcGrpToFilterGrp", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnSpanFilterGrpName": "%s"
								
			}
		}
	}`, "spanRsSrcGrpToFilterGrp", dn,tnSpanFilterGrpName))

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

func (sm *ServiceManager) DeleteRelationspanRsSrcGrpToFilterGrpFromSPANSourceGroup(parentDn string) error{
	dn := fmt.Sprintf("%s/rssrcGrpToFilterGrp", parentDn)
	return sm.DeleteByDn(dn , "spanRsSrcGrpToFilterGrp")
}

func (sm *ServiceManager) ReadRelationspanRsSrcGrpToFilterGrpFromSPANSourceGroup( parentDn string) (interface{},error) {
	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/%s/%s.json",baseurlStr,parentDn,"spanRsSrcGrpToFilterGrp")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont,"spanRsSrcGrpToFilterGrp")
	
	if len(contList) > 0 {
		dat := models.G(contList[0], "tnSpanFilterGrpName")
		return dat, err
	} else {
		return nil,err
	}
		





}


package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	



	


)









func (sm *ServiceManager) CreateAttachableAccessEntityProfile(name string , description string, infraAttEntityPattr models.AttachableAccessEntityProfileAttributes) (*models.AttachableAccessEntityProfile, error) {	
	rn := fmt.Sprintf("infra/attentp-%s",name)
	parentDn := fmt.Sprintf("uni")
	infraAttEntityP := models.NewAttachableAccessEntityProfile(rn, parentDn, description, infraAttEntityPattr)
	err := sm.Save(infraAttEntityP)
	return infraAttEntityP, err
}

func (sm *ServiceManager) ReadAttachableAccessEntityProfile(name string ) (*models.AttachableAccessEntityProfile, error) {
	dn := fmt.Sprintf("uni/infra/attentp-%s", name )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraAttEntityP := models.AttachableAccessEntityProfileFromContainer(cont)
	return infraAttEntityP, nil
}

func (sm *ServiceManager) DeleteAttachableAccessEntityProfile(name string ) error {
	dn := fmt.Sprintf("uni/infra/attentp-%s", name )
	return sm.DeleteByDn(dn, models.InfraattentitypClassName)
}

func (sm *ServiceManager) UpdateAttachableAccessEntityProfile(name string  ,description string, infraAttEntityPattr models.AttachableAccessEntityProfileAttributes) (*models.AttachableAccessEntityProfile, error) {
	rn := fmt.Sprintf("infra/attentp-%s",name)
	parentDn := fmt.Sprintf("uni")
	infraAttEntityP := models.NewAttachableAccessEntityProfile(rn, parentDn, description, infraAttEntityPattr)

    infraAttEntityP.Status = "modified"
	err := sm.Save(infraAttEntityP)
	return infraAttEntityP, err

}

func (sm *ServiceManager) ListAttachableAccessEntityProfile() ([]*models.AttachableAccessEntityProfile, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/infraAttEntityP.json", baseurlStr )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.AttachableAccessEntityProfileListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsDomPFromAttachableAccessEntityProfile( parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsdomP-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s"				
			}
		}
	}`, "infraRsDomP", dn))

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

func (sm *ServiceManager) DeleteRelationinfraRsDomPFromAttachableAccessEntityProfile(parentDn , tDn string) error{
	dn := fmt.Sprintf("%s/rsdomP-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn , "infraRsDomP")
}

func (sm *ServiceManager) ReadRelationinfraRsDomPFromAttachableAccessEntityProfile( parentDn string) (interface{},error) {
	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/%s/%s.json",baseurlStr,parentDn,"infraRsDomP")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont,"infraRsDomP")
	
	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList{
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err
			





}


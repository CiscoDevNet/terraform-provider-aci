package client



import (
	
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/ciscoecosystem/aci-go-client/container"
	
)


func (sm *ServiceManager) CreateUserManagement(description string, nameAlias string, aaaUserEpAttr models.UserManagementAttributes) (*models.UserManagement, error) {	
	rn := fmt.Sprintf(models.RnaaaUserEp , )
	parentDn := fmt.Sprintf(models.ParentDnaaaUserEp, )
	aaaUserEp := models.NewUserManagement(rn, parentDn, description, nameAlias, aaaUserEpAttr)
	err := sm.Save(aaaUserEp)
	return aaaUserEp, err
}

func (sm *ServiceManager) ReadUserManagement() (*models.UserManagement, error) {
	dn := fmt.Sprintf(models.DnaaaUserEp, )
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaUserEp := models.UserManagementFromContainer(cont)
	return aaaUserEp, nil
}

func (sm *ServiceManager) DeleteUserManagement() error {
	dn := fmt.Sprintf(models.DnaaaUserEp, )
	return sm.DeleteByDn(dn, models.AaauserepClassName)
}

func (sm *ServiceManager) UpdateUserManagement(description string, nameAlias string, aaaUserEpAttr models.UserManagementAttributes) (*models.UserManagement, error) {
	rn := fmt.Sprintf(models.RnaaaUserEp , )
	parentDn := fmt.Sprintf(models.ParentDnaaaUserEp, )
	aaaUserEp := models.NewUserManagement(rn, parentDn, description, nameAlias, aaaUserEpAttr)
    aaaUserEp.Status = "modified"
	err := sm.Save(aaaUserEp)
	return aaaUserEp, err
}

func (sm *ServiceManager) ListUserManagement() ([]*models.UserManagement, error) {	
	dnUrl := fmt.Sprintf("%s/uni/aaaUserEp.json", models.BaseurlStr)
    cont, err := sm.GetViaURL(dnUrl)
	list := models.UserManagementListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationaaaRsToUserEp(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rstoUserEp", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "aaaRsToUserEp", dn, annotation, tDn))

	
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

func (sm *ServiceManager) DeleteRelationaaaRsToUserEp(parentDn string) error{
	dn := fmt.Sprintf("%s/rstoUserEp", parentDn)
	return sm.DeleteByDn(dn , "aaaRsToUserEp")
}

func (sm *ServiceManager) ReadRelationaaaRsToUserEp(parentDn string) (interface{},error) {	
	dnUrl := fmt.Sprintf("%s/%s/%s.json",models.BaseurlStr,parentDn,"aaaRsToUserEp")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont,"aaaRsToUserEp")
	
		if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}


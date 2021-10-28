package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateAAAAuthentication(description string, nameAlias string, aaaAuthRealmAttr models.AAAAuthenticationAttributes) (*models.AAAAuthentication, error) {
	rn := fmt.Sprintf(models.RnaaaAuthRealm)
	parentDn := fmt.Sprintf(models.ParentDnaaaAuthRealm)
	aaaAuthRealm := models.NewAAAAuthentication(rn, parentDn, description, nameAlias, aaaAuthRealmAttr)
	err := sm.Save(aaaAuthRealm)
	return aaaAuthRealm, err
}

func (sm *ServiceManager) ReadAAAAuthentication() (*models.AAAAuthentication, error) {
	dn := fmt.Sprintf(models.DnaaaAuthRealm)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaAuthRealm := models.AAAAuthenticationFromContainer(cont)
	return aaaAuthRealm, nil
}

func (sm *ServiceManager) DeleteAAAAuthentication() error {
	dn := fmt.Sprintf(models.DnaaaAuthRealm)
	return sm.DeleteByDn(dn, models.AaaauthrealmClassName)
}

func (sm *ServiceManager) UpdateAAAAuthentication(description string, nameAlias string, aaaAuthRealmAttr models.AAAAuthenticationAttributes) (*models.AAAAuthentication, error) {
	rn := fmt.Sprintf(models.RnaaaAuthRealm)
	parentDn := fmt.Sprintf(models.ParentDnaaaAuthRealm)
	aaaAuthRealm := models.NewAAAAuthentication(rn, parentDn, description, nameAlias, aaaAuthRealmAttr)
	aaaAuthRealm.Status = "modified"
	err := sm.Save(aaaAuthRealm)
	return aaaAuthRealm, err
}

func (sm *ServiceManager) ListAAAAuthentication() ([]*models.AAAAuthentication, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/aaaAuthRealm.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.AAAAuthenticationListFromContainer(cont)
	return list, err
}

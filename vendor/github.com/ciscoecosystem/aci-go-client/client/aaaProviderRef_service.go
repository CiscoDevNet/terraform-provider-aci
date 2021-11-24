package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)


func (sm *ServiceManager) CreateProviderGroupMember(name string, duo_provider_group string, description string, nameAlias string, aaaProviderRefAttr models.ProviderGroupMemberAttributes) (*models.ProviderGroupMember, error) {	
	rn := fmt.Sprintf(models.RnaaaProviderRef , name)
	parentDn := fmt.Sprintf(models.ParentDnaaaProviderRef, duo_provider_group )
	aaaProviderRef := models.NewProviderGroupMember(rn, parentDn, description, nameAlias, aaaProviderRefAttr)
	err := sm.Save(aaaProviderRef)
	return aaaProviderRef, err
}

func (sm *ServiceManager) ReadProviderGroupMember(name string, duo_provider_group string, ) (*models.ProviderGroupMember, error) {
	dn := fmt.Sprintf(models.DnaaaProviderRef, duo_provider_group,name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaProviderRef := models.ProviderGroupMemberFromContainer(cont)
	return aaaProviderRef, nil
}

func (sm *ServiceManager) DeleteProviderGroupMember(name string, duo_provider_group string, ) error {
	dn := fmt.Sprintf(models.DnaaaProviderRef, duo_provider_group,name)
	return sm.DeleteByDn(dn, models.AaaproviderrefClassName)
}

func (sm *ServiceManager) UpdateProviderGroupMember(name string, duo_provider_group string, description string, nameAlias string, aaaProviderRefAttr models.ProviderGroupMemberAttributes) (*models.ProviderGroupMember, error) {
	rn := fmt.Sprintf(models.RnaaaProviderRef , name)
	parentDn := fmt.Sprintf(models.ParentDnaaaProviderRef, duo_provider_group )
	aaaProviderRef := models.NewProviderGroupMember(rn, parentDn, description, nameAlias, aaaProviderRefAttr)
    aaaProviderRef.Status = "modified"
	err := sm.Save(aaaProviderRef)
	return aaaProviderRef, err
}

func (sm *ServiceManager) ListProviderGroupMember(duo_provider_group string ) ([]*models.ProviderGroupMember, error) {	
	dnUrl := fmt.Sprintf("%s/uni/userext/duoext/duoprovidergroup-%s/aaaProviderRef.json", models.BaseurlStr, duo_provider_group )
    cont, err := sm.GetViaURL(dnUrl)
	list := models.ProviderGroupMemberListFromContainer(cont)
	return list, err
}


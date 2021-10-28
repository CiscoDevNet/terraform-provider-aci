package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLDAPGroupMap(name string, description string, nameAlias string, aaaLdapGroupMapAttr models.LDAPGroupMapAttributes) (*models.LDAPGroupMap, error) {
	rn := fmt.Sprintf(models.RnaaaLdapGroupMap, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaLdapGroupMap)
	aaaLdapGroupMap := models.NewLDAPGroupMap(rn, parentDn, description, nameAlias, aaaLdapGroupMapAttr)
	err := sm.Save(aaaLdapGroupMap)
	return aaaLdapGroupMap, err
}

func (sm *ServiceManager) ReadLDAPGroupMap(name string) (*models.LDAPGroupMap, error) {
	dn := fmt.Sprintf(models.DnaaaLdapGroupMap, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaLdapGroupMap := models.LDAPGroupMapFromContainer(cont)
	return aaaLdapGroupMap, nil
}

func (sm *ServiceManager) DeleteLDAPGroupMap(name string) error {
	dn := fmt.Sprintf(models.DnaaaLdapGroupMap, name)
	return sm.DeleteByDn(dn, models.AaaldapgroupmapClassName)
}

func (sm *ServiceManager) UpdateLDAPGroupMap(name string, description string, nameAlias string, aaaLdapGroupMapAttr models.LDAPGroupMapAttributes) (*models.LDAPGroupMap, error) {
	rn := fmt.Sprintf(models.RnaaaLdapGroupMap, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaLdapGroupMap)
	aaaLdapGroupMap := models.NewLDAPGroupMap(rn, parentDn, description, nameAlias, aaaLdapGroupMapAttr)
	aaaLdapGroupMap.Status = "modified"
	err := sm.Save(aaaLdapGroupMap)
	return aaaLdapGroupMap, err
}

func (sm *ServiceManager) ListLDAPGroupMap() ([]*models.LDAPGroupMap, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/ldapext/aaaLdapGroupMap.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.LDAPGroupMapListFromContainer(cont)
	return list, err
}

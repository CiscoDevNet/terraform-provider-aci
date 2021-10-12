package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLDAPGroupMapruleref(name string, ldap_group_map string, description string, nameAlias string, aaaLdapGroupMapRuleRefAttr models.LDAPGroupMaprulerefAttributes) (*models.LDAPGroupMapruleref, error) {
	rn := fmt.Sprintf(models.RnaaaLdapGroupMapRuleRef, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaLdapGroupMapRuleRef, ldap_group_map)
	aaaLdapGroupMapRuleRef := models.NewLDAPGroupMapruleref(rn, parentDn, description, nameAlias, aaaLdapGroupMapRuleRefAttr)
	err := sm.Save(aaaLdapGroupMapRuleRef)
	return aaaLdapGroupMapRuleRef, err
}

func (sm *ServiceManager) ReadLDAPGroupMapruleref(name string, ldap_group_map string) (*models.LDAPGroupMapruleref, error) {
	dn := fmt.Sprintf(models.DnaaaLdapGroupMapRuleRef, ldap_group_map, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaLdapGroupMapRuleRef := models.LDAPGroupMaprulerefFromContainer(cont)
	return aaaLdapGroupMapRuleRef, nil
}

func (sm *ServiceManager) DeleteLDAPGroupMapruleref(name string, ldap_group_map string) error {
	dn := fmt.Sprintf(models.DnaaaLdapGroupMapRuleRef, ldap_group_map, name)
	return sm.DeleteByDn(dn, models.AaaldapgroupmaprulerefClassName)
}

func (sm *ServiceManager) UpdateLDAPGroupMapruleref(name string, ldap_group_map string, description string, nameAlias string, aaaLdapGroupMapRuleRefAttr models.LDAPGroupMaprulerefAttributes) (*models.LDAPGroupMapruleref, error) {
	rn := fmt.Sprintf(models.RnaaaLdapGroupMapRuleRef, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaLdapGroupMapRuleRef, ldap_group_map)
	aaaLdapGroupMapRuleRef := models.NewLDAPGroupMapruleref(rn, parentDn, description, nameAlias, aaaLdapGroupMapRuleRefAttr)
	aaaLdapGroupMapRuleRef.Status = "modified"
	err := sm.Save(aaaLdapGroupMapRuleRef)
	return aaaLdapGroupMapRuleRef, err
}

func (sm *ServiceManager) ListLDAPGroupMapruleref(ldap_group_map string) ([]*models.LDAPGroupMapruleref, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/ldapext/ldapgroupmap-%s/aaaLdapGroupMapRuleRef.json", models.BaseurlStr, ldap_group_map)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.LDAPGroupMaprulerefListFromContainer(cont)
	return list, err
}

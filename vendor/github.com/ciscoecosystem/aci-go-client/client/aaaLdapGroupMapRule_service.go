package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLDAPGroupMapRule(name string, description string, nameAlias string, aaaLdapGroupMapRuleAttr models.LDAPGroupMapRuleAttributes) (*models.LDAPGroupMapRule, error) {
	rn := fmt.Sprintf(models.RnaaaLdapGroupMapRule, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaLdapGroupMapRule)
	aaaLdapGroupMapRule := models.NewLDAPGroupMapRule(rn, parentDn, description, nameAlias, aaaLdapGroupMapRuleAttr)
	err := sm.Save(aaaLdapGroupMapRule)
	return aaaLdapGroupMapRule, err
}

func (sm *ServiceManager) ReadLDAPGroupMapRule(name string) (*models.LDAPGroupMapRule, error) {
	dn := fmt.Sprintf(models.DnaaaLdapGroupMapRule, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaLdapGroupMapRule := models.LDAPGroupMapRuleFromContainer(cont)
	return aaaLdapGroupMapRule, nil
}

func (sm *ServiceManager) DeleteLDAPGroupMapRule(name string) error {
	dn := fmt.Sprintf(models.DnaaaLdapGroupMapRule, name)
	return sm.DeleteByDn(dn, models.AaaldapgroupmapruleClassName)
}

func (sm *ServiceManager) UpdateLDAPGroupMapRule(name string, description string, nameAlias string, aaaLdapGroupMapRuleAttr models.LDAPGroupMapRuleAttributes) (*models.LDAPGroupMapRule, error) {
	rn := fmt.Sprintf(models.RnaaaLdapGroupMapRule, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaLdapGroupMapRule)
	aaaLdapGroupMapRule := models.NewLDAPGroupMapRule(rn, parentDn, description, nameAlias, aaaLdapGroupMapRuleAttr)
	aaaLdapGroupMapRule.Status = "modified"
	err := sm.Save(aaaLdapGroupMapRule)
	return aaaLdapGroupMapRule, err
}

func (sm *ServiceManager) ListLDAPGroupMapRule() ([]*models.LDAPGroupMapRule, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/duoext/aaaLdapGroupMapRule.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.LDAPGroupMapRuleListFromContainer(cont)
	return list, err
}

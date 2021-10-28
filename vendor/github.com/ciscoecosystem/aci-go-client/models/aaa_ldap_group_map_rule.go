package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaLdapGroupMapRule        = "uni/userext/duoext/ldapgroupmaprule-%s"
	RnaaaLdapGroupMapRule        = "ldapgroupmaprule-%s"
	ParentDnaaaLdapGroupMapRule  = "uni/userext/duoext"
	AaaldapgroupmapruleClassName = "aaaLdapGroupMapRule"
)

type LDAPGroupMapRule struct {
	BaseAttributes
	NameAliasAttribute
	LDAPGroupMapRuleAttributes
}

type LDAPGroupMapRuleAttributes struct {
	Annotation string `json:",omitempty"`
	Groupdn    string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewLDAPGroupMapRule(aaaLdapGroupMapRuleRn, parentDn, description, nameAlias string, aaaLdapGroupMapRuleAttr LDAPGroupMapRuleAttributes) *LDAPGroupMapRule {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaLdapGroupMapRuleRn)
	return &LDAPGroupMapRule{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaaldapgroupmapruleClassName,
			Rn:                aaaLdapGroupMapRuleRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		LDAPGroupMapRuleAttributes: aaaLdapGroupMapRuleAttr,
	}
}

func (aaaLdapGroupMapRule *LDAPGroupMapRule) ToMap() (map[string]string, error) {
	aaaLdapGroupMapRuleMap, err := aaaLdapGroupMapRule.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaLdapGroupMapRule.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaLdapGroupMapRuleMap, key, value)
	}
	A(aaaLdapGroupMapRuleMap, "annotation", aaaLdapGroupMapRule.Annotation)
	A(aaaLdapGroupMapRuleMap, "groupdn", aaaLdapGroupMapRule.Groupdn)
	A(aaaLdapGroupMapRuleMap, "name", aaaLdapGroupMapRule.Name)
	return aaaLdapGroupMapRuleMap, err
}

func LDAPGroupMapRuleFromContainerList(cont *container.Container, index int) *LDAPGroupMapRule {
	LDAPGroupMapRuleCont := cont.S("imdata").Index(index).S(AaaldapgroupmapruleClassName, "attributes")
	return &LDAPGroupMapRule{
		BaseAttributes{
			DistinguishedName: G(LDAPGroupMapRuleCont, "dn"),
			Description:       G(LDAPGroupMapRuleCont, "descr"),
			Status:            G(LDAPGroupMapRuleCont, "status"),
			ClassName:         AaaldapgroupmapruleClassName,
			Rn:                G(LDAPGroupMapRuleCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(LDAPGroupMapRuleCont, "nameAlias"),
		},
		LDAPGroupMapRuleAttributes{
			Annotation: G(LDAPGroupMapRuleCont, "annotation"),
			Groupdn:    G(LDAPGroupMapRuleCont, "groupdn"),
			Name:       G(LDAPGroupMapRuleCont, "name"),
		},
	}
}

func LDAPGroupMapRuleFromContainer(cont *container.Container) *LDAPGroupMapRule {
	return LDAPGroupMapRuleFromContainerList(cont, 0)
}

func LDAPGroupMapRuleListFromContainer(cont *container.Container) []*LDAPGroupMapRule {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*LDAPGroupMapRule, length)
	for i := 0; i < length; i++ {
		arr[i] = LDAPGroupMapRuleFromContainerList(cont, i)
	}
	return arr
}

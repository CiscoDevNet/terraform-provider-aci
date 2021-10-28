package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaLdapGroupMapRuleRef        = "uni/userext/ldapext/ldapgroupmap-%s/ldapgroupmapruleref-%s"
	RnaaaLdapGroupMapRuleRef        = "ldapgroupmapruleref-%s"
	ParentDnaaaLdapGroupMapRuleRef  = "uni/userext/ldapext/ldapgroupmap-%s"
	AaaldapgroupmaprulerefClassName = "aaaLdapGroupMapRuleRef"
)

type LDAPGroupMapruleref struct {
	BaseAttributes
	NameAliasAttribute
	LDAPGroupMaprulerefAttributes
}

type LDAPGroupMaprulerefAttributes struct {
	Name       string `json:",omitempty"`
	Annotation string `json:",omitempty"`
}

func NewLDAPGroupMapruleref(aaaLdapGroupMapRuleRefRn, parentDn, description, nameAlias string, aaaLdapGroupMapRuleRefAttr LDAPGroupMaprulerefAttributes) *LDAPGroupMapruleref {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaLdapGroupMapRuleRefRn)
	return &LDAPGroupMapruleref{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaaldapgroupmaprulerefClassName,
			Rn:                aaaLdapGroupMapRuleRefRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		LDAPGroupMaprulerefAttributes: aaaLdapGroupMapRuleRefAttr,
	}
}

func (aaaLdapGroupMapRuleRef *LDAPGroupMapruleref) ToMap() (map[string]string, error) {
	aaaLdapGroupMapRuleRefMap, err := aaaLdapGroupMapRuleRef.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaLdapGroupMapRuleRef.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaLdapGroupMapRuleRefMap, key, value)
	}
	A(aaaLdapGroupMapRuleRefMap, "name", aaaLdapGroupMapRuleRef.Name)
	A(aaaLdapGroupMapRuleRefMap, "annotation", aaaLdapGroupMapRuleRef.Annotation)
	return aaaLdapGroupMapRuleRefMap, err
}

func LDAPGroupMaprulerefFromContainerList(cont *container.Container, index int) *LDAPGroupMapruleref {
	LDAPGroupMaprulerefCont := cont.S("imdata").Index(index).S(AaaldapgroupmaprulerefClassName, "attributes")
	return &LDAPGroupMapruleref{
		BaseAttributes{
			DistinguishedName: G(LDAPGroupMaprulerefCont, "dn"),
			Description:       G(LDAPGroupMaprulerefCont, "descr"),
			Status:            G(LDAPGroupMaprulerefCont, "status"),
			ClassName:         AaaldapgroupmaprulerefClassName,
			Rn:                G(LDAPGroupMaprulerefCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(LDAPGroupMaprulerefCont, "nameAlias"),
		},
		LDAPGroupMaprulerefAttributes{
			Name:       G(LDAPGroupMaprulerefCont, "name"),
			Annotation: G(LDAPGroupMaprulerefCont, "annotation"),
		},
	}
}

func LDAPGroupMaprulerefFromContainer(cont *container.Container) *LDAPGroupMapruleref {
	return LDAPGroupMaprulerefFromContainerList(cont, 0)
}

func LDAPGroupMaprulerefListFromContainer(cont *container.Container) []*LDAPGroupMapruleref {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*LDAPGroupMapruleref, length)
	for i := 0; i < length; i++ {
		arr[i] = LDAPGroupMaprulerefFromContainerList(cont, i)
	}
	return arr
}

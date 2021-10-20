package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaLdapGroupMap        = "uni/userext/ldapext/ldapgroupmap-%s"
	RnaaaLdapGroupMap        = "ldapgroupmap-%s"
	ParentDnaaaLdapGroupMap  = "uni/userext/ldapext"
	AaaldapgroupmapClassName = "aaaLdapGroupMap"
)

type LDAPGroupMap struct {
	BaseAttributes
	NameAliasAttribute
	LDAPGroupMapAttributes
}

type LDAPGroupMapAttributes struct {
	Name       string `json:",omitempty"`
	Annotation string `json:",omitempty"`
}

func NewLDAPGroupMap(aaaLdapGroupMapRn, parentDn, description, nameAlias string, aaaLdapGroupMapAttr LDAPGroupMapAttributes) *LDAPGroupMap {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaLdapGroupMapRn)
	return &LDAPGroupMap{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaaldapgroupmapClassName,
			Rn:                aaaLdapGroupMapRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		LDAPGroupMapAttributes: aaaLdapGroupMapAttr,
	}
}

func (aaaLdapGroupMap *LDAPGroupMap) ToMap() (map[string]string, error) {
	aaaLdapGroupMapMap, err := aaaLdapGroupMap.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaLdapGroupMap.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaLdapGroupMapMap, key, value)
	}
	A(aaaLdapGroupMapMap, "name", aaaLdapGroupMap.Name)
	return aaaLdapGroupMapMap, err
}

func LDAPGroupMapFromContainerList(cont *container.Container, index int) *LDAPGroupMap {
	LDAPGroupMapCont := cont.S("imdata").Index(index).S(AaaldapgroupmapClassName, "attributes")
	return &LDAPGroupMap{
		BaseAttributes{
			DistinguishedName: G(LDAPGroupMapCont, "dn"),
			Description:       G(LDAPGroupMapCont, "descr"),
			Status:            G(LDAPGroupMapCont, "status"),
			ClassName:         AaaldapgroupmapClassName,
			Rn:                G(LDAPGroupMapCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(LDAPGroupMapCont, "nameAlias"),
		},
		LDAPGroupMapAttributes{
			Name:       G(LDAPGroupMapCont, "name"),
			Annotation: G(LDAPGroupMapCont, "annotation"),
		},
	}
}

func LDAPGroupMapFromContainer(cont *container.Container) *LDAPGroupMap {
	return LDAPGroupMapFromContainerList(cont, 0)
}

func LDAPGroupMapListFromContainer(cont *container.Container) []*LDAPGroupMap {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*LDAPGroupMap, length)
	for i := 0; i < length; i++ {
		arr[i] = LDAPGroupMapFromContainerList(cont, i)
	}
	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaDuoProviderGroup        = "uni/userext/duoext/duoprovidergroup-%s"
	RnaaaDuoProviderGroup        = "duoprovidergroup-%s"
	ParentDnaaaDuoProviderGroup  = "uni/userext/duoext"
	AaaduoprovidergroupClassName = "aaaDuoProviderGroup"
)

type DuoProviderGroup struct {
	BaseAttributes
	NameAliasAttribute
	DuoProviderGroupAttributes
}

type DuoProviderGroupAttributes struct {
	Annotation        string `json:",omitempty"`
	AuthChoice        string `json:",omitempty"`
	LdapGroupMapRef   string `json:",omitempty"`
	Name              string `json:",omitempty"`
	ProviderType      string `json:",omitempty"`
	SecFacAuthMethods string `json:",omitempty"`
}

func NewDuoProviderGroup(aaaDuoProviderGroupRn, parentDn, description, nameAlias string, aaaDuoProviderGroupAttr DuoProviderGroupAttributes) *DuoProviderGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaDuoProviderGroupRn)
	return &DuoProviderGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaaduoprovidergroupClassName,
			Rn:                aaaDuoProviderGroupRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		DuoProviderGroupAttributes: aaaDuoProviderGroupAttr,
	}
}

func (aaaDuoProviderGroup *DuoProviderGroup) ToMap() (map[string]string, error) {
	aaaDuoProviderGroupMap, err := aaaDuoProviderGroup.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaDuoProviderGroup.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaDuoProviderGroupMap, key, value)
	}
	A(aaaDuoProviderGroupMap, "annotation", aaaDuoProviderGroup.Annotation)
	A(aaaDuoProviderGroupMap, "authChoice", aaaDuoProviderGroup.AuthChoice)
	A(aaaDuoProviderGroupMap, "ldapGroupMapRef", aaaDuoProviderGroup.LdapGroupMapRef)
	A(aaaDuoProviderGroupMap, "name", aaaDuoProviderGroup.Name)
	A(aaaDuoProviderGroupMap, "providerType", aaaDuoProviderGroup.ProviderType)
	A(aaaDuoProviderGroupMap, "secFacAuthMethods", aaaDuoProviderGroup.SecFacAuthMethods)
	return aaaDuoProviderGroupMap, err
}

func DuoProviderGroupFromContainerList(cont *container.Container, index int) *DuoProviderGroup {
	DuoProviderGroupCont := cont.S("imdata").Index(index).S(AaaduoprovidergroupClassName, "attributes")
	return &DuoProviderGroup{
		BaseAttributes{
			DistinguishedName: G(DuoProviderGroupCont, "dn"),
			Description:       G(DuoProviderGroupCont, "descr"),
			Status:            G(DuoProviderGroupCont, "status"),
			ClassName:         AaaduoprovidergroupClassName,
			Rn:                G(DuoProviderGroupCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(DuoProviderGroupCont, "nameAlias"),
		},
		DuoProviderGroupAttributes{
			Annotation:        G(DuoProviderGroupCont, "annotation"),
			AuthChoice:        G(DuoProviderGroupCont, "authChoice"),
			LdapGroupMapRef:   G(DuoProviderGroupCont, "ldapGroupMapRef"),
			Name:              G(DuoProviderGroupCont, "name"),
			ProviderType:      G(DuoProviderGroupCont, "providerType"),
			SecFacAuthMethods: G(DuoProviderGroupCont, "secFacAuthMethods"),
		},
	}
}

func DuoProviderGroupFromContainer(cont *container.Container) *DuoProviderGroup {
	return DuoProviderGroupFromContainerList(cont, 0)
}

func DuoProviderGroupListFromContainer(cont *container.Container) []*DuoProviderGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*DuoProviderGroup, length)
	for i := 0; i < length; i++ {
		arr[i] = DuoProviderGroupFromContainerList(cont, i)
	}
	return arr
}

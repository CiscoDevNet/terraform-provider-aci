package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaProviderRef        = "uni/userext/duoext/duoprovidergroup-%s/providerref-%s"
	RnaaaProviderRef        = "providerref-%s"
	ParentDnaaaProviderRef  = "uni/userext/duoext/duoprovidergroup-%s"
	AaaproviderrefClassName = "aaaProviderRef"
)

type ProviderGroupMember struct {
	BaseAttributes
	NameAliasAttribute
	ProviderGroupMemberAttributes
}

type ProviderGroupMemberAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Order      string `json:",omitempty"`
}

func NewProviderGroupMember(aaaProviderRefRn, parentDn, description, nameAlias string, aaaProviderRefAttr ProviderGroupMemberAttributes) *ProviderGroupMember {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaProviderRefRn)
	return &ProviderGroupMember{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaaproviderrefClassName,
			Rn:                aaaProviderRefRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		ProviderGroupMemberAttributes: aaaProviderRefAttr,
	}
}

func (aaaProviderRef *ProviderGroupMember) ToMap() (map[string]string, error) {
	aaaProviderRefMap, err := aaaProviderRef.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaProviderRef.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaProviderRefMap, key, value)
	}
	A(aaaProviderRefMap, "annotation", aaaProviderRef.Annotation)
	A(aaaProviderRefMap, "name", aaaProviderRef.Name)
	A(aaaProviderRefMap, "order", aaaProviderRef.Order)
	return aaaProviderRefMap, err
}

func ProviderGroupMemberFromContainerList(cont *container.Container, index int) *ProviderGroupMember {
	ProviderGroupMemberCont := cont.S("imdata").Index(index).S(AaaproviderrefClassName, "attributes")
	return &ProviderGroupMember{
		BaseAttributes{
			DistinguishedName: G(ProviderGroupMemberCont, "dn"),
			Description:       G(ProviderGroupMemberCont, "descr"),
			Status:            G(ProviderGroupMemberCont, "status"),
			ClassName:         AaaproviderrefClassName,
			Rn:                G(ProviderGroupMemberCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(ProviderGroupMemberCont, "nameAlias"),
		},
		ProviderGroupMemberAttributes{
			Annotation: G(ProviderGroupMemberCont, "annotation"),
			Name:       G(ProviderGroupMemberCont, "name"),
			Order:      G(ProviderGroupMemberCont, "order"),
		},
	}
}

func ProviderGroupMemberFromContainer(cont *container.Container) *ProviderGroupMember {
	return ProviderGroupMemberFromContainerList(cont, 0)
}

func ProviderGroupMemberListFromContainer(cont *container.Container) []*ProviderGroupMember {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*ProviderGroupMember, length)
	for i := 0; i < length; i++ {
		arr[i] = ProviderGroupMemberFromContainerList(cont, i)
	}
	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaTacacsPlusProviderGroup        = "uni/userext/tacacsext/tacacsplusprovidergroup-%s"
	RnaaaTacacsPlusProviderGroup        = "tacacsplusprovidergroup-%s"
	ParentDnaaaTacacsPlusProviderGroup  = "uni/userext/tacacsext"
	AaatacacsplusprovidergroupClassName = "aaaTacacsPlusProviderGroup"
)

type TACACSPlusProviderGroup struct {
	BaseAttributes
	NameAliasAttribute
	TACACSPlusProviderGroupAttributes
}

type TACACSPlusProviderGroupAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewTACACSPlusProviderGroup(aaaTacacsPlusProviderGroupRn, parentDn, description, nameAlias string, aaaTacacsPlusProviderGroupAttr TACACSPlusProviderGroupAttributes) *TACACSPlusProviderGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaTacacsPlusProviderGroupRn)
	return &TACACSPlusProviderGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaatacacsplusprovidergroupClassName,
			Rn:                aaaTacacsPlusProviderGroupRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		TACACSPlusProviderGroupAttributes: aaaTacacsPlusProviderGroupAttr,
	}
}

func (aaaTacacsPlusProviderGroup *TACACSPlusProviderGroup) ToMap() (map[string]string, error) {
	aaaTacacsPlusProviderGroupMap, err := aaaTacacsPlusProviderGroup.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaTacacsPlusProviderGroup.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaTacacsPlusProviderGroupMap, key, value)
	}
	A(aaaTacacsPlusProviderGroupMap, "annotation", aaaTacacsPlusProviderGroup.Annotation)
	A(aaaTacacsPlusProviderGroupMap, "name", aaaTacacsPlusProviderGroup.Name)
	return aaaTacacsPlusProviderGroupMap, err
}

func TACACSPlusProviderGroupFromContainerList(cont *container.Container, index int) *TACACSPlusProviderGroup {
	TACACSPlusProviderGroupCont := cont.S("imdata").Index(index).S(AaatacacsplusprovidergroupClassName, "attributes")
	return &TACACSPlusProviderGroup{
		BaseAttributes{
			DistinguishedName: G(TACACSPlusProviderGroupCont, "dn"),
			Description:       G(TACACSPlusProviderGroupCont, "descr"),
			Status:            G(TACACSPlusProviderGroupCont, "status"),
			ClassName:         AaatacacsplusprovidergroupClassName,
			Rn:                G(TACACSPlusProviderGroupCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(TACACSPlusProviderGroupCont, "nameAlias"),
		},
		TACACSPlusProviderGroupAttributes{
			Annotation: G(TACACSPlusProviderGroupCont, "annotation"),
			Name:       G(TACACSPlusProviderGroupCont, "name"),
		},
	}
}

func TACACSPlusProviderGroupFromContainer(cont *container.Container) *TACACSPlusProviderGroup {
	return TACACSPlusProviderGroupFromContainerList(cont, 0)
}

func TACACSPlusProviderGroupListFromContainer(cont *container.Container) []*TACACSPlusProviderGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*TACACSPlusProviderGroup, length)
	for i := 0; i < length; i++ {
		arr[i] = TACACSPlusProviderGroupFromContainerList(cont, i)
	}
	return arr
}

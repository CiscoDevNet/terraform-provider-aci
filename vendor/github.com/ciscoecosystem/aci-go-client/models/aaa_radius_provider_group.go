package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaRadiusProviderGroup        = "uni/userext/radiusext/radiusprovidergroup-%s"
	RnaaaRadiusProviderGroup        = "radiusprovidergroup-%s"
	ParentDnaaaRadiusProviderGroup  = "uni/userext/radiusext"
	AaaradiusprovidergroupClassName = "aaaRadiusProviderGroup"
)

type RADIUSProviderGroup struct {
	BaseAttributes
	NameAliasAttribute
	RADIUSProviderGroupAttributes
}

type RADIUSProviderGroupAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewRADIUSProviderGroup(aaaRadiusProviderGroupRn, parentDn, description, nameAlias string, aaaRadiusProviderGroupAttr RADIUSProviderGroupAttributes) *RADIUSProviderGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaRadiusProviderGroupRn)
	return &RADIUSProviderGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaaradiusprovidergroupClassName,
			Rn:                aaaRadiusProviderGroupRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RADIUSProviderGroupAttributes: aaaRadiusProviderGroupAttr,
	}
}

func (aaaRadiusProviderGroup *RADIUSProviderGroup) ToMap() (map[string]string, error) {
	aaaRadiusProviderGroupMap, err := aaaRadiusProviderGroup.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaRadiusProviderGroup.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaRadiusProviderGroupMap, key, value)
	}
	A(aaaRadiusProviderGroupMap, "annotation", aaaRadiusProviderGroup.Annotation)
	A(aaaRadiusProviderGroupMap, "name", aaaRadiusProviderGroup.Name)
	return aaaRadiusProviderGroupMap, err
}

func RADIUSProviderGroupFromContainerList(cont *container.Container, index int) *RADIUSProviderGroup {
	RADIUSProviderGroupCont := cont.S("imdata").Index(index).S(AaaradiusprovidergroupClassName, "attributes")
	return &RADIUSProviderGroup{
		BaseAttributes{
			DistinguishedName: G(RADIUSProviderGroupCont, "dn"),
			Description:       G(RADIUSProviderGroupCont, "descr"),
			Status:            G(RADIUSProviderGroupCont, "status"),
			ClassName:         AaaradiusprovidergroupClassName,
			Rn:                G(RADIUSProviderGroupCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RADIUSProviderGroupCont, "nameAlias"),
		},
		RADIUSProviderGroupAttributes{
			Annotation: G(RADIUSProviderGroupCont, "annotation"),
			Name:       G(RADIUSProviderGroupCont, "name"),
		},
	}
}

func RADIUSProviderGroupFromContainer(cont *container.Container) *RADIUSProviderGroup {
	return RADIUSProviderGroupFromContainerList(cont, 0)
}

func RADIUSProviderGroupListFromContainer(cont *container.Container) []*RADIUSProviderGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RADIUSProviderGroup, length)
	for i := 0; i < length; i++ {
		arr[i] = RADIUSProviderGroupFromContainerList(cont, i)
	}
	return arr
}

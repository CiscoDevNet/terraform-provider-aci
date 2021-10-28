package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaSamlProviderGroup        = "uni/userext/samlext/samlprovidergroup-%s"
	RnaaaSamlProviderGroup        = "samlprovidergroup-%s"
	ParentDnaaaSamlProviderGroup  = "uni/userext/samlext"
	AaasamlprovidergroupClassName = "aaaSamlProviderGroup"
)

type SAMLProviderGroup struct {
	BaseAttributes
	NameAliasAttribute
	SAMLProviderGroupAttributes
}

type SAMLProviderGroupAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewSAMLProviderGroup(aaaSamlProviderGroupRn, parentDn, description, nameAlias string, aaaSamlProviderGroupAttr SAMLProviderGroupAttributes) *SAMLProviderGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaSamlProviderGroupRn)
	return &SAMLProviderGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaasamlprovidergroupClassName,
			Rn:                aaaSamlProviderGroupRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		SAMLProviderGroupAttributes: aaaSamlProviderGroupAttr,
	}
}

func (aaaSamlProviderGroup *SAMLProviderGroup) ToMap() (map[string]string, error) {
	aaaSamlProviderGroupMap, err := aaaSamlProviderGroup.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaSamlProviderGroup.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaSamlProviderGroupMap, key, value)
	}
	A(aaaSamlProviderGroupMap, "annotation", aaaSamlProviderGroup.Annotation)
	A(aaaSamlProviderGroupMap, "name", aaaSamlProviderGroup.Name)
	return aaaSamlProviderGroupMap, err
}

func SAMLProviderGroupFromContainerList(cont *container.Container, index int) *SAMLProviderGroup {
	SAMLProviderGroupCont := cont.S("imdata").Index(index).S(AaasamlprovidergroupClassName, "attributes")
	return &SAMLProviderGroup{
		BaseAttributes{
			DistinguishedName: G(SAMLProviderGroupCont, "dn"),
			Description:       G(SAMLProviderGroupCont, "descr"),
			Status:            G(SAMLProviderGroupCont, "status"),
			ClassName:         AaasamlprovidergroupClassName,
			Rn:                G(SAMLProviderGroupCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(SAMLProviderGroupCont, "nameAlias"),
		},
		SAMLProviderGroupAttributes{
			Annotation: G(SAMLProviderGroupCont, "annotation"),
			Name:       G(SAMLProviderGroupCont, "name"),
		},
	}
}

func SAMLProviderGroupFromContainer(cont *container.Container) *SAMLProviderGroup {
	return SAMLProviderGroupFromContainerList(cont, 0)
}

func SAMLProviderGroupListFromContainer(cont *container.Container) []*SAMLProviderGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*SAMLProviderGroup, length)
	for i := 0; i < length; i++ {
		arr[i] = SAMLProviderGroupFromContainerList(cont, i)
	}
	return arr
}

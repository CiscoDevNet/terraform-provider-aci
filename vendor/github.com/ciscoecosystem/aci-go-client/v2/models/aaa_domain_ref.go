package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnaaaDomainRef        = "%s/domain-%s"
	RnaaaDomainRef        = "domain-%s"
	AaadomainrefClassName = "aaaDomainRef"
)

type AaaDomainRef struct {
	BaseAttributes
	NameAliasAttribute
	AaaDomainRefAttributes
}

type AaaDomainRefAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewAaaDomainRef(aaaDomainRefRn, parentDn, nameAlias string, aaaDomainRefAttr AaaDomainRefAttributes) *AaaDomainRef {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaDomainRefRn)
	return &AaaDomainRef{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         AaadomainrefClassName,
			Rn:                aaaDomainRefRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		AaaDomainRefAttributes: aaaDomainRefAttr,
	}
}

func (aaaDomainRef *AaaDomainRef) ToMap() (map[string]string, error) {
	aaaDomainRefMap, err := aaaDomainRef.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := aaaDomainRef.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(aaaDomainRefMap, key, value)
	}

	A(aaaDomainRefMap, "name", aaaDomainRef.Name)
	return aaaDomainRefMap, err
}

func AaaDomainRefFromContainerList(cont *container.Container, index int) *AaaDomainRef {
	AaaDomainRefCont := cont.S("imdata").Index(index).S(AaadomainrefClassName, "attributes")
	return &AaaDomainRef{
		BaseAttributes{
			DistinguishedName: G(AaaDomainRefCont, "dn"),
			Status:            G(AaaDomainRefCont, "status"),
			ClassName:         AaadomainrefClassName,
			Rn:                G(AaaDomainRefCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(AaaDomainRefCont, "nameAlias"),
		},
		AaaDomainRefAttributes{
			Name: G(AaaDomainRefCont, "name"),
		},
	}
}

func AaaDomainRefFromContainer(cont *container.Container) *AaaDomainRef {
	return AaaDomainRefFromContainerList(cont, 0)
}

func AaaDomainRefListFromContainer(cont *container.Container) []*AaaDomainRef {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*AaaDomainRef, length)

	for i := 0; i < length; i++ {
		arr[i] = AaaDomainRefFromContainerList(cont, i)
	}

	return arr
}

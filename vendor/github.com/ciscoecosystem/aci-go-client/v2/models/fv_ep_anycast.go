package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnfvEpAnycast        = "epAnycast-%s"
	FvepanycastClassName = "fvEpAnycast"
)

type AnycastEndpoint struct {
	BaseAttributes
	NameAliasAttribute
	AnycastEndpointAttributes
}

type AnycastEndpointAttributes struct {
	Annotation string `json:",omitempty"`
	Mac        string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewAnycastEndpoint(fvEpAnycastRn, parentDn, description, nameAlias string, fvEpAnycastAttr AnycastEndpointAttributes) *AnycastEndpoint {
	dn := fmt.Sprintf("%s/%s", parentDn, fvEpAnycastRn)
	return &AnycastEndpoint{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvepanycastClassName,
			Rn:                fvEpAnycastRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		AnycastEndpointAttributes: fvEpAnycastAttr,
	}
}

func (fvEpAnycast *AnycastEndpoint) ToMap() (map[string]string, error) {
	fvEpAnycastMap, err := fvEpAnycast.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := fvEpAnycast.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(fvEpAnycastMap, key, value)
	}

	A(fvEpAnycastMap, "mac", fvEpAnycast.Mac)
	A(fvEpAnycastMap, "name", fvEpAnycast.Name)
	return fvEpAnycastMap, err
}

func AnycastEndpointFromContainerList(cont *container.Container, index int) *AnycastEndpoint {
	AnycastEndpointCont := cont.S("imdata").Index(index).S(FvepanycastClassName, "attributes")
	return &AnycastEndpoint{
		BaseAttributes{
			DistinguishedName: G(AnycastEndpointCont, "dn"),
			Description:       G(AnycastEndpointCont, "descr"),
			Status:            G(AnycastEndpointCont, "status"),
			ClassName:         FvepanycastClassName,
			Rn:                G(AnycastEndpointCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(AnycastEndpointCont, "nameAlias"),
		},
		AnycastEndpointAttributes{
			Mac:  G(AnycastEndpointCont, "mac"),
			Name: G(AnycastEndpointCont, "name"),
		},
	}
}

func AnycastEndpointFromContainer(cont *container.Container) *AnycastEndpoint {
	return AnycastEndpointFromContainerList(cont, 0)
}

func AnycastEndpointListFromContainer(cont *container.Container) []*AnycastEndpoint {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*AnycastEndpoint, length)

	for i := 0; i < length; i++ {
		arr[i] = AnycastEndpointFromContainerList(cont, i)
	}

	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	RnfvEpNlb        = "epnlb"
	FvepnlbClassName = "fvEpNlb"
)

type NlbEndpoint struct {
	BaseAttributes
	NameAliasAttribute
	NlbEndpointAttributes
}

type NlbEndpointAttributes struct {
	Annotation string `json:",omitempty"`
	Group      string `json:",omitempty"`
	Mac        string `json:",omitempty"`
	Mode       string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewNlbEndpoint(fvEpNlbRn, parentDn, description, nameAlias string, fvEpNlbAttr NlbEndpointAttributes) *NlbEndpoint {
	dn := fmt.Sprintf("%s/%s", parentDn, fvEpNlbRn)
	return &NlbEndpoint{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvepnlbClassName,
			Rn:                fvEpNlbRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		NlbEndpointAttributes: fvEpNlbAttr,
	}
}

func (fvEpNlb *NlbEndpoint) ToMap() (map[string]string, error) {
	fvEpNlbMap, err := fvEpNlb.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := fvEpNlb.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(fvEpNlbMap, key, value)
	}

	A(fvEpNlbMap, "group", fvEpNlb.Group)
	A(fvEpNlbMap, "mac", fvEpNlb.Mac)
	A(fvEpNlbMap, "mode", fvEpNlb.Mode)
	A(fvEpNlbMap, "name", fvEpNlb.Name)
	return fvEpNlbMap, err
}

func NlbEndpointFromContainerList(cont *container.Container, index int) *NlbEndpoint {
	NlbEndpointCont := cont.S("imdata").Index(index).S(FvepnlbClassName, "attributes")
	return &NlbEndpoint{
		BaseAttributes{
			DistinguishedName: G(NlbEndpointCont, "dn"),
			Description:       G(NlbEndpointCont, "descr"),
			Status:            G(NlbEndpointCont, "status"),
			ClassName:         FvepnlbClassName,
			Rn:                G(NlbEndpointCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(NlbEndpointCont, "nameAlias"),
		},
		NlbEndpointAttributes{
			Group: G(NlbEndpointCont, "group"),
			Mac:   G(NlbEndpointCont, "mac"),
			Mode:  G(NlbEndpointCont, "mode"),
			Name:  G(NlbEndpointCont, "name"),
		},
	}
}

func NlbEndpointFromContainer(cont *container.Container) *NlbEndpoint {
	return NlbEndpointFromContainerList(cont, 0)
}

func NlbEndpointListFromContainer(cont *container.Container) []*NlbEndpoint {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*NlbEndpoint, length)

	for i := 0; i < length; i++ {
		arr[i] = NlbEndpointFromContainerList(cont, i)
	}

	return arr
}

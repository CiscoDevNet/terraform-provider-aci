package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnfvEPgSelector        = "uni/tn-%s/ap-%s/esg-%s/epgselector-[%s]"
	RnfvEPgSelector        = "epgselector-[%s]"
	ParentDnfvEPgSelector  = "uni/tn-%s/ap-%s/esg-%s"
	FvepgselectorClassName = "fvEPgSelector"
)

type EndpointSecurityGroupEPgSelector struct {
	BaseAttributes
	NameAliasAttribute
	EndpointSecurityGroupEPgSelectorAttributes
}

type EndpointSecurityGroupEPgSelectorAttributes struct {
	Annotation string `json:",omitempty"`
	MatchEpgDn string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewEndpointSecurityGroupEPgSelector(fvEPgSelectorRn, parentDn, description, nameAlias string, fvEPgSelectorAttr EndpointSecurityGroupEPgSelectorAttributes) *EndpointSecurityGroupEPgSelector {
	dn := fmt.Sprintf("%s/%s", parentDn, fvEPgSelectorRn)
	return &EndpointSecurityGroupEPgSelector{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvepgselectorClassName,
			Rn:                fvEPgSelectorRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		EndpointSecurityGroupEPgSelectorAttributes: fvEPgSelectorAttr,
	}
}

func (fvEPgSelector *EndpointSecurityGroupEPgSelector) ToMap() (map[string]string, error) {
	fvEPgSelectorMap, err := fvEPgSelector.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := fvEPgSelector.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(fvEPgSelectorMap, key, value)
	}
	A(fvEPgSelectorMap, "annotation", fvEPgSelector.Annotation)
	A(fvEPgSelectorMap, "matchEpgDn", fvEPgSelector.MatchEpgDn)
	A(fvEPgSelectorMap, "name", fvEPgSelector.Name)
	return fvEPgSelectorMap, err
}

func EndpointSecurityGroupEPgSelectorFromContainerList(cont *container.Container, index int) *EndpointSecurityGroupEPgSelector {
	EndpointSecurityGroupEPgSelectorCont := cont.S("imdata").Index(index).S(FvepgselectorClassName, "attributes")
	return &EndpointSecurityGroupEPgSelector{
		BaseAttributes{
			DistinguishedName: G(EndpointSecurityGroupEPgSelectorCont, "dn"),
			Description:       G(EndpointSecurityGroupEPgSelectorCont, "descr"),
			Status:            G(EndpointSecurityGroupEPgSelectorCont, "status"),
			ClassName:         FvepgselectorClassName,
			Rn:                G(EndpointSecurityGroupEPgSelectorCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(EndpointSecurityGroupEPgSelectorCont, "nameAlias"),
		},
		EndpointSecurityGroupEPgSelectorAttributes{
			Annotation: G(EndpointSecurityGroupEPgSelectorCont, "annotation"),
			MatchEpgDn: G(EndpointSecurityGroupEPgSelectorCont, "matchEpgDn"),
			Name:       G(EndpointSecurityGroupEPgSelectorCont, "name"),
		},
	}
}

func EndpointSecurityGroupEPgSelectorFromContainer(cont *container.Container) *EndpointSecurityGroupEPgSelector {
	return EndpointSecurityGroupEPgSelectorFromContainerList(cont, 0)
}

func EndpointSecurityGroupEPgSelectorListFromContainer(cont *container.Container) []*EndpointSecurityGroupEPgSelector {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*EndpointSecurityGroupEPgSelector, length)
	for i := 0; i < length; i++ {
		arr[i] = EndpointSecurityGroupEPgSelectorFromContainerList(cont, i)
	}
	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnfvEPSelector        = "uni/tn-%s/ap-%s/esg-%s/epselector-[%s]"
	RnfvEPSelector        = "epselector-[%s]"
	ParentDnfvEPSelector  = "uni/tn-%s/ap-%s/esg-%s"
	FvepselectorClassName = "fvEPSelector"
)

type EndpointSecurityGroupSelector struct {
	BaseAttributes
	NameAliasAttribute
	EndpointSecurityGroupSelectorAttributes
}

type EndpointSecurityGroupSelectorAttributes struct {
	Annotation      string `json:",omitempty"`
	MatchExpression string `json:",omitempty"`
	Name            string `json:",omitempty"`
}

func NewEndpointSecurityGroupSelector(fvEPSelectorRn, parentDn, description, nameAlias string, fvEPSelectorAttr EndpointSecurityGroupSelectorAttributes) *EndpointSecurityGroupSelector {
	dn := fmt.Sprintf("%s/%s", parentDn, fvEPSelectorRn)
	return &EndpointSecurityGroupSelector{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvepselectorClassName,
			Rn:                fvEPSelectorRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		EndpointSecurityGroupSelectorAttributes: fvEPSelectorAttr,
	}
}

func (fvEPSelector *EndpointSecurityGroupSelector) ToMap() (map[string]string, error) {
	fvEPSelectorMap, err := fvEPSelector.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := fvEPSelector.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(fvEPSelectorMap, key, value)
	}
	A(fvEPSelectorMap, "annotation", fvEPSelector.Annotation)
	A(fvEPSelectorMap, "matchExpression", fvEPSelector.MatchExpression)
	A(fvEPSelectorMap, "name", fvEPSelector.Name)
	return fvEPSelectorMap, err
}

func EndpointSecurityGroupSelectorFromContainerList(cont *container.Container, index int) *EndpointSecurityGroupSelector {
	EndpointSecurityGroupSelectorCont := cont.S("imdata").Index(index).S(FvepselectorClassName, "attributes")
	return &EndpointSecurityGroupSelector{
		BaseAttributes{
			DistinguishedName: G(EndpointSecurityGroupSelectorCont, "dn"),
			Description:       G(EndpointSecurityGroupSelectorCont, "descr"),
			Status:            G(EndpointSecurityGroupSelectorCont, "status"),
			ClassName:         FvepselectorClassName,
			Rn:                G(EndpointSecurityGroupSelectorCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(EndpointSecurityGroupSelectorCont, "nameAlias"),
		},
		EndpointSecurityGroupSelectorAttributes{
			Annotation:      G(EndpointSecurityGroupSelectorCont, "annotation"),
			MatchExpression: G(EndpointSecurityGroupSelectorCont, "matchExpression"),
			Name:            G(EndpointSecurityGroupSelectorCont, "name"),
		},
	}
}

func EndpointSecurityGroupSelectorFromContainer(cont *container.Container) *EndpointSecurityGroupSelector {
	return EndpointSecurityGroupSelectorFromContainerList(cont, 0)
}

func EndpointSecurityGroupSelectorListFromContainer(cont *container.Container) []*EndpointSecurityGroupSelector {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*EndpointSecurityGroupSelector, length)
	for i := 0; i < length; i++ {
		arr[i] = EndpointSecurityGroupSelectorFromContainerList(cont, i)
	}
	return arr
}

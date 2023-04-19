package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnfvTagSelector        = "uni/tn-%s/ap-%s/esg-%s/tagselectorkey-[%s]-value-[%s]"
	RnfvTagSelector        = "tagselectorkey-[%s]-value-[%s]"
	ParentDnfvTagSelector  = "uni/tn-%s/ap-%s/esg-%s"
	FvtagselectorClassName = "fvTagSelector"
)

type EndpointSecurityGroupTagSelector struct {
	BaseAttributes
	NameAliasAttribute
	EndpointSecurityGroupTagSelectorAttributes
}

type EndpointSecurityGroupTagSelectorAttributes struct {
	Annotation    string `json:",omitempty"`
	MatchKey      string `json:",omitempty"`
	MatchValue    string `json:",omitempty"`
	Name          string `json:",omitempty"`
	ValueOperator string `json:",omitempty"`
}

func NewEndpointSecurityGroupTagSelector(fvTagSelectorRn, parentDn, description, nameAlias string, fvTagSelectorAttr EndpointSecurityGroupTagSelectorAttributes) *EndpointSecurityGroupTagSelector {
	dn := fmt.Sprintf("%s/%s", parentDn, fvTagSelectorRn)
	return &EndpointSecurityGroupTagSelector{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvtagselectorClassName,
			Rn:                fvTagSelectorRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		EndpointSecurityGroupTagSelectorAttributes: fvTagSelectorAttr,
	}
}

func (fvTagSelector *EndpointSecurityGroupTagSelector) ToMap() (map[string]string, error) {
	fvTagSelectorMap, err := fvTagSelector.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := fvTagSelector.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(fvTagSelectorMap, key, value)
	}
	A(fvTagSelectorMap, "annotation", fvTagSelector.Annotation)
	A(fvTagSelectorMap, "matchKey", fvTagSelector.MatchKey)
	A(fvTagSelectorMap, "matchValue", fvTagSelector.MatchValue)
	A(fvTagSelectorMap, "name", fvTagSelector.Name)
	A(fvTagSelectorMap, "valueOperator", fvTagSelector.ValueOperator)
	return fvTagSelectorMap, err
}

func EndpointSecurityGroupTagSelectorFromContainerList(cont *container.Container, index int) *EndpointSecurityGroupTagSelector {
	EndpointSecurityGroupTagSelectorCont := cont.S("imdata").Index(index).S(FvtagselectorClassName, "attributes")
	return &EndpointSecurityGroupTagSelector{
		BaseAttributes{
			DistinguishedName: G(EndpointSecurityGroupTagSelectorCont, "dn"),
			Description:       G(EndpointSecurityGroupTagSelectorCont, "descr"),
			Status:            G(EndpointSecurityGroupTagSelectorCont, "status"),
			ClassName:         FvtagselectorClassName,
			Rn:                G(EndpointSecurityGroupTagSelectorCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(EndpointSecurityGroupTagSelectorCont, "nameAlias"),
		},
		EndpointSecurityGroupTagSelectorAttributes{
			Annotation:    G(EndpointSecurityGroupTagSelectorCont, "annotation"),
			MatchKey:      G(EndpointSecurityGroupTagSelectorCont, "matchKey"),
			MatchValue:    G(EndpointSecurityGroupTagSelectorCont, "matchValue"),
			Name:          G(EndpointSecurityGroupTagSelectorCont, "name"),
			ValueOperator: G(EndpointSecurityGroupTagSelectorCont, "valueOperator"),
		},
	}
}

func EndpointSecurityGroupTagSelectorFromContainer(cont *container.Container) *EndpointSecurityGroupTagSelector {
	return EndpointSecurityGroupTagSelectorFromContainerList(cont, 0)
}

func EndpointSecurityGroupTagSelectorListFromContainer(cont *container.Container) []*EndpointSecurityGroupTagSelector {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*EndpointSecurityGroupTagSelector, length)
	for i := 0; i < length; i++ {
		arr[i] = EndpointSecurityGroupTagSelectorFromContainerList(cont, i)
	}
	return arr
}

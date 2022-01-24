package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DninfraSHPortS        = "uni/infra/spaccportprof-%s/shports-%s-typ-%s"
	RninfraSHPortS        = "shports-%s-typ-%s"
	ParentDninfraSHPortS  = "uni/infra/spaccportprof-%s"
	InfrashportsClassName = "infraSHPortS"
)

type SpineAccessPortSelector struct {
	BaseAttributes
	NameAliasAttribute
	SpineAccessPortSelectorAttributes
}

type SpineAccessPortSelectorAttributes struct {
	Annotation                   string `json:",omitempty"`
	Name                         string `json:",omitempty"`
	SpineAccessPortSelector_type string `json:",omitempty"`
}

func NewSpineAccessPortSelector(infraSHPortSRn, parentDn, description, nameAlias string, infraSHPortSAttr SpineAccessPortSelectorAttributes) *SpineAccessPortSelector {
	dn := fmt.Sprintf("%s/%s", parentDn, infraSHPortSRn)
	return &SpineAccessPortSelector{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfrashportsClassName,
			Rn:                infraSHPortSRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		SpineAccessPortSelectorAttributes: infraSHPortSAttr,
	}
}

func (infraSHPortS *SpineAccessPortSelector) ToMap() (map[string]string, error) {
	infraSHPortSMap, err := infraSHPortS.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := infraSHPortS.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(infraSHPortSMap, key, value)
	}
	A(infraSHPortSMap, "annotation", infraSHPortS.Annotation)
	A(infraSHPortSMap, "name", infraSHPortS.Name)
	A(infraSHPortSMap, "type", infraSHPortS.SpineAccessPortSelector_type)
	return infraSHPortSMap, err
}

func SpineAccessPortSelectorFromContainerList(cont *container.Container, index int) *SpineAccessPortSelector {
	SpineAccessPortSelectorCont := cont.S("imdata").Index(index).S(InfrashportsClassName, "attributes")
	return &SpineAccessPortSelector{
		BaseAttributes{
			DistinguishedName: G(SpineAccessPortSelectorCont, "dn"),
			Description:       G(SpineAccessPortSelectorCont, "descr"),
			Status:            G(SpineAccessPortSelectorCont, "status"),
			ClassName:         InfrashportsClassName,
			Rn:                G(SpineAccessPortSelectorCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(SpineAccessPortSelectorCont, "nameAlias"),
		},
		SpineAccessPortSelectorAttributes{
			Annotation:                   G(SpineAccessPortSelectorCont, "annotation"),
			Name:                         G(SpineAccessPortSelectorCont, "name"),
			SpineAccessPortSelector_type: G(SpineAccessPortSelectorCont, "type"),
		},
	}
}

func SpineAccessPortSelectorFromContainer(cont *container.Container) *SpineAccessPortSelector {
	return SpineAccessPortSelectorFromContainerList(cont, 0)
}

func SpineAccessPortSelectorListFromContainer(cont *container.Container) []*SpineAccessPortSelector {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*SpineAccessPortSelector, length)
	for i := 0; i < length; i++ {
		arr[i] = SpineAccessPortSelectorFromContainerList(cont, i)
	}
	return arr
}

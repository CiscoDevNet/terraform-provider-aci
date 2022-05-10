package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnvzRsFiltAtt        = "%s/rsfiltAtt-%s"
	RnvzRsFiltAtt        = "rsfiltAtt-%s"
	VzrsfiltattClassName = "vzRsFiltAtt"
)

type FilterRelationship struct {
	BaseAttributes
	FilterRelationshipAttributes
}

type FilterRelationshipAttributes struct {
	Annotation       string `json:",omitempty"`
	Action           string `json:",omitempty"`
	Directives       string `json:",omitempty"`
	PriorityOverride string `json:",omitempty"`
	TDn              string `json:",omitempty"`
	TnVzFilterName   string `json:",omitempty"`
}

func NewFilterRelationship(vzRsFiltAttRn, parentDn string, vzRsFiltAttAttr FilterRelationshipAttributes) *FilterRelationship {
	dn := fmt.Sprintf("%s/%s", parentDn, vzRsFiltAttRn)
	return &FilterRelationship{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         VzrsfiltattClassName,
			Rn:                vzRsFiltAttRn,
		},
		FilterRelationshipAttributes: vzRsFiltAttAttr,
	}
}

func (vzRsFiltAtt *FilterRelationship) ToMap() (map[string]string, error) {
	vzRsFiltAttMap, err := vzRsFiltAtt.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vzRsFiltAttMap, "action", vzRsFiltAtt.Action)
	A(vzRsFiltAttMap, "directives", vzRsFiltAtt.Directives)
	A(vzRsFiltAttMap, "priorityOverride", vzRsFiltAtt.PriorityOverride)
	A(vzRsFiltAttMap, "annotation", vzRsFiltAtt.Annotation)
	A(vzRsFiltAttMap, "tDn", vzRsFiltAtt.TDn)
	A(vzRsFiltAttMap, "tnVzFilterName", vzRsFiltAtt.TnVzFilterName)
	return vzRsFiltAttMap, err
}

func FilterRelationshipFromContainerList(cont *container.Container, index int) *FilterRelationship {
	FilterCont := cont.S("imdata").Index(index).S(VzrsfiltattClassName, "attributes")
	return &FilterRelationship{
		BaseAttributes{
			DistinguishedName: G(FilterCont, "dn"),
			Status:            G(FilterCont, "status"),
			ClassName:         VzrsfiltattClassName,
			Rn:                G(FilterCont, "rn"),
		},
		FilterRelationshipAttributes{
			Action:           G(FilterCont, "action"),
			Directives:       G(FilterCont, "directives"),
			PriorityOverride: G(FilterCont, "priorityOverride"),
			TDn:              G(FilterCont, "tDn"),
			TnVzFilterName:   G(FilterCont, "tnVzFilterName"),
		},
	}
}

func FilterRelationshipFromContainer(cont *container.Container) *FilterRelationship {
	return FilterRelationshipFromContainerList(cont, 0)
}

func FilterRelationshipListFromContainer(cont *container.Container) []*FilterRelationship {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*FilterRelationship, length)

	for i := 0; i < length; i++ {
		arr[i] = FilterRelationshipFromContainerList(cont, i)
	}

	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnvzRsSubjFiltAtt        = "uni/tn-%s/brc-%s/subj-%s/rssubjFiltAtt-%s"
	RnvzRsSubjFiltAtt        = "rssubjFiltAtt-%s"
	ParentDnvzRsSubjFiltAtt  = "uni/tn-%s/brc-%s/subj-%s"
	VzrssubjfiltattClassName = "vzRsSubjFiltAtt"
)

type SubjectFilter struct {
	BaseAttributes
	SubjectFilterAttributes
}

type SubjectFilterAttributes struct {
	Annotation       string `json:",omitempty"`
	Action           string `json:",omitempty"`
	Directives       string `json:",omitempty"`
	PriorityOverride string `json:",omitempty"`
	TDn              string `json:",omitempty"`
	TnVzFilterName   string `json:",omitempty"`
}

func NewSubjectFilter(vzRsSubjFiltAttRn, parentDn string, vzRsSubjFiltAttAttr SubjectFilterAttributes) *SubjectFilter {
	dn := fmt.Sprintf("%s/%s", parentDn, vzRsSubjFiltAttRn)
	return &SubjectFilter{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         VzrssubjfiltattClassName,
			Rn:                vzRsSubjFiltAttRn,
		},
		SubjectFilterAttributes: vzRsSubjFiltAttAttr,
	}
}

func (vzRsSubjFiltAtt *SubjectFilter) ToMap() (map[string]string, error) {
	vzRsSubjFiltAttMap, err := vzRsSubjFiltAtt.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vzRsSubjFiltAttMap, "action", vzRsSubjFiltAtt.Action)
	A(vzRsSubjFiltAttMap, "directives", vzRsSubjFiltAtt.Directives)
	A(vzRsSubjFiltAttMap, "priorityOverride", vzRsSubjFiltAtt.PriorityOverride)
	A(vzRsSubjFiltAttMap, "annotation", vzRsSubjFiltAtt.Annotation)
	A(vzRsSubjFiltAttMap, "tDn", vzRsSubjFiltAtt.TDn)
	A(vzRsSubjFiltAttMap, "tnVzFilterName", vzRsSubjFiltAtt.TnVzFilterName)
	return vzRsSubjFiltAttMap, err
}

func SubjectFilterFromContainerList(cont *container.Container, index int) *SubjectFilter {
	SubjectFilterCont := cont.S("imdata").Index(index).S(VzrssubjfiltattClassName, "attributes")
	return &SubjectFilter{
		BaseAttributes{
			DistinguishedName: G(SubjectFilterCont, "dn"),
			Status:            G(SubjectFilterCont, "status"),
			ClassName:         VzrssubjfiltattClassName,
			Rn:                G(SubjectFilterCont, "rn"),
		},
		SubjectFilterAttributes{
			Action:           G(SubjectFilterCont, "action"),
			Directives:       G(SubjectFilterCont, "directives"),
			PriorityOverride: G(SubjectFilterCont, "priorityOverride"),
			TDn:              G(SubjectFilterCont, "tDn"),
			TnVzFilterName:   G(SubjectFilterCont, "tnVzFilterName"),
		},
	}
}

func SubjectFilterFromContainer(cont *container.Container) *SubjectFilter {
	return SubjectFilterFromContainerList(cont, 0)
}

func SubjectFilterListFromContainer(cont *container.Container) []*SubjectFilter {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*SubjectFilter, length)

	for i := 0; i < length; i++ {
		arr[i] = SubjectFilterFromContainerList(cont, i)
	}

	return arr
}

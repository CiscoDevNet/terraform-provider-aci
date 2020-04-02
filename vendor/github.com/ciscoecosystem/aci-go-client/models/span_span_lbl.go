package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const SpanspanlblClassName = "spanSpanLbl"

type SPANSourcedestinationGroupMatchLabel struct {
	BaseAttributes
	SPANSourcedestinationGroupMatchLabelAttributes
}

type SPANSourcedestinationGroupMatchLabelAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Tag string `json:",omitempty"`
}

func NewSPANSourcedestinationGroupMatchLabel(spanSpanLblRn, parentDn, description string, spanSpanLblattr SPANSourcedestinationGroupMatchLabelAttributes) *SPANSourcedestinationGroupMatchLabel {
	dn := fmt.Sprintf("%s/%s", parentDn, spanSpanLblRn)
	return &SPANSourcedestinationGroupMatchLabel{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         SpanspanlblClassName,
			Rn:                spanSpanLblRn,
		},

		SPANSourcedestinationGroupMatchLabelAttributes: spanSpanLblattr,
	}
}

func (spanSpanLbl *SPANSourcedestinationGroupMatchLabel) ToMap() (map[string]string, error) {
	spanSpanLblMap, err := spanSpanLbl.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(spanSpanLblMap, "name", spanSpanLbl.Name)

	A(spanSpanLblMap, "annotation", spanSpanLbl.Annotation)

	A(spanSpanLblMap, "nameAlias", spanSpanLbl.NameAlias)

	A(spanSpanLblMap, "tag", spanSpanLbl.Tag)

	return spanSpanLblMap, err
}

func SPANSourcedestinationGroupMatchLabelFromContainerList(cont *container.Container, index int) *SPANSourcedestinationGroupMatchLabel {

	SPANSourcedestinationGroupMatchLabelCont := cont.S("imdata").Index(index).S(SpanspanlblClassName, "attributes")
	return &SPANSourcedestinationGroupMatchLabel{
		BaseAttributes{
			DistinguishedName: G(SPANSourcedestinationGroupMatchLabelCont, "dn"),
			Description:       G(SPANSourcedestinationGroupMatchLabelCont, "descr"),
			Status:            G(SPANSourcedestinationGroupMatchLabelCont, "status"),
			ClassName:         SpanspanlblClassName,
			Rn:                G(SPANSourcedestinationGroupMatchLabelCont, "rn"),
		},

		SPANSourcedestinationGroupMatchLabelAttributes{

			Name: G(SPANSourcedestinationGroupMatchLabelCont, "name"),

			Annotation: G(SPANSourcedestinationGroupMatchLabelCont, "annotation"),

			NameAlias: G(SPANSourcedestinationGroupMatchLabelCont, "nameAlias"),

			Tag: G(SPANSourcedestinationGroupMatchLabelCont, "tag"),
		},
	}
}

func SPANSourcedestinationGroupMatchLabelFromContainer(cont *container.Container) *SPANSourcedestinationGroupMatchLabel {

	return SPANSourcedestinationGroupMatchLabelFromContainerList(cont, 0)
}

func SPANSourcedestinationGroupMatchLabelListFromContainer(cont *container.Container) []*SPANSourcedestinationGroupMatchLabel {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*SPANSourcedestinationGroupMatchLabel, length)

	for i := 0; i < length; i++ {

		arr[i] = SPANSourcedestinationGroupMatchLabelFromContainerList(cont, i)
	}

	return arr
}

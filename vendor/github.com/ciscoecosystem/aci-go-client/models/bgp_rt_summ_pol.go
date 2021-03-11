package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BgprtsummpolClassName = "bgpRtSummPol"

type BgpRouteSummarization struct {
	BaseAttributes
	BgpRouteSummarizationAttributes
}

type BgpRouteSummarizationAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Attrmap string `json:",omitempty"`

	Ctrl string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewBgpRouteSummarization(bgpRtSummPolRn, parentDn, description string, bgpRtSummPolattr BgpRouteSummarizationAttributes) *BgpRouteSummarization {
	dn := fmt.Sprintf("%s/%s", parentDn, bgpRtSummPolRn)
	return &BgpRouteSummarization{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BgprtsummpolClassName,
			Rn:                bgpRtSummPolRn,
		},

		BgpRouteSummarizationAttributes: bgpRtSummPolattr,
	}
}

func (bgpRtSummPol *BgpRouteSummarization) ToMap() (map[string]string, error) {
	bgpRtSummPolMap, err := bgpRtSummPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bgpRtSummPolMap, "name", bgpRtSummPol.Name)

	A(bgpRtSummPolMap, "annotation", bgpRtSummPol.Annotation)

	A(bgpRtSummPolMap, "attrmap", bgpRtSummPol.Attrmap)

	A(bgpRtSummPolMap, "ctrl", bgpRtSummPol.Ctrl)

	A(bgpRtSummPolMap, "nameAlias", bgpRtSummPol.NameAlias)

	return bgpRtSummPolMap, err
}

func BgpRouteSummarizationFromContainerList(cont *container.Container, index int) *BgpRouteSummarization {

	BgpRouteSummarizationCont := cont.S("imdata").Index(index).S(BgprtsummpolClassName, "attributes")
	return &BgpRouteSummarization{
		BaseAttributes{
			DistinguishedName: G(BgpRouteSummarizationCont, "dn"),
			Description:       G(BgpRouteSummarizationCont, "descr"),
			Status:            G(BgpRouteSummarizationCont, "status"),
			ClassName:         BgprtsummpolClassName,
			Rn:                G(BgpRouteSummarizationCont, "rn"),
		},

		BgpRouteSummarizationAttributes{

			Name: G(BgpRouteSummarizationCont, "name"),

			Annotation: G(BgpRouteSummarizationCont, "annotation"),

			Attrmap: G(BgpRouteSummarizationCont, "attrmap"),

			Ctrl: G(BgpRouteSummarizationCont, "ctrl"),

			NameAlias: G(BgpRouteSummarizationCont, "nameAlias"),
		},
	}
}

func BgpRouteSummarizationFromContainer(cont *container.Container) *BgpRouteSummarization {

	return BgpRouteSummarizationFromContainerList(cont, 0)
}

func BgpRouteSummarizationListFromContainer(cont *container.Container) []*BgpRouteSummarization {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*BgpRouteSummarization, length)

	for i := 0; i < length; i++ {

		arr[i] = BgpRouteSummarizationFromContainerList(cont, i)
	}

	return arr
}

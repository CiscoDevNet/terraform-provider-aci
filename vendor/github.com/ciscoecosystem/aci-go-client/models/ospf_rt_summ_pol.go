package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const OspfrtsummpolClassName = "ospfRtSummPol"

type OspfRouteSummarization struct {
	BaseAttributes
	OspfRouteSummarizationAttributes
}

type OspfRouteSummarizationAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Cost string `json:",omitempty"`

	InterAreaEnabled string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Tag string `json:",omitempty"`
}

func NewOspfRouteSummarization(ospfRtSummPolRn, parentDn, description string, ospfRtSummPolattr OspfRouteSummarizationAttributes) *OspfRouteSummarization {
	dn := fmt.Sprintf("%s/%s", parentDn, ospfRtSummPolRn)
	return &OspfRouteSummarization{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         OspfrtsummpolClassName,
			Rn:                ospfRtSummPolRn,
		},

		OspfRouteSummarizationAttributes: ospfRtSummPolattr,
	}
}

func (ospfRtSummPol *OspfRouteSummarization) ToMap() (map[string]string, error) {
	ospfRtSummPolMap, err := ospfRtSummPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(ospfRtSummPolMap, "name", ospfRtSummPol.Name)

	A(ospfRtSummPolMap, "annotation", ospfRtSummPol.Annotation)

	A(ospfRtSummPolMap, "cost", ospfRtSummPol.Cost)

	A(ospfRtSummPolMap, "interAreaEnabled", ospfRtSummPol.InterAreaEnabled)

	A(ospfRtSummPolMap, "nameAlias", ospfRtSummPol.NameAlias)

	A(ospfRtSummPolMap, "tag", ospfRtSummPol.Tag)

	return ospfRtSummPolMap, err
}

func OspfRouteSummarizationFromContainerList(cont *container.Container, index int) *OspfRouteSummarization {

	OspfRouteSummarizationCont := cont.S("imdata").Index(index).S(OspfrtsummpolClassName, "attributes")
	return &OspfRouteSummarization{
		BaseAttributes{
			DistinguishedName: G(OspfRouteSummarizationCont, "dn"),
			Description:       G(OspfRouteSummarizationCont, "descr"),
			Status:            G(OspfRouteSummarizationCont, "status"),
			ClassName:         OspfrtsummpolClassName,
			Rn:                G(OspfRouteSummarizationCont, "rn"),
		},

		OspfRouteSummarizationAttributes{

			Name: G(OspfRouteSummarizationCont, "name"),

			Annotation: G(OspfRouteSummarizationCont, "annotation"),

			Cost: G(OspfRouteSummarizationCont, "cost"),

			InterAreaEnabled: G(OspfRouteSummarizationCont, "interAreaEnabled"),

			NameAlias: G(OspfRouteSummarizationCont, "nameAlias"),

			Tag: G(OspfRouteSummarizationCont, "tag"),
		},
	}
}

func OspfRouteSummarizationFromContainer(cont *container.Container) *OspfRouteSummarization {

	return OspfRouteSummarizationFromContainerList(cont, 0)
}

func OspfRouteSummarizationListFromContainer(cont *container.Container) []*OspfRouteSummarization {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*OspfRouteSummarization, length)

	for i := 0; i < length; i++ {

		arr[i] = OspfRouteSummarizationFromContainerList(cont, i)
	}

	return arr
}

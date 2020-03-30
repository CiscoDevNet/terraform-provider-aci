package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const SpandestgrpClassName = "spanDestGrp"

type SPANDestinationGroup struct {
	BaseAttributes
	SPANDestinationGroupAttributes
}

type SPANDestinationGroupAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewSPANDestinationGroup(spanDestGrpRn, parentDn, description string, spanDestGrpattr SPANDestinationGroupAttributes) *SPANDestinationGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, spanDestGrpRn)
	return &SPANDestinationGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         SpandestgrpClassName,
			Rn:                spanDestGrpRn,
		},

		SPANDestinationGroupAttributes: spanDestGrpattr,
	}
}

func (spanDestGrp *SPANDestinationGroup) ToMap() (map[string]string, error) {
	spanDestGrpMap, err := spanDestGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(spanDestGrpMap, "name", spanDestGrp.Name)

	A(spanDestGrpMap, "annotation", spanDestGrp.Annotation)

	A(spanDestGrpMap, "nameAlias", spanDestGrp.NameAlias)

	return spanDestGrpMap, err
}

func SPANDestinationGroupFromContainerList(cont *container.Container, index int) *SPANDestinationGroup {

	SPANDestinationGroupCont := cont.S("imdata").Index(index).S(SpandestgrpClassName, "attributes")
	return &SPANDestinationGroup{
		BaseAttributes{
			DistinguishedName: G(SPANDestinationGroupCont, "dn"),
			Description:       G(SPANDestinationGroupCont, "descr"),
			Status:            G(SPANDestinationGroupCont, "status"),
			ClassName:         SpandestgrpClassName,
			Rn:                G(SPANDestinationGroupCont, "rn"),
		},

		SPANDestinationGroupAttributes{

			Name: G(SPANDestinationGroupCont, "name"),

			Annotation: G(SPANDestinationGroupCont, "annotation"),

			NameAlias: G(SPANDestinationGroupCont, "nameAlias"),
		},
	}
}

func SPANDestinationGroupFromContainer(cont *container.Container) *SPANDestinationGroup {

	return SPANDestinationGroupFromContainerList(cont, 0)
}

func SPANDestinationGroupListFromContainer(cont *container.Container) []*SPANDestinationGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*SPANDestinationGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = SPANDestinationGroupFromContainerList(cont, i)
	}

	return arr
}

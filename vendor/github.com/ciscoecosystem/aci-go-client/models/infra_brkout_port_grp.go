package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfrabrkoutportgrpClassName = "infraBrkoutPortGrp"

type LeafBreakoutPortGroup struct {
	BaseAttributes
	LeafBreakoutPortGroupAttributes
}

type LeafBreakoutPortGroupAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	BrkoutMap string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewLeafBreakoutPortGroup(infraBrkoutPortGrpRn, parentDn, description string, infraBrkoutPortGrpattr LeafBreakoutPortGroupAttributes) *LeafBreakoutPortGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, infraBrkoutPortGrpRn)
	return &LeafBreakoutPortGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfrabrkoutportgrpClassName,
			Rn:                infraBrkoutPortGrpRn,
		},

		LeafBreakoutPortGroupAttributes: infraBrkoutPortGrpattr,
	}
}

func (infraBrkoutPortGrp *LeafBreakoutPortGroup) ToMap() (map[string]string, error) {
	infraBrkoutPortGrpMap, err := infraBrkoutPortGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraBrkoutPortGrpMap, "name", infraBrkoutPortGrp.Name)

	A(infraBrkoutPortGrpMap, "annotation", infraBrkoutPortGrp.Annotation)

	A(infraBrkoutPortGrpMap, "brkoutMap", infraBrkoutPortGrp.BrkoutMap)

	A(infraBrkoutPortGrpMap, "nameAlias", infraBrkoutPortGrp.NameAlias)

	return infraBrkoutPortGrpMap, err
}

func LeafBreakoutPortGroupFromContainerList(cont *container.Container, index int) *LeafBreakoutPortGroup {

	LeafBreakoutPortGroupCont := cont.S("imdata").Index(index).S(InfrabrkoutportgrpClassName, "attributes")
	return &LeafBreakoutPortGroup{
		BaseAttributes{
			DistinguishedName: G(LeafBreakoutPortGroupCont, "dn"),
			Description:       G(LeafBreakoutPortGroupCont, "descr"),
			Status:            G(LeafBreakoutPortGroupCont, "status"),
			ClassName:         InfrabrkoutportgrpClassName,
			Rn:                G(LeafBreakoutPortGroupCont, "rn"),
		},

		LeafBreakoutPortGroupAttributes{

			Name: G(LeafBreakoutPortGroupCont, "name"),

			Annotation: G(LeafBreakoutPortGroupCont, "annotation"),

			BrkoutMap: G(LeafBreakoutPortGroupCont, "brkoutMap"),

			NameAlias: G(LeafBreakoutPortGroupCont, "nameAlias"),
		},
	}
}

func LeafBreakoutPortGroupFromContainer(cont *container.Container) *LeafBreakoutPortGroup {

	return LeafBreakoutPortGroupFromContainerList(cont, 0)
}

func LeafBreakoutPortGroupListFromContainer(cont *container.Container) []*LeafBreakoutPortGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LeafBreakoutPortGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = LeafBreakoutPortGroupFromContainerList(cont, i)
	}

	return arr
}

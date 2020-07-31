package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfrafexbndlgrpClassName = "infraFexBndlGrp"

type FexBundleGroup struct {
	BaseAttributes
	FexBundleGroupAttributes
}

type FexBundleGroupAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewFexBundleGroup(infraFexBndlGrpRn, parentDn, description string, infraFexBndlGrpattr FexBundleGroupAttributes) *FexBundleGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, infraFexBndlGrpRn)
	return &FexBundleGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfrafexbndlgrpClassName,
			Rn:                infraFexBndlGrpRn,
		},

		FexBundleGroupAttributes: infraFexBndlGrpattr,
	}
}

func (infraFexBndlGrp *FexBundleGroup) ToMap() (map[string]string, error) {
	infraFexBndlGrpMap, err := infraFexBndlGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraFexBndlGrpMap, "name", infraFexBndlGrp.Name)

	A(infraFexBndlGrpMap, "annotation", infraFexBndlGrp.Annotation)

	A(infraFexBndlGrpMap, "nameAlias", infraFexBndlGrp.NameAlias)

	return infraFexBndlGrpMap, err
}

func FexBundleGroupFromContainerList(cont *container.Container, index int) *FexBundleGroup {

	FexBundleGroupCont := cont.S("imdata").Index(index).S(InfrafexbndlgrpClassName, "attributes")
	return &FexBundleGroup{
		BaseAttributes{
			DistinguishedName: G(FexBundleGroupCont, "dn"),
			Description:       G(FexBundleGroupCont, "descr"),
			Status:            G(FexBundleGroupCont, "status"),
			ClassName:         InfrafexbndlgrpClassName,
			Rn:                G(FexBundleGroupCont, "rn"),
		},

		FexBundleGroupAttributes{

			Name: G(FexBundleGroupCont, "name"),

			Annotation: G(FexBundleGroupCont, "annotation"),

			NameAlias: G(FexBundleGroupCont, "nameAlias"),
		},
	}
}

func FexBundleGroupFromContainer(cont *container.Container) *FexBundleGroup {

	return FexBundleGroupFromContainerList(cont, 0)
}

func FexBundleGroupListFromContainer(cont *container.Container) []*FexBundleGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FexBundleGroup, length)

	for i := 0; i < length; i++ {

		arr[i] = FexBundleGroupFromContainerList(cont, i)
	}

	return arr
}

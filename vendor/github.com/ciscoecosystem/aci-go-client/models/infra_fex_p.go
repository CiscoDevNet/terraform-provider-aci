package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfrafexpClassName = "infraFexP"

type FEXProfile struct {
	BaseAttributes
	FEXProfileAttributes
}

type FEXProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewFEXProfile(infraFexPRn, parentDn, description string, infraFexPattr FEXProfileAttributes) *FEXProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, infraFexPRn)
	return &FEXProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfrafexpClassName,
			Rn:                infraFexPRn,
		},

		FEXProfileAttributes: infraFexPattr,
	}
}

func (infraFexP *FEXProfile) ToMap() (map[string]string, error) {
	infraFexPMap, err := infraFexP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraFexPMap, "name", infraFexP.Name)

	A(infraFexPMap, "annotation", infraFexP.Annotation)

	A(infraFexPMap, "nameAlias", infraFexP.NameAlias)

	return infraFexPMap, err
}

func FEXProfileFromContainerList(cont *container.Container, index int) *FEXProfile {

	FEXProfileCont := cont.S("imdata").Index(index).S(InfrafexpClassName, "attributes")
	return &FEXProfile{
		BaseAttributes{
			DistinguishedName: G(FEXProfileCont, "dn"),
			Description:       G(FEXProfileCont, "descr"),
			Status:            G(FEXProfileCont, "status"),
			ClassName:         InfrafexpClassName,
			Rn:                G(FEXProfileCont, "rn"),
		},

		FEXProfileAttributes{

			Name: G(FEXProfileCont, "name"),

			Annotation: G(FEXProfileCont, "annotation"),

			NameAlias: G(FEXProfileCont, "nameAlias"),
		},
	}
}

func FEXProfileFromContainer(cont *container.Container) *FEXProfile {

	return FEXProfileFromContainerList(cont, 0)
}

func FEXProfileListFromContainer(cont *container.Container) []*FEXProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FEXProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = FEXProfileFromContainerList(cont, i)
	}

	return arr
}

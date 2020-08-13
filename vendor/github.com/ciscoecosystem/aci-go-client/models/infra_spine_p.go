package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfraspinepClassName = "infraSpineP"

type SpineProfile struct {
	BaseAttributes
	SpineProfileAttributes
}

type SpineProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewSpineProfile(infraSpinePRn, parentDn, description string, infraSpinePattr SpineProfileAttributes) *SpineProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, infraSpinePRn)
	return &SpineProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraspinepClassName,
			Rn:                infraSpinePRn,
		},

		SpineProfileAttributes: infraSpinePattr,
	}
}

func (infraSpineP *SpineProfile) ToMap() (map[string]string, error) {
	infraSpinePMap, err := infraSpineP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraSpinePMap, "name", infraSpineP.Name)

	A(infraSpinePMap, "annotation", infraSpineP.Annotation)

	A(infraSpinePMap, "nameAlias", infraSpineP.NameAlias)

	return infraSpinePMap, err
}

func SpineProfileFromContainerList(cont *container.Container, index int) *SpineProfile {

	SpineProfileCont := cont.S("imdata").Index(index).S(InfraspinepClassName, "attributes")
	return &SpineProfile{
		BaseAttributes{
			DistinguishedName: G(SpineProfileCont, "dn"),
			Description:       G(SpineProfileCont, "descr"),
			Status:            G(SpineProfileCont, "status"),
			ClassName:         InfraspinepClassName,
			Rn:                G(SpineProfileCont, "rn"),
		},

		SpineProfileAttributes{

			Name: G(SpineProfileCont, "name"),

			Annotation: G(SpineProfileCont, "annotation"),

			NameAlias: G(SpineProfileCont, "nameAlias"),
		},
	}
}

func SpineProfileFromContainer(cont *container.Container) *SpineProfile {

	return SpineProfileFromContainerList(cont, 0)
}

func SpineProfileListFromContainer(cont *container.Container) []*SpineProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*SpineProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = SpineProfileFromContainerList(cont, i)
	}

	return arr
}

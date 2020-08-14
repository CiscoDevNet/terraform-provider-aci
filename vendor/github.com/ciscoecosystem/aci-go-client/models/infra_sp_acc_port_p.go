package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfraspaccportpClassName = "infraSpAccPortP"

type SpineInterfaceProfile struct {
	BaseAttributes
	SpineInterfaceProfileAttributes
}

type SpineInterfaceProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewSpineInterfaceProfile(infraSpAccPortPRn, parentDn, description string, infraSpAccPortPattr SpineInterfaceProfileAttributes) *SpineInterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, infraSpAccPortPRn)
	return &SpineInterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraspaccportpClassName,
			Rn:                infraSpAccPortPRn,
		},

		SpineInterfaceProfileAttributes: infraSpAccPortPattr,
	}
}

func (infraSpAccPortP *SpineInterfaceProfile) ToMap() (map[string]string, error) {
	infraSpAccPortPMap, err := infraSpAccPortP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraSpAccPortPMap, "name", infraSpAccPortP.Name)

	A(infraSpAccPortPMap, "annotation", infraSpAccPortP.Annotation)

	A(infraSpAccPortPMap, "nameAlias", infraSpAccPortP.NameAlias)

	return infraSpAccPortPMap, err
}

func SpineInterfaceProfileFromContainerList(cont *container.Container, index int) *SpineInterfaceProfile {

	SpineInterfaceProfileCont := cont.S("imdata").Index(index).S(InfraspaccportpClassName, "attributes")
	return &SpineInterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(SpineInterfaceProfileCont, "dn"),
			Description:       G(SpineInterfaceProfileCont, "descr"),
			Status:            G(SpineInterfaceProfileCont, "status"),
			ClassName:         InfraspaccportpClassName,
			Rn:                G(SpineInterfaceProfileCont, "rn"),
		},

		SpineInterfaceProfileAttributes{

			Name: G(SpineInterfaceProfileCont, "name"),

			Annotation: G(SpineInterfaceProfileCont, "annotation"),

			NameAlias: G(SpineInterfaceProfileCont, "nameAlias"),
		},
	}
}

func SpineInterfaceProfileFromContainer(cont *container.Container) *SpineInterfaceProfile {

	return SpineInterfaceProfileFromContainerList(cont, 0)
}

func SpineInterfaceProfileListFromContainer(cont *container.Container) []*SpineInterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*SpineInterfaceProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = SpineInterfaceProfileFromContainerList(cont, i)
	}

	return arr
}

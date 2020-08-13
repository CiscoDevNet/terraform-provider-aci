package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfrarsspaccportpClassName = "infraRsSpAccPortP"

type InterfaceProfile struct {
	BaseAttributes
	InterfaceProfileAttributes
}

type InterfaceProfileAttributes struct {
	TDn string `json:",omitempty"`

	Annotation string `json:",omitempty"`
}

func NewInterfaceProfile(infraRsSpAccPortPRn, parentDn, description string, infraRsSpAccPortPattr InterfaceProfileAttributes) *InterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, infraRsSpAccPortPRn)
	return &InterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfrarsspaccportpClassName,
			Rn:                infraRsSpAccPortPRn,
		},

		InterfaceProfileAttributes: infraRsSpAccPortPattr,
	}
}

func (infraRsSpAccPortP *InterfaceProfile) ToMap() (map[string]string, error) {
	infraRsSpAccPortPMap, err := infraRsSpAccPortP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraRsSpAccPortPMap, "tDn", infraRsSpAccPortP.TDn)

	A(infraRsSpAccPortPMap, "annotation", infraRsSpAccPortP.Annotation)

	return infraRsSpAccPortPMap, err
}

func InterfaceProfileFromContainerList(cont *container.Container, index int) *InterfaceProfile {

	InterfaceProfileCont := cont.S("imdata").Index(index).S(InfrarsspaccportpClassName, "attributes")
	return &InterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(InterfaceProfileCont, "dn"),
			Description:       G(InterfaceProfileCont, "descr"),
			Status:            G(InterfaceProfileCont, "status"),
			ClassName:         InfrarsspaccportpClassName,
			Rn:                G(InterfaceProfileCont, "rn"),
		},

		InterfaceProfileAttributes{

			TDn: G(InterfaceProfileCont, "tDn"),

			Annotation: G(InterfaceProfileCont, "annotation"),
		},
	}
}

func InterfaceProfileFromContainer(cont *container.Container) *InterfaceProfile {

	return InterfaceProfileFromContainerList(cont, 0)
}

func InterfaceProfileListFromContainer(cont *container.Container) []*InterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*InterfaceProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = InterfaceProfileFromContainerList(cont, i)
	}

	return arr
}

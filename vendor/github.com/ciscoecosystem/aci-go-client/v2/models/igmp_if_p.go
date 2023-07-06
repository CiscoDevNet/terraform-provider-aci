package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnIgmpIfP        = "igmpIfP"
	IgmpIfPClassName = "igmpIfP"
)

type IGMPInterfaceProfile struct {
	BaseAttributes
	IGMPInterfaceProfileAttributes
}

type IGMPInterfaceProfileAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewIGMPInterfaceProfile(parentDn, description string, igmpIfPAttr IGMPInterfaceProfileAttributes) *IGMPInterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, RnIgmpIfP)
	return &IGMPInterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         IgmpIfPClassName,
			Rn:                RnIgmpIfP,
		},
		IGMPInterfaceProfileAttributes: igmpIfPAttr,
	}
}

func (igmpIfP *IGMPInterfaceProfile) ToMap() (map[string]string, error) {
	igmpIfPMap, err := igmpIfP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(igmpIfPMap, "annotation", igmpIfP.Annotation)
	A(igmpIfPMap, "name", igmpIfP.Name)
	A(igmpIfPMap, "nameAlias", igmpIfP.NameAlias)
	return igmpIfPMap, err
}

func IGMPInterfaceProfileFromContainerList(cont *container.Container, index int) *IGMPInterfaceProfile {
	InterfaceProfileCont := cont.S("imdata").Index(index).S(IgmpIfPClassName, "attributes")
	return &IGMPInterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(InterfaceProfileCont, "dn"),
			Description:       G(InterfaceProfileCont, "descr"),
			Status:            G(InterfaceProfileCont, "status"),
			ClassName:         IgmpIfPClassName,
			Rn:                G(InterfaceProfileCont, "rn"),
		},
		IGMPInterfaceProfileAttributes{
			Annotation: G(InterfaceProfileCont, "annotation"),
			Name:       G(InterfaceProfileCont, "name"),
			NameAlias:  G(InterfaceProfileCont, "nameAlias"),
		},
	}
}

func IGMPInterfaceProfileFromContainer(cont *container.Container) *IGMPInterfaceProfile {
	return IGMPInterfaceProfileFromContainerList(cont, 0)
}

func IGMPInterfaceProfileListFromContainer(cont *container.Container) []*IGMPInterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*IGMPInterfaceProfile, length)

	for i := 0; i < length; i++ {
		arr[i] = IGMPInterfaceProfileFromContainerList(cont, i)
	}

	return arr
}

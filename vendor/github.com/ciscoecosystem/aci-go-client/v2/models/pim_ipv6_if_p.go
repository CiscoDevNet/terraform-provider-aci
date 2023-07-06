package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnPimIPV6IfP        = "pimipv6ifp"
	DnPimIPV6IfP        = "uni/tn-%s/out-%s/lnodep-%s/lifp-%s/pimipv6ifp"
	ParentDnPimIPV6IfP  = "uni/tn-%s/out-%s/lnodep-%s/lifp-%s"
	PimIPV6IfPClassName = "pimIPV6IfP"
)

type PIMIPv6InterfaceProfile struct {
	BaseAttributes
	PIMIPv6InterfaceProfileAttributes
}

type PIMIPv6InterfaceProfileAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewPIMIPv6InterfaceProfile(parentDn, description string, pimIPV6IfPAttr PIMIPv6InterfaceProfileAttributes) *PIMIPv6InterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, RnPimIPV6IfP)
	return &PIMIPv6InterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PimIPV6IfPClassName,
			Rn:                RnPimIPV6IfP,
		},
		PIMIPv6InterfaceProfileAttributes: pimIPV6IfPAttr,
	}
}

func (pimIPV6IfP *PIMIPv6InterfaceProfile) ToMap() (map[string]string, error) {
	pimIPV6IfPMap, err := pimIPV6IfP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(pimIPV6IfPMap, "annotation", pimIPV6IfP.Annotation)
	A(pimIPV6IfPMap, "name", pimIPV6IfP.Name)
	A(pimIPV6IfPMap, "nameAlias", pimIPV6IfP.NameAlias)
	return pimIPV6IfPMap, err
}

func PIMIPv6InterfaceProfileFromContainerList(cont *container.Container, index int) *PIMIPv6InterfaceProfile {
	InterfaceProfileCont := cont.S("imdata").Index(index).S(PimIPV6IfPClassName, "attributes")
	return &PIMIPv6InterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(InterfaceProfileCont, "dn"),
			Description:       G(InterfaceProfileCont, "descr"),
			Status:            G(InterfaceProfileCont, "status"),
			ClassName:         PimIPV6IfPClassName,
			Rn:                G(InterfaceProfileCont, "rn"),
		},
		PIMIPv6InterfaceProfileAttributes{
			Annotation: G(InterfaceProfileCont, "annotation"),
			Name:       G(InterfaceProfileCont, "name"),
			NameAlias:  G(InterfaceProfileCont, "nameAlias"),
		},
	}
}

func PIMIPv6InterfaceProfileFromContainer(cont *container.Container) *PIMIPv6InterfaceProfile {
	return PIMIPv6InterfaceProfileFromContainerList(cont, 0)
}

func PIMIPv6InterfaceProfileListFromContainer(cont *container.Container) []*PIMIPv6InterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PIMIPv6InterfaceProfile, length)

	for i := 0; i < length; i++ {
		arr[i] = PIMIPv6InterfaceProfileFromContainerList(cont, i)
	}

	return arr
}

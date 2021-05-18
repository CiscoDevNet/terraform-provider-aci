package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extloopbackifpClassName = "l3extLoopBackIfP"

type LoopBackInterfaceProfile struct {
	BaseAttributes
	LoopBackInterfaceProfileAttributes
}

type LoopBackInterfaceProfileAttributes struct {
	Addr string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewLoopBackInterfaceProfile(l3extLoopBackIfPRn, parentDn, description string, l3extLoopBackIfPattr LoopBackInterfaceProfileAttributes) *LoopBackInterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extLoopBackIfPRn)
	return &LoopBackInterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extloopbackifpClassName,
			Rn:                l3extLoopBackIfPRn,
		},

		LoopBackInterfaceProfileAttributes: l3extLoopBackIfPattr,
	}
}

func (l3extLoopBackIfP *LoopBackInterfaceProfile) ToMap() (map[string]string, error) {
	l3extLoopBackIfPMap, err := l3extLoopBackIfP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extLoopBackIfPMap, "addr", l3extLoopBackIfP.Addr)

	A(l3extLoopBackIfPMap, "annotation", l3extLoopBackIfP.Annotation)

	A(l3extLoopBackIfPMap, "nameAlias", l3extLoopBackIfP.NameAlias)

	return l3extLoopBackIfPMap, err
}

func LoopBackInterfaceProfileFromContainerList(cont *container.Container, index int) *LoopBackInterfaceProfile {

	LoopBackInterfaceProfileCont := cont.S("imdata").Index(index).S(L3extloopbackifpClassName, "attributes")
	return &LoopBackInterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(LoopBackInterfaceProfileCont, "dn"),
			Description:       G(LoopBackInterfaceProfileCont, "descr"),
			Status:            G(LoopBackInterfaceProfileCont, "status"),
			ClassName:         L3extloopbackifpClassName,
			Rn:                G(LoopBackInterfaceProfileCont, "rn"),
		},

		LoopBackInterfaceProfileAttributes{

			Addr: G(LoopBackInterfaceProfileCont, "addr"),

			Annotation: G(LoopBackInterfaceProfileCont, "annotation"),

			NameAlias: G(LoopBackInterfaceProfileCont, "nameAlias"),
		},
	}
}

func LoopBackInterfaceProfileFromContainer(cont *container.Container) *LoopBackInterfaceProfile {

	return LoopBackInterfaceProfileFromContainerList(cont, 0)
}

func LoopBackInterfaceProfileListFromContainer(cont *container.Container) []*LoopBackInterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LoopBackInterfaceProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = LoopBackInterfaceProfileFromContainerList(cont, i)
	}

	return arr
}

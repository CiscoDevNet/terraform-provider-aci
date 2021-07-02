package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extlifpClassName = "l3extLIfP"

type LogicalInterfaceProfile struct {
	BaseAttributes
	LogicalInterfaceProfileAttributes
}

type LogicalInterfaceProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Prio string `json:",omitempty"`

	Tag string `json:",omitempty"`
}

func NewLogicalInterfaceProfile(l3extLIfPRn, parentDn, description string, l3extLIfPattr LogicalInterfaceProfileAttributes) *LogicalInterfaceProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extLIfPRn)
	return &LogicalInterfaceProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extlifpClassName,
			Rn:                l3extLIfPRn,
		},

		LogicalInterfaceProfileAttributes: l3extLIfPattr,
	}
}

func (l3extLIfP *LogicalInterfaceProfile) ToMap() (map[string]string, error) {
	l3extLIfPMap, err := l3extLIfP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extLIfPMap, "name", l3extLIfP.Name)

	A(l3extLIfPMap, "annotation", l3extLIfP.Annotation)

	A(l3extLIfPMap, "nameAlias", l3extLIfP.NameAlias)

	A(l3extLIfPMap, "prio", l3extLIfP.Prio)

	A(l3extLIfPMap, "tag", l3extLIfP.Tag)

	return l3extLIfPMap, err
}

func LogicalInterfaceProfileFromContainerList(cont *container.Container, index int) *LogicalInterfaceProfile {

	LogicalInterfaceProfileCont := cont.S("imdata").Index(index).S(L3extlifpClassName, "attributes")
	return &LogicalInterfaceProfile{
		BaseAttributes{
			DistinguishedName: G(LogicalInterfaceProfileCont, "dn"),
			Description:       G(LogicalInterfaceProfileCont, "descr"),
			Status:            G(LogicalInterfaceProfileCont, "status"),
			ClassName:         L3extlifpClassName,
			Rn:                G(LogicalInterfaceProfileCont, "rn"),
		},

		LogicalInterfaceProfileAttributes{

			Name: G(LogicalInterfaceProfileCont, "name"),

			Annotation: G(LogicalInterfaceProfileCont, "annotation"),

			NameAlias: G(LogicalInterfaceProfileCont, "nameAlias"),

			Prio: G(LogicalInterfaceProfileCont, "prio"),

			Tag: G(LogicalInterfaceProfileCont, "tag"),
		},
	}
}

func LogicalInterfaceProfileFromContainer(cont *container.Container) *LogicalInterfaceProfile {

	return LogicalInterfaceProfileFromContainerList(cont, 0)
}

func LogicalInterfaceProfileListFromContainer(cont *container.Container) []*LogicalInterfaceProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LogicalInterfaceProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = LogicalInterfaceProfileFromContainerList(cont, i)
	}

	return arr
}

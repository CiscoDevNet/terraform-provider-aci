package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const RtctrlprofileClassName = "rtctrlProfile"

type RouteControlProfile struct {
	BaseAttributes
	RouteControlProfileAttributes
}

type RouteControlProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	RouteControlProfileType string `json:",omitempty"`
}

func NewRouteControlProfile(rtctrlProfileRn, parentDn, description string, rtctrlProfileattr RouteControlProfileAttributes) *RouteControlProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlProfileRn)
	return &RouteControlProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlprofileClassName,
			Rn:                rtctrlProfileRn,
		},

		RouteControlProfileAttributes: rtctrlProfileattr,
	}
}

func (rtctrlProfile *RouteControlProfile) ToMap() (map[string]string, error) {
	rtctrlProfileMap, err := rtctrlProfile.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(rtctrlProfileMap, "name", rtctrlProfile.Name)

	A(rtctrlProfileMap, "annotation", rtctrlProfile.Annotation)

	A(rtctrlProfileMap, "nameAlias", rtctrlProfile.NameAlias)

	A(rtctrlProfileMap, "type", rtctrlProfile.RouteControlProfileType)

	return rtctrlProfileMap, err
}

func RouteControlProfileFromContainerList(cont *container.Container, index int) *RouteControlProfile {

	RouteControlProfileCont := cont.S("imdata").Index(index).S(RtctrlprofileClassName, "attributes")
	return &RouteControlProfile{
		BaseAttributes{
			DistinguishedName: G(RouteControlProfileCont, "dn"),
			Description:       G(RouteControlProfileCont, "descr"),
			Status:            G(RouteControlProfileCont, "status"),
			ClassName:         RtctrlprofileClassName,
			Rn:                G(RouteControlProfileCont, "rn"),
		},

		RouteControlProfileAttributes{

			Name: G(RouteControlProfileCont, "name"),

			Annotation: G(RouteControlProfileCont, "annotation"),

			NameAlias: G(RouteControlProfileCont, "nameAlias"),

			RouteControlProfileType: G(RouteControlProfileCont, "type"),
		},
	}
}

func RouteControlProfileFromContainer(cont *container.Container) *RouteControlProfile {

	return RouteControlProfileFromContainerList(cont, 0)
}

func RouteControlProfileListFromContainer(cont *container.Container) []*RouteControlProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*RouteControlProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = RouteControlProfileFromContainerList(cont, i)
	}

	return arr
}

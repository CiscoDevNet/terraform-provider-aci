package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const HsrpgrouppClassName = "hsrpGroupP"

type HSRPGroupProfile struct {
	BaseAttributes
	HSRPGroupProfileAttributes
}

type HSRPGroupProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ConfigIssues string `json:",omitempty"`

	GroupAf string `json:",omitempty"`

	GroupId string `json:",omitempty"`

	GroupName string `json:",omitempty"`

	Ip string `json:",omitempty"`

	IpObtainMode string `json:",omitempty"`

	Mac string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewHSRPGroupProfile(hsrpGroupPRn, parentDn, description string, hsrpGroupPattr HSRPGroupProfileAttributes) *HSRPGroupProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, hsrpGroupPRn)
	return &HSRPGroupProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         HsrpgrouppClassName,
			Rn:                hsrpGroupPRn,
		},

		HSRPGroupProfileAttributes: hsrpGroupPattr,
	}
}

func (hsrpGroupP *HSRPGroupProfile) ToMap() (map[string]string, error) {
	hsrpGroupPMap, err := hsrpGroupP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(hsrpGroupPMap, "name", hsrpGroupP.Name)

	A(hsrpGroupPMap, "annotation", hsrpGroupP.Annotation)

	A(hsrpGroupPMap, "configIssues", hsrpGroupP.ConfigIssues)

	A(hsrpGroupPMap, "groupAf", hsrpGroupP.GroupAf)

	A(hsrpGroupPMap, "groupId", hsrpGroupP.GroupId)

	A(hsrpGroupPMap, "groupName", hsrpGroupP.GroupName)

	A(hsrpGroupPMap, "ip", hsrpGroupP.Ip)

	A(hsrpGroupPMap, "ipObtainMode", hsrpGroupP.IpObtainMode)

	A(hsrpGroupPMap, "mac", hsrpGroupP.Mac)

	A(hsrpGroupPMap, "nameAlias", hsrpGroupP.NameAlias)

	return hsrpGroupPMap, err
}

func HSRPGroupProfileFromContainerList(cont *container.Container, index int) *HSRPGroupProfile {

	HSRPGroupProfileCont := cont.S("imdata").Index(index).S(HsrpgrouppClassName, "attributes")
	return &HSRPGroupProfile{
		BaseAttributes{
			DistinguishedName: G(HSRPGroupProfileCont, "dn"),
			Description:       G(HSRPGroupProfileCont, "descr"),
			Status:            G(HSRPGroupProfileCont, "status"),
			ClassName:         HsrpgrouppClassName,
			Rn:                G(HSRPGroupProfileCont, "rn"),
		},

		HSRPGroupProfileAttributes{

			Name: G(HSRPGroupProfileCont, "name"),

			Annotation: G(HSRPGroupProfileCont, "annotation"),

			ConfigIssues: G(HSRPGroupProfileCont, "configIssues"),

			GroupAf: G(HSRPGroupProfileCont, "groupAf"),

			GroupId: G(HSRPGroupProfileCont, "groupId"),

			GroupName: G(HSRPGroupProfileCont, "groupName"),

			Ip: G(HSRPGroupProfileCont, "ip"),

			IpObtainMode: G(HSRPGroupProfileCont, "ipObtainMode"),

			Mac: G(HSRPGroupProfileCont, "mac"),

			NameAlias: G(HSRPGroupProfileCont, "nameAlias"),
		},
	}
}

func HSRPGroupProfileFromContainer(cont *container.Container) *HSRPGroupProfile {

	return HSRPGroupProfileFromContainerList(cont, 0)
}

func HSRPGroupProfileListFromContainer(cont *container.Container) []*HSRPGroupProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*HSRPGroupProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = HSRPGroupProfileFromContainerList(cont, i)
	}

	return arr
}

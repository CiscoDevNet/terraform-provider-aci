package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BgpctxpolClassName = "bgpCtxPol"

type BGPTimersPolicy struct {
	BaseAttributes
	BGPTimersPolicyAttributes
}

type BGPTimersPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	GrCtrl string `json:",omitempty"`

	HoldIntvl string `json:",omitempty"`

	KaIntvl string `json:",omitempty"`

	MaxAsLimit string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	StaleIntvl string `json:",omitempty"`
}

func NewBGPTimersPolicy(bgpCtxPolRn, parentDn, description string, bgpCtxPolattr BGPTimersPolicyAttributes) *BGPTimersPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, bgpCtxPolRn)
	return &BGPTimersPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BgpctxpolClassName,
			Rn:                bgpCtxPolRn,
		},

		BGPTimersPolicyAttributes: bgpCtxPolattr,
	}
}

func (bgpCtxPol *BGPTimersPolicy) ToMap() (map[string]string, error) {
	bgpCtxPolMap, err := bgpCtxPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bgpCtxPolMap, "name", bgpCtxPol.Name)

	A(bgpCtxPolMap, "annotation", bgpCtxPol.Annotation)

	A(bgpCtxPolMap, "grCtrl", bgpCtxPol.GrCtrl)

	A(bgpCtxPolMap, "holdIntvl", bgpCtxPol.HoldIntvl)

	A(bgpCtxPolMap, "kaIntvl", bgpCtxPol.KaIntvl)

	A(bgpCtxPolMap, "maxAsLimit", bgpCtxPol.MaxAsLimit)

	A(bgpCtxPolMap, "nameAlias", bgpCtxPol.NameAlias)

	A(bgpCtxPolMap, "staleIntvl", bgpCtxPol.StaleIntvl)

	return bgpCtxPolMap, err
}

func BGPTimersPolicyFromContainerList(cont *container.Container, index int) *BGPTimersPolicy {

	BGPTimersPolicyCont := cont.S("imdata").Index(index).S(BgpctxpolClassName, "attributes")
	return &BGPTimersPolicy{
		BaseAttributes{
			DistinguishedName: G(BGPTimersPolicyCont, "dn"),
			Description:       G(BGPTimersPolicyCont, "descr"),
			Status:            G(BGPTimersPolicyCont, "status"),
			ClassName:         BgpctxpolClassName,
			Rn:                G(BGPTimersPolicyCont, "rn"),
		},

		BGPTimersPolicyAttributes{

			Name: G(BGPTimersPolicyCont, "name"),

			Annotation: G(BGPTimersPolicyCont, "annotation"),

			GrCtrl: G(BGPTimersPolicyCont, "grCtrl"),

			HoldIntvl: G(BGPTimersPolicyCont, "holdIntvl"),

			KaIntvl: G(BGPTimersPolicyCont, "kaIntvl"),

			MaxAsLimit: G(BGPTimersPolicyCont, "maxAsLimit"),

			NameAlias: G(BGPTimersPolicyCont, "nameAlias"),

			StaleIntvl: G(BGPTimersPolicyCont, "staleIntvl"),
		},
	}
}

func BGPTimersPolicyFromContainer(cont *container.Container) *BGPTimersPolicy {

	return BGPTimersPolicyFromContainerList(cont, 0)
}

func BGPTimersPolicyListFromContainer(cont *container.Container) []*BGPTimersPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*BGPTimersPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = BGPTimersPolicyFromContainerList(cont, i)
	}

	return arr
}

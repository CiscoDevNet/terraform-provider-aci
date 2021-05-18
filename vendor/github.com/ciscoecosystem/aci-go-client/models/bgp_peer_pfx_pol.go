package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BgppeerpfxpolClassName = "bgpPeerPfxPol"

type BGPPeerPrefixPolicy struct {
	BaseAttributes
	BGPPeerPrefixPolicyAttributes
}

type BGPPeerPrefixPolicyAttributes struct {
	Name string `json:",omitempty"`

	Action string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	MaxPfx string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	RestartTime string `json:",omitempty"`

	Thresh string `json:",omitempty"`
}

func NewBGPPeerPrefixPolicy(bgpPeerPfxPolRn, parentDn, description string, bgpPeerPfxPolattr BGPPeerPrefixPolicyAttributes) *BGPPeerPrefixPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, bgpPeerPfxPolRn)
	return &BGPPeerPrefixPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BgppeerpfxpolClassName,
			Rn:                bgpPeerPfxPolRn,
		},

		BGPPeerPrefixPolicyAttributes: bgpPeerPfxPolattr,
	}
}

func (bgpPeerPfxPol *BGPPeerPrefixPolicy) ToMap() (map[string]string, error) {
	bgpPeerPfxPolMap, err := bgpPeerPfxPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bgpPeerPfxPolMap, "name", bgpPeerPfxPol.Name)

	A(bgpPeerPfxPolMap, "action", bgpPeerPfxPol.Action)

	A(bgpPeerPfxPolMap, "annotation", bgpPeerPfxPol.Annotation)

	A(bgpPeerPfxPolMap, "maxPfx", bgpPeerPfxPol.MaxPfx)

	A(bgpPeerPfxPolMap, "nameAlias", bgpPeerPfxPol.NameAlias)

	A(bgpPeerPfxPolMap, "restartTime", bgpPeerPfxPol.RestartTime)

	A(bgpPeerPfxPolMap, "thresh", bgpPeerPfxPol.Thresh)

	return bgpPeerPfxPolMap, err
}

func BGPPeerPrefixPolicyFromContainerList(cont *container.Container, index int) *BGPPeerPrefixPolicy {

	BGPPeerPrefixPolicyCont := cont.S("imdata").Index(index).S(BgppeerpfxpolClassName, "attributes")
	return &BGPPeerPrefixPolicy{
		BaseAttributes{
			DistinguishedName: G(BGPPeerPrefixPolicyCont, "dn"),
			Description:       G(BGPPeerPrefixPolicyCont, "descr"),
			Status:            G(BGPPeerPrefixPolicyCont, "status"),
			ClassName:         BgppeerpfxpolClassName,
			Rn:                G(BGPPeerPrefixPolicyCont, "rn"),
		},

		BGPPeerPrefixPolicyAttributes{

			Name: G(BGPPeerPrefixPolicyCont, "name"),

			Action: G(BGPPeerPrefixPolicyCont, "action"),

			Annotation: G(BGPPeerPrefixPolicyCont, "annotation"),

			MaxPfx: G(BGPPeerPrefixPolicyCont, "maxPfx"),

			NameAlias: G(BGPPeerPrefixPolicyCont, "nameAlias"),

			RestartTime: G(BGPPeerPrefixPolicyCont, "restartTime"),

			Thresh: G(BGPPeerPrefixPolicyCont, "thresh"),
		},
	}
}

func BGPPeerPrefixPolicyFromContainer(cont *container.Container) *BGPPeerPrefixPolicy {

	return BGPPeerPrefixPolicyFromContainerList(cont, 0)
}

func BGPPeerPrefixPolicyListFromContainer(cont *container.Container) []*BGPPeerPrefixPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*BGPPeerPrefixPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = BGPPeerPrefixPolicyFromContainerList(cont, i)
	}

	return arr
}

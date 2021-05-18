package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BgpctxafpolClassName = "bgpCtxAfPol"

type BGPAddressFamilyContextPolicy struct {
	BaseAttributes
	BGPAddressFamilyContextPolicyAttributes
}

type BGPAddressFamilyContextPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Ctrl string `json:",omitempty"`

	EDist string `json:",omitempty"`

	IDist string `json:",omitempty"`

	LocalDist string `json:",omitempty"`

	MaxEcmp string `json:",omitempty"`

	MaxEcmpIbgp string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewBGPAddressFamilyContextPolicy(bgpCtxAfPolRn, parentDn, description string, bgpCtxAfPolattr BGPAddressFamilyContextPolicyAttributes) *BGPAddressFamilyContextPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, bgpCtxAfPolRn)
	return &BGPAddressFamilyContextPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BgpctxafpolClassName,
			Rn:                bgpCtxAfPolRn,
		},

		BGPAddressFamilyContextPolicyAttributes: bgpCtxAfPolattr,
	}
}

func (bgpCtxAfPol *BGPAddressFamilyContextPolicy) ToMap() (map[string]string, error) {
	bgpCtxAfPolMap, err := bgpCtxAfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bgpCtxAfPolMap, "name", bgpCtxAfPol.Name)

	A(bgpCtxAfPolMap, "annotation", bgpCtxAfPol.Annotation)

	A(bgpCtxAfPolMap, "ctrl", bgpCtxAfPol.Ctrl)

	A(bgpCtxAfPolMap, "eDist", bgpCtxAfPol.EDist)

	A(bgpCtxAfPolMap, "iDist", bgpCtxAfPol.IDist)

	A(bgpCtxAfPolMap, "localDist", bgpCtxAfPol.LocalDist)

	A(bgpCtxAfPolMap, "maxEcmp", bgpCtxAfPol.MaxEcmp)

	A(bgpCtxAfPolMap, "maxEcmpIbgp", bgpCtxAfPol.MaxEcmpIbgp)

	A(bgpCtxAfPolMap, "nameAlias", bgpCtxAfPol.NameAlias)

	return bgpCtxAfPolMap, err
}

func BGPAddressFamilyContextPolicyFromContainerList(cont *container.Container, index int) *BGPAddressFamilyContextPolicy {

	BGPAddressFamilyContextPolicyCont := cont.S("imdata").Index(index).S(BgpctxafpolClassName, "attributes")
	return &BGPAddressFamilyContextPolicy{
		BaseAttributes{
			DistinguishedName: G(BGPAddressFamilyContextPolicyCont, "dn"),
			Description:       G(BGPAddressFamilyContextPolicyCont, "descr"),
			Status:            G(BGPAddressFamilyContextPolicyCont, "status"),
			ClassName:         BgpctxafpolClassName,
			Rn:                G(BGPAddressFamilyContextPolicyCont, "rn"),
		},

		BGPAddressFamilyContextPolicyAttributes{

			Name: G(BGPAddressFamilyContextPolicyCont, "name"),

			Annotation: G(BGPAddressFamilyContextPolicyCont, "annotation"),

			Ctrl: G(BGPAddressFamilyContextPolicyCont, "ctrl"),

			EDist: G(BGPAddressFamilyContextPolicyCont, "eDist"),

			IDist: G(BGPAddressFamilyContextPolicyCont, "iDist"),

			LocalDist: G(BGPAddressFamilyContextPolicyCont, "localDist"),

			MaxEcmp: G(BGPAddressFamilyContextPolicyCont, "maxEcmp"),

			MaxEcmpIbgp: G(BGPAddressFamilyContextPolicyCont, "maxEcmpIbgp"),

			NameAlias: G(BGPAddressFamilyContextPolicyCont, "nameAlias"),
		},
	}
}

func BGPAddressFamilyContextPolicyFromContainer(cont *container.Container) *BGPAddressFamilyContextPolicy {

	return BGPAddressFamilyContextPolicyFromContainerList(cont, 0)
}

func BGPAddressFamilyContextPolicyListFromContainer(cont *container.Container) []*BGPAddressFamilyContextPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*BGPAddressFamilyContextPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = BGPAddressFamilyContextPolicyFromContainerList(cont, i)
	}

	return arr
}

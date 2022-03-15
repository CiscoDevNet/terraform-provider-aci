package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnfvRsCtxToBgpCtxAfPol        = "uni/tn-%s/ctx-%s/rsctxToBgpCtxAfPol-[%s]-%s"
	RnfvRsCtxToBgpCtxAfPol        = "rsctxToBgpCtxAfPol-[%s]-%s"
	ParentDnfvRsCtxToBgpCtxAfPol  = "uni/tn-%s/ctx-%s"
	FvrsctxtobgpctxafpolClassName = "fvRsCtxToBgpCtxAfPol"
)

type BGPAddressFamilyContextPolicyRelationship struct {
	BaseAttributes
	BGPAddressFamilyContextPolicyRelationshipAttributes
}

type BGPAddressFamilyContextPolicyRelationshipAttributes struct {
	Annotation        string `json:",omitempty"`
	Af                string `json:",omitempty"`
	TnBgpCtxAfPolName string `json:",omitempty"`
	TDn               string `json:",omitempty"`
}

func NewBGPAddressFamilyContextPolicyRelationship(fvRsCtxToBgpCtxAfPolRn, parentDn string, fvRsCtxToBgpCtxAfPolAttr BGPAddressFamilyContextPolicyRelationshipAttributes) *BGPAddressFamilyContextPolicyRelationship {
	dn := fmt.Sprintf("%s/%s", parentDn, fvRsCtxToBgpCtxAfPolRn)
	return &BGPAddressFamilyContextPolicyRelationship{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         FvrsctxtobgpctxafpolClassName,
			Rn:                fvRsCtxToBgpCtxAfPolRn,
		},
		BGPAddressFamilyContextPolicyRelationshipAttributes: fvRsCtxToBgpCtxAfPolAttr,
	}
}

func (fvRsCtxToBgpCtxAfPol *BGPAddressFamilyContextPolicyRelationship) ToMap() (map[string]string, error) {
	fvRsCtxToBgpCtxAfPolMap, err := fvRsCtxToBgpCtxAfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvRsCtxToBgpCtxAfPolMap, "af", fvRsCtxToBgpCtxAfPol.Af)
	A(fvRsCtxToBgpCtxAfPolMap, "annotation", fvRsCtxToBgpCtxAfPol.Annotation)
	A(fvRsCtxToBgpCtxAfPolMap, "tnBgpCtxAfPolName", fvRsCtxToBgpCtxAfPol.TnBgpCtxAfPolName)
	A(fvRsCtxToBgpCtxAfPolMap, "tDn", fvRsCtxToBgpCtxAfPol.TDn)
	return fvRsCtxToBgpCtxAfPolMap, err
}

func BGPAddressFamilyContextPolicyRelationshipFromContainerList(cont *container.Container, index int) *BGPAddressFamilyContextPolicyRelationship {
	BGPAddressFamilyContextPolicyRelationshipCont := cont.S("imdata").Index(index).S(FvrsctxtobgpctxafpolClassName, "attributes")
	return &BGPAddressFamilyContextPolicyRelationship{
		BaseAttributes{
			DistinguishedName: G(BGPAddressFamilyContextPolicyRelationshipCont, "dn"),
			Status:            G(BGPAddressFamilyContextPolicyRelationshipCont, "status"),
			ClassName:         FvrsctxtobgpctxafpolClassName,
			Rn:                G(BGPAddressFamilyContextPolicyRelationshipCont, "rn"),
		},
		BGPAddressFamilyContextPolicyRelationshipAttributes{
			Af:                G(BGPAddressFamilyContextPolicyRelationshipCont, "af"),
			Annotation:        G(BGPAddressFamilyContextPolicyRelationshipCont, "annotation"),
			TnBgpCtxAfPolName: G(BGPAddressFamilyContextPolicyRelationshipCont, "tnBgpCtxAfPolName"),
			TDn:               G(BGPAddressFamilyContextPolicyRelationshipCont, "tDn"),
		},
	}
}

func BGPAddressFamilyContextPolicyRelationshipFromContainer(cont *container.Container) *BGPAddressFamilyContextPolicyRelationship {
	return BGPAddressFamilyContextPolicyRelationshipFromContainerList(cont, 0)
}

func BGPAddressFamilyContextPolicyRelationshipListFromContainer(cont *container.Container) []*BGPAddressFamilyContextPolicyRelationship {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*BGPAddressFamilyContextPolicyRelationship, length)

	for i := 0; i < length; i++ {
		arr[i] = BGPAddressFamilyContextPolicyRelationshipFromContainerList(cont, i)
	}

	return arr
}

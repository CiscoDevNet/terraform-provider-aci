package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	Dnl3extVirtualLIfPLagPolAtt        = "uni/tn-%s/out-%s/lnodep-%s/lifp-%s/vlifp-[%s]-[%s]/rsdynPathAtt-[%s]/vlifplagpolatt"
	Rnl3extVirtualLIfPLagPolAtt        = "vlifplagpolatt"
	ParentDnl3extVirtualLIfPLagPolAtt  = "uni/tn-%s/out-%s/lnodep-%s/lifp-%s/vlifp-[%s]-[%s]/rsdynPathAtt-[%s]"
	L3extvirtuallifplagpolattClassName = "l3extVirtualLIfPLagPolAtt"
)

type L3extVirtualLIfPLagPolicy struct {
	BaseAttributes
	L3extVirtualLIfPLagPolicyAttributes
}

type L3extVirtualLIfPLagPolicyAttributes struct {
	Annotation string `json:",omitempty"`
}

func NewL3extVirtualLIfPLagPolicy(l3extVirtualLIfPLagPolAttRn, parentDn, description string, l3extVirtualLIfPLagPolAttAttr L3extVirtualLIfPLagPolicyAttributes) *L3extVirtualLIfPLagPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extVirtualLIfPLagPolAttRn)
	return &L3extVirtualLIfPLagPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extvirtuallifplagpolattClassName,
			Rn:                l3extVirtualLIfPLagPolAttRn,
		},
		L3extVirtualLIfPLagPolicyAttributes: l3extVirtualLIfPLagPolAttAttr,
	}
}

func (l3extVirtualLIfPLagPolAtt *L3extVirtualLIfPLagPolicy) ToMap() (map[string]string, error) {
	l3extVirtualLIfPLagPolAttMap, err := l3extVirtualLIfPLagPolAtt.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extVirtualLIfPLagPolAttMap, "annotation", l3extVirtualLIfPLagPolAtt.Annotation)
	return l3extVirtualLIfPLagPolAttMap, err
}

func L3extVirtualLIfPLagPolicyFromContainerList(cont *container.Container, index int) *L3extVirtualLIfPLagPolicy {
	L3extVirtualLIfPLagPolicyCont := cont.S("imdata").Index(index).S(L3extvirtuallifplagpolattClassName, "attributes")
	return &L3extVirtualLIfPLagPolicy{
		BaseAttributes{
			DistinguishedName: G(L3extVirtualLIfPLagPolicyCont, "dn"),
			Description:       G(L3extVirtualLIfPLagPolicyCont, "descr"),
			Status:            G(L3extVirtualLIfPLagPolicyCont, "status"),
			ClassName:         L3extvirtuallifplagpolattClassName,
			Rn:                G(L3extVirtualLIfPLagPolicyCont, "rn"),
		},
		L3extVirtualLIfPLagPolicyAttributes{
			Annotation: G(L3extVirtualLIfPLagPolicyCont, "annotation"),
		},
	}
}

func L3extVirtualLIfPLagPolicyFromContainer(cont *container.Container) *L3extVirtualLIfPLagPolicy {
	return L3extVirtualLIfPLagPolicyFromContainerList(cont, 0)
}

func L3extVirtualLIfPLagPolicyListFromContainer(cont *container.Container) []*L3extVirtualLIfPLagPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*L3extVirtualLIfPLagPolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = L3extVirtualLIfPLagPolicyFromContainerList(cont, i)
	}

	return arr
}

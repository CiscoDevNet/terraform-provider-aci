package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnL3extDefaultRouteLeakP        = "defrtleak"
	DnL3extDefaultRouteLeakP        = "uni/tn-%s/out-%s/defrtleak"
	ParentDnL3extDefaultRouteLeakP  = "uni/tn-%s/out-%s"
	L3extDefaultRouteLeakPClassName = "l3extDefaultRouteLeakP"
)

type DefaultRouteLeakPolicy struct {
	BaseAttributes
	DefaultRouteLeakPolicyAttributes
}

type DefaultRouteLeakPolicyAttributes struct {
	Always     string `json:",omitempty"`
	Annotation string `json:",omitempty"`
	Criteria   string `json:",omitempty"`
	Scope      string `json:",omitempty"`
}

func NewDefaultRouteLeakPolicy(l3extDefaultRouteLeakPRn, parentDn string, l3extDefaultRouteLeakPAttr DefaultRouteLeakPolicyAttributes) *DefaultRouteLeakPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extDefaultRouteLeakPRn)
	return &DefaultRouteLeakPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         L3extDefaultRouteLeakPClassName,
			Rn:                l3extDefaultRouteLeakPRn,
		},
		DefaultRouteLeakPolicyAttributes: l3extDefaultRouteLeakPAttr,
	}
}

func (l3extDefaultRouteLeakP *DefaultRouteLeakPolicy) ToMap() (map[string]string, error) {
	l3extDefaultRouteLeakPMap, err := l3extDefaultRouteLeakP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extDefaultRouteLeakPMap, "always", l3extDefaultRouteLeakP.Always)
	A(l3extDefaultRouteLeakPMap, "annotation", l3extDefaultRouteLeakP.Annotation)
	A(l3extDefaultRouteLeakPMap, "criteria", l3extDefaultRouteLeakP.Criteria)
	A(l3extDefaultRouteLeakPMap, "scope", l3extDefaultRouteLeakP.Scope)
	return l3extDefaultRouteLeakPMap, err
}

func DefaultRouteLeakPolicyFromContainerList(cont *container.Container, index int) *DefaultRouteLeakPolicy {
	DefaultRouteLeakPolicyCont := cont.S("imdata").Index(index).S(L3extDefaultRouteLeakPClassName, "attributes")
	return &DefaultRouteLeakPolicy{
		BaseAttributes{
			DistinguishedName: G(DefaultRouteLeakPolicyCont, "dn"),
			Status:            G(DefaultRouteLeakPolicyCont, "status"),
			ClassName:         L3extDefaultRouteLeakPClassName,
			Rn:                G(DefaultRouteLeakPolicyCont, "rn"),
		},
		DefaultRouteLeakPolicyAttributes{
			Always:     G(DefaultRouteLeakPolicyCont, "always"),
			Annotation: G(DefaultRouteLeakPolicyCont, "annotation"),
			Criteria:   G(DefaultRouteLeakPolicyCont, "criteria"),
			Scope:      G(DefaultRouteLeakPolicyCont, "scope"),
		},
	}
}

func DefaultRouteLeakPolicyFromContainer(cont *container.Container) *DefaultRouteLeakPolicy {
	return DefaultRouteLeakPolicyFromContainerList(cont, 0)
}

func DefaultRouteLeakPolicyListFromContainer(cont *container.Container) []*DefaultRouteLeakPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*DefaultRouteLeakPolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = DefaultRouteLeakPolicyFromContainerList(cont, i)
	}

	return arr
}

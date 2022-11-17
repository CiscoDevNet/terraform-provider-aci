package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const L3extroutetagpolClassName = "l3extRouteTagPol"

type L3outRouteTagPolicy struct {
	BaseAttributes
	L3outRouteTagPolicyAttributes
}

type L3outRouteTagPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Tag string `json:",omitempty"`
}

func NewL3outRouteTagPolicy(l3extRouteTagPolRn, parentDn, description string, l3extRouteTagPolattr L3outRouteTagPolicyAttributes) *L3outRouteTagPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extRouteTagPolRn)
	return &L3outRouteTagPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extroutetagpolClassName,
			Rn:                l3extRouteTagPolRn,
		},

		L3outRouteTagPolicyAttributes: l3extRouteTagPolattr,
	}
}

func (l3extRouteTagPol *L3outRouteTagPolicy) ToMap() (map[string]string, error) {
	l3extRouteTagPolMap, err := l3extRouteTagPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extRouteTagPolMap, "name", l3extRouteTagPol.Name)

	A(l3extRouteTagPolMap, "annotation", l3extRouteTagPol.Annotation)

	A(l3extRouteTagPolMap, "nameAlias", l3extRouteTagPol.NameAlias)

	A(l3extRouteTagPolMap, "tag", l3extRouteTagPol.Tag)

	return l3extRouteTagPolMap, err
}

func L3outRouteTagPolicyFromContainerList(cont *container.Container, index int) *L3outRouteTagPolicy {

	L3outRouteTagPolicyCont := cont.S("imdata").Index(index).S(L3extroutetagpolClassName, "attributes")
	return &L3outRouteTagPolicy{
		BaseAttributes{
			DistinguishedName: G(L3outRouteTagPolicyCont, "dn"),
			Description:       G(L3outRouteTagPolicyCont, "descr"),
			Status:            G(L3outRouteTagPolicyCont, "status"),
			ClassName:         L3extroutetagpolClassName,
			Rn:                G(L3outRouteTagPolicyCont, "rn"),
		},

		L3outRouteTagPolicyAttributes{

			Name: G(L3outRouteTagPolicyCont, "name"),

			Annotation: G(L3outRouteTagPolicyCont, "annotation"),

			NameAlias: G(L3outRouteTagPolicyCont, "nameAlias"),

			Tag: G(L3outRouteTagPolicyCont, "tag"),
		},
	}
}

func L3outRouteTagPolicyFromContainer(cont *container.Container) *L3outRouteTagPolicy {

	return L3outRouteTagPolicyFromContainerList(cont, 0)
}

func L3outRouteTagPolicyListFromContainer(cont *container.Container) []*L3outRouteTagPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3outRouteTagPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = L3outRouteTagPolicyFromContainerList(cont, i)
	}

	return arr
}

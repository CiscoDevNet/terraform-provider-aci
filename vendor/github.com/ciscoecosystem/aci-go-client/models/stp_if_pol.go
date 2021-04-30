package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnstpIfPol        = "uni/infra/ifPol-%s"
	RnstpIfPol        = "ifPol-%s"
	ParentDnstpIfPol  = "uni/infra"
	StpifpolClassName = "stpIfPol"
)

type SpanningTreeInterfacePolicy struct {
	BaseAttributes
	NameAliasAttribute
	SpanningTreeInterfacePolicyAttributes
}

type SpanningTreeInterfacePolicyAttributes struct {
	Annotation string `json:",omitempty"`
	Ctrl       string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewSpanningTreeInterfacePolicy(stpIfPolRn, parentDn, description, nameAlias string, stpIfPolAttr SpanningTreeInterfacePolicyAttributes) *SpanningTreeInterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, stpIfPolRn)
	return &SpanningTreeInterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         StpifpolClassName,
			Rn:                stpIfPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		SpanningTreeInterfacePolicyAttributes: stpIfPolAttr,
	}
}

func (stpIfPol *SpanningTreeInterfacePolicy) ToMap() (map[string]string, error) {
	stpIfPolMap, err := stpIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := stpIfPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(stpIfPolMap, key, value)
	}
	A(stpIfPolMap, "annotation", stpIfPol.Annotation)
	A(stpIfPolMap, "ctrl", stpIfPol.Ctrl)
	A(stpIfPolMap, "name", stpIfPol.Name)
	return stpIfPolMap, err
}

func SpanningTreeInterfacePolicyFromContainerList(cont *container.Container, index int) *SpanningTreeInterfacePolicy {
	SpanningTreeInterfacePolicyCont := cont.S("imdata").Index(index).S(StpifpolClassName, "attributes")
	return &SpanningTreeInterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(SpanningTreeInterfacePolicyCont, "dn"),
			Description:       G(SpanningTreeInterfacePolicyCont, "descr"),
			Status:            G(SpanningTreeInterfacePolicyCont, "status"),
			ClassName:         StpifpolClassName,
			Rn:                G(SpanningTreeInterfacePolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(SpanningTreeInterfacePolicyCont, "nameAlias"),
		},
		SpanningTreeInterfacePolicyAttributes{
			Annotation: G(SpanningTreeInterfacePolicyCont, "annotation"),
			Ctrl:       G(SpanningTreeInterfacePolicyCont, "ctrl"),
			Name:       G(SpanningTreeInterfacePolicyCont, "name"),
		},
	}
}

func SpanningTreeInterfacePolicyFromContainer(cont *container.Container) *SpanningTreeInterfacePolicy {
	return SpanningTreeInterfacePolicyFromContainerList(cont, 0)
}

func SpanningTreeInterfacePolicyListFromContainer(cont *container.Container) []*SpanningTreeInterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*SpanningTreeInterfacePolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = SpanningTreeInterfacePolicyFromContainerList(cont, i)
	}
	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3ifpolClassName = "l3IfPol"

type L3InterfacePolicy struct {
	BaseAttributes
	L3InterfacePolicyAttributes
}

type L3InterfacePolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	BfdIsis string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewL3InterfacePolicy(l3IfPolRn, parentDn, description string, l3IfPolattr L3InterfacePolicyAttributes) *L3InterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, l3IfPolRn)
	return &L3InterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3ifpolClassName,
			Rn:                l3IfPolRn,
		},

		L3InterfacePolicyAttributes: l3IfPolattr,
	}
}

func (l3IfPol *L3InterfacePolicy) ToMap() (map[string]string, error) {
	l3IfPolMap, err := l3IfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3IfPolMap, "name", l3IfPol.Name)

	A(l3IfPolMap, "annotation", l3IfPol.Annotation)

	A(l3IfPolMap, "bfdIsis", l3IfPol.BfdIsis)

	A(l3IfPolMap, "nameAlias", l3IfPol.NameAlias)

	return l3IfPolMap, err
}

func L3InterfacePolicyFromContainerList(cont *container.Container, index int) *L3InterfacePolicy {

	L3InterfacePolicyCont := cont.S("imdata").Index(index).S(L3ifpolClassName, "attributes")
	return &L3InterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(L3InterfacePolicyCont, "dn"),
			Description:       G(L3InterfacePolicyCont, "descr"),
			Status:            G(L3InterfacePolicyCont, "status"),
			ClassName:         L3ifpolClassName,
			Rn:                G(L3InterfacePolicyCont, "rn"),
		},

		L3InterfacePolicyAttributes{

			Name: G(L3InterfacePolicyCont, "name"),

			Annotation: G(L3InterfacePolicyCont, "annotation"),

			BfdIsis: G(L3InterfacePolicyCont, "bfdIsis"),

			NameAlias: G(L3InterfacePolicyCont, "nameAlias"),
		},
	}
}

func L3InterfacePolicyFromContainer(cont *container.Container) *L3InterfacePolicy {

	return L3InterfacePolicyFromContainerList(cont, 0)
}

func L3InterfacePolicyListFromContainer(cont *container.Container) []*L3InterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3InterfacePolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = L3InterfacePolicyFromContainerList(cont, i)
	}

	return arr
}

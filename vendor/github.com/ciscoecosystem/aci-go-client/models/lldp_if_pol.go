package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const LldpifpolClassName = "lldpIfPol"

type LLDPInterfacePolicy struct {
	BaseAttributes
	LLDPInterfacePolicyAttributes
}

type LLDPInterfacePolicyAttributes struct {
	Name string `json:",omitempty"`

	AdminRxSt string `json:",omitempty"`

	AdminTxSt string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewLLDPInterfacePolicy(lldpIfPolRn, parentDn, description string, lldpIfPolattr LLDPInterfacePolicyAttributes) *LLDPInterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, lldpIfPolRn)
	return &LLDPInterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         LldpifpolClassName,
			Rn:                lldpIfPolRn,
		},

		LLDPInterfacePolicyAttributes: lldpIfPolattr,
	}
}

func (lldpIfPol *LLDPInterfacePolicy) ToMap() (map[string]string, error) {
	lldpIfPolMap, err := lldpIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(lldpIfPolMap, "name", lldpIfPol.Name)

	A(lldpIfPolMap, "adminRxSt", lldpIfPol.AdminRxSt)

	A(lldpIfPolMap, "adminTxSt", lldpIfPol.AdminTxSt)

	A(lldpIfPolMap, "annotation", lldpIfPol.Annotation)

	A(lldpIfPolMap, "nameAlias", lldpIfPol.NameAlias)

	return lldpIfPolMap, err
}

func LLDPInterfacePolicyFromContainerList(cont *container.Container, index int) *LLDPInterfacePolicy {

	LLDPInterfacePolicyCont := cont.S("imdata").Index(index).S(LldpifpolClassName, "attributes")
	return &LLDPInterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(LLDPInterfacePolicyCont, "dn"),
			Description:       G(LLDPInterfacePolicyCont, "descr"),
			Status:            G(LLDPInterfacePolicyCont, "status"),
			ClassName:         LldpifpolClassName,
			Rn:                G(LLDPInterfacePolicyCont, "rn"),
		},

		LLDPInterfacePolicyAttributes{

			Name: G(LLDPInterfacePolicyCont, "name"),

			AdminRxSt: G(LLDPInterfacePolicyCont, "adminRxSt"),

			AdminTxSt: G(LLDPInterfacePolicyCont, "adminTxSt"),

			Annotation: G(LLDPInterfacePolicyCont, "annotation"),

			NameAlias: G(LLDPInterfacePolicyCont, "nameAlias"),
		},
	}
}

func LLDPInterfacePolicyFromContainer(cont *container.Container) *LLDPInterfacePolicy {

	return LLDPInterfacePolicyFromContainerList(cont, 0)
}

func LLDPInterfacePolicyListFromContainer(cont *container.Container) []*LLDPInterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LLDPInterfacePolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = LLDPInterfacePolicyFromContainerList(cont, i)
	}

	return arr
}

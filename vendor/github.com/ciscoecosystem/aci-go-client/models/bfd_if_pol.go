package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BfdifpolClassName = "bfdIfPol"

type BFDInterfacePolicy struct {
	BaseAttributes
	BFDInterfacePolicyAttributes
}

type BFDInterfacePolicyAttributes struct {
	Name string `json:",omitempty"`

	AdminSt string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Ctrl string `json:",omitempty"`

	DetectMult string `json:",omitempty"`

	EchoAdminSt string `json:",omitempty"`

	EchoRxIntvl string `json:",omitempty"`

	MinRxIntvl string `json:",omitempty"`

	MinTxIntvl string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewBFDInterfacePolicy(bfdIfPolRn, parentDn, description string, bfdIfPolattr BFDInterfacePolicyAttributes) *BFDInterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, bfdIfPolRn)
	return &BFDInterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BfdifpolClassName,
			Rn:                bfdIfPolRn,
		},

		BFDInterfacePolicyAttributes: bfdIfPolattr,
	}
}

func (bfdIfPol *BFDInterfacePolicy) ToMap() (map[string]string, error) {
	bfdIfPolMap, err := bfdIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bfdIfPolMap, "name", bfdIfPol.Name)

	A(bfdIfPolMap, "adminSt", bfdIfPol.AdminSt)

	A(bfdIfPolMap, "annotation", bfdIfPol.Annotation)

	A(bfdIfPolMap, "ctrl", bfdIfPol.Ctrl)

	A(bfdIfPolMap, "detectMult", bfdIfPol.DetectMult)

	A(bfdIfPolMap, "echoAdminSt", bfdIfPol.EchoAdminSt)

	A(bfdIfPolMap, "echoRxIntvl", bfdIfPol.EchoRxIntvl)

	A(bfdIfPolMap, "minRxIntvl", bfdIfPol.MinRxIntvl)

	A(bfdIfPolMap, "minTxIntvl", bfdIfPol.MinTxIntvl)

	A(bfdIfPolMap, "nameAlias", bfdIfPol.NameAlias)

	return bfdIfPolMap, err
}

func BFDInterfacePolicyFromContainerList(cont *container.Container, index int) *BFDInterfacePolicy {

	BFDInterfacePolicyCont := cont.S("imdata").Index(index).S(BfdifpolClassName, "attributes")
	return &BFDInterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(BFDInterfacePolicyCont, "dn"),
			Description:       G(BFDInterfacePolicyCont, "descr"),
			Status:            G(BFDInterfacePolicyCont, "status"),
			ClassName:         BfdifpolClassName,
			Rn:                G(BFDInterfacePolicyCont, "rn"),
		},

		BFDInterfacePolicyAttributes{

			Name: G(BFDInterfacePolicyCont, "name"),

			AdminSt: G(BFDInterfacePolicyCont, "adminSt"),

			Annotation: G(BFDInterfacePolicyCont, "annotation"),

			Ctrl: G(BFDInterfacePolicyCont, "ctrl"),

			DetectMult: G(BFDInterfacePolicyCont, "detectMult"),

			EchoAdminSt: G(BFDInterfacePolicyCont, "echoAdminSt"),

			EchoRxIntvl: G(BFDInterfacePolicyCont, "echoRxIntvl"),

			MinRxIntvl: G(BFDInterfacePolicyCont, "minRxIntvl"),

			MinTxIntvl: G(BFDInterfacePolicyCont, "minTxIntvl"),

			NameAlias: G(BFDInterfacePolicyCont, "nameAlias"),
		},
	}
}

func BFDInterfacePolicyFromContainer(cont *container.Container) *BFDInterfacePolicy {

	return BFDInterfacePolicyFromContainerList(cont, 0)
}

func BFDInterfacePolicyListFromContainer(cont *container.Container) []*BFDInterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*BFDInterfacePolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = BFDInterfacePolicyFromContainerList(cont, i)
	}

	return arr
}

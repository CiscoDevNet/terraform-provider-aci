package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnbfdMhIfPol        = "uni/tn-%s/bfdMhIfPol-%s"
	RnbfdMhIfPol        = "bfdMhIfPol-%s"
	ParentDnbfdMhIfPol  = "uni/tn-%s"
	BfdmhifpolClassName = "bfdMhIfPol"
)

type AciBfdMultihopInterfacePolicy struct {
	BaseAttributes
	NameAliasAttribute
	AciBfdMultihopInterfacePolicyAttributes
}

type AciBfdMultihopInterfacePolicyAttributes struct {
	Annotation string `json:",omitempty"`
	AdminSt    string `json:",omitempty"`
	DetectMult string `json:",omitempty"`
	MinRxIntvl string `json:",omitempty"`
	MinTxIntvl string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewAciBfdMultihopInterfacePolicy(bfdMhIfPolRn, parentDn, description, nameAlias string, bfdMhIfPolAttr AciBfdMultihopInterfacePolicyAttributes) *AciBfdMultihopInterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, bfdMhIfPolRn)
	return &AciBfdMultihopInterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BfdmhifpolClassName,
			Rn:                bfdMhIfPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		AciBfdMultihopInterfacePolicyAttributes: bfdMhIfPolAttr,
	}
}

func (bfdMhIfPol *AciBfdMultihopInterfacePolicy) ToMap() (map[string]string, error) {
	bfdMhIfPolMap, err := bfdMhIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := bfdMhIfPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(bfdMhIfPolMap, key, value)
	}

	A(bfdMhIfPolMap, "adminSt", bfdMhIfPol.AdminSt)
	A(bfdMhIfPolMap, "annotation", bfdMhIfPol.Annotation)
	A(bfdMhIfPolMap, "detectMult", bfdMhIfPol.DetectMult)
	A(bfdMhIfPolMap, "minRxIntvl", bfdMhIfPol.MinRxIntvl)
	A(bfdMhIfPolMap, "minTxIntvl", bfdMhIfPol.MinTxIntvl)
	A(bfdMhIfPolMap, "name", bfdMhIfPol.Name)
	return bfdMhIfPolMap, err
}

func AciBfdMultihopInterfacePolicyFromContainerList(cont *container.Container, index int) *AciBfdMultihopInterfacePolicy {
	AciBfdMultihopInterfacePolicyCont := cont.S("imdata").Index(index).S(BfdmhifpolClassName, "attributes")
	return &AciBfdMultihopInterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(AciBfdMultihopInterfacePolicyCont, "dn"),
			Description:       G(AciBfdMultihopInterfacePolicyCont, "descr"),
			Status:            G(AciBfdMultihopInterfacePolicyCont, "status"),
			ClassName:         BfdmhifpolClassName,
			Rn:                G(AciBfdMultihopInterfacePolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(AciBfdMultihopInterfacePolicyCont, "nameAlias"),
		},
		AciBfdMultihopInterfacePolicyAttributes{
			AdminSt:    G(AciBfdMultihopInterfacePolicyCont, "adminSt"),
			Annotation: G(AciBfdMultihopInterfacePolicyCont, "annotation"),
			DetectMult: G(AciBfdMultihopInterfacePolicyCont, "detectMult"),
			MinRxIntvl: G(AciBfdMultihopInterfacePolicyCont, "minRxIntvl"),
			MinTxIntvl: G(AciBfdMultihopInterfacePolicyCont, "minTxIntvl"),
			Name:       G(AciBfdMultihopInterfacePolicyCont, "name"),
		},
	}
}

func AciBfdMultihopInterfacePolicyFromContainer(cont *container.Container) *AciBfdMultihopInterfacePolicy {
	return AciBfdMultihopInterfacePolicyFromContainerList(cont, 0)
}

func AciBfdMultihopInterfacePolicyListFromContainer(cont *container.Container) []*AciBfdMultihopInterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*AciBfdMultihopInterfacePolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = AciBfdMultihopInterfacePolicyFromContainerList(cont, i)
	}

	return arr
}

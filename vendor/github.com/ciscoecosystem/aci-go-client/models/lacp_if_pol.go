package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnlacpIfPol        = "uni/infra/lacpifp-%s"
	RnlacpIfPol        = "lacpifp-%s"
	ParentDnlacpIfPol  = "uni/infra"
	LacpifpolClassName = "lacpIfPol"
)

type LACPMemberPolicy struct {
	BaseAttributes
	LACPMemberPolicyAttributes
}

type LACPMemberPolicyAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
	Prio       string `json:",omitempty"`
	TxRate     string `json:",omitempty"`
}

func NewLACPMemberPolicy(lacpIfPolRn, parentDn, description string, lacpIfPolAttr LACPMemberPolicyAttributes) *LACPMemberPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, lacpIfPolRn)
	return &LACPMemberPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         LacpifpolClassName,
			Rn:                lacpIfPolRn,
		},
		LACPMemberPolicyAttributes: lacpIfPolAttr,
	}
}

func (lacpIfPol *LACPMemberPolicy) ToMap() (map[string]string, error) {
	lacpIfPolMap, err := lacpIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(lacpIfPolMap, "annotation", lacpIfPol.Annotation)
	A(lacpIfPolMap, "name", lacpIfPol.Name)
	A(lacpIfPolMap, "nameAlias", lacpIfPol.NameAlias)
	A(lacpIfPolMap, "prio", lacpIfPol.Prio)
	A(lacpIfPolMap, "txRate", lacpIfPol.TxRate)
	return lacpIfPolMap, err
}

func LACPMemberPolicyFromContainerList(cont *container.Container, index int) *LACPMemberPolicy {
	LACPMemberPolicyCont := cont.S("imdata").Index(index).S(LacpifpolClassName, "attributes")
	return &LACPMemberPolicy{
		BaseAttributes{
			DistinguishedName: G(LACPMemberPolicyCont, "dn"),
			Description:       G(LACPMemberPolicyCont, "descr"),
			Status:            G(LACPMemberPolicyCont, "status"),
			ClassName:         LacpifpolClassName,
			Rn:                G(LACPMemberPolicyCont, "rn"),
		},
		LACPMemberPolicyAttributes{
			Annotation: G(LACPMemberPolicyCont, "annotation"),
			Name:       G(LACPMemberPolicyCont, "name"),
			NameAlias:  G(LACPMemberPolicyCont, "nameAlias"),
			Prio:       G(LACPMemberPolicyCont, "prio"),
			TxRate:     G(LACPMemberPolicyCont, "txRate"),
		},
	}
}

func LACPMemberPolicyFromContainer(cont *container.Container) *LACPMemberPolicy {
	return LACPMemberPolicyFromContainerList(cont, 0)
}

func LACPMemberPolicyListFromContainer(cont *container.Container) []*LACPMemberPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*LACPMemberPolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = LACPMemberPolicyFromContainerList(cont, i)
	}

	return arr
}

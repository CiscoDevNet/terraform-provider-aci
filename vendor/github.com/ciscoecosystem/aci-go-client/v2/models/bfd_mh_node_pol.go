package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnBfdMhNodePol        = "bfdMhNodePol-%s"
	DnBfdMhNodePol        = "uni/tn-%s/bfdMhNodePol-%s"
	ParentDnBfdMhNodePol  = "uni/tn-%s"
	BfdMhNodePolClassName = "bfdMhNodePol"
)

type BFDMultihopNodePolicy struct {
	BaseAttributes
	BFDMultihopNodePolicyAttributes
}

type BFDMultihopNodePolicyAttributes struct {
	AdminSt    string `json:",omitempty"`
	Annotation string `json:",omitempty"`
	DetectMult string `json:",omitempty"`
	MinRxIntvl string `json:",omitempty"`
	MinTxIntvl string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewBFDMultihopNodePolicy(bfdMhNodePolRn, parentDn, description string, bfdMhNodePolAttr BFDMultihopNodePolicyAttributes) *BFDMultihopNodePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, bfdMhNodePolRn)
	return &BFDMultihopNodePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BfdMhNodePolClassName,
			Rn:                bfdMhNodePolRn,
		},
		BFDMultihopNodePolicyAttributes: bfdMhNodePolAttr,
	}
}

func (bfdMhNodePol *BFDMultihopNodePolicy) ToMap() (map[string]string, error) {
	bfdMhNodePolMap, err := bfdMhNodePol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bfdMhNodePolMap, "adminSt", bfdMhNodePol.AdminSt)
	A(bfdMhNodePolMap, "annotation", bfdMhNodePol.Annotation)
	A(bfdMhNodePolMap, "detectMult", bfdMhNodePol.DetectMult)
	A(bfdMhNodePolMap, "minRxIntvl", bfdMhNodePol.MinRxIntvl)
	A(bfdMhNodePolMap, "minTxIntvl", bfdMhNodePol.MinTxIntvl)
	A(bfdMhNodePolMap, "name", bfdMhNodePol.Name)
	A(bfdMhNodePolMap, "nameAlias", bfdMhNodePol.NameAlias)
	return bfdMhNodePolMap, err
}

func BFDMultihopNodePolicyFromContainerList(cont *container.Container, index int) *BFDMultihopNodePolicy {
	BFDMultihopNodePolicyCont := cont.S("imdata").Index(index).S(BfdMhNodePolClassName, "attributes")
	return &BFDMultihopNodePolicy{
		BaseAttributes{
			DistinguishedName: G(BFDMultihopNodePolicyCont, "dn"),
			Description:       G(BFDMultihopNodePolicyCont, "descr"),
			Status:            G(BFDMultihopNodePolicyCont, "status"),
			ClassName:         BfdMhNodePolClassName,
			Rn:                G(BFDMultihopNodePolicyCont, "rn"),
		},
		BFDMultihopNodePolicyAttributes{
			AdminSt:    G(BFDMultihopNodePolicyCont, "adminSt"),
			Annotation: G(BFDMultihopNodePolicyCont, "annotation"),
			DetectMult: G(BFDMultihopNodePolicyCont, "detectMult"),
			MinRxIntvl: G(BFDMultihopNodePolicyCont, "minRxIntvl"),
			MinTxIntvl: G(BFDMultihopNodePolicyCont, "minTxIntvl"),
			Name:       G(BFDMultihopNodePolicyCont, "name"),
			NameAlias:  G(BFDMultihopNodePolicyCont, "nameAlias"),
		},
	}
}

func BFDMultihopNodePolicyFromContainer(cont *container.Container) *BFDMultihopNodePolicy {
	return BFDMultihopNodePolicyFromContainerList(cont, 0)
}

func BFDMultihopNodePolicyListFromContainer(cont *container.Container) []*BFDMultihopNodePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*BFDMultihopNodePolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = BFDMultihopNodePolicyFromContainerList(cont, i)
	}

	return arr
}

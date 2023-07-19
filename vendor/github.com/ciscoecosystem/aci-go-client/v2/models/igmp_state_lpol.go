package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnIgmpStateLPol        = "igmpstateLPol"
	DnIgmpStateLPol        = "uni/tn-%s/igmpIfPol-%s/igmpstateLPol"
	ParentDnIgmpStateLPol  = "uni/tn-%s/igmpIfPol-%s"
	IgmpStateLPolClassName = "igmpStateLPol"
)

type IGMPStateLimitPolicy struct {
	BaseAttributes
	IGMPStateLimitPolicyAttributes
}

type IGMPStateLimitPolicyAttributes struct {
	Annotation string `json:",omitempty"`
	Max        string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
	Rsvd       string `json:",omitempty"`
}

func NewIGMPStateLimitPolicy(parentDn, description string, igmpStateLPolAttr IGMPStateLimitPolicyAttributes) *IGMPStateLimitPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, RnIgmpStateLPol)
	return &IGMPStateLimitPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         IgmpStateLPolClassName,
			Rn:                RnIgmpStateLPol,
		},
		IGMPStateLimitPolicyAttributes: igmpStateLPolAttr,
	}
}

func (igmpStateLPol *IGMPStateLimitPolicy) ToMap() (map[string]string, error) {
	igmpStateLPolMap, err := igmpStateLPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(igmpStateLPolMap, "annotation", igmpStateLPol.Annotation)
	A(igmpStateLPolMap, "max", igmpStateLPol.Max)
	A(igmpStateLPolMap, "name", igmpStateLPol.Name)
	A(igmpStateLPolMap, "nameAlias", igmpStateLPol.NameAlias)
	A(igmpStateLPolMap, "rsvd", igmpStateLPol.Rsvd)
	return igmpStateLPolMap, err
}

func IGMPStateLimitPolicyFromContainerList(cont *container.Container, index int) *IGMPStateLimitPolicy {
	IGMPStateLimitPolicyCont := cont.S("imdata").Index(index).S(IgmpStateLPolClassName, "attributes")
	return &IGMPStateLimitPolicy{
		BaseAttributes{
			DistinguishedName: G(IGMPStateLimitPolicyCont, "dn"),
			Description:       G(IGMPStateLimitPolicyCont, "descr"),
			Status:            G(IGMPStateLimitPolicyCont, "status"),
			ClassName:         IgmpStateLPolClassName,
			Rn:                G(IGMPStateLimitPolicyCont, "rn"),
		},
		IGMPStateLimitPolicyAttributes{
			Annotation: G(IGMPStateLimitPolicyCont, "annotation"),
			Max:        G(IGMPStateLimitPolicyCont, "max"),
			Name:       G(IGMPStateLimitPolicyCont, "name"),
			NameAlias:  G(IGMPStateLimitPolicyCont, "nameAlias"),
			Rsvd:       G(IGMPStateLimitPolicyCont, "rsvd"),
		},
	}
}

func IGMPStateLimitPolicyFromContainer(cont *container.Container) *IGMPStateLimitPolicy {
	return IGMPStateLimitPolicyFromContainerList(cont, 0)
}

func IGMPStateLimitPolicyListFromContainer(cont *container.Container) []*IGMPStateLimitPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*IGMPStateLimitPolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = IGMPStateLimitPolicyFromContainerList(cont, i)
	}

	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnIgmpRepPol        = "igmprepPol"
	DnIgmpRepPol        = "uni/tn-%s/igmpIfPol-%s/igmprepPol"
	ParentDnIgmpRepPol  = "uni/tn-%s/igmpIfPol-%s"
	IgmpRepPolClassName = "igmpRepPol"
)

type IGMPReportPolicy struct {
	BaseAttributes
	IGMPReportPolicyAttributes
}

type IGMPReportPolicyAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewIGMPReportPolicy(parentDn, description string, igmpRepPolAttr IGMPReportPolicyAttributes) *IGMPReportPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, RnIgmpRepPol)
	return &IGMPReportPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         IgmpRepPolClassName,
			Rn:                RnIgmpRepPol,
		},
		IGMPReportPolicyAttributes: igmpRepPolAttr,
	}
}

func (igmpRepPol *IGMPReportPolicy) ToMap() (map[string]string, error) {
	igmpRepPolMap, err := igmpRepPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(igmpRepPolMap, "annotation", igmpRepPol.Annotation)
	A(igmpRepPolMap, "name", igmpRepPol.Name)
	A(igmpRepPolMap, "nameAlias", igmpRepPol.NameAlias)
	return igmpRepPolMap, err
}

func IGMPReportPolicyFromContainerList(cont *container.Container, index int) *IGMPReportPolicy {
	IGMPReportPolicyCont := cont.S("imdata").Index(index).S(IgmpRepPolClassName, "attributes")
	return &IGMPReportPolicy{
		BaseAttributes{
			DistinguishedName: G(IGMPReportPolicyCont, "dn"),
			Description:       G(IGMPReportPolicyCont, "descr"),
			Status:            G(IGMPReportPolicyCont, "status"),
			ClassName:         IgmpRepPolClassName,
			Rn:                G(IGMPReportPolicyCont, "rn"),
		},
		IGMPReportPolicyAttributes{
			Annotation: G(IGMPReportPolicyCont, "annotation"),
			Name:       G(IGMPReportPolicyCont, "name"),
			NameAlias:  G(IGMPReportPolicyCont, "nameAlias"),
		},
	}
}

func IGMPReportPolicyFromContainer(cont *container.Container) *IGMPReportPolicy {
	return IGMPReportPolicyFromContainerList(cont, 0)
}

func IGMPReportPolicyListFromContainer(cont *container.Container) []*IGMPReportPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*IGMPReportPolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = IGMPReportPolicyFromContainerList(cont, i)
	}

	return arr
}

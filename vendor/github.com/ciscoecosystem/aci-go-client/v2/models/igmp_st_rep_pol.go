package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnIgmpStRepPol        = "igmpstrepPol-static-group"
	DnIgmpStRepPol        = "uni/tn-%s/igmpIfPol-%s/igmpstrepPol-static-group"
	ParentDnIgmpStRepPol  = "uni/tn-%s/igmpIfPol-%s"
	IgmpStRepPolClassName = "igmpStRepPol"
)

type IGMPStaticReportPolicy struct {
	BaseAttributes
	IGMPStaticReportPolicyAttributes
}

type IGMPStaticReportPolicyAttributes struct {
	Annotation string `json:",omitempty"`
	JoinType   string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewIGMPStaticReportPolicy(parentDn, description string, igmpStRepPolAttr IGMPStaticReportPolicyAttributes) *IGMPStaticReportPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, RnIgmpStRepPol)
	return &IGMPStaticReportPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         IgmpStRepPolClassName,
			Rn:                RnIgmpStRepPol,
		},
		IGMPStaticReportPolicyAttributes: igmpStRepPolAttr,
	}
}

func (igmpStRepPol *IGMPStaticReportPolicy) ToMap() (map[string]string, error) {
	igmpStRepPolMap, err := igmpStRepPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(igmpStRepPolMap, "annotation", igmpStRepPol.Annotation)
	A(igmpStRepPolMap, "name", igmpStRepPol.Name)
	A(igmpStRepPolMap, "nameAlias", igmpStRepPol.NameAlias)
	return igmpStRepPolMap, err
}

func IGMPStaticReportPolicyFromContainerList(cont *container.Container, index int) *IGMPStaticReportPolicy {
	IGMPStaticReportPolicyCont := cont.S("imdata").Index(index).S(IgmpStRepPolClassName, "attributes")
	return &IGMPStaticReportPolicy{
		BaseAttributes{
			DistinguishedName: G(IGMPStaticReportPolicyCont, "dn"),
			Description:       G(IGMPStaticReportPolicyCont, "descr"),
			Status:            G(IGMPStaticReportPolicyCont, "status"),
			ClassName:         IgmpStRepPolClassName,
			Rn:                G(IGMPStaticReportPolicyCont, "rn"),
		},
		IGMPStaticReportPolicyAttributes{
			Annotation: G(IGMPStaticReportPolicyCont, "annotation"),
			Name:       G(IGMPStaticReportPolicyCont, "name"),
			NameAlias:  G(IGMPStaticReportPolicyCont, "nameAlias"),
		},
	}
}

func IGMPStaticReportPolicyFromContainer(cont *container.Container) *IGMPStaticReportPolicy {
	return IGMPStaticReportPolicyFromContainerList(cont, 0)
}

func IGMPStaticReportPolicyListFromContainer(cont *container.Container) []*IGMPStaticReportPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*IGMPStaticReportPolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = IGMPStaticReportPolicyFromContainerList(cont, i)
	}

	return arr
}

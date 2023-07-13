package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnPimJPInbFilterPol        = "jpinbfilter"
	DnPimJPInbFilterPol        = "uni/tn-%s/pimifpol-%s/jpinbfilter"
	ParentDnPimJPInbFilterPol  = "uni/tn-%s/pimifpol-%s"
	PimJPInbFilterPolClassName = "pimJPInbFilterPol"
)

type PIMJPInboundFilterPolicy struct {
	BaseAttributes
	PIMJPInboundFilterPolicyAttributes
}

type PIMJPInboundFilterPolicyAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewPIMJPInboundFilterPolicy(parentDn, description string, pimJPInbFilterPolAttr PIMJPInboundFilterPolicyAttributes) *PIMJPInboundFilterPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, RnPimJPInbFilterPol)
	return &PIMJPInboundFilterPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PimJPInbFilterPolClassName,
			Rn:                RnPimJPInbFilterPol,
		},
		PIMJPInboundFilterPolicyAttributes: pimJPInbFilterPolAttr,
	}
}

func (pimJPInbFilterPol *PIMJPInboundFilterPolicy) ToMap() (map[string]string, error) {
	pimJPInbFilterPolMap, err := pimJPInbFilterPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(pimJPInbFilterPolMap, "annotation", pimJPInbFilterPol.Annotation)
	A(pimJPInbFilterPolMap, "name", pimJPInbFilterPol.Name)
	A(pimJPInbFilterPolMap, "nameAlias", pimJPInbFilterPol.NameAlias)
	return pimJPInbFilterPolMap, err
}

func PIMJPInboundFilterPolicyFromContainerList(cont *container.Container, index int) *PIMJPInboundFilterPolicy {
	PIMJPInboundFilterPolicyCont := cont.S("imdata").Index(index).S(PimJPInbFilterPolClassName, "attributes")
	return &PIMJPInboundFilterPolicy{
		BaseAttributes{
			DistinguishedName: G(PIMJPInboundFilterPolicyCont, "dn"),
			Description:       G(PIMJPInboundFilterPolicyCont, "descr"),
			Status:            G(PIMJPInboundFilterPolicyCont, "status"),
			ClassName:         PimJPInbFilterPolClassName,
			Rn:                G(PIMJPInboundFilterPolicyCont, "rn"),
		},
		PIMJPInboundFilterPolicyAttributes{
			Annotation: G(PIMJPInboundFilterPolicyCont, "annotation"),
			Name:       G(PIMJPInboundFilterPolicyCont, "name"),
			NameAlias:  G(PIMJPInboundFilterPolicyCont, "nameAlias"),
		},
	}
}

func PIMJPInboundFilterPolicyFromContainer(cont *container.Container) *PIMJPInboundFilterPolicy {
	return PIMJPInboundFilterPolicyFromContainerList(cont, 0)
}

func PIMJPInboundFilterPolicyListFromContainer(cont *container.Container) []*PIMJPInboundFilterPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PIMJPInboundFilterPolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = PIMJPInboundFilterPolicyFromContainerList(cont, i)
	}

	return arr
}

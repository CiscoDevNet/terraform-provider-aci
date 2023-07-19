package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnPimJPOutbFilterPol        = "jpoutbfilter"
	DnPimJPOutbFilterPol        = "uni/tn-%s/pimifpol-%s/jpoutbfilter"
	ParentDnPimJPOutbFilterPol  = "uni/tn-%s/pimifpol-%s"
	PimJPOutbFilterPolClassName = "pimJPOutbFilterPol"
)

type PIMJPOutboundFilterPolicy struct {
	BaseAttributes
	PIMJPOutboundFilterPolicyAttributes
}

type PIMJPOutboundFilterPolicyAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewPIMJPOutboundFilterPolicy(parentDn, description string, pimJPOutbFilterPolAttr PIMJPOutboundFilterPolicyAttributes) *PIMJPOutboundFilterPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, RnPimJPOutbFilterPol)
	return &PIMJPOutboundFilterPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PimJPOutbFilterPolClassName,
			Rn:                RnPimJPOutbFilterPol,
		},
		PIMJPOutboundFilterPolicyAttributes: pimJPOutbFilterPolAttr,
	}
}

func (pimJPOutbFilterPol *PIMJPOutboundFilterPolicy) ToMap() (map[string]string, error) {
	pimJPOutbFilterPolMap, err := pimJPOutbFilterPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(pimJPOutbFilterPolMap, "annotation", pimJPOutbFilterPol.Annotation)
	A(pimJPOutbFilterPolMap, "name", pimJPOutbFilterPol.Name)
	A(pimJPOutbFilterPolMap, "nameAlias", pimJPOutbFilterPol.NameAlias)
	return pimJPOutbFilterPolMap, err
}

func PIMJPOutboundFilterPolicyFromContainerList(cont *container.Container, index int) *PIMJPOutboundFilterPolicy {
	PIMJPOutboundFilterPolicyCont := cont.S("imdata").Index(index).S(PimJPOutbFilterPolClassName, "attributes")
	return &PIMJPOutboundFilterPolicy{
		BaseAttributes{
			DistinguishedName: G(PIMJPOutboundFilterPolicyCont, "dn"),
			Description:       G(PIMJPOutboundFilterPolicyCont, "descr"),
			Status:            G(PIMJPOutboundFilterPolicyCont, "status"),
			ClassName:         PimJPOutbFilterPolClassName,
			Rn:                G(PIMJPOutboundFilterPolicyCont, "rn"),
		},
		PIMJPOutboundFilterPolicyAttributes{
			Annotation: G(PIMJPOutboundFilterPolicyCont, "annotation"),
			Name:       G(PIMJPOutboundFilterPolicyCont, "name"),
			NameAlias:  G(PIMJPOutboundFilterPolicyCont, "nameAlias"),
		},
	}
}

func PIMJPOutboundFilterPolicyFromContainer(cont *container.Container) *PIMJPOutboundFilterPolicy {
	return PIMJPOutboundFilterPolicyFromContainerList(cont, 0)
}

func PIMJPOutboundFilterPolicyListFromContainer(cont *container.Container) []*PIMJPOutboundFilterPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PIMJPOutboundFilterPolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = PIMJPOutboundFilterPolicyFromContainerList(cont, i)
	}

	return arr
}

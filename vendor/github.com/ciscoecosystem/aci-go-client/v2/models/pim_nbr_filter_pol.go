package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnPimNbrFilterPol        = "nbrfilter"
	DnPimNbrFilterPol        = "uni/tn-%s/pimifpol-%s/nbrfilter"
	ParentDnPimNbrFilterPol  = "uni/tn-%s/pimifpol-%s"
	PimNbrFilterPolClassName = "pimNbrFilterPol"
)

type PIMNeighborFiterPolicy struct {
	BaseAttributes
	PIMNeighborFiterPolicyAttributes
}

type PIMNeighborFiterPolicyAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewPIMNeighborFiterPolicy(parentDn, description string, pimNbrFilterPolAttr PIMNeighborFiterPolicyAttributes) *PIMNeighborFiterPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, RnPimNbrFilterPol)
	return &PIMNeighborFiterPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PimNbrFilterPolClassName,
			Rn:                RnPimNbrFilterPol,
		},
		PIMNeighborFiterPolicyAttributes: pimNbrFilterPolAttr,
	}
}

func (pimNbrFilterPol *PIMNeighborFiterPolicy) ToMap() (map[string]string, error) {
	pimNbrFilterPolMap, err := pimNbrFilterPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(pimNbrFilterPolMap, "annotation", pimNbrFilterPol.Annotation)
	A(pimNbrFilterPolMap, "name", pimNbrFilterPol.Name)
	A(pimNbrFilterPolMap, "nameAlias", pimNbrFilterPol.NameAlias)
	return pimNbrFilterPolMap, err
}

func PIMNeighborFiterPolicyFromContainerList(cont *container.Container, index int) *PIMNeighborFiterPolicy {
	PIMNeighborFiterPolicyCont := cont.S("imdata").Index(index).S(PimNbrFilterPolClassName, "attributes")
	return &PIMNeighborFiterPolicy{
		BaseAttributes{
			DistinguishedName: G(PIMNeighborFiterPolicyCont, "dn"),
			Description:       G(PIMNeighborFiterPolicyCont, "descr"),
			Status:            G(PIMNeighborFiterPolicyCont, "status"),
			ClassName:         PimNbrFilterPolClassName,
			Rn:                G(PIMNeighborFiterPolicyCont, "rn"),
		},
		PIMNeighborFiterPolicyAttributes{
			Annotation: G(PIMNeighborFiterPolicyCont, "annotation"),
			Name:       G(PIMNeighborFiterPolicyCont, "name"),
			NameAlias:  G(PIMNeighborFiterPolicyCont, "nameAlias"),
		},
	}
}

func PIMNeighborFiterPolicyFromContainer(cont *container.Container) *PIMNeighborFiterPolicy {
	return PIMNeighborFiterPolicyFromContainerList(cont, 0)
}

func PIMNeighborFiterPolicyListFromContainer(cont *container.Container) []*PIMNeighborFiterPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PIMNeighborFiterPolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = PIMNeighborFiterPolicyFromContainerList(cont, i)
	}

	return arr
}

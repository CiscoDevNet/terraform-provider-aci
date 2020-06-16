package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VmmsecpClassName = "vmmSecP"

type VMMSecurityPolicy struct {
	BaseAttributes
	VMMSecurityPolicyAttributes
}

type VMMSecurityPolicyAttributes struct {
	AllowPromiscuous string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ForgedTransmits string `json:",omitempty"`

	MacChanges string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewVMMSecurityPolicy(vmmSecPRn, parentDn, description string, vmmSecPattr VMMSecurityPolicyAttributes) *VMMSecurityPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, vmmSecPRn)
	return &VMMSecurityPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VmmsecpClassName,
			Rn:                vmmSecPRn,
		},

		VMMSecurityPolicyAttributes: vmmSecPattr,
	}
}

func (vmmSecP *VMMSecurityPolicy) ToMap() (map[string]string, error) {
	vmmSecPMap, err := vmmSecP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vmmSecPMap, "allowPromiscuous", vmmSecP.AllowPromiscuous)

	A(vmmSecPMap, "annotation", vmmSecP.Annotation)

	A(vmmSecPMap, "forgedTransmits", vmmSecP.ForgedTransmits)

	A(vmmSecPMap, "macChanges", vmmSecP.MacChanges)

	A(vmmSecPMap, "nameAlias", vmmSecP.NameAlias)

	return vmmSecPMap, err
}

func VMMSecurityPolicyFromContainerList(cont *container.Container, index int) *VMMSecurityPolicy {

	VMMSecurityPolicyCont := cont.S("imdata").Index(index).S(VmmsecpClassName, "attributes")
	return &VMMSecurityPolicy{
		BaseAttributes{
			DistinguishedName: G(VMMSecurityPolicyCont, "dn"),
			Description:       G(VMMSecurityPolicyCont, "descr"),
			Status:            G(VMMSecurityPolicyCont, "status"),
			ClassName:         VmmsecpClassName,
			Rn:                G(VMMSecurityPolicyCont, "rn"),
		},

		VMMSecurityPolicyAttributes{

			AllowPromiscuous: G(VMMSecurityPolicyCont, "allowPromiscuous"),

			Annotation: G(VMMSecurityPolicyCont, "annotation"),

			ForgedTransmits: G(VMMSecurityPolicyCont, "forgedTransmits"),

			MacChanges: G(VMMSecurityPolicyCont, "macChanges"),

			NameAlias: G(VMMSecurityPolicyCont, "nameAlias"),
		},
	}
}

func VMMSecurityPolicyFromContainer(cont *container.Container) *VMMSecurityPolicy {

	return VMMSecurityPolicyFromContainerList(cont, 0)
}

func VMMSecurityPolicyListFromContainer(cont *container.Container) []*VMMSecurityPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*VMMSecurityPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = VMMSecurityPolicyFromContainerList(cont, i)
	}

	return arr
}

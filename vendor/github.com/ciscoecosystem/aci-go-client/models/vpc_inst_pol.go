package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnvpcInstPol        = "uni/fabric/vpcInst-%s"
	RnvpcInstPol        = "vpcInst-%s"
	ParentDnvpcInstPol  = "uni/fabric"
	VpcinstpolClassName = "vpcInstPol"
)

type VPCDomainPolicy struct {
	BaseAttributes
	NameAliasAttribute
	VPCDomainPolicyAttributes
}

type VPCDomainPolicyAttributes struct {
	Annotation string `json:",omitempty"`
	DeadIntvl  string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewVPCDomainPolicy(vpcInstPolRn, parentDn, description, nameAlias string, vpcInstPolAttr VPCDomainPolicyAttributes) *VPCDomainPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, vpcInstPolRn)
	return &VPCDomainPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VpcinstpolClassName,
			Rn:                vpcInstPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		VPCDomainPolicyAttributes: vpcInstPolAttr,
	}
}

func (vpcInstPol *VPCDomainPolicy) ToMap() (map[string]string, error) {
	vpcInstPolMap, err := vpcInstPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := vpcInstPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(vpcInstPolMap, key, value)
	}
	A(vpcInstPolMap, "annotation", vpcInstPol.Annotation)
	A(vpcInstPolMap, "deadIntvl", vpcInstPol.DeadIntvl)
	A(vpcInstPolMap, "name", vpcInstPol.Name)
	return vpcInstPolMap, err
}

func VPCDomainPolicyFromContainerList(cont *container.Container, index int) *VPCDomainPolicy {
	VPCDomainPolicyCont := cont.S("imdata").Index(index).S(VpcinstpolClassName, "attributes")
	return &VPCDomainPolicy{
		BaseAttributes{
			DistinguishedName: G(VPCDomainPolicyCont, "dn"),
			Description:       G(VPCDomainPolicyCont, "descr"),
			Status:            G(VPCDomainPolicyCont, "status"),
			ClassName:         VpcinstpolClassName,
			Rn:                G(VPCDomainPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(VPCDomainPolicyCont, "nameAlias"),
		},
		VPCDomainPolicyAttributes{
			Annotation: G(VPCDomainPolicyCont, "annotation"),
			DeadIntvl:  G(VPCDomainPolicyCont, "deadIntvl"),
			Name:       G(VPCDomainPolicyCont, "name"),
		},
	}
}

func VPCDomainPolicyFromContainer(cont *container.Container) *VPCDomainPolicy {
	return VPCDomainPolicyFromContainerList(cont, 0)
}

func VPCDomainPolicyListFromContainer(cont *container.Container) []*VPCDomainPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*VPCDomainPolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = VPCDomainPolicyFromContainerList(cont, i)
	}
	return arr
}

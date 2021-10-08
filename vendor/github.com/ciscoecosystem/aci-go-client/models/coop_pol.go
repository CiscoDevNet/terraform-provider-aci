package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DncoopPol        = "uni/fabric/pol-%s"
	RncoopPol        = "pol-%s"
	ParentDncoopPol  = "uni/fabric"
	CooppolClassName = "coopPol"
)

type COOPGroupPolicy struct {
	BaseAttributes
	NameAliasAttribute
	COOPGroupPolicyAttributes
}

type COOPGroupPolicyAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Type       string `json:",omitempty"`
}

func NewCOOPGroupPolicy(coopPolRn, parentDn, description, nameAlias string, coopPolAttr COOPGroupPolicyAttributes) *COOPGroupPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, coopPolRn)
	return &COOPGroupPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CooppolClassName,
			Rn:                coopPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		COOPGroupPolicyAttributes: coopPolAttr,
	}
}

func (coopPol *COOPGroupPolicy) ToMap() (map[string]string, error) {
	coopPolMap, err := coopPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := coopPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(coopPolMap, key, value)
	}
	A(coopPolMap, "annotation", coopPol.Annotation)
	A(coopPolMap, "name", coopPol.Name)
	A(coopPolMap, "type", coopPol.Type)
	return coopPolMap, err
}

func COOPGroupPolicyFromContainerList(cont *container.Container, index int) *COOPGroupPolicy {
	COOPGroupPolicyCont := cont.S("imdata").Index(index).S(CooppolClassName, "attributes")
	return &COOPGroupPolicy{
		BaseAttributes{
			DistinguishedName: G(COOPGroupPolicyCont, "dn"),
			Description:       G(COOPGroupPolicyCont, "descr"),
			Status:            G(COOPGroupPolicyCont, "status"),
			ClassName:         CooppolClassName,
			Rn:                G(COOPGroupPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(COOPGroupPolicyCont, "nameAlias"),
		},
		COOPGroupPolicyAttributes{
			Annotation: G(COOPGroupPolicyCont, "annotation"),
			Name:       G(COOPGroupPolicyCont, "name"),
			Type:       G(COOPGroupPolicyCont, "type"),
		},
	}
}

func COOPGroupPolicyFromContainer(cont *container.Container) *COOPGroupPolicy {
	return COOPGroupPolicyFromContainerList(cont, 0)
}

func COOPGroupPolicyListFromContainer(cont *container.Container) []*COOPGroupPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*COOPGroupPolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = COOPGroupPolicyFromContainerList(cont, i)
	}
	return arr
}

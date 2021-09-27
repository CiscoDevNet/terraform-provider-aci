package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnepLoopProtectP        = "uni/infra/epLoopProtectP-%s"
	RnepLoopProtectP        = "epLoopProtectP-%s"
	ParentDnepLoopProtectP  = "uni/infra"
	EploopprotectpClassName = "epLoopProtectP"
)

type EPLoopProtectionPolicy struct {
	BaseAttributes
	NameAliasAttribute
	EPLoopProtectionPolicyAttributes
}

type EPLoopProtectionPolicyAttributes struct {
	Action          string `json:",omitempty"`
	AdminSt         string `json:",omitempty"`
	Annotation      string `json:",omitempty"`
	LoopDetectIntvl string `json:",omitempty"`
	LoopDetectMult  string `json:",omitempty"`
	Name            string `json:",omitempty"`
}

func NewEPLoopProtectionPolicy(epLoopProtectPRn, parentDn, description, nameAlias string, epLoopProtectPAttr EPLoopProtectionPolicyAttributes) *EPLoopProtectionPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, epLoopProtectPRn)
	return &EPLoopProtectionPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         EploopprotectpClassName,
			Rn:                epLoopProtectPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		EPLoopProtectionPolicyAttributes: epLoopProtectPAttr,
	}
}

func (epLoopProtectP *EPLoopProtectionPolicy) ToMap() (map[string]string, error) {
	epLoopProtectPMap, err := epLoopProtectP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := epLoopProtectP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(epLoopProtectPMap, key, value)
	}
	A(epLoopProtectPMap, "action", epLoopProtectP.Action)
	A(epLoopProtectPMap, "adminSt", epLoopProtectP.AdminSt)
	A(epLoopProtectPMap, "annotation", epLoopProtectP.Annotation)
	A(epLoopProtectPMap, "loopDetectIntvl", epLoopProtectP.LoopDetectIntvl)
	A(epLoopProtectPMap, "loopDetectMult", epLoopProtectP.LoopDetectMult)
	A(epLoopProtectPMap, "name", epLoopProtectP.Name)
	return epLoopProtectPMap, err
}

func EPLoopProtectionPolicyFromContainerList(cont *container.Container, index int) *EPLoopProtectionPolicy {
	EPLoopProtectionPolicyCont := cont.S("imdata").Index(index).S(EploopprotectpClassName, "attributes")
	return &EPLoopProtectionPolicy{
		BaseAttributes{
			DistinguishedName: G(EPLoopProtectionPolicyCont, "dn"),
			Description:       G(EPLoopProtectionPolicyCont, "descr"),
			Status:            G(EPLoopProtectionPolicyCont, "status"),
			ClassName:         EploopprotectpClassName,
			Rn:                G(EPLoopProtectionPolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(EPLoopProtectionPolicyCont, "nameAlias"),
		},
		EPLoopProtectionPolicyAttributes{
			Action:          G(EPLoopProtectionPolicyCont, "action"),
			AdminSt:         G(EPLoopProtectionPolicyCont, "adminSt"),
			Annotation:      G(EPLoopProtectionPolicyCont, "annotation"),
			LoopDetectIntvl: G(EPLoopProtectionPolicyCont, "loopDetectIntvl"),
			LoopDetectMult:  G(EPLoopProtectionPolicyCont, "loopDetectMult"),
			Name:            G(EPLoopProtectionPolicyCont, "name"),
		},
	}
}

func EPLoopProtectionPolicyFromContainer(cont *container.Container) *EPLoopProtectionPolicy {
	return EPLoopProtectionPolicyFromContainerList(cont, 0)
}

func EPLoopProtectionPolicyListFromContainer(cont *container.Container) []*EPLoopProtectionPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*EPLoopProtectionPolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = EPLoopProtectionPolicyFromContainerList(cont, i)
	}
	return arr
}

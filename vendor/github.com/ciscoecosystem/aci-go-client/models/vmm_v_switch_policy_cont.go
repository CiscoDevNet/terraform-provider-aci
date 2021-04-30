package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnvmmVSwitchPolicyCont        = "uni/vmmp-%s/dom-%s/vswitchpolcont"
	RnvmmVSwitchPolicyCont        = "vswitchpolcont"
	ParentDnvmmVSwitchPolicyCont  = "uni/vmmp-%s/dom-%s"
	VmmvswitchpolicycontClassName = "vmmVSwitchPolicyCont"
)

type VSwitchPolicyGroup struct {
	BaseAttributes
	NameAliasAttribute
	VSwitchPolicyGroupAttributes
}

type VSwitchPolicyGroupAttributes struct {
	Annotation string `json:",omitempty"`
}

func NewVSwitchPolicyGroup(vmmVSwitchPolicyContRn, parentDn, description, nameAlias string, vmmVSwitchPolicyContAttr VSwitchPolicyGroupAttributes) *VSwitchPolicyGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, vmmVSwitchPolicyContRn)
	return &VSwitchPolicyGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VmmvswitchpolicycontClassName,
			Rn:                vmmVSwitchPolicyContRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		VSwitchPolicyGroupAttributes: vmmVSwitchPolicyContAttr,
	}
}

func (vmmVSwitchPolicyCont *VSwitchPolicyGroup) ToMap() (map[string]string, error) {
	vmmVSwitchPolicyContMap, err := vmmVSwitchPolicyCont.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := vmmVSwitchPolicyCont.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(vmmVSwitchPolicyContMap, key, value)
	}
	A(vmmVSwitchPolicyContMap, "annotation", vmmVSwitchPolicyCont.Annotation)
	return vmmVSwitchPolicyContMap, err
}

func VSwitchPolicyGroupFromContainerList(cont *container.Container, index int) *VSwitchPolicyGroup {
	VSwitchPolicyGroupCont := cont.S("imdata").Index(index).S(VmmvswitchpolicycontClassName, "attributes")
	return &VSwitchPolicyGroup{
		BaseAttributes{
			DistinguishedName: G(VSwitchPolicyGroupCont, "dn"),
			Description:       G(VSwitchPolicyGroupCont, "descr"),
			Status:            G(VSwitchPolicyGroupCont, "status"),
			ClassName:         VmmvswitchpolicycontClassName,
			Rn:                G(VSwitchPolicyGroupCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(VSwitchPolicyGroupCont, "nameAlias"),
		},
		VSwitchPolicyGroupAttributes{
			Annotation: G(VSwitchPolicyGroupCont, "annotation"),
		},
	}
}

func VSwitchPolicyGroupFromContainer(cont *container.Container) *VSwitchPolicyGroup {
	return VSwitchPolicyGroupFromContainerList(cont, 0)
}

func VSwitchPolicyGroupListFromContainer(cont *container.Container) []*VSwitchPolicyGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*VSwitchPolicyGroup, length)
	for i := 0; i < length; i++ {
		arr[i] = VSwitchPolicyGroupFromContainerList(cont, i)
	}
	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DninfraSpineAccNodePGrp        = "uni/infra/funcprof/spaccnodepgrp-%s"
	RninfraSpineAccNodePGrp        = "spaccnodepgrp-%s"
	ParentDninfraSpineAccNodePGrp  = "uni/infra/funcprof"
	InfraspineaccnodepgrpClassName = "infraSpineAccNodePGrp"
)

type SpineSwitchPolicyGroup struct {
	BaseAttributes
	NameAliasAttribute
	SpineSwitchPolicyGroupAttributes
}

type SpineSwitchPolicyGroupAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewSpineSwitchPolicyGroup(infraSpineAccNodePGrpRn, parentDn, description, nameAlias string, infraSpineAccNodePGrpAttr SpineSwitchPolicyGroupAttributes) *SpineSwitchPolicyGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, infraSpineAccNodePGrpRn)
	return &SpineSwitchPolicyGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraspineaccnodepgrpClassName,
			Rn:                infraSpineAccNodePGrpRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		SpineSwitchPolicyGroupAttributes: infraSpineAccNodePGrpAttr,
	}
}

func (infraSpineAccNodePGrp *SpineSwitchPolicyGroup) ToMap() (map[string]string, error) {
	infraSpineAccNodePGrpMap, err := infraSpineAccNodePGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := infraSpineAccNodePGrp.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(infraSpineAccNodePGrpMap, key, value)
	}
	A(infraSpineAccNodePGrpMap, "annotation", infraSpineAccNodePGrp.Annotation)
	A(infraSpineAccNodePGrpMap, "name", infraSpineAccNodePGrp.Name)
	return infraSpineAccNodePGrpMap, err
}

func SpineSwitchPolicyGroupFromContainerList(cont *container.Container, index int) *SpineSwitchPolicyGroup {
	SpineSwitchPolicyGroupCont := cont.S("imdata").Index(index).S(InfraspineaccnodepgrpClassName, "attributes")
	return &SpineSwitchPolicyGroup{
		BaseAttributes{
			DistinguishedName: G(SpineSwitchPolicyGroupCont, "dn"),
			Description:       G(SpineSwitchPolicyGroupCont, "descr"),
			Status:            G(SpineSwitchPolicyGroupCont, "status"),
			ClassName:         InfraspineaccnodepgrpClassName,
			Rn:                G(SpineSwitchPolicyGroupCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(SpineSwitchPolicyGroupCont, "nameAlias"),
		},
		SpineSwitchPolicyGroupAttributes{
			Annotation: G(SpineSwitchPolicyGroupCont, "annotation"),
			Name:       G(SpineSwitchPolicyGroupCont, "name"),
		},
	}
}

func SpineSwitchPolicyGroupFromContainer(cont *container.Container) *SpineSwitchPolicyGroup {
	return SpineSwitchPolicyGroupFromContainerList(cont, 0)
}

func SpineSwitchPolicyGroupListFromContainer(cont *container.Container) []*SpineSwitchPolicyGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*SpineSwitchPolicyGroup, length)
	for i := 0; i < length; i++ {
		arr[i] = SpineSwitchPolicyGroupFromContainerList(cont, i)
	}
	return arr
}

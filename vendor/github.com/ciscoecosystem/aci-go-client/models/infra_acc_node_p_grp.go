package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DninfraAccNodePGrp        = "uni/infra/funcprof/accnodepgrp-%s"
	RninfraAccNodePGrp        = "accnodepgrp-%s"
	ParentDninfraAccNodePGrp  = "uni/infra/funcprof"
	InfraaccnodepgrpClassName = "infraAccNodePGrp"
)

type AccessSwitchPolicyGroup struct {
	BaseAttributes
	NameAliasAttribute
	AccessSwitchPolicyGroupAttributes
}

type AccessSwitchPolicyGroupAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewAccessSwitchPolicyGroup(infraAccNodePGrpRn, parentDn, description, nameAlias string, infraAccNodePGrpAttr AccessSwitchPolicyGroupAttributes) *AccessSwitchPolicyGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, infraAccNodePGrpRn)
	return &AccessSwitchPolicyGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraaccnodepgrpClassName,
			Rn:                infraAccNodePGrpRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		AccessSwitchPolicyGroupAttributes: infraAccNodePGrpAttr,
	}
}

func (infraAccNodePGrp *AccessSwitchPolicyGroup) ToMap() (map[string]string, error) {
	infraAccNodePGrpMap, err := infraAccNodePGrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := infraAccNodePGrp.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(infraAccNodePGrpMap, key, value)
	}
	A(infraAccNodePGrpMap, "annotation", infraAccNodePGrp.Annotation)
	A(infraAccNodePGrpMap, "name", infraAccNodePGrp.Name)
	return infraAccNodePGrpMap, err
}

func AccessSwitchPolicyGroupFromContainerList(cont *container.Container, index int) *AccessSwitchPolicyGroup {
	AccessSwitchPolicyGroupCont := cont.S("imdata").Index(index).S(InfraaccnodepgrpClassName, "attributes")
	return &AccessSwitchPolicyGroup{
		BaseAttributes{
			DistinguishedName: G(AccessSwitchPolicyGroupCont, "dn"),
			Description:       G(AccessSwitchPolicyGroupCont, "descr"),
			Status:            G(AccessSwitchPolicyGroupCont, "status"),
			ClassName:         InfraaccnodepgrpClassName,
			Rn:                G(AccessSwitchPolicyGroupCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(AccessSwitchPolicyGroupCont, "nameAlias"),
		},
		AccessSwitchPolicyGroupAttributes{
			Annotation: G(AccessSwitchPolicyGroupCont, "annotation"),
			Name:       G(AccessSwitchPolicyGroupCont, "name"),
		},
	}
}

func AccessSwitchPolicyGroupFromContainer(cont *container.Container) *AccessSwitchPolicyGroup {
	return AccessSwitchPolicyGroupFromContainerList(cont, 0)
}

func AccessSwitchPolicyGroupListFromContainer(cont *container.Container) []*AccessSwitchPolicyGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*AccessSwitchPolicyGroup, length)
	for i := 0; i < length; i++ {
		arr[i] = AccessSwitchPolicyGroupFromContainerList(cont, i)
	}
	return arr
}

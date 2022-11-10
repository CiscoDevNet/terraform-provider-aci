package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DninfraAccBndlSubgrp        = "uni/infra/funcprof/accbundle-%s/accsubbndl-%s"
	RninfraAccBndlSubgrp        = "accsubbndl-%s"
	ParentDninfraAccBndlSubgrp  = "uni/infra/funcprof/accbundle-%s"
	InfraaccbndlsubgrpClassName = "infraAccBndlSubgrp"
)

type OverridePCVPCPolicyGroup struct {
	BaseAttributes
	OverridePCVPCPolicyGroupAttributes
}

type OverridePCVPCPolicyGroupAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewOverridePCVPCPolicyGroup(infraAccBndlSubgrpRn, parentDn, description string, infraAccBndlSubgrpAttr OverridePCVPCPolicyGroupAttributes) *OverridePCVPCPolicyGroup {
	dn := fmt.Sprintf("%s/%s", parentDn, infraAccBndlSubgrpRn)
	return &OverridePCVPCPolicyGroup{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfraaccbndlsubgrpClassName,
			Rn:                infraAccBndlSubgrpRn,
		},
		OverridePCVPCPolicyGroupAttributes: infraAccBndlSubgrpAttr,
	}
}

func (infraAccBndlSubgrp *OverridePCVPCPolicyGroup) ToMap() (map[string]string, error) {
	infraAccBndlSubgrpMap, err := infraAccBndlSubgrp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraAccBndlSubgrpMap, "annotation", infraAccBndlSubgrp.Annotation)
	A(infraAccBndlSubgrpMap, "name", infraAccBndlSubgrp.Name)
	A(infraAccBndlSubgrpMap, "nameAlias", infraAccBndlSubgrp.NameAlias)
	return infraAccBndlSubgrpMap, err
}

func OverridePCVPCPolicyGroupFromContainerList(cont *container.Container, index int) *OverridePCVPCPolicyGroup {
	OverridePCVPCPolicyGroupCont := cont.S("imdata").Index(index).S(InfraaccbndlsubgrpClassName, "attributes")
	return &OverridePCVPCPolicyGroup{
		BaseAttributes{
			DistinguishedName: G(OverridePCVPCPolicyGroupCont, "dn"),
			Description:       G(OverridePCVPCPolicyGroupCont, "descr"),
			Status:            G(OverridePCVPCPolicyGroupCont, "status"),
			ClassName:         InfraaccbndlsubgrpClassName,
			Rn:                G(OverridePCVPCPolicyGroupCont, "rn"),
		},
		OverridePCVPCPolicyGroupAttributes{
			Annotation: G(OverridePCVPCPolicyGroupCont, "annotation"),
			Name:       G(OverridePCVPCPolicyGroupCont, "name"),
			NameAlias:  G(OverridePCVPCPolicyGroupCont, "nameAlias"),
		},
	}
}

func OverridePCVPCPolicyGroupFromContainer(cont *container.Container) *OverridePCVPCPolicyGroup {
	return OverridePCVPCPolicyGroupFromContainerList(cont, 0)
}

func OverridePCVPCPolicyGroupListFromContainer(cont *container.Container) []*OverridePCVPCPolicyGroup {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*OverridePCVPCPolicyGroup, length)

	for i := 0; i < length; i++ {
		arr[i] = OverridePCVPCPolicyGroupFromContainerList(cont, i)
	}

	return arr
}

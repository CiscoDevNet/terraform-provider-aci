package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnfabricNodeControl        = "uni/fabric/nodecontrol-%s"
	RnfabricNodeControl        = "nodecontrol-%s"
	ParentDnfabricNodeControl  = "uni/fabric"
	FabricnodecontrolClassName = "fabricNodeControl"
)

type FabricNodeControl struct {
	BaseAttributes
	NameAliasAttribute
	FabricNodeControlAttributes
}

type FabricNodeControlAttributes struct {
	Annotation string `json:",omitempty"`
	Control    string `json:",omitempty"`
	FeatureSel string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewFabricNodeControl(fabricNodeControlRn, parentDn, description, nameAlias string, fabricNodeControlAttr FabricNodeControlAttributes) *FabricNodeControl {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricNodeControlRn)
	return &FabricNodeControl{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FabricnodecontrolClassName,
			Rn:                fabricNodeControlRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		FabricNodeControlAttributes: fabricNodeControlAttr,
	}
}

func (fabricNodeControl *FabricNodeControl) ToMap() (map[string]string, error) {
	fabricNodeControlMap, err := fabricNodeControl.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := fabricNodeControl.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(fabricNodeControlMap, key, value)
	}
	A(fabricNodeControlMap, "annotation", fabricNodeControl.Annotation)
	A(fabricNodeControlMap, "control", fabricNodeControl.Control)
	A(fabricNodeControlMap, "featureSel", fabricNodeControl.FeatureSel)
	A(fabricNodeControlMap, "name", fabricNodeControl.Name)
	return fabricNodeControlMap, err
}

func FabricNodeControlFromContainerList(cont *container.Container, index int) *FabricNodeControl {
	FabricNodeControlCont := cont.S("imdata").Index(index).S(FabricnodecontrolClassName, "attributes")
	return &FabricNodeControl{
		BaseAttributes{
			DistinguishedName: G(FabricNodeControlCont, "dn"),
			Description:       G(FabricNodeControlCont, "descr"),
			Status:            G(FabricNodeControlCont, "status"),
			ClassName:         FabricnodecontrolClassName,
			Rn:                G(FabricNodeControlCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(FabricNodeControlCont, "nameAlias"),
		},
		FabricNodeControlAttributes{
			Annotation: G(FabricNodeControlCont, "annotation"),
			Control:    G(FabricNodeControlCont, "control"),
			FeatureSel: G(FabricNodeControlCont, "featureSel"),
			Name:       G(FabricNodeControlCont, "name"),
		},
	}
}

func FabricNodeControlFromContainer(cont *container.Container) *FabricNodeControl {
	return FabricNodeControlFromContainerList(cont, 0)
}

func FabricNodeControlListFromContainer(cont *container.Container) []*FabricNodeControl {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*FabricNodeControl, length)
	for i := 0; i < length; i++ {
		arr[i] = FabricNodeControlFromContainerList(cont, i)
	}
	return arr
}

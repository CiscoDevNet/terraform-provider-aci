package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FabricnodeidentpClassName = "fabricNodeIdentP"

type FabricNodeMember struct {
	BaseAttributes
	FabricNodeMemberAttributes
}

type FabricNodeMemberAttributes struct {
	Serial string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ExtPoolId string `json:",omitempty"`

	FabricId string `json:",omitempty"`

	Name string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	NodeId string `json:",omitempty"`

	NodeType string `json:",omitempty"`

	PodId string `json:",omitempty"`

	Role string `json:",omitempty"`
}

func NewFabricNodeMember(fabricNodeIdentPRn, parentDn, description string, fabricNodeIdentPattr FabricNodeMemberAttributes) *FabricNodeMember {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricNodeIdentPRn)
	return &FabricNodeMember{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FabricnodeidentpClassName,
			Rn:                fabricNodeIdentPRn,
		},

		FabricNodeMemberAttributes: fabricNodeIdentPattr,
	}
}

func (fabricNodeIdentP *FabricNodeMember) ToMap() (map[string]string, error) {
	fabricNodeIdentPMap, err := fabricNodeIdentP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fabricNodeIdentPMap, "serial", fabricNodeIdentP.Serial)

	A(fabricNodeIdentPMap, "annotation", fabricNodeIdentP.Annotation)

	A(fabricNodeIdentPMap, "extPoolId", fabricNodeIdentP.ExtPoolId)

	A(fabricNodeIdentPMap, "fabricId", fabricNodeIdentP.FabricId)

	A(fabricNodeIdentPMap, "name", fabricNodeIdentP.Name)

	A(fabricNodeIdentPMap, "nameAlias", fabricNodeIdentP.NameAlias)

	A(fabricNodeIdentPMap, "nodeId", fabricNodeIdentP.NodeId)

	A(fabricNodeIdentPMap, "nodeType", fabricNodeIdentP.NodeType)

	A(fabricNodeIdentPMap, "podId", fabricNodeIdentP.PodId)

	A(fabricNodeIdentPMap, "role", fabricNodeIdentP.Role)

	return fabricNodeIdentPMap, err
}

func FabricNodeMemberFromContainerList(cont *container.Container, index int) *FabricNodeMember {

	FabricNodeMemberCont := cont.S("imdata").Index(index).S(FabricnodeidentpClassName, "attributes")
	return &FabricNodeMember{
		BaseAttributes{
			DistinguishedName: G(FabricNodeMemberCont, "dn"),
			Description:       G(FabricNodeMemberCont, "descr"),
			Status:            G(FabricNodeMemberCont, "status"),
			ClassName:         FabricnodeidentpClassName,
			Rn:                G(FabricNodeMemberCont, "rn"),
		},

		FabricNodeMemberAttributes{

			Serial: G(FabricNodeMemberCont, "serial"),

			Annotation: G(FabricNodeMemberCont, "annotation"),

			ExtPoolId: G(FabricNodeMemberCont, "extPoolId"),

			FabricId: G(FabricNodeMemberCont, "fabricId"),

			NameAlias: G(FabricNodeMemberCont, "nameAlias"),

			Name: G(FabricNodeMemberCont, "name"),

			NodeId: G(FabricNodeMemberCont, "nodeId"),

			NodeType: G(FabricNodeMemberCont, "nodeType"),

			PodId: G(FabricNodeMemberCont, "podId"),

			Role: G(FabricNodeMemberCont, "role"),
		},
	}
}

func FabricNodeMemberFromContainer(cont *container.Container) *FabricNodeMember {

	return FabricNodeMemberFromContainerList(cont, 0)
}

func FabricNodeMemberListFromContainer(cont *container.Container) []*FabricNodeMember {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FabricNodeMember, length)

	for i := 0; i < length; i++ {

		arr[i] = FabricNodeMemberFromContainerList(cont, i)
	}

	return arr
}

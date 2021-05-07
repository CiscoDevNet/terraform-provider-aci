package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VnsabsnodeClassName = "vnsAbsNode"

type FunctionNode struct {
	BaseAttributes
	FunctionNodeAttributes
}

type FunctionNodeAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	FuncTemplateType string `json:",omitempty"`

	FuncType string `json:",omitempty"`

	IsCopy string `json:",omitempty"`

	Managed string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	RoutingMode string `json:",omitempty"`

	SequenceNumber string `json:",omitempty"`

	ShareEncap string `json:",omitempty"`
}

func NewFunctionNode(vnsAbsNodeRn, parentDn, description string, vnsAbsNodeattr FunctionNodeAttributes) *FunctionNode {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsAbsNodeRn)
	return &FunctionNode{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsabsnodeClassName,
			Rn:                vnsAbsNodeRn,
		},

		FunctionNodeAttributes: vnsAbsNodeattr,
	}
}

func (vnsAbsNode *FunctionNode) ToMap() (map[string]string, error) {
	vnsAbsNodeMap, err := vnsAbsNode.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vnsAbsNodeMap, "name", vnsAbsNode.Name)

	A(vnsAbsNodeMap, "annotation", vnsAbsNode.Annotation)

	A(vnsAbsNodeMap, "funcTemplateType", vnsAbsNode.FuncTemplateType)

	A(vnsAbsNodeMap, "funcType", vnsAbsNode.FuncType)

	A(vnsAbsNodeMap, "isCopy", vnsAbsNode.IsCopy)

	A(vnsAbsNodeMap, "managed", vnsAbsNode.Managed)

	A(vnsAbsNodeMap, "nameAlias", vnsAbsNode.NameAlias)

	A(vnsAbsNodeMap, "routingMode", vnsAbsNode.RoutingMode)

	A(vnsAbsNodeMap, "sequenceNumber", vnsAbsNode.SequenceNumber)

	A(vnsAbsNodeMap, "shareEncap", vnsAbsNode.ShareEncap)

	return vnsAbsNodeMap, err
}

func FunctionNodeFromContainerList(cont *container.Container, index int) *FunctionNode {

	FunctionNodeCont := cont.S("imdata").Index(index).S(VnsabsnodeClassName, "attributes")
	return &FunctionNode{
		BaseAttributes{
			DistinguishedName: G(FunctionNodeCont, "dn"),
			Description:       G(FunctionNodeCont, "descr"),
			Status:            G(FunctionNodeCont, "status"),
			ClassName:         VnsabsnodeClassName,
			Rn:                G(FunctionNodeCont, "rn"),
		},

		FunctionNodeAttributes{

			Name: G(FunctionNodeCont, "name"),

			Annotation: G(FunctionNodeCont, "annotation"),

			FuncTemplateType: G(FunctionNodeCont, "funcTemplateType"),

			FuncType: G(FunctionNodeCont, "funcType"),

			IsCopy: G(FunctionNodeCont, "isCopy"),

			Managed: G(FunctionNodeCont, "managed"),

			NameAlias: G(FunctionNodeCont, "nameAlias"),

			RoutingMode: G(FunctionNodeCont, "routingMode"),

			SequenceNumber: G(FunctionNodeCont, "sequenceNumber"),

			ShareEncap: G(FunctionNodeCont, "shareEncap"),
		},
	}
}

func FunctionNodeFromContainer(cont *container.Container) *FunctionNode {

	return FunctionNodeFromContainerList(cont, 0)
}

func FunctionNodeListFromContainer(cont *container.Container) []*FunctionNode {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FunctionNode, length)

	for i := 0; i < length; i++ {

		arr[i] = FunctionNodeFromContainerList(cont, i)
	}

	return arr
}

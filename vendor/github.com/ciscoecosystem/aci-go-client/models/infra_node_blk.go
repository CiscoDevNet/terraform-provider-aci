package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfranodeblkClassName = "infraNodeBlk"

type NodeBlock struct {
	BaseAttributes
	NodeBlockAttributes
}

type NodeBlockAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	From_ string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	To_ string `json:",omitempty"`
}

func NewNodeBlock(infraNodeBlkRn, parentDn, description string, infraNodeBlkattr NodeBlockAttributes) *NodeBlock {
	dn := fmt.Sprintf("%s/%s", parentDn, infraNodeBlkRn)
	return &NodeBlock{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfranodeblkClassName,
			Rn:                infraNodeBlkRn,
		},

		NodeBlockAttributes: infraNodeBlkattr,
	}
}

func (infraNodeBlk *NodeBlock) ToMap() (map[string]string, error) {
	infraNodeBlkMap, err := infraNodeBlk.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraNodeBlkMap, "name", infraNodeBlk.Name)

	A(infraNodeBlkMap, "annotation", infraNodeBlk.Annotation)

	A(infraNodeBlkMap, "from_", infraNodeBlk.From_)

	A(infraNodeBlkMap, "nameAlias", infraNodeBlk.NameAlias)

	A(infraNodeBlkMap, "to_", infraNodeBlk.To_)

	return infraNodeBlkMap, err
}

func NodeBlockFromContainerListBLK(cont *container.Container, index int) *NodeBlock {

	NodeBlockCont := cont.S("imdata").Index(index).S(InfranodeblkClassName, "attributes")
	return &NodeBlock{
		BaseAttributes{
			DistinguishedName: G(NodeBlockCont, "dn"),
			Description:       G(NodeBlockCont, "descr"),
			Status:            G(NodeBlockCont, "status"),
			ClassName:         InfranodeblkClassName,
			Rn:                G(NodeBlockCont, "rn"),
		},

		NodeBlockAttributes{

			Name: G(NodeBlockCont, "name"),

			Annotation: G(NodeBlockCont, "annotation"),

			From_: G(NodeBlockCont, "from_"),

			NameAlias: G(NodeBlockCont, "nameAlias"),

			To_: G(NodeBlockCont, "to_"),
		},
	}
}

func NodeBlockFromContainerBLK(cont *container.Container) *NodeBlock {

	return NodeBlockFromContainerListBLK(cont, 0)
}

func NodeBlockListFromContainerBLK(cont *container.Container) []*NodeBlock {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*NodeBlock, length)

	for i := 0; i < length; i++ {

		arr[i] = NodeBlockFromContainerListBLK(cont, i)
	}

	return arr
}

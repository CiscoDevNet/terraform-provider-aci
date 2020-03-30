package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FabricnodeblkClassNameFW = "fabricNodeBlk"

type NodeBlockFW struct {
	BaseAttributes
	NodeBlockAttributesFW
}

type NodeBlockAttributesFW struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	From_ string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	To_ string `json:",omitempty"`
}

func NewNodeBlockFW(fabricNodeBlkRn, parentDn, description string, fabricNodeBlkattr NodeBlockAttributesFW) *NodeBlockFW {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricNodeBlkRn)
	return &NodeBlockFW{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FabricnodeblkClassNameFW,
			Rn:                fabricNodeBlkRn,
		},

		NodeBlockAttributesFW: fabricNodeBlkattr,
	}
}

func (fabricNodeBlk *NodeBlockFW) ToMap() (map[string]string, error) {
	fabricNodeBlkMap, err := fabricNodeBlk.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fabricNodeBlkMap, "name", fabricNodeBlk.Name)

	A(fabricNodeBlkMap, "annotation", fabricNodeBlk.Annotation)

	A(fabricNodeBlkMap, "from_", fabricNodeBlk.From_)

	A(fabricNodeBlkMap, "nameAlias", fabricNodeBlk.NameAlias)

	A(fabricNodeBlkMap, "to_", fabricNodeBlk.To_)

	return fabricNodeBlkMap, err
}

func NodeBlockFromContainerListFW(cont *container.Container, index int) *NodeBlockFW {

	NodeBlockCont := cont.S("imdata").Index(index).S(FabricnodeblkClassNameFW, "attributes")
	return &NodeBlockFW{
		BaseAttributes{
			DistinguishedName: G(NodeBlockCont, "dn"),
			Description:       G(NodeBlockCont, "descr"),
			Status:            G(NodeBlockCont, "status"),
			ClassName:         FabricnodeblkClassNameFW,
			Rn:                G(NodeBlockCont, "rn"),
		},

		NodeBlockAttributesFW{

			Name: G(NodeBlockCont, "name"),

			Annotation: G(NodeBlockCont, "annotation"),

			From_: G(NodeBlockCont, "from_"),

			NameAlias: G(NodeBlockCont, "nameAlias"),

			To_: G(NodeBlockCont, "to_"),
		},
	}
}

func NodeBlockFromContainer(cont *container.Container) *NodeBlockFW {

	return NodeBlockFromContainerListFW(cont, 0)
}

func NodeBlockListFromContainer(cont *container.Container) []*NodeBlockFW {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*NodeBlockFW, length)

	for i := 0; i < length; i++ {

		arr[i] = NodeBlockFromContainerListFW(cont, i)
	}

	return arr
}

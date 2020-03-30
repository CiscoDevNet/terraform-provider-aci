package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FabricnodeblkClassNameMG = "fabricNodeBlk"

type NodeBlockMG struct {
	BaseAttributes
	NodeBlockAttributesMG
}

type NodeBlockAttributesMG struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	From_ string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	To_ string `json:",omitempty"`
}

func NewNodeBlockMG(fabricNodeBlkRn, parentDn, description string, fabricNodeBlkattr NodeBlockAttributesMG) *NodeBlockMG {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricNodeBlkRn)
	return &NodeBlockMG{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FabricnodeblkClassNameMG,
			Rn:                fabricNodeBlkRn,
		},

		NodeBlockAttributesMG: fabricNodeBlkattr,
	}
}

func (fabricNodeBlk *NodeBlockMG) ToMap() (map[string]string, error) {
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

func NodeBlockFromContainerListMG(cont *container.Container, index int) *NodeBlockMG {

	NodeBlockCont := cont.S("imdata").Index(index).S(FabricnodeblkClassNameMG, "attributes")
	return &NodeBlockMG{
		BaseAttributes{
			DistinguishedName: G(NodeBlockCont, "dn"),
			Description:       G(NodeBlockCont, "descr"),
			Status:            G(NodeBlockCont, "status"),
			ClassName:         FabricnodeblkClassNameMG,
			Rn:                G(NodeBlockCont, "rn"),
		},

		NodeBlockAttributesMG{

			Name: G(NodeBlockCont, "name"),

			Annotation: G(NodeBlockCont, "annotation"),

			From_: G(NodeBlockCont, "from_"),

			NameAlias: G(NodeBlockCont, "nameAlias"),

			To_: G(NodeBlockCont, "to_"),
		},
	}
}

func NodeBlockFromContainerMG(cont *container.Container) *NodeBlockMG {

	return NodeBlockFromContainerListMG(cont, 0)
}

func NodeBlockListFromContainerMG(cont *container.Container) []*NodeBlockMG {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*NodeBlockMG, length)

	for i := 0; i < length; i++ {

		arr[i] = NodeBlockFromContainerListMG(cont, i)
	}

	return arr
}

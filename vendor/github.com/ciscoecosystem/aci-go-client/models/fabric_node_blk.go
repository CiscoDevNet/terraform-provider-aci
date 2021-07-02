package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FabricNodeBlkClassName = "fabricNodeBlk"

type NodeBlk struct {
	BaseAttributes
	NodeBlkAttributes
}

type NodeBlkAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
	From_      string `json:",omitempty"`
	To_        string `json:",omitempty"`
}

func NewNodeBlk(fabricNodeBlkRn, parentDn, description string, fabricNodeBlkattr NodeBlkAttributes) *NodeBlk {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricNodeBlkRn)
	return &NodeBlk{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "",
			ClassName:         FabricNodeBlkClassName,
			Rn:                fabricNodeBlkRn,
		},

		NodeBlkAttributes: fabricNodeBlkattr,
	}
}

func (fabricNodeBlk *NodeBlk) ToMap() (map[string]string, error) {
	fabricNodeBlkMap, err := fabricNodeBlk.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fabricNodeBlkMap, "annotation", fabricNodeBlk.Annotation)
	A(fabricNodeBlkMap, "name", fabricNodeBlk.Name)
	A(fabricNodeBlkMap, "nameAlias", fabricNodeBlk.NameAlias)
	A(fabricNodeBlkMap, "from_", fabricNodeBlk.From_)
	A(fabricNodeBlkMap, "to_", fabricNodeBlk.To_)

	return fabricNodeBlkMap, err
}

func NodeBlkFromContainerList(cont *container.Container, index int) *NodeBlk {

	NodeBlkCont := cont.S("imdata").Index(index).S(FabricNodeBlkClassName, "attributes")
	return &NodeBlk{
		BaseAttributes{
			DistinguishedName: G(NodeBlkCont, "dn"),
			Description:       G(NodeBlkCont, "descr"),
			Status:            G(NodeBlkCont, "status"),
			ClassName:         FabricNodeBlkClassName,
			Rn:                G(NodeBlkCont, "rn"),
		},

		NodeBlkAttributes{
			Annotation: G(NodeBlkCont, "annotation"),
			Name:       G(NodeBlkCont, "name"),
			NameAlias:  G(NodeBlkCont, "nameAlias"),
			From_:      G(NodeBlkCont, "from_"),
			To_:        G(NodeBlkCont, "to_"),
		},
	}
}

func NodeBlkFromContainer(cont *container.Container) *NodeBlk {

	return NodeBlkFromContainerList(cont, 0)
}

func NodeBlkListFromContainer(cont *container.Container) []*NodeBlk {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*NodeBlk, length)

	for i := 0; i < length; i++ {

		arr[i] = NodeBlkFromContainerList(cont, i)
	}

	return arr
}

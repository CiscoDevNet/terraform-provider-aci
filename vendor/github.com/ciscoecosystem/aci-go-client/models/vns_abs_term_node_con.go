package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VnsabstermnodeconClassName = "vnsAbsTermNodeCon"

type ConsumerTerminalNode struct {
	BaseAttributes
	ConsumerTerminalNodeAttributes
}

type ConsumerTerminalNodeAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewConsumerTerminalNode(vnsAbsTermNodeConRn, parentDn, description string, vnsAbsTermNodeConattr ConsumerTerminalNodeAttributes) *ConsumerTerminalNode {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsAbsTermNodeConRn)
	return &ConsumerTerminalNode{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsabstermnodeconClassName,
			Rn:                vnsAbsTermNodeConRn,
		},

		ConsumerTerminalNodeAttributes: vnsAbsTermNodeConattr,
	}
}

func (vnsAbsTermNodeCon *ConsumerTerminalNode) ToMap() (map[string]string, error) {
	vnsAbsTermNodeConMap, err := vnsAbsTermNodeCon.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vnsAbsTermNodeConMap, "name", vnsAbsTermNodeCon.Name)

	A(vnsAbsTermNodeConMap, "annotation", vnsAbsTermNodeCon.Annotation)

	A(vnsAbsTermNodeConMap, "nameAlias", vnsAbsTermNodeCon.NameAlias)

	return vnsAbsTermNodeConMap, err
}

func ConsumerTerminalNodeFromContainerList(cont *container.Container, index int) *ConsumerTerminalNode {

	ConsumerTerminalNodeCont := cont.S("imdata").Index(index).S(VnsabstermnodeconClassName, "attributes")
	return &ConsumerTerminalNode{
		BaseAttributes{
			DistinguishedName: G(ConsumerTerminalNodeCont, "dn"),
			Description:       G(ConsumerTerminalNodeCont, "descr"),
			Status:            G(ConsumerTerminalNodeCont, "status"),
			ClassName:         VnsabstermnodeconClassName,
			Rn:                G(ConsumerTerminalNodeCont, "rn"),
		},

		ConsumerTerminalNodeAttributes{

			Name: G(ConsumerTerminalNodeCont, "name"),

			Annotation: G(ConsumerTerminalNodeCont, "annotation"),

			NameAlias: G(ConsumerTerminalNodeCont, "nameAlias"),
		},
	}
}

func ConsumerTerminalNodeFromContainer(cont *container.Container) *ConsumerTerminalNode {

	return ConsumerTerminalNodeFromContainerList(cont, 0)
}

func ConsumerTerminalNodeListFromContainer(cont *container.Container) []*ConsumerTerminalNode {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ConsumerTerminalNode, length)

	for i := 0; i < length; i++ {

		arr[i] = ConsumerTerminalNodeFromContainerList(cont, i)
	}

	return arr
}

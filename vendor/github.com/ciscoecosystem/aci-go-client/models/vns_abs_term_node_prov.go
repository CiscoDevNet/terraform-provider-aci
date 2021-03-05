package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VnsabstermnodeprovClassName = "vnsAbsTermNodeProv"

type ProviderTerminalNode struct {
	BaseAttributes
	ProviderTerminalNodeAttributes
}

type ProviderTerminalNodeAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewProviderTerminalNode(vnsAbsTermNodeProvRn, parentDn, description string, vnsAbsTermNodeProvattr ProviderTerminalNodeAttributes) *ProviderTerminalNode {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsAbsTermNodeProvRn)
	return &ProviderTerminalNode{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsabstermnodeprovClassName,
			Rn:                vnsAbsTermNodeProvRn,
		},

		ProviderTerminalNodeAttributes: vnsAbsTermNodeProvattr,
	}
}

func (vnsAbsTermNodeProv *ProviderTerminalNode) ToMap() (map[string]string, error) {
	vnsAbsTermNodeProvMap, err := vnsAbsTermNodeProv.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vnsAbsTermNodeProvMap, "name", vnsAbsTermNodeProv.Name)

	A(vnsAbsTermNodeProvMap, "annotation", vnsAbsTermNodeProv.Annotation)

	A(vnsAbsTermNodeProvMap, "nameAlias", vnsAbsTermNodeProv.NameAlias)

	return vnsAbsTermNodeProvMap, err
}

func ProviderTerminalNodeFromContainerList(cont *container.Container, index int) *ProviderTerminalNode {

	ProviderTerminalNodeCont := cont.S("imdata").Index(index).S(VnsabstermnodeprovClassName, "attributes")
	return &ProviderTerminalNode{
		BaseAttributes{
			DistinguishedName: G(ProviderTerminalNodeCont, "dn"),
			Description:       G(ProviderTerminalNodeCont, "descr"),
			Status:            G(ProviderTerminalNodeCont, "status"),
			ClassName:         VnsabstermnodeprovClassName,
			Rn:                G(ProviderTerminalNodeCont, "rn"),
		},

		ProviderTerminalNodeAttributes{

			Name: G(ProviderTerminalNodeCont, "name"),

			Annotation: G(ProviderTerminalNodeCont, "annotation"),

			NameAlias: G(ProviderTerminalNodeCont, "nameAlias"),
		},
	}
}

func ProviderTerminalNodeFromContainer(cont *container.Container) *ProviderTerminalNode {

	return ProviderTerminalNodeFromContainerList(cont, 0)
}

func ProviderTerminalNodeListFromContainer(cont *container.Container) []*ProviderTerminalNode {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ProviderTerminalNode, length)

	for i := 0; i < length; i++ {

		arr[i] = ProviderTerminalNodeFromContainerList(cont, i)
	}

	return arr
}

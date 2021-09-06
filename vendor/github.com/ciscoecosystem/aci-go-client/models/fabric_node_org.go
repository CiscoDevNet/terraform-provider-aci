package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FabricnodeClassName = "fabricNode"

type OrgFabricNode struct {
	BaseAttributes
	OrgFabricNodeAttributes
}

type OrgFabricNodeAttributes struct {
	Fabric_node_id string `json:",omitempty"`

	AdSt string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ApicType string `json:",omitempty"`

	FabricSt string `json:",omitempty"`

	LastStateModTs string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func OrgNewFabricNode(fabricNodeRn, parentDn, description string, fabricNodeattr OrgFabricNodeAttributes) *OrgFabricNode {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricNodeRn)
	return &OrgFabricNode{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FabricnodeClassName,
			Rn:                fabricNodeRn,
		},

		OrgFabricNodeAttributes: fabricNodeattr,
	}
}

func (fabricNode *OrgFabricNode) ToMap() (map[string]string, error) {
	fabricNodeMap, err := fabricNode.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fabricNodeMap, "id", fabricNode.Fabric_node_id)

	A(fabricNodeMap, "adSt", fabricNode.AdSt)

	A(fabricNodeMap, "annotation", fabricNode.Annotation)

	A(fabricNodeMap, "apicType", fabricNode.ApicType)

	A(fabricNodeMap, "fabricSt", fabricNode.FabricSt)

	A(fabricNodeMap, "lastStateModTs", fabricNode.LastStateModTs)

	A(fabricNodeMap, "nameAlias", fabricNode.NameAlias)

	return fabricNodeMap, err
}

func OrgFabricNodeFromContainerList(cont *container.Container, index int) *OrgFabricNode {

	FabricNodeCont := cont.S("imdata").Index(index).S(FabricnodeClassName, "attributes")
	return &OrgFabricNode{
		BaseAttributes{
			DistinguishedName: G(FabricNodeCont, "dn"),
			Description:       G(FabricNodeCont, "descr"),
			Status:            G(FabricNodeCont, "status"),
			ClassName:         FabricnodeClassName,
			Rn:                G(FabricNodeCont, "rn"),
		},

		OrgFabricNodeAttributes{

			Fabric_node_id: G(FabricNodeCont, "id"),

			AdSt: G(FabricNodeCont, "adSt"),

			Annotation: G(FabricNodeCont, "annotation"),

			ApicType: G(FabricNodeCont, "apicType"),

			FabricSt: G(FabricNodeCont, "fabricSt"),

			LastStateModTs: G(FabricNodeCont, "lastStateModTs"),

			NameAlias: G(FabricNodeCont, "nameAlias"),
		},
	}
}

func OrgFabricNodeFromContainer(cont *container.Container) *OrgFabricNode {

	return OrgFabricNodeFromContainerList(cont, 0)
}

func OrgFabricNodeListFromContainer(cont *container.Container) []*OrgFabricNode {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*OrgFabricNode, length)

	for i := 0; i < length; i++ {

		arr[i] = OrgFabricNodeFromContainerList(cont, i)
	}

	return arr
}

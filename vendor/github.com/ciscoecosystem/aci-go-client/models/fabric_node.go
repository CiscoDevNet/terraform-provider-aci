package models

import (
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const TopologyFabricNodeClassName = "fabricNode"

type TopologyFabricNode struct {
	BaseAttributes
	TopologyFabricNodeAttributes
}

type TopologyFabricNodeAttributes struct {
	AdSt             string `json:",omitempty"`
	Address          string `json:",omitempty"`
	Annotation       string `json:",omitempty"`
	ApicType         string `json:",omitempty"`
	DelayedHeartbeat string `json:",omitempty"`
	ExtMngdBy        string `json:",omitempty"`
	FabricSt         string `json:",omitempty"`
	Id               string `json:",omitempty"`
	LastStateModTs   string `json:",omitempty"`
	ModTs            string `json:",omitempty"`
	Model            string `json:",omitempty"`
	MonPolDn         string `json:",omitempty"`
	Name             string `json:",omitempty"`
	NameAlias        string `json:",omitempty"`
	NodeType         string `json:",omitempty"`
	Role             string `json:",omitempty"`
	Serial           string `json:",omitempty"`
	Uid              string `json:",omitempty"`
	Userdom          string `json:",omitempty"`
	Vendor           string `json:",omitempty"`
	Version          string `json:",omitempty"`
}

/*
 * No NewTopologyFabricNode as this is a non-configurable MO
func NewTopologyFabricNode(fabricNodeRn, parentDn, description string, fabricNodeattr TopologyFabricNodeAttributes) *TopologyFabricNode {
    dn := fmt.Sprintf("%s/%s", parentDn, fabricNodeRn)
    return &TopologyFabricNode{
        BaseAttributes: BaseAttributes{
            DistinguishedName: dn,
            Description:       description,
            Status:            "created, modified",
            ClassName:         TopologyFabricNodeClassName,
            Rn:                fabricNodeRn,
        },

        TopologyFabricNodeAttributes: fabricNodeattr,

    }
}
*/

func (fabricNode *TopologyFabricNode) ToMap() (map[string]string, error) {
	fabricNodeMap, err := fabricNode.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fabricNodeMap, "adSt", fabricNode.AdSt)
	A(fabricNodeMap, "address", fabricNode.Address)
	A(fabricNodeMap, "annotation", fabricNode.Annotation)
	A(fabricNodeMap, "apicType", fabricNode.ApicType)
	A(fabricNodeMap, "delayedHeartbeat", fabricNode.DelayedHeartbeat)
	A(fabricNodeMap, "extMngdBy", fabricNode.ExtMngdBy)
	A(fabricNodeMap, "fabricSt", fabricNode.FabricSt)
	A(fabricNodeMap, "id", fabricNode.Id)
	A(fabricNodeMap, "lastStateModTs", fabricNode.LastStateModTs)
	A(fabricNodeMap, "modTs", fabricNode.ModTs)
	A(fabricNodeMap, "model", fabricNode.Model)
	A(fabricNodeMap, "monPolDn", fabricNode.MonPolDn)
	A(fabricNodeMap, "name", fabricNode.Name)
	A(fabricNodeMap, "nameAlias", fabricNode.NameAlias)
	A(fabricNodeMap, "nodeType", fabricNode.NodeType)
	A(fabricNodeMap, "role", fabricNode.Role)
	A(fabricNodeMap, "serial", fabricNode.Serial)
	A(fabricNodeMap, "uid", fabricNode.Uid)
	A(fabricNodeMap, "userdom", fabricNode.Userdom)
	A(fabricNodeMap, "vendor", fabricNode.Vendor)
	A(fabricNodeMap, "version", fabricNode.Version)

	return fabricNodeMap, err
}

func TopologyFabricNodeFromContainerList(cont *container.Container, index int) *TopologyFabricNode {

	TopologyFabricNodeCont := cont.S("imdata").Index(index).S(TopologyFabricNodeClassName, "attributes")
	return &TopologyFabricNode{
		BaseAttributes{
			DistinguishedName: G(TopologyFabricNodeCont, "dn"),
			Description:       G(TopologyFabricNodeCont, "descr"),
			Status:            G(TopologyFabricNodeCont, "status"),
			ClassName:         TopologyFabricNodeClassName,
			Rn:                G(TopologyFabricNodeCont, "rn"),
		},

		TopologyFabricNodeAttributes{
			AdSt:             G(TopologyFabricNodeCont, "adSt"),
			Address:          G(TopologyFabricNodeCont, "address"),
			Annotation:       G(TopologyFabricNodeCont, "annotation"),
			ApicType:         G(TopologyFabricNodeCont, "apicType"),
			DelayedHeartbeat: G(TopologyFabricNodeCont, "delayedHeartbeat"),
			ExtMngdBy:        G(TopologyFabricNodeCont, "extMngdBy"),
			FabricSt:         G(TopologyFabricNodeCont, "fabricSt"),
			Id:               G(TopologyFabricNodeCont, "id"),
			LastStateModTs:   G(TopologyFabricNodeCont, "lastStateModTs"),
			ModTs:            G(TopologyFabricNodeCont, "modTs"),
			Model:            G(TopologyFabricNodeCont, "model"),
			MonPolDn:         G(TopologyFabricNodeCont, "monPolDn"),
			Name:             G(TopologyFabricNodeCont, "name"),
			NameAlias:        G(TopologyFabricNodeCont, "nameAlias"),
			NodeType:         G(TopologyFabricNodeCont, "nodeType"),
			Role:             G(TopologyFabricNodeCont, "role"),
			Serial:           G(TopologyFabricNodeCont, "serial"),
			Uid:              G(TopologyFabricNodeCont, "uid"),
			Userdom:          G(TopologyFabricNodeCont, "userdom"),
			Vendor:           G(TopologyFabricNodeCont, "vendor"),
			Version:          G(TopologyFabricNodeCont, "version"),
		},
	}
}

func TopologyFabricNodeFromContainer(cont *container.Container) *TopologyFabricNode {

	return TopologyFabricNodeFromContainerList(cont, 0)
}

func TopologyFabricNodeListFromContainer(cont *container.Container) []*TopologyFabricNode {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*TopologyFabricNode, length)

	for i := 0; i < length; i++ {

		arr[i] = TopologyFabricNodeFromContainerList(cont, i)
	}

	return arr
}

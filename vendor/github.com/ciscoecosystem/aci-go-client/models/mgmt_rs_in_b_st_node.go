package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const MgmtrsinbstnodeClassName = "mgmtRsInBStNode"

type InbandStaticNode struct {
	BaseAttributes
	InbandStaticNodeAttributes
}

type InbandStaticNodeAttributes struct {
	TDn string `json:",omitempty"`

	Addr string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Gw string `json:",omitempty"`

	V6Addr string `json:",omitempty"`

	V6Gw string `json:",omitempty"`
}

func NewInbandStaticNode(mgmtRsInBStNodeRn, parentDn, description string, mgmtRsInBStNodeattr InbandStaticNodeAttributes) *InbandStaticNode {
	dn := fmt.Sprintf("%s/%s", parentDn, mgmtRsInBStNodeRn)
	return &InbandStaticNode{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         MgmtrsinbstnodeClassName,
			Rn:                mgmtRsInBStNodeRn,
		},

		InbandStaticNodeAttributes: mgmtRsInBStNodeattr,
	}
}

func (mgmtRsInBStNode *InbandStaticNode) ToMap() (map[string]string, error) {
	mgmtRsInBStNodeMap, err := mgmtRsInBStNode.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(mgmtRsInBStNodeMap, "tDn", mgmtRsInBStNode.TDn)

	A(mgmtRsInBStNodeMap, "addr", mgmtRsInBStNode.Addr)

	A(mgmtRsInBStNodeMap, "annotation", mgmtRsInBStNode.Annotation)

	A(mgmtRsInBStNodeMap, "gw", mgmtRsInBStNode.Gw)

	A(mgmtRsInBStNodeMap, "v6Addr", mgmtRsInBStNode.V6Addr)

	A(mgmtRsInBStNodeMap, "v6Gw", mgmtRsInBStNode.V6Gw)

	return mgmtRsInBStNodeMap, err
}

func InbandStaticNodeFromContainerList(cont *container.Container, index int) *InbandStaticNode {

	InbandStaticNodeCont := cont.S("imdata").Index(index).S(MgmtrsinbstnodeClassName, "attributes")
	return &InbandStaticNode{
		BaseAttributes{
			DistinguishedName: G(InbandStaticNodeCont, "dn"),
			Description:       G(InbandStaticNodeCont, "descr"),
			Status:            G(InbandStaticNodeCont, "status"),
			ClassName:         MgmtrsinbstnodeClassName,
			Rn:                G(InbandStaticNodeCont, "rn"),
		},

		InbandStaticNodeAttributes{

			TDn: G(InbandStaticNodeCont, "tDn"),

			Addr: G(InbandStaticNodeCont, "addr"),

			Annotation: G(InbandStaticNodeCont, "annotation"),

			Gw: G(InbandStaticNodeCont, "gw"),

			V6Addr: G(InbandStaticNodeCont, "v6Addr"),

			V6Gw: G(InbandStaticNodeCont, "v6Gw"),
		},
	}
}

func InbandStaticNodeFromContainer(cont *container.Container) *InbandStaticNode {

	return InbandStaticNodeFromContainerList(cont, 0)
}

func InbandStaticNodeListFromContainer(cont *container.Container) []*InbandStaticNode {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*InbandStaticNode, length)

	for i := 0; i < length; i++ {

		arr[i] = InbandStaticNodeFromContainerList(cont, i)
	}

	return arr
}

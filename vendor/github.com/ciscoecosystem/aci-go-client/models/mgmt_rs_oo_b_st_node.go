package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const MgmtrsoobstnodeClassName = "mgmtRsOoBStNode"

type OutofbandStaticNode struct {
	BaseAttributes
	OutofbandStaticNodeAttributes
}

type OutofbandStaticNodeAttributes struct {
	TDn string `json:",omitempty"`

	Addr string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Gw string `json:",omitempty"`

	V6Addr string `json:",omitempty"`

	V6Gw string `json:",omitempty"`
}

func NewOutofbandStaticNode(mgmtRsOoBStNodeRn, parentDn, description string, mgmtRsOoBStNodeattr OutofbandStaticNodeAttributes) *OutofbandStaticNode {
	dn := fmt.Sprintf("%s/%s", parentDn, mgmtRsOoBStNodeRn)
	return &OutofbandStaticNode{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         MgmtrsoobstnodeClassName,
			Rn:                mgmtRsOoBStNodeRn,
		},

		OutofbandStaticNodeAttributes: mgmtRsOoBStNodeattr,
	}
}

func (mgmtRsOoBStNode *OutofbandStaticNode) ToMap() (map[string]string, error) {
	mgmtRsOoBStNodeMap, err := mgmtRsOoBStNode.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(mgmtRsOoBStNodeMap, "tDn", mgmtRsOoBStNode.TDn)

	A(mgmtRsOoBStNodeMap, "addr", mgmtRsOoBStNode.Addr)

	A(mgmtRsOoBStNodeMap, "annotation", mgmtRsOoBStNode.Annotation)

	A(mgmtRsOoBStNodeMap, "gw", mgmtRsOoBStNode.Gw)

	A(mgmtRsOoBStNodeMap, "v6Addr", mgmtRsOoBStNode.V6Addr)

	A(mgmtRsOoBStNodeMap, "v6Gw", mgmtRsOoBStNode.V6Gw)

	return mgmtRsOoBStNodeMap, err
}

func OutofbandStaticNodeFromContainerList(cont *container.Container, index int) *OutofbandStaticNode {

	OutofbandStaticNodeCont := cont.S("imdata").Index(index).S(MgmtrsoobstnodeClassName, "attributes")
	return &OutofbandStaticNode{
		BaseAttributes{
			DistinguishedName: G(OutofbandStaticNodeCont, "dn"),
			Description:       G(OutofbandStaticNodeCont, "descr"),
			Status:            G(OutofbandStaticNodeCont, "status"),
			ClassName:         MgmtrsoobstnodeClassName,
			Rn:                G(OutofbandStaticNodeCont, "rn"),
		},

		OutofbandStaticNodeAttributes{

			TDn: G(OutofbandStaticNodeCont, "tDn"),

			Addr: G(OutofbandStaticNodeCont, "addr"),

			Annotation: G(OutofbandStaticNodeCont, "annotation"),

			Gw: G(OutofbandStaticNodeCont, "gw"),

			V6Addr: G(OutofbandStaticNodeCont, "v6Addr"),

			V6Gw: G(OutofbandStaticNodeCont, "v6Gw"),
		},
	}
}

func OutofbandStaticNodeFromContainer(cont *container.Container) *OutofbandStaticNode {

	return OutofbandStaticNodeFromContainerList(cont, 0)
}

func OutofbandStaticNodeListFromContainer(cont *container.Container) []*OutofbandStaticNode {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*OutofbandStaticNode, length)

	for i := 0; i < length; i++ {

		arr[i] = OutofbandStaticNodeFromContainerList(cont, i)
	}

	return arr
}

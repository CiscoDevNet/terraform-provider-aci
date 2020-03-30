package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FabricnodepepClassName = "fabricNodePEp"

type NodePolicyEndPoint struct {
	BaseAttributes
	NodePolicyEndPointAttributes
}

type NodePolicyEndPointAttributes struct {
	Node_policy_end_point_id string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	PodId string `json:",omitempty"`
}

func NewNodePolicyEndPoint(fabricNodePEpRn, parentDn, description string, fabricNodePEpattr NodePolicyEndPointAttributes) *NodePolicyEndPoint {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricNodePEpRn)
	return &NodePolicyEndPoint{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FabricnodepepClassName,
			Rn:                fabricNodePEpRn,
		},

		NodePolicyEndPointAttributes: fabricNodePEpattr,
	}
}

func (fabricNodePEp *NodePolicyEndPoint) ToMap() (map[string]string, error) {
	fabricNodePEpMap, err := fabricNodePEp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fabricNodePEpMap, "id", fabricNodePEp.Node_policy_end_point_id)

	A(fabricNodePEpMap, "annotation", fabricNodePEp.Annotation)

	A(fabricNodePEpMap, "nameAlias", fabricNodePEp.NameAlias)

	A(fabricNodePEpMap, "podId", fabricNodePEp.PodId)

	return fabricNodePEpMap, err
}

func NodePolicyEndPointFromContainerList(cont *container.Container, index int) *NodePolicyEndPoint {

	NodePolicyEndPointCont := cont.S("imdata").Index(index).S(FabricnodepepClassName, "attributes")
	return &NodePolicyEndPoint{
		BaseAttributes{
			DistinguishedName: G(NodePolicyEndPointCont, "dn"),
			Description:       G(NodePolicyEndPointCont, "descr"),
			Status:            G(NodePolicyEndPointCont, "status"),
			ClassName:         FabricnodepepClassName,
			Rn:                G(NodePolicyEndPointCont, "rn"),
		},

		NodePolicyEndPointAttributes{

			Node_policy_end_point_id: G(NodePolicyEndPointCont, "id"),

			Annotation: G(NodePolicyEndPointCont, "annotation"),

			NameAlias: G(NodePolicyEndPointCont, "nameAlias"),

			PodId: G(NodePolicyEndPointCont, "podId"),
		},
	}
}

func NodePolicyEndPointFromContainer(cont *container.Container) *NodePolicyEndPoint {

	return NodePolicyEndPointFromContainerList(cont, 0)
}

func NodePolicyEndPointListFromContainer(cont *container.Container) []*NodePolicyEndPoint {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*NodePolicyEndPoint, length)

	for i := 0; i < length; i++ {

		arr[i] = NodePolicyEndPointFromContainerList(cont, i)
	}

	return arr
}

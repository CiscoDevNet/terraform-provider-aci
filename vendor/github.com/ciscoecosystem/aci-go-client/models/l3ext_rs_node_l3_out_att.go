package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extrsnodel3outattClassName = "l3extRsNodeL3OutAtt"

type FabricNode struct {
	BaseAttributes
	FabricNodeAttributes
}

type FabricNodeAttributes struct {
	TDn string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ConfigIssues string `json:",omitempty"`

	RtrId string `json:",omitempty"`

	RtrIdLoopBack string `json:",omitempty"`
}

func NewFabricNode(l3extRsNodeL3OutAttRn, parentDn, description string, l3extRsNodeL3OutAttattr FabricNodeAttributes) *FabricNode {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extRsNodeL3OutAttRn)
	return &FabricNode{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extrsnodel3outattClassName,
			Rn:                l3extRsNodeL3OutAttRn,
		},

		FabricNodeAttributes: l3extRsNodeL3OutAttattr,
	}
}

func (l3extRsNodeL3OutAtt *FabricNode) ToMap() (map[string]string, error) {
	l3extRsNodeL3OutAttMap, err := l3extRsNodeL3OutAtt.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extRsNodeL3OutAttMap, "tDn", l3extRsNodeL3OutAtt.TDn)

	A(l3extRsNodeL3OutAttMap, "annotation", l3extRsNodeL3OutAtt.Annotation)

	A(l3extRsNodeL3OutAttMap, "configIssues", l3extRsNodeL3OutAtt.ConfigIssues)

	A(l3extRsNodeL3OutAttMap, "rtrId", l3extRsNodeL3OutAtt.RtrId)

	A(l3extRsNodeL3OutAttMap, "rtrIdLoopBack", l3extRsNodeL3OutAtt.RtrIdLoopBack)

	return l3extRsNodeL3OutAttMap, err
}

func FabricNodeFromContainerList(cont *container.Container, index int) *FabricNode {

	FabricNodeCont := cont.S("imdata").Index(index).S(L3extrsnodel3outattClassName, "attributes")
	return &FabricNode{
		BaseAttributes{
			DistinguishedName: G(FabricNodeCont, "dn"),
			Description:       G(FabricNodeCont, "descr"),
			Status:            G(FabricNodeCont, "status"),
			ClassName:         L3extrsnodel3outattClassName,
			Rn:                G(FabricNodeCont, "rn"),
		},

		FabricNodeAttributes{

			TDn: G(FabricNodeCont, "tDn"),

			Annotation: G(FabricNodeCont, "annotation"),

			ConfigIssues: G(FabricNodeCont, "configIssues"),

			RtrId: G(FabricNodeCont, "rtrId"),

			RtrIdLoopBack: G(FabricNodeCont, "rtrIdLoopBack"),
		},
	}
}

func FabricNodeFromContainer(cont *container.Container) *FabricNode {

	return FabricNodeFromContainerList(cont, 0)
}

func FabricNodeListFromContainer(cont *container.Container) []*FabricNode {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FabricNode, length)

	for i := 0; i < length; i++ {

		arr[i] = FabricNodeFromContainerList(cont, i)
	}

	return arr
}

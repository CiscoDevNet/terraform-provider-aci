package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extlnodepClassName = "l3extLNodeP"

type LogicalNodeProfile struct {
	BaseAttributes
	LogicalNodeProfileAttributes
}

type LogicalNodeProfileAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ConfigIssues string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Tag string `json:",omitempty"`

	TargetDscp string `json:",omitempty"`
}

func NewLogicalNodeProfile(l3extLNodePRn, parentDn, description string, l3extLNodePattr LogicalNodeProfileAttributes) *LogicalNodeProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extLNodePRn)
	return &LogicalNodeProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extlnodepClassName,
			Rn:                l3extLNodePRn,
		},

		LogicalNodeProfileAttributes: l3extLNodePattr,
	}
}

func (l3extLNodeP *LogicalNodeProfile) ToMap() (map[string]string, error) {
	l3extLNodePMap, err := l3extLNodeP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extLNodePMap, "name", l3extLNodeP.Name)

	A(l3extLNodePMap, "annotation", l3extLNodeP.Annotation)

	A(l3extLNodePMap, "configIssues", l3extLNodeP.ConfigIssues)

	A(l3extLNodePMap, "nameAlias", l3extLNodeP.NameAlias)

	A(l3extLNodePMap, "tag", l3extLNodeP.Tag)

	A(l3extLNodePMap, "targetDscp", l3extLNodeP.TargetDscp)

	return l3extLNodePMap, err
}

func LogicalNodeProfileFromContainerList(cont *container.Container, index int) *LogicalNodeProfile {

	LogicalNodeProfileCont := cont.S("imdata").Index(index).S(L3extlnodepClassName, "attributes")
	return &LogicalNodeProfile{
		BaseAttributes{
			DistinguishedName: G(LogicalNodeProfileCont, "dn"),
			Description:       G(LogicalNodeProfileCont, "descr"),
			Status:            G(LogicalNodeProfileCont, "status"),
			ClassName:         L3extlnodepClassName,
			Rn:                G(LogicalNodeProfileCont, "rn"),
		},

		LogicalNodeProfileAttributes{

			Name: G(LogicalNodeProfileCont, "name"),

			Annotation: G(LogicalNodeProfileCont, "annotation"),

			ConfigIssues: G(LogicalNodeProfileCont, "configIssues"),

			NameAlias: G(LogicalNodeProfileCont, "nameAlias"),

			Tag: G(LogicalNodeProfileCont, "tag"),

			TargetDscp: G(LogicalNodeProfileCont, "targetDscp"),
		},
	}
}

func LogicalNodeProfileFromContainer(cont *container.Container) *LogicalNodeProfile {

	return LogicalNodeProfileFromContainerList(cont, 0)
}

func LogicalNodeProfileListFromContainer(cont *container.Container) []*LogicalNodeProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LogicalNodeProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = LogicalNodeProfileFromContainerList(cont, i)
	}

	return arr
}

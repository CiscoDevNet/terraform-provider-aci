package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnvnsLIf        = "%s/lIf-%s"
	RnvnsLIf        = "lIf-%s"
	VnslifClassName = "vnsLIf"
)

type LogicalInterface struct {
	BaseAttributes
	NameAliasAttribute
	LogicalInterfaceAttributes
}

type LogicalInterfaceAttributes struct {
	Encap         string `json:",omitempty"`
	LagPolicyName string `json:",omitempty"`
	Name          string `json:",omitempty"`
	Annotation    string `json:",omitempty"`
}

func NewLogicalInterface(vnsLIfRn, parentDn, nameAlias string, vnsLIfAttr LogicalInterfaceAttributes) *LogicalInterface {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsLIfRn)
	return &LogicalInterface{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         VnslifClassName,
			Rn:                vnsLIfRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		LogicalInterfaceAttributes: vnsLIfAttr,
	}
}

func (vnsLIf *LogicalInterface) ToMap() (map[string]string, error) {
	vnsLIfMap, err := vnsLIf.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := vnsLIf.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(vnsLIfMap, key, value)
	}

	A(vnsLIfMap, "encap", vnsLIf.Encap)
	A(vnsLIfMap, "lagPolicyName", vnsLIf.LagPolicyName)
	A(vnsLIfMap, "name", vnsLIf.Name)
	A(vnsLIfMap, "annotation", vnsLIf.Annotation)
	return vnsLIfMap, err
}

func LogicalInterfaceFromContainerList(cont *container.Container, index int) *LogicalInterface {
	LogicalInterfaceCont := cont.S("imdata").Index(index).S(VnslifClassName, "attributes")
	return &LogicalInterface{
		BaseAttributes{
			DistinguishedName: G(LogicalInterfaceCont, "dn"),
			Status:            G(LogicalInterfaceCont, "status"),
			ClassName:         VnslifClassName,
			Rn:                G(LogicalInterfaceCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(LogicalInterfaceCont, "nameAlias"),
		},
		LogicalInterfaceAttributes{
			Encap:         G(LogicalInterfaceCont, "encap"),
			LagPolicyName: G(LogicalInterfaceCont, "lagPolicyName"),
			Name:          G(LogicalInterfaceCont, "name"),
			Annotation:    G(LogicalInterfaceCont, "annotation"),
		},
	}
}

func LogicalInterfaceFromContainer(cont *container.Container) *LogicalInterface {
	return LogicalInterfaceFromContainerList(cont, 0)
}

func LogicalInterfaceListFromContainer(cont *container.Container) []*LogicalInterface {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*LogicalInterface, length)

	for i := 0; i < length; i++ {
		arr[i] = LogicalInterfaceFromContainerList(cont, i)
	}

	return arr
}

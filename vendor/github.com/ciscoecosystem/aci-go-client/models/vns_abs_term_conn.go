package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VnsabstermconnClassName = "vnsAbsTermConn"

type TerminalConnector struct {
	BaseAttributes
	TerminalConnectorAttributes
}

type TerminalConnectorAttributes struct {
	Annotation string `json:",omitempty"`

	AttNotify string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewTerminalConnector(vnsAbsTermConnRn, parentDn, description string, vnsAbsTermConnattr TerminalConnectorAttributes) *TerminalConnector {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsAbsTermConnRn)
	return &TerminalConnector{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsabstermconnClassName,
			Rn:                vnsAbsTermConnRn,
		},

		TerminalConnectorAttributes: vnsAbsTermConnattr,
	}
}

func (vnsAbsTermConn *TerminalConnector) ToMap() (map[string]string, error) {
	vnsAbsTermConnMap, err := vnsAbsTermConn.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vnsAbsTermConnMap, "annotation", vnsAbsTermConn.Annotation)

	A(vnsAbsTermConnMap, "attNotify", vnsAbsTermConn.AttNotify)

	A(vnsAbsTermConnMap, "nameAlias", vnsAbsTermConn.NameAlias)

	return vnsAbsTermConnMap, err
}

func TerminalConnectorFromContainerList(cont *container.Container, index int) *TerminalConnector {

	TerminalConnectorCont := cont.S("imdata").Index(index).S(VnsabstermconnClassName, "attributes")
	return &TerminalConnector{
		BaseAttributes{
			DistinguishedName: G(TerminalConnectorCont, "dn"),
			Description:       G(TerminalConnectorCont, "descr"),
			Status:            G(TerminalConnectorCont, "status"),
			ClassName:         VnsabstermconnClassName,
			Rn:                G(TerminalConnectorCont, "rn"),
		},

		TerminalConnectorAttributes{

			Annotation: G(TerminalConnectorCont, "annotation"),

			AttNotify: G(TerminalConnectorCont, "attNotify"),

			NameAlias: G(TerminalConnectorCont, "nameAlias"),
		},
	}
}

func TerminalConnectorFromContainer(cont *container.Container) *TerminalConnector {

	return TerminalConnectorFromContainerList(cont, 0)
}

func TerminalConnectorListFromContainer(cont *container.Container) []*TerminalConnector {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*TerminalConnector, length)

	for i := 0; i < length; i++ {

		arr[i] = TerminalConnectorFromContainerList(cont, i)
	}

	return arr
}

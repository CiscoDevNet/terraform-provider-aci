package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const VnsabsfuncconnClassName = "vnsAbsFuncConn"

type FunctionConnector struct {
	BaseAttributes
	FunctionConnectorAttributes
}

type FunctionConnectorAttributes struct {
	Name          string `json:",omitempty"`
	Annotation    string `json:",omitempty"`
	AttNotify     string `json:",omitempty"`
	ConnType      string `json:",omitempty"`
	NameAlias     string `json:",omitempty"`
	DeviceLIfName string `json:",omitempty"`
}

func NewFunctionConnector(vnsAbsFuncConnRn, parentDn, description string, vnsAbsFuncConnattr FunctionConnectorAttributes) *FunctionConnector {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsAbsFuncConnRn)
	return &FunctionConnector{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsabsfuncconnClassName,
			Rn:                vnsAbsFuncConnRn,
		},

		FunctionConnectorAttributes: vnsAbsFuncConnattr,
	}
}

func (vnsAbsFuncConn *FunctionConnector) ToMap() (map[string]string, error) {
	vnsAbsFuncConnMap, err := vnsAbsFuncConn.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vnsAbsFuncConnMap, "name", vnsAbsFuncConn.Name)
	A(vnsAbsFuncConnMap, "annotation", vnsAbsFuncConn.Annotation)
	A(vnsAbsFuncConnMap, "attNotify", vnsAbsFuncConn.AttNotify)
	A(vnsAbsFuncConnMap, "connType", vnsAbsFuncConn.ConnType)
	A(vnsAbsFuncConnMap, "nameAlias", vnsAbsFuncConn.NameAlias)
	A(vnsAbsFuncConnMap, "deviceLIfName", vnsAbsFuncConn.DeviceLIfName)

	return vnsAbsFuncConnMap, err
}

func FunctionConnectorFromContainerList(cont *container.Container, index int) *FunctionConnector {

	FunctionConnectorCont := cont.S("imdata").Index(index).S(VnsabsfuncconnClassName, "attributes")
	return &FunctionConnector{
		BaseAttributes{
			DistinguishedName: G(FunctionConnectorCont, "dn"),
			Description:       G(FunctionConnectorCont, "descr"),
			Status:            G(FunctionConnectorCont, "status"),
			ClassName:         VnsabsfuncconnClassName,
			Rn:                G(FunctionConnectorCont, "rn"),
		},

		FunctionConnectorAttributes{
			Name:          G(FunctionConnectorCont, "name"),
			Annotation:    G(FunctionConnectorCont, "annotation"),
			AttNotify:     G(FunctionConnectorCont, "attNotify"),
			ConnType:      G(FunctionConnectorCont, "connType"),
			NameAlias:     G(FunctionConnectorCont, "nameAlias"),
			DeviceLIfName: G(FunctionConnectorCont, "deviceLIfName"),
		},
	}
}

func FunctionConnectorFromContainer(cont *container.Container) *FunctionConnector {

	return FunctionConnectorFromContainerList(cont, 0)
}

func FunctionConnectorListFromContainer(cont *container.Container) []*FunctionConnector {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FunctionConnector, length)

	for i := 0; i < length; i++ {

		arr[i] = FunctionConnectorFromContainerList(cont, i)
	}

	return arr
}

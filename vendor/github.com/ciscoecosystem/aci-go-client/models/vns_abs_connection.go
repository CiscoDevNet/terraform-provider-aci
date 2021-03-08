package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VnsabsconnectionClassName = "vnsAbsConnection"

type Connection struct {
	BaseAttributes
	ConnectionAttributes
}

type ConnectionAttributes struct {
	Name string `json:",omitempty"`

	AdjType string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ConnDir string `json:",omitempty"`

	ConnType string `json:",omitempty"`

	DirectConnect string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	UnicastRoute string `json:",omitempty"`
}

func NewConnection(vnsAbsConnectionRn, parentDn, description string, vnsAbsConnectionattr ConnectionAttributes) *Connection {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsAbsConnectionRn)
	return &Connection{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsabsconnectionClassName,
			Rn:                vnsAbsConnectionRn,
		},

		ConnectionAttributes: vnsAbsConnectionattr,
	}
}

func (vnsAbsConnection *Connection) ToMap() (map[string]string, error) {
	vnsAbsConnectionMap, err := vnsAbsConnection.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vnsAbsConnectionMap, "name", vnsAbsConnection.Name)

	A(vnsAbsConnectionMap, "adjType", vnsAbsConnection.AdjType)

	A(vnsAbsConnectionMap, "annotation", vnsAbsConnection.Annotation)

	A(vnsAbsConnectionMap, "connDir", vnsAbsConnection.ConnDir)

	A(vnsAbsConnectionMap, "connType", vnsAbsConnection.ConnType)

	A(vnsAbsConnectionMap, "directConnect", vnsAbsConnection.DirectConnect)

	A(vnsAbsConnectionMap, "nameAlias", vnsAbsConnection.NameAlias)

	A(vnsAbsConnectionMap, "unicastRoute", vnsAbsConnection.UnicastRoute)

	return vnsAbsConnectionMap, err
}

func ConnectionFromContainerList(cont *container.Container, index int) *Connection {

	ConnectionCont := cont.S("imdata").Index(index).S(VnsabsconnectionClassName, "attributes")
	return &Connection{
		BaseAttributes{
			DistinguishedName: G(ConnectionCont, "dn"),
			Description:       G(ConnectionCont, "descr"),
			Status:            G(ConnectionCont, "status"),
			ClassName:         VnsabsconnectionClassName,
			Rn:                G(ConnectionCont, "rn"),
		},

		ConnectionAttributes{

			Name: G(ConnectionCont, "name"),

			AdjType: G(ConnectionCont, "adjType"),

			Annotation: G(ConnectionCont, "annotation"),

			ConnDir: G(ConnectionCont, "connDir"),

			ConnType: G(ConnectionCont, "connType"),

			DirectConnect: G(ConnectionCont, "directConnect"),

			NameAlias: G(ConnectionCont, "nameAlias"),

			UnicastRoute: G(ConnectionCont, "unicastRoute"),
		},
	}
}

func ConnectionFromContainer(cont *container.Container) *Connection {

	return ConnectionFromContainerList(cont, 0)
}

func ConnectionListFromContainer(cont *container.Container) []*Connection {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*Connection, length)

	for i := 0; i < length; i++ {

		arr[i] = ConnectionFromContainerList(cont, i)
	}

	return arr
}

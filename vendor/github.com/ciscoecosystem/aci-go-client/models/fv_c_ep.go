package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvcepClassName = "fvCEp"

type ClientEndPoint struct {
	BaseAttributes
	ClientEndPointAttributes
}

type ClientEndPointAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ClientEndPoint_id string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewClientEndPoint(fvCEpRn, parentDn, description string, fvCEpattr ClientEndPointAttributes) *ClientEndPoint {
	dn := fmt.Sprintf("%s/%s", parentDn, fvCEpRn)
	return &ClientEndPoint{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvcepClassName,
			Rn:                fvCEpRn,
		},

		ClientEndPointAttributes: fvCEpattr,
	}
}

func (fvCEp *ClientEndPoint) ToMap() (map[string]string, error) {
	fvCEpMap, err := fvCEp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvCEpMap, "name", fvCEp.Name)

	A(fvCEpMap, "annotation", fvCEp.Annotation)

	A(fvCEpMap, "id", fvCEp.ClientEndPoint_id)

	A(fvCEpMap, "nameAlias", fvCEp.NameAlias)

	return fvCEpMap, err
}

func ClientEndPointFromContainerList(cont *container.Container, index int) *ClientEndPoint {

	ClientEndPointCont := cont.S("imdata").Index(index).S(FvcepClassName, "attributes")
	return &ClientEndPoint{
		BaseAttributes{
			DistinguishedName: G(ClientEndPointCont, "dn"),
			Description:       G(ClientEndPointCont, "descr"),
			Status:            G(ClientEndPointCont, "status"),
			ClassName:         FvcepClassName,
			Rn:                G(ClientEndPointCont, "rn"),
		},

		ClientEndPointAttributes{

			Name: G(ClientEndPointCont, "name"),

			Annotation: G(ClientEndPointCont, "annotation"),

			ClientEndPoint_id: G(ClientEndPointCont, "id"),

			NameAlias: G(ClientEndPointCont, "nameAlias"),
		},
	}
}

func ClientEndPointFromContainer(cont *container.Container) *ClientEndPoint {

	return ClientEndPointFromContainerList(cont, 0)
}

func ClientEndPointListFromContainer(cont *container.Container) []*ClientEndPoint {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ClientEndPoint, length)

	for i := 0; i < length; i++ {

		arr[i] = ClientEndPointFromContainerList(cont, i)
	}

	return arr
}

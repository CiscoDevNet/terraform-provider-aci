package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FabricpathepClassName = "fabricPathEp"

type FabricPathEndpoint struct {
	BaseAttributes
	FabricPathEndpointAttributes
}

type FabricPathEndpointAttributes struct {
	Name string `json:",omitempty"`
}

func NewFabricPathEndpoint(fabricPathEpRn, parentDn, description string, fabricPathEpattr FabricPathEndpointAttributes) *FabricPathEndpoint {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricPathEpRn)
	return &FabricPathEndpoint{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FabricpathepClassName,
			Rn:                fabricPathEpRn,
		},

		FabricPathEndpointAttributes: fabricPathEpattr,
	}
}

func (fabricPathEp *FabricPathEndpoint) ToMap() (map[string]string, error) {
	fabricPathEpMap, err := fabricPathEp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fabricPathEpMap, "name", fabricPathEp.Name)

	return fabricPathEpMap, err
}

func FabricPathEndpointFromContainerList(cont *container.Container, index int) *FabricPathEndpoint {

	FabricPathEndpointCont := cont.S("imdata").Index(index).S(FabricpathepClassName, "attributes")
	return &FabricPathEndpoint{
		BaseAttributes{
			DistinguishedName: G(FabricPathEndpointCont, "dn"),
			Description:       G(FabricPathEndpointCont, "descr"),
			Status:            G(FabricPathEndpointCont, "status"),
			ClassName:         FabricpathepClassName,
			Rn:                G(FabricPathEndpointCont, "rn"),
		},

		FabricPathEndpointAttributes{
			Name: G(FabricPathEndpointCont, "name"),
		},
	}
}

func FabricPathEndpointFromContainer(cont *container.Container) *FabricPathEndpoint {

	return FabricPathEndpointFromContainerList(cont, 0)
}

func FabricPathEndpointListFromContainer(cont *container.Container) []*FabricPathEndpoint {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FabricPathEndpoint, length)

	for i := 0; i < length; i++ {

		arr[i] = FabricPathEndpointFromContainerList(cont, i)
	}

	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudrouterpClassName = "cloudRouterP"

type CloudVpnGateway struct {
	BaseAttributes
	CloudVpnGatewayAttributes
}

type CloudVpnGatewayAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	NumInstances string `json:",omitempty"`

	CloudVpnGateway_type string `json:",omitempty"`
}

func NewCloudVpnGateway(cloudRouterPRn, parentDn, description string, cloudRouterPattr CloudVpnGatewayAttributes) *CloudVpnGateway {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudRouterPRn)
	return &CloudVpnGateway{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudrouterpClassName,
			Rn:                cloudRouterPRn,
		},

		CloudVpnGatewayAttributes: cloudRouterPattr,
	}
}

func (cloudRouterP *CloudVpnGateway) ToMap() (map[string]string, error) {
	cloudRouterPMap, err := cloudRouterP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudRouterPMap, "name", cloudRouterP.Name)

	A(cloudRouterPMap, "annotation", cloudRouterP.Annotation)

	A(cloudRouterPMap, "nameAlias", cloudRouterP.NameAlias)

	A(cloudRouterPMap, "numInstances", cloudRouterP.NumInstances)

	A(cloudRouterPMap, "type", cloudRouterP.CloudVpnGateway_type)

	return cloudRouterPMap, err
}

func CloudVpnGatewayFromContainerList(cont *container.Container, index int) *CloudVpnGateway {

	CloudVpnGatewayCont := cont.S("imdata").Index(index).S(CloudrouterpClassName, "attributes")
	return &CloudVpnGateway{
		BaseAttributes{
			DistinguishedName: G(CloudVpnGatewayCont, "dn"),
			Description:       G(CloudVpnGatewayCont, "descr"),
			Status:            G(CloudVpnGatewayCont, "status"),
			ClassName:         CloudrouterpClassName,
			Rn:                G(CloudVpnGatewayCont, "rn"),
		},

		CloudVpnGatewayAttributes{

			Name: G(CloudVpnGatewayCont, "name"),

			Annotation: G(CloudVpnGatewayCont, "annotation"),

			NameAlias: G(CloudVpnGatewayCont, "nameAlias"),

			NumInstances: G(CloudVpnGatewayCont, "numInstances"),

			CloudVpnGateway_type: G(CloudVpnGatewayCont, "type"),
		},
	}
}

func CloudVpnGatewayFromContainer(cont *container.Container) *CloudVpnGateway {

	return CloudVpnGatewayFromContainerList(cont, 0)
}

func CloudVpnGatewayListFromContainer(cont *container.Container) []*CloudVpnGateway {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudVpnGateway, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudVpnGatewayFromContainerList(cont, i)
	}

	return arr
}

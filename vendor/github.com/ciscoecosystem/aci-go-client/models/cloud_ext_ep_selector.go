package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudextepselectorClassName = "cloudExtEPSelector"

type CloudEndpointSelectorforExternalEPgs struct {
	BaseAttributes
	CloudEndpointSelectorforExternalEPgsAttributes
}

type CloudEndpointSelectorforExternalEPgsAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	IsShared string `json:",omitempty"`

	MatchExpression string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Subnet string `json:",omitempty"`
}

func NewCloudEndpointSelectorforExternalEPgs(cloudExtEPSelectorRn, parentDn, description string, cloudExtEPSelectorattr CloudEndpointSelectorforExternalEPgsAttributes) *CloudEndpointSelectorforExternalEPgs {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudExtEPSelectorRn)
	return &CloudEndpointSelectorforExternalEPgs{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudextepselectorClassName,
			Rn:                cloudExtEPSelectorRn,
		},

		CloudEndpointSelectorforExternalEPgsAttributes: cloudExtEPSelectorattr,
	}
}

func (cloudExtEPSelector *CloudEndpointSelectorforExternalEPgs) ToMap() (map[string]string, error) {
	cloudExtEPSelectorMap, err := cloudExtEPSelector.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudExtEPSelectorMap, "name", cloudExtEPSelector.Name)

	A(cloudExtEPSelectorMap, "annotation", cloudExtEPSelector.Annotation)

	A(cloudExtEPSelectorMap, "isShared", cloudExtEPSelector.IsShared)

	A(cloudExtEPSelectorMap, "matchExpression", cloudExtEPSelector.MatchExpression)

	A(cloudExtEPSelectorMap, "nameAlias", cloudExtEPSelector.NameAlias)

	A(cloudExtEPSelectorMap, "subnet", cloudExtEPSelector.Subnet)

	return cloudExtEPSelectorMap, err
}

func CloudEndpointSelectorforExternalEPgsFromContainerList(cont *container.Container, index int) *CloudEndpointSelectorforExternalEPgs {

	CloudEndpointSelectorforExternalEPgsCont := cont.S("imdata").Index(index).S(CloudextepselectorClassName, "attributes")
	return &CloudEndpointSelectorforExternalEPgs{
		BaseAttributes{
			DistinguishedName: G(CloudEndpointSelectorforExternalEPgsCont, "dn"),
			Description:       G(CloudEndpointSelectorforExternalEPgsCont, "descr"),
			Status:            G(CloudEndpointSelectorforExternalEPgsCont, "status"),
			ClassName:         CloudextepselectorClassName,
			Rn:                G(CloudEndpointSelectorforExternalEPgsCont, "rn"),
		},

		CloudEndpointSelectorforExternalEPgsAttributes{

			Name: G(CloudEndpointSelectorforExternalEPgsCont, "name"),

			Annotation: G(CloudEndpointSelectorforExternalEPgsCont, "annotation"),

			IsShared: G(CloudEndpointSelectorforExternalEPgsCont, "isShared"),

			MatchExpression: G(CloudEndpointSelectorforExternalEPgsCont, "matchExpression"),

			NameAlias: G(CloudEndpointSelectorforExternalEPgsCont, "nameAlias"),

			Subnet: G(CloudEndpointSelectorforExternalEPgsCont, "subnet"),
		},
	}
}

func CloudEndpointSelectorforExternalEPgsFromContainer(cont *container.Container) *CloudEndpointSelectorforExternalEPgs {

	return CloudEndpointSelectorforExternalEPgsFromContainerList(cont, 0)
}

func CloudEndpointSelectorforExternalEPgsListFromContainer(cont *container.Container) []*CloudEndpointSelectorforExternalEPgs {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudEndpointSelectorforExternalEPgs, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudEndpointSelectorforExternalEPgsFromContainerList(cont, i)
	}

	return arr
}

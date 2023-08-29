package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnCloudSvcEPSelector        = "svcepselector-%s"
	DnCloudSvcEPSelector        = "uni/tn-%s/cloudapp-%s/cloudsvcepg-%s/svcepselector-%s"
	ParentDnCloudSvcEPSelector  = "uni/tn-%s/cloudapp-%s/cloudsvcepg-%s"
	CloudSvcEPSelectorClassName = "cloudSvcEPSelector"
)

type CloudServiceEndpointSelector struct {
	BaseAttributes
	CloudServiceEndpointSelectorAttributes
}

type CloudServiceEndpointSelectorAttributes struct {
	Annotation      string `json:",omitempty"`
	MatchExpression string `json:",omitempty"`
	Name            string `json:",omitempty"`
	NameAlias       string `json:",omitempty"`
}

func NewCloudServiceEndpointSelector(cloudSvcEPSelectorRn, parentDn, description string, cloudSvcEPSelectorAttr CloudServiceEndpointSelectorAttributes) *CloudServiceEndpointSelector {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudSvcEPSelectorRn)
	return &CloudServiceEndpointSelector{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudSvcEPSelectorClassName,
			Rn:                cloudSvcEPSelectorRn,
		},
		CloudServiceEndpointSelectorAttributes: cloudSvcEPSelectorAttr,
	}
}

func (cloudSvcEPSelector *CloudServiceEndpointSelector) ToMap() (map[string]string, error) {
	cloudSvcEPSelectorMap, err := cloudSvcEPSelector.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudSvcEPSelectorMap, "annotation", cloudSvcEPSelector.Annotation)
	A(cloudSvcEPSelectorMap, "matchExpression", cloudSvcEPSelector.MatchExpression)
	A(cloudSvcEPSelectorMap, "name", cloudSvcEPSelector.Name)
	A(cloudSvcEPSelectorMap, "nameAlias", cloudSvcEPSelector.NameAlias)
	return cloudSvcEPSelectorMap, err
}

func CloudServiceEndpointSelectorFromContainerList(cont *container.Container, index int) *CloudServiceEndpointSelector {
	CloudServiceEndpointSelectorCont := cont.S("imdata").Index(index).S(CloudSvcEPSelectorClassName, "attributes")
	return &CloudServiceEndpointSelector{
		BaseAttributes{
			DistinguishedName: G(CloudServiceEndpointSelectorCont, "dn"),
			Description:       G(CloudServiceEndpointSelectorCont, "descr"),
			Status:            G(CloudServiceEndpointSelectorCont, "status"),
			ClassName:         CloudSvcEPSelectorClassName,
			Rn:                G(CloudServiceEndpointSelectorCont, "rn"),
		},
		CloudServiceEndpointSelectorAttributes{
			Annotation:      G(CloudServiceEndpointSelectorCont, "annotation"),
			MatchExpression: G(CloudServiceEndpointSelectorCont, "matchExpression"),
			Name:            G(CloudServiceEndpointSelectorCont, "name"),
			NameAlias:       G(CloudServiceEndpointSelectorCont, "nameAlias"),
		},
	}
}

func CloudServiceEndpointSelectorFromContainer(cont *container.Container) *CloudServiceEndpointSelector {
	return CloudServiceEndpointSelectorFromContainerList(cont, 0)
}

func CloudServiceEndpointSelectorListFromContainer(cont *container.Container) []*CloudServiceEndpointSelector {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudServiceEndpointSelector, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudServiceEndpointSelectorFromContainerList(cont, i)
	}

	return arr
}

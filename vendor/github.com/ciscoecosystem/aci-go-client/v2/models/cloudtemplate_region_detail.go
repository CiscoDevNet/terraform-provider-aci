package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnCloudtemplateRegionDetail        = "regiondetail"
	CloudtemplateRegionDetailClassName = "cloudtemplateRegionDetail"
)

type CloudTemplateRegion struct {
	BaseAttributes
	CloudTemplateRegionAttributes
}

type CloudTemplateRegionAttributes struct {
	Annotation           string `json:",omitempty"`
	HubNetworkingEnabled string `json:",omitempty"`
}

func NewCloudTemplateRegion(cloudtemplateRegionDetailRn, parentDn string, cloudtemplateRegionDetailAttr CloudTemplateRegionAttributes) *CloudTemplateRegion {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudtemplateRegionDetailRn)
	return &CloudTemplateRegion{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "modified",
			ClassName:         CloudtemplateRegionDetailClassName,
			Rn:                cloudtemplateRegionDetailRn,
		},
		CloudTemplateRegionAttributes: cloudtemplateRegionDetailAttr,
	}
}

func (cloudtemplateRegionDetail *CloudTemplateRegion) ToMap() (map[string]string, error) {
	cloudtemplateRegionDetailMap, err := cloudtemplateRegionDetail.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudtemplateRegionDetailMap, "annotation", cloudtemplateRegionDetail.Annotation)
	A(cloudtemplateRegionDetailMap, "hubNetworkingEnabled", cloudtemplateRegionDetail.HubNetworkingEnabled)
	return cloudtemplateRegionDetailMap, err
}

func CloudTemplateRegionFromContainerList(cont *container.Container, index int) *CloudTemplateRegion {
	CloudTemplateRegionCont := cont.S("imdata").Index(index).S(CloudtemplateRegionDetailClassName, "attributes")
	return &CloudTemplateRegion{
		BaseAttributes{
			DistinguishedName: G(CloudTemplateRegionCont, "dn"),
			Status:            G(CloudTemplateRegionCont, "status"),
			ClassName:         CloudtemplateRegionDetailClassName,
			Rn:                G(CloudTemplateRegionCont, "rn"),
		},
		CloudTemplateRegionAttributes{
			Annotation:           G(CloudTemplateRegionCont, "annotation"),
			HubNetworkingEnabled: G(CloudTemplateRegionCont, "hubNetworkingEnabled"),
		},
	}
}

func CloudTemplateRegionFromContainer(cont *container.Container) *CloudTemplateRegion {
	return CloudTemplateRegionFromContainerList(cont, 0)
}

func CloudTemplateRegionListFromContainer(cont *container.Container) []*CloudTemplateRegion {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudTemplateRegion, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudTemplateRegionFromContainerList(cont, i)
	}

	return arr
}

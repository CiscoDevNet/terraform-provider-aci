package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DncloudRegionName        = "uni/tn-%s/infranetwork-%s/extnetwork-%s/provider-%s-region-%s"
	RncloudRegionName        = "provider-%s-region-%s"
	ParentDncloudRegionName  = "uni/tn-%s/infranetwork-%s/extnetwork-%s"
	CloudregionnameClassName = "cloudRegionName"
)

type CloudProviderandRegionNames struct {
	BaseAttributes
	CloudProviderandRegionNamesAttributes
}

type CloudProviderandRegionNamesAttributes struct {
	Annotation string `json:",omitempty"`
	Provider   string `json:",omitempty"`
	Region     string `json:",omitempty"`
}

func NewCloudProviderandRegionNames(cloudRegionNameRn, parentDn string, cloudRegionNameAttr CloudProviderandRegionNamesAttributes) *CloudProviderandRegionNames {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudRegionNameRn)
	return &CloudProviderandRegionNames{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         CloudregionnameClassName,
			Rn:                cloudRegionNameRn,
		},
		CloudProviderandRegionNamesAttributes: cloudRegionNameAttr,
	}
}

func (cloudRegionName *CloudProviderandRegionNames) ToMap() (map[string]string, error) {
	cloudRegionNameMap, err := cloudRegionName.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudRegionNameMap, "provider", cloudRegionName.Provider)
	A(cloudRegionNameMap, "region", cloudRegionName.Region)
	return cloudRegionNameMap, err
}

func CloudProviderandRegionNamesFromContainerList(cont *container.Container, index int) *CloudProviderandRegionNames {
	CloudProviderandRegionNamesCont := cont.S("imdata").Index(index).S(CloudregionnameClassName, "attributes")
	return &CloudProviderandRegionNames{
		BaseAttributes{
			DistinguishedName: G(CloudProviderandRegionNamesCont, "dn"),
			Status:            G(CloudProviderandRegionNamesCont, "status"),
			ClassName:         CloudregionnameClassName,
			Rn:                G(CloudProviderandRegionNamesCont, "rn"),
		},
		CloudProviderandRegionNamesAttributes{
			Provider: G(CloudProviderandRegionNamesCont, "provider"),
			Region:   G(CloudProviderandRegionNamesCont, "region"),
		},
	}
}

func CloudProviderandRegionNamesFromContainer(cont *container.Container) *CloudProviderandRegionNames {
	return CloudProviderandRegionNamesFromContainerList(cont, 0)
}

func CloudProviderandRegionNamesListFromContainer(cont *container.Container) []*CloudProviderandRegionNames {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudProviderandRegionNames, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudProviderandRegionNamesFromContainerList(cont, i)
	}

	return arr
}

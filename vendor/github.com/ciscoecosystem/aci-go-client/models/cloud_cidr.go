package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudcidrClassName = "cloudCidr"

type CloudCIDRPool struct {
	BaseAttributes
	CloudCIDRPoolAttributes
}

type CloudCIDRPoolAttributes struct {
	Addr string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Primary string `json:",omitempty"`
}

func NewCloudCIDRPool(cloudCidrRn, parentDn, description string, cloudCidrattr CloudCIDRPoolAttributes) *CloudCIDRPool {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudCidrRn)
	return &CloudCIDRPool{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudcidrClassName,
			Rn:                cloudCidrRn,
		},

		CloudCIDRPoolAttributes: cloudCidrattr,
	}
}

func (cloudCidr *CloudCIDRPool) ToMap() (map[string]string, error) {
	cloudCidrMap, err := cloudCidr.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudCidrMap, "addr", cloudCidr.Addr)

	A(cloudCidrMap, "annotation", cloudCidr.Annotation)

	A(cloudCidrMap, "nameAlias", cloudCidr.NameAlias)

	A(cloudCidrMap, "primary", cloudCidr.Primary)

	return cloudCidrMap, err
}

func CloudCIDRPoolFromContainerList(cont *container.Container, index int) *CloudCIDRPool {

	CloudCIDRPoolCont := cont.S("imdata").Index(index).S(CloudcidrClassName, "attributes")
	return &CloudCIDRPool{
		BaseAttributes{
			DistinguishedName: G(CloudCIDRPoolCont, "dn"),
			Description:       G(CloudCIDRPoolCont, "descr"),
			Status:            G(CloudCIDRPoolCont, "status"),
			ClassName:         CloudcidrClassName,
			Rn:                G(CloudCIDRPoolCont, "rn"),
		},

		CloudCIDRPoolAttributes{

			Addr: G(CloudCIDRPoolCont, "addr"),

			Annotation: G(CloudCIDRPoolCont, "annotation"),

			NameAlias: G(CloudCIDRPoolCont, "nameAlias"),

			Primary: G(CloudCIDRPoolCont, "primary"),
		},
	}
}

func CloudCIDRPoolFromContainer(cont *container.Container) *CloudCIDRPool {

	return CloudCIDRPoolFromContainerList(cont, 0)
}

func CloudCIDRPoolListFromContainer(cont *container.Container) []*CloudCIDRPool {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudCIDRPool, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudCIDRPoolFromContainerList(cont, i)
	}

	return arr
}

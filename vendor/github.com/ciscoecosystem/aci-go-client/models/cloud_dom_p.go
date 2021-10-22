package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const ClouddompClassName = "cloudDomP"

type CloudDomainProfile struct {
	BaseAttributes
	CloudDomainProfileAttributes
}

type CloudDomainProfileAttributes struct {
	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	SiteId string `json:",omitempty"`
}

func NewCloudDomainProfile(cloudDomPRn, parentDn, description string, cloudDomPattr CloudDomainProfileAttributes) *CloudDomainProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudDomPRn)
	return &CloudDomainProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         ClouddompClassName,
			Rn:                cloudDomPRn,
		},

		CloudDomainProfileAttributes: cloudDomPattr,
	}
}

func (cloudDomP *CloudDomainProfile) ToMap() (map[string]string, error) {
	cloudDomPMap, err := cloudDomP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudDomPMap, "annotation", cloudDomP.Annotation)

	A(cloudDomPMap, "nameAlias", cloudDomP.NameAlias)

	A(cloudDomPMap, "siteId", cloudDomP.SiteId)

	return cloudDomPMap, err
}

func CloudDomainProfileFromContainerList(cont *container.Container, index int) *CloudDomainProfile {

	CloudDomainProfileCont := cont.S("imdata").Index(index).S(ClouddompClassName, "attributes")
	return &CloudDomainProfile{
		BaseAttributes{
			DistinguishedName: G(CloudDomainProfileCont, "dn"),
			Description:       G(CloudDomainProfileCont, "descr"),
			Status:            G(CloudDomainProfileCont, "status"),
			ClassName:         ClouddompClassName,
			Rn:                G(CloudDomainProfileCont, "rn"),
		},

		CloudDomainProfileAttributes{

			Annotation: G(CloudDomainProfileCont, "annotation"),

			NameAlias: G(CloudDomainProfileCont, "nameAlias"),

			SiteId: G(CloudDomainProfileCont, "siteId"),
		},
	}
}

func CloudDomainProfileFromContainer(cont *container.Container) *CloudDomainProfile {

	return CloudDomainProfileFromContainerList(cont, 0)
}

func CloudDomainProfileListFromContainer(cont *container.Container) []*CloudDomainProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudDomainProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudDomainProfileFromContainerList(cont, i)
	}

	return arr
}

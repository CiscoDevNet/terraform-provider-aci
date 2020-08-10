package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudprovpClassName = "cloudProvP"

type CloudProviderProfile struct {
	BaseAttributes
	CloudProviderProfileAttributes
}

type CloudProviderProfileAttributes struct {
	Vendor string `json:",omitempty"`

	Annotation string `json:",omitempty"`
}

func NewCloudProviderProfile(cloudProvPRn, parentDn, description string, cloudProvPattr CloudProviderProfileAttributes) *CloudProviderProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudProvPRn)
	return &CloudProviderProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudprovpClassName,
			// Rn:                cloudProvPRn,
		},

		CloudProviderProfileAttributes: cloudProvPattr,
	}
}

func (cloudProvP *CloudProviderProfile) ToMap() (map[string]string, error) {
	cloudProvPMap, err := cloudProvP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudProvPMap, "vendor", cloudProvP.Vendor)

	A(cloudProvPMap, "annotation", cloudProvP.Annotation)

	return cloudProvPMap, err
}

func CloudProviderProfileFromContainerList(cont *container.Container, index int) *CloudProviderProfile {

	CloudProviderProfileCont := cont.S("imdata").Index(index).S(CloudprovpClassName, "attributes")
	return &CloudProviderProfile{
		BaseAttributes{
			DistinguishedName: G(CloudProviderProfileCont, "dn"),
			Description:       G(CloudProviderProfileCont, "descr"),
			Status:            G(CloudProviderProfileCont, "status"),
			ClassName:         CloudprovpClassName,
			Rn:                G(CloudProviderProfileCont, "rn"),
		},

		CloudProviderProfileAttributes{

			Vendor: G(CloudProviderProfileCont, "vendor"),

			Annotation: G(CloudProviderProfileCont, "annotation"),
		},
	}
}

func CloudProviderProfileFromContainer(cont *container.Container) *CloudProviderProfile {

	return CloudProviderProfileFromContainerList(cont, 0)
}

func CloudProviderProfileListFromContainer(cont *container.Container) []*CloudProviderProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudProviderProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudProviderProfileFromContainerList(cont, i)
	}

	return arr
}

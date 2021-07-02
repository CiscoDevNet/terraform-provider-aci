package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudregionClassName = "cloudRegion"

type CloudProvidersRegion struct {
	BaseAttributes
	CloudProvidersRegionAttributes
}

type CloudProvidersRegionAttributes struct {
	Name string `json:",omitempty"`

	AdminSt string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewCloudProvidersRegion(cloudRegionRn, parentDn, description string, cloudRegionattr CloudProvidersRegionAttributes) *CloudProvidersRegion {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudRegionRn)
	return &CloudProvidersRegion{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudregionClassName,
			Rn:                cloudRegionRn,
		},

		CloudProvidersRegionAttributes: cloudRegionattr,
	}
}

func (cloudRegion *CloudProvidersRegion) ToMap() (map[string]string, error) {
	cloudRegionMap, err := cloudRegion.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudRegionMap, "name", cloudRegion.Name)

	A(cloudRegionMap, "adminSt", cloudRegion.AdminSt)

	A(cloudRegionMap, "annotation", cloudRegion.Annotation)

	A(cloudRegionMap, "nameAlias", cloudRegion.NameAlias)

	return cloudRegionMap, err
}

func CloudProvidersRegionFromContainerList(cont *container.Container, index int) *CloudProvidersRegion {

	CloudProvidersRegionCont := cont.S("imdata").Index(index).S(CloudregionClassName, "attributes")
	return &CloudProvidersRegion{
		BaseAttributes{
			DistinguishedName: G(CloudProvidersRegionCont, "dn"),
			Description:       G(CloudProvidersRegionCont, "descr"),
			Status:            G(CloudProvidersRegionCont, "status"),
			ClassName:         CloudregionClassName,
			Rn:                G(CloudProvidersRegionCont, "rn"),
		},

		CloudProvidersRegionAttributes{

			Name: G(CloudProvidersRegionCont, "name"),

			AdminSt: G(CloudProvidersRegionCont, "adminSt"),

			Annotation: G(CloudProvidersRegionCont, "annotation"),

			NameAlias: G(CloudProvidersRegionCont, "nameAlias"),
		},
	}
}

func CloudProvidersRegionFromContainer(cont *container.Container) *CloudProvidersRegion {

	return CloudProvidersRegionFromContainerList(cont, 0)
}

func CloudProvidersRegionListFromContainer(cont *container.Container) []*CloudProvidersRegion {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudProvidersRegion, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudProvidersRegionFromContainerList(cont, i)
	}

	return arr
}

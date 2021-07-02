package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudzoneClassName = "cloudZone"

type CloudAvailabilityZone struct {
	BaseAttributes
	CloudAvailabilityZoneAttributes
}

type CloudAvailabilityZoneAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewCloudAvailabilityZone(cloudZoneRn, parentDn, description string, cloudZoneattr CloudAvailabilityZoneAttributes) *CloudAvailabilityZone {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudZoneRn)
	return &CloudAvailabilityZone{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudzoneClassName,
			Rn:                cloudZoneRn,
		},

		CloudAvailabilityZoneAttributes: cloudZoneattr,
	}
}

func (cloudZone *CloudAvailabilityZone) ToMap() (map[string]string, error) {
	cloudZoneMap, err := cloudZone.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudZoneMap, "name", cloudZone.Name)

	A(cloudZoneMap, "annotation", cloudZone.Annotation)

	A(cloudZoneMap, "nameAlias", cloudZone.NameAlias)

	return cloudZoneMap, err
}

func CloudAvailabilityZoneFromContainerList(cont *container.Container, index int) *CloudAvailabilityZone {

	CloudAvailabilityZoneCont := cont.S("imdata").Index(index).S(CloudzoneClassName, "attributes")
	return &CloudAvailabilityZone{
		BaseAttributes{
			DistinguishedName: G(CloudAvailabilityZoneCont, "dn"),
			Description:       G(CloudAvailabilityZoneCont, "descr"),
			Status:            G(CloudAvailabilityZoneCont, "status"),
			ClassName:         CloudzoneClassName,
			Rn:                G(CloudAvailabilityZoneCont, "rn"),
		},

		CloudAvailabilityZoneAttributes{

			Name: G(CloudAvailabilityZoneCont, "name"),

			Annotation: G(CloudAvailabilityZoneCont, "annotation"),

			NameAlias: G(CloudAvailabilityZoneCont, "nameAlias"),
		},
	}
}

func CloudAvailabilityZoneFromContainer(cont *container.Container) *CloudAvailabilityZone {

	return CloudAvailabilityZoneFromContainerList(cont, 0)
}

func CloudAvailabilityZoneListFromContainer(cont *container.Container) []*CloudAvailabilityZone {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudAvailabilityZone, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudAvailabilityZoneFromContainerList(cont, i)
	}

	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudappClassName = "cloudApp"

type CloudApplicationcontainer struct {
	BaseAttributes
	CloudApplicationcontainerAttributes
}

type CloudApplicationcontainerAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewCloudApplicationcontainer(cloudAppRn, parentDn, description string, cloudAppattr CloudApplicationcontainerAttributes) *CloudApplicationcontainer {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudAppRn)
	return &CloudApplicationcontainer{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudappClassName,
			Rn:                cloudAppRn,
		},

		CloudApplicationcontainerAttributes: cloudAppattr,
	}
}

func (cloudApp *CloudApplicationcontainer) ToMap() (map[string]string, error) {
	cloudAppMap, err := cloudApp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudAppMap, "name", cloudApp.Name)

	A(cloudAppMap, "annotation", cloudApp.Annotation)

	A(cloudAppMap, "nameAlias", cloudApp.NameAlias)

	return cloudAppMap, err
}

func CloudApplicationcontainerFromContainerList(cont *container.Container, index int) *CloudApplicationcontainer {

	CloudApplicationcontainerCont := cont.S("imdata").Index(index).S(CloudappClassName, "attributes")
	return &CloudApplicationcontainer{
		BaseAttributes{
			DistinguishedName: G(CloudApplicationcontainerCont, "dn"),
			Description:       G(CloudApplicationcontainerCont, "descr"),
			Status:            G(CloudApplicationcontainerCont, "status"),
			ClassName:         CloudappClassName,
			Rn:                G(CloudApplicationcontainerCont, "rn"),
		},

		CloudApplicationcontainerAttributes{

			Name: G(CloudApplicationcontainerCont, "name"),

			Annotation: G(CloudApplicationcontainerCont, "annotation"),

			NameAlias: G(CloudApplicationcontainerCont, "nameAlias"),
		},
	}
}

func CloudApplicationcontainerFromContainer(cont *container.Container) *CloudApplicationcontainer {

	return CloudApplicationcontainerFromContainerList(cont, 0)
}

func CloudApplicationcontainerListFromContainer(cont *container.Container) []*CloudApplicationcontainer {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudApplicationcontainer, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudApplicationcontainerFromContainerList(cont, i)
	}

	return arr
}

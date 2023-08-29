package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnCloudPrivateLinkLabel        = "privatelinklabel-%s"
	CloudPrivateLinkLabelClassName = "cloudPrivateLinkLabel"
)

type CloudPrivateLinkLabel struct {
	BaseAttributes
	CloudPrivateLinkLabelAttributes
}

type CloudPrivateLinkLabelAttributes struct {
	Annotation string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewCloudPrivateLinkLabel(cloudPrivateLinkLabelRn, parentDn, description string, cloudPrivateLinkLabelAttr CloudPrivateLinkLabelAttributes) *CloudPrivateLinkLabel {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudPrivateLinkLabelRn)
	return &CloudPrivateLinkLabel{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudPrivateLinkLabelClassName,
			Rn:                cloudPrivateLinkLabelRn,
		},
		CloudPrivateLinkLabelAttributes: cloudPrivateLinkLabelAttr,
	}
}

func (cloudPrivateLinkLabel *CloudPrivateLinkLabel) ToMap() (map[string]string, error) {
	cloudPrivateLinkLabelMap, err := cloudPrivateLinkLabel.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudPrivateLinkLabelMap, "annotation", cloudPrivateLinkLabel.Annotation)
	A(cloudPrivateLinkLabelMap, "name", cloudPrivateLinkLabel.Name)
	A(cloudPrivateLinkLabelMap, "nameAlias", cloudPrivateLinkLabel.NameAlias)
	return cloudPrivateLinkLabelMap, err
}

func CloudPrivateLinkLabelFromContainerList(cont *container.Container, index int) *CloudPrivateLinkLabel {
	CloudPrivateLinkLabelCont := cont.S("imdata").Index(index).S(CloudPrivateLinkLabelClassName, "attributes")
	return &CloudPrivateLinkLabel{
		BaseAttributes{
			DistinguishedName: G(CloudPrivateLinkLabelCont, "dn"),
			Description:       G(CloudPrivateLinkLabelCont, "descr"),
			Status:            G(CloudPrivateLinkLabelCont, "status"),
			ClassName:         CloudPrivateLinkLabelClassName,
			Rn:                G(CloudPrivateLinkLabelCont, "rn"),
		},
		CloudPrivateLinkLabelAttributes{
			Annotation: G(CloudPrivateLinkLabelCont, "annotation"),
			Name:       G(CloudPrivateLinkLabelCont, "name"),
			NameAlias:  G(CloudPrivateLinkLabelCont, "nameAlias"),
		},
	}
}

func CloudPrivateLinkLabelFromContainer(cont *container.Container) *CloudPrivateLinkLabel {
	return CloudPrivateLinkLabelFromContainerList(cont, 0)
}

func CloudPrivateLinkLabelListFromContainer(cont *container.Container) []*CloudPrivateLinkLabel {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudPrivateLinkLabel, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudPrivateLinkLabelFromContainerList(cont, i)
	}

	return arr
}

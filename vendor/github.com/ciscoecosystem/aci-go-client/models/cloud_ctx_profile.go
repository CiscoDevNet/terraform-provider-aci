package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DncloudCtxProfile        = "uni/tn-%s/ctxprofile-%s"
	RncloudCtxProfile        = "ctxprofile-%s"
	ParentDncloudCtxProfile  = "uni/tn-%s"
	CloudctxprofileClassName = "cloudCtxProfile"
)

type CloudContextProfile struct {
	BaseAttributes
	CloudContextProfileAttributes
}

type CloudContextProfileAttributes struct {
	Annotation string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
	Type       string `json:",omitempty"`
	VpcGroup   string `json:",omitempty"`
}

func NewCloudContextProfile(cloudCtxProfileRn, parentDn, description string, cloudCtxProfileattr CloudContextProfileAttributes) *CloudContextProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudCtxProfileRn)
	return &CloudContextProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudctxprofileClassName,
			Rn:                cloudCtxProfileRn,
		},

		CloudContextProfileAttributes: cloudCtxProfileattr,
	}
}

func (cloudCtxProfile *CloudContextProfile) ToMap() (map[string]string, error) {
	cloudCtxProfileMap, err := cloudCtxProfile.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudCtxProfileMap, "annotation", cloudCtxProfile.Annotation)
	A(cloudCtxProfileMap, "nameAlias", cloudCtxProfile.NameAlias)
	A(cloudCtxProfileMap, "type", cloudCtxProfile.Type)
	A(cloudCtxProfileMap, "vpcGroup", cloudCtxProfile.VpcGroup)

	return cloudCtxProfileMap, err
}

func CloudContextProfileFromContainerList(cont *container.Container, index int) *CloudContextProfile {

	CloudContextProfileCont := cont.S("imdata").Index(index).S(CloudctxprofileClassName, "attributes")
	return &CloudContextProfile{
		BaseAttributes{
			DistinguishedName: G(CloudContextProfileCont, "dn"),
			Description:       G(CloudContextProfileCont, "descr"),
			Status:            G(CloudContextProfileCont, "status"),
			ClassName:         CloudctxprofileClassName,
			Rn:                G(CloudContextProfileCont, "rn"),
		},

		CloudContextProfileAttributes{
			Annotation: G(CloudContextProfileCont, "annotation"),
			NameAlias:  G(CloudContextProfileCont, "nameAlias"),
			Type:       G(CloudContextProfileCont, "type"),
			VpcGroup:   G(CloudContextProfileCont, "vpcGroup"),
		},
	}
}

func CloudContextProfileFromContainer(cont *container.Container) *CloudContextProfile {
	return CloudContextProfileFromContainerList(cont, 0)
}

func CloudContextProfileListFromContainer(cont *container.Container) []*CloudContextProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudContextProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudContextProfileFromContainerList(cont, i)
	}

	return arr
}

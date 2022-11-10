package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DncloudtemplateExtNetwork        = "uni/tn-%s/infranetwork-%s/extnetwork-%s"
	RncloudtemplateExtNetwork        = "extnetwork-%s"
	ParentDncloudtemplateExtNetwork  = "uni/tn-%s/infranetwork-%s"
	CloudtemplateextnetworkClassName = "cloudtemplateExtNetwork"
)

type CloudTemplateforExternalNetwork struct {
	BaseAttributes
	NameAliasAttribute
	CloudTemplateforExternalNetworkAttributes
}

type CloudTemplateforExternalNetworkAttributes struct {
	Annotation     string `json:",omitempty"`
	HubNetworkName string `json:",omitempty"`
	Name           string `json:",omitempty"`
	VrfName        string `json:",omitempty"`
	AllRegion      string `json:",omitempty"`
	HostRouterName string `json:",omitempty"`
}

func NewCloudTemplateforExternalNetwork(cloudtemplateExtNetworkRn, parentDn, nameAlias string, cloudtemplateExtNetworkAttr CloudTemplateforExternalNetworkAttributes) *CloudTemplateforExternalNetwork {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudtemplateExtNetworkRn)
	return &CloudTemplateforExternalNetwork{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         CloudtemplateextnetworkClassName,
			Rn:                cloudtemplateExtNetworkRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		CloudTemplateforExternalNetworkAttributes: cloudtemplateExtNetworkAttr,
	}
}

func (cloudtemplateExtNetwork *CloudTemplateforExternalNetwork) ToMap() (map[string]string, error) {
	cloudtemplateExtNetworkMap, err := cloudtemplateExtNetwork.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := cloudtemplateExtNetwork.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(cloudtemplateExtNetworkMap, key, value)
	}

	A(cloudtemplateExtNetworkMap, "annotation", cloudtemplateExtNetwork.Annotation)
	A(cloudtemplateExtNetworkMap, "hubNetworkName", cloudtemplateExtNetwork.HubNetworkName)
	A(cloudtemplateExtNetworkMap, "name", cloudtemplateExtNetwork.Name)
	A(cloudtemplateExtNetworkMap, "vrfName", cloudtemplateExtNetwork.VrfName)
	A(cloudtemplateExtNetworkMap, "allRegion", cloudtemplateExtNetwork.AllRegion)
	A(cloudtemplateExtNetworkMap, "hostRouterName", cloudtemplateExtNetwork.HostRouterName)
	return cloudtemplateExtNetworkMap, err
}

func CloudTemplateforExternalNetworkFromContainerList(cont *container.Container, index int) *CloudTemplateforExternalNetwork {
	CloudTemplateforExternalNetworkCont := cont.S("imdata").Index(index).S(CloudtemplateextnetworkClassName, "attributes")
	return &CloudTemplateforExternalNetwork{
		BaseAttributes{
			DistinguishedName: G(CloudTemplateforExternalNetworkCont, "dn"),
			Status:            G(CloudTemplateforExternalNetworkCont, "status"),
			ClassName:         CloudtemplateextnetworkClassName,
			Rn:                G(CloudTemplateforExternalNetworkCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(CloudTemplateforExternalNetworkCont, "nameAlias"),
		},
		CloudTemplateforExternalNetworkAttributes{
			Annotation:     G(CloudTemplateforExternalNetworkCont, "annotation"),
			HubNetworkName: G(CloudTemplateforExternalNetworkCont, "hubNetworkName"),
			Name:           G(CloudTemplateforExternalNetworkCont, "name"),
			VrfName:        G(CloudTemplateforExternalNetworkCont, "vrfName"),
			AllRegion:      G(CloudTemplateforExternalNetworkCont, "allRegion"),
			HostRouterName: G(CloudTemplateforExternalNetworkCont, "hostRouterName"),
		},
	}
}

func CloudTemplateforExternalNetworkFromContainer(cont *container.Container) *CloudTemplateforExternalNetwork {
	return CloudTemplateforExternalNetworkFromContainerList(cont, 0)
}

func CloudTemplateforExternalNetworkListFromContainer(cont *container.Container) []*CloudTemplateforExternalNetwork {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudTemplateforExternalNetwork, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudTemplateforExternalNetworkFromContainerList(cont, i)
	}

	return arr
}

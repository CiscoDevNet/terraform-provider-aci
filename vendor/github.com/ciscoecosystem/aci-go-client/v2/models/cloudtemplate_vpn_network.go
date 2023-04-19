package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DncloudtemplateVpnNetwork        = "uni/tn-%s/infranetwork-%s/extnetwork-%s/vpnnetwork-%s"
	RncloudtemplateVpnNetwork        = "vpnnetwork-%s"
	ParentDncloudtemplateVpnNetwork  = "uni/tn-%s/infranetwork-%s/extnetwork-%s"
	CloudtemplatevpnnetworkClassName = "cloudtemplateVpnNetwork"
)

type CloudTemplateforVPNNetwork struct {
	BaseAttributes
	NameAliasAttribute
	CloudTemplateforVPNNetworkAttributes
}

type CloudTemplateforVPNNetworkAttributes struct {
	Annotation     string `json:",omitempty"`
	Name           string `json:",omitempty"`
	RemoteSiteId   string `json:",omitempty"`
	RemoteSiteName string `json:",omitempty"`
}

func NewCloudTemplateforVPNNetwork(cloudtemplateVpnNetworkRn, parentDn, nameAlias string, cloudtemplateVpnNetworkAttr CloudTemplateforVPNNetworkAttributes) *CloudTemplateforVPNNetwork {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudtemplateVpnNetworkRn)
	return &CloudTemplateforVPNNetwork{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         CloudtemplatevpnnetworkClassName,
			Rn:                cloudtemplateVpnNetworkRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		CloudTemplateforVPNNetworkAttributes: cloudtemplateVpnNetworkAttr,
	}
}

func (cloudtemplateVpnNetwork *CloudTemplateforVPNNetwork) ToMap() (map[string]string, error) {
	cloudtemplateVpnNetworkMap, err := cloudtemplateVpnNetwork.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := cloudtemplateVpnNetwork.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(cloudtemplateVpnNetworkMap, key, value)
	}

	A(cloudtemplateVpnNetworkMap, "annotation", cloudtemplateVpnNetwork.Annotation)
	A(cloudtemplateVpnNetworkMap, "name", cloudtemplateVpnNetwork.Name)
	A(cloudtemplateVpnNetworkMap, "remoteSiteId", cloudtemplateVpnNetwork.RemoteSiteId)
	A(cloudtemplateVpnNetworkMap, "remoteSiteName", cloudtemplateVpnNetwork.RemoteSiteName)
	return cloudtemplateVpnNetworkMap, err
}

func CloudTemplateforVPNNetworkFromContainerList(cont *container.Container, index int) *CloudTemplateforVPNNetwork {
	CloudTemplateforVPNNetworkCont := cont.S("imdata").Index(index).S(CloudtemplatevpnnetworkClassName, "attributes")
	return &CloudTemplateforVPNNetwork{
		BaseAttributes{
			DistinguishedName: G(CloudTemplateforVPNNetworkCont, "dn"),
			Status:            G(CloudTemplateforVPNNetworkCont, "status"),
			ClassName:         CloudtemplatevpnnetworkClassName,
			Rn:                G(CloudTemplateforVPNNetworkCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(CloudTemplateforVPNNetworkCont, "nameAlias"),
		},
		CloudTemplateforVPNNetworkAttributes{
			Annotation:     G(CloudTemplateforVPNNetworkCont, "annotation"),
			Name:           G(CloudTemplateforVPNNetworkCont, "name"),
			RemoteSiteId:   G(CloudTemplateforVPNNetworkCont, "remoteSiteId"),
			RemoteSiteName: G(CloudTemplateforVPNNetworkCont, "remoteSiteName"),
		},
	}
}

func CloudTemplateforVPNNetworkFromContainer(cont *container.Container) *CloudTemplateforVPNNetwork {
	return CloudTemplateforVPNNetworkFromContainerList(cont, 0)
}

func CloudTemplateforVPNNetworkListFromContainer(cont *container.Container) []*CloudTemplateforVPNNetwork {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudTemplateforVPNNetwork, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudTemplateforVPNNetworkFromContainerList(cont, i)
	}

	return arr
}

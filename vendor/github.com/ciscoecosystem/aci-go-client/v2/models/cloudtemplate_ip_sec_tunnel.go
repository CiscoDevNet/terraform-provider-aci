package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DncloudtemplateIpSecTunnel        = "uni/tn-%s/infranetwork-%s/extnetwork-%s/vpnnetwork-%s/ipsec-[%s]"
	RncloudtemplateIpSecTunnel        = "ipsec-[%s]"
	ParentDncloudtemplateIpSecTunnel  = "uni/tn-%s/infranetwork-%s/extnetwork-%s/vpnnetwork-%s"
	CloudtemplateipsectunnelClassName = "cloudtemplateIpSecTunnel"
)

type CloudTemplateforIpSectunnel struct {
	BaseAttributes
	CloudTemplateforIpSectunnelAttributes
}

type CloudTemplateforIpSectunnelAttributes struct {
	Annotation   string `json:",omitempty"`
	IkeVersion   string `json:",omitempty"`
	Peeraddr     string `json:",omitempty"`
	Poolname     string `json:",omitempty"`
	PreSharedKey string `json:",omitempty"`
}

func NewCloudTemplateforIpSectunnel(cloudtemplateIpSecTunnelRn, parentDn string, cloudtemplateIpSecTunnelAttr CloudTemplateforIpSectunnelAttributes) *CloudTemplateforIpSectunnel {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudtemplateIpSecTunnelRn)
	return &CloudTemplateforIpSectunnel{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         CloudtemplateipsectunnelClassName,
			Rn:                cloudtemplateIpSecTunnelRn,
		},
		CloudTemplateforIpSectunnelAttributes: cloudtemplateIpSecTunnelAttr,
	}
}

func (cloudtemplateIpSecTunnel *CloudTemplateforIpSectunnel) ToMap() (map[string]string, error) {
	cloudtemplateIpSecTunnelMap, err := cloudtemplateIpSecTunnel.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudtemplateIpSecTunnelMap, "annotation", cloudtemplateIpSecTunnel.Annotation)
	A(cloudtemplateIpSecTunnelMap, "ikeVersion", cloudtemplateIpSecTunnel.IkeVersion)
	A(cloudtemplateIpSecTunnelMap, "peeraddr", cloudtemplateIpSecTunnel.Peeraddr)
	A(cloudtemplateIpSecTunnelMap, "poolname", cloudtemplateIpSecTunnel.Poolname)
	A(cloudtemplateIpSecTunnelMap, "preSharedKey", cloudtemplateIpSecTunnel.PreSharedKey)
	return cloudtemplateIpSecTunnelMap, err
}

func CloudTemplateforIpSectunnelFromContainerList(cont *container.Container, index int) *CloudTemplateforIpSectunnel {
	CloudTemplateforIpSectunnelCont := cont.S("imdata").Index(index).S(CloudtemplateipsectunnelClassName, "attributes")
	return &CloudTemplateforIpSectunnel{
		BaseAttributes{
			DistinguishedName: G(CloudTemplateforIpSectunnelCont, "dn"),
			Status:            G(CloudTemplateforIpSectunnelCont, "status"),
			ClassName:         CloudtemplateipsectunnelClassName,
			Rn:                G(CloudTemplateforIpSectunnelCont, "rn"),
		},
		CloudTemplateforIpSectunnelAttributes{
			Annotation:   G(CloudTemplateforIpSectunnelCont, "annotation"),
			IkeVersion:   G(CloudTemplateforIpSectunnelCont, "ikeVersion"),
			Peeraddr:     G(CloudTemplateforIpSectunnelCont, "peeraddr"),
			Poolname:     G(CloudTemplateforIpSectunnelCont, "poolname"),
			PreSharedKey: G(CloudTemplateforIpSectunnelCont, "preSharedKey"),
		},
	}
}

func CloudTemplateforIpSectunnelFromContainer(cont *container.Container) *CloudTemplateforIpSectunnel {
	return CloudTemplateforIpSectunnelFromContainerList(cont, 0)
}

func CloudTemplateforIpSectunnelListFromContainer(cont *container.Container) []*CloudTemplateforIpSectunnel {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudTemplateforIpSectunnel, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudTemplateforIpSectunnelFromContainerList(cont, i)
	}

	return arr
}

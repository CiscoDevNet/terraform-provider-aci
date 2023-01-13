package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DncloudtemplateIpSecTunnelSourceInterface        = "uni/tn-%s/infranetwork-%s/extnetwork-%s/vpnnetwork-%s/ipsec-[%s]/ipsecsrcif-%s"
	RncloudtemplateIpSecTunnelSourceInterface        = "ipsecsrcif-%s"
	ParentDncloudtemplateIpSecTunnelSourceInterface  = "uni/tn-%s/infranetwork-%s/extnetwork-%s/vpnnetwork-%s/ipsec-[%s]"
	CloudtemplateipsectunnelsourceinterfaceClassName = "cloudtemplateIpSecTunnelSourceInterface"
)

type CloudTemplateforIpSectunnelSourceInterface struct {
	BaseAttributes
	CloudTemplateforIpSectunnelSourceInterfaceAttributes
}

type CloudTemplateforIpSectunnelSourceInterfaceAttributes struct {
	Annotation        string `json:",omitempty"`
	SourceInterfaceId string `json:",omitempty"`
}

func NewCloudTemplateIpSecTunnelSourceInterface(cloudtemplateIpSecTunnelSourceInterfaceRn, parentDn string, cloudtemplateIpSecTunnelSourceInterfaceAttr CloudTemplateforIpSectunnelSourceInterfaceAttributes) *CloudTemplateforIpSectunnelSourceInterface {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudtemplateIpSecTunnelSourceInterfaceRn)
	return &CloudTemplateforIpSectunnelSourceInterface{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         CloudtemplateipsectunnelsourceinterfaceClassName,
			Rn:                cloudtemplateIpSecTunnelSourceInterfaceRn,
		},
		CloudTemplateforIpSectunnelSourceInterfaceAttributes: cloudtemplateIpSecTunnelSourceInterfaceAttr,
	}
}

func (cloudtemplateIpSecTunnelSourceInterface *CloudTemplateforIpSectunnelSourceInterface) ToMap() (map[string]string, error) {
	cloudtemplateIpSecTunnelSourceInterfaceMap, err := cloudtemplateIpSecTunnelSourceInterface.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudtemplateIpSecTunnelSourceInterfaceMap, "annotation", cloudtemplateIpSecTunnelSourceInterface.Annotation)
	A(cloudtemplateIpSecTunnelSourceInterfaceMap, "sourceInterfaceId", cloudtemplateIpSecTunnelSourceInterface.SourceInterfaceId)
	return cloudtemplateIpSecTunnelSourceInterfaceMap, err
}

func CloudTemplateIpSecTunnelSourceInterfaceList(cont *container.Container, index int) *CloudTemplateforIpSectunnelSourceInterface {
	CloudTemplateforIpSectunnelSourceInterfaceCont := cont.S("imdata").Index(index).S(CloudtemplateipsectunnelsourceinterfaceClassName, "attributes")
	return &CloudTemplateforIpSectunnelSourceInterface{
		BaseAttributes{
			DistinguishedName: G(CloudTemplateforIpSectunnelSourceInterfaceCont, "dn"),
			Status:            G(CloudTemplateforIpSectunnelSourceInterfaceCont, "status"),
			ClassName:         CloudtemplateipsectunnelsourceinterfaceClassName,
			Rn:                G(CloudTemplateforIpSectunnelSourceInterfaceCont, "rn"),
		},
		CloudTemplateforIpSectunnelSourceInterfaceAttributes{
			Annotation:        G(CloudTemplateforIpSectunnelSourceInterfaceCont, "annotation"),
			SourceInterfaceId: G(CloudTemplateforIpSectunnelSourceInterfaceCont, "sourceInterfaceId"),
		},
	}
}

func CloudTemplateIpSecTunnelSourceInterfaceFromContainer(cont *container.Container) *CloudTemplateforIpSectunnelSourceInterface {
	return CloudTemplateIpSecTunnelSourceInterfaceList(cont, 0)
}

func CloudTemplateIpSecTunnelSourceInterfaceListFromContainer(cont *container.Container) []*CloudTemplateforIpSectunnelSourceInterface {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudTemplateforIpSectunnelSourceInterface, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudTemplateIpSecTunnelSourceInterfaceList(cont, i)
	}

	return arr
}

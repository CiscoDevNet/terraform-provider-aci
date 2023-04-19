package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DncloudtemplateIpSecTunnelSubnetPool        = "uni/tn-%s/infranetwork-%s/ipsecsubnetpool-[%s]"
	RncloudtemplateIpSecTunnelSubnetPool        = "ipsecsubnetpool-[%s]"
	ParentDncloudtemplateIpSecTunnelSubnetPool  = "uni/tn-%s/infranetwork-%s"
	CloudtemplateipsectunnelsubnetpoolClassName = "cloudtemplateIpSecTunnelSubnetPool"
)

type SubnetPoolforIpSecTunnels struct {
	BaseAttributes
	SubnetPoolforIpSecTunnelsAttributes
}

type SubnetPoolforIpSecTunnelsAttributes struct {
	Annotation string `json:",omitempty"`
	Poolname   string `json:",omitempty"`
	Subnetpool string `json:",omitempty"`
}

func NewSubnetPoolforIpSecTunnels(cloudtemplateIpSecTunnelSubnetPoolRn, parentDn string, cloudtemplateIpSecTunnelSubnetPoolAttr SubnetPoolforIpSecTunnelsAttributes) *SubnetPoolforIpSecTunnels {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudtemplateIpSecTunnelSubnetPoolRn)
	return &SubnetPoolforIpSecTunnels{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         CloudtemplateipsectunnelsubnetpoolClassName,
			Rn:                cloudtemplateIpSecTunnelSubnetPoolRn,
		},
		SubnetPoolforIpSecTunnelsAttributes: cloudtemplateIpSecTunnelSubnetPoolAttr,
	}
}

func (cloudtemplateIpSecTunnelSubnetPool *SubnetPoolforIpSecTunnels) ToMap() (map[string]string, error) {
	cloudtemplateIpSecTunnelSubnetPoolMap, err := cloudtemplateIpSecTunnelSubnetPool.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudtemplateIpSecTunnelSubnetPoolMap, "annotation", cloudtemplateIpSecTunnelSubnetPool.Annotation)
	A(cloudtemplateIpSecTunnelSubnetPoolMap, "poolname", cloudtemplateIpSecTunnelSubnetPool.Poolname)
	A(cloudtemplateIpSecTunnelSubnetPoolMap, "subnetpool", cloudtemplateIpSecTunnelSubnetPool.Subnetpool)
	return cloudtemplateIpSecTunnelSubnetPoolMap, err
}

func SubnetPoolforIpSecTunnelsFromContainerList(cont *container.Container, index int) *SubnetPoolforIpSecTunnels {
	SubnetPoolforIpSecTunnelsCont := cont.S("imdata").Index(index).S(CloudtemplateipsectunnelsubnetpoolClassName, "attributes")
	return &SubnetPoolforIpSecTunnels{
		BaseAttributes{
			DistinguishedName: G(SubnetPoolforIpSecTunnelsCont, "dn"),
			Status:            G(SubnetPoolforIpSecTunnelsCont, "status"),
			ClassName:         CloudtemplateipsectunnelsubnetpoolClassName,
			Rn:                G(SubnetPoolforIpSecTunnelsCont, "rn"),
		},
		SubnetPoolforIpSecTunnelsAttributes{
			Annotation: G(SubnetPoolforIpSecTunnelsCont, "annotation"),
			Poolname:   G(SubnetPoolforIpSecTunnelsCont, "poolname"),
			Subnetpool: G(SubnetPoolforIpSecTunnelsCont, "subnetpool"),
		},
	}
}

func SubnetPoolforIpSecTunnelsFromContainer(cont *container.Container) *SubnetPoolforIpSecTunnels {
	return SubnetPoolforIpSecTunnelsFromContainerList(cont, 0)
}

func SubnetPoolforIpSecTunnelsListFromContainer(cont *container.Container) []*SubnetPoolforIpSecTunnels {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*SubnetPoolforIpSecTunnels, length)

	for i := 0; i < length; i++ {
		arr[i] = SubnetPoolforIpSecTunnelsFromContainerList(cont, i)
	}

	return arr
}

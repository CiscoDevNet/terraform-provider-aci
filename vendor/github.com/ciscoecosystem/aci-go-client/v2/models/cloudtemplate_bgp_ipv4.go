package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DncloudtemplateBgpIpv4        = "uni/tn-%s/infranetwork-%s/extnetwork-%s/vpnnetwork-%s/ipsec-[%s]/bgpipv4-[%s]"
	RncloudtemplateBgpIpv4        = "bgpipv4-[%s]"
	ParentDncloudtemplateBgpIpv4  = "uni/tn-%s/infranetwork-%s/extnetwork-%s/vpnnetwork-%s/ipsec-[%s]"
	Cloudtemplatebgpipv4ClassName = "cloudtemplateBgpIpv4"
)

type CloudTemplateBGPIPv4Peer struct {
	BaseAttributes
	CloudTemplateBGPIPv4PeerAttributes
}

type CloudTemplateBGPIPv4PeerAttributes struct {
	Annotation string `json:",omitempty"`
	Peeraddr   string `json:",omitempty"`
	Peerasn    string `json:",omitempty"`
	Asn        string `json:",omitempty"`
}

func NewCloudTemplateBGPIPv4Peer(cloudtemplateBgpIpv4Rn, parentDn string, cloudtemplateBgpIpv4Attr CloudTemplateBGPIPv4PeerAttributes) *CloudTemplateBGPIPv4Peer {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudtemplateBgpIpv4Rn)
	return &CloudTemplateBGPIPv4Peer{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         Cloudtemplatebgpipv4ClassName,
			Rn:                cloudtemplateBgpIpv4Rn,
		},
		CloudTemplateBGPIPv4PeerAttributes: cloudtemplateBgpIpv4Attr,
	}
}

func (cloudtemplateBgpIpv4 *CloudTemplateBGPIPv4Peer) ToMap() (map[string]string, error) {
	cloudtemplateBgpIpv4Map, err := cloudtemplateBgpIpv4.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudtemplateBgpIpv4Map, "annotation", cloudtemplateBgpIpv4.Annotation)
	A(cloudtemplateBgpIpv4Map, "peeraddr", cloudtemplateBgpIpv4.Peeraddr)
	A(cloudtemplateBgpIpv4Map, "peerasn", cloudtemplateBgpIpv4.Peerasn)
	A(cloudtemplateBgpIpv4Map, "asn", cloudtemplateBgpIpv4.Asn)
	return cloudtemplateBgpIpv4Map, err
}

func CloudTemplateBGPIPv4PeerFromContainerList(cont *container.Container, index int) *CloudTemplateBGPIPv4Peer {
	CloudTemplateBGPIPv4PeerCont := cont.S("imdata").Index(index).S(Cloudtemplatebgpipv4ClassName, "attributes")
	return &CloudTemplateBGPIPv4Peer{
		BaseAttributes{
			DistinguishedName: G(CloudTemplateBGPIPv4PeerCont, "dn"),
			Status:            G(CloudTemplateBGPIPv4PeerCont, "status"),
			ClassName:         Cloudtemplatebgpipv4ClassName,
			Rn:                G(CloudTemplateBGPIPv4PeerCont, "rn"),
		},
		CloudTemplateBGPIPv4PeerAttributes{
			Annotation: G(CloudTemplateBGPIPv4PeerCont, "annotation"),
			Peeraddr:   G(CloudTemplateBGPIPv4PeerCont, "peeraddr"),
			Peerasn:    G(CloudTemplateBGPIPv4PeerCont, "peerasn"),
			Asn:        G(CloudTemplateBGPIPv4PeerCont, "asn"),
		},
	}
}

func CloudTemplateBGPIPv4PeerFromContainer(cont *container.Container) *CloudTemplateBGPIPv4Peer {
	return CloudTemplateBGPIPv4PeerFromContainerList(cont, 0)
}

func CloudTemplateBGPIPv4PeerListFromContainer(cont *container.Container) []*CloudTemplateBGPIPv4Peer {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudTemplateBGPIPv4Peer, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudTemplateBGPIPv4PeerFromContainerList(cont, i)
	}

	return arr
}

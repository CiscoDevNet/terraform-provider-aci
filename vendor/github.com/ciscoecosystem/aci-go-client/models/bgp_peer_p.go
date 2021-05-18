package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BgppeerpClassName = "bgpPeerP"

type BgpPeerConnectivityProfile struct {
	BaseAttributes
	BgpPeerConnectivityProfileAttributes
}

type BgpPeerConnectivityProfileAttributes struct {
	Addr string `json:",omitempty"`

	AddrTCtrl string `json:",omitempty"`

	AllowedSelfAsCnt string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Ctrl string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Password string `json:",omitempty"`

	PeerCtrl string `json:",omitempty"`

	PrivateASctrl string `json:",omitempty"`

	Ttl string `json:",omitempty"`

	Weight string `json:",omitempty"`
}

func NewBgpPeerConnectivityProfile(bgpPeerPRn, parentDn, description string, bgpPeerPattr BgpPeerConnectivityProfileAttributes) *BgpPeerConnectivityProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, bgpPeerPRn)
	return &BgpPeerConnectivityProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BgppeerpClassName,
			Rn:                bgpPeerPRn,
		},

		BgpPeerConnectivityProfileAttributes: bgpPeerPattr,
	}
}

func (bgpPeerP *BgpPeerConnectivityProfile) ToMap() (map[string]string, error) {
	bgpPeerPMap, err := bgpPeerP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bgpPeerPMap, "addr", bgpPeerP.Addr)

	A(bgpPeerPMap, "addrTCtrl", bgpPeerP.AddrTCtrl)

	A(bgpPeerPMap, "allowedSelfAsCnt", bgpPeerP.AllowedSelfAsCnt)

	A(bgpPeerPMap, "annotation", bgpPeerP.Annotation)

	A(bgpPeerPMap, "ctrl", bgpPeerP.Ctrl)

	A(bgpPeerPMap, "nameAlias", bgpPeerP.NameAlias)

	A(bgpPeerPMap, "password", bgpPeerP.Password)

	A(bgpPeerPMap, "peerCtrl", bgpPeerP.PeerCtrl)

	A(bgpPeerPMap, "privateASctrl", bgpPeerP.PrivateASctrl)

	A(bgpPeerPMap, "ttl", bgpPeerP.Ttl)

	A(bgpPeerPMap, "weight", bgpPeerP.Weight)

	return bgpPeerPMap, err
}

func BgpPeerConnectivityProfileFromContainerList(cont *container.Container, index int) *BgpPeerConnectivityProfile {

	BgpPeerConnectivityProfileCont := cont.S("imdata").Index(index).S(BgppeerpClassName, "attributes")
	return &BgpPeerConnectivityProfile{
		BaseAttributes{
			DistinguishedName: G(BgpPeerConnectivityProfileCont, "dn"),
			Description:       G(BgpPeerConnectivityProfileCont, "descr"),
			Status:            G(BgpPeerConnectivityProfileCont, "status"),
			ClassName:         BgppeerpClassName,
			Rn:                G(BgpPeerConnectivityProfileCont, "rn"),
		},

		BgpPeerConnectivityProfileAttributes{

			Addr: G(BgpPeerConnectivityProfileCont, "addr"),

			AddrTCtrl: G(BgpPeerConnectivityProfileCont, "addrTCtrl"),

			AllowedSelfAsCnt: G(BgpPeerConnectivityProfileCont, "allowedSelfAsCnt"),

			Annotation: G(BgpPeerConnectivityProfileCont, "annotation"),

			Ctrl: G(BgpPeerConnectivityProfileCont, "ctrl"),

			NameAlias: G(BgpPeerConnectivityProfileCont, "nameAlias"),

			Password: G(BgpPeerConnectivityProfileCont, "password"),

			PeerCtrl: G(BgpPeerConnectivityProfileCont, "peerCtrl"),

			PrivateASctrl: G(BgpPeerConnectivityProfileCont, "privateASctrl"),

			Ttl: G(BgpPeerConnectivityProfileCont, "ttl"),

			Weight: G(BgpPeerConnectivityProfileCont, "weight"),
		},
	}
}

func BgpPeerConnectivityProfileFromContainer(cont *container.Container) *BgpPeerConnectivityProfile {

	return BgpPeerConnectivityProfileFromContainerList(cont, 0)
}

func BgpPeerConnectivityProfileListFromContainer(cont *container.Container) []*BgpPeerConnectivityProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*BgpPeerConnectivityProfile, length)

	for i := 0; i < length; i++ {

		arr[i] = BgpPeerConnectivityProfileFromContainerList(cont, i)
	}

	return arr
}

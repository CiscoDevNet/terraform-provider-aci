package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extipClassName = "l3extIp"

type L3outPathAttachmentSecondaryIp struct {
	BaseAttributes
	L3outPathAttachmentSecondaryIpAttributes
}

type L3outPathAttachmentSecondaryIpAttributes struct {
	Addr string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Ipv6Dad string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewL3outPathAttachmentSecondaryIp(l3extIpRn, parentDn, description string, l3extIpattr L3outPathAttachmentSecondaryIpAttributes) *L3outPathAttachmentSecondaryIp {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extIpRn)
	return &L3outPathAttachmentSecondaryIp{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extipClassName,
			Rn:                l3extIpRn,
		},

		L3outPathAttachmentSecondaryIpAttributes: l3extIpattr,
	}
}

func (l3extIp *L3outPathAttachmentSecondaryIp) ToMap() (map[string]string, error) {
	l3extIpMap, err := l3extIp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extIpMap, "addr", l3extIp.Addr)

	A(l3extIpMap, "annotation", l3extIp.Annotation)

	A(l3extIpMap, "ipv6Dad", l3extIp.Ipv6Dad)

	A(l3extIpMap, "nameAlias", l3extIp.NameAlias)

	return l3extIpMap, err
}

func L3outPathAttachmentSecondaryIpFromContainerList(cont *container.Container, index int) *L3outPathAttachmentSecondaryIp {

	L3outPathAttachmentSecondaryIpCont := cont.S("imdata").Index(index).S(L3extipClassName, "attributes")
	return &L3outPathAttachmentSecondaryIp{
		BaseAttributes{
			DistinguishedName: G(L3outPathAttachmentSecondaryIpCont, "dn"),
			Description:       G(L3outPathAttachmentSecondaryIpCont, "descr"),
			Status:            G(L3outPathAttachmentSecondaryIpCont, "status"),
			ClassName:         L3extipClassName,
			Rn:                G(L3outPathAttachmentSecondaryIpCont, "rn"),
		},

		L3outPathAttachmentSecondaryIpAttributes{

			Addr: G(L3outPathAttachmentSecondaryIpCont, "addr"),

			Annotation: G(L3outPathAttachmentSecondaryIpCont, "annotation"),

			Ipv6Dad: G(L3outPathAttachmentSecondaryIpCont, "ipv6Dad"),

			NameAlias: G(L3outPathAttachmentSecondaryIpCont, "nameAlias"),
		},
	}
}

func L3outPathAttachmentSecondaryIpFromContainer(cont *container.Container) *L3outPathAttachmentSecondaryIp {

	return L3outPathAttachmentSecondaryIpFromContainerList(cont, 0)
}

func L3outPathAttachmentSecondaryIpListFromContainer(cont *container.Container) []*L3outPathAttachmentSecondaryIp {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3outPathAttachmentSecondaryIp, length)

	for i := 0; i < length; i++ {

		arr[i] = L3outPathAttachmentSecondaryIpFromContainerList(cont, i)
	}

	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extmemberClassName = "l3extMember"

type L3outVPCMember struct {
	BaseAttributes
	L3outVPCMemberAttributes
}

type L3outVPCMemberAttributes struct {
	Side string `json:",omitempty"`

	Addr string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Ipv6Dad string `json:",omitempty"`

	LlAddr string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewL3outVPCMember(l3extMemberRn, parentDn, description string, l3extMemberattr L3outVPCMemberAttributes) *L3outVPCMember {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extMemberRn)
	return &L3outVPCMember{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extmemberClassName,
			Rn:                l3extMemberRn,
		},

		L3outVPCMemberAttributes: l3extMemberattr,
	}
}

func (l3extMember *L3outVPCMember) ToMap() (map[string]string, error) {
	l3extMemberMap, err := l3extMember.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extMemberMap, "side", l3extMember.Side)

	A(l3extMemberMap, "addr", l3extMember.Addr)

	A(l3extMemberMap, "annotation", l3extMember.Annotation)

	A(l3extMemberMap, "ipv6Dad", l3extMember.Ipv6Dad)

	A(l3extMemberMap, "llAddr", l3extMember.LlAddr)

	A(l3extMemberMap, "nameAlias", l3extMember.NameAlias)

	return l3extMemberMap, err
}

func L3outVPCMemberFromContainerList(cont *container.Container, index int) *L3outVPCMember {

	L3outVPCMemberCont := cont.S("imdata").Index(index).S(L3extmemberClassName, "attributes")
	return &L3outVPCMember{
		BaseAttributes{
			DistinguishedName: G(L3outVPCMemberCont, "dn"),
			Description:       G(L3outVPCMemberCont, "descr"),
			Status:            G(L3outVPCMemberCont, "status"),
			ClassName:         L3extmemberClassName,
			Rn:                G(L3outVPCMemberCont, "rn"),
		},

		L3outVPCMemberAttributes{

			Side: G(L3outVPCMemberCont, "side"),

			Addr: G(L3outVPCMemberCont, "addr"),

			Annotation: G(L3outVPCMemberCont, "annotation"),

			Ipv6Dad: G(L3outVPCMemberCont, "ipv6Dad"),

			LlAddr: G(L3outVPCMemberCont, "llAddr"),

			NameAlias: G(L3outVPCMemberCont, "nameAlias"),
		},
	}
}

func L3outVPCMemberFromContainer(cont *container.Container) *L3outVPCMember {

	return L3outVPCMemberFromContainerList(cont, 0)
}

func L3outVPCMemberListFromContainer(cont *container.Container) []*L3outVPCMember {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3outVPCMember, length)

	for i := 0; i < length; i++ {

		arr[i] = L3outVPCMemberFromContainerList(cont, i)
	}

	return arr
}

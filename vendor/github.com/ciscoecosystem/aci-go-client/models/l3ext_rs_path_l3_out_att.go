package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L3extrspathl3outattClassName = "l3extRsPathL3OutAtt"

type L3outPathAttachment struct {
	BaseAttributes
	L3outPathAttachmentAttributes
}

type L3outPathAttachmentAttributes struct {
	TDn string `json:",omitempty"`

	Addr string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Autostate string `json:",omitempty"`

	Encap string `json:",omitempty"`

	EncapScope string `json:",omitempty"`

	IfInstT string `json:",omitempty"`

	Ipv6Dad string `json:",omitempty"`

	LlAddr string `json:",omitempty"`

	Mac string `json:",omitempty"`

	Mode string `json:",omitempty"`

	Mtu string `json:",omitempty"`

	TargetDscp string `json:",omitempty"`
}

func NewL3outPathAttachment(l3extRsPathL3OutAttRn, parentDn, description string, l3extRsPathL3OutAttattr L3outPathAttachmentAttributes) *L3outPathAttachment {
	dn := fmt.Sprintf("%s/%s", parentDn, l3extRsPathL3OutAttRn)
	return &L3outPathAttachment{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extrspathl3outattClassName,
			Rn:                l3extRsPathL3OutAttRn,
		},

		L3outPathAttachmentAttributes: l3extRsPathL3OutAttattr,
	}
}

func (l3extRsPathL3OutAtt *L3outPathAttachment) ToMap() (map[string]string, error) {
	l3extRsPathL3OutAttMap, err := l3extRsPathL3OutAtt.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l3extRsPathL3OutAttMap, "tDn", l3extRsPathL3OutAtt.TDn)

	A(l3extRsPathL3OutAttMap, "addr", l3extRsPathL3OutAtt.Addr)

	A(l3extRsPathL3OutAttMap, "annotation", l3extRsPathL3OutAtt.Annotation)

	A(l3extRsPathL3OutAttMap, "autostate", l3extRsPathL3OutAtt.Autostate)

	A(l3extRsPathL3OutAttMap, "encap", l3extRsPathL3OutAtt.Encap)

	A(l3extRsPathL3OutAttMap, "encapScope", l3extRsPathL3OutAtt.EncapScope)

	A(l3extRsPathL3OutAttMap, "ifInstT", l3extRsPathL3OutAtt.IfInstT)

	A(l3extRsPathL3OutAttMap, "ipv6Dad", l3extRsPathL3OutAtt.Ipv6Dad)

	A(l3extRsPathL3OutAttMap, "llAddr", l3extRsPathL3OutAtt.LlAddr)

	A(l3extRsPathL3OutAttMap, "mac", l3extRsPathL3OutAtt.Mac)

	A(l3extRsPathL3OutAttMap, "mode", l3extRsPathL3OutAtt.Mode)

	A(l3extRsPathL3OutAttMap, "mtu", l3extRsPathL3OutAtt.Mtu)

	A(l3extRsPathL3OutAttMap, "targetDscp", l3extRsPathL3OutAtt.TargetDscp)

	return l3extRsPathL3OutAttMap, err
}

func L3outPathAttachmentFromContainerList(cont *container.Container, index int) *L3outPathAttachment {

	L3outPathAttachmentCont := cont.S("imdata").Index(index).S(L3extrspathl3outattClassName, "attributes")
	return &L3outPathAttachment{
		BaseAttributes{
			DistinguishedName: G(L3outPathAttachmentCont, "dn"),
			Description:       G(L3outPathAttachmentCont, "descr"),
			Status:            G(L3outPathAttachmentCont, "status"),
			ClassName:         L3extrspathl3outattClassName,
			Rn:                G(L3outPathAttachmentCont, "rn"),
		},

		L3outPathAttachmentAttributes{

			TDn: G(L3outPathAttachmentCont, "tDn"),

			Addr: G(L3outPathAttachmentCont, "addr"),

			Annotation: G(L3outPathAttachmentCont, "annotation"),

			Autostate: G(L3outPathAttachmentCont, "autostate"),

			Encap: G(L3outPathAttachmentCont, "encap"),

			EncapScope: G(L3outPathAttachmentCont, "encapScope"),

			IfInstT: G(L3outPathAttachmentCont, "ifInstT"),

			Ipv6Dad: G(L3outPathAttachmentCont, "ipv6Dad"),

			LlAddr: G(L3outPathAttachmentCont, "llAddr"),

			Mac: G(L3outPathAttachmentCont, "mac"),

			Mode: G(L3outPathAttachmentCont, "mode"),

			Mtu: G(L3outPathAttachmentCont, "mtu"),

			TargetDscp: G(L3outPathAttachmentCont, "targetDscp"),
		},
	}
}

func L3outPathAttachmentFromContainer(cont *container.Container) *L3outPathAttachment {

	return L3outPathAttachmentFromContainerList(cont, 0)
}

func L3outPathAttachmentListFromContainer(cont *container.Container) []*L3outPathAttachment {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3outPathAttachment, length)

	for i := 0; i < length; i++ {

		arr[i] = L3outPathAttachmentFromContainerList(cont, i)
	}

	return arr
}

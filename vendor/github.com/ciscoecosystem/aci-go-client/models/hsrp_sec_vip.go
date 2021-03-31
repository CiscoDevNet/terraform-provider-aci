package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const HsrpsecvipClassName = "hsrpSecVip"

type L3outHSRPSecondaryVIP struct {
	BaseAttributes
	L3outHSRPSecondaryVIPAttributes
}

type L3outHSRPSecondaryVIPAttributes struct {
	Ip string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ConfigIssues string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewL3outHSRPSecondaryVIP(hsrpSecVipRn, parentDn, description string, hsrpSecVipattr L3outHSRPSecondaryVIPAttributes) *L3outHSRPSecondaryVIP {
	dn := fmt.Sprintf("%s/%s", parentDn, hsrpSecVipRn)
	return &L3outHSRPSecondaryVIP{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         HsrpsecvipClassName,
			Rn:                hsrpSecVipRn,
		},

		L3outHSRPSecondaryVIPAttributes: hsrpSecVipattr,
	}
}

func (hsrpSecVip *L3outHSRPSecondaryVIP) ToMap() (map[string]string, error) {
	hsrpSecVipMap, err := hsrpSecVip.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(hsrpSecVipMap, "ip", hsrpSecVip.Ip)

	A(hsrpSecVipMap, "annotation", hsrpSecVip.Annotation)

	A(hsrpSecVipMap, "configIssues", hsrpSecVip.ConfigIssues)

	A(hsrpSecVipMap, "nameAlias", hsrpSecVip.NameAlias)

	return hsrpSecVipMap, err
}

func L3outHSRPSecondaryVIPFromContainerList(cont *container.Container, index int) *L3outHSRPSecondaryVIP {

	L3outHSRPSecondaryVIPCont := cont.S("imdata").Index(index).S(HsrpsecvipClassName, "attributes")
	return &L3outHSRPSecondaryVIP{
		BaseAttributes{
			DistinguishedName: G(L3outHSRPSecondaryVIPCont, "dn"),
			Description:       G(L3outHSRPSecondaryVIPCont, "descr"),
			Status:            G(L3outHSRPSecondaryVIPCont, "status"),
			ClassName:         HsrpsecvipClassName,
			Rn:                G(L3outHSRPSecondaryVIPCont, "rn"),
		},

		L3outHSRPSecondaryVIPAttributes{

			Ip: G(L3outHSRPSecondaryVIPCont, "ip"),

			Annotation: G(L3outHSRPSecondaryVIPCont, "annotation"),

			ConfigIssues: G(L3outHSRPSecondaryVIPCont, "configIssues"),

			NameAlias: G(L3outHSRPSecondaryVIPCont, "nameAlias"),
		},
	}
}

func L3outHSRPSecondaryVIPFromContainer(cont *container.Container) *L3outHSRPSecondaryVIP {

	return L3outHSRPSecondaryVIPFromContainerList(cont, 0)
}

func L3outHSRPSecondaryVIPListFromContainer(cont *container.Container) []*L3outHSRPSecondaryVIP {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3outHSRPSecondaryVIP, length)

	for i := 0; i < length; i++ {

		arr[i] = L3outHSRPSecondaryVIPFromContainerList(cont, i)
	}

	return arr
}

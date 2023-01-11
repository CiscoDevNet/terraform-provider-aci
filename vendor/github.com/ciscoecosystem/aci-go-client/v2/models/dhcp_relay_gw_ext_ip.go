package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RndhcpRelayGwExtIp        = "relayGwExtIp"
	DhcprelaygwextipClassName = "dhcpRelayGwExtIp"
)

type DhcpRelayGwExtIp struct {
	BaseAttributes
	DhcpRelayGwExtIpAttributes
}

type DhcpRelayGwExtIpAttributes struct {
	Name string `json:",omitempty"`
}

func NewDhcpRelayGwExtIp(dhcpRelayGwExtIpRn, parentDn, description string, dhcpRelayGwExtIpAttr DhcpRelayGwExtIpAttributes) *DhcpRelayGwExtIp {
	dn := fmt.Sprintf("%s/%s", parentDn, dhcpRelayGwExtIpRn)
	return &DhcpRelayGwExtIp{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         DhcprelaygwextipClassName,
			Rn:                dhcpRelayGwExtIpRn,
		},
		DhcpRelayGwExtIpAttributes: dhcpRelayGwExtIpAttr,
	}
}

func (dhcpRelayGwExtIp *DhcpRelayGwExtIp) ToMap() (map[string]string, error) {
	dhcpRelayGwExtIpMap, err := dhcpRelayGwExtIp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(dhcpRelayGwExtIpMap, "name", dhcpRelayGwExtIp.Name)
	return dhcpRelayGwExtIpMap, err
}

func DhcpRelayGwExtIpFromContainerList(cont *container.Container, index int) *DhcpRelayGwExtIp {
	DhcpRelayGwExtIpCont := cont.S("imdata").Index(index).S(DhcprelaygwextipClassName, "attributes")
	return &DhcpRelayGwExtIp{
		BaseAttributes{
			DistinguishedName: G(DhcpRelayGwExtIpCont, "dn"),
			Description:       G(DhcpRelayGwExtIpCont, "descr"),
			Status:            G(DhcpRelayGwExtIpCont, "status"),
			ClassName:         DhcprelaygwextipClassName,
			Rn:                G(DhcpRelayGwExtIpCont, "rn"),
		},
		DhcpRelayGwExtIpAttributes{
			Name: G(DhcpRelayGwExtIpCont, "name"),
		},
	}
}

func DhcpRelayGwExtIpFromContainer(cont *container.Container) *DhcpRelayGwExtIp {
	return DhcpRelayGwExtIpFromContainerList(cont, 0)
}

func DhcpRelayGwExtIpListFromContainer(cont *container.Container) []*DhcpRelayGwExtIp {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*DhcpRelayGwExtIp, length)

	for i := 0; i < length; i++ {
		arr[i] = DhcpRelayGwExtIpFromContainerList(cont, i)
	}

	return arr
}

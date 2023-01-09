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

type UsetheexternalsecondaryaddressforDHCPrelaygateway struct {
	BaseAttributes
	UsetheexternalsecondaryaddressforDHCPrelaygatewayAttributes
}

type UsetheexternalsecondaryaddressforDHCPrelaygatewayAttributes struct {
	Name string `json:",omitempty"`
}

func NewUsetheexternalsecondaryaddressforDHCPrelaygateway(dhcpRelayGwExtIpRn, parentDn, description string, dhcpRelayGwExtIpAttr UsetheexternalsecondaryaddressforDHCPrelaygatewayAttributes) *UsetheexternalsecondaryaddressforDHCPrelaygateway {
	dn := fmt.Sprintf("%s/%s", parentDn, dhcpRelayGwExtIpRn)
	return &UsetheexternalsecondaryaddressforDHCPrelaygateway{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         DhcprelaygwextipClassName,
			Rn:                dhcpRelayGwExtIpRn,
		},
		UsetheexternalsecondaryaddressforDHCPrelaygatewayAttributes: dhcpRelayGwExtIpAttr,
	}
}

func (dhcpRelayGwExtIp *UsetheexternalsecondaryaddressforDHCPrelaygateway) ToMap() (map[string]string, error) {
	dhcpRelayGwExtIpMap, err := dhcpRelayGwExtIp.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(dhcpRelayGwExtIpMap, "name", dhcpRelayGwExtIp.Name)
	return dhcpRelayGwExtIpMap, err
}

func UsetheexternalsecondaryaddressforDHCPrelaygatewayFromContainerList(cont *container.Container, index int) *UsetheexternalsecondaryaddressforDHCPrelaygateway {
	UsetheexternalsecondaryaddressforDHCPrelaygatewayCont := cont.S("imdata").Index(index).S(DhcprelaygwextipClassName, "attributes")
	return &UsetheexternalsecondaryaddressforDHCPrelaygateway{
		BaseAttributes{
			DistinguishedName: G(UsetheexternalsecondaryaddressforDHCPrelaygatewayCont, "dn"),
			Description:       G(UsetheexternalsecondaryaddressforDHCPrelaygatewayCont, "descr"),
			Status:            G(UsetheexternalsecondaryaddressforDHCPrelaygatewayCont, "status"),
			ClassName:         DhcprelaygwextipClassName,
			Rn:                G(UsetheexternalsecondaryaddressforDHCPrelaygatewayCont, "rn"),
		},
		UsetheexternalsecondaryaddressforDHCPrelaygatewayAttributes{
			Name: G(UsetheexternalsecondaryaddressforDHCPrelaygatewayCont, "name"),
		},
	}
}

func UsetheexternalsecondaryaddressforDHCPrelaygatewayFromContainer(cont *container.Container) *UsetheexternalsecondaryaddressforDHCPrelaygateway {
	return UsetheexternalsecondaryaddressforDHCPrelaygatewayFromContainerList(cont, 0)
}

func UsetheexternalsecondaryaddressforDHCPrelaygatewayListFromContainer(cont *container.Container) []*UsetheexternalsecondaryaddressforDHCPrelaygateway {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*UsetheexternalsecondaryaddressforDHCPrelaygateway, length)

	for i := 0; i < length; i++ {
		arr[i] = UsetheexternalsecondaryaddressforDHCPrelaygatewayFromContainerList(cont, i)
	}

	return arr
}

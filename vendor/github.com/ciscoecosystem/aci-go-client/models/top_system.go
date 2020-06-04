package models

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/container"
	"strconv"
)

const TopSystemClassName = "topSystem"

type System struct {
	BaseAttributes
	SystemAttributes
}

type SystemAttributes struct {
	Address                 string `json:",omitempty"`
	BootstrapState          string `json:",omitempty"`
	ChildAction             string `json:",omitempty"`
	ConfigIssues            string `json:",omitempty"`
	ControlPlaneMTU         string `json:",omitempty"`
	CurrentTime             string `json:",omitempty"`
	EnforceSubnetCheck      string `json:",omitempty"`
	EtepAddr                string `json:",omitempty"`
	FabricDomain            string `json:",omitempty"`
	FabricId                string `json:",omitempty"`
	FabricMAC               string `json:",omitempty"`
	Id                      string `json:",omitempty"`
	InbMgmtAddr             string `json:",omitempty"`
	InbMgmtAddr6            string `json:",omitempty"`
	InbMgmtAddr6Mask        string `json:",omitempty"`
	InbMgmtAddrMask         string `json:",omitempty"`
	InbMgmtGateway          string `json:",omitempty"`
	InbMgmtGateway6         string `json:",omitempty"`
	LastRebootTime          string `json:",omitempty"`
	LastResetReason         string `json:",omitempty"`
	LcOwn                   string `json:",omitempty"`
	ModTs                   string `json:",omitempty"`
	Mode                    string `json:",omitempty"`
	MonPolDn                string `json:",omitempty"`
	Name                    string `json:",omitempty"`
	NameAlias               string `json:",omitempty"`
	NodeType                string `json:",omitempty"`
	OobMgmtAddr             string `json:",omitempty"`
	OobMgmtAddr6            string `json:",omitempty"`
	OobMgmtAddr6Mask        string `json:",omitempty"`
	OobMgmtAddrMask         string `json:",omitempty"`
	OobMgmtGateway          string `json:",omitempty"`
	OobMgmtGateway6         string `json:",omitempty"`
	PodId                   string `json:",omitempty"`
	RemoteNetworkId         string `json:",omitempty"`
	RemoteNode              string `json:",omitempty"`
	RlOperPodId             string `json:",omitempty"`
	RlRoutableMode          string `json:",omitempty"`
	RldirectMode            string `json:",omitempty"`
	Role                    string `json:",omitempty"`
	Serial                  string `json:",omitempty"`
	ServerType              string `json:",omitempty"`
	SiteId                  string `json:",omitempty"`
	State                   string `json:",omitempty"`
	SystemUpTime            string `json:",omitempty"`
	TepPool                 string `json:",omitempty"`
	UnicastXrEpLearnDisable string `json:",omitempty"`
	Version                 string `json:",omitempty"`
	VirtualMode             string `json:",omitempty"`
}

func NewSystem(topSystemRn, parentDn, description string, topSystemattr SystemAttributes) *System {
	dn := fmt.Sprintf("%s/%s", parentDn, topSystemRn)
	return &System{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "",
			ClassName:         TopSystemClassName,
			Rn:                topSystemRn,
		},

		SystemAttributes: topSystemattr,
	}
}

func (topSystem *System) ToMap() (map[string]string, error) {
	topSystemMap, err := topSystem.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(topSystemMap, "address", topSystem.Address)
	A(topSystemMap, "bootstrapState", topSystem.BootstrapState)
	A(topSystemMap, "childAction", topSystem.ChildAction)
	A(topSystemMap, "configIssues", topSystem.ConfigIssues)
	A(topSystemMap, "controlPlaneMTU", topSystem.ControlPlaneMTU)
	A(topSystemMap, "currentTime", topSystem.CurrentTime)
	A(topSystemMap, "enforceSubnetCheck", topSystem.EnforceSubnetCheck)
	A(topSystemMap, "etepAddr", topSystem.EtepAddr)
	A(topSystemMap, "fabricDomain", topSystem.FabricDomain)
	A(topSystemMap, "fabricId", topSystem.FabricId)
	A(topSystemMap, "fabricMAC", topSystem.FabricMAC)
	A(topSystemMap, "id", topSystem.Id)
	A(topSystemMap, "inbMgmtAddr", topSystem.InbMgmtAddr)
	A(topSystemMap, "inbMgmtAddr6", topSystem.InbMgmtAddr6)
	A(topSystemMap, "inbMgmtAddr6Mask", topSystem.InbMgmtAddr6Mask)
	A(topSystemMap, "inbMgmtAddrMask", topSystem.InbMgmtAddrMask)
	A(topSystemMap, "inbMgmtGateway", topSystem.InbMgmtGateway)
	A(topSystemMap, "inbMgmtGateway6", topSystem.InbMgmtGateway6)
	A(topSystemMap, "lastRebootTime", topSystem.LastRebootTime)
	A(topSystemMap, "lastResetReason", topSystem.LastResetReason)
	A(topSystemMap, "lcOwn", topSystem.LcOwn)
	A(topSystemMap, "modTs", topSystem.ModTs)
	A(topSystemMap, "mode", topSystem.Mode)
	A(topSystemMap, "monPolDn", topSystem.MonPolDn)
	A(topSystemMap, "name", topSystem.Name)
	A(topSystemMap, "nameAlias", topSystem.NameAlias)
	A(topSystemMap, "nodeType", topSystem.NodeType)
	A(topSystemMap, "oobMgmtAddr", topSystem.OobMgmtAddr)
	A(topSystemMap, "oobMgmtAddr6", topSystem.OobMgmtAddr6)
	A(topSystemMap, "oobMgmtAddr6Mask", topSystem.OobMgmtAddr6Mask)
	A(topSystemMap, "oobMgmtAddrMask", topSystem.OobMgmtAddrMask)
	A(topSystemMap, "oobMgmtGateway", topSystem.OobMgmtGateway)
	A(topSystemMap, "oobMgmtGateway6", topSystem.OobMgmtGateway6)
	A(topSystemMap, "podId", topSystem.PodId)
	A(topSystemMap, "remoteNetworkId", topSystem.RemoteNetworkId)
	A(topSystemMap, "remoteNode", topSystem.RemoteNode)
	A(topSystemMap, "rlOperPodId", topSystem.RlOperPodId)
	A(topSystemMap, "rlRoutableMode", topSystem.RlRoutableMode)
	A(topSystemMap, "rldirectMode", topSystem.RldirectMode)
	A(topSystemMap, "role", topSystem.Role)
	A(topSystemMap, "serial", topSystem.Serial)
	A(topSystemMap, "serverType", topSystem.ServerType)
	A(topSystemMap, "siteId", topSystem.SiteId)
	A(topSystemMap, "state", topSystem.State)
	A(topSystemMap, "systemUpTime", topSystem.SystemUpTime)
	A(topSystemMap, "tepPool", topSystem.TepPool)
	A(topSystemMap, "unicastXrEpLearnDisable", topSystem.UnicastXrEpLearnDisable)
	A(topSystemMap, "version", topSystem.Version)
	A(topSystemMap, "virtualMode", topSystem.VirtualMode)

	return topSystemMap, err
}

func SystemFromContainerList(cont *container.Container, index int) *System {

	SystemCont := cont.S("imdata").Index(index).S(TopSystemClassName, "attributes")
	return &System{
		BaseAttributes{
			DistinguishedName: G(SystemCont, "dn"),
			Status:            G(SystemCont, "status"),
			ClassName:         TopSystemClassName,
			Rn:                G(SystemCont, "rn"),
		},

		SystemAttributes{
			BootstrapState:          G(SystemCont, "bootstrapState"),
			ChildAction:             G(SystemCont, "childAction"),
			ConfigIssues:            G(SystemCont, "configIssues"),
			ControlPlaneMTU:         G(SystemCont, "controlPlaneMTU"),
			CurrentTime:             G(SystemCont, "currentTime"),
			EnforceSubnetCheck:      G(SystemCont, "enforceSubnetCheck"),
			EtepAddr:                G(SystemCont, "etepAddr"),
			FabricDomain:            G(SystemCont, "fabricDomain"),
			FabricId:                G(SystemCont, "fabricId"),
			FabricMAC:               G(SystemCont, "fabricMAC"),
			Id:                      G(SystemCont, "id"),
			InbMgmtAddr:             G(SystemCont, "inbMgmtAddr"),
			InbMgmtAddr6:            G(SystemCont, "inbMgmtAddr6"),
			InbMgmtAddr6Mask:        G(SystemCont, "inbMgmtAddr6Mask"),
			InbMgmtAddrMask:         G(SystemCont, "inbMgmtAddrMask"),
			InbMgmtGateway:          G(SystemCont, "inbMgmtGateway"),
			InbMgmtGateway6:         G(SystemCont, "inbMgmtGateway6"),
			LastRebootTime:          G(SystemCont, "lastRebootTime"),
			LastResetReason:         G(SystemCont, "lastResetReason"),
			LcOwn:                   G(SystemCont, "lcOwn"),
			ModTs:                   G(SystemCont, "modTs"),
			Mode:                    G(SystemCont, "mode"),
			MonPolDn:                G(SystemCont, "monPolDn"),
			Name:                    G(SystemCont, "name"),
			NameAlias:               G(SystemCont, "nameAlias"),
			NodeType:                G(SystemCont, "nodeType"),
			OobMgmtAddr:             G(SystemCont, "oobMgmtAddr"),
			OobMgmtAddr6:            G(SystemCont, "oobMgmtAddr6"),
			OobMgmtAddr6Mask:        G(SystemCont, "oobMgmtAddr6Mask"),
			OobMgmtAddrMask:         G(SystemCont, "oobMgmtAddrMask"),
			OobMgmtGateway:          G(SystemCont, "oobMgmtGateway"),
			OobMgmtGateway6:         G(SystemCont, "oobMgmtGateway6"),
			PodId:                   G(SystemCont, "podId"),
			RemoteNetworkId:         G(SystemCont, "remoteNetworkId"),
			RemoteNode:              G(SystemCont, "remoteNode"),
			RlOperPodId:             G(SystemCont, "rlOperPodId"),
			RlRoutableMode:          G(SystemCont, "rlRoutableMode"),
			RldirectMode:            G(SystemCont, "rldirectMode"),
			Role:                    G(SystemCont, "role"),
			Serial:                  G(SystemCont, "serial"),
			ServerType:              G(SystemCont, "serverType"),
			SiteId:                  G(SystemCont, "siteId"),
			State:                   G(SystemCont, "state"),
			SystemUpTime:            G(SystemCont, "systemUpTime"),
			TepPool:                 G(SystemCont, "tepPool"),
			UnicastXrEpLearnDisable: G(SystemCont, "unicastXrEpLearnDisable"),
			Version:                 G(SystemCont, "version"),
			VirtualMode:             G(SystemCont, "virtualMode"),
		},
	}
}

func SystemFromContainer(cont *container.Container) *System {

	return SystemFromContainerList(cont, 0)
}

func SystemListFromContainer(cont *container.Container) []*System {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*System, length)

	for i := 0; i < length; i++ {

		arr[i] = SystemFromContainerList(cont, i)
	}

	return arr
}

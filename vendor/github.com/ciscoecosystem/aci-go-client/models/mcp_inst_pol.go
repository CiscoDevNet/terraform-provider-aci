package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnmcpInstPol        = "uni/infra/mcpInstP-%s"
	RnmcpInstPol        = "mcpInstP-%s"
	ParentDnmcpInstPol  = "uni/infra"
	McpinstpolClassName = "mcpInstPol"
)

type MiscablingProtocolInstancePolicy struct {
	BaseAttributes
	NameAliasAttribute
	MiscablingProtocolInstancePolicyAttributes
}

type MiscablingProtocolInstancePolicyAttributes struct {
	AdminSt        string `json:",omitempty"`
	Annotation     string `json:",omitempty"`
	Ctrl           string `json:",omitempty"`
	InitDelayTime  string `json:",omitempty"`
	Key            string `json:",omitempty"`
	LoopDetectMult string `json:",omitempty"`
	LoopProtectAct string `json:",omitempty"`
	Name           string `json:",omitempty"`
	TxFreq         string `json:",omitempty"`
	TxFreqMsec     string `json:",omitempty"`
}

func NewMiscablingProtocolInstancePolicy(mcpInstPolRn, parentDn, description, nameAlias string, mcpInstPolAttr MiscablingProtocolInstancePolicyAttributes) *MiscablingProtocolInstancePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, mcpInstPolRn)
	return &MiscablingProtocolInstancePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         McpinstpolClassName,
			Rn:                mcpInstPolRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		MiscablingProtocolInstancePolicyAttributes: mcpInstPolAttr,
	}
}

func (mcpInstPol *MiscablingProtocolInstancePolicy) ToMap() (map[string]string, error) {
	mcpInstPolMap, err := mcpInstPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := mcpInstPol.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(mcpInstPolMap, key, value)
	}
	A(mcpInstPolMap, "adminSt", mcpInstPol.AdminSt)
	A(mcpInstPolMap, "annotation", mcpInstPol.Annotation)
	A(mcpInstPolMap, "ctrl", mcpInstPol.Ctrl)
	A(mcpInstPolMap, "initDelayTime", mcpInstPol.InitDelayTime)
	A(mcpInstPolMap, "key", mcpInstPol.Key)
	A(mcpInstPolMap, "loopDetectMult", mcpInstPol.LoopDetectMult)
	A(mcpInstPolMap, "loopProtectAct", mcpInstPol.LoopProtectAct)
	A(mcpInstPolMap, "name", mcpInstPol.Name)
	A(mcpInstPolMap, "txFreq", mcpInstPol.TxFreq)
	A(mcpInstPolMap, "txFreqMsec", mcpInstPol.TxFreqMsec)
	return mcpInstPolMap, err
}

func MiscablingProtocolInstancePolicyFromContainerList(cont *container.Container, index int) *MiscablingProtocolInstancePolicy {
	MiscablingProtocolInstancePolicyCont := cont.S("imdata").Index(index).S(McpinstpolClassName, "attributes")
	return &MiscablingProtocolInstancePolicy{
		BaseAttributes{
			DistinguishedName: G(MiscablingProtocolInstancePolicyCont, "dn"),
			Description:       G(MiscablingProtocolInstancePolicyCont, "descr"),
			Status:            G(MiscablingProtocolInstancePolicyCont, "status"),
			ClassName:         McpinstpolClassName,
			Rn:                G(MiscablingProtocolInstancePolicyCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(MiscablingProtocolInstancePolicyCont, "nameAlias"),
		},
		MiscablingProtocolInstancePolicyAttributes{
			AdminSt:        G(MiscablingProtocolInstancePolicyCont, "adminSt"),
			Annotation:     G(MiscablingProtocolInstancePolicyCont, "annotation"),
			Ctrl:           G(MiscablingProtocolInstancePolicyCont, "ctrl"),
			InitDelayTime:  G(MiscablingProtocolInstancePolicyCont, "initDelayTime"),
			Key:            G(MiscablingProtocolInstancePolicyCont, "key"),
			LoopDetectMult: G(MiscablingProtocolInstancePolicyCont, "loopDetectMult"),
			LoopProtectAct: G(MiscablingProtocolInstancePolicyCont, "loopProtectAct"),
			Name:           G(MiscablingProtocolInstancePolicyCont, "name"),
			TxFreq:         G(MiscablingProtocolInstancePolicyCont, "txFreq"),
			TxFreqMsec:     G(MiscablingProtocolInstancePolicyCont, "txFreqMsec"),
		},
	}
}

func MiscablingProtocolInstancePolicyFromContainer(cont *container.Container) *MiscablingProtocolInstancePolicy {
	return MiscablingProtocolInstancePolicyFromContainerList(cont, 0)
}

func MiscablingProtocolInstancePolicyListFromContainer(cont *container.Container) []*MiscablingProtocolInstancePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*MiscablingProtocolInstancePolicy, length)
	for i := 0; i < length; i++ {
		arr[i] = MiscablingProtocolInstancePolicyFromContainerList(cont, i)
	}
	return arr
}

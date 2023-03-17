package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnInfraPortConfig        = "portconfnode-%s-card-%s-port-%s-sub-%s"
	DnInfraPortConfig        = "uni/infra/portconfnode-%s-card-%s-port-%s-sub-%s"
	ParentDnInfraPortConfig  = "uni/infra"
	InfraPortConfigClassName = "infraPortConfig"
)

type InfraPortConfiguration struct {
	BaseAttributes
	InfraPortConfigurationAttributes
}

type InfraPortConfigurationAttributes struct {
	Annotation   string `json:",omitempty"`
	AssocGrp     string `json:",omitempty"`
	BrkoutMap    string `json:",omitempty"`
	Card         string `json:",omitempty"`
	ConnectedFex string `json:",omitempty"`
	Descr        string `json:",omitempty"`
	Node         string `json:",omitempty"`
	PcMember     string `json:",omitempty"`
	Port         string `json:",omitempty"`
	Role         string `json:",omitempty"`
	Shutdown     string `json:",omitempty"`
	SubPort      string `json:",omitempty"`
}

func NewInfraPortConfiguration(infraPortConfigRn, parentDn, description string, infraPortConfigAttr InfraPortConfigurationAttributes) *InfraPortConfiguration {
	dn := fmt.Sprintf("%s/%s", parentDn, infraPortConfigRn)
	return &InfraPortConfiguration{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         InfraPortConfigClassName,
			Rn:                infraPortConfigRn,
		},
		InfraPortConfigurationAttributes: infraPortConfigAttr,
	}
}

func (infraPortConfig *InfraPortConfiguration) ToMap() (map[string]string, error) {
	infraPortConfigMap, err := infraPortConfig.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraPortConfigMap, "annotation", infraPortConfig.Annotation)
	A(infraPortConfigMap, "assocGrp", infraPortConfig.AssocGrp)
	A(infraPortConfigMap, "brkoutMap", infraPortConfig.BrkoutMap)
	A(infraPortConfigMap, "card", infraPortConfig.Card)
	A(infraPortConfigMap, "connectedFex", infraPortConfig.ConnectedFex)
	A(infraPortConfigMap, "description", infraPortConfig.Descr)
	A(infraPortConfigMap, "node", infraPortConfig.Node)
	A(infraPortConfigMap, "pcMember", infraPortConfig.PcMember)
	A(infraPortConfigMap, "port", infraPortConfig.Port)
	A(infraPortConfigMap, "role", infraPortConfig.Role)
	A(infraPortConfigMap, "shutdown", infraPortConfig.Shutdown)
	A(infraPortConfigMap, "subPort", infraPortConfig.SubPort)
	return infraPortConfigMap, err
}

func InfraPortConfigurationFromContainerList(cont *container.Container, index int) *InfraPortConfiguration {
	InfraPortConfigurationCont := cont.S("imdata").Index(index).S(InfraPortConfigClassName, "attributes")
	return &InfraPortConfiguration{
		BaseAttributes{
			DistinguishedName: G(InfraPortConfigurationCont, "dn"),
			Status:            G(InfraPortConfigurationCont, "status"),
			ClassName:         InfraPortConfigClassName,
			Rn:                G(InfraPortConfigurationCont, "rn"),
		},
		InfraPortConfigurationAttributes{
			Annotation:   G(InfraPortConfigurationCont, "annotation"),
			AssocGrp:     G(InfraPortConfigurationCont, "assocGrp"),
			BrkoutMap:    G(InfraPortConfigurationCont, "brkoutMap"),
			Card:         G(InfraPortConfigurationCont, "card"),
			ConnectedFex: G(InfraPortConfigurationCont, "connectedFex"),
			Descr:        G(InfraPortConfigurationCont, "description"),
			Node:         G(InfraPortConfigurationCont, "node"),
			PcMember:     G(InfraPortConfigurationCont, "pcMember"),
			Port:         G(InfraPortConfigurationCont, "port"),
			Role:         G(InfraPortConfigurationCont, "role"),
			Shutdown:     G(InfraPortConfigurationCont, "shutdown"),
			SubPort:      G(InfraPortConfigurationCont, "subPort"),
		},
	}
}

func InfraPortConfigurationFromContainer(cont *container.Container) *InfraPortConfiguration {
	return InfraPortConfigurationFromContainerList(cont, 0)
}

func InfraPortConfigurationListFromContainer(cont *container.Container) []*InfraPortConfiguration {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*InfraPortConfiguration, length)

	for i := 0; i < length; i++ {
		arr[i] = InfraPortConfigurationFromContainerList(cont, i)
	}

	return arr
}

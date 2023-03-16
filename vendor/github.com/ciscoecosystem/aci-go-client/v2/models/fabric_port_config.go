package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnFabricPortConfig        = "portconfnode-%s-card-%s-port-%s-sub-%s"
	DnFabricPortConfig        = "uni/fabric/portconfnode-%s-card-%s-port-%s-sub-%s"
	ParentDnFabricPortConfig  = "uni/fabric"
	FabricPortConfigClassName = "fabricPortConfig"
)

type FabricPortConfiguration struct {
	BaseAttributes
	FabricPortConfigurationAttributes
}

type FabricPortConfigurationAttributes struct {
	Annotation string `json:",omitempty"`
	AssocGrp   string `json:",omitempty"`
	Card       string `json:",omitempty"`
	Descr      string `json:",omitempty"`
	Node       string `json:",omitempty"`
	Port       string `json:",omitempty"`
	Role       string `json:",omitempty"`
	Shutdown   string `json:",omitempty"`
	SubPort    string `json:",omitempty"`
}

func NewFabricPortConfiguration(fabricPortConfigRn, parentDn, description string, fabricPortConfigAttr FabricPortConfigurationAttributes) *FabricPortConfiguration {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricPortConfigRn)
	return &FabricPortConfiguration{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         FabricPortConfigClassName,
			Rn:                fabricPortConfigRn,
		},
		FabricPortConfigurationAttributes: fabricPortConfigAttr,
	}
}

func (fabricPortConfig *FabricPortConfiguration) ToMap() (map[string]string, error) {
	fabricPortConfigMap, err := fabricPortConfig.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fabricPortConfigMap, "annotation", fabricPortConfig.Annotation)
	A(fabricPortConfigMap, "assocGrp", fabricPortConfig.AssocGrp)
	A(fabricPortConfigMap, "card", fabricPortConfig.Card)
	A(fabricPortConfigMap, "description", fabricPortConfig.Descr)
	A(fabricPortConfigMap, "node", fabricPortConfig.Node)
	A(fabricPortConfigMap, "port", fabricPortConfig.Port)
	A(fabricPortConfigMap, "role", fabricPortConfig.Role)
	A(fabricPortConfigMap, "shutdown", fabricPortConfig.Shutdown)
	A(fabricPortConfigMap, "subPort", fabricPortConfig.SubPort)
	return fabricPortConfigMap, err
}

func FabricPortConfigurationFromContainerList(cont *container.Container, index int) *FabricPortConfiguration {
	FabricPortConfigurationCont := cont.S("imdata").Index(index).S(FabricPortConfigClassName, "attributes")
	return &FabricPortConfiguration{
		BaseAttributes{
			DistinguishedName: G(FabricPortConfigurationCont, "dn"),
			Status:            G(FabricPortConfigurationCont, "status"),
			ClassName:         FabricPortConfigClassName,
			Rn:                G(FabricPortConfigurationCont, "rn"),
		},
		FabricPortConfigurationAttributes{
			Annotation: G(FabricPortConfigurationCont, "annotation"),
			AssocGrp:   G(FabricPortConfigurationCont, "assocGrp"),
			Card:       G(FabricPortConfigurationCont, "card"),
			Descr:      G(FabricPortConfigurationCont, "description"),
			Node:       G(FabricPortConfigurationCont, "node"),
			Port:       G(FabricPortConfigurationCont, "port"),
			Role:       G(FabricPortConfigurationCont, "role"),
			Shutdown:   G(FabricPortConfigurationCont, "shutdown"),
			SubPort:    G(FabricPortConfigurationCont, "subPort"),
		},
	}
}

func FabricPortConfigurationFromContainer(cont *container.Container) *FabricPortConfiguration {
	return FabricPortConfigurationFromContainerList(cont, 0)
}

func FabricPortConfigurationListFromContainer(cont *container.Container) []*FabricPortConfiguration {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*FabricPortConfiguration, length)

	for i := 0; i < length; i++ {
		arr[i] = FabricPortConfigurationFromContainerList(cont, i)
	}

	return arr
}

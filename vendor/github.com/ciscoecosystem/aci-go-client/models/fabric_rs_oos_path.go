package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnfabricRsOosPath        = "uni/fabric/outofsvc/rsoosPath-[%s]"
	RnfabricRsOosPath        = "rsoosPath-[%s]"
	ParentDnfabricRsOosPath  = "uni/fabric/outofsvc"
	FabricrsoospathClassName = "fabricRsOosPath"
)

type OutofServiceFabricPath struct {
	BaseAttributes
	OutofServiceFabricPathAttributes
}

type OutofServiceFabricPathAttributes struct {
	Annotation string `json:",omitempty"`
	Lc         string `json:",omitempty"`
	TDn        string `json:",omitempty"`
}

func NewOutofServiceFabricPath(fabricRsOosPathRn, parentDn string, fabricRsOosPathAttr OutofServiceFabricPathAttributes) *OutofServiceFabricPath {
	dn := fmt.Sprintf("%s/%s", parentDn, fabricRsOosPathRn)
	return &OutofServiceFabricPath{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         FabricrsoospathClassName,
			Rn:                fabricRsOosPathRn,
		},
		OutofServiceFabricPathAttributes: fabricRsOosPathAttr,
	}
}

func (fabricRsOosPath *OutofServiceFabricPath) ToMap() (map[string]string, error) {
	fabricRsOosPathMap, err := fabricRsOosPath.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range fabricRsOosPathMap {
		A(fabricRsOosPathMap, key, value)
	}
	A(fabricRsOosPathMap, "annotation", fabricRsOosPath.Annotation)
	A(fabricRsOosPathMap, "lc", fabricRsOosPath.Lc)
	A(fabricRsOosPathMap, "tDn", fabricRsOosPath.TDn)
	return fabricRsOosPathMap, err
}

func OutofServiceFabricPathFromContainerList(cont *container.Container, index int) *OutofServiceFabricPath {
	OutofServiceFabricPathCont := cont.S("imdata").Index(index).S(FabricrsoospathClassName, "attributes")
	return &OutofServiceFabricPath{
		BaseAttributes{
			DistinguishedName: G(OutofServiceFabricPathCont, "dn"),
			Status:            G(OutofServiceFabricPathCont, "status"),
			ClassName:         FabricrsoospathClassName,
			Rn:                G(OutofServiceFabricPathCont, "rn"),
		},
		OutofServiceFabricPathAttributes{
			Annotation: G(OutofServiceFabricPathCont, "annotation"),
			Lc:         G(OutofServiceFabricPathCont, "lc"),
			TDn:        G(OutofServiceFabricPathCont, "tDn"),
		},
	}
}

func OutofServiceFabricPathFromContainer(cont *container.Container) *OutofServiceFabricPath {
	return OutofServiceFabricPathFromContainerList(cont, 0)
}

func OutofServiceFabricPathListFromContainer(cont *container.Container) []*OutofServiceFabricPath {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*OutofServiceFabricPath, length)
	for i := 0; i < length; i++ {
		arr[i] = OutofServiceFabricPathFromContainerList(cont, i)
	}
	return arr
}

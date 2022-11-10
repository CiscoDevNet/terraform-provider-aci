package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DninfraRsDomP        = "%s/%s"
	RninfraRsDomP        = "rsdomP-[%s]"
	ParentDninfraRsDomP  = "uni/infra/attentp-%s"
	InfrarsdompClassName = "infraRsDomP"
)

type InfraRsDomP struct {
	BaseAttributes
	InfraRsDomPAttributes
}

type InfraRsDomPAttributes struct {
	Annotation string `json:",omitempty"`
	TDn        string `json:",omitempty"`
}

func NewInfraRsDomP(infraRsDomPRn string, parentDn string, infraRsDomPAttr InfraRsDomPAttributes) *InfraRsDomP {
	dn := fmt.Sprintf("%s/%s", parentDn, infraRsDomPRn)
	return &InfraRsDomP{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         InfrarsdompClassName,
			Rn:                infraRsDomPRn,
		},

		InfraRsDomPAttributes: infraRsDomPAttr,
	}
}

func (infraRsDomP *InfraRsDomP) ToMap() (map[string]string, error) {
	infraRsDomPMap, err := infraRsDomP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraRsDomPMap, "annotation", infraRsDomP.Annotation)
	A(infraRsDomPMap, "tDn", infraRsDomP.TDn)
	return infraRsDomPMap, err
}

func InfraRsDomPFromContainerList(cont *container.Container, index int) *InfraRsDomP {
	InfraRsDomPCont := cont.S("imdata").Index(index).S(InfrarsdompClassName, "attributes")
	return &InfraRsDomP{
		BaseAttributes{
			DistinguishedName: G(InfraRsDomPCont, "dn"),
			Status:            G(InfraRsDomPCont, "status"),
			ClassName:         InfrarsdompClassName,
			Rn:                G(InfraRsDomPCont, "rn"),
		},

		InfraRsDomPAttributes{
			Annotation: G(InfraRsDomPCont, "annotation"),
			TDn:        G(InfraRsDomPCont, "tDn"),
		},
	}
}

func InfraRsDomPFromContainer(cont *container.Container) *InfraRsDomP {
	return InfraRsDomPFromContainerList(cont, 0)
}

func InfraRsDomPListFromContainer(cont *container.Container) []*InfraRsDomP {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*InfraRsDomP, length)

	for i := 0; i < length; i++ {
		arr[i] = InfraRsDomPFromContainerList(cont, i)
	}

	return arr
}

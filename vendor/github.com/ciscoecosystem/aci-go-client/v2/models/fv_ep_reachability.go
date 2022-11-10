package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnfvEpReachability        = "epReach"
	FvepreachabilityClassName = "fvEpReachability"
)

type EpReachability struct {
	BaseAttributes
	EpReachabilityAttributes
}

type EpReachabilityAttributes struct {
	Annotation string `json:",omitempty"`
}

func NewEpReachability(fvEpReachabilityRn, parentDn string, fvEpReachabilityAttr EpReachabilityAttributes) *EpReachability {
	dn := fmt.Sprintf("%s/%s", parentDn, fvEpReachabilityRn)
	return &EpReachability{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         FvepreachabilityClassName,
			Rn:                fvEpReachabilityRn,
		},
	}
}

func (fvEpReachability *EpReachability) ToMap() (map[string]string, error) {
	fvEpReachabilityMap, err := fvEpReachability.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	return fvEpReachabilityMap, err
}

func EpReachabilityFromContainerList(cont *container.Container, index int) *EpReachability {
	EpReachabilityCont := cont.S("imdata").Index(index).S(FvepreachabilityClassName, "attributes")
	return &EpReachability{
		BaseAttributes{
			DistinguishedName: G(EpReachabilityCont, "dn"),
			Status:            G(EpReachabilityCont, "status"),
			ClassName:         FvepreachabilityClassName,
			Rn:                G(EpReachabilityCont, "rn"),
		},
		EpReachabilityAttributes{
			Annotation: G(EpReachabilityCont, "annotation"),
		},
	}
}

func EpReachabilityFromContainer(cont *container.Container) *EpReachability {
	return EpReachabilityFromContainerList(cont, 0)
}

func EpReachabilityListFromContainer(cont *container.Container) []*EpReachability {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*EpReachability, length)

	for i := 0; i < length; i++ {
		arr[i] = EpReachabilityFromContainerList(cont, i)
	}

	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnPsuInstPol        = "uni/fabric/psuInstP-%s"
	RnPsuInstPol        = "psuInstP-%s"
	ParentDnPsuInstPol  = "uni/fabric"
	PsuInstPolClassName = "psuInstPol"
)

type PsuInstPol struct {
	BaseAttributes
	PsuInstPolAttributes
}

type PsuInstPolAttributes struct {
	AdminRdnM  string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Annotation string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewPowerSupplyRedundancyPolicy(psuInstPolRn, description string, psuInstPolAttr PsuInstPolAttributes) *PsuInstPol {
	dn := fmt.Sprintf("%s/%s", ParentDnPsuInstPol, psuInstPolRn)
	return &PsuInstPol{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PsuInstPolClassName,
			Rn:                psuInstPolRn,
		},

		PsuInstPolAttributes: psuInstPolAttr,
	}
}

func (psuInstPol *PsuInstPol) ToMap() (map[string]string, error) {
	psuInstPolMap, err := psuInstPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(psuInstPolMap, "adminRdnM", psuInstPol.AdminRdnM)
	A(psuInstPolMap, "name", psuInstPol.Name)
	A(psuInstPolMap, "annotation", psuInstPol.Annotation)
	A(psuInstPolMap, "nameAlias", psuInstPol.NameAlias)

	return psuInstPolMap, err
}

func PsuInstPolFromContainerList(cont *container.Container, index int) *PsuInstPol {
	PsuInstPolCont := cont.S("imdata").Index(index).S(PsuInstPolClassName, "attributes")
	return &PsuInstPol{
		BaseAttributes{
			DistinguishedName: G(PsuInstPolCont, "dn"),
			Description:       G(PsuInstPolCont, "descr"),
			Status:            G(PsuInstPolCont, "status"),
			ClassName:         PsuInstPolClassName,
			Rn:                G(PsuInstPolCont, "rn"),
		},
		PsuInstPolAttributes{
			AdminRdnM:  G(PsuInstPolCont, "adminRdnM"),
			Name:       G(PsuInstPolCont, "name"),
			Annotation: G(PsuInstPolCont, "annotation"),
			NameAlias:  G(PsuInstPolCont, "nameAlias"),
		},
	}
}

func PsuInstPolFromContainer(cont *container.Container) *PsuInstPol {
	return PsuInstPolFromContainerList(cont, 0)
}

func PsuInstPolListFromContainer(cont *container.Container) []*PsuInstPol {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PsuInstPol, length)
	for i := 0; i < length; i++ {
		arr[i] = PsuInstPolFromContainerList(cont, i)
	}
	return arr
}

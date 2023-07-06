package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnPimExtP        = "pimextp"
	DnPimExtP        = "uni/tn-%s/out-%s/pimextp"
	ParentDnPimExtP  = "uni/tn-%s/out-%s"
	PimExtPClassName = "pimExtP"
)

type PIMExternalProfile struct {
	BaseAttributes
	PIMExternalProfileAttributes
}

type PIMExternalProfileAttributes struct {
	Annotation string `json:",omitempty"`
	EnabledAf  string `json:",omitempty"`
	Name       string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
}

func NewPIMExternalProfile(parentDn, description string, pimExtPAttr PIMExternalProfileAttributes) *PIMExternalProfile {
	dn := fmt.Sprintf("%s/%s", parentDn, RnPimExtP)
	return &PIMExternalProfile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PimExtPClassName,
			Rn:                RnPimExtP,
		},
		PIMExternalProfileAttributes: pimExtPAttr,
	}
}

func (pimExtP *PIMExternalProfile) ToMap() (map[string]string, error) {
	pimExtPMap, err := pimExtP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(pimExtPMap, "annotation", pimExtP.Annotation)
	A(pimExtPMap, "enabledAf", pimExtP.EnabledAf)
	A(pimExtPMap, "name", pimExtP.Name)
	A(pimExtPMap, "nameAlias", pimExtP.NameAlias)
	return pimExtPMap, err
}

func PIMExternalProfileFromContainerList(cont *container.Container, index int) *PIMExternalProfile {
	ExternalProfileCont := cont.S("imdata").Index(index).S(PimExtPClassName, "attributes")
	return &PIMExternalProfile{
		BaseAttributes{
			DistinguishedName: G(ExternalProfileCont, "dn"),
			Description:       G(ExternalProfileCont, "descr"),
			Status:            G(ExternalProfileCont, "status"),
			ClassName:         PimExtPClassName,
			Rn:                G(ExternalProfileCont, "rn"),
		},
		PIMExternalProfileAttributes{
			Annotation: G(ExternalProfileCont, "annotation"),
			EnabledAf:  G(ExternalProfileCont, "enabledAf"),
			Name:       G(ExternalProfileCont, "name"),
			NameAlias:  G(ExternalProfileCont, "nameAlias"),
		},
	}
}

func PIMExternalProfileFromContainer(cont *container.Container) *PIMExternalProfile {
	return PIMExternalProfileFromContainerList(cont, 0)
}

func PIMExternalProfileListFromContainer(cont *container.Container) []*PIMExternalProfile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*PIMExternalProfile, length)

	for i := 0; i < length; i++ {
		arr[i] = PIMExternalProfileFromContainerList(cont, i)
	}

	return arr
}

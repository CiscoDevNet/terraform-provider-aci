package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const MgmtinbClassName = "mgmtInB"

type InBandManagementEPg struct {
	BaseAttributes
	InBandManagementEPgAttributes
}

type InBandManagementEPgAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Encap string `json:",omitempty"`

	ExceptionTag string `json:",omitempty"`

	FloodOnEncap string `json:",omitempty"`

	MatchT string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	PrefGrMemb string `json:",omitempty"`

	Prio string `json:",omitempty"`
}

func NewInBandManagementEPg(mgmtInBRn, parentDn, description string, mgmtInBattr InBandManagementEPgAttributes) *InBandManagementEPg {
	dn := fmt.Sprintf("%s/%s", parentDn, mgmtInBRn)
	return &InBandManagementEPg{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         MgmtinbClassName,
			Rn:                mgmtInBRn,
		},

		InBandManagementEPgAttributes: mgmtInBattr,
	}
}

func (mgmtInB *InBandManagementEPg) ToMap() (map[string]string, error) {
	mgmtInBMap, err := mgmtInB.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(mgmtInBMap, "name", mgmtInB.Name)

	A(mgmtInBMap, "annotation", mgmtInB.Annotation)

	A(mgmtInBMap, "encap", mgmtInB.Encap)

	A(mgmtInBMap, "exceptionTag", mgmtInB.ExceptionTag)

	A(mgmtInBMap, "floodOnEncap", mgmtInB.FloodOnEncap)

	A(mgmtInBMap, "matchT", mgmtInB.MatchT)

	A(mgmtInBMap, "nameAlias", mgmtInB.NameAlias)

	A(mgmtInBMap, "prefGrMemb", mgmtInB.PrefGrMemb)

	A(mgmtInBMap, "prio", mgmtInB.Prio)

	return mgmtInBMap, err
}

func InBandManagementEPgFromContainerList(cont *container.Container, index int) *InBandManagementEPg {

	InBandManagementEPgCont := cont.S("imdata").Index(index).S(MgmtinbClassName, "attributes")
	return &InBandManagementEPg{
		BaseAttributes{
			DistinguishedName: G(InBandManagementEPgCont, "dn"),
			Description:       G(InBandManagementEPgCont, "descr"),
			Status:            G(InBandManagementEPgCont, "status"),
			ClassName:         MgmtinbClassName,
			Rn:                G(InBandManagementEPgCont, "rn"),
		},

		InBandManagementEPgAttributes{

			Name: G(InBandManagementEPgCont, "name"),

			Annotation: G(InBandManagementEPgCont, "annotation"),

			Encap: G(InBandManagementEPgCont, "encap"),

			ExceptionTag: G(InBandManagementEPgCont, "exceptionTag"),

			FloodOnEncap: G(InBandManagementEPgCont, "floodOnEncap"),

			MatchT: G(InBandManagementEPgCont, "matchT"),

			NameAlias: G(InBandManagementEPgCont, "nameAlias"),

			PrefGrMemb: G(InBandManagementEPgCont, "prefGrMemb"),

			Prio: G(InBandManagementEPgCont, "prio"),
		},
	}
}

func InBandManagementEPgFromContainer(cont *container.Container) *InBandManagementEPg {

	return InBandManagementEPgFromContainerList(cont, 0)
}

func InBandManagementEPgListFromContainer(cont *container.Container) []*InBandManagementEPg {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*InBandManagementEPg, length)

	for i := 0; i < length; i++ {

		arr[i] = InBandManagementEPgFromContainerList(cont, i)
	}

	return arr
}

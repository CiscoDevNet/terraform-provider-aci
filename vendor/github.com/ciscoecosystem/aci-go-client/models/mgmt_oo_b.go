package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const MgmtoobClassName = "mgmtOoB"

type OutOfBandManagementEPg struct {
	BaseAttributes
	OutOfBandManagementEPgAttributes
}

type OutOfBandManagementEPgAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Prio string `json:",omitempty"`
}

func NewOutOfBandManagementEPg(mgmtOoBRn, parentDn, description string, mgmtOoBattr OutOfBandManagementEPgAttributes) *OutOfBandManagementEPg {
	dn := fmt.Sprintf("%s/%s", parentDn, mgmtOoBRn)
	return &OutOfBandManagementEPg{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         MgmtoobClassName,
			Rn:                mgmtOoBRn,
		},

		OutOfBandManagementEPgAttributes: mgmtOoBattr,
	}
}

func (mgmtOoB *OutOfBandManagementEPg) ToMap() (map[string]string, error) {
	mgmtOoBMap, err := mgmtOoB.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(mgmtOoBMap, "name", mgmtOoB.Name)

	A(mgmtOoBMap, "annotation", mgmtOoB.Annotation)

	A(mgmtOoBMap, "nameAlias", mgmtOoB.NameAlias)

	A(mgmtOoBMap, "prio", mgmtOoB.Prio)

	return mgmtOoBMap, err
}

func OutOfBandManagementEPgFromContainerList(cont *container.Container, index int) *OutOfBandManagementEPg {

	OutOfBandManagementEPgCont := cont.S("imdata").Index(index).S(MgmtoobClassName, "attributes")
	return &OutOfBandManagementEPg{
		BaseAttributes{
			DistinguishedName: G(OutOfBandManagementEPgCont, "dn"),
			Description:       G(OutOfBandManagementEPgCont, "descr"),
			Status:            G(OutOfBandManagementEPgCont, "status"),
			ClassName:         MgmtoobClassName,
			Rn:                G(OutOfBandManagementEPgCont, "rn"),
		},

		OutOfBandManagementEPgAttributes{

			Name: G(OutOfBandManagementEPgCont, "name"),

			Annotation: G(OutOfBandManagementEPgCont, "annotation"),

			NameAlias: G(OutOfBandManagementEPgCont, "nameAlias"),

			Prio: G(OutOfBandManagementEPgCont, "prio"),
		},
	}
}

func OutOfBandManagementEPgFromContainer(cont *container.Container) *OutOfBandManagementEPg {

	return OutOfBandManagementEPgFromContainerList(cont, 0)
}

func OutOfBandManagementEPgListFromContainer(cont *container.Container) []*OutOfBandManagementEPg {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*OutOfBandManagementEPg, length)

	for i := 0; i < length; i++ {

		arr[i] = OutOfBandManagementEPgFromContainerList(cont, i)
	}

	return arr
}

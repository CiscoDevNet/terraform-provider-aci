package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const L2extinstpClassName = "l2extInstP"

type L2outExternalEpg struct {
	BaseAttributes
	L2outExternalEpgAttributes
}

type L2outExternalEpgAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ExceptionTag string `json:",omitempty"`

	FloodOnEncap string `json:",omitempty"`

	MatchT string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	PrefGrMemb string `json:",omitempty"`

	Prio string `json:",omitempty"`

	TargetDscp string `json:",omitempty"`
}

func NewL2outExternalEpg(l2extInstPRn, parentDn, description string, l2extInstPattr L2outExternalEpgAttributes) *L2outExternalEpg {
	dn := fmt.Sprintf("%s/%s", parentDn, l2extInstPRn)
	return &L2outExternalEpg{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L2extinstpClassName,
			Rn:                l2extInstPRn,
		},

		L2outExternalEpgAttributes: l2extInstPattr,
	}
}

func (l2extInstP *L2outExternalEpg) ToMap() (map[string]string, error) {
	l2extInstPMap, err := l2extInstP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(l2extInstPMap, "name", l2extInstP.Name)

	A(l2extInstPMap, "annotation", l2extInstP.Annotation)

	A(l2extInstPMap, "exceptionTag", l2extInstP.ExceptionTag)

	A(l2extInstPMap, "floodOnEncap", l2extInstP.FloodOnEncap)

	A(l2extInstPMap, "matchT", l2extInstP.MatchT)

	A(l2extInstPMap, "nameAlias", l2extInstP.NameAlias)

	A(l2extInstPMap, "prefGrMemb", l2extInstP.PrefGrMemb)

	A(l2extInstPMap, "prio", l2extInstP.Prio)

	A(l2extInstPMap, "targetDscp", l2extInstP.TargetDscp)

	return l2extInstPMap, err
}

func L2outExternalEpgFromContainerList(cont *container.Container, index int) *L2outExternalEpg {

	L2outExternalEpgCont := cont.S("imdata").Index(index).S(L2extinstpClassName, "attributes")
	return &L2outExternalEpg{
		BaseAttributes{
			DistinguishedName: G(L2outExternalEpgCont, "dn"),
			Description:       G(L2outExternalEpgCont, "descr"),
			Status:            G(L2outExternalEpgCont, "status"),
			ClassName:         L2extinstpClassName,
			Rn:                G(L2outExternalEpgCont, "rn"),
		},

		L2outExternalEpgAttributes{

			Name: G(L2outExternalEpgCont, "name"),

			Annotation: G(L2outExternalEpgCont, "annotation"),

			ExceptionTag: G(L2outExternalEpgCont, "exceptionTag"),

			FloodOnEncap: G(L2outExternalEpgCont, "floodOnEncap"),

			MatchT: G(L2outExternalEpgCont, "matchT"),

			NameAlias: G(L2outExternalEpgCont, "nameAlias"),

			PrefGrMemb: G(L2outExternalEpgCont, "prefGrMemb"),

			Prio: G(L2outExternalEpgCont, "prio"),

			TargetDscp: G(L2outExternalEpgCont, "targetDscp"),
		},
	}
}

func L2outExternalEpgFromContainer(cont *container.Container) *L2outExternalEpg {

	return L2outExternalEpgFromContainerList(cont, 0)
}

func L2outExternalEpgListFromContainer(cont *container.Container) []*L2outExternalEpg {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L2outExternalEpg, length)

	for i := 0; i < length; i++ {

		arr[i] = L2outExternalEpgFromContainerList(cont, i)
	}

	return arr
}

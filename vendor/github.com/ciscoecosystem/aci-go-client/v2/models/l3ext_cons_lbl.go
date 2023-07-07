package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnL3extConsLbl        = "uni/tn-%s/out-%s/conslbl-%s"
	RnL3extConsLbl        = "conslbl-%s"
	ParentDnL3ExtConsLbl  = "uni/tn-%s/out-%s"
	L3extConsLblClassName = "l3extConsLbl"
)

type L3ExtConsLbl struct {
	BaseAttributes
	L3ExtConsLblAttributes
}

type L3ExtConsLblAttributes struct {
	Name       string `json:",omitempty"`
	Annotation string `json:",omitempty"`
	NameAlias  string `json:",omitempty"`
	Owner      string `json:",omitempty"`
	Tag        string `json:",omitempty"`
}

func NewL3ExtConsLbl(consLblRn, parentDn, description string, consLblAttr L3ExtConsLblAttributes) *L3ExtConsLbl {
	dn := fmt.Sprintf("%s/%s", parentDn, consLblRn)
	return &L3ExtConsLbl{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         L3extConsLblClassName,
			Rn:                consLblRn,
		},

		L3ExtConsLblAttributes: consLblAttr,
	}
}

func (consLbl *L3ExtConsLbl) ToMap() (map[string]string, error) {
	consLblMap, err := consLbl.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(consLblMap, "name", consLbl.Name)
	A(consLblMap, "annotation", consLbl.Annotation)
	A(consLblMap, "nameAlias", consLbl.NameAlias)
	A(consLblMap, "owner", consLbl.Owner)
	A(consLblMap, "tag", consLbl.Tag)
	return consLblMap, err
}

func L3ExtConsLblFromContainerList(cont *container.Container, index int) *L3ExtConsLbl {
	L3ExtConsLblCont := cont.S("imdata").Index(index).S(L3extConsLblClassName, "attributes")
	return &L3ExtConsLbl{
		BaseAttributes{
			DistinguishedName: G(L3ExtConsLblCont, "dn"),
			Description:       G(L3ExtConsLblCont, "descr"),
			Status:            G(L3ExtConsLblCont, "status"),
			ClassName:         L3extConsLblClassName,
			Rn:                G(L3ExtConsLblCont, "rn"),
		},

		L3ExtConsLblAttributes{
			Name: G(L3ExtConsLblCont, "name"),
			Annotation: G(L3ExtConsLblCont, "annotation"),
			NameAlias: G(L3ExtConsLblCont, "nameAlias"),
			Owner: G(L3ExtConsLblCont, "owner"),
			Tag: G(L3ExtConsLblCont, "tag"),
		},
	}
}

func L3ExtConsLblFromContainer(cont *container.Container) *L3ExtConsLbl {
	return L3ExtConsLblFromContainerList(cont, 0)
}

func L3ExtConsLblListFromContainer(cont *container.Container) []*L3ExtConsLbl {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3ExtConsLbl, length)
	for i := 0; i < length; i++ {
		arr[i] = L3ExtConsLblFromContainerList(cont, i)
	}

	return arr
}

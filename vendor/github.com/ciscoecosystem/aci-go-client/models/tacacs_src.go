package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DntacacsSrc        = "uni/tn-%s/monepg-%s/tarepg-%s/tacacssrc-%s"
	RntacacsSrc        = "tacacssrc-%s"
	ParentDntacacsSrc  = "uni/tn-%s/monepg-%s/tarepg-%s"
	TacacssrcClassName = "tacacsSrc"
)

type TACACSSource struct {
	BaseAttributes
	NameAliasAttribute
	TACACSSourceAttributes
}

type TACACSSourceAttributes struct {
	Annotation string `json:",omitempty"`
	Incl       string `json:",omitempty"`
	MinSev     string `json:",omitempty"`
	Name       string `json:",omitempty"`
}

func NewTACACSSource(tacacsSrcRn, parentDn, description, nameAlias string, tacacsSrcAttr TACACSSourceAttributes) *TACACSSource {
	dn := fmt.Sprintf("%s/%s", parentDn, tacacsSrcRn)
	return &TACACSSource{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         TacacssrcClassName,
			Rn:                tacacsSrcRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		TACACSSourceAttributes: tacacsSrcAttr,
	}
}

func (tacacsSrc *TACACSSource) ToMap() (map[string]string, error) {
	tacacsSrcMap, err := tacacsSrc.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := tacacsSrc.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(tacacsSrcMap, key, value)
	}
	A(tacacsSrcMap, "annotation", tacacsSrc.Annotation)
	A(tacacsSrcMap, "incl", tacacsSrc.Incl)
	A(tacacsSrcMap, "minSev", tacacsSrc.MinSev)
	A(tacacsSrcMap, "name", tacacsSrc.Name)
	return tacacsSrcMap, err
}

func TACACSSourceFromContainerList(cont *container.Container, index int) *TACACSSource {
	TACACSSourceCont := cont.S("imdata").Index(index).S(TacacssrcClassName, "attributes")
	return &TACACSSource{
		BaseAttributes{
			DistinguishedName: G(TACACSSourceCont, "dn"),
			Description:       G(TACACSSourceCont, "descr"),
			Status:            G(TACACSSourceCont, "status"),
			ClassName:         TacacssrcClassName,
			Rn:                G(TACACSSourceCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(TACACSSourceCont, "nameAlias"),
		},
		TACACSSourceAttributes{
			Annotation: G(TACACSSourceCont, "annotation"),
			Incl:       G(TACACSSourceCont, "incl"),
			MinSev:     G(TACACSSourceCont, "minSev"),
			Name:       G(TACACSSourceCont, "name"),
		},
	}
}

func TACACSSourceFromContainer(cont *container.Container) *TACACSSource {
	return TACACSSourceFromContainerList(cont, 0)
}

func TACACSSourceListFromContainer(cont *container.Container) []*TACACSSource {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*TACACSSource, length)
	for i := 0; i < length; i++ {
		arr[i] = TACACSSourceFromContainerList(cont, i)
	}
	return arr
}

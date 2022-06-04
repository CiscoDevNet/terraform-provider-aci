package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnvnsEPgDef        = "%s/EPgDef-%s"
	RnvnsEPgDef        = "EPgDef-%s"
	VnsepgdefClassName = "vnsEPgDef"
)

type EPgDef struct {
	BaseAttributes
	NameAliasAttribute
	EPgDefAttributes
}

type EPgDefAttributes struct {
	Name        string `json:",omitempty"`
	Encap       string `json:",omitempty"`
	FabricEncap string `json:",omitempty"`
	IsDelPbr    string `json:",omitempty"`
	LIfCtxDn    string `json:",omitempty"`
	MemberType  string `json:",omitempty"`
	RtrId       string `json:",omitempty"`
}

func NewEPgDef(vnsEPgDefRn, parentDn, description, nameAlias string, vnsEPgDefAttr EPgDefAttributes) *EPgDef {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsEPgDefRn)
	return &EPgDef{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsepgdefClassName,
			Rn:                vnsEPgDefRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		EPgDefAttributes: vnsEPgDefAttr,
	}
}

func (vnsEPgDef *EPgDef) ToMap() (map[string]string, error) {
	vnsEPgDefMap, err := vnsEPgDef.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := vnsEPgDef.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(vnsEPgDefMap, key, value)
	}

	A(vnsEPgDefMap, "name", vnsEPgDef.Name)
	A(vnsEPgDefMap, "encap", vnsEPgDef.Encap)
	A(vnsEPgDefMap, "fabEncap", vnsEPgDef.FabricEncap)
	A(vnsEPgDefMap, "isDelPbr", vnsEPgDef.IsDelPbr)
	A(vnsEPgDefMap, "membType", vnsEPgDef.MemberType)
	A(vnsEPgDefMap, "lIfCtxDn", vnsEPgDef.LIfCtxDn)
	A(vnsEPgDefMap, "rtrId", vnsEPgDef.RtrId)
	return vnsEPgDefMap, err
}

func EPgDefFromContainerList(cont *container.Container, index int) *EPgDef {
	EPgDefCont := cont.S("imdata").Index(index).S(VnsepgdefClassName, "attributes")
	return &EPgDef{
		BaseAttributes{
			DistinguishedName: G(EPgDefCont, "dn"),
			Description:       G(EPgDefCont, "descr"),
			Status:            G(EPgDefCont, "status"),
			ClassName:         VnsepgdefClassName,
			Rn:                G(EPgDefCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(EPgDefCont, "nameAlias"),
		},
		EPgDefAttributes{
			Name:        G(EPgDefCont, "name"),
			Encap:       G(EPgDefCont, "encap"),
			FabricEncap: G(EPgDefCont, "fabEncap"),
			IsDelPbr:    G(EPgDefCont, "isDelPbr"),
			MemberType:  G(EPgDefCont, "membType"),
			LIfCtxDn:    G(EPgDefCont, "lIfCtxDn"),
			RtrId:       G(EPgDefCont, "rtrId"),
		},
	}
}

func EPgDefFromContainer(cont *container.Container) *EPgDef {
	return EPgDefFromContainerList(cont, 0)
}

func EPgDefListFromContainer(cont *container.Container) []*EPgDef {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*EPgDef, length)

	for i := 0; i < length; i++ {
		arr[i] = EPgDefFromContainerList(cont, i)
	}

	return arr
}

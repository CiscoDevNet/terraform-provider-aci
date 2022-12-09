package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnvnsL1L2RedirectDest        = "L1L2RedirectDest-%s"
	Vnsl1l2redirectdestClassName = "vnsL1L2RedirectDest"
)

type L1L2RedirectDestTraffic struct {
	BaseAttributes
	NameAliasAttribute
	L1L2RedirectDestTrafficAttributes
}

type L1L2RedirectDestTrafficAttributes struct {
	Annotation string `json:",omitempty"`
	DestName   string `json:",omitempty"`
	Mac        string `json:",omitempty"`
	Name       string `json:",omitempty"`
	PodId      string `json:",omitempty"`
}

func NewL1L2RedirectDestTraffic(vnsL1L2RedirectDestRn, parentDn, description, nameAlias string, vnsL1L2RedirectDestAttr L1L2RedirectDestTrafficAttributes) *L1L2RedirectDestTraffic {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsL1L2RedirectDestRn)
	return &L1L2RedirectDestTraffic{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         Vnsl1l2redirectdestClassName,
			Rn:                vnsL1L2RedirectDestRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		L1L2RedirectDestTrafficAttributes: vnsL1L2RedirectDestAttr,
	}
}

func (vnsL1L2RedirectDest *L1L2RedirectDestTraffic) ToMap() (map[string]string, error) {
	vnsL1L2RedirectDestMap, err := vnsL1L2RedirectDest.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := vnsL1L2RedirectDest.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(vnsL1L2RedirectDestMap, key, value)
	}

	A(vnsL1L2RedirectDestMap, "annotation", vnsL1L2RedirectDest.Annotation)
	A(vnsL1L2RedirectDestMap, "destName", vnsL1L2RedirectDest.DestName)
	A(vnsL1L2RedirectDestMap, "mac", vnsL1L2RedirectDest.Mac)
	A(vnsL1L2RedirectDestMap, "name", vnsL1L2RedirectDest.Name)
	A(vnsL1L2RedirectDestMap, "podId", vnsL1L2RedirectDest.PodId)
	return vnsL1L2RedirectDestMap, err
}

func L1L2RedirectDestTrafficFromContainerList(cont *container.Container, index int) *L1L2RedirectDestTraffic {
	L1L2RedirectDestTrafficCont := cont.S("imdata").Index(index).S(Vnsl1l2redirectdestClassName, "attributes")
	return &L1L2RedirectDestTraffic{
		BaseAttributes{
			DistinguishedName: G(L1L2RedirectDestTrafficCont, "dn"),
			Description:       G(L1L2RedirectDestTrafficCont, "descr"),
			Status:            G(L1L2RedirectDestTrafficCont, "status"),
			ClassName:         Vnsl1l2redirectdestClassName,
			Rn:                G(L1L2RedirectDestTrafficCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(L1L2RedirectDestTrafficCont, "nameAlias"),
		},
		L1L2RedirectDestTrafficAttributes{
			Annotation: G(L1L2RedirectDestTrafficCont, "annotation"),
			DestName:   G(L1L2RedirectDestTrafficCont, "destName"),
			Mac:        G(L1L2RedirectDestTrafficCont, "mac"),
			Name:       G(L1L2RedirectDestTrafficCont, "name"),
			PodId:      G(L1L2RedirectDestTrafficCont, "podId"),
		},
	}
}

func L1L2RedirectDestTrafficFromContainer(cont *container.Container) *L1L2RedirectDestTraffic {
	return L1L2RedirectDestTrafficFromContainerList(cont, 0)
}

func L1L2RedirectDestTrafficListFromContainer(cont *container.Container) []*L1L2RedirectDestTraffic {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*L1L2RedirectDestTraffic, length)

	for i := 0; i < length; i++ {
		arr[i] = L1L2RedirectDestTrafficFromContainerList(cont, i)
	}

	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvnsencapblkClassName = "fvnsEncapBlk"

type Ranges struct {
	BaseAttributes
	RangesAttributes
}

type RangesAttributes struct {
	_from string `json:",omitempty"`

	To string `json:",omitempty"`

	AllocMode string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	From string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Role string `json:",omitempty"`
}

func NewRanges(fvnsEncapBlkRn, parentDn, description string, fvnsEncapBlkattr RangesAttributes) *Ranges {
	dn := fmt.Sprintf("%s/%s", parentDn, fvnsEncapBlkRn) // Comment
	return &Ranges{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvnsencapblkClassName,
			//Rn:                fvnsEncapBlkRn,
		},

		RangesAttributes: fvnsEncapBlkattr,
	}
}

func (fvnsEncapBlk *Ranges) ToMap() (map[string]string, error) {
	fvnsEncapBlkMap, err := fvnsEncapBlk.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvnsEncapBlkMap, "_from", fvnsEncapBlk._from)

	A(fvnsEncapBlkMap, "to", fvnsEncapBlk.To)

	A(fvnsEncapBlkMap, "allocMode", fvnsEncapBlk.AllocMode)

	A(fvnsEncapBlkMap, "annotation", fvnsEncapBlk.Annotation)

	A(fvnsEncapBlkMap, "from", fvnsEncapBlk.From)

	A(fvnsEncapBlkMap, "nameAlias", fvnsEncapBlk.NameAlias)

	A(fvnsEncapBlkMap, "role", fvnsEncapBlk.Role)

	return fvnsEncapBlkMap, err
}

func RangesFromContainerList(cont *container.Container, index int) *Ranges {

	RangesCont := cont.S("imdata").Index(index).S(FvnsencapblkClassName, "attributes")
	return &Ranges{
		BaseAttributes{
			DistinguishedName: G(RangesCont, "dn"),
			Description:       G(RangesCont, "descr"),
			Status:            G(RangesCont, "status"),
			ClassName:         FvnsencapblkClassName,
			//Rn:                G(RangesCont, "rn"),
		},

		RangesAttributes{

			_from: G(RangesCont, "_from"),

			To: G(RangesCont, "to"),

			AllocMode: G(RangesCont, "allocMode"),

			Annotation: G(RangesCont, "annotation"),

			From: G(RangesCont, "from"),

			NameAlias: G(RangesCont, "nameAlias"),

			Role: G(RangesCont, "role"),
		},
	}
}

func RangesFromContainer(cont *container.Container) *Ranges {

	return RangesFromContainerList(cont, 0)
}

func RangesListFromContainer(cont *container.Container) []*Ranges {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*Ranges, length)

	for i := 0; i < length; i++ {

		arr[i] = RangesFromContainerList(cont, i)
	}

	return arr
}

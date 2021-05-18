package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VnslifctxClassName = "vnsLIfCtx"

type LogicalInterfaceContext struct {
	BaseAttributes
	LogicalInterfaceContextAttributes
}

type LogicalInterfaceContextAttributes struct {
	ConnNameOrLbl string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	L3Dest string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	PermitLog string `json:",omitempty"`
}

func NewLogicalInterfaceContext(vnsLIfCtxRn, parentDn, description string, vnsLIfCtxattr LogicalInterfaceContextAttributes) *LogicalInterfaceContext {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsLIfCtxRn)
	return &LogicalInterfaceContext{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnslifctxClassName,
			Rn:                vnsLIfCtxRn,
		},

		LogicalInterfaceContextAttributes: vnsLIfCtxattr,
	}
}

func (vnsLIfCtx *LogicalInterfaceContext) ToMap() (map[string]string, error) {
	vnsLIfCtxMap, err := vnsLIfCtx.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vnsLIfCtxMap, "connNameOrLbl", vnsLIfCtx.ConnNameOrLbl)

	A(vnsLIfCtxMap, "annotation", vnsLIfCtx.Annotation)

	A(vnsLIfCtxMap, "connNameOrLbl", vnsLIfCtx.ConnNameOrLbl)

	A(vnsLIfCtxMap, "l3Dest", vnsLIfCtx.L3Dest)

	A(vnsLIfCtxMap, "nameAlias", vnsLIfCtx.NameAlias)

	A(vnsLIfCtxMap, "permitLog", vnsLIfCtx.PermitLog)

	return vnsLIfCtxMap, err
}

func LogicalInterfaceContextFromContainerList(cont *container.Container, index int) *LogicalInterfaceContext {

	LogicalInterfaceContextCont := cont.S("imdata").Index(index).S(VnslifctxClassName, "attributes")
	return &LogicalInterfaceContext{
		BaseAttributes{
			DistinguishedName: G(LogicalInterfaceContextCont, "dn"),
			Description:       G(LogicalInterfaceContextCont, "descr"),
			Status:            G(LogicalInterfaceContextCont, "status"),
			ClassName:         VnslifctxClassName,
			Rn:                G(LogicalInterfaceContextCont, "rn"),
		},

		LogicalInterfaceContextAttributes{

			ConnNameOrLbl: G(LogicalInterfaceContextCont, "connNameOrLbl"),

			Annotation: G(LogicalInterfaceContextCont, "annotation"),

			L3Dest: G(LogicalInterfaceContextCont, "l3Dest"),

			NameAlias: G(LogicalInterfaceContextCont, "nameAlias"),

			PermitLog: G(LogicalInterfaceContextCont, "permitLog"),
		},
	}
}

func LogicalInterfaceContextFromContainer(cont *container.Container) *LogicalInterfaceContext {

	return LogicalInterfaceContextFromContainerList(cont, 0)
}

func LogicalInterfaceContextListFromContainer(cont *container.Container) []*LogicalInterfaceContext {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LogicalInterfaceContext, length)

	for i := 0; i < length; i++ {

		arr[i] = LogicalInterfaceContextFromContainerList(cont, i)
	}

	return arr
}

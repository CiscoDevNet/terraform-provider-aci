package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VnsldevctxClassName = "vnsLDevCtx"

type LogicalDeviceContext struct {
	BaseAttributes
	LogicalDeviceContextAttributes
}

type LogicalDeviceContextAttributes struct {
	CtrctNameOrLbl string `json:",omitempty"`

	GraphNameOrLbl string `json:",omitempty"`

	NodeNameOrLbl string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Context string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewLogicalDeviceContext(vnsLDevCtxRn, parentDn, description string, vnsLDevCtxattr LogicalDeviceContextAttributes) *LogicalDeviceContext {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsLDevCtxRn)
	return &LogicalDeviceContext{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsldevctxClassName,
			Rn:                vnsLDevCtxRn,
		},

		LogicalDeviceContextAttributes: vnsLDevCtxattr,
	}
}

func (vnsLDevCtx *LogicalDeviceContext) ToMap() (map[string]string, error) {
	vnsLDevCtxMap, err := vnsLDevCtx.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vnsLDevCtxMap, "ctrctNameOrLbl", vnsLDevCtx.CtrctNameOrLbl)

	A(vnsLDevCtxMap, "graphNameOrLbl", vnsLDevCtx.GraphNameOrLbl)

	A(vnsLDevCtxMap, "nodeNameOrLbl", vnsLDevCtx.NodeNameOrLbl)

	A(vnsLDevCtxMap, "annotation", vnsLDevCtx.Annotation)

	A(vnsLDevCtxMap, "context", vnsLDevCtx.Context)

	A(vnsLDevCtxMap, "nameAlias", vnsLDevCtx.NameAlias)

	return vnsLDevCtxMap, err
}

func LogicalDeviceContextFromContainerList(cont *container.Container, index int) *LogicalDeviceContext {

	LogicalDeviceContextCont := cont.S("imdata").Index(index).S(VnsldevctxClassName, "attributes")
	return &LogicalDeviceContext{
		BaseAttributes{
			DistinguishedName: G(LogicalDeviceContextCont, "dn"),
			Description:       G(LogicalDeviceContextCont, "descr"),
			Status:            G(LogicalDeviceContextCont, "status"),
			ClassName:         VnsldevctxClassName,
			Rn:                G(LogicalDeviceContextCont, "rn"),
		},

		LogicalDeviceContextAttributes{

			CtrctNameOrLbl: G(LogicalDeviceContextCont, "ctrctNameOrLbl"),

			GraphNameOrLbl: G(LogicalDeviceContextCont, "graphNameOrLbl"),

			NodeNameOrLbl: G(LogicalDeviceContextCont, "nodeNameOrLbl"),

			Annotation: G(LogicalDeviceContextCont, "annotation"),

			Context: G(LogicalDeviceContextCont, "context"),

			NameAlias: G(LogicalDeviceContextCont, "nameAlias"),
		},
	}
}

func LogicalDeviceContextFromContainer(cont *container.Container) *LogicalDeviceContext {

	return LogicalDeviceContextFromContainerList(cont, 0)
}

func LogicalDeviceContextListFromContainer(cont *container.Container) []*LogicalDeviceContext {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*LogicalDeviceContext, length)

	for i := 0; i < length; i++ {

		arr[i] = LogicalDeviceContextFromContainerList(cont, i)
	}

	return arr
}

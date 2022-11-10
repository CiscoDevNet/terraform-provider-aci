package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlSetRtMetricType        = "uni/tn-%s/attr-%s/smetrict"
	RnrtctrlSetRtMetricType        = "smetrict"
	ParentDnrtctrlSetRtMetricType  = "uni/tn-%s/attr-%s"
	RtctrlsetrtmetrictypeClassName = "rtctrlSetRtMetricType"
)

type RtctrlSetRtMetricType struct {
	BaseAttributes
	NameAliasAttribute
	RtctrlSetRtMetricTypeAttributes
}

type RtctrlSetRtMetricTypeAttributes struct {
	Annotation string `json:",omitempty"`
	MetricType string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Type       string `json:",omitempty"`
}

func NewRtctrlSetRtMetricType(rtctrlSetRtMetricTypeRn, parentDn, description, nameAlias string, rtctrlSetRtMetricTypeAttr RtctrlSetRtMetricTypeAttributes) *RtctrlSetRtMetricType {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetRtMetricTypeRn)
	return &RtctrlSetRtMetricType{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetrtmetrictypeClassName,
			Rn:                rtctrlSetRtMetricTypeRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RtctrlSetRtMetricTypeAttributes: rtctrlSetRtMetricTypeAttr,
	}
}

func (rtctrlSetRtMetricType *RtctrlSetRtMetricType) ToMap() (map[string]string, error) {
	rtctrlSetRtMetricTypeMap, err := rtctrlSetRtMetricType.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetRtMetricType.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetRtMetricTypeMap, key, value)
	}

	A(rtctrlSetRtMetricTypeMap, "annotation", rtctrlSetRtMetricType.Annotation)
	A(rtctrlSetRtMetricTypeMap, "metricType", rtctrlSetRtMetricType.MetricType)
	A(rtctrlSetRtMetricTypeMap, "name", rtctrlSetRtMetricType.Name)
	A(rtctrlSetRtMetricTypeMap, "type", rtctrlSetRtMetricType.Type)
	return rtctrlSetRtMetricTypeMap, err
}

func RtctrlSetRtMetricTypeFromContainerList(cont *container.Container, index int) *RtctrlSetRtMetricType {
	RtctrlSetRtMetricTypeCont := cont.S("imdata").Index(index).S(RtctrlsetrtmetrictypeClassName, "attributes")
	return &RtctrlSetRtMetricType{
		BaseAttributes{
			DistinguishedName: G(RtctrlSetRtMetricTypeCont, "dn"),
			Description:       G(RtctrlSetRtMetricTypeCont, "descr"),
			Status:            G(RtctrlSetRtMetricTypeCont, "status"),
			ClassName:         RtctrlsetrtmetrictypeClassName,
			Rn:                G(RtctrlSetRtMetricTypeCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RtctrlSetRtMetricTypeCont, "nameAlias"),
		},
		RtctrlSetRtMetricTypeAttributes{
			Annotation: G(RtctrlSetRtMetricTypeCont, "annotation"),
			MetricType: G(RtctrlSetRtMetricTypeCont, "metricType"),
			Name:       G(RtctrlSetRtMetricTypeCont, "name"),
			Type:       G(RtctrlSetRtMetricTypeCont, "type"),
		},
	}
}

func RtctrlSetRtMetricTypeFromContainer(cont *container.Container) *RtctrlSetRtMetricType {
	return RtctrlSetRtMetricTypeFromContainerList(cont, 0)
}

func RtctrlSetRtMetricTypeListFromContainer(cont *container.Container) []*RtctrlSetRtMetricType {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RtctrlSetRtMetricType, length)

	for i := 0; i < length; i++ {
		arr[i] = RtctrlSetRtMetricTypeFromContainerList(cont, i)
	}

	return arr
}

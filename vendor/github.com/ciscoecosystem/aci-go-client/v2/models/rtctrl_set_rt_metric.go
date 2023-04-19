package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnrtctrlSetRtMetric        = "uni/tn-%s/attr-%s/smetric"
	RnrtctrlSetRtMetric        = "smetric"
	ParentDnrtctrlSetRtMetric  = "uni/tn-%s/attr-%s"
	RtctrlsetrtmetricClassName = "rtctrlSetRtMetric"
)

type RtctrlSetRtMetric struct {
	BaseAttributes
	NameAliasAttribute
	RtctrlSetRtMetricAttributes
}

type RtctrlSetRtMetricAttributes struct {
	Annotation string `json:",omitempty"`
	Metric     string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Type       string `json:",omitempty"`
}

func NewRtctrlSetRtMetric(rtctrlSetRtMetricRn, parentDn, description, nameAlias string, rtctrlSetRtMetricAttr RtctrlSetRtMetricAttributes) *RtctrlSetRtMetric {
	dn := fmt.Sprintf("%s/%s", parentDn, rtctrlSetRtMetricRn)
	return &RtctrlSetRtMetric{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         RtctrlsetrtmetricClassName,
			Rn:                rtctrlSetRtMetricRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RtctrlSetRtMetricAttributes: rtctrlSetRtMetricAttr,
	}
}

func (rtctrlSetRtMetric *RtctrlSetRtMetric) ToMap() (map[string]string, error) {
	rtctrlSetRtMetricMap, err := rtctrlSetRtMetric.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := rtctrlSetRtMetric.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(rtctrlSetRtMetricMap, key, value)
	}

	A(rtctrlSetRtMetricMap, "annotation", rtctrlSetRtMetric.Annotation)
	A(rtctrlSetRtMetricMap, "metric", rtctrlSetRtMetric.Metric)
	A(rtctrlSetRtMetricMap, "name", rtctrlSetRtMetric.Name)
	A(rtctrlSetRtMetricMap, "type", rtctrlSetRtMetric.Type)
	return rtctrlSetRtMetricMap, err
}

func RtctrlSetRtMetricFromContainerList(cont *container.Container, index int) *RtctrlSetRtMetric {
	RtctrlSetRtMetricCont := cont.S("imdata").Index(index).S(RtctrlsetrtmetricClassName, "attributes")
	return &RtctrlSetRtMetric{
		BaseAttributes{
			DistinguishedName: G(RtctrlSetRtMetricCont, "dn"),
			Description:       G(RtctrlSetRtMetricCont, "descr"),
			Status:            G(RtctrlSetRtMetricCont, "status"),
			ClassName:         RtctrlsetrtmetricClassName,
			Rn:                G(RtctrlSetRtMetricCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RtctrlSetRtMetricCont, "nameAlias"),
		},
		RtctrlSetRtMetricAttributes{
			Annotation: G(RtctrlSetRtMetricCont, "annotation"),
			Metric:     G(RtctrlSetRtMetricCont, "metric"),
			Name:       G(RtctrlSetRtMetricCont, "name"),
			Type:       G(RtctrlSetRtMetricCont, "type"),
		},
	}
}

func RtctrlSetRtMetricFromContainer(cont *container.Container) *RtctrlSetRtMetric {
	return RtctrlSetRtMetricFromContainerList(cont, 0)
}

func RtctrlSetRtMetricListFromContainer(cont *container.Container) []*RtctrlSetRtMetric {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RtctrlSetRtMetric, length)

	for i := 0; i < length; i++ {
		arr[i] = RtctrlSetRtMetricFromContainerList(cont, i)
	}

	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnTagAnnotation        = "%s/annotationKey-[%s]"
	RnTagAnnotation        = "annotationKey-[%s]"
	TagAnnotationClassName = "tagAnnotation"
)

type Annotation struct {
	BaseAttributes
	AnnotationAttributes
}

type AnnotationAttributes struct {
	Key   string `json:",omitempty"`
	Value string `json:",omitempty"`
}

func NewAnnotation(tagAnnotationRn, parentDn string, tagAnnotationAttr AnnotationAttributes) *Annotation {
	dn := fmt.Sprintf("%s/%s", parentDn, tagAnnotationRn)
	return &Annotation{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         TagAnnotationClassName,
			Rn:                tagAnnotationRn,
		},

		AnnotationAttributes: tagAnnotationAttr,
	}
}

func (tagAnnotation *Annotation) ToMap() (map[string]string, error) {
	tagAnnotationMap, err := tagAnnotation.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(tagAnnotationMap, "key", tagAnnotation.Key)
	A(tagAnnotationMap, "value", tagAnnotation.Value)

	return tagAnnotationMap, err
}

func AnnotationFromContainerList(cont *container.Container, index int) *Annotation {
	AnnotationCont := cont.S("imdata").Index(index).S(TagAnnotationClassName, "attributes")
	return &Annotation{
		BaseAttributes{
			DistinguishedName: G(AnnotationCont, "dn"),
			Status:            G(AnnotationCont, "status"),
			ClassName:         TagAnnotationClassName,
		},
		AnnotationAttributes{
			Key:   G(AnnotationCont, "key"),
			Value: G(AnnotationCont, "value"),
		},
	}
}

func AnnotationFromContainer(cont *container.Container) *Annotation {
	return AnnotationFromContainerList(cont, 0)
}

func AnnotationListFromContainer(cont *container.Container) []*Annotation {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*Annotation, length)

	for i := 0; i < length; i++ {
		arr[i] = AnnotationFromContainerList(cont, i)
	}

	return arr
}

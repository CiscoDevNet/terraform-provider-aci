package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DnTagTag        = "%s/tagKey-%s"
	RnTagTag        = "tagKey-%s"
	TagTagClassName = "tagTag"
)

type Tag struct {
	BaseAttributes
	TagAttributes
}

type TagAttributes struct {
	Key   string `json:",omitempty"`
	Value string `json:",omitempty"`
}

func NewTag(tagTagRn, parentDn string, tagTagAttr TagAttributes) *Tag {
	dn := fmt.Sprintf("%s/%s", parentDn, tagTagRn)
	return &Tag{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         TagTagClassName,
			Rn:                tagTagRn,
		},
		TagAttributes: tagTagAttr,
	}
}

func (tagTag *Tag) ToMap() (map[string]string, error) {
	tagTagMap, err := tagTag.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	A(tagTagMap, "key", tagTag.Key)
	A(tagTagMap, "value", tagTag.Value)
	return tagTagMap, err
}

func TagFromContainerList(cont *container.Container, index int) *Tag {
	TagCont := cont.S("imdata").Index(index).S(TagTagClassName, "attributes")
	return &Tag{
		BaseAttributes{
			DistinguishedName: G(TagCont, "dn"),
			Status:            G(TagCont, "status"),
			ClassName:         TagTagClassName,
			Rn:                G(TagCont, "rn"),
		},
		TagAttributes{
			Key:   G(TagCont, "key"),
			Value: G(TagCont, "value"),
		},
	}
}

func TagFromContainer(cont *container.Container) *Tag {
	return TagFromContainerList(cont, 0)
}

func TagListFromContainer(cont *container.Container) []*Tag {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*Tag, length)
	for i := 0; i < length; i++ {
		arr[i] = TagFromContainerList(cont, i)
	}
	return arr
}

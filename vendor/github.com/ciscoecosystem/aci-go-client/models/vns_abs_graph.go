package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const VnsabsgraphClassName = "vnsAbsGraph"

type L4L7ServiceGraphTemplate struct {
	BaseAttributes
	L4L7ServiceGraphTemplateAttributes
}

type L4L7ServiceGraphTemplateAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	L4L7ServiceGraphTemplate_type string `json:",omitempty"`

	UiTemplateType string `json:",omitempty"`
}

func NewL4L7ServiceGraphTemplate(vnsAbsGraphRn, parentDn, description string, vnsAbsGraphattr L4L7ServiceGraphTemplateAttributes) *L4L7ServiceGraphTemplate {
	dn := fmt.Sprintf("%s/%s", parentDn, vnsAbsGraphRn)
	return &L4L7ServiceGraphTemplate{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         VnsabsgraphClassName,
			Rn:                vnsAbsGraphRn,
		},

		L4L7ServiceGraphTemplateAttributes: vnsAbsGraphattr,
	}
}

func (vnsAbsGraph *L4L7ServiceGraphTemplate) ToMap() (map[string]string, error) {
	vnsAbsGraphMap, err := vnsAbsGraph.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(vnsAbsGraphMap, "name", vnsAbsGraph.Name)

	A(vnsAbsGraphMap, "annotation", vnsAbsGraph.Annotation)

	A(vnsAbsGraphMap, "nameAlias", vnsAbsGraph.NameAlias)

	A(vnsAbsGraphMap, "type", vnsAbsGraph.L4L7ServiceGraphTemplate_type)

	A(vnsAbsGraphMap, "uiTemplateType", vnsAbsGraph.UiTemplateType)

	return vnsAbsGraphMap, err
}

func L4L7ServiceGraphTemplateFromContainerList(cont *container.Container, index int) *L4L7ServiceGraphTemplate {

	L4L7ServiceGraphTemplateCont := cont.S("imdata").Index(index).S(VnsabsgraphClassName, "attributes")
	return &L4L7ServiceGraphTemplate{
		BaseAttributes{
			DistinguishedName: G(L4L7ServiceGraphTemplateCont, "dn"),
			Description:       G(L4L7ServiceGraphTemplateCont, "descr"),
			Status:            G(L4L7ServiceGraphTemplateCont, "status"),
			ClassName:         VnsabsgraphClassName,
			Rn:                G(L4L7ServiceGraphTemplateCont, "rn"),
		},

		L4L7ServiceGraphTemplateAttributes{

			Name: G(L4L7ServiceGraphTemplateCont, "name"),

			Annotation: G(L4L7ServiceGraphTemplateCont, "annotation"),

			NameAlias: G(L4L7ServiceGraphTemplateCont, "nameAlias"),

			L4L7ServiceGraphTemplate_type: G(L4L7ServiceGraphTemplateCont, "type"),

			UiTemplateType: G(L4L7ServiceGraphTemplateCont, "uiTemplateType"),
		},
	}
}

func L4L7ServiceGraphTemplateFromContainer(cont *container.Container) *L4L7ServiceGraphTemplate {

	return L4L7ServiceGraphTemplateFromContainerList(cont, 0)
}

func L4L7ServiceGraphTemplateListFromContainer(cont *container.Container) []*L4L7ServiceGraphTemplate {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L4L7ServiceGraphTemplate, length)

	for i := 0; i < length; i++ {

		arr[i] = L4L7ServiceGraphTemplateFromContainerList(cont, i)
	}

	return arr
}

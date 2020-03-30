package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const TrigschedpClassName = "trigSchedP"

type TriggerScheduler struct {
	BaseAttributes
	TriggerSchedulerAttributes
}

type TriggerSchedulerAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewTriggerScheduler(trigSchedPRn, parentDn, description string, trigSchedPattr TriggerSchedulerAttributes) *TriggerScheduler {
	dn := fmt.Sprintf("%s/%s", parentDn, trigSchedPRn)
	return &TriggerScheduler{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         TrigschedpClassName,
			Rn:                trigSchedPRn,
		},

		TriggerSchedulerAttributes: trigSchedPattr,
	}
}

func (trigSchedP *TriggerScheduler) ToMap() (map[string]string, error) {
	trigSchedPMap, err := trigSchedP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(trigSchedPMap, "name", trigSchedP.Name)

	A(trigSchedPMap, "annotation", trigSchedP.Annotation)

	A(trigSchedPMap, "nameAlias", trigSchedP.NameAlias)

	return trigSchedPMap, err
}

func TriggerSchedulerFromContainerList(cont *container.Container, index int) *TriggerScheduler {

	TriggerSchedulerCont := cont.S("imdata").Index(index).S(TrigschedpClassName, "attributes")
	return &TriggerScheduler{
		BaseAttributes{
			DistinguishedName: G(TriggerSchedulerCont, "dn"),
			Description:       G(TriggerSchedulerCont, "descr"),
			Status:            G(TriggerSchedulerCont, "status"),
			ClassName:         TrigschedpClassName,
			Rn:                G(TriggerSchedulerCont, "rn"),
		},

		TriggerSchedulerAttributes{

			Name: G(TriggerSchedulerCont, "name"),

			Annotation: G(TriggerSchedulerCont, "annotation"),

			NameAlias: G(TriggerSchedulerCont, "nameAlias"),
		},
	}
}

func TriggerSchedulerFromContainer(cont *container.Container) *TriggerScheduler {

	return TriggerSchedulerFromContainerList(cont, 0)
}

func TriggerSchedulerListFromContainer(cont *container.Container) []*TriggerScheduler {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*TriggerScheduler, length)

	for i := 0; i < length; i++ {

		arr[i] = TriggerSchedulerFromContainerList(cont, i)
	}

	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DntrigRecurrWindowP        = "uni/fabric/schedp-%s/recurrwinp-%s"
	RntrigRecurrWindowP        = "recurrwinp-%s"
	ParentDntrigRecurrWindowP  = "uni/fabric/schedp-%s"
	TrigrecurrwindowpClassName = "trigRecurrWindowP"
)

type RecurringWindow struct {
	BaseAttributes
	NameAliasAttribute
	RecurringWindowAttributes
}

type RecurringWindowAttributes struct {
	ConcurCap       string `json:",omitempty"`
	Day             string `json:",omitempty"`
	Hour            string `json:",omitempty"`
	Minute          string `json:",omitempty"`
	Name            string `json:",omitempty"`
	NodeUpgInterval string `json:",omitempty"`
	ProcBreak       string `json:",omitempty"`
	ProcCap         string `json:",omitempty"`
	TimeCap         string `json:",omitempty"`
	Annotation      string `json:",omitempty"`
}

func NewRecurringWindow(trigRecurrWindowPRn, parentDn, nameAlias string, trigRecurrWindowPAttr RecurringWindowAttributes) *RecurringWindow {
	dn := fmt.Sprintf("%s/%s", parentDn, trigRecurrWindowPRn)
	return &RecurringWindow{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			//Description:       description,
			Status:    "created, modified",
			ClassName: TrigrecurrwindowpClassName,
			Rn:        trigRecurrWindowPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RecurringWindowAttributes: trigRecurrWindowPAttr,
	}
}

func (trigRecurrWindowP *RecurringWindow) ToMap() (map[string]string, error) {
	trigRecurrWindowPMap, err := trigRecurrWindowP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := trigRecurrWindowP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(trigRecurrWindowPMap, key, value)
	}
	A(trigRecurrWindowPMap, "concurCap", trigRecurrWindowP.ConcurCap)
	A(trigRecurrWindowPMap, "day", trigRecurrWindowP.Day)
	A(trigRecurrWindowPMap, "hour", trigRecurrWindowP.Hour)
	A(trigRecurrWindowPMap, "minute", trigRecurrWindowP.Minute)
	A(trigRecurrWindowPMap, "name", trigRecurrWindowP.Name)
	A(trigRecurrWindowPMap, "nodeUpgInterval", trigRecurrWindowP.NodeUpgInterval)
	A(trigRecurrWindowPMap, "procBreak", trigRecurrWindowP.ProcBreak)
	A(trigRecurrWindowPMap, "procCap", trigRecurrWindowP.ProcCap)
	A(trigRecurrWindowPMap, "timeCap", trigRecurrWindowP.TimeCap)
	A(trigRecurrWindowPMap, "annotation", trigRecurrWindowP.Annotation)
	return trigRecurrWindowPMap, err
}

func RecurringWindowFromContainerList(cont *container.Container, index int) *RecurringWindow {
	RecurringWindowCont := cont.S("imdata").Index(index).S(TrigrecurrwindowpClassName, "attributes")
	return &RecurringWindow{
		BaseAttributes{
			DistinguishedName: G(RecurringWindowCont, "dn"),
			//Description:       G(RecurringWindowCont, "descr"),
			Status:    G(RecurringWindowCont, "status"),
			ClassName: TrigrecurrwindowpClassName,
			Rn:        G(RecurringWindowCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RecurringWindowCont, "nameAlias"),
		},
		RecurringWindowAttributes{
			ConcurCap:       G(RecurringWindowCont, "concurCap"),
			Day:             G(RecurringWindowCont, "day"),
			Hour:            G(RecurringWindowCont, "hour"),
			Minute:          G(RecurringWindowCont, "minute"),
			Name:            G(RecurringWindowCont, "name"),
			NodeUpgInterval: G(RecurringWindowCont, "nodeUpgInterval"),
			ProcBreak:       G(RecurringWindowCont, "procBreak"),
			ProcCap:         G(RecurringWindowCont, "procCap"),
			TimeCap:         G(RecurringWindowCont, "timeCap"),
			Annotation:      G(RecurringWindowCont, "annotation"),
		},
	}
}

func RecurringWindowFromContainer(cont *container.Container) *RecurringWindow {
	return RecurringWindowFromContainerList(cont, 0)
}

func RecurringWindowListFromContainer(cont *container.Container) []*RecurringWindow {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RecurringWindow, length)
	for i := 0; i < length; i++ {
		arr[i] = RecurringWindowFromContainerList(cont, i)
	}
	return arr
}

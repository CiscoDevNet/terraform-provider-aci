package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnedrEventP        = "uni/infra/edrErrDisRecoverPol-%s/edrEventP-%s"
	RnedrEventP        = "edrEventP-%s"
	ParentDnedrEventP  = "uni/infra/edrErrDisRecoverPol-%s"
	EdreventpClassName = "edrEventP"
)

type ErrorDisabledRecoveryEvent struct {
	BaseAttributes
	NameAliasAttribute
	ErrorDisabledRecoveryEventAttributes
}

type ErrorDisabledRecoveryEventAttributes struct {
	Annotation string `json:",omitempty"`
	Event      string `json:",omitempty"`
	Name       string `json:",omitempty"`
	Recover    string `json:",omitempty"`
}

func NewErrorDisabledRecoveryEvent(edrEventPRn, parentDn, description, nameAlias string, edrEventPAttr ErrorDisabledRecoveryEventAttributes) *ErrorDisabledRecoveryEvent {
	dn := fmt.Sprintf("%s/%s", parentDn, edrEventPRn)
	return &ErrorDisabledRecoveryEvent{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         EdreventpClassName,
			Rn:                edrEventPRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		ErrorDisabledRecoveryEventAttributes: edrEventPAttr,
	}
}

func (edrEventP *ErrorDisabledRecoveryEvent) ToMap() (map[string]string, error) {
	edrEventPMap, err := edrEventP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := edrEventP.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(edrEventPMap, key, value)
	}
	A(edrEventPMap, "annotation", edrEventP.Annotation)
	A(edrEventPMap, "event", edrEventP.Event)
	A(edrEventPMap, "name", edrEventP.Name)
	A(edrEventPMap, "recover", edrEventP.Recover)
	return edrEventPMap, err
}

func ErrorDisabledRecoveryEventFromContainerList(cont *container.Container, index int) *ErrorDisabledRecoveryEvent {
	ErrorDisabledRecoveryEventCont := cont.S("imdata").Index(index).S(EdreventpClassName, "attributes")
	return &ErrorDisabledRecoveryEvent{
		BaseAttributes{
			DistinguishedName: G(ErrorDisabledRecoveryEventCont, "dn"),
			Description:       G(ErrorDisabledRecoveryEventCont, "descr"),
			Status:            G(ErrorDisabledRecoveryEventCont, "status"),
			ClassName:         EdreventpClassName,
			Rn:                G(ErrorDisabledRecoveryEventCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(ErrorDisabledRecoveryEventCont, "nameAlias"),
		},
		ErrorDisabledRecoveryEventAttributes{
			Annotation: G(ErrorDisabledRecoveryEventCont, "annotation"),
			Event:      G(ErrorDisabledRecoveryEventCont, "event"),
			Name:       G(ErrorDisabledRecoveryEventCont, "name"),
			Recover:    G(ErrorDisabledRecoveryEventCont, "recover"),
		},
	}
}

func ErrorDisabledRecoveryEventFromContainer(cont *container.Container) *ErrorDisabledRecoveryEvent {
	return ErrorDisabledRecoveryEventFromContainerList(cont, 0)
}

func ErrorDisabledRecoveryEventListFromContainer(cont *container.Container) []*ErrorDisabledRecoveryEvent {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*ErrorDisabledRecoveryEvent, length)
	for i := 0; i < length; i++ {
		arr[i] = ErrorDisabledRecoveryEventFromContainerList(cont, i)
	}
	return arr
}

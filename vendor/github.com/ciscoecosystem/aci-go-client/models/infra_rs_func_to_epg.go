package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const InfrarsfunctoepgClassName = "infraRsFuncToEpg"

type EPGsUsingFunction struct {
	BaseAttributes
	EPGsUsingFunctionAttributes
}

type EPGsUsingFunctionAttributes struct {
	TDn string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Encap string `json:",omitempty"`

	InstrImedcy string `json:",omitempty"`

	Mode string `json:",omitempty"`

	PrimaryEncap string `json:",omitempty"`
}

func NewEPGsUsingFunction(infraRsFuncToEpgRn, parentDn, description string, infraRsFuncToEpgattr EPGsUsingFunctionAttributes) *EPGsUsingFunction {
	dn := fmt.Sprintf("%s/%s", parentDn, infraRsFuncToEpgRn)
	return &EPGsUsingFunction{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         InfrarsfunctoepgClassName,
			Rn:                infraRsFuncToEpgRn,
		},

		EPGsUsingFunctionAttributes: infraRsFuncToEpgattr,
	}
}

func (infraRsFuncToEpg *EPGsUsingFunction) ToMap() (map[string]string, error) {
	infraRsFuncToEpgMap, err := infraRsFuncToEpg.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(infraRsFuncToEpgMap, "tDn", infraRsFuncToEpg.TDn)

	A(infraRsFuncToEpgMap, "annotation", infraRsFuncToEpg.Annotation)

	A(infraRsFuncToEpgMap, "encap", infraRsFuncToEpg.Encap)

	A(infraRsFuncToEpgMap, "instrImedcy", infraRsFuncToEpg.InstrImedcy)

	A(infraRsFuncToEpgMap, "mode", infraRsFuncToEpg.Mode)

	A(infraRsFuncToEpgMap, "primaryEncap", infraRsFuncToEpg.PrimaryEncap)

	return infraRsFuncToEpgMap, err
}

func EPGsUsingFunctionFromContainerList(cont *container.Container, index int) *EPGsUsingFunction {

	EPGsUsingFunctionCont := cont.S("imdata").Index(index).S(InfrarsfunctoepgClassName, "attributes")
	return &EPGsUsingFunction{
		BaseAttributes{
			DistinguishedName: G(EPGsUsingFunctionCont, "dn"),
			Description:       G(EPGsUsingFunctionCont, "descr"),
			Status:            G(EPGsUsingFunctionCont, "status"),
			ClassName:         InfrarsfunctoepgClassName,
			Rn:                G(EPGsUsingFunctionCont, "rn"),
		},

		EPGsUsingFunctionAttributes{

			TDn: G(EPGsUsingFunctionCont, "tDn"),

			Annotation: G(EPGsUsingFunctionCont, "annotation"),

			Encap: G(EPGsUsingFunctionCont, "encap"),

			InstrImedcy: G(EPGsUsingFunctionCont, "instrImedcy"),

			Mode: G(EPGsUsingFunctionCont, "mode"),

			PrimaryEncap: G(EPGsUsingFunctionCont, "primaryEncap"),
		},
	}
}

func EPGsUsingFunctionFromContainer(cont *container.Container) *EPGsUsingFunction {

	return EPGsUsingFunctionFromContainerList(cont, 0)
}

func EPGsUsingFunctionListFromContainer(cont *container.Container) []*EPGsUsingFunction {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*EPGsUsingFunction, length)

	for i := 0; i < length; i++ {

		arr[i] = EPGsUsingFunctionFromContainerList(cont, i)
	}

	return arr
}

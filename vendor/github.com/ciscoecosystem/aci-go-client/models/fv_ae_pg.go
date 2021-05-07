package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvaepgClassName = "fvAEPg"

type ApplicationEPG struct {
	BaseAttributes
	ApplicationEPGAttributes
}

type ApplicationEPGAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ExceptionTag string `json:",omitempty"`

	FloodOnEncap string `json:",omitempty"`

	FwdCtrl string `json:",omitempty"`

	HasMcastSource string `json:",omitempty"`

	IsAttrBasedEPg string `json:",omitempty"`

	MatchT string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	PcEnfPref string `json:",omitempty"`

	PrefGrMemb string `json:",omitempty"`

	Prio string `json:",omitempty"`

	Shutdown string `json:",omitempty"`
}

func NewApplicationEPG(fvAEPgRn, parentDn, description string, fvAEPgattr ApplicationEPGAttributes) *ApplicationEPG {
	dn := fmt.Sprintf("%s/%s", parentDn, fvAEPgRn)
	return &ApplicationEPG{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvaepgClassName,
			Rn:                fvAEPgRn,
		},

		ApplicationEPGAttributes: fvAEPgattr,
	}
}

func (fvAEPg *ApplicationEPG) ToMap() (map[string]string, error) {
	fvAEPgMap, err := fvAEPg.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvAEPgMap, "name", fvAEPg.Name)

	A(fvAEPgMap, "annotation", fvAEPg.Annotation)

	A(fvAEPgMap, "exceptionTag", fvAEPg.ExceptionTag)

	A(fvAEPgMap, "floodOnEncap", fvAEPg.FloodOnEncap)

	A(fvAEPgMap, "fwdCtrl", fvAEPg.FwdCtrl)

	A(fvAEPgMap, "hasMcastSource", fvAEPg.HasMcastSource)

	A(fvAEPgMap, "isAttrBasedEPg", fvAEPg.IsAttrBasedEPg)

	A(fvAEPgMap, "matchT", fvAEPg.MatchT)

	A(fvAEPgMap, "nameAlias", fvAEPg.NameAlias)

	A(fvAEPgMap, "pcEnfPref", fvAEPg.PcEnfPref)

	A(fvAEPgMap, "prefGrMemb", fvAEPg.PrefGrMemb)

	A(fvAEPgMap, "prio", fvAEPg.Prio)

	A(fvAEPgMap, "shutdown", fvAEPg.Shutdown)

	return fvAEPgMap, err
}

func ApplicationEPGFromContainerList(cont *container.Container, index int) *ApplicationEPG {

	ApplicationEPGCont := cont.S("imdata").Index(index).S(FvaepgClassName, "attributes")
	return &ApplicationEPG{
		BaseAttributes{
			DistinguishedName: G(ApplicationEPGCont, "dn"),
			Description:       G(ApplicationEPGCont, "descr"),
			Status:            G(ApplicationEPGCont, "status"),
			ClassName:         FvaepgClassName,
			Rn:                G(ApplicationEPGCont, "rn"),
		},

		ApplicationEPGAttributes{

			Name: G(ApplicationEPGCont, "name"),

			Annotation: G(ApplicationEPGCont, "annotation"),

			ExceptionTag: G(ApplicationEPGCont, "exceptionTag"),

			FloodOnEncap: G(ApplicationEPGCont, "floodOnEncap"),

			FwdCtrl: G(ApplicationEPGCont, "fwdCtrl"),

			HasMcastSource: G(ApplicationEPGCont, "hasMcastSource"),

			IsAttrBasedEPg: G(ApplicationEPGCont, "isAttrBasedEPg"),

			MatchT: G(ApplicationEPGCont, "matchT"),

			NameAlias: G(ApplicationEPGCont, "nameAlias"),

			PcEnfPref: G(ApplicationEPGCont, "pcEnfPref"),

			PrefGrMemb: G(ApplicationEPGCont, "prefGrMemb"),

			Prio: G(ApplicationEPGCont, "prio"),

			Shutdown: G(ApplicationEPGCont, "shutdown"),
		},
	}
}

func ApplicationEPGFromContainer(cont *container.Container) *ApplicationEPG {

	return ApplicationEPGFromContainerList(cont, 0)
}

func ApplicationEPGListFromContainer(cont *container.Container) []*ApplicationEPG {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*ApplicationEPG, length)

	for i := 0; i < length; i++ {

		arr[i] = ApplicationEPGFromContainerList(cont, i)
	}

	return arr
}

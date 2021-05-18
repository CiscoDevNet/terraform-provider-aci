package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const OspfctxpolClassName = "ospfCtxPol"

type OSPFTimersPolicy struct {
	BaseAttributes
	OSPFTimersPolicyAttributes
}

type OSPFTimersPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	BwRef string `json:",omitempty"`

	Ctrl string `json:",omitempty"`

	Dist string `json:",omitempty"`

	GrCtrl string `json:",omitempty"`

	LsaArrivalIntvl string `json:",omitempty"`

	LsaGpPacingIntvl string `json:",omitempty"`

	LsaHoldIntvl string `json:",omitempty"`

	LsaMaxIntvl string `json:",omitempty"`

	LsaStartIntvl string `json:",omitempty"`

	MaxEcmp string `json:",omitempty"`

	MaxLsaAction string `json:",omitempty"`

	MaxLsaNum string `json:",omitempty"`

	MaxLsaResetIntvl string `json:",omitempty"`

	MaxLsaSleepCnt string `json:",omitempty"`

	MaxLsaSleepIntvl string `json:",omitempty"`

	MaxLsaThresh string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	SpfHoldIntvl string `json:",omitempty"`

	SpfInitIntvl string `json:",omitempty"`

	SpfMaxIntvl string `json:",omitempty"`
}

func NewOSPFTimersPolicy(ospfCtxPolRn, parentDn, description string, ospfCtxPolattr OSPFTimersPolicyAttributes) *OSPFTimersPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, ospfCtxPolRn)
	return &OSPFTimersPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         OspfctxpolClassName,
			Rn:                ospfCtxPolRn,
		},

		OSPFTimersPolicyAttributes: ospfCtxPolattr,
	}
}

func (ospfCtxPol *OSPFTimersPolicy) ToMap() (map[string]string, error) {
	ospfCtxPolMap, err := ospfCtxPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(ospfCtxPolMap, "name", ospfCtxPol.Name)

	A(ospfCtxPolMap, "annotation", ospfCtxPol.Annotation)

	A(ospfCtxPolMap, "bwRef", ospfCtxPol.BwRef)

	A(ospfCtxPolMap, "ctrl", ospfCtxPol.Ctrl)

	A(ospfCtxPolMap, "dist", ospfCtxPol.Dist)

	A(ospfCtxPolMap, "grCtrl", ospfCtxPol.GrCtrl)

	A(ospfCtxPolMap, "lsaArrivalIntvl", ospfCtxPol.LsaArrivalIntvl)

	A(ospfCtxPolMap, "lsaGpPacingIntvl", ospfCtxPol.LsaGpPacingIntvl)

	A(ospfCtxPolMap, "lsaHoldIntvl", ospfCtxPol.LsaHoldIntvl)

	A(ospfCtxPolMap, "lsaMaxIntvl", ospfCtxPol.LsaMaxIntvl)

	A(ospfCtxPolMap, "lsaStartIntvl", ospfCtxPol.LsaStartIntvl)

	A(ospfCtxPolMap, "maxEcmp", ospfCtxPol.MaxEcmp)

	A(ospfCtxPolMap, "maxLsaAction", ospfCtxPol.MaxLsaAction)

	A(ospfCtxPolMap, "maxLsaNum", ospfCtxPol.MaxLsaNum)

	A(ospfCtxPolMap, "maxLsaResetIntvl", ospfCtxPol.MaxLsaResetIntvl)

	A(ospfCtxPolMap, "maxLsaSleepCnt", ospfCtxPol.MaxLsaSleepCnt)

	A(ospfCtxPolMap, "maxLsaSleepIntvl", ospfCtxPol.MaxLsaSleepIntvl)

	A(ospfCtxPolMap, "maxLsaThresh", ospfCtxPol.MaxLsaThresh)

	A(ospfCtxPolMap, "nameAlias", ospfCtxPol.NameAlias)

	A(ospfCtxPolMap, "spfHoldIntvl", ospfCtxPol.SpfHoldIntvl)

	A(ospfCtxPolMap, "spfInitIntvl", ospfCtxPol.SpfInitIntvl)

	A(ospfCtxPolMap, "spfMaxIntvl", ospfCtxPol.SpfMaxIntvl)

	return ospfCtxPolMap, err
}

func OSPFTimersPolicyFromContainerList(cont *container.Container, index int) *OSPFTimersPolicy {

	OSPFTimersPolicyCont := cont.S("imdata").Index(index).S(OspfctxpolClassName, "attributes")
	return &OSPFTimersPolicy{
		BaseAttributes{
			DistinguishedName: G(OSPFTimersPolicyCont, "dn"),
			Description:       G(OSPFTimersPolicyCont, "descr"),
			Status:            G(OSPFTimersPolicyCont, "status"),
			ClassName:         OspfctxpolClassName,
			Rn:                G(OSPFTimersPolicyCont, "rn"),
		},

		OSPFTimersPolicyAttributes{

			Name: G(OSPFTimersPolicyCont, "name"),

			Annotation: G(OSPFTimersPolicyCont, "annotation"),

			BwRef: G(OSPFTimersPolicyCont, "bwRef"),

			Ctrl: G(OSPFTimersPolicyCont, "ctrl"),

			Dist: G(OSPFTimersPolicyCont, "dist"),

			GrCtrl: G(OSPFTimersPolicyCont, "grCtrl"),

			LsaArrivalIntvl: G(OSPFTimersPolicyCont, "lsaArrivalIntvl"),

			LsaGpPacingIntvl: G(OSPFTimersPolicyCont, "lsaGpPacingIntvl"),

			LsaHoldIntvl: G(OSPFTimersPolicyCont, "lsaHoldIntvl"),

			LsaMaxIntvl: G(OSPFTimersPolicyCont, "lsaMaxIntvl"),

			LsaStartIntvl: G(OSPFTimersPolicyCont, "lsaStartIntvl"),

			MaxEcmp: G(OSPFTimersPolicyCont, "maxEcmp"),

			MaxLsaAction: G(OSPFTimersPolicyCont, "maxLsaAction"),

			MaxLsaNum: G(OSPFTimersPolicyCont, "maxLsaNum"),

			MaxLsaResetIntvl: G(OSPFTimersPolicyCont, "maxLsaResetIntvl"),

			MaxLsaSleepCnt: G(OSPFTimersPolicyCont, "maxLsaSleepCnt"),

			MaxLsaSleepIntvl: G(OSPFTimersPolicyCont, "maxLsaSleepIntvl"),

			MaxLsaThresh: G(OSPFTimersPolicyCont, "maxLsaThresh"),

			NameAlias: G(OSPFTimersPolicyCont, "nameAlias"),

			SpfHoldIntvl: G(OSPFTimersPolicyCont, "spfHoldIntvl"),

			SpfInitIntvl: G(OSPFTimersPolicyCont, "spfInitIntvl"),

			SpfMaxIntvl: G(OSPFTimersPolicyCont, "spfMaxIntvl"),
		},
	}
}

func OSPFTimersPolicyFromContainer(cont *container.Container) *OSPFTimersPolicy {

	return OSPFTimersPolicyFromContainerList(cont, 0)
}

func OSPFTimersPolicyListFromContainer(cont *container.Container) []*OSPFTimersPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*OSPFTimersPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = OSPFTimersPolicyFromContainerList(cont, i)
	}

	return arr
}

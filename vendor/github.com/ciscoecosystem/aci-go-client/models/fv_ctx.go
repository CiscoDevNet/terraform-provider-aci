package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvctxClassName = "fvCtx"

type VRF struct {
	BaseAttributes
	VRFAttributes
}

type VRFAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	BdEnforcedEnable string `json:",omitempty"`

	IpDataPlaneLearning string `json:",omitempty"`

	KnwMcastAct string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	PcEnfDir string `json:",omitempty"`

	PcEnfPref string `json:",omitempty"`
}

func NewVRF(fvCtxRn, parentDn, description string, fvCtxattr VRFAttributes) *VRF {
	dn := fmt.Sprintf("%s/%s", parentDn, fvCtxRn)
	return &VRF{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvctxClassName,
			Rn:                fvCtxRn,
		},

		VRFAttributes: fvCtxattr,
	}
}

func (fvCtx *VRF) ToMap() (map[string]string, error) {
	fvCtxMap, err := fvCtx.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvCtxMap, "name", fvCtx.Name)

	A(fvCtxMap, "annotation", fvCtx.Annotation)

	A(fvCtxMap, "bdEnforcedEnable", fvCtx.BdEnforcedEnable)

	A(fvCtxMap, "ipDataPlaneLearning", fvCtx.IpDataPlaneLearning)

	A(fvCtxMap, "knwMcastAct", fvCtx.KnwMcastAct)

	A(fvCtxMap, "nameAlias", fvCtx.NameAlias)

	A(fvCtxMap, "pcEnfDir", fvCtx.PcEnfDir)

	A(fvCtxMap, "pcEnfPref", fvCtx.PcEnfPref)

	return fvCtxMap, err
}

func VRFFromContainerList(cont *container.Container, index int) *VRF {

	VRFCont := cont.S("imdata").Index(index).S(FvctxClassName, "attributes")
	return &VRF{
		BaseAttributes{
			DistinguishedName: G(VRFCont, "dn"),
			Description:       G(VRFCont, "descr"),
			Status:            G(VRFCont, "status"),
			ClassName:         FvctxClassName,
			Rn:                G(VRFCont, "rn"),
		},

		VRFAttributes{

			Name: G(VRFCont, "name"),

			Annotation: G(VRFCont, "annotation"),

			BdEnforcedEnable: G(VRFCont, "bdEnforcedEnable"),

			IpDataPlaneLearning: G(VRFCont, "ipDataPlaneLearning"),

			KnwMcastAct: G(VRFCont, "knwMcastAct"),

			NameAlias: G(VRFCont, "nameAlias"),

			PcEnfDir: G(VRFCont, "pcEnfDir"),

			PcEnfPref: G(VRFCont, "pcEnfPref"),
		},
	}
}

func VRFFromContainer(cont *container.Container) *VRF {

	return VRFFromContainerList(cont, 0)
}

func VRFListFromContainer(cont *container.Container) []*VRF {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*VRF, length)

	for i := 0; i < length; i++ {

		arr[i] = VRFFromContainerList(cont, i)
	}

	return arr
}

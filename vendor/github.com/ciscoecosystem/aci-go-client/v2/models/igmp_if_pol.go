package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnIgmpIfPol        = "igmpIfPol-%s"
	DnIgmpIfPol        = "uni/tn-%s/igmpIfPol-%s"
	ParentDnIgmpIfPol  = "uni/tn-%s"
	IgmpIfPolClassName = "igmpIfPol"
)

type IGMPInterfacePolicy struct {
	BaseAttributes
	IGMPInterfacePolicyAttributes
}

type IGMPInterfacePolicyAttributes struct {
	Annotation      string `json:",omitempty"`
	GrpTimeout      string `json:",omitempty"`
	IfCtrl          string `json:",omitempty"`
	LastMbrCnt      string `json:",omitempty"`
	LastMbrRespTime string `json:",omitempty"`
	Name            string `json:",omitempty"`
	NameAlias       string `json:",omitempty"`
	QuerierTimeout  string `json:",omitempty"`
	QueryIntvl      string `json:",omitempty"`
	RobustFac       string `json:",omitempty"`
	RspIntvl        string `json:",omitempty"`
	StartQueryCnt   string `json:",omitempty"`
	StartQueryIntvl string `json:",omitempty"`
	Ver             string `json:",omitempty"`
}

func NewIGMPInterfacePolicy(igmpIfPolRn, parentDn, description string, igmpIfPolAttr IGMPInterfacePolicyAttributes) *IGMPInterfacePolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, igmpIfPolRn)
	return &IGMPInterfacePolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         IgmpIfPolClassName,
			Rn:                igmpIfPolRn,
		},
		IGMPInterfacePolicyAttributes: igmpIfPolAttr,
	}
}

func (igmpIfPol *IGMPInterfacePolicy) ToMap() (map[string]string, error) {
	igmpIfPolMap, err := igmpIfPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(igmpIfPolMap, "annotation", igmpIfPol.Annotation)
	A(igmpIfPolMap, "grpTimeout", igmpIfPol.GrpTimeout)
	A(igmpIfPolMap, "ifCtrl", igmpIfPol.IfCtrl)
	A(igmpIfPolMap, "lastMbrCnt", igmpIfPol.LastMbrCnt)
	A(igmpIfPolMap, "lastMbrRespTime", igmpIfPol.LastMbrRespTime)
	A(igmpIfPolMap, "name", igmpIfPol.Name)
	A(igmpIfPolMap, "nameAlias", igmpIfPol.NameAlias)
	A(igmpIfPolMap, "querierTimeout", igmpIfPol.QuerierTimeout)
	A(igmpIfPolMap, "queryIntvl", igmpIfPol.QueryIntvl)
	A(igmpIfPolMap, "robustFac", igmpIfPol.RobustFac)
	A(igmpIfPolMap, "rspIntvl", igmpIfPol.RspIntvl)
	A(igmpIfPolMap, "startQueryCnt", igmpIfPol.StartQueryCnt)
	A(igmpIfPolMap, "startQueryIntvl", igmpIfPol.StartQueryIntvl)
	A(igmpIfPolMap, "ver", igmpIfPol.Ver)
	return igmpIfPolMap, err
}

func IGMPInterfacePolicyFromContainerList(cont *container.Container, index int) *IGMPInterfacePolicy {
	IGMPInterfacePolicyCont := cont.S("imdata").Index(index).S(IgmpIfPolClassName, "attributes")
	return &IGMPInterfacePolicy{
		BaseAttributes{
			DistinguishedName: G(IGMPInterfacePolicyCont, "dn"),
			Description:       G(IGMPInterfacePolicyCont, "descr"),
			Status:            G(IGMPInterfacePolicyCont, "status"),
			ClassName:         IgmpIfPolClassName,
			Rn:                G(IGMPInterfacePolicyCont, "rn"),
		},
		IGMPInterfacePolicyAttributes{
			Annotation:      G(IGMPInterfacePolicyCont, "annotation"),
			GrpTimeout:      G(IGMPInterfacePolicyCont, "grpTimeout"),
			IfCtrl:          G(IGMPInterfacePolicyCont, "ifCtrl"),
			LastMbrCnt:      G(IGMPInterfacePolicyCont, "lastMbrCnt"),
			LastMbrRespTime: G(IGMPInterfacePolicyCont, "lastMbrRespTime"),
			Name:            G(IGMPInterfacePolicyCont, "name"),
			NameAlias:       G(IGMPInterfacePolicyCont, "nameAlias"),
			QuerierTimeout:  G(IGMPInterfacePolicyCont, "querierTimeout"),
			QueryIntvl:      G(IGMPInterfacePolicyCont, "queryIntvl"),
			RobustFac:       G(IGMPInterfacePolicyCont, "robustFac"),
			RspIntvl:        G(IGMPInterfacePolicyCont, "rspIntvl"),
			StartQueryCnt:   G(IGMPInterfacePolicyCont, "startQueryCnt"),
			StartQueryIntvl: G(IGMPInterfacePolicyCont, "startQueryIntvl"),
			Ver:             G(IGMPInterfacePolicyCont, "ver"),
		},
	}
}

func IGMPInterfacePolicyFromContainer(cont *container.Container) *IGMPInterfacePolicy {
	return IGMPInterfacePolicyFromContainerList(cont, 0)
}

func IGMPInterfacePolicyListFromContainer(cont *container.Container) []*IGMPInterfacePolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*IGMPInterfacePolicy, length)

	for i := 0; i < length; i++ {
		arr[i] = IGMPInterfacePolicyFromContainerList(cont, i)
	}

	return arr
}

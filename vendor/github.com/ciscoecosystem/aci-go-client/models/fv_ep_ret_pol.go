package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FvepretpolClassName = "fvEpRetPol"

type EndPointRetentionPolicy struct {
	BaseAttributes
	EndPointRetentionPolicyAttributes
}

type EndPointRetentionPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	BounceAgeIntvl string `json:",omitempty"`

	BounceTrig string `json:",omitempty"`

	HoldIntvl string `json:",omitempty"`

	LocalEpAgeIntvl string `json:",omitempty"`

	MoveFreq string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	RemoteEpAgeIntvl string `json:",omitempty"`
}

func NewEndPointRetentionPolicy(fvEpRetPolRn, parentDn, description string, fvEpRetPolattr EndPointRetentionPolicyAttributes) *EndPointRetentionPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, fvEpRetPolRn)
	return &EndPointRetentionPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FvepretpolClassName,
			Rn:                fvEpRetPolRn,
		},

		EndPointRetentionPolicyAttributes: fvEpRetPolattr,
	}
}

func (fvEpRetPol *EndPointRetentionPolicy) ToMap() (map[string]string, error) {
	fvEpRetPolMap, err := fvEpRetPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(fvEpRetPolMap, "name", fvEpRetPol.Name)

	A(fvEpRetPolMap, "annotation", fvEpRetPol.Annotation)

	A(fvEpRetPolMap, "bounceAgeIntvl", fvEpRetPol.BounceAgeIntvl)

	A(fvEpRetPolMap, "bounceTrig", fvEpRetPol.BounceTrig)

	A(fvEpRetPolMap, "holdIntvl", fvEpRetPol.HoldIntvl)

	A(fvEpRetPolMap, "localEpAgeIntvl", fvEpRetPol.LocalEpAgeIntvl)

	A(fvEpRetPolMap, "moveFreq", fvEpRetPol.MoveFreq)

	A(fvEpRetPolMap, "nameAlias", fvEpRetPol.NameAlias)

	A(fvEpRetPolMap, "remoteEpAgeIntvl", fvEpRetPol.RemoteEpAgeIntvl)

	return fvEpRetPolMap, err
}

func EndPointRetentionPolicyFromContainerList(cont *container.Container, index int) *EndPointRetentionPolicy {

	EndPointRetentionPolicyCont := cont.S("imdata").Index(index).S(FvepretpolClassName, "attributes")
	return &EndPointRetentionPolicy{
		BaseAttributes{
			DistinguishedName: G(EndPointRetentionPolicyCont, "dn"),
			Description:       G(EndPointRetentionPolicyCont, "descr"),
			Status:            G(EndPointRetentionPolicyCont, "status"),
			ClassName:         FvepretpolClassName,
			Rn:                G(EndPointRetentionPolicyCont, "rn"),
		},

		EndPointRetentionPolicyAttributes{

			Name: G(EndPointRetentionPolicyCont, "name"),

			Annotation: G(EndPointRetentionPolicyCont, "annotation"),

			BounceAgeIntvl: G(EndPointRetentionPolicyCont, "bounceAgeIntvl"),

			BounceTrig: G(EndPointRetentionPolicyCont, "bounceTrig"),

			HoldIntvl: G(EndPointRetentionPolicyCont, "holdIntvl"),

			LocalEpAgeIntvl: G(EndPointRetentionPolicyCont, "localEpAgeIntvl"),

			MoveFreq: G(EndPointRetentionPolicyCont, "moveFreq"),

			NameAlias: G(EndPointRetentionPolicyCont, "nameAlias"),

			RemoteEpAgeIntvl: G(EndPointRetentionPolicyCont, "remoteEpAgeIntvl"),
		},
	}
}

func EndPointRetentionPolicyFromContainer(cont *container.Container) *EndPointRetentionPolicy {

	return EndPointRetentionPolicyFromContainerList(cont, 0)
}

func EndPointRetentionPolicyListFromContainer(cont *container.Container) []*EndPointRetentionPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*EndPointRetentionPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = EndPointRetentionPolicyFromContainerList(cont, i)
	}

	return arr
}

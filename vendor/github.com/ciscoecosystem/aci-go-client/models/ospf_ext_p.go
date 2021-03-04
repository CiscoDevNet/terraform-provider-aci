package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const OspfextpClassName = "ospfExtP"

type L3outOspfExternalPolicy struct {
	BaseAttributes
	L3outOspfExternalPolicyAttributes
}

type L3outOspfExternalPolicyAttributes struct {
	Annotation string `json:",omitempty"`

	AreaCost string `json:",omitempty"`

	AreaCtrl string `json:",omitempty"`

	AreaId string `json:",omitempty"`

	AreaType string `json:",omitempty"`

	MultipodInternal string `json:",omitempty"`

	NameAlias string `json:",omitempty"`
}

func NewL3outOspfExternalPolicy(ospfExtPRn, parentDn, description string, ospfExtPattr L3outOspfExternalPolicyAttributes) *L3outOspfExternalPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, ospfExtPRn)
	return &L3outOspfExternalPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         OspfextpClassName,
			Rn:                ospfExtPRn,
		},

		L3outOspfExternalPolicyAttributes: ospfExtPattr,
	}
}

func (ospfExtP *L3outOspfExternalPolicy) ToMap() (map[string]string, error) {
	ospfExtPMap, err := ospfExtP.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(ospfExtPMap, "annotation", ospfExtP.Annotation)

	A(ospfExtPMap, "areaCost", ospfExtP.AreaCost)

	A(ospfExtPMap, "areaCtrl", ospfExtP.AreaCtrl)

	A(ospfExtPMap, "areaId", ospfExtP.AreaId)

	A(ospfExtPMap, "areaType", ospfExtP.AreaType)

	A(ospfExtPMap, "multipodInternal", ospfExtP.MultipodInternal)

	A(ospfExtPMap, "nameAlias", ospfExtP.NameAlias)

	return ospfExtPMap, err
}

func L3outOspfExternalPolicyFromContainerList(cont *container.Container, index int) *L3outOspfExternalPolicy {

	L3outOspfExternalPolicyCont := cont.S("imdata").Index(index).S(OspfextpClassName, "attributes")
	return &L3outOspfExternalPolicy{
		BaseAttributes{
			DistinguishedName: G(L3outOspfExternalPolicyCont, "dn"),
			Description:       G(L3outOspfExternalPolicyCont, "descr"),
			Status:            G(L3outOspfExternalPolicyCont, "status"),
			ClassName:         OspfextpClassName,
			Rn:                G(L3outOspfExternalPolicyCont, "rn"),
		},

		L3outOspfExternalPolicyAttributes{

			Annotation: G(L3outOspfExternalPolicyCont, "annotation"),

			AreaCost: G(L3outOspfExternalPolicyCont, "areaCost"),

			AreaCtrl: G(L3outOspfExternalPolicyCont, "areaCtrl"),

			AreaId: G(L3outOspfExternalPolicyCont, "areaId"),

			AreaType: G(L3outOspfExternalPolicyCont, "areaType"),

			MultipodInternal: G(L3outOspfExternalPolicyCont, "multipodInternal"),

			NameAlias: G(L3outOspfExternalPolicyCont, "nameAlias"),
		},
	}
}

func L3outOspfExternalPolicyFromContainer(cont *container.Container) *L3outOspfExternalPolicy {

	return L3outOspfExternalPolicyFromContainerList(cont, 0)
}

func L3outOspfExternalPolicyListFromContainer(cont *container.Container) []*L3outOspfExternalPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*L3outOspfExternalPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = L3outOspfExternalPolicyFromContainerList(cont, i)
	}

	return arr
}

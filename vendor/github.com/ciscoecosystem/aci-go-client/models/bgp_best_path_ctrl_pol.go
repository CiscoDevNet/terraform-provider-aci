package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const BgpbestpathctrlpolClassName = "bgpBestPathCtrlPol"

type BgpBestPathPolicy struct {
	BaseAttributes
	BgpBestPathPolicyAttributes
}

type BgpBestPathPolicyAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	Ctrl string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Userdom string `json:",omitempty"`
}

func NewBgpBestPathPolicy(bgpBestPathCtrlPolRn, parentDn, description string, bgpBestPathCtrlPolattr BgpBestPathPolicyAttributes) *BgpBestPathPolicy {
	dn := fmt.Sprintf("%s/%s", parentDn, bgpBestPathCtrlPolRn)
	return &BgpBestPathPolicy{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         BgpbestpathctrlpolClassName,
			Rn:                bgpBestPathCtrlPolRn,
		},

		BgpBestPathPolicyAttributes: bgpBestPathCtrlPolattr,
	}
}

func (bgpBestPathCtrlPol *BgpBestPathPolicy) ToMap() (map[string]string, error) {
	bgpBestPathCtrlPolMap, err := bgpBestPathCtrlPol.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(bgpBestPathCtrlPolMap, "name", bgpBestPathCtrlPol.Name)

	A(bgpBestPathCtrlPolMap, "annotation", bgpBestPathCtrlPol.Annotation)

	A(bgpBestPathCtrlPolMap, "ctrl", bgpBestPathCtrlPol.Ctrl)

	A(bgpBestPathCtrlPolMap, "nameAlias", bgpBestPathCtrlPol.NameAlias)

	A(bgpBestPathCtrlPolMap, "userdom", bgpBestPathCtrlPol.Userdom)

	return bgpBestPathCtrlPolMap, err
}

func BgpBestPathPolicyFromContainerList(cont *container.Container, index int) *BgpBestPathPolicy {

	BgpBestPathPolicyCont := cont.S("imdata").Index(index).S(BgpbestpathctrlpolClassName, "attributes")
	return &BgpBestPathPolicy{
		BaseAttributes{
			DistinguishedName: G(BgpBestPathPolicyCont, "dn"),
			Description:       G(BgpBestPathPolicyCont, "descr"),
			Status:            G(BgpBestPathPolicyCont, "status"),
			ClassName:         BgpbestpathctrlpolClassName,
			Rn:                G(BgpBestPathPolicyCont, "rn"),
		},

		BgpBestPathPolicyAttributes{

			Name: G(BgpBestPathPolicyCont, "name"),

			Annotation: G(BgpBestPathPolicyCont, "annotation"),

			Ctrl: G(BgpBestPathPolicyCont, "ctrl"),

			NameAlias: G(BgpBestPathPolicyCont, "nameAlias"),

			Userdom: G(BgpBestPathPolicyCont, "userdom"),
		},
	}
}

func BgpBestPathPolicyFromContainer(cont *container.Container) *BgpBestPathPolicy {

	return BgpBestPathPolicyFromContainerList(cont, 0)
}

func BgpBestPathPolicyListFromContainer(cont *container.Container) []*BgpBestPathPolicy {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*BgpBestPathPolicy, length)

	for i := 0; i < length; i++ {

		arr[i] = BgpBestPathPolicyFromContainerList(cont, i)
	}

	return arr
}

package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const CloudepgClassName = "cloudEPg"

type CloudEPg struct {
	BaseAttributes
	CloudEPgAttributes
}

type CloudEPgAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	ExceptionTag string `json:",omitempty"`

	FloodOnEncap string `json:",omitempty"`

	MatchT string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	PrefGrMemb string `json:",omitempty"`

	Prio string `json:",omitempty"`
}

func NewCloudEPg(cloudEPgRn, parentDn, description string, cloudEPgattr CloudEPgAttributes) *CloudEPg {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudEPgRn)
	return &CloudEPg{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudepgClassName,
			Rn:                cloudEPgRn,
		},

		CloudEPgAttributes: cloudEPgattr,
	}
}

func (cloudEPg *CloudEPg) ToMap() (map[string]string, error) {
	cloudEPgMap, err := cloudEPg.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudEPgMap, "name", cloudEPg.Name)

	A(cloudEPgMap, "annotation", cloudEPg.Annotation)

	A(cloudEPgMap, "exceptionTag", cloudEPg.ExceptionTag)

	A(cloudEPgMap, "floodOnEncap", cloudEPg.FloodOnEncap)

	A(cloudEPgMap, "matchT", cloudEPg.MatchT)

	A(cloudEPgMap, "nameAlias", cloudEPg.NameAlias)

	A(cloudEPgMap, "prefGrMemb", cloudEPg.PrefGrMemb)

	A(cloudEPgMap, "prio", cloudEPg.Prio)

	return cloudEPgMap, err
}

func CloudEPgFromContainerList(cont *container.Container, index int) *CloudEPg {

	CloudEPgCont := cont.S("imdata").Index(index).S(CloudepgClassName, "attributes")
	return &CloudEPg{
		BaseAttributes{
			DistinguishedName: G(CloudEPgCont, "dn"),
			Description:       G(CloudEPgCont, "descr"),
			Status:            G(CloudEPgCont, "status"),
			ClassName:         CloudepgClassName,
			Rn:                G(CloudEPgCont, "rn"),
		},

		CloudEPgAttributes{

			Name: G(CloudEPgCont, "name"),

			Annotation: G(CloudEPgCont, "annotation"),

			ExceptionTag: G(CloudEPgCont, "exceptionTag"),

			FloodOnEncap: G(CloudEPgCont, "floodOnEncap"),

			MatchT: G(CloudEPgCont, "matchT"),

			NameAlias: G(CloudEPgCont, "nameAlias"),

			PrefGrMemb: G(CloudEPgCont, "prefGrMemb"),

			Prio: G(CloudEPgCont, "prio"),
		},
	}
}

func CloudEPgFromContainer(cont *container.Container) *CloudEPg {

	return CloudEPgFromContainerList(cont, 0)
}

func CloudEPgListFromContainer(cont *container.Container) []*CloudEPg {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*CloudEPg, length)

	for i := 0; i < length; i++ {

		arr[i] = CloudEPgFromContainerList(cont, i)
	}

	return arr
}

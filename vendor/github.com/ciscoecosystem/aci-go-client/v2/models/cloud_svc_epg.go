package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	RnCloudSvcEPg        = "cloudsvcepg-%s"
	DnCloudSvcEPg        = "uni/tn-%s/cloudapp-%s/cloudsvcepg-%s"
	ParentDnCloudSvcEPg  = "uni/tn-%s/cloudapp-%s"
	CloudSvcEPgClassName = "cloudSvcEPg"
)

type CloudServiceEPg struct {
	BaseAttributes
	CloudServiceEPgAttributes
}

type CloudServiceEPgAttributes struct {
	AccessType           string `json:",omitempty"`
	Annotation           string `json:",omitempty"`
	AzPrivateEndpoint    string `json:",omitempty"`
	CustomSvcType        string `json:",omitempty"`
	DeploymentType       string `json:",omitempty"`
	ExceptionTag         string `json:",omitempty"`
	FloodOnEncap         string `json:",omitempty"`
	MatchT               string `json:",omitempty"`
	Name                 string `json:",omitempty"`
	NameAlias            string `json:",omitempty"`
	PrefGrMemb           string `json:",omitempty"`
	Prio                 string `json:",omitempty"`
	CloudServiceEPg_type string `json:",omitempty"`
}

func NewCloudServiceEPg(cloudSvcEPgRn, parentDn, description string, cloudSvcEPgAttr CloudServiceEPgAttributes) *CloudServiceEPg {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudSvcEPgRn)
	return &CloudServiceEPg{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         CloudSvcEPgClassName,
			Rn:                cloudSvcEPgRn,
		},
		CloudServiceEPgAttributes: cloudSvcEPgAttr,
	}
}

func (cloudSvcEPg *CloudServiceEPg) ToMap() (map[string]string, error) {
	cloudSvcEPgMap, err := cloudSvcEPg.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(cloudSvcEPgMap, "accessType", cloudSvcEPg.AccessType)
	A(cloudSvcEPgMap, "annotation", cloudSvcEPg.Annotation)
	A(cloudSvcEPgMap, "azPrivateEndpoint", cloudSvcEPg.AzPrivateEndpoint)
	A(cloudSvcEPgMap, "customSvcType", cloudSvcEPg.CustomSvcType)
	A(cloudSvcEPgMap, "deploymentType", cloudSvcEPg.DeploymentType)
	A(cloudSvcEPgMap, "exceptionTag", cloudSvcEPg.ExceptionTag)
	A(cloudSvcEPgMap, "floodOnEncap", cloudSvcEPg.FloodOnEncap)
	A(cloudSvcEPgMap, "matchT", cloudSvcEPg.MatchT)
	A(cloudSvcEPgMap, "name", cloudSvcEPg.Name)
	A(cloudSvcEPgMap, "nameAlias", cloudSvcEPg.NameAlias)
	A(cloudSvcEPgMap, "prefGrMemb", cloudSvcEPg.PrefGrMemb)
	A(cloudSvcEPgMap, "prio", cloudSvcEPg.Prio)
	A(cloudSvcEPgMap, "type", cloudSvcEPg.CloudServiceEPg_type)
	return cloudSvcEPgMap, err
}

func CloudServiceEPgFromContainerList(cont *container.Container, index int) *CloudServiceEPg {
	CloudServiceEPgCont := cont.S("imdata").Index(index).S(CloudSvcEPgClassName, "attributes")
	return &CloudServiceEPg{
		BaseAttributes{
			DistinguishedName: G(CloudServiceEPgCont, "dn"),
			Description:       G(CloudServiceEPgCont, "descr"),
			Status:            G(CloudServiceEPgCont, "status"),
			ClassName:         CloudSvcEPgClassName,
			Rn:                G(CloudServiceEPgCont, "rn"),
		},
		CloudServiceEPgAttributes{
			AccessType:           G(CloudServiceEPgCont, "accessType"),
			Annotation:           G(CloudServiceEPgCont, "annotation"),
			AzPrivateEndpoint:    G(CloudServiceEPgCont, "azPrivateEndpoint"),
			CustomSvcType:        G(CloudServiceEPgCont, "customSvcType"),
			DeploymentType:       G(CloudServiceEPgCont, "deploymentType"),
			ExceptionTag:         G(CloudServiceEPgCont, "exceptionTag"),
			FloodOnEncap:         G(CloudServiceEPgCont, "floodOnEncap"),
			MatchT:               G(CloudServiceEPgCont, "matchT"),
			Name:                 G(CloudServiceEPgCont, "name"),
			NameAlias:            G(CloudServiceEPgCont, "nameAlias"),
			PrefGrMemb:           G(CloudServiceEPgCont, "prefGrMemb"),
			Prio:                 G(CloudServiceEPgCont, "prio"),
			CloudServiceEPg_type: G(CloudServiceEPgCont, "type"),
		},
	}
}

func CloudServiceEPgFromContainer(cont *container.Container) *CloudServiceEPg {
	return CloudServiceEPgFromContainerList(cont, 0)
}

func CloudServiceEPgListFromContainer(cont *container.Container) []*CloudServiceEPg {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudServiceEPg, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudServiceEPgFromContainerList(cont, i)
	}

	return arr
}

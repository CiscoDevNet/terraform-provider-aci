package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

const (
	DncloudAD        = "uni/tn-%s/ad-%s"
	RncloudAD        = "ad-%s"
	ParentDncloudAD  = "uni/tn-%s"
	CloudadClassName = "cloudAD"
)

type CloudActiveDirectory struct {
	BaseAttributes
	NameAliasAttribute
	CloudActiveDirectoryAttributes
}

type CloudActiveDirectoryAttributes struct {
	Annotation         string `json:",omitempty"`
	ActiveDirectory_id string `json:",omitempty"`
	Name               string `json:",omitempty"`
}

func NewCloudActiveDirectory(cloudADRn, parentDn, nameAlias string, cloudADAttr CloudActiveDirectoryAttributes) *CloudActiveDirectory {
	dn := fmt.Sprintf("%s/%s", parentDn, cloudADRn)
	return &CloudActiveDirectory{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Status:            "created, modified",
			ClassName:         CloudadClassName,
			Rn:                cloudADRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		CloudActiveDirectoryAttributes: cloudADAttr,
	}
}

func (cloudAD *CloudActiveDirectory) ToMap() (map[string]string, error) {
	cloudADMap, err := cloudAD.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	alias, err := cloudAD.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}

	for key, value := range alias {
		A(cloudADMap, key, value)
	}

	A(cloudADMap, "id", cloudAD.ActiveDirectory_id)
	A(cloudADMap, "name", cloudAD.Name)
	return cloudADMap, err
}

func CloudActiveDirectoryFromContainerList(cont *container.Container, index int) *CloudActiveDirectory {
	CloudActiveDirectoryCont := cont.S("imdata").Index(index).S(CloudadClassName, "attributes")
	return &CloudActiveDirectory{
		BaseAttributes{
			DistinguishedName: G(CloudActiveDirectoryCont, "dn"),
			Status:            G(CloudActiveDirectoryCont, "status"),
			ClassName:         CloudadClassName,
			Rn:                G(CloudActiveDirectoryCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(CloudActiveDirectoryCont, "nameAlias"),
		},
		CloudActiveDirectoryAttributes{
			ActiveDirectory_id: G(CloudActiveDirectoryCont, "id"),
			Name:               G(CloudActiveDirectoryCont, "name"),
		},
	}
}

func CloudActiveDirectoryFromContainer(cont *container.Container) *CloudActiveDirectory {
	return CloudActiveDirectoryFromContainerList(cont, 0)
}

func CloudActiveDirectoryListFromContainer(cont *container.Container) []*CloudActiveDirectory {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*CloudActiveDirectory, length)

	for i := 0; i < length; i++ {
		arr[i] = CloudActiveDirectoryFromContainerList(cont, i)
	}

	return arr
}

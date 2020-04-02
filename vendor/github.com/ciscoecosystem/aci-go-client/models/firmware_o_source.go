package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FirmwareosourceClassName = "firmwareOSource"

type FirmwareDownloadTask struct {
	BaseAttributes
	FirmwareDownloadTaskAttributes
}

type FirmwareDownloadTaskAttributes struct {
	Name string `json:",omitempty"`

	Annotation string `json:",omitempty"`

	AuthPass string `json:",omitempty"`

	AuthType string `json:",omitempty"`

	DnldTaskFlip string `json:",omitempty"`

	IdentityPrivateKeyContents string `json:",omitempty"`

	IdentityPrivateKeyPassphrase string `json:",omitempty"`

	IdentityPublicKeyContents string `json:",omitempty"`

	LoadCatalogIfExistsAndNewer string `json:",omitempty"`

	NameAlias string `json:",omitempty"`

	Password string `json:",omitempty"`

	PollingInterval string `json:",omitempty"`

	Proto string `json:",omitempty"`

	Url string `json:",omitempty"`

	User string `json:",omitempty"`
}

func NewFirmwareDownloadTask(firmwareOSourceRn, parentDn, description string, firmwareOSourceattr FirmwareDownloadTaskAttributes) *FirmwareDownloadTask {
	dn := fmt.Sprintf("%s/%s", parentDn, firmwareOSourceRn)
	return &FirmwareDownloadTask{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FirmwareosourceClassName,
			Rn:                firmwareOSourceRn,
		},

		FirmwareDownloadTaskAttributes: firmwareOSourceattr,
	}
}

func (firmwareOSource *FirmwareDownloadTask) ToMap() (map[string]string, error) {
	firmwareOSourceMap, err := firmwareOSource.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(firmwareOSourceMap, "name", firmwareOSource.Name)

	A(firmwareOSourceMap, "annotation", firmwareOSource.Annotation)

	A(firmwareOSourceMap, "authPass", firmwareOSource.AuthPass)

	A(firmwareOSourceMap, "authType", firmwareOSource.AuthType)

	A(firmwareOSourceMap, "dnldTaskFlip", firmwareOSource.DnldTaskFlip)

	A(firmwareOSourceMap, "identityPrivateKeyContents", firmwareOSource.IdentityPrivateKeyContents)

	A(firmwareOSourceMap, "identityPrivateKeyPassphrase", firmwareOSource.IdentityPrivateKeyPassphrase)

	A(firmwareOSourceMap, "identityPublicKeyContents", firmwareOSource.IdentityPublicKeyContents)

	A(firmwareOSourceMap, "loadCatalogIfExistsAndNewer", firmwareOSource.LoadCatalogIfExistsAndNewer)

	A(firmwareOSourceMap, "nameAlias", firmwareOSource.NameAlias)

	A(firmwareOSourceMap, "password", firmwareOSource.Password)

	A(firmwareOSourceMap, "pollingInterval", firmwareOSource.PollingInterval)

	A(firmwareOSourceMap, "proto", firmwareOSource.Proto)

	A(firmwareOSourceMap, "url", firmwareOSource.Url)

	A(firmwareOSourceMap, "user", firmwareOSource.User)

	return firmwareOSourceMap, err
}

func FirmwareDownloadTaskFromContainerList(cont *container.Container, index int) *FirmwareDownloadTask {

	FirmwareDownloadTaskCont := cont.S("imdata").Index(index).S(FirmwareosourceClassName, "attributes")
	return &FirmwareDownloadTask{
		BaseAttributes{
			DistinguishedName: G(FirmwareDownloadTaskCont, "dn"),
			Description:       G(FirmwareDownloadTaskCont, "descr"),
			Status:            G(FirmwareDownloadTaskCont, "status"),
			ClassName:         FirmwareosourceClassName,
			Rn:                G(FirmwareDownloadTaskCont, "rn"),
		},

		FirmwareDownloadTaskAttributes{

			Name: G(FirmwareDownloadTaskCont, "name"),

			Annotation: G(FirmwareDownloadTaskCont, "annotation"),

			AuthPass: G(FirmwareDownloadTaskCont, "authPass"),

			AuthType: G(FirmwareDownloadTaskCont, "authType"),

			DnldTaskFlip: G(FirmwareDownloadTaskCont, "dnldTaskFlip"),

			IdentityPrivateKeyContents: G(FirmwareDownloadTaskCont, "identityPrivateKeyContents"),

			IdentityPrivateKeyPassphrase: G(FirmwareDownloadTaskCont, "identityPrivateKeyPassphrase"),

			IdentityPublicKeyContents: G(FirmwareDownloadTaskCont, "identityPublicKeyContents"),

			LoadCatalogIfExistsAndNewer: G(FirmwareDownloadTaskCont, "loadCatalogIfExistsAndNewer"),

			NameAlias: G(FirmwareDownloadTaskCont, "nameAlias"),

			Password: G(FirmwareDownloadTaskCont, "password"),

			PollingInterval: G(FirmwareDownloadTaskCont, "pollingInterval"),

			Proto: G(FirmwareDownloadTaskCont, "proto"),

			Url: G(FirmwareDownloadTaskCont, "url"),

			User: G(FirmwareDownloadTaskCont, "user"),
		},
	}
}

func FirmwareDownloadTaskFromContainer(cont *container.Container) *FirmwareDownloadTask {

	return FirmwareDownloadTaskFromContainerList(cont, 0)
}

func FirmwareDownloadTaskListFromContainer(cont *container.Container) []*FirmwareDownloadTask {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*FirmwareDownloadTask, length)

	for i := 0; i < length; i++ {

		arr[i] = FirmwareDownloadTaskFromContainerList(cont, i)
	}

	return arr
}

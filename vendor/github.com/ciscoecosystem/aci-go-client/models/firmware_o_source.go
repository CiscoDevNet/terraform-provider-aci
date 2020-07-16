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

// START: Variable/Struct/Fuction Naming per ACI SDK Model Definitions
const FirmwareOSourceClassName = "firmwareOSource"

type OSource struct {
	BaseAttributes
	OSourceAttributes
}

type OSourceAttributes struct {
	Name                         string `json:",omitempty"`
	Annotation                   string `json:",omitempty"`
	Url                          string `json:",omitempty"`
	Proto                        string `json:",omitempty"`
	User                         string `json:",omitempty"`
	AuthType                     string `json:",omitempty"`
	AuthPass                     string `json:",omitempty"`
	Password                     string `json:",omitempty"`
	IdentityPrivateKeyContents   string `json:",omitempty"`
	IdentityPrivateKeyPassphrase string `json:",omitempty"`
	IdentityPublicKeyContents    string `json:",omitempty"`
	DnldTaskFlip                 string `json:",omitempty"`
}

func NewOSource(firmwareOSourceRn, parentDn, description string, firmwareOSourceAttr OSourceAttributes) *OSource {
	dn := fmt.Sprintf("%s/%s", parentDn, firmwareOSourceRn)
	return &OSource{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "",
			ClassName:         FirmwareOSourceClassName,
			Rn:                firmwareOSourceRn,
		},
		OSourceAttributes: firmwareOSourceAttr,
	}
}

func (firmwareOSource *OSource) ToMap() (map[string]string, error) {
	firmwareOSourceMap, err := firmwareOSource.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(firmwareOSourceMap, "name", firmwareOSource.Name)
	A(firmwareOSourceMap, "annotation", firmwareOSource.Annotation)
	A(firmwareOSourceMap, "url", firmwareOSource.Url)
	A(firmwareOSourceMap, "proto", firmwareOSource.Proto)
	A(firmwareOSourceMap, "user", firmwareOSource.User)
	A(firmwareOSourceMap, "authType", firmwareOSource.AuthType)
	A(firmwareOSourceMap, "authPass", firmwareOSource.AuthPass)
	A(firmwareOSourceMap, "password", firmwareOSource.Password)
	A(firmwareOSourceMap, "identityPrivateKeyContents", firmwareOSource.IdentityPrivateKeyContents)
	A(firmwareOSourceMap, "identityPrivateKeyPassphrase", firmwareOSource.IdentityPrivateKeyPassphrase)
	A(firmwareOSourceMap, "identityPublicKeyContents", firmwareOSource.IdentityPublicKeyContents)
	A(firmwareOSourceMap, "dnldTaskFlip", firmwareOSource.DnldTaskFlip)

	return firmwareOSourceMap, err
}

func OSourceFromContainerList(cont *container.Container, index int) *OSource {

	OSourceCont := cont.S("imdata").Index(index).S(FirmwareOSourceClassName, "attributes")
	return &OSource{
		BaseAttributes{
			DistinguishedName: G(OSourceCont, "dn"),
			Description:       G(OSourceCont, "descr"),
			Status:            G(OSourceCont, "status"),
			ClassName:         FirmwareOSourceClassName,
			Rn:                G(OSourceCont, "rn"),
		},

		OSourceAttributes{
			Name:                         G(OSourceCont, "name"),
			Annotation:                   G(OSourceCont, "annotation"),
			Url:                          G(OSourceCont, "url"),
			Proto:                        G(OSourceCont, "proto"),
			User:                         G(OSourceCont, "user"),
			AuthType:                     G(OSourceCont, "authType"),
			AuthPass:                     G(OSourceCont, "AuthPass"),
			Password:                     G(OSourceCont, "password"),
			IdentityPrivateKeyContents:   G(OSourceCont, "identityPrivateKeyContents"),
			IdentityPrivateKeyPassphrase: G(OSourceCont, "identityPrivateKeyPassphrase"),
			IdentityPublicKeyContents:    G(OSourceCont, "identityPublicKeyContents"),
			DnldTaskFlip:                 G(OSourceCont, "dnldTaskFlip"),
		},
	}
}

func OSourceFromContainer(cont *container.Container) *OSource {

	return OSourceFromContainerList(cont, 0)
}

func OSourceListFromContainer(cont *container.Container) []*OSource {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*OSource, length)

	for i := 0; i < length; i++ {

		arr[i] = OSourceFromContainerList(cont, i)
	}

	return arr
}

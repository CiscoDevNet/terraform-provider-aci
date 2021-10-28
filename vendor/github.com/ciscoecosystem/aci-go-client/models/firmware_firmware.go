package models

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/container"
	"strconv"
)

const FirmwareFirmwareClassName = "firmwareFirmware"

type Firmware struct {
	BaseAttributes
	FirmwareAttributes
}

type FirmwareAttributes struct {
	Name            string `json:",omitempty"`
	Annotation      string `json:",omitempty"`
	FullVersion     string `json:",omitempty"`
	Version         string `json:",omitempty"`
	MinorVersion    string `json:",omitempty"`
	Isoname         string `json:",omitempty"`
	Type            string `json:",omitempty"`
	Size            string `json:",omitempty"`
	Size64          string `json:",omitempty"`
	Checksum        string `json:",omitempty"`
	Latest          string `json:",omitempty"`
	DeleteIt        string `json:",omitempty"`
	DownloadDate    string `json:",omitempty"`
	ReleaseDate     string `json:",omitempty"`
	Url             string `json:",omitempty"`
	IUrl            string `json:",omitempty"`
	DnldStatus      string `json:",omitempty"`
	AutoloadCatalog string `json:",omitempty"`
}

func NewFirmware(firmwareFirmwareRn, parentDn, description string, firmwareFirmwareAttr FirmwareAttributes) *Firmware {
	dn := fmt.Sprintf("%s/%s", parentDn, firmwareFirmwareRn)
	return &Firmware{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FirmwareFirmwareClassName,
			Rn:                firmwareFirmwareRn,
		},
		FirmwareAttributes: firmwareFirmwareAttr,
	}
}

func (firmwareFirmware *Firmware) ToMap() (map[string]string, error) {
	firmwareFirmwareMap, err := firmwareFirmware.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(firmwareFirmwareMap, "name", firmwareFirmware.Name)
	A(firmwareFirmwareMap, "annotation", firmwareFirmware.Annotation)
	A(firmwareFirmwareMap, "fullVersion", firmwareFirmware.FullVersion)
	A(firmwareFirmwareMap, "version", firmwareFirmware.Version)
	A(firmwareFirmwareMap, "minorVersion", firmwareFirmware.MinorVersion)
	A(firmwareFirmwareMap, "isoname", firmwareFirmware.Isoname)
	A(firmwareFirmwareMap, "type", firmwareFirmware.Type)
	A(firmwareFirmwareMap, "size", firmwareFirmware.Size)
	A(firmwareFirmwareMap, "size64", firmwareFirmware.Size64)
	A(firmwareFirmwareMap, "checksum", firmwareFirmware.Checksum)
	A(firmwareFirmwareMap, "latest", firmwareFirmware.Latest)
	A(firmwareFirmwareMap, "deleteIt", firmwareFirmware.DeleteIt)
	A(firmwareFirmwareMap, "downloadDate", firmwareFirmware.DownloadDate)
	A(firmwareFirmwareMap, "releaseDate", firmwareFirmware.ReleaseDate)
	A(firmwareFirmwareMap, "url", firmwareFirmware.Url)
	A(firmwareFirmwareMap, "iUrl", firmwareFirmware.IUrl)
	A(firmwareFirmwareMap, "dnldStatus", firmwareFirmware.DnldStatus)
	A(firmwareFirmwareMap, "autoloadCatalog", firmwareFirmware.AutoloadCatalog)

	return firmwareFirmwareMap, err
}

func FirmwareFromContainerList(cont *container.Container, index int) *Firmware {

	FirmwareCont := cont.S("imdata").Index(index).S(FirmwareFirmwareClassName, "attributes")
	return &Firmware{
		BaseAttributes{
			DistinguishedName: G(FirmwareCont, "dn"),
			Description:       G(FirmwareCont, "descr"),
			Status:            G(FirmwareCont, "status"),
			ClassName:         FirmwareFirmwareClassName,
			Rn:                G(FirmwareCont, "rn"),
		},

		FirmwareAttributes{
			Name:            G(FirmwareCont, "name"),
			Annotation:      G(FirmwareCont, "annotation"),
			FullVersion:     G(FirmwareCont, "fullVersion"),
			Version:         G(FirmwareCont, "version"),
			MinorVersion:    G(FirmwareCont, "minorVersion"),
			Isoname:         G(FirmwareCont, "isoname"),
			Type:            G(FirmwareCont, "type"),
			Size:            G(FirmwareCont, "size"),
			Size64:          G(FirmwareCont, "size64"),
			Checksum:        G(FirmwareCont, "checksum"),
			Latest:          G(FirmwareCont, "latest"),
			DeleteIt:        G(FirmwareCont, "deleteIt"),
			DownloadDate:    G(FirmwareCont, "downloadDate"),
			ReleaseDate:     G(FirmwareCont, "releaseDate"),
			Url:             G(FirmwareCont, "url"),
			IUrl:            G(FirmwareCont, "iUrl"),
			DnldStatus:      G(FirmwareCont, "dnldStatus"),
			AutoloadCatalog: G(FirmwareCont, "autoloadCatalog"),
		},
	}
}

func FirmwareFromContainer(cont *container.Container) *Firmware {

	return FirmwareFromContainerList(cont, 0)
}

func FirmwareListFromContainer(cont *container.Container) []*Firmware {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*Firmware, length)

	for i := 0; i < length; i++ {

		arr[i] = FirmwareFromContainerList(cont, i)
	}

	return arr
}

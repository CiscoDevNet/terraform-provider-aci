package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const FirmwareDownloadClassName = "firmwareDownload"

type Download struct {
	BaseAttributes
	DownloadAttributes
}

type DownloadAttributes struct {
	Name        string `json:",omitempty"`
	Annotation  string `json:",omitempty"`
	Url         string `json:",omitempty"`
	LastPolled  string `json:",omitempty"`
	OperSt      string `json:",omitempty"`
	OperQual    string `json:",omitempty"`
	OperQualStr string `json:",omitempty"`
	DnldPercent string `json:",omitempty"`
}

func NewDownload(firmwareDownloadRn, parentDn, description string, firmwareDownloadAttr DownloadAttributes) *Download {
	dn := fmt.Sprintf("%s/%s", parentDn, firmwareDownloadRn)
	return &Download{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FirmwareDownloadClassName,
			Rn:                firmwareDownloadRn,
		},
		DownloadAttributes: firmwareDownloadAttr,
	}
}

func (firmwareDownload *Download) ToMap() (map[string]string, error) {
	firmwareDownloadMap, err := firmwareDownload.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}

	A(firmwareDownloadMap, "name", firmwareDownload.Name)
	A(firmwareDownloadMap, "annotation", firmwareDownload.Annotation)
	A(firmwareDownloadMap, "url", firmwareDownload.Url)
	A(firmwareDownloadMap, "lastPolled", firmwareDownload.LastPolled)
	A(firmwareDownloadMap, "operSt", firmwareDownload.OperSt)
	A(firmwareDownloadMap, "operQual", firmwareDownload.OperQual)
	A(firmwareDownloadMap, "operQualStr", firmwareDownload.OperQualStr)
	A(firmwareDownloadMap, "dnldPercent", firmwareDownload.DnldPercent)

	return firmwareDownloadMap, err
}

func DownloadFromContainerList(cont *container.Container, index int) *Download {

	DownloadCont := cont.S("imdata").Index(index).S(FirmwareDownloadClassName, "attributes")
	return &Download{
		BaseAttributes{
			DistinguishedName: G(DownloadCont, "dn"),
			Description:       G(DownloadCont, "descr"),
			Status:            G(DownloadCont, "status"),
			ClassName:         FirmwareDownloadClassName,
			Rn:                G(DownloadCont, "rn"),
		},

		DownloadAttributes{
			Name:        G(DownloadCont, "name"),
			Annotation:  G(DownloadCont, "annotation"),
			Url:         G(DownloadCont, "url"),
			LastPolled:  G(DownloadCont, "lastPolled"),
			OperSt:      G(DownloadCont, "operSt"),
			OperQual:    G(DownloadCont, "operQual"),
			OperQualStr: G(DownloadCont, "operQualStr"),
			DnldPercent: G(DownloadCont, "dnldPercent"),
		},
	}
}

func DownloadFromContainer(cont *container.Container) *Download {

	return DownloadFromContainerList(cont, 0)
}

func DownloadListFromContainer(cont *container.Container) []*Download {
	length, _ := strconv.Atoi(G(cont, "totalCount"))

	arr := make([]*Download, length)

	for i := 0; i < length; i++ {

		arr[i] = DownloadFromContainerList(cont, i)
	}

	return arr
}

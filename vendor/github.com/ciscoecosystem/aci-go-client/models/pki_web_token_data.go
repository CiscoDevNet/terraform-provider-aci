package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnpkiWebTokenData        = "uni/userext/pkiext/webtokendata"
	RnpkiWebTokenData        = "webtokendata"
	ParentDnpkiWebTokenData  = "uni/userext/pkiext"
	PkiwebtokendataClassName = "pkiWebTokenData"
)

type WebTokenData struct {
	BaseAttributes
	NameAliasAttribute
	WebTokenDataAttributes
}

type WebTokenDataAttributes struct {
	Annotation             string `json:",omitempty"`
	JwtApiKey              string `json:",omitempty"`
	JwtPrivateKey          string `json:",omitempty"`
	JwtPublicKey           string `json:",omitempty"`
	MaximumValidityPeriod  string `json:",omitempty"`
	Name                   string `json:",omitempty"`
	SessionRecordFlags     string `json:",omitempty"`
	UiIdleTimeoutSeconds   string `json:",omitempty"`
	WebtokenTimeoutSeconds string `json:",omitempty"`
}

func NewWebTokenData(pkiWebTokenDataRn, parentDn, description, nameAlias string, pkiWebTokenDataAttr WebTokenDataAttributes) *WebTokenData {
	dn := fmt.Sprintf("%s/%s", parentDn, pkiWebTokenDataRn)
	return &WebTokenData{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PkiwebtokendataClassName,
			Rn:                pkiWebTokenDataRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		WebTokenDataAttributes: pkiWebTokenDataAttr,
	}
}

func (pkiWebTokenData *WebTokenData) ToMap() (map[string]string, error) {
	pkiWebTokenDataMap, err := pkiWebTokenData.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := pkiWebTokenData.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(pkiWebTokenDataMap, key, value)
	}
	A(pkiWebTokenDataMap, "annotation", pkiWebTokenData.Annotation)
	A(pkiWebTokenDataMap, "jwtApiKey", pkiWebTokenData.JwtApiKey)
	A(pkiWebTokenDataMap, "jwtPrivateKey", pkiWebTokenData.JwtPrivateKey)
	A(pkiWebTokenDataMap, "jwtPublicKey", pkiWebTokenData.JwtPublicKey)
	A(pkiWebTokenDataMap, "maximumValidityPeriod", pkiWebTokenData.MaximumValidityPeriod)
	A(pkiWebTokenDataMap, "name", pkiWebTokenData.Name)
	A(pkiWebTokenDataMap, "sessionRecordFlags", pkiWebTokenData.SessionRecordFlags)
	A(pkiWebTokenDataMap, "uiIdleTimeoutSeconds", pkiWebTokenData.UiIdleTimeoutSeconds)
	A(pkiWebTokenDataMap, "webtokenTimeoutSeconds", pkiWebTokenData.WebtokenTimeoutSeconds)
	return pkiWebTokenDataMap, err
}

func WebTokenDataFromContainerList(cont *container.Container, index int) *WebTokenData {
	WebTokenDataCont := cont.S("imdata").Index(index).S(PkiwebtokendataClassName, "attributes")
	return &WebTokenData{
		BaseAttributes{
			DistinguishedName: G(WebTokenDataCont, "dn"),
			Description:       G(WebTokenDataCont, "descr"),
			Status:            G(WebTokenDataCont, "status"),
			ClassName:         PkiwebtokendataClassName,
			Rn:                G(WebTokenDataCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(WebTokenDataCont, "nameAlias"),
		},
		WebTokenDataAttributes{
			Annotation:             G(WebTokenDataCont, "annotation"),
			JwtApiKey:              G(WebTokenDataCont, "jwtApiKey"),
			JwtPrivateKey:          G(WebTokenDataCont, "jwtPrivateKey"),
			JwtPublicKey:           G(WebTokenDataCont, "jwtPublicKey"),
			MaximumValidityPeriod:  G(WebTokenDataCont, "maximumValidityPeriod"),
			Name:                   G(WebTokenDataCont, "name"),
			SessionRecordFlags:     G(WebTokenDataCont, "sessionRecordFlags"),
			UiIdleTimeoutSeconds:   G(WebTokenDataCont, "uiIdleTimeoutSeconds"),
			WebtokenTimeoutSeconds: G(WebTokenDataCont, "webtokenTimeoutSeconds"),
		},
	}
}

func WebTokenDataFromContainer(cont *container.Container) *WebTokenData {
	return WebTokenDataFromContainerList(cont, 0)
}

func WebTokenDataListFromContainer(cont *container.Container) []*WebTokenData {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*WebTokenData, length)
	for i := 0; i < length; i++ {
		arr[i] = WebTokenDataFromContainerList(cont, i)
	}
	return arr
}

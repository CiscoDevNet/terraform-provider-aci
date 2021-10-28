package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnaaaSamlEncCert        = "uni/userext/samlext/samlenccert-%s"
	RnaaaSamlEncCert        = "samlenccert-%s"
	ParentDnaaaSamlEncCert  = "uni/userext/samlext"
	AaasamlenccertClassName = "aaaSamlEncCert"
)

type KeypairforSAMLEncryption struct {
	BaseAttributes
	NameAliasAttribute
	KeypairforSAMLEncryptionAttributes
}

type KeypairforSAMLEncryptionAttributes struct {
	Annotation     string `json:",omitempty"`
	Name           string `json:",omitempty"`
	Regenerate     string `json:",omitempty"`
	Cert           string `json:",omitempty"`
	CertValidUntil string `json:",omitempty"`
	ExpState       string `json:",omitempty"`
}

func NewKeypairforSAMLEncryption(aaaSamlEncCertRn, parentDn, description, nameAlias string, aaaSamlEncCertAttr KeypairforSAMLEncryptionAttributes) *KeypairforSAMLEncryption {
	dn := fmt.Sprintf("%s/%s", parentDn, aaaSamlEncCertRn)
	return &KeypairforSAMLEncryption{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         AaasamlenccertClassName,
			Rn:                aaaSamlEncCertRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		KeypairforSAMLEncryptionAttributes: aaaSamlEncCertAttr,
	}
}

func (aaaSamlEncCert *KeypairforSAMLEncryption) ToMap() (map[string]string, error) {
	aaaSamlEncCertMap, err := aaaSamlEncCert.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := aaaSamlEncCert.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(aaaSamlEncCertMap, key, value)
	}
	A(aaaSamlEncCertMap, "annotation", aaaSamlEncCert.Annotation)
	A(aaaSamlEncCertMap, "name", aaaSamlEncCert.Name)
	A(aaaSamlEncCertMap, "regenerate", aaaSamlEncCert.Regenerate)
	A(aaaSamlEncCertMap, "cert", aaaSamlEncCert.Cert)
	A(aaaSamlEncCertMap, "certValidUntil", aaaSamlEncCert.CertValidUntil)
	A(aaaSamlEncCertMap, "expState", aaaSamlEncCert.ExpState)
	return aaaSamlEncCertMap, err
}

func KeypairforSAMLEncryptionFromContainerList(cont *container.Container, index int) *KeypairforSAMLEncryption {
	KeypairforSAMLEncryptionCont := cont.S("imdata").Index(index).S(AaasamlenccertClassName, "attributes")
	return &KeypairforSAMLEncryption{
		BaseAttributes{
			DistinguishedName: G(KeypairforSAMLEncryptionCont, "dn"),
			Description:       G(KeypairforSAMLEncryptionCont, "descr"),
			Status:            G(KeypairforSAMLEncryptionCont, "status"),
			ClassName:         AaasamlenccertClassName,
			Rn:                G(KeypairforSAMLEncryptionCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(KeypairforSAMLEncryptionCont, "nameAlias"),
		},
		KeypairforSAMLEncryptionAttributes{
			Annotation:     G(KeypairforSAMLEncryptionCont, "annotation"),
			Name:           G(KeypairforSAMLEncryptionCont, "name"),
			Regenerate:     G(KeypairforSAMLEncryptionCont, "regenerate"),
			Cert:           G(KeypairforSAMLEncryptionCont, "cert"),
			CertValidUntil: G(KeypairforSAMLEncryptionCont, "certValidUntil"),
			ExpState:       G(KeypairforSAMLEncryptionCont, "expState"),
		},
	}
}

func KeypairforSAMLEncryptionFromContainer(cont *container.Container) *KeypairforSAMLEncryption {
	return KeypairforSAMLEncryptionFromContainerList(cont, 0)
}

func KeypairforSAMLEncryptionListFromContainer(cont *container.Container) []*KeypairforSAMLEncryption {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*KeypairforSAMLEncryption, length)
	for i := 0; i < length; i++ {
		arr[i] = KeypairforSAMLEncryptionFromContainerList(cont, i)
	}
	return arr
}

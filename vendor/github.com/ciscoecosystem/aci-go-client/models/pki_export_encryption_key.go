package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnpkiExportEncryptionKey        = "uni/exportcryptkey"
	RnpkiExportEncryptionKey        = "exportcryptkey"
	ParentDnpkiExportEncryptionKey  = "uni"
	PkiexportencryptionkeyClassName = "pkiExportEncryptionKey"
)

type AESEncryptionPassphraseandKeysforConfigExportImport struct {
	BaseAttributes
	NameAliasAttribute
	AESEncryptionPassphraseandKeysforConfigExportImportAttributes
}

type AESEncryptionPassphraseandKeysforConfigExportImportAttributes struct {
	Annotation                     string `json:",omitempty"`
	ClearEncryptionKey             string `json:",omitempty"`
	Name                           string `json:",omitempty"`
	Passphrase                     string `json:",omitempty"`
	PassphraseKeyDerivationVersion string `json:",omitempty"`
	StrongEncryptionEnabled        string `json:",omitempty"`
	KeyConfigured                  string `json:",omitempty"`
}

func NewAESEncryptionPassphraseandKeysforConfigExportImport(pkiExportEncryptionKeyRn, parentDn, description, nameAlias string, pkiExportEncryptionKeyAttr AESEncryptionPassphraseandKeysforConfigExportImportAttributes) *AESEncryptionPassphraseandKeysforConfigExportImport {
	dn := fmt.Sprintf("%s/%s", parentDn, pkiExportEncryptionKeyRn)
	return &AESEncryptionPassphraseandKeysforConfigExportImport{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         PkiexportencryptionkeyClassName,
			Rn:                pkiExportEncryptionKeyRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		AESEncryptionPassphraseandKeysforConfigExportImportAttributes: pkiExportEncryptionKeyAttr,
	}
}

func (pkiExportEncryptionKey *AESEncryptionPassphraseandKeysforConfigExportImport) ToMap() (map[string]string, error) {
	pkiExportEncryptionKeyMap, err := pkiExportEncryptionKey.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := pkiExportEncryptionKey.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(pkiExportEncryptionKeyMap, key, value)
	}
	A(pkiExportEncryptionKeyMap, "annotation", pkiExportEncryptionKey.Annotation)
	A(pkiExportEncryptionKeyMap, "clearEncryptionKey", pkiExportEncryptionKey.ClearEncryptionKey)
	A(pkiExportEncryptionKeyMap, "name", pkiExportEncryptionKey.Name)
	A(pkiExportEncryptionKeyMap, "passphrase", pkiExportEncryptionKey.Passphrase)
	A(pkiExportEncryptionKeyMap, "passphraseKeyDerivationVersion", pkiExportEncryptionKey.PassphraseKeyDerivationVersion)
	A(pkiExportEncryptionKeyMap, "strongEncryptionEnabled", pkiExportEncryptionKey.StrongEncryptionEnabled)
	A(pkiExportEncryptionKeyMap, "keyConfigured", pkiExportEncryptionKey.KeyConfigured)
	return pkiExportEncryptionKeyMap, err
}

func AESEncryptionPassphraseandKeysforConfigExportImportFromContainerList(cont *container.Container, index int) *AESEncryptionPassphraseandKeysforConfigExportImport {
	AESEncryptionPassphraseandKeysforConfigExportImportCont := cont.S("imdata").Index(index).S(PkiexportencryptionkeyClassName, "attributes")
	return &AESEncryptionPassphraseandKeysforConfigExportImport{
		BaseAttributes{
			DistinguishedName: G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "dn"),
			Description:       G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "descr"),
			Status:            G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "status"),
			ClassName:         PkiexportencryptionkeyClassName,
			Rn:                G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "nameAlias"),
		},
		AESEncryptionPassphraseandKeysforConfigExportImportAttributes{
			Annotation:                     G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "annotation"),
			ClearEncryptionKey:             G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "clearEncryptionKey"),
			Name:                           G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "name"),
			Passphrase:                     G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "passphrase"),
			PassphraseKeyDerivationVersion: G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "passphraseKeyDerivationVersion"),
			StrongEncryptionEnabled:        G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "strongEncryptionEnabled"),
			KeyConfigured:                  G(AESEncryptionPassphraseandKeysforConfigExportImportCont, "keyConfigured"),
		},
	}
}

func AESEncryptionPassphraseandKeysforConfigExportImportFromContainer(cont *container.Container) *AESEncryptionPassphraseandKeysforConfigExportImport {
	return AESEncryptionPassphraseandKeysforConfigExportImportFromContainerList(cont, 0)
}

func AESEncryptionPassphraseandKeysforConfigExportImportListFromContainer(cont *container.Container) []*AESEncryptionPassphraseandKeysforConfigExportImport {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*AESEncryptionPassphraseandKeysforConfigExportImport, length)
	for i := 0; i < length; i++ {
		arr[i] = AESEncryptionPassphraseandKeysforConfigExportImportFromContainerList(cont, i)
	}
	return arr
}

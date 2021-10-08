package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateAESEncryptionPassphraseandKeysforConfigExportImport(description string, nameAlias string, pkiExportEncryptionKeyAttr models.AESEncryptionPassphraseandKeysforConfigExportImportAttributes) (*models.AESEncryptionPassphraseandKeysforConfigExportImport, error) {
	rn := fmt.Sprintf(models.RnpkiExportEncryptionKey)
	parentDn := fmt.Sprintf(models.ParentDnpkiExportEncryptionKey)
	pkiExportEncryptionKey := models.NewAESEncryptionPassphraseandKeysforConfigExportImport(rn, parentDn, description, nameAlias, pkiExportEncryptionKeyAttr)
	err := sm.Save(pkiExportEncryptionKey)
	return pkiExportEncryptionKey, err
}

func (sm *ServiceManager) ReadAESEncryptionPassphraseandKeysforConfigExportImport() (*models.AESEncryptionPassphraseandKeysforConfigExportImport, error) {
	dn := fmt.Sprintf(models.DnpkiExportEncryptionKey)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	pkiExportEncryptionKey := models.AESEncryptionPassphraseandKeysforConfigExportImportFromContainer(cont)
	return pkiExportEncryptionKey, nil
}

func (sm *ServiceManager) DeleteAESEncryptionPassphraseandKeysforConfigExportImport() error {
	dn := fmt.Sprintf(models.DnpkiExportEncryptionKey)
	return sm.DeleteByDn(dn, models.PkiexportencryptionkeyClassName)
}

func (sm *ServiceManager) UpdateAESEncryptionPassphraseandKeysforConfigExportImport(description string, nameAlias string, pkiExportEncryptionKeyAttr models.AESEncryptionPassphraseandKeysforConfigExportImportAttributes) (*models.AESEncryptionPassphraseandKeysforConfigExportImport, error) {
	rn := fmt.Sprintf(models.RnpkiExportEncryptionKey)
	parentDn := fmt.Sprintf(models.ParentDnpkiExportEncryptionKey)
	pkiExportEncryptionKey := models.NewAESEncryptionPassphraseandKeysforConfigExportImport(rn, parentDn, description, nameAlias, pkiExportEncryptionKeyAttr)
	pkiExportEncryptionKey.Status = "modified"
	err := sm.Save(pkiExportEncryptionKey)
	return pkiExportEncryptionKey, err
}

func (sm *ServiceManager) ListAESEncryptionPassphraseandKeysforConfigExportImport() ([]*models.AESEncryptionPassphraseandKeysforConfigExportImport, error) {
	dnUrl := fmt.Sprintf("%s/uni/pkiExportEncryptionKey.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.AESEncryptionPassphraseandKeysforConfigExportImportListFromContainer(cont)
	return list, err
}

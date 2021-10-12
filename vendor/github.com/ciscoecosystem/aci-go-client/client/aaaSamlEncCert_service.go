package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateKeypairforSAMLEncryption(name string, description string, nameAlias string, aaaSamlEncCertAttr models.KeypairforSAMLEncryptionAttributes) (*models.KeypairforSAMLEncryption, error) {
	rn := fmt.Sprintf(models.RnaaaSamlEncCert, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaSamlEncCert)
	aaaSamlEncCert := models.NewKeypairforSAMLEncryption(rn, parentDn, description, nameAlias, aaaSamlEncCertAttr)
	err := sm.Save(aaaSamlEncCert)
	return aaaSamlEncCert, err
}

func (sm *ServiceManager) ReadKeypairforSAMLEncryption(name string) (*models.KeypairforSAMLEncryption, error) {
	dn := fmt.Sprintf(models.DnaaaSamlEncCert, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaSamlEncCert := models.KeypairforSAMLEncryptionFromContainer(cont)
	return aaaSamlEncCert, nil
}

func (sm *ServiceManager) DeleteKeypairforSAMLEncryption(name string) error {
	dn := fmt.Sprintf(models.DnaaaSamlEncCert, name)
	return sm.DeleteByDn(dn, models.AaasamlenccertClassName)
}

func (sm *ServiceManager) UpdateKeypairforSAMLEncryption(name string, description string, nameAlias string, aaaSamlEncCertAttr models.KeypairforSAMLEncryptionAttributes) (*models.KeypairforSAMLEncryption, error) {
	rn := fmt.Sprintf(models.RnaaaSamlEncCert, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaSamlEncCert)
	aaaSamlEncCert := models.NewKeypairforSAMLEncryption(rn, parentDn, description, nameAlias, aaaSamlEncCertAttr)
	aaaSamlEncCert.Status = "modified"
	err := sm.Save(aaaSamlEncCert)
	return aaaSamlEncCert, err
}

func (sm *ServiceManager) ListKeypairforSAMLEncryption() ([]*models.KeypairforSAMLEncryption, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/samlext/aaaSamlEncCert.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.KeypairforSAMLEncryptionListFromContainer(cont)
	return list, err
}

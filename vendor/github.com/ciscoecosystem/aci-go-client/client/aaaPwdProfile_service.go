package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreatePasswordChangeExpirationPolicy(description string, nameAlias string, aaaPwdProfileAttr models.PasswordChangeExpirationPolicyAttributes) (*models.PasswordChangeExpirationPolicy, error) {
	rn := fmt.Sprintf(models.RnaaaPwdProfile)
	parentDn := fmt.Sprintf(models.ParentDnaaaPwdProfile)
	aaaPwdProfile := models.NewPasswordChangeExpirationPolicy(rn, parentDn, description, nameAlias, aaaPwdProfileAttr)
	err := sm.Save(aaaPwdProfile)
	return aaaPwdProfile, err
}

func (sm *ServiceManager) ReadPasswordChangeExpirationPolicy() (*models.PasswordChangeExpirationPolicy, error) {
	dn := fmt.Sprintf(models.DnaaaPwdProfile)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaPwdProfile := models.PasswordChangeExpirationPolicyFromContainer(cont)
	return aaaPwdProfile, nil
}

func (sm *ServiceManager) DeletePasswordChangeExpirationPolicy() error {
	dn := fmt.Sprintf(models.DnaaaPwdProfile)
	return sm.DeleteByDn(dn, models.AaapwdprofileClassName)
}

func (sm *ServiceManager) UpdatePasswordChangeExpirationPolicy(description string, nameAlias string, aaaPwdProfileAttr models.PasswordChangeExpirationPolicyAttributes) (*models.PasswordChangeExpirationPolicy, error) {
	rn := fmt.Sprintf(models.RnaaaPwdProfile)
	parentDn := fmt.Sprintf(models.ParentDnaaaPwdProfile)
	aaaPwdProfile := models.NewPasswordChangeExpirationPolicy(rn, parentDn, description, nameAlias, aaaPwdProfileAttr)
	aaaPwdProfile.Status = "modified"
	err := sm.Save(aaaPwdProfile)
	return aaaPwdProfile, err
}

func (sm *ServiceManager) ListPasswordChangeExpirationPolicy() ([]*models.PasswordChangeExpirationPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/aaaPwdProfile.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.PasswordChangeExpirationPolicyListFromContainer(cont)
	return list, err
}

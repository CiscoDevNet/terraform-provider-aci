package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateBlockUserLoginsPolicy(description string, nameAlias string, aaaBlockLoginProfileAttr models.BlockUserLoginsPolicyAttributes) (*models.BlockUserLoginsPolicy, error) {
	rn := fmt.Sprintf(models.RnaaaBlockLoginProfile)
	parentDn := fmt.Sprintf(models.ParentDnaaaBlockLoginProfile)
	aaaBlockLoginProfile := models.NewBlockUserLoginsPolicy(rn, parentDn, description, nameAlias, aaaBlockLoginProfileAttr)
	err := sm.Save(aaaBlockLoginProfile)
	return aaaBlockLoginProfile, err
}

func (sm *ServiceManager) ReadBlockUserLoginsPolicy() (*models.BlockUserLoginsPolicy, error) {
	dn := fmt.Sprintf(models.DnaaaBlockLoginProfile)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaBlockLoginProfile := models.BlockUserLoginsPolicyFromContainer(cont)
	return aaaBlockLoginProfile, nil
}

func (sm *ServiceManager) DeleteBlockUserLoginsPolicy() error {
	dn := fmt.Sprintf(models.DnaaaBlockLoginProfile)
	return sm.DeleteByDn(dn, models.AaablockloginprofileClassName)
}

func (sm *ServiceManager) UpdateBlockUserLoginsPolicy(description string, nameAlias string, aaaBlockLoginProfileAttr models.BlockUserLoginsPolicyAttributes) (*models.BlockUserLoginsPolicy, error) {
	rn := fmt.Sprintf(models.RnaaaBlockLoginProfile)
	parentDn := fmt.Sprintf(models.ParentDnaaaBlockLoginProfile)
	aaaBlockLoginProfile := models.NewBlockUserLoginsPolicy(rn, parentDn, description, nameAlias, aaaBlockLoginProfileAttr)
	aaaBlockLoginProfile.Status = "modified"
	err := sm.Save(aaaBlockLoginProfile)
	return aaaBlockLoginProfile, err
}

func (sm *ServiceManager) ListBlockUserLoginsPolicy() ([]*models.BlockUserLoginsPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/aaaBlockLoginProfile.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.BlockUserLoginsPolicyListFromContainer(cont)
	return list, err
}

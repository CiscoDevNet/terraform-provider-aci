package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateDefaultRadiusAuthenticationSettings(description string, nameAlias string, aaaPingEpAttr models.DefaultRadiusAuthenticationSettingsAttributes) (*models.DefaultRadiusAuthenticationSettings, error) {
	rn := fmt.Sprintf(models.RnaaaPingEp)
	parentDn := fmt.Sprintf(models.ParentDnaaaPingEp)
	aaaPingEp := models.NewDefaultRadiusAuthenticationSettings(rn, parentDn, description, nameAlias, aaaPingEpAttr)
	err := sm.Save(aaaPingEp)
	return aaaPingEp, err
}

func (sm *ServiceManager) ReadDefaultRadiusAuthenticationSettings() (*models.DefaultRadiusAuthenticationSettings, error) {
	dn := fmt.Sprintf(models.DnaaaPingEp)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaPingEp := models.DefaultRadiusAuthenticationSettingsFromContainer(cont)
	return aaaPingEp, nil
}

func (sm *ServiceManager) DeleteDefaultRadiusAuthenticationSettings() error {
	dn := fmt.Sprintf(models.DnaaaPingEp)
	return sm.DeleteByDn(dn, models.AaapingepClassName)
}

func (sm *ServiceManager) UpdateDefaultRadiusAuthenticationSettings(description string, nameAlias string, aaaPingEpAttr models.DefaultRadiusAuthenticationSettingsAttributes) (*models.DefaultRadiusAuthenticationSettings, error) {
	rn := fmt.Sprintf(models.RnaaaPingEp)
	parentDn := fmt.Sprintf(models.ParentDnaaaPingEp)
	aaaPingEp := models.NewDefaultRadiusAuthenticationSettings(rn, parentDn, description, nameAlias, aaaPingEpAttr)
	aaaPingEp.Status = "modified"
	err := sm.Save(aaaPingEp)
	return aaaPingEp, err
}

func (sm *ServiceManager) ListDefaultRadiusAuthenticationSettings() ([]*models.DefaultRadiusAuthenticationSettings, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/aaaPingEp.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.DefaultRadiusAuthenticationSettingsListFromContainer(cont)
	return list, err
}

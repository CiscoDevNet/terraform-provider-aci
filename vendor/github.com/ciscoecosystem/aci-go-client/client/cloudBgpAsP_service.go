package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateAutonomousSystemProfile(description string, cloudBgpAsPattr models.AutonomousSystemProfileAttributes) (*models.AutonomousSystemProfile, error) {
	rn := fmt.Sprintf("clouddomp/as")
	parentDn := fmt.Sprintf("uni")
	cloudBgpAsP := models.NewAutonomousSystemProfile(rn, parentDn, description, cloudBgpAsPattr)
	err := sm.Save(cloudBgpAsP)
	return cloudBgpAsP, err
}

func (sm *ServiceManager) ReadAutonomousSystemProfile() (*models.AutonomousSystemProfile, error) {
	dn := fmt.Sprintf("uni/clouddomp/as")
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudBgpAsP := models.AutonomousSystemProfileFromContainer(cont)
	return cloudBgpAsP, nil
}

func (sm *ServiceManager) DeleteAutonomousSystemProfile() error {
	dn := fmt.Sprintf("uni/clouddomp/as")
	return sm.DeleteByDn(dn, models.CloudbgpaspClassName)
}

func (sm *ServiceManager) UpdateAutonomousSystemProfile(description string, cloudBgpAsPattr models.AutonomousSystemProfileAttributes) (*models.AutonomousSystemProfile, error) {
	rn := fmt.Sprintf("clouddomp/as")
	parentDn := fmt.Sprintf("uni")
	cloudBgpAsP := models.NewAutonomousSystemProfile(rn, parentDn, description, cloudBgpAsPattr)

	cloudBgpAsP.Status = "modified"
	err := sm.Save(cloudBgpAsP)
	return cloudBgpAsP, err

}

func (sm *ServiceManager) ListAutonomousSystemProfile() ([]*models.AutonomousSystemProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/cloudBgpAsP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.AutonomousSystemProfileListFromContainer(cont)

	return list, err
}

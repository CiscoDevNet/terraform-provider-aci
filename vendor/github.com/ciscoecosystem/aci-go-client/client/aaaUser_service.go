package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLocalUser(name string, description string, aaaUserattr models.LocalUserAttributes) (*models.LocalUser, error) {
	rn := fmt.Sprintf("userext/user-%s", name)
	parentDn := fmt.Sprintf("uni")
	aaaUser := models.NewLocalUser(rn, parentDn, description, aaaUserattr)
	err := sm.Save(aaaUser)
	return aaaUser, err
}

func (sm *ServiceManager) ReadLocalUser(name string) (*models.LocalUser, error) {
	dn := fmt.Sprintf("uni/userext/user-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	aaaUser := models.LocalUserFromContainer(cont)
	return aaaUser, nil
}

func (sm *ServiceManager) DeleteLocalUser(name string) error {
	dn := fmt.Sprintf("uni/userext/user-%s", name)
	return sm.DeleteByDn(dn, models.AaauserClassName)
}

func (sm *ServiceManager) UpdateLocalUser(name string, description string, aaaUserattr models.LocalUserAttributes) (*models.LocalUser, error) {
	rn := fmt.Sprintf("userext/user-%s", name)
	parentDn := fmt.Sprintf("uni")
	aaaUser := models.NewLocalUser(rn, parentDn, description, aaaUserattr)

	aaaUser.Status = "modified"
	err := sm.Save(aaaUser)
	return aaaUser, err

}

func (sm *ServiceManager) ListLocalUser() ([]*models.LocalUser, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/aaaUser.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LocalUserListFromContainer(cont)

	return list, err
}

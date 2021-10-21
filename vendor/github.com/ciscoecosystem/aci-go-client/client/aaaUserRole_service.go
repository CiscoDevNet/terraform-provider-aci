package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateUserRole(name string, user_domain string, local_user string, description string, nameAlias string, aaaUserRoleAttr models.UserRoleAttributes) (*models.UserRole, error) {
	rn := fmt.Sprintf(models.RnaaaUserRole, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaUserRole, local_user, user_domain)
	aaaUserRole := models.NewUserRole(rn, parentDn, description, nameAlias, aaaUserRoleAttr)
	err := sm.Save(aaaUserRole)
	return aaaUserRole, err
}

func (sm *ServiceManager) ReadUserRole(name string, user_domain string, local_user string) (*models.UserRole, error) {
	dn := fmt.Sprintf(models.DnaaaUserRole, local_user, user_domain, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaUserRole := models.UserRoleFromContainer(cont)
	return aaaUserRole, nil
}

func (sm *ServiceManager) DeleteUserRole(name string, user_domain string, local_user string) error {
	dn := fmt.Sprintf(models.DnaaaUserRole, local_user, user_domain, name)
	return sm.DeleteByDn(dn, models.AaauserroleClassName)
}

func (sm *ServiceManager) UpdateUserRole(name string, user_domain string, local_user string, description string, nameAlias string, aaaUserRoleAttr models.UserRoleAttributes) (*models.UserRole, error) {
	rn := fmt.Sprintf(models.RnaaaUserRole, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaUserRole, local_user, user_domain)
	aaaUserRole := models.NewUserRole(rn, parentDn, description, nameAlias, aaaUserRoleAttr)
	aaaUserRole.Status = "modified"
	err := sm.Save(aaaUserRole)
	return aaaUserRole, err
}

func (sm *ServiceManager) ListUserRole(user_domain string, local_user string) ([]*models.UserRole, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/user-%s/userdomain-%s/aaaUserRole.json", models.BaseurlStr, local_user, user_domain)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.UserRoleListFromContainer(cont)
	return list, err
}

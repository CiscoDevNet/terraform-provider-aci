package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateAaaDomainRef(name string, parentDn string, nameAlias string, aaaDomainRefAttr models.AaaDomainRefAttributes) (*models.AaaDomainRef, error) {
	rn := fmt.Sprintf(models.RnaaaDomainRef, name)
	aaaDomainRef := models.NewAaaDomainRef(rn, parentDn, nameAlias, aaaDomainRefAttr)
	err := sm.Save(aaaDomainRef)
	return aaaDomainRef, err
}

func (sm *ServiceManager) ReadAaaDomainRef(name string, parentDn string) (*models.AaaDomainRef, error) {
	dn := fmt.Sprintf(models.DnaaaDomainRef, parentDn, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	aaaDomainRef := models.AaaDomainRefFromContainer(cont)
	return aaaDomainRef, nil
}

func (sm *ServiceManager) DeleteAaaDomainRef(name string, parentDn string) error {
	dn := fmt.Sprintf(models.DnaaaDomainRef, parentDn, name)
	return sm.DeleteByDn(dn, models.AaadomainrefClassName)
}

func (sm *ServiceManager) UpdateAaaDomainRef(name string, parentDn string, nameAlias string, aaaDomainRefAttr models.AaaDomainRefAttributes) (*models.AaaDomainRef, error) {
	rn := fmt.Sprintf(models.RnaaaDomainRef, name)
	aaaDomainRef := models.NewAaaDomainRef(rn, parentDn, nameAlias, aaaDomainRefAttr)
	aaaDomainRef.Status = "modified"
	err := sm.Save(aaaDomainRef)
	return aaaDomainRef, err
}

func (sm *ServiceManager) ListAaaDomainRef(parentDn string) ([]*models.AaaDomainRef, error) {
	dnUrl := fmt.Sprintf("%s/%s/aaaDomainRef.json", models.BaseurlStr, parentDn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.AaaDomainRefListFromContainer(cont)
	return list, err
}

package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateMulticastAddressPool(name string, description string, fvnsMcastAddrInstPAttr models.MulticastAddressPoolAttributes) (*models.MulticastAddressPool, error) {
	rn := fmt.Sprintf(models.RnfvnsMcastAddrInstP, name)
	parentDn := fmt.Sprintf(models.ParentDnfvnsMcastAddrInstP)
	fvnsMcastAddrInstP := models.NewMulticastAddressPool(rn, parentDn, description, fvnsMcastAddrInstPAttr)
	err := sm.Save(fvnsMcastAddrInstP)
	return fvnsMcastAddrInstP, err
}

func (sm *ServiceManager) ReadMulticastAddressPool(name string) (*models.MulticastAddressPool, error) {
	dn := fmt.Sprintf(models.DnfvnsMcastAddrInstP, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvnsMcastAddrInstP := models.MulticastAddressPoolFromContainer(cont)
	return fvnsMcastAddrInstP, nil
}

func (sm *ServiceManager) DeleteMulticastAddressPool(name string) error {
	dn := fmt.Sprintf(models.DnfvnsMcastAddrInstP, name)
	return sm.DeleteByDn(dn, models.FvnsmcastaddrinstpClassName)
}

func (sm *ServiceManager) UpdateMulticastAddressPool(name string, description string, fvnsMcastAddrInstPAttr models.MulticastAddressPoolAttributes) (*models.MulticastAddressPool, error) {
	rn := fmt.Sprintf(models.RnfvnsMcastAddrInstP, name)
	parentDn := fmt.Sprintf(models.ParentDnfvnsMcastAddrInstP)
	fvnsMcastAddrInstP := models.NewMulticastAddressPool(rn, parentDn, description, fvnsMcastAddrInstPAttr)
	fvnsMcastAddrInstP.Status = "modified"
	err := sm.Save(fvnsMcastAddrInstP)
	return fvnsMcastAddrInstP, err
}

func (sm *ServiceManager) ListMulticastAddressPool() ([]*models.MulticastAddressPool, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/fvnsMcastAddrInstP.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MulticastAddressPoolListFromContainer(cont)
	return list, err
}

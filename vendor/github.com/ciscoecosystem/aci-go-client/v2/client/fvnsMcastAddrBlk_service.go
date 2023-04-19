package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateMulticastAddressBlock(to string, from string, multicast_address_pool string, description string, fvnsMcastAddrBlkAttr models.MulticastAddressBlockAttributes) (*models.MulticastAddressBlock, error) {
	rn := fmt.Sprintf(models.RnfvnsMcastAddrBlk, from, to)
	parentDn := fmt.Sprintf(models.ParentDnfvnsMcastAddrBlk, multicast_address_pool)
	fvnsMcastAddrBlk := models.NewMulticastAddressBlock(rn, parentDn, description, fvnsMcastAddrBlkAttr)
	err := sm.Save(fvnsMcastAddrBlk)
	return fvnsMcastAddrBlk, err
}

func (sm *ServiceManager) ReadMulticastAddressBlock(to string, from string, multicast_address_pool string) (*models.MulticastAddressBlock, error) {
	dn := fmt.Sprintf(models.DnfvnsMcastAddrBlk, multicast_address_pool, from, to)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvnsMcastAddrBlk := models.MulticastAddressBlockFromContainer(cont)
	return fvnsMcastAddrBlk, nil
}

func (sm *ServiceManager) DeleteMulticastAddressBlock(to string, from string, multicast_address_pool string) error {
	dn := fmt.Sprintf(models.DnfvnsMcastAddrBlk, multicast_address_pool, from, to)
	return sm.DeleteByDn(dn, models.FvnsmcastaddrblkClassName)
}

func (sm *ServiceManager) UpdateMulticastAddressBlock(to string, from string, multicast_address_pool string, description string, fvnsMcastAddrBlkAttr models.MulticastAddressBlockAttributes) (*models.MulticastAddressBlock, error) {
	rn := fmt.Sprintf(models.RnfvnsMcastAddrBlk, from, to)
	parentDn := fmt.Sprintf(models.ParentDnfvnsMcastAddrBlk, multicast_address_pool)
	fvnsMcastAddrBlk := models.NewMulticastAddressBlock(rn, parentDn, description, fvnsMcastAddrBlkAttr)
	fvnsMcastAddrBlk.Status = "modified"
	err := sm.Save(fvnsMcastAddrBlk)
	return fvnsMcastAddrBlk, err
}

func (sm *ServiceManager) ListMulticastAddressBlock(multicast_address_pool string) ([]*models.MulticastAddressBlock, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/maddrns-%s/fvnsMcastAddrBlk.json", models.BaseurlStr, multicast_address_pool)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MulticastAddressBlockListFromContainer(cont)
	return list, err
}

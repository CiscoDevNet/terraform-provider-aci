package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateVLANPool(allocMode string, name string, description string, fvnsVlanInstPattr models.VLANPoolAttributes) (*models.VLANPool, error) {
	rn := fmt.Sprintf("infra/vlanns-[%s]-%s", name, allocMode)
	parentDn := fmt.Sprintf("uni")
	fvnsVlanInstP := models.NewVLANPool(rn, parentDn, description, fvnsVlanInstPattr)
	err := sm.Save(fvnsVlanInstP)
	return fvnsVlanInstP, err
}

func (sm *ServiceManager) ReadVLANPool(allocMode string, name string) (*models.VLANPool, error) {
	dn := fmt.Sprintf("uni/infra/vlanns-[%s]-%s", name, allocMode)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvnsVlanInstP := models.VLANPoolFromContainer(cont)
	return fvnsVlanInstP, nil
}

func (sm *ServiceManager) DeleteVLANPool(allocMode string, name string) error {
	dn := fmt.Sprintf("uni/infra/vlanns-[%s]-%s", name, allocMode)
	return sm.DeleteByDn(dn, models.FvnsvlaninstpClassName)
}

func (sm *ServiceManager) UpdateVLANPool(allocMode string, name string, description string, fvnsVlanInstPattr models.VLANPoolAttributes) (*models.VLANPool, error) {
	rn := fmt.Sprintf("infra/vlanns-[%s]-%s", name, allocMode)
	parentDn := fmt.Sprintf("uni")
	fvnsVlanInstP := models.NewVLANPool(rn, parentDn, description, fvnsVlanInstPattr)

	fvnsVlanInstP.Status = "modified"
	err := sm.Save(fvnsVlanInstP)
	return fvnsVlanInstP, err

}

func (sm *ServiceManager) ListVLANPool() ([]*models.VLANPool, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/fvnsVlanInstP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.VLANPoolListFromContainer(cont)

	return list, err
}

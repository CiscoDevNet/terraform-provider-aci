package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateVXLANPool(name string, description string, fvnsVxlanInstPattr models.VXLANPoolAttributes) (*models.VXLANPool, error) {
	rn := fmt.Sprintf("infra/vxlanns-%s", name)
	parentDn := fmt.Sprintf("uni")
	fvnsVxlanInstP := models.NewVXLANPool(rn, parentDn, description, fvnsVxlanInstPattr)
	err := sm.Save(fvnsVxlanInstP)
	return fvnsVxlanInstP, err
}

func (sm *ServiceManager) ReadVXLANPool(name string) (*models.VXLANPool, error) {
	dn := fmt.Sprintf("uni/infra/vxlanns-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvnsVxlanInstP := models.VXLANPoolFromContainer(cont)
	return fvnsVxlanInstP, nil
}

func (sm *ServiceManager) DeleteVXLANPool(name string) error {
	dn := fmt.Sprintf("uni/infra/vxlanns-%s", name)
	return sm.DeleteByDn(dn, models.FvnsvxlaninstpClassName)
}

func (sm *ServiceManager) UpdateVXLANPool(name string, description string, fvnsVxlanInstPattr models.VXLANPoolAttributes) (*models.VXLANPool, error) {
	rn := fmt.Sprintf("infra/vxlanns-%s", name)
	parentDn := fmt.Sprintf("uni")
	fvnsVxlanInstP := models.NewVXLANPool(rn, parentDn, description, fvnsVxlanInstPattr)

	fvnsVxlanInstP.Status = "modified"
	err := sm.Save(fvnsVxlanInstP)
	return fvnsVxlanInstP, err

}

func (sm *ServiceManager) ListVXLANPool() ([]*models.VXLANPool, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/fvnsVxlanInstP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.VXLANPoolListFromContainer(cont)

	return list, err
}

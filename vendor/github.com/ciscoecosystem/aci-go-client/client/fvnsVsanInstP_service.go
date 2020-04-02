package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateVSANPool(allocMode string, name string, description string, fvnsVsanInstPattr models.VSANPoolAttributes) (*models.VSANPool, error) {
	rn := fmt.Sprintf("infra/vsanns-[%s]-%s", name, allocMode)
	parentDn := fmt.Sprintf("uni")
	fvnsVsanInstP := models.NewVSANPool(rn, parentDn, description, fvnsVsanInstPattr)
	err := sm.Save(fvnsVsanInstP)
	return fvnsVsanInstP, err
}

func (sm *ServiceManager) ReadVSANPool(allocMode string, name string) (*models.VSANPool, error) {
	dn := fmt.Sprintf("uni/infra/vsanns-[%s]-%s", name, allocMode)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvnsVsanInstP := models.VSANPoolFromContainer(cont)
	return fvnsVsanInstP, nil
}

func (sm *ServiceManager) DeleteVSANPool(allocMode string, name string) error {
	dn := fmt.Sprintf("uni/infra/vsanns-[%s]-%s", name, allocMode)
	return sm.DeleteByDn(dn, models.FvnsvsaninstpClassName)
}

func (sm *ServiceManager) UpdateVSANPool(allocMode string, name string, description string, fvnsVsanInstPattr models.VSANPoolAttributes) (*models.VSANPool, error) {
	rn := fmt.Sprintf("infra/vsanns-[%s]-%s", name, allocMode)
	parentDn := fmt.Sprintf("uni")
	fvnsVsanInstP := models.NewVSANPool(rn, parentDn, description, fvnsVsanInstPattr)

	fvnsVsanInstP.Status = "modified"
	err := sm.Save(fvnsVsanInstP)
	return fvnsVsanInstP, err

}

func (sm *ServiceManager) ListVSANPool() ([]*models.VSANPool, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/fvnsVsanInstP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.VSANPoolListFromContainer(cont)

	return list, err
}

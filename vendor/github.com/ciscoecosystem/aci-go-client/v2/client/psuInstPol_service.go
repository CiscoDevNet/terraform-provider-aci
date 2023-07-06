package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreatePowerSupplyRedundancyPolicy(name string, description string, psuInstPolAttr models.PsuInstPolAttributes) (*models.PsuInstPol, error) {
	rn := fmt.Sprintf(models.RnPsuInstPol, name)
	psuInstPol := models.NewPowerSupplyRedundancyPolicy(rn, description, psuInstPolAttr)
	err := sm.Save(psuInstPol)
	return psuInstPol, err
}

func (sm *ServiceManager) ReadPowerSupplyRedundancyPolicy(name string) (*models.PsuInstPol, error) {
	dn := fmt.Sprintf(models.DnPsuInstPol, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	psuInstPol := models.PsuInstPolFromContainer(cont)
	return psuInstPol, nil
}

func (sm *ServiceManager) DeletePowerSupplyRedundancyPolicy(name string) error {
	dn := fmt.Sprintf(models.DnPsuInstPol, name)
	return sm.DeleteByDn(dn, models.PsuInstPolClassName)
}

func (sm *ServiceManager) UpdatePowerSupplyRedundancyPolicy(name string, description string, psuInstPolAttr models.PsuInstPolAttributes) (*models.PsuInstPol, error) {
	rn := fmt.Sprintf(models.RnPsuInstPol, name)
	psuInstPol := models.NewPowerSupplyRedundancyPolicy(rn, description, psuInstPolAttr)
	psuInstPol.Status = "modified"
	err := sm.Save(psuInstPol)
	return psuInstPol, err
}

func (sm *ServiceManager) ListPowerSupplyRedundancyPolicy() ([]*models.PsuInstPol, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/psuInstPol.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.PsuInstPolListFromContainer(cont)
	return list, err
}

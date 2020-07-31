package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFEXProfile(name string, description string, infraFexPattr models.FEXProfileAttributes) (*models.FEXProfile, error) {
	rn := fmt.Sprintf("infra/fexprof-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraFexP := models.NewFEXProfile(rn, parentDn, description, infraFexPattr)
	err := sm.Save(infraFexP)
	return infraFexP, err
}

func (sm *ServiceManager) ReadFEXProfile(name string) (*models.FEXProfile, error) {
	dn := fmt.Sprintf("uni/infra/fexprof-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraFexP := models.FEXProfileFromContainer(cont)
	return infraFexP, nil
}

func (sm *ServiceManager) DeleteFEXProfile(name string) error {
	dn := fmt.Sprintf("uni/infra/fexprof-%s", name)
	return sm.DeleteByDn(dn, models.InfrafexpClassName)
}

func (sm *ServiceManager) UpdateFEXProfile(name string, description string, infraFexPattr models.FEXProfileAttributes) (*models.FEXProfile, error) {
	rn := fmt.Sprintf("infra/fexprof-%s", name)
	parentDn := fmt.Sprintf("uni")
	infraFexP := models.NewFEXProfile(rn, parentDn, description, infraFexPattr)

	infraFexP.Status = "modified"
	err := sm.Save(infraFexP)
	return infraFexP, err

}

func (sm *ServiceManager) ListFEXProfile() ([]*models.FEXProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infraFexP.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FEXProfileListFromContainer(cont)

	return list, err
}

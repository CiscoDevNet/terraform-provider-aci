package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLinkLevelPolicy(name string, description string, fabricHIfPolattr models.LinkLevelPolicyAttributes) (*models.LinkLevelPolicy, error) {
	rn := fmt.Sprintf("infra/hintfpol-%s", name)
	parentDn := fmt.Sprintf("uni")
	fabricHIfPol := models.NewLinkLevelPolicy(rn, parentDn, description, fabricHIfPolattr)
	err := sm.Save(fabricHIfPol)
	return fabricHIfPol, err
}

func (sm *ServiceManager) ReadLinkLevelPolicy(name string) (*models.LinkLevelPolicy, error) {
	dn := fmt.Sprintf("uni/infra/hintfpol-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricHIfPol := models.LinkLevelPolicyFromContainer(cont)
	return fabricHIfPol, nil
}

func (sm *ServiceManager) DeleteLinkLevelPolicy(name string) error {
	dn := fmt.Sprintf("uni/infra/hintfpol-%s", name)
	return sm.DeleteByDn(dn, models.FabrichifpolClassName)
}

func (sm *ServiceManager) UpdateLinkLevelPolicy(name string, description string, fabricHIfPolattr models.LinkLevelPolicyAttributes) (*models.LinkLevelPolicy, error) {
	rn := fmt.Sprintf("infra/hintfpol-%s", name)
	parentDn := fmt.Sprintf("uni")
	fabricHIfPol := models.NewLinkLevelPolicy(rn, parentDn, description, fabricHIfPolattr)

	fabricHIfPol.Status = "modified"
	err := sm.Save(fabricHIfPol)
	return fabricHIfPol, err

}

func (sm *ServiceManager) ListLinkLevelPolicy() ([]*models.LinkLevelPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/fabricHIfPol.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.LinkLevelPolicyListFromContainer(cont)

	return list, err
}

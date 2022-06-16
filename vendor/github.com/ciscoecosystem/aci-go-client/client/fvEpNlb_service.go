package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateNlbEndpoint(parent_dn string, description string, nameAlias string, fvEpNlbAttr models.NlbEndpointAttributes) (*models.NlbEndpoint, error) {
	rn := fmt.Sprintf(models.RnfvEpNlb)
	fvEpNlb := models.NewNlbEndpoint(rn, parent_dn, description, nameAlias, fvEpNlbAttr)
	err := sm.Save(fvEpNlb)
	return fvEpNlb, err
}

func (sm *ServiceManager) ReadNlbEndpoint(parent_dn string) (*models.NlbEndpoint, error) {
	dn := fmt.Sprintf("%s/%s", parent_dn, models.RnfvEpNlb)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvEpNlb := models.NlbEndpointFromContainer(cont)
	return fvEpNlb, nil
}

func (sm *ServiceManager) DeleteNlbEndpoint(parent_dn string) error {
	dn := fmt.Sprintf("%s/%s", parent_dn, models.RnfvEpNlb)
	return sm.DeleteByDn(dn, models.FvepnlbClassName)
}

func (sm *ServiceManager) UpdateNlbEndpoint(parent_dn string, description string, nameAlias string, fvEpNlbAttr models.NlbEndpointAttributes) (*models.NlbEndpoint, error) {
	rn := fmt.Sprintf(models.RnfvEpNlb)
	fvEpNlb := models.NewNlbEndpoint(rn, parent_dn, description, nameAlias, fvEpNlbAttr)
	fvEpNlb.Status = "modified"
	err := sm.Save(fvEpNlb)
	return fvEpNlb, err
}

func (sm *ServiceManager) ListNlbEndpoint(parent_dn string) ([]*models.NlbEndpoint, error) {
	dnUrl := fmt.Sprintf("%s/%s/fvEpNlb.json", models.BaseurlStr, parent_dn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.NlbEndpointListFromContainer(cont)
	return list, err
}

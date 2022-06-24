package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateEpReachability(parent_dn string, fvEpReachabilityAttr models.EpReachabilityAttributes) (*models.EpReachability, error) {
	rn := fmt.Sprintf(models.RnfvEpReachability)
	fvEpReachability := models.NewEpReachability(rn, parent_dn, fvEpReachabilityAttr)
	err := sm.Save(fvEpReachability)
	return fvEpReachability, err
}

func (sm *ServiceManager) ReadEpReachability(parent_dn string) (*models.EpReachability, error) {
	dn := "%s/%s" + fmt.Sprintf(parent_dn, models.RnfvEpReachability)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvEpReachability := models.EpReachabilityFromContainer(cont)
	return fvEpReachability, nil
}

func (sm *ServiceManager) DeleteEpReachability(parent_dn string) error {
	dn := "%s/%s" + fmt.Sprintf(parent_dn, models.RnfvEpReachability)
	return sm.DeleteByDn(dn, models.FvepreachabilityClassName)
}

func (sm *ServiceManager) UpdateEpReachability(parent_dn string, fvEpReachabilityAttr models.EpReachabilityAttributes) (*models.EpReachability, error) {
	rn := fmt.Sprintf(models.RnfvEpReachability)
	fvEpReachability := models.NewEpReachability(rn, parent_dn, fvEpReachabilityAttr)
	fvEpReachability.Status = "modified"
	err := sm.Save(fvEpReachability)
	return fvEpReachability, err
}

func (sm *ServiceManager) ListEpReachability(parent_dn string) ([]*models.EpReachability, error) {
	dnUrl := fmt.Sprintf("%s/%s/fvEpReachability.json", models.BaseurlStr, parent_dn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.EpReachabilityListFromContainer(cont)
	return list, err
}

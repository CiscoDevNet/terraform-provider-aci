package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateAnycastEndpoint(mac string, parent_dn string, description string, nameAlias string, fvEpAnycastAttr models.AnycastEndpointAttributes) (*models.AnycastEndpoint, error) {
	rn := fmt.Sprintf(models.RnfvEpAnycast, mac)
	fvEpAnycast := models.NewAnycastEndpoint(rn, parent_dn, description, nameAlias, fvEpAnycastAttr)
	err := sm.Save(fvEpAnycast)
	return fvEpAnycast, err
}

func (sm *ServiceManager) ReadAnycastEndpoint(mac string, parent_dn string) (*models.AnycastEndpoint, error) {
	dn := parent_dn + "/" + fmt.Sprintf(models.RnfvEpAnycast, mac)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvEpAnycast := models.AnycastEndpointFromContainer(cont)
	return fvEpAnycast, nil
}

func (sm *ServiceManager) DeleteAnycastEndpoint(mac string, parent_dn string) error {
	dn := parent_dn + "/" + fmt.Sprintf(models.RnfvEpAnycast, mac)
	return sm.DeleteByDn(dn, models.FvepanycastClassName)
}

func (sm *ServiceManager) UpdateAnycastEndpoint(mac string, parent_dn string, description string, nameAlias string, fvEpAnycastAttr models.AnycastEndpointAttributes) (*models.AnycastEndpoint, error) {
	rn := fmt.Sprintf(models.RnfvEpAnycast, mac)
	fvEpAnycast := models.NewAnycastEndpoint(rn, parent_dn, description, nameAlias, fvEpAnycastAttr)
	fvEpAnycast.Status = "modified"
	err := sm.Save(fvEpAnycast)
	return fvEpAnycast, err
}

func (sm *ServiceManager) ListAnycastEndpoint(parent_dn string) ([]*models.AnycastEndpoint, error) {
	dnUrl := fmt.Sprintf("%s/%s/fvEpAnycast.json", models.BaseurlStr, parent_dn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.AnycastEndpointListFromContainer(cont)
	return list, err
}

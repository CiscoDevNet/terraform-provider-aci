package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateInfraRsDomP(tDn string, attachable_access_entity_profile string, infraRsDomPAttr models.InfraRsDomPAttributes) (*models.InfraRsDomP, error) {
	rn := fmt.Sprintf(models.RninfraRsDomP, tDn)
	parentDn := fmt.Sprintf(models.ParentDninfraRsDomP, attachable_access_entity_profile)
	infraRsDomP := models.NewInfraRsDomP(rn, parentDn, infraRsDomPAttr)
	err := sm.Save(infraRsDomP)
	return infraRsDomP, err
}

func (sm *ServiceManager) ReadInfraRsDomP(tDn string, attachable_access_entity_profile string) (*models.InfraRsDomP, error) {
	dn := fmt.Sprintf(models.DninfraRsDomP, attachable_access_entity_profile, tDn)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraRsDomP := models.InfraRsDomPFromContainer(cont)
	return infraRsDomP, nil
}

func (sm *ServiceManager) DeleteInfraRsDomP(tDn string, attachable_access_entity_profile string) error {
	dn := fmt.Sprintf(models.DninfraRsDomP, attachable_access_entity_profile, tDn)
	return sm.DeleteByDn(dn, models.InfrarsdompClassName)
}

func (sm *ServiceManager) ListInfraRsDomP(attachable_access_entity_profile string) ([]*models.InfraRsDomP, error) {
	parentDn := fmt.Sprintf(models.ParentDninfraRsDomP, attachable_access_entity_profile)
	dnUrl := fmt.Sprintf("%s/%s/infraRsDomP.json", models.BaseurlStr, parentDn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.InfraRsDomPListFromContainer(cont)
	return list, err
}

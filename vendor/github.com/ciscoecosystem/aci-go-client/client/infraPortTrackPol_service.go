package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreatePortTracking(name string, description string, nameAlias string, infraPortTrackPolAttr models.PortTrackingAttributes) (*models.PortTracking, error) {
	rn := fmt.Sprintf(models.RninfraPortTrackPol, name)
	parentDn := fmt.Sprintf(models.ParentDninfraPortTrackPol)
	infraPortTrackPol := models.NewPortTracking(rn, parentDn, description, nameAlias, infraPortTrackPolAttr)
	err := sm.Save(infraPortTrackPol)
	return infraPortTrackPol, err
}

func (sm *ServiceManager) ReadPortTracking(name string) (*models.PortTracking, error) {
	dn := fmt.Sprintf(models.DninfraPortTrackPol, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	infraPortTrackPol := models.PortTrackingFromContainer(cont)
	return infraPortTrackPol, nil
}

func (sm *ServiceManager) DeletePortTracking(name string) error {
	dn := fmt.Sprintf(models.DninfraPortTrackPol, name)
	return sm.DeleteByDn(dn, models.InfraporttrackpolClassName)
}

func (sm *ServiceManager) UpdatePortTracking(name string, description string, nameAlias string, infraPortTrackPolAttr models.PortTrackingAttributes) (*models.PortTracking, error) {
	rn := fmt.Sprintf(models.RninfraPortTrackPol, name)
	parentDn := fmt.Sprintf(models.ParentDninfraPortTrackPol)
	infraPortTrackPol := models.NewPortTracking(rn, parentDn, description, nameAlias, infraPortTrackPolAttr)
	infraPortTrackPol.Status = "modified"
	err := sm.Save(infraPortTrackPol)
	return infraPortTrackPol, err
}

func (sm *ServiceManager) ListPortTracking() ([]*models.PortTracking, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/infraPortTrackPol.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.PortTrackingListFromContainer(cont)
	return list, err
}

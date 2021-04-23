package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateRouteControlProfile(name string, tenant string, description string, rtctrlProfileattr models.RouteControlProfileAttributes) (*models.RouteControlProfile, error) {
	rn := fmt.Sprintf("prof-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	rtctrlProfile := models.NewRouteControlProfile(rn, parentDn, description, rtctrlProfileattr)
	err := sm.Save(rtctrlProfile)
	return rtctrlProfile, err
}

func (sm *ServiceManager) ReadRouteControlProfile(name string, tenant string) (*models.RouteControlProfile, error) {
	dn := fmt.Sprintf("uni/tn-%s/prof-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlProfile := models.RouteControlProfileFromContainer(cont)
	return rtctrlProfile, nil
}

func (sm *ServiceManager) DeleteRouteControlProfile(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/prof-%s", tenant, name)
	return sm.DeleteByDn(dn, models.RtctrlprofileClassName)
}

func (sm *ServiceManager) UpdateRouteControlProfile(name string, tenant string, description string, rtctrlProfileattr models.RouteControlProfileAttributes) (*models.RouteControlProfile, error) {
	rn := fmt.Sprintf("prof-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	rtctrlProfile := models.NewRouteControlProfile(rn, parentDn, description, rtctrlProfileattr)

	rtctrlProfile.Status = "modified"
	err := sm.Save(rtctrlProfile)
	return rtctrlProfile, err

}

func (sm *ServiceManager) ListRouteControlProfile(tenant string) ([]*models.RouteControlProfile, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/rtctrlProfile.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.RouteControlProfileListFromContainer(cont)

	return list, err
}

package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateNexthopEpPReachability(nhAddr string, parent_dn string, description string, nameAlias string, ipNexthopEpPAttr models.NexthopEpPReachabilityAttributes) (*models.NexthopEpPReachability, error) {
	rn := fmt.Sprintf(models.RnipNexthopEpP, nhAddr)
	ipNexthopEpP := models.NewNexthopEpPReachability(rn, parent_dn, description, nameAlias, ipNexthopEpPAttr)
	err := sm.Save(ipNexthopEpP)
	return ipNexthopEpP, err
}

func (sm *ServiceManager) ReadNexthopEpPReachability(nhAddr string, parent_dn string) (*models.NexthopEpPReachability, error) {
	dn := parent_dn + "/" + fmt.Sprintf(models.RnipNexthopEpP, nhAddr)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	ipNexthopEpP := models.NexthopEpPReachabilityFromContainer(cont)
	return ipNexthopEpP, nil
}

func (sm *ServiceManager) DeleteNexthopEpPReachability(nhAddr string, parent_dn string) error {
	dn := parent_dn + "/" + fmt.Sprintf(models.RnipNexthopEpP, nhAddr)
	return sm.DeleteByDn(dn, models.IpnexthopeppClassName)
}

func (sm *ServiceManager) UpdateNexthopEpPReachability(nhAddr string, parent_dn string, description string, nameAlias string, ipNexthopEpPAttr models.NexthopEpPReachabilityAttributes) (*models.NexthopEpPReachability, error) {
	rn := fmt.Sprintf(models.RnipNexthopEpP, nhAddr)
	ipNexthopEpP := models.NewNexthopEpPReachability(rn, parent_dn, description, nameAlias, ipNexthopEpPAttr)
	ipNexthopEpP.Status = "modified"
	err := sm.Save(ipNexthopEpP)
	return ipNexthopEpP, err
}

func (sm *ServiceManager) ListNexthopEpPReachability(parent_dn string) ([]*models.NexthopEpPReachability, error) {
	dnUrl := fmt.Sprintf("%s/%s/ipNexthopEpP.json", models.BaseurlStr, parent_dn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.NexthopEpPReachabilityListFromContainer(cont)
	return list, err
}

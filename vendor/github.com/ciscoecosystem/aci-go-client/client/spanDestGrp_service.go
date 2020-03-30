package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSPANDestinationGroup(name string, tenant string, description string, spanDestGrpattr models.SPANDestinationGroupAttributes) (*models.SPANDestinationGroup, error) {
	rn := fmt.Sprintf("destgrp-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	spanDestGrp := models.NewSPANDestinationGroup(rn, parentDn, description, spanDestGrpattr)
	err := sm.Save(spanDestGrp)
	return spanDestGrp, err
}

func (sm *ServiceManager) ReadSPANDestinationGroup(name string, tenant string) (*models.SPANDestinationGroup, error) {
	dn := fmt.Sprintf("uni/tn-%s/destgrp-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	spanDestGrp := models.SPANDestinationGroupFromContainer(cont)
	return spanDestGrp, nil
}

func (sm *ServiceManager) DeleteSPANDestinationGroup(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/destgrp-%s", tenant, name)
	return sm.DeleteByDn(dn, models.SpandestgrpClassName)
}

func (sm *ServiceManager) UpdateSPANDestinationGroup(name string, tenant string, description string, spanDestGrpattr models.SPANDestinationGroupAttributes) (*models.SPANDestinationGroup, error) {
	rn := fmt.Sprintf("destgrp-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	spanDestGrp := models.NewSPANDestinationGroup(rn, parentDn, description, spanDestGrpattr)

	spanDestGrp.Status = "modified"
	err := sm.Save(spanDestGrp)
	return spanDestGrp, err

}

func (sm *ServiceManager) ListSPANDestinationGroup(tenant string) ([]*models.SPANDestinationGroup, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/spanDestGrp.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.SPANDestinationGroupListFromContainer(cont)

	return list, err
}

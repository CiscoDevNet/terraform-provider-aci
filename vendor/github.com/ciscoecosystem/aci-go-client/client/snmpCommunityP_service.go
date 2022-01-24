package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSNMPCommunity(name string, parent_dn string, description string, nameAlias string, snmpCommunityPAttr models.SNMPCommunityAttributes) (*models.SNMPCommunity, error) {
	rn := fmt.Sprintf(models.RnsnmpCommunityP, name)
	snmpCommunityP := models.NewSNMPCommunity(rn, parent_dn, description, nameAlias, snmpCommunityPAttr)
	err := sm.Save(snmpCommunityP)
	return snmpCommunityP, err
}

func (sm *ServiceManager) ReadSNMPCommunity(name string, parent_dn string) (*models.SNMPCommunity, error) {
	dn := fmt.Sprintf(models.DnsnmpCommunityP, parent_dn, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	snmpCommunityP := models.SNMPCommunityFromContainer(cont)
	return snmpCommunityP, nil
}

func (sm *ServiceManager) DeleteSNMPCommunity(name string, parent_dn string) error {
	dn := fmt.Sprintf(models.DnsnmpCommunityP, parent_dn, name)
	return sm.DeleteByDn(dn, models.SnmpcommunitypClassName)
}

func (sm *ServiceManager) UpdateSNMPCommunity(name string, parent_dn string, description string, nameAlias string, snmpCommunityPAttr models.SNMPCommunityAttributes) (*models.SNMPCommunity, error) {
	rn := fmt.Sprintf(models.RnsnmpCommunityP, name)
	snmpCommunityP := models.NewSNMPCommunity(rn, parent_dn, description, nameAlias, snmpCommunityPAttr)
	snmpCommunityP.Status = "modified"
	err := sm.Save(snmpCommunityP)
	return snmpCommunityP, err
}

func (sm *ServiceManager) ListSNMPCommunity(parent_dn string) ([]*models.SNMPCommunity, error) {
	dnUrl := fmt.Sprintf("%s/%s/snmpCommunityP.json", models.BaseurlStr, parent_dn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SNMPCommunityListFromContainer(cont)
	return list, err
}

package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSNMPCommunity(name string, vrf string, tenant string, description string, nameAlias string, snmpCommunityPAttr models.SNMPCommunityAttributes) (*models.SNMPCommunity, error) {
	rn := fmt.Sprintf(models.RnsnmpCommunityP, name)
	parentDn := fmt.Sprintf(models.ParentDnsnmpCommunityP, tenant, vrf)
	snmpCommunityP := models.NewSNMPCommunity(rn, parentDn, description, nameAlias, snmpCommunityPAttr)
	err := sm.Save(snmpCommunityP)
	return snmpCommunityP, err
}

func (sm *ServiceManager) ReadSNMPCommunity(name string, vrf string, tenant string) (*models.SNMPCommunity, error) {
	dn := fmt.Sprintf(models.DnsnmpCommunityP, tenant, vrf, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	snmpCommunityP := models.SNMPCommunityFromContainer(cont)
	return snmpCommunityP, nil
}

func (sm *ServiceManager) DeleteSNMPCommunity(name string, vrf string, tenant string) error {
	dn := fmt.Sprintf(models.DnsnmpCommunityP, tenant, vrf, name)
	return sm.DeleteByDn(dn, models.SnmpcommunitypClassName)
}

func (sm *ServiceManager) UpdateSNMPCommunity(name string, vrf string, tenant string, description string, nameAlias string, snmpCommunityPAttr models.SNMPCommunityAttributes) (*models.SNMPCommunity, error) {
	rn := fmt.Sprintf(models.RnsnmpCommunityP, name)
	parentDn := fmt.Sprintf(models.ParentDnsnmpCommunityP, tenant, vrf)
	snmpCommunityP := models.NewSNMPCommunity(rn, parentDn, description, nameAlias, snmpCommunityPAttr)
	snmpCommunityP.Status = "modified"
	err := sm.Save(snmpCommunityP)
	return snmpCommunityP, err
}

func (sm *ServiceManager) ListSNMPCommunity(vrf string, tenant string) ([]*models.SNMPCommunity, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ctx-%s/snmpctx/snmpCommunityP.json", models.BaseurlStr, tenant, vrf)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SNMPCommunityListFromContainer(cont)
	return list, err
}

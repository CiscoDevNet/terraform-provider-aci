package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSNMPContextProfile(vrf string, tenant string, nameAlias string, snmpCtxPAttr models.SNMPContextProfileAttributes) (*models.SNMPContextProfile, error) {
	rn := fmt.Sprintf(models.RnsnmpCtxP)
	parentDn := fmt.Sprintf(models.ParentDnsnmpCtxP, tenant, vrf)
	snmpCtxP := models.NewSNMPContextProfile(rn, parentDn, nameAlias, snmpCtxPAttr)
	err := sm.Save(snmpCtxP)
	return snmpCtxP, err
}

func (sm *ServiceManager) ReadSNMPContextProfile(vrf string, tenant string) (*models.SNMPContextProfile, error) {
	dn := fmt.Sprintf(models.DnsnmpCtxP, tenant, vrf)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	snmpCtxP := models.SNMPContextProfileFromContainer(cont)
	return snmpCtxP, nil
}

func (sm *ServiceManager) DeleteSNMPContextProfile(vrf string, tenant string) error {
	dn := fmt.Sprintf(models.DnsnmpCtxP, tenant, vrf)
	return sm.DeleteByDn(dn, models.SnmpctxpClassName)
}

func (sm *ServiceManager) UpdateSNMPContextProfile(vrf string, tenant string, nameAlias string, snmpCtxPAttr models.SNMPContextProfileAttributes) (*models.SNMPContextProfile, error) {
	rn := fmt.Sprintf(models.RnsnmpCtxP)
	parentDn := fmt.Sprintf(models.ParentDnsnmpCtxP, tenant, vrf)
	snmpCtxP := models.NewSNMPContextProfile(rn, parentDn, nameAlias, snmpCtxPAttr)
	snmpCtxP.Status = "modified"
	err := sm.Save(snmpCtxP)
	return snmpCtxP, err
}

func (sm *ServiceManager) ListSNMPContextProfile(vrf string, tenant string) ([]*models.SNMPContextProfile, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ctx-%s/snmpCtxP.json", models.BaseurlStr, tenant, vrf)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SNMPContextProfileListFromContainer(cont)
	return list, err
}

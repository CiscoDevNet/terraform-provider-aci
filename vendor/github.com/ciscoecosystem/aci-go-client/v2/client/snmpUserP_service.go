package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateSnmpUserProfile(name string, snmp_policy string, description string, snmpUserPAttr models.SnmpUserProfileAttributes) (*models.SnmpUserProfile, error) {

	rn := fmt.Sprintf(models.RnSnmpUserP, name)

	parentDn := fmt.Sprintf(models.ParentDnSnmpUserP, snmp_policy)
	snmpUserP := models.NewSnmpUserProfile(rn, parentDn, description, snmpUserPAttr)

	err := sm.Save(snmpUserP)
	return snmpUserP, err
}

func (sm *ServiceManager) ReadSnmpUserProfile(name string, snmp_policy string) (*models.SnmpUserProfile, error) {

	rn := fmt.Sprintf(models.RnSnmpUserP, name)

	parentDn := fmt.Sprintf(models.ParentDnSnmpUserP, snmp_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	snmpUserP := models.SnmpUserProfileFromContainer(cont)
	return snmpUserP, nil
}

func (sm *ServiceManager) DeleteSnmpUserProfile(name string, snmp_policy string) error {

	rn := fmt.Sprintf(models.RnSnmpUserP, name)

	parentDn := fmt.Sprintf(models.ParentDnSnmpUserP, snmp_policy)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)

	return sm.DeleteByDn(dn, models.SnmpUserPClassName)
}

func (sm *ServiceManager) UpdateSnmpUserProfile(name string, snmp_policy string, description string, snmpUserPAttr models.SnmpUserProfileAttributes) (*models.SnmpUserProfile, error) {

	rn := fmt.Sprintf(models.RnSnmpUserP, name)

	parentDn := fmt.Sprintf(models.ParentDnSnmpUserP, snmp_policy)
	snmpUserP := models.NewSnmpUserProfile(rn, parentDn, description, snmpUserPAttr)

	snmpUserP.Status = "modified"
	err := sm.Save(snmpUserP)
	return snmpUserP, err
}

func (sm *ServiceManager) ListSnmpUserProfile(snmp_policy string) ([]*models.SnmpUserProfile, error) {

	parentDn := fmt.Sprintf(models.ParentDnSnmpUserP, snmp_policy)
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, models.SnmpUserPClassName)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.SnmpUserProfileListFromContainer(cont)
	return list, err
}

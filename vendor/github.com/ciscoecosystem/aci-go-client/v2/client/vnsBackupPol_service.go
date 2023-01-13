package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreatePBRBackupPolicy(name string, tenant string, description string, nameAlias string, vnsBackupPolAttr models.PBRBackupPolicyAttributes) (*models.PBRBackupPolicy, error) {
	rn := fmt.Sprintf(models.RnvnsBackupPol, name)
	parentDn := fmt.Sprintf(models.ParentDnvnsBackupPol, tenant)
	vnsBackupPol := models.NewPBRBackupPolicy(rn, parentDn, description, nameAlias, vnsBackupPolAttr)
	err := sm.Save(vnsBackupPol)
	return vnsBackupPol, err
}

func (sm *ServiceManager) ReadPBRBackupPolicy(name string, tenant string) (*models.PBRBackupPolicy, error) {
	dn := fmt.Sprintf(models.DnvnsBackupPol, tenant, name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsBackupPol := models.PBRBackupPolicyFromContainer(cont)
	return vnsBackupPol, nil
}

func (sm *ServiceManager) DeletePBRBackupPolicy(name string, tenant string) error {
	dn := fmt.Sprintf(models.DnvnsBackupPol, tenant, name)
	return sm.DeleteByDn(dn, models.VnsbackuppolClassName)
}

func (sm *ServiceManager) UpdatePBRBackupPolicy(name string, tenant string, description string, nameAlias string, vnsBackupPolAttr models.PBRBackupPolicyAttributes) (*models.PBRBackupPolicy, error) {
	rn := fmt.Sprintf(models.RnvnsBackupPol, name)
	parentDn := fmt.Sprintf(models.ParentDnvnsBackupPol, tenant)
	vnsBackupPol := models.NewPBRBackupPolicy(rn, parentDn, description, nameAlias, vnsBackupPolAttr)
	vnsBackupPol.Status = "modified"
	err := sm.Save(vnsBackupPol)
	return vnsBackupPol, err
}

func (sm *ServiceManager) ListPBRBackupPolicy(tenant string) ([]*models.PBRBackupPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/svcCont/vnsBackupPol.json", models.BaseurlStr, tenant)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.PBRBackupPolicyListFromContainer(cont)
	return list, err
}

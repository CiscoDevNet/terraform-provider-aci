package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateEPLoopProtectionPolicy(name string, description string, nameAlias string, epLoopProtectPAttr models.EPLoopProtectionPolicyAttributes) (*models.EPLoopProtectionPolicy, error) {
	rn := fmt.Sprintf(models.RnepLoopProtectP, name)
	parentDn := fmt.Sprintf(models.ParentDnepLoopProtectP)
	epLoopProtectP := models.NewEPLoopProtectionPolicy(rn, parentDn, description, nameAlias, epLoopProtectPAttr)
	err := sm.Save(epLoopProtectP)
	return epLoopProtectP, err
}

func (sm *ServiceManager) ReadEPLoopProtectionPolicy(name string) (*models.EPLoopProtectionPolicy, error) {
	dn := fmt.Sprintf(models.DnepLoopProtectP, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	epLoopProtectP := models.EPLoopProtectionPolicyFromContainer(cont)
	return epLoopProtectP, nil
}

func (sm *ServiceManager) DeleteEPLoopProtectionPolicy(name string) error {
	dn := fmt.Sprintf(models.DnepLoopProtectP, name)
	return sm.DeleteByDn(dn, models.EploopprotectpClassName)
}

func (sm *ServiceManager) UpdateEPLoopProtectionPolicy(name string, description string, nameAlias string, epLoopProtectPAttr models.EPLoopProtectionPolicyAttributes) (*models.EPLoopProtectionPolicy, error) {
	rn := fmt.Sprintf(models.RnepLoopProtectP, name)
	parentDn := fmt.Sprintf(models.ParentDnepLoopProtectP)
	epLoopProtectP := models.NewEPLoopProtectionPolicy(rn, parentDn, description, nameAlias, epLoopProtectPAttr)
	epLoopProtectP.Status = "modified"
	err := sm.Save(epLoopProtectP)
	return epLoopProtectP, err
}

func (sm *ServiceManager) ListEPLoopProtectionPolicy() ([]*models.EPLoopProtectionPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/epLoopProtectP.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.EPLoopProtectionPolicyListFromContainer(cont)
	return list, err
}

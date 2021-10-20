package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateFabricWideSettingsPolicy(description string, nameAlias string, infraSetPolAttr models.FabricWideSettingsPolicyAttributes) (*models.FabricWideSettingsPolicy, error) {
	rn := fmt.Sprintf(models.RninfraSetPol)
	parentDn := fmt.Sprintf(models.ParentDninfraSetPol)
	infraSetPol := models.NewFabricWideSettingsPolicy(rn, parentDn, description, nameAlias, infraSetPolAttr)
	err := sm.Save(infraSetPol)
	return infraSetPol, err
}

func (sm *ServiceManager) ReadFabricWideSettingsPolicy() (*models.FabricWideSettingsPolicy, error) {
	dn := fmt.Sprintf(models.DninfraSetPol)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	infraSetPol := models.FabricWideSettingsPolicyFromContainer(cont)
	return infraSetPol, nil
}

func (sm *ServiceManager) DeleteFabricWideSettingsPolicy() error {
	dn := fmt.Sprintf(models.DninfraSetPol)
	return sm.DeleteByDn(dn, models.InfrasetpolClassName)
}

func (sm *ServiceManager) UpdateFabricWideSettingsPolicy(description string, nameAlias string, infraSetPolAttr models.FabricWideSettingsPolicyAttributes) (*models.FabricWideSettingsPolicy, error) {
	rn := fmt.Sprintf(models.RninfraSetPol)
	parentDn := fmt.Sprintf(models.ParentDninfraSetPol)
	infraSetPol := models.NewFabricWideSettingsPolicy(rn, parentDn, description, nameAlias, infraSetPolAttr)
	infraSetPol.Status = "modified"
	err := sm.Save(infraSetPol)
	return infraSetPol, err
}

func (sm *ServiceManager) ListFabricWideSettingsPolicy() ([]*models.FabricWideSettingsPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/infraSetPol.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.FabricWideSettingsPolicyListFromContainer(cont)
	return list, err
}

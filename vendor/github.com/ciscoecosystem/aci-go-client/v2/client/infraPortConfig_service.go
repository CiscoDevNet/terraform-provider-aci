package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateInfraPortConfiguration(subPort string, port string, card string, node string, description string, infraPortConfigAttr models.InfraPortConfigurationAttributes) (*models.InfraPortConfiguration, error) {
	rn := fmt.Sprintf(models.RnInfraPortConfig, node, card, port, subPort)
	infraPortConfig := models.NewInfraPortConfiguration(rn, models.ParentDnInfraPortConfig, description, infraPortConfigAttr)
	err := sm.Save(infraPortConfig)
	return infraPortConfig, err
}

func (sm *ServiceManager) ReadInfraPortConfiguration(subPort string, port string, card string, node string) (*models.InfraPortConfiguration, error) {
	rn := fmt.Sprintf(models.RnInfraPortConfig, node, card, port, subPort)
	dn := fmt.Sprintf("%s/%s", models.ParentDnInfraPortConfig, rn)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	infraPortConfig := models.InfraPortConfigurationFromContainer(cont)
	return infraPortConfig, nil
}

func (sm *ServiceManager) DeleteInfraPortConfiguration(subPort string, port string, card string, node string) error {
	rn := fmt.Sprintf(models.RnInfraPortConfig, node, card, port, subPort)
	dn := fmt.Sprintf("%s/%s", models.ParentDnInfraPortConfig, rn)
	return sm.DeleteByDn(dn, models.InfraPortConfigClassName)
}

func (sm *ServiceManager) UpdateInfraPortConfiguration(subPort string, port string, card string, node string, description string, infraPortConfigAttr models.InfraPortConfigurationAttributes) (*models.InfraPortConfiguration, error) {
	rn := fmt.Sprintf(models.RnInfraPortConfig, node, card, port, subPort)
	infraPortConfig := models.NewInfraPortConfiguration(rn, models.ParentDnInfraPortConfig, description, infraPortConfigAttr)
	infraPortConfig.Status = "modified"
	err := sm.Save(infraPortConfig)
	return infraPortConfig, err
}

func (sm *ServiceManager) ListInfraPortConfiguration() ([]*models.InfraPortConfiguration, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, models.ParentDnInfraPortConfig, models.InfraPortConfigClassName)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.InfraPortConfigurationListFromContainer(cont)
	return list, err
}

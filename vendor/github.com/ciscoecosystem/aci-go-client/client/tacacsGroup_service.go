package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateTACACSMonitoringDestinationGroup(name string, description string, nameAlias string, tacacsGroupAttr models.TACACSMonitoringDestinationGroupAttributes) (*models.TACACSMonitoringDestinationGroup, error) {
	rn := fmt.Sprintf(models.RntacacsGroup, name)
	parentDn := fmt.Sprintf(models.ParentDntacacsGroup)
	tacacsGroup := models.NewTACACSMonitoringDestinationGroup(rn, parentDn, description, nameAlias, tacacsGroupAttr)
	err := sm.Save(tacacsGroup)
	return tacacsGroup, err
}

func (sm *ServiceManager) ReadTACACSMonitoringDestinationGroup(name string) (*models.TACACSMonitoringDestinationGroup, error) {
	dn := fmt.Sprintf(models.DntacacsGroup, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	tacacsGroup := models.TACACSMonitoringDestinationGroupFromContainer(cont)
	return tacacsGroup, nil
}

func (sm *ServiceManager) DeleteTACACSMonitoringDestinationGroup(name string) error {
	dn := fmt.Sprintf(models.DntacacsGroup, name)
	return sm.DeleteByDn(dn, models.TacacsgroupClassName)
}

func (sm *ServiceManager) UpdateTACACSMonitoringDestinationGroup(name string, description string, nameAlias string, tacacsGroupAttr models.TACACSMonitoringDestinationGroupAttributes) (*models.TACACSMonitoringDestinationGroup, error) {
	rn := fmt.Sprintf(models.RntacacsGroup, name)
	parentDn := fmt.Sprintf(models.ParentDntacacsGroup)
	tacacsGroup := models.NewTACACSMonitoringDestinationGroup(rn, parentDn, description, nameAlias, tacacsGroupAttr)
	tacacsGroup.Status = "modified"
	err := sm.Save(tacacsGroup)
	return tacacsGroup, err
}

func (sm *ServiceManager) ListTACACSMonitoringDestinationGroup() ([]*models.TACACSMonitoringDestinationGroup, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/tacacsGroup.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.TACACSMonitoringDestinationGroupListFromContainer(cont)
	return list, err
}

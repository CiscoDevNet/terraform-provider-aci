package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateTACACSDestination(port string, host string, tacacs_monitoring_destination_group string, description string, nameAlias string, tacacsTacacsDestAttr models.TACACSDestinationAttributes) (*models.TACACSDestination, error) {
	rn := fmt.Sprintf(models.RntacacsTacacsDest, host, port)
	parentDn := fmt.Sprintf(models.ParentDntacacsTacacsDest, tacacs_monitoring_destination_group)
	tacacsTacacsDest := models.NewTACACSDestination(rn, parentDn, description, nameAlias, tacacsTacacsDestAttr)
	err := sm.Save(tacacsTacacsDest)
	return tacacsTacacsDest, err
}

func (sm *ServiceManager) ReadTACACSDestination(port string, host string, tacacs_monitoring_destination_group string) (*models.TACACSDestination, error) {
	dn := fmt.Sprintf(models.DntacacsTacacsDest, tacacs_monitoring_destination_group, host, port)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	tacacsTacacsDest := models.TACACSDestinationFromContainer(cont)
	return tacacsTacacsDest, nil
}

func (sm *ServiceManager) DeleteTACACSDestination(port string, host string, tacacs_monitoring_destination_group string) error {
	dn := fmt.Sprintf(models.DntacacsTacacsDest, tacacs_monitoring_destination_group, host, port)
	return sm.DeleteByDn(dn, models.TacacstacacsdestClassName)
}

func (sm *ServiceManager) UpdateTACACSDestination(port string, host string, tacacs_monitoring_destination_group string, description string, nameAlias string, tacacsTacacsDestAttr models.TACACSDestinationAttributes) (*models.TACACSDestination, error) {
	rn := fmt.Sprintf(models.RntacacsTacacsDest, host, port)
	parentDn := fmt.Sprintf(models.ParentDntacacsTacacsDest, tacacs_monitoring_destination_group)
	tacacsTacacsDest := models.NewTACACSDestination(rn, parentDn, description, nameAlias, tacacsTacacsDestAttr)
	tacacsTacacsDest.Status = "modified"
	err := sm.Save(tacacsTacacsDest)
	return tacacsTacacsDest, err
}

func (sm *ServiceManager) ListTACACSDestination(tacacs_monitoring_destination_group string) ([]*models.TACACSDestination, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/tacacsgroup-%s/tacacsTacacsDest.json", models.BaseurlStr, tacacs_monitoring_destination_group)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.TACACSDestinationListFromContainer(cont)
	return list, err
}

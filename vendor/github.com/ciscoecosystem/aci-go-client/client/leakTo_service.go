package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateTenantandVRFdestinationforInterVRFLeakedRoutes(destinationCtxName string, destinationTenantName string, parentDn string, description string, nameAlias string, leakToAttr models.TenantandVRFdestinationforInterVRFLeakedRoutesAttributes) (*models.TenantandVRFdestinationforInterVRFLeakedRoutes, error) {
	rn := fmt.Sprintf(models.RnleakTo, destinationTenantName, destinationCtxName)
	leakTo := models.NewTenantandVRFdestinationforInterVRFLeakedRoutes(rn, parentDn, description, nameAlias, leakToAttr)
	err := sm.Save(leakTo)
	return leakTo, err
}

func (sm *ServiceManager) ReadTenantandVRFdestinationforInterVRFLeakedRoutes(destinationCtxName string, destinationTenantName string, parentDn string) (*models.TenantandVRFdestinationforInterVRFLeakedRoutes, error) {
	dn := parentDn + "/" + fmt.Sprintf(models.RnleakTo, destinationTenantName, destinationCtxName)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	leakTo := models.TenantandVRFdestinationforInterVRFLeakedRoutesFromContainer(cont)
	return leakTo, nil
}

func (sm *ServiceManager) DeleteTenantandVRFdestinationforInterVRFLeakedRoutes(destinationCtxName string, destinationTenantName string, parentDn string) error {
	dn := parentDn + "/" + fmt.Sprintf(models.RnleakTo, destinationTenantName, destinationCtxName)
	return sm.DeleteByDn(dn, models.LeaktoClassName)
}

func (sm *ServiceManager) UpdateTenantandVRFdestinationforInterVRFLeakedRoutes(destinationCtxName string, destinationTenantName string, parentDn string, description string, nameAlias string, leakToAttr models.TenantandVRFdestinationforInterVRFLeakedRoutesAttributes) (*models.TenantandVRFdestinationforInterVRFLeakedRoutes, error) {
	rn := fmt.Sprintf(models.RnleakTo, destinationTenantName, destinationCtxName)
	leakTo := models.NewTenantandVRFdestinationforInterVRFLeakedRoutes(rn, parentDn, description, nameAlias, leakToAttr)
	leakTo.Status = "modified"
	err := sm.Save(leakTo)
	return leakTo, err
}

func (sm *ServiceManager) ListTenantandVRFdestinationforInterVRFLeakedRoutes(parentDn string) ([]*models.TenantandVRFdestinationforInterVRFLeakedRoutes, error) {
	dnUrl := fmt.Sprintf("%s/%s/leakTo.json", models.BaseurlStr, parentDn)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.TenantandVRFdestinationforInterVRFLeakedRoutesListFromContainer(cont)
	return list, err
}

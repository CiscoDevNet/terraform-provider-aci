package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateInterVRFLeakedRoutesContainer(vrf string, tenant string, description string, nameAlias string, leakRoutesAttr models.InterVRFLeakedRoutesContainerAttributes) (*models.InterVRFLeakedRoutesContainer, error) {
	rn := fmt.Sprintf(models.RnleakRoutes)
	parentDn := fmt.Sprintf(models.ParentDnleakRoutes, tenant, vrf)
	leakRoutes := models.NewInterVRFLeakedRoutesContainer(rn, parentDn, description, nameAlias, leakRoutesAttr)
	err := sm.Save(leakRoutes)
	return leakRoutes, err
}

func (sm *ServiceManager) ReadInterVRFLeakedRoutesContainer(vrf string, tenant string) (*models.InterVRFLeakedRoutesContainer, error) {
	dn := fmt.Sprintf(models.DnleakRoutes, tenant, vrf)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	leakRoutes := models.InterVRFLeakedRoutesContainerFromContainer(cont)
	return leakRoutes, nil
}

func (sm *ServiceManager) DeleteInterVRFLeakedRoutesContainer(vrf string, tenant string) error {
	dn := fmt.Sprintf(models.DnleakRoutes, tenant, vrf)
	return sm.DeleteByDn(dn, models.LeakroutesClassName)
}

func (sm *ServiceManager) UpdateInterVRFLeakedRoutesContainer(vrf string, tenant string, description string, nameAlias string, leakRoutesAttr models.InterVRFLeakedRoutesContainerAttributes) (*models.InterVRFLeakedRoutesContainer, error) {
	rn := fmt.Sprintf(models.RnleakRoutes)
	parentDn := fmt.Sprintf(models.ParentDnleakRoutes, tenant, vrf)
	leakRoutes := models.NewInterVRFLeakedRoutesContainer(rn, parentDn, description, nameAlias, leakRoutesAttr)
	leakRoutes.Status = "modified"
	err := sm.Save(leakRoutes)
	return leakRoutes, err
}

func (sm *ServiceManager) ListInterVRFLeakedRoutesContainer(vrf string, tenant string) ([]*models.InterVRFLeakedRoutesContainer, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ctx-%s/leakRoutes.json", models.BaseurlStr, tenant, vrf)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.InterVRFLeakedRoutesContainerListFromContainer(cont)
	return list, err
}

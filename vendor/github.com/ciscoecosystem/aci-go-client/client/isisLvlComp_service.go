package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateISISLevel(isis_level_type string, isis_domain_policy string, description string, nameAlias string, isisLvlCompAttr models.ISISLevelAttributes) (*models.ISISLevel, error) {
	rn := fmt.Sprintf(models.RnisisLvlComp, isis_level_type)
	parentDn := fmt.Sprintf(models.ParentDnisisLvlComp, isis_domain_policy)
	isisLvlComp := models.NewISISLevel(rn, parentDn, description, nameAlias, isisLvlCompAttr)
	err := sm.Save(isisLvlComp)
	return isisLvlComp, err
}

func (sm *ServiceManager) ReadISISLevel(isis_level_type string, isis_domain_policy string) (*models.ISISLevel, error) {
	dn := fmt.Sprintf(models.DnisisLvlComp, isis_domain_policy, isis_level_type)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	isisLvlComp := models.ISISLevelFromContainer(cont)
	return isisLvlComp, nil
}

func (sm *ServiceManager) DeleteISISLevel(isis_level_type string, isis_domain_policy string) error {
	dn := fmt.Sprintf(models.DnisisLvlComp, isis_domain_policy, isis_level_type)
	return sm.DeleteByDn(dn, models.IsislvlcompClassName)
}

func (sm *ServiceManager) UpdateISISLevel(isis_level_type string, isis_domain_policy string, description string, nameAlias string, isisLvlCompAttr models.ISISLevelAttributes) (*models.ISISLevel, error) {
	rn := fmt.Sprintf(models.RnisisLvlComp, isis_level_type)
	parentDn := fmt.Sprintf(models.ParentDnisisLvlComp, isis_domain_policy)
	isisLvlComp := models.NewISISLevel(rn, parentDn, description, nameAlias, isisLvlCompAttr)
	isisLvlComp.Status = "modified"
	err := sm.Save(isisLvlComp)
	return isisLvlComp, err
}

func (sm *ServiceManager) ListISISLevel(isis_domain_policy string) ([]*models.ISISLevel, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/isisDomP-%s/isisLvlComp.json", models.BaseurlStr, isis_domain_policy)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.ISISLevelListFromContainer(cont)
	return list, err
}

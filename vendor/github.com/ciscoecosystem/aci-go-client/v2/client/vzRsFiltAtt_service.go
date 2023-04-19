package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateFilterRelationship(tnVzFilterName string, parentDn string, vzRsFiltAttAttr models.FilterRelationshipAttributes) (*models.FilterRelationship, error) {
	rn := fmt.Sprintf(models.RnvzRsFiltAtt, tnVzFilterName)
	vzRsFiltAtt := models.NewFilterRelationship(rn, parentDn, vzRsFiltAttAttr)
	err := sm.Save(vzRsFiltAtt)
	return vzRsFiltAtt, err
}

func (sm *ServiceManager) ReadFilterRelationship(tnVzFilterName string, parentDn string) (*models.FilterRelationship, error) {
	dn := fmt.Sprintf(models.DnvzRsFiltAtt, parentDn, tnVzFilterName)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vzRsFiltAtt := models.FilterRelationshipFromContainer(cont)
	return vzRsFiltAtt, nil
}

func (sm *ServiceManager) DeleteFilterRelationship(tnVzFilterName string, parentDn string) error {
	dn := fmt.Sprintf(models.DnvzRsFiltAtt, parentDn, tnVzFilterName)
	return sm.DeleteByDn(dn, models.VzrsfiltattClassName)
}

func (sm *ServiceManager) UpdateFilterRelationship(tnVzFilterName string, parentDn string, vzRsFiltAttAttr models.FilterRelationshipAttributes) (*models.FilterRelationship, error) {
	rn := fmt.Sprintf(models.RnvzRsFiltAtt, tnVzFilterName)
	vzRsFiltAtt := models.NewFilterRelationship(rn, parentDn, vzRsFiltAttAttr)
	vzRsFiltAtt.Status = "modified"
	err := sm.Save(vzRsFiltAtt)
	return vzRsFiltAtt, err
}

func (sm *ServiceManager) ListFilterRelationship(parentDn string) ([]*models.FilterRelationship, error) {
	dnUrl := fmt.Sprintf("%s/%s/vzRsFiltAtt.json", models.BaseurlStr, parentDn)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.FilterRelationshipListFromContainer(cont)
	return list, err
}

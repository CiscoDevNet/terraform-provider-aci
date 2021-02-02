package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL4L7ServiceGraphTemplate(name string, tenant string, description string, vnsAbsGraphattr models.L4L7ServiceGraphTemplateAttributes) (*models.L4L7ServiceGraphTemplate, error) {
	rn := fmt.Sprintf("AbsGraph-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vnsAbsGraph := models.NewL4L7ServiceGraphTemplate(rn, parentDn, description, vnsAbsGraphattr)
	err := sm.Save(vnsAbsGraph)
	return vnsAbsGraph, err
}

func (sm *ServiceManager) ReadL4L7ServiceGraphTemplate(name string, tenant string) (*models.L4L7ServiceGraphTemplate, error) {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsGraph := models.L4L7ServiceGraphTemplateFromContainer(cont)
	return vnsAbsGraph, nil
}

func (sm *ServiceManager) DeleteL4L7ServiceGraphTemplate(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/AbsGraph-%s", tenant, name)
	return sm.DeleteByDn(dn, models.VnsabsgraphClassName)
}

func (sm *ServiceManager) UpdateL4L7ServiceGraphTemplate(name string, tenant string, description string, vnsAbsGraphattr models.L4L7ServiceGraphTemplateAttributes) (*models.L4L7ServiceGraphTemplate, error) {
	rn := fmt.Sprintf("AbsGraph-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	vnsAbsGraph := models.NewL4L7ServiceGraphTemplate(rn, parentDn, description, vnsAbsGraphattr)

	vnsAbsGraph.Status = "modified"
	err := sm.Save(vnsAbsGraph)
	return vnsAbsGraph, err

}

func (sm *ServiceManager) ListL4L7ServiceGraphTemplate(tenant string) ([]*models.L4L7ServiceGraphTemplate, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/vnsAbsGraph.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L4L7ServiceGraphTemplateListFromContainer(cont)

	return list, err
}

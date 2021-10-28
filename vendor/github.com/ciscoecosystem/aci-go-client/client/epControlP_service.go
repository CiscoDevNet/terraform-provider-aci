package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateEndpointControlPolicy(name string, description string, nameAlias string, epControlPAttr models.EndpointControlPolicyAttributes) (*models.EndpointControlPolicy, error) {
	rn := fmt.Sprintf(models.RnepControlP, name)
	parentDn := fmt.Sprintf(models.ParentDnepControlP)
	epControlP := models.NewEndpointControlPolicy(rn, parentDn, description, nameAlias, epControlPAttr)
	err := sm.Save(epControlP)
	return epControlP, err
}

func (sm *ServiceManager) ReadEndpointControlPolicy(name string) (*models.EndpointControlPolicy, error) {
	dn := fmt.Sprintf(models.DnepControlP, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	epControlP := models.EndpointControlPolicyFromContainer(cont)
	return epControlP, nil
}

func (sm *ServiceManager) DeleteEndpointControlPolicy(name string) error {
	dn := fmt.Sprintf(models.DnepControlP, name)
	return sm.DeleteByDn(dn, models.EpcontrolpClassName)
}

func (sm *ServiceManager) UpdateEndpointControlPolicy(name string, description string, nameAlias string, epControlPAttr models.EndpointControlPolicyAttributes) (*models.EndpointControlPolicy, error) {
	rn := fmt.Sprintf(models.RnepControlP, name)
	parentDn := fmt.Sprintf(models.ParentDnepControlP)
	epControlP := models.NewEndpointControlPolicy(rn, parentDn, description, nameAlias, epControlPAttr)
	epControlP.Status = "modified"
	err := sm.Save(epControlP)
	return epControlP, err
}

func (sm *ServiceManager) ListEndpointControlPolicy() ([]*models.EndpointControlPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/epControlP.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.EndpointControlPolicyListFromContainer(cont)
	return list, err
}

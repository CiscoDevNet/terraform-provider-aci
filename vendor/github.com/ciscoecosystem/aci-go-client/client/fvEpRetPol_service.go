package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateEndPointRetentionPolicy(name string, tenant string, description string, fvEpRetPolattr models.EndPointRetentionPolicyAttributes) (*models.EndPointRetentionPolicy, error) {
	rn := fmt.Sprintf("epRPol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	fvEpRetPol := models.NewEndPointRetentionPolicy(rn, parentDn, description, fvEpRetPolattr)
	err := sm.Save(fvEpRetPol)
	return fvEpRetPol, err
}

func (sm *ServiceManager) ReadEndPointRetentionPolicy(name string, tenant string) (*models.EndPointRetentionPolicy, error) {
	dn := fmt.Sprintf("uni/tn-%s/epRPol-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvEpRetPol := models.EndPointRetentionPolicyFromContainer(cont)
	return fvEpRetPol, nil
}

func (sm *ServiceManager) DeleteEndPointRetentionPolicy(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/epRPol-%s", tenant, name)
	return sm.DeleteByDn(dn, models.FvepretpolClassName)
}

func (sm *ServiceManager) UpdateEndPointRetentionPolicy(name string, tenant string, description string, fvEpRetPolattr models.EndPointRetentionPolicyAttributes) (*models.EndPointRetentionPolicy, error) {
	rn := fmt.Sprintf("epRPol-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	fvEpRetPol := models.NewEndPointRetentionPolicy(rn, parentDn, description, fvEpRetPolattr)

	fvEpRetPol.Status = "modified"
	err := sm.Save(fvEpRetPol)
	return fvEpRetPol, err

}

func (sm *ServiceManager) ListEndPointRetentionPolicy(tenant string) ([]*models.EndPointRetentionPolicy, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/fvEpRetPol.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.EndPointRetentionPolicyListFromContainer(cont)

	return list, err
}

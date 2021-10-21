package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateQOSInstancePolicy(name string, description string, nameAlias string, qosInstPolAttr models.QOSInstancePolicyAttributes) (*models.QOSInstancePolicy, error) {
	rn := fmt.Sprintf(models.RnqosInstPol, name)
	parentDn := fmt.Sprintf(models.ParentDnqosInstPol)
	qosInstPol := models.NewQOSInstancePolicy(rn, parentDn, description, nameAlias, qosInstPolAttr)
	err := sm.Save(qosInstPol)
	return qosInstPol, err
}

func (sm *ServiceManager) ReadQOSInstancePolicy(name string) (*models.QOSInstancePolicy, error) {
	dn := fmt.Sprintf(models.DnqosInstPol, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	qosInstPol := models.QOSInstancePolicyFromContainer(cont)
	return qosInstPol, nil
}

func (sm *ServiceManager) DeleteQOSInstancePolicy(name string) error {
	dn := fmt.Sprintf(models.DnqosInstPol, name)
	return sm.DeleteByDn(dn, models.QosinstpolClassName)
}

func (sm *ServiceManager) UpdateQOSInstancePolicy(name string, description string, nameAlias string, qosInstPolAttr models.QOSInstancePolicyAttributes) (*models.QOSInstancePolicy, error) {
	rn := fmt.Sprintf(models.RnqosInstPol, name)
	parentDn := fmt.Sprintf(models.ParentDnqosInstPol)
	qosInstPol := models.NewQOSInstancePolicy(rn, parentDn, description, nameAlias, qosInstPolAttr)
	qosInstPol.Status = "modified"
	err := sm.Save(qosInstPol)
	return qosInstPol, err
}

func (sm *ServiceManager) ListQOSInstancePolicy() ([]*models.QOSInstancePolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/qosInstPol.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.QOSInstancePolicyListFromContainer(cont)
	return list, err
}

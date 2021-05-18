package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateVPCDomainPolicy(name string, description string, nameAlias string, vpcInstPolAttr models.VPCDomainPolicyAttributes) (*models.VPCDomainPolicy, error) {
	rn := fmt.Sprintf(models.RnvpcInstPol, name)
	parentDn := fmt.Sprintf(models.ParentDnvpcInstPol)
	vpcInstPol := models.NewVPCDomainPolicy(rn, parentDn, description, nameAlias, vpcInstPolAttr)
	err := sm.Save(vpcInstPol)
	return vpcInstPol, err
}

func (sm *ServiceManager) ReadVPCDomainPolicy(name string) (*models.VPCDomainPolicy, error) {
	dn := fmt.Sprintf(models.DnvpcInstPol, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	vpcInstPol := models.VPCDomainPolicyFromContainer(cont)
	return vpcInstPol, nil
}

func (sm *ServiceManager) DeleteVPCDomainPolicy(name string) error {
	dn := fmt.Sprintf(models.DnvpcInstPol, name)
	return sm.DeleteByDn(dn, models.VpcinstpolClassName)
}

func (sm *ServiceManager) UpdateVPCDomainPolicy(name string, description string, nameAlias string, vpcInstPolAttr models.VPCDomainPolicyAttributes) (*models.VPCDomainPolicy, error) {
	rn := fmt.Sprintf(models.RnvpcInstPol, name)
	parentDn := fmt.Sprintf(models.ParentDnvpcInstPol)
	vpcInstPol := models.NewVPCDomainPolicy(rn, parentDn, description, nameAlias, vpcInstPolAttr)
	vpcInstPol.Status = "modified"
	err := sm.Save(vpcInstPol)
	return vpcInstPol, err
}

func (sm *ServiceManager) ListVPCDomainPolicy() ([]*models.VPCDomainPolicy, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/vpcInstPol.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.VPCDomainPolicyListFromContainer(cont)
	return list, err
}

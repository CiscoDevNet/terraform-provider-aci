package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateBGPAddressFamilyContextPolicyRelationship(af string, tnBgpCtxAfPolName string, vrf string, tenant string, fvRsCtxToBgpCtxAfPolAttr models.BGPAddressFamilyContextPolicyRelationshipAttributes) (*models.BGPAddressFamilyContextPolicyRelationship, error) {
	rn := fmt.Sprintf(models.RnfvRsCtxToBgpCtxAfPol, tnBgpCtxAfPolName, af)
	parentDn := fmt.Sprintf(models.ParentDnfvRsCtxToBgpCtxAfPol, tenant, vrf)
	fvRsCtxToBgpCtxAfPol := models.NewBGPAddressFamilyContextPolicyRelationship(rn, parentDn, fvRsCtxToBgpCtxAfPolAttr)
	err := sm.Save(fvRsCtxToBgpCtxAfPol)
	return fvRsCtxToBgpCtxAfPol, err
}

func (sm *ServiceManager) ReadBGPAddressFamilyContextPolicyRelationship(af string, tnBgpCtxAfPolName string, vrf string, tenant string) (*models.BGPAddressFamilyContextPolicyRelationship, error) {
	dn := fmt.Sprintf(models.DnfvRsCtxToBgpCtxAfPol, tenant, vrf, tnBgpCtxAfPolName, af)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvRsCtxToBgpCtxAfPol := models.BGPAddressFamilyContextPolicyRelationshipFromContainer(cont)
	return fvRsCtxToBgpCtxAfPol, nil
}

func (sm *ServiceManager) DeleteBGPAddressFamilyContextPolicyRelationship(af string, tnBgpCtxAfPolName string, vrf string, tenant string) error {
	dn := fmt.Sprintf(models.DnfvRsCtxToBgpCtxAfPol, tenant, vrf, tnBgpCtxAfPolName, af)
	return sm.DeleteByDn(dn, models.FvrsctxtobgpctxafpolClassName)
}

func (sm *ServiceManager) UpdateBGPAddressFamilyContextPolicyRelationship(af string, tnBgpCtxAfPolName string, vrf string, tenant string, fvRsCtxToBgpCtxAfPolAttr models.BGPAddressFamilyContextPolicyRelationshipAttributes) (*models.BGPAddressFamilyContextPolicyRelationship, error) {
	rn := fmt.Sprintf(models.RnfvRsCtxToBgpCtxAfPol, tnBgpCtxAfPolName, af)
	parentDn := fmt.Sprintf(models.ParentDnfvRsCtxToBgpCtxAfPol, tenant, vrf)
	fvRsCtxToBgpCtxAfPol := models.NewBGPAddressFamilyContextPolicyRelationship(rn, parentDn, fvRsCtxToBgpCtxAfPolAttr)
	fvRsCtxToBgpCtxAfPol.Status = "modified"
	err := sm.Save(fvRsCtxToBgpCtxAfPol)
	return fvRsCtxToBgpCtxAfPol, err
}

func (sm *ServiceManager) ListBGPAddressFamilyContextPolicyRelationship(vrf string, tenant string) ([]*models.BGPAddressFamilyContextPolicyRelationship, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ctx-%s/fvRsCtxToBgpCtxAfPol.json", models.BaseurlStr, tenant, vrf)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.BGPAddressFamilyContextPolicyRelationshipListFromContainer(cont)
	return list, err
}

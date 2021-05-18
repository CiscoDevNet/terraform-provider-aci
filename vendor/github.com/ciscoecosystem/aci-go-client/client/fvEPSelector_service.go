package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateEndpointSecurityGroupSelector(matchExpression string, endpoint_security_group string, application_profile string, tenant string, description string, nameAlias string, fvEPSelectorAttr models.EndpointSecurityGroupSelectorAttributes) (*models.EndpointSecurityGroupSelector, error) {
	rn := fmt.Sprintf(models.RnfvEPSelector, matchExpression)
	parentDn := fmt.Sprintf(models.ParentDnfvEPSelector, tenant, application_profile, endpoint_security_group)
	fvEPSelector := models.NewEndpointSecurityGroupSelector(rn, parentDn, description, nameAlias, fvEPSelectorAttr)
	err := sm.Save(fvEPSelector)
	return fvEPSelector, err
}

func (sm *ServiceManager) ReadEndpointSecurityGroupSelector(matchExpression string, endpoint_security_group string, application_profile string, tenant string) (*models.EndpointSecurityGroupSelector, error) {
	dn := fmt.Sprintf(models.DnfvEPSelector, tenant, application_profile, endpoint_security_group, matchExpression)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	fvEPSelector := models.EndpointSecurityGroupSelectorFromContainer(cont)
	return fvEPSelector, nil
}

func (sm *ServiceManager) DeleteEndpointSecurityGroupSelector(matchExpression string, endpoint_security_group string, application_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnfvEPSelector, tenant, application_profile, endpoint_security_group, matchExpression)
	return sm.DeleteByDn(dn, models.FvepselectorClassName)
}

func (sm *ServiceManager) UpdateEndpointSecurityGroupSelector(matchExpression string, endpoint_security_group string, application_profile string, tenant string, description string, nameAlias string, fvEPSelectorAttr models.EndpointSecurityGroupSelectorAttributes) (*models.EndpointSecurityGroupSelector, error) {
	rn := fmt.Sprintf(models.RnfvEPSelector, matchExpression)
	parentDn := fmt.Sprintf(models.ParentDnfvEPSelector, tenant, application_profile, endpoint_security_group)
	fvEPSelector := models.NewEndpointSecurityGroupSelector(rn, parentDn, description, nameAlias, fvEPSelectorAttr)
	fvEPSelector.Status = "modified"
	err := sm.Save(fvEPSelector)
	return fvEPSelector, err
}

func (sm *ServiceManager) ListEndpointSecurityGroupSelector(endpoint_security_group string, application_profile string, tenant string) ([]*models.EndpointSecurityGroupSelector, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/esg-%s/fvEPSelector.json", models.BaseurlStr, tenant, application_profile, endpoint_security_group)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.EndpointSecurityGroupSelectorListFromContainer(cont)
	return list, err
}

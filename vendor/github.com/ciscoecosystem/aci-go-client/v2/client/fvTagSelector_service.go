package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateEndpointSecurityGroupTagSelector(matchValue string, matchKey string, endpoint_security_group string, application_profile string, tenant string, description string, nameAlias string, fvTagSelectorAttr models.EndpointSecurityGroupTagSelectorAttributes) (*models.EndpointSecurityGroupTagSelector, error) {
	rn := fmt.Sprintf(models.RnfvTagSelector, matchKey, matchValue)
	parentDn := fmt.Sprintf(models.ParentDnfvTagSelector, tenant, application_profile, endpoint_security_group)
	fvTagSelector := models.NewEndpointSecurityGroupTagSelector(rn, parentDn, description, nameAlias, fvTagSelectorAttr)
	err := sm.Save(fvTagSelector)
	return fvTagSelector, err
}

func (sm *ServiceManager) ReadEndpointSecurityGroupTagSelector(matchValue string, matchKey string, endpoint_security_group string, application_profile string, tenant string) (*models.EndpointSecurityGroupTagSelector, error) {
	dn := fmt.Sprintf(models.DnfvTagSelector, tenant, application_profile, endpoint_security_group, matchKey, matchValue)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	fvTagSelector := models.EndpointSecurityGroupTagSelectorFromContainer(cont)
	return fvTagSelector, nil
}

func (sm *ServiceManager) DeleteEndpointSecurityGroupTagSelector(matchValue string, matchKey string, endpoint_security_group string, application_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnfvTagSelector, tenant, application_profile, endpoint_security_group, matchKey, matchValue)
	return sm.DeleteByDn(dn, models.FvtagselectorClassName)
}

func (sm *ServiceManager) UpdateEndpointSecurityGroupTagSelector(matchValue string, matchKey string, endpoint_security_group string, application_profile string, tenant string, description string, nameAlias string, fvTagSelectorAttr models.EndpointSecurityGroupTagSelectorAttributes) (*models.EndpointSecurityGroupTagSelector, error) {
	rn := fmt.Sprintf(models.RnfvTagSelector, matchKey, matchValue)
	parentDn := fmt.Sprintf(models.ParentDnfvTagSelector, tenant, application_profile, endpoint_security_group)
	fvTagSelector := models.NewEndpointSecurityGroupTagSelector(rn, parentDn, description, nameAlias, fvTagSelectorAttr)
	fvTagSelector.Status = "modified"
	err := sm.Save(fvTagSelector)
	return fvTagSelector, err
}

func (sm *ServiceManager) ListEndpointSecurityGroupTagSelector(endpoint_security_group string, application_profile string, tenant string) ([]*models.EndpointSecurityGroupTagSelector, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/esg-%s/fvTagSelector.json", models.BaseurlStr, tenant, application_profile, endpoint_security_group)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.EndpointSecurityGroupTagSelectorListFromContainer(cont)
	return list, err
}

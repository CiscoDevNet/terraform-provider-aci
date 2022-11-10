package client

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateEndpointSecurityGroupEPgSelector(matchEpgDn string, endpoint_security_group string, application_profile string, tenant string, description string, nameAlias string, fvEPgSelectorAttr models.EndpointSecurityGroupEPgSelectorAttributes) (*models.EndpointSecurityGroupEPgSelector, error) {
	rn := fmt.Sprintf(models.RnfvEPgSelector, matchEpgDn)
	parentDn := fmt.Sprintf(models.ParentDnfvEPgSelector, tenant, application_profile, endpoint_security_group)
	fvEPgSelector := models.NewEndpointSecurityGroupEPgSelector(rn, parentDn, description, nameAlias, fvEPgSelectorAttr)
	err := sm.Save(fvEPgSelector)
	return fvEPgSelector, err
}

func (sm *ServiceManager) ReadEndpointSecurityGroupEPgSelector(matchEpgDn string, endpoint_security_group string, application_profile string, tenant string) (*models.EndpointSecurityGroupEPgSelector, error) {
	dn := fmt.Sprintf(models.DnfvEPgSelector, tenant, application_profile, endpoint_security_group, matchEpgDn)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	fvEPgSelector := models.EndpointSecurityGroupEPgSelectorFromContainer(cont)
	return fvEPgSelector, nil
}

func (sm *ServiceManager) DeleteEndpointSecurityGroupEPgSelector(matchEpgDn string, endpoint_security_group string, application_profile string, tenant string) error {
	dn := fmt.Sprintf(models.DnfvEPgSelector, tenant, application_profile, endpoint_security_group, matchEpgDn)
	return sm.DeleteByDn(dn, models.FvepgselectorClassName)
}

func (sm *ServiceManager) UpdateEndpointSecurityGroupEPgSelector(matchEpgDn string, endpoint_security_group string, application_profile string, tenant string, description string, nameAlias string, fvEPgSelectorAttr models.EndpointSecurityGroupEPgSelectorAttributes) (*models.EndpointSecurityGroupEPgSelector, error) {
	rn := fmt.Sprintf(models.RnfvEPgSelector, matchEpgDn)
	parentDn := fmt.Sprintf(models.ParentDnfvEPgSelector, tenant, application_profile, endpoint_security_group)
	fvEPgSelector := models.NewEndpointSecurityGroupEPgSelector(rn, parentDn, description, nameAlias, fvEPgSelectorAttr)
	fvEPgSelector.Status = "modified"
	err := sm.Save(fvEPgSelector)
	return fvEPgSelector, err
}

func (sm *ServiceManager) ListEndpointSecurityGroupEPgSelector(endpoint_security_group string, application_profile string, tenant string) ([]*models.EndpointSecurityGroupEPgSelector, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/ap-%s/esg-%s/fvEPgSelector.json", models.BaseurlStr, tenant, application_profile, endpoint_security_group)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.EndpointSecurityGroupEPgSelectorListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationfvRsMatchEPg(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsmatchEPg", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "fvRsMatchEPg", dn, annotation, tDn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return err
	}
	log.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationfvRsMatchEPg(parentDn string) error {
	dn := fmt.Sprintf("%s/rsmatchEPg", parentDn)
	return sm.DeleteByDn(dn, "fvRsMatchEPg")
}

func (sm *ServiceManager) ReadRelationfvRsMatchEPg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "fvRsMatchEPg")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "fvRsMatchEPg")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

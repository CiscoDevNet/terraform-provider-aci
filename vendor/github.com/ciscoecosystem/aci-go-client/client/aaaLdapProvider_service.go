package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateLDAPProvider(name string, description string, nameAlias string, aaaLdapProviderAttr models.LDAPProviderAttributes) (*models.LDAPProvider, error) {
	rn := fmt.Sprintf(models.RnaaaLdapProvider, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaLdapProvider)
	aaaLdapProvider := models.NewLDAPProvider(rn, parentDn, description, nameAlias, aaaLdapProviderAttr)
	err := sm.Save(aaaLdapProvider)
	return aaaLdapProvider, err
}

func (sm *ServiceManager) ReadLDAPProvider(name string) (*models.LDAPProvider, error) {
	dn := fmt.Sprintf(models.DnaaaLdapProvider, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaLdapProvider := models.LDAPProviderFromContainer(cont)
	return aaaLdapProvider, nil
}

func (sm *ServiceManager) DeleteLDAPProvider(name string) error {
	dn := fmt.Sprintf(models.DnaaaLdapProvider, name)
	return sm.DeleteByDn(dn, models.AaaldapproviderClassName)
}

func (sm *ServiceManager) UpdateLDAPProvider(name string, description string, nameAlias string, aaaLdapProviderAttr models.LDAPProviderAttributes) (*models.LDAPProvider, error) {
	rn := fmt.Sprintf(models.RnaaaLdapProvider, name)
	parentDn := fmt.Sprintf(models.ParentDnaaaLdapProvider)
	aaaLdapProvider := models.NewLDAPProvider(rn, parentDn, description, nameAlias, aaaLdapProviderAttr)
	aaaLdapProvider.Status = "modified"
	err := sm.Save(aaaLdapProvider)
	return aaaLdapProvider, err
}

func (sm *ServiceManager) ListLDAPProvider() ([]*models.LDAPProvider, error) {
	dnUrl := fmt.Sprintf("%s/uni/userext/duoext/aaaLdapProvider.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.LDAPProviderListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationaaaRsProvToEpp(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsProvToEpp", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "aaaRsProvToEpp", dn, annotation, tDn))

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
	fmt.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationaaaRsProvToEpp(parentDn string) error {
	dn := fmt.Sprintf("%s/rsProvToEpp", parentDn)
	return sm.DeleteByDn(dn, "aaaRsProvToEpp")
}

func (sm *ServiceManager) ReadRelationaaaRsProvToEpp(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "aaaRsProvToEpp")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "aaaRsProvToEpp")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationaaaRsSecProvToEpg(parentDn, annotation, tDn string) error {
	dn := fmt.Sprintf("%s/rsSecProvToEpg", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tDn": "%s"
			}
		}
	}`, "aaaRsSecProvToEpg", dn, annotation, tDn))

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
	fmt.Printf("%+v", cont)
	return nil
}

func (sm *ServiceManager) DeleteRelationaaaRsSecProvToEpg(parentDn string) error {
	dn := fmt.Sprintf("%s/rsSecProvToEpg", parentDn)
	return sm.DeleteByDn(dn, "aaaRsSecProvToEpg")
}

func (sm *ServiceManager) ReadRelationaaaRsSecProvToEpg(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "aaaRsSecProvToEpg")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "aaaRsSecProvToEpg")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}
}

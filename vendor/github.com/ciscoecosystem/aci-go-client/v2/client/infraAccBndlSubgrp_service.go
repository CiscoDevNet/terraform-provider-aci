package client



import (
	
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	
)


func (sm *ServiceManager) CreateOverridePCVPCPolicyGroup(name string, leaf_access_bundle_policy_group string, description string, infraAccBndlSubgrpAttr models.OverridePCVPCPolicyGroupAttributes) (*models.OverridePCVPCPolicyGroup, error) {	
	rn := fmt.Sprintf(models.RninfraAccBndlSubgrp , name)
	parentDn := fmt.Sprintf(models.ParentDninfraAccBndlSubgrp, leaf_access_bundle_policy_group )
	infraAccBndlSubgrp := models.NewOverridePCVPCPolicyGroup(rn, parentDn, description, infraAccBndlSubgrpAttr)
	err := sm.Save(infraAccBndlSubgrp)
	return infraAccBndlSubgrp, err
}

func (sm *ServiceManager) ReadOverridePCVPCPolicyGroup(name string, leaf_access_bundle_policy_group string, ) (*models.OverridePCVPCPolicyGroup, error) {
	dn := fmt.Sprintf(models.DninfraAccBndlSubgrp, leaf_access_bundle_policy_group,name)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraAccBndlSubgrp := models.OverridePCVPCPolicyGroupFromContainer(cont)
	return infraAccBndlSubgrp, nil
}

func (sm *ServiceManager) DeleteOverridePCVPCPolicyGroup(name string, leaf_access_bundle_policy_group string, ) error {
	dn := fmt.Sprintf(models.DninfraAccBndlSubgrp, leaf_access_bundle_policy_group,name)
	return sm.DeleteByDn(dn, models.InfraaccbndlsubgrpClassName)
}

func (sm *ServiceManager) UpdateOverridePCVPCPolicyGroup(name string, leaf_access_bundle_policy_group string, description string, infraAccBndlSubgrpAttr models.OverridePCVPCPolicyGroupAttributes) (*models.OverridePCVPCPolicyGroup, error) {
	rn := fmt.Sprintf(models.RninfraAccBndlSubgrp , name)
	parentDn := fmt.Sprintf(models.ParentDninfraAccBndlSubgrp, leaf_access_bundle_policy_group )
	infraAccBndlSubgrp := models.NewOverridePCVPCPolicyGroup(rn, parentDn, description, infraAccBndlSubgrpAttr)
    infraAccBndlSubgrp.Status = "modified"
	err := sm.Save(infraAccBndlSubgrp)
	return infraAccBndlSubgrp, err
}

func (sm *ServiceManager) ListOverridePCVPCPolicyGroup(leaf_access_bundle_policy_group string ) ([]*models.OverridePCVPCPolicyGroup, error) {	
	dnUrl := fmt.Sprintf("%s/uni/infra/funcprof/accbundle-%s/infraAccBndlSubgrp.json", models.BaseurlStr, leaf_access_bundle_policy_group )
    cont, err := sm.GetViaURL(dnUrl)
	list := models.OverridePCVPCPolicyGroupListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsLacpIfPol(parentDn, annotation, tnLacpIfPolName string) error {
	dn := fmt.Sprintf("%s/rslacpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnLacpIfPolName": "%s"
			}
		}
	}`, "infraRsLacpIfPol", dn, annotation, tnLacpIfPolName))

	
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

func (sm *ServiceManager) DeleteRelationinfraRsLacpIfPol(parentDn string) error{
	dn := fmt.Sprintf("%s/rslacpIfPol", parentDn)
	return sm.DeleteByDn(dn , "infraRsLacpIfPol")
}

func (sm *ServiceManager) ReadRelationinfraRsLacpIfPol(parentDn string) (interface{},error) {	
	dnUrl := fmt.Sprintf("%s/%s/%s.json",models.BaseurlStr,parentDn,"infraRsLacpIfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont,"infraRsLacpIfPol")
	
		if len(contList) > 0 {
		dat := models.G(contList[0], "tnLacpIfPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsLacpInterfacePol(parentDn, annotation, tnLacpIfPolName string) error {
	dn := fmt.Sprintf("%s/rslacpInterfacePol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnLacpIfPolName": "%s"
			}
		}
	}`, "infraRsLacpInterfacePol", dn, annotation, tnLacpIfPolName))

	
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

func (sm *ServiceManager) DeleteRelationinfraRsLacpInterfacePol(parentDn string) error{
	dn := fmt.Sprintf("%s/rslacpInterfacePol", parentDn)
	return sm.DeleteByDn(dn , "infraRsLacpInterfacePol")
}

func (sm *ServiceManager) ReadRelationinfraRsLacpInterfacePol(parentDn string) (interface{},error) {	
	dnUrl := fmt.Sprintf("%s/%s/%s.json",models.BaseurlStr,parentDn,"infraRsLacpInterfacePol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont,"infraRsLacpInterfacePol")
	
		if len(contList) > 0 {
		dat := models.G(contList[0], "tnLacpIfPolName")
		return dat, err
	} else {
		return nil, err
	}
}


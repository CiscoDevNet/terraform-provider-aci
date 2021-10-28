package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateAccessSwitchPolicyGroup(name string, description string, nameAlias string, infraAccNodePGrpAttr models.AccessSwitchPolicyGroupAttributes) (*models.AccessSwitchPolicyGroup, error) {
	rn := fmt.Sprintf(models.RninfraAccNodePGrp, name)
	parentDn := fmt.Sprintf(models.ParentDninfraAccNodePGrp)
	infraAccNodePGrp := models.NewAccessSwitchPolicyGroup(rn, parentDn, description, nameAlias, infraAccNodePGrpAttr)
	err := sm.Save(infraAccNodePGrp)
	return infraAccNodePGrp, err
}

func (sm *ServiceManager) ReadAccessSwitchPolicyGroup(name string) (*models.AccessSwitchPolicyGroup, error) {
	dn := fmt.Sprintf(models.DninfraAccNodePGrp, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	infraAccNodePGrp := models.AccessSwitchPolicyGroupFromContainer(cont)
	return infraAccNodePGrp, nil
}

func (sm *ServiceManager) DeleteAccessSwitchPolicyGroup(name string) error {
	dn := fmt.Sprintf(models.DninfraAccNodePGrp, name)
	return sm.DeleteByDn(dn, models.InfraaccnodepgrpClassName)
}

func (sm *ServiceManager) UpdateAccessSwitchPolicyGroup(name string, description string, nameAlias string, infraAccNodePGrpAttr models.AccessSwitchPolicyGroupAttributes) (*models.AccessSwitchPolicyGroup, error) {
	rn := fmt.Sprintf(models.RninfraAccNodePGrp, name)
	parentDn := fmt.Sprintf(models.ParentDninfraAccNodePGrp)
	infraAccNodePGrp := models.NewAccessSwitchPolicyGroup(rn, parentDn, description, nameAlias, infraAccNodePGrpAttr)
	infraAccNodePGrp.Status = "modified"
	err := sm.Save(infraAccNodePGrp)
	return infraAccNodePGrp, err
}

func (sm *ServiceManager) ListAccessSwitchPolicyGroup() ([]*models.AccessSwitchPolicyGroup, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/funcprof/infraAccNodePGrp.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.AccessSwitchPolicyGroupListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsBfdIpv4InstPol(parentDn, annotation, tnBfdIpv4InstPolName string) error {
	dn := fmt.Sprintf("%s/rsbfdIpv4InstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnBfdIpv4InstPolName": "%s"
			}
		}
	}`, "infraRsBfdIpv4InstPol", dn, annotation, tnBfdIpv4InstPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsBfdIpv4InstPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsbfdIpv4InstPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsBfdIpv4InstPol")
}

func (sm *ServiceManager) ReadRelationinfraRsBfdIpv4InstPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsBfdIpv4InstPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsBfdIpv4InstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnBfdIpv4InstPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsBfdIpv6InstPol(parentDn, annotation, tnBfdIpv6InstPolName string) error {
	dn := fmt.Sprintf("%s/rsbfdIpv6InstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnBfdIpv6InstPolName": "%s"
			}
		}
	}`, "infraRsBfdIpv6InstPol", dn, annotation, tnBfdIpv6InstPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsBfdIpv6InstPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsbfdIpv6InstPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsBfdIpv6InstPol")
}

func (sm *ServiceManager) ReadRelationinfraRsBfdIpv6InstPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsBfdIpv6InstPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsBfdIpv6InstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnBfdIpv6InstPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsBfdMhIpv4InstPol(parentDn, annotation, tnBfdMhIpv4InstPolName string) error {
	dn := fmt.Sprintf("%s/rsbfdMhIpv4InstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnBfdMhIpv4InstPolName": "%s"
			}
		}
	}`, "infraRsBfdMhIpv4InstPol", dn, annotation, tnBfdMhIpv4InstPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsBfdMhIpv4InstPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsbfdMhIpv4InstPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsBfdMhIpv4InstPol")
}

func (sm *ServiceManager) ReadRelationinfraRsBfdMhIpv4InstPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsBfdMhIpv4InstPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsBfdMhIpv4InstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnBfdMhIpv4InstPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsBfdMhIpv6InstPol(parentDn, annotation, tnBfdMhIpv6InstPolName string) error {
	dn := fmt.Sprintf("%s/rsbfdMhIpv6InstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnBfdMhIpv6InstPolName": "%s"
			}
		}
	}`, "infraRsBfdMhIpv6InstPol", dn, annotation, tnBfdMhIpv6InstPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsBfdMhIpv6InstPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsbfdMhIpv6InstPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsBfdMhIpv6InstPol")
}

func (sm *ServiceManager) ReadRelationinfraRsBfdMhIpv6InstPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsBfdMhIpv6InstPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsBfdMhIpv6InstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnBfdMhIpv6InstPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsEquipmentFlashConfigPol(parentDn, annotation, tnEquipmentFlashConfigPolName string) error {
	dn := fmt.Sprintf("%s/rsequipmentFlashConfigPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnEquipmentFlashConfigPolName": "%s"
			}
		}
	}`, "infraRsEquipmentFlashConfigPol", dn, annotation, tnEquipmentFlashConfigPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsEquipmentFlashConfigPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsequipmentFlashConfigPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsEquipmentFlashConfigPol")
}

func (sm *ServiceManager) ReadRelationinfraRsEquipmentFlashConfigPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsEquipmentFlashConfigPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsEquipmentFlashConfigPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnEquipmentFlashConfigPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsFcFabricPol(parentDn, annotation, tnFcFabricPolName string) error {
	dn := fmt.Sprintf("%s/rsfcFabricPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnFcFabricPolName": "%s"
			}
		}
	}`, "infraRsFcFabricPol", dn, annotation, tnFcFabricPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsFcFabricPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsfcFabricPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsFcFabricPol")
}

func (sm *ServiceManager) ReadRelationinfraRsFcFabricPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsFcFabricPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsFcFabricPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnFcFabricPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsFcInstPol(parentDn, annotation, tnFcInstPolName string) error {
	dn := fmt.Sprintf("%s/rsfcInstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnFcInstPolName": "%s"
			}
		}
	}`, "infraRsFcInstPol", dn, annotation, tnFcInstPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsFcInstPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsfcInstPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsFcInstPol")
}

func (sm *ServiceManager) ReadRelationinfraRsFcInstPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsFcInstPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsFcInstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnFcInstPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsIaclLeafProfile(parentDn, annotation, tnIaclLeafProfileName string) error {
	dn := fmt.Sprintf("%s/rsiaclLeafProfile", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnIaclLeafProfileName": "%s"
			}
		}
	}`, "infraRsIaclLeafProfile", dn, annotation, tnIaclLeafProfileName))

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

func (sm *ServiceManager) DeleteRelationinfraRsIaclLeafProfile(parentDn string) error {
	dn := fmt.Sprintf("%s/rsiaclLeafProfile", parentDn)
	return sm.DeleteByDn(dn, "infraRsIaclLeafProfile")
}

func (sm *ServiceManager) ReadRelationinfraRsIaclLeafProfile(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsIaclLeafProfile")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsIaclLeafProfile")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnIaclLeafProfileName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsL2NodeAuthPol(parentDn, annotation, tnL2NodeAuthPolName string) error {
	dn := fmt.Sprintf("%s/rsl2NodeAuthPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnL2NodeAuthPolName": "%s"
			}
		}
	}`, "infraRsL2NodeAuthPol", dn, annotation, tnL2NodeAuthPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsL2NodeAuthPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsl2NodeAuthPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsL2NodeAuthPol")
}

func (sm *ServiceManager) ReadRelationinfraRsL2NodeAuthPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsL2NodeAuthPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsL2NodeAuthPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnL2NodeAuthPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsLeafCoppProfile(parentDn, annotation, tnCoppLeafProfileName string) error {
	dn := fmt.Sprintf("%s/rsleafCoppProfile", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnCoppLeafProfileName": "%s"
			}
		}
	}`, "infraRsLeafCoppProfile", dn, annotation, tnCoppLeafProfileName))

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

func (sm *ServiceManager) DeleteRelationinfraRsLeafCoppProfile(parentDn string) error {
	dn := fmt.Sprintf("%s/rsleafCoppProfile", parentDn)
	return sm.DeleteByDn(dn, "infraRsLeafCoppProfile")
}

func (sm *ServiceManager) ReadRelationinfraRsLeafCoppProfile(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsLeafCoppProfile")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsLeafCoppProfile")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnCoppLeafProfileName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsLeafPGrpToCdpIfPol(parentDn, annotation, tnCdpIfPolName string) error {
	dn := fmt.Sprintf("%s/rsleafPGrpToCdpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnCdpIfPolName": "%s"
			}
		}
	}`, "infraRsLeafPGrpToCdpIfPol", dn, annotation, tnCdpIfPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsLeafPGrpToCdpIfPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsleafPGrpToCdpIfPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsLeafPGrpToCdpIfPol")
}

func (sm *ServiceManager) ReadRelationinfraRsLeafPGrpToCdpIfPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsLeafPGrpToCdpIfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsLeafPGrpToCdpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnCdpIfPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsLeafPGrpToLldpIfPol(parentDn, annotation, tnLldpIfPolName string) error {
	dn := fmt.Sprintf("%s/rsleafPGrpToLldpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnLldpIfPolName": "%s"
			}
		}
	}`, "infraRsLeafPGrpToLldpIfPol", dn, annotation, tnLldpIfPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsLeafPGrpToLldpIfPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsleafPGrpToLldpIfPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsLeafPGrpToLldpIfPol")
}

func (sm *ServiceManager) ReadRelationinfraRsLeafPGrpToLldpIfPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsLeafPGrpToLldpIfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsLeafPGrpToLldpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnLldpIfPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsMonNodeInfraPol(parentDn, annotation, tnMonInfraPolName string) error {
	dn := fmt.Sprintf("%s/rsmonNodeInfraPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnMonInfraPolName": "%s"
			}
		}
	}`, "infraRsMonNodeInfraPol", dn, annotation, tnMonInfraPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsMonNodeInfraPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsmonNodeInfraPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsMonNodeInfraPol")
}

func (sm *ServiceManager) ReadRelationinfraRsMonNodeInfraPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsMonNodeInfraPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsMonNodeInfraPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnMonInfraPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsMstInstPol(parentDn, annotation, tnStpInstPolName string) error {
	dn := fmt.Sprintf("%s/rsmstInstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnStpInstPolName": "%s"
			}
		}
	}`, "infraRsMstInstPol", dn, annotation, tnStpInstPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsMstInstPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsmstInstPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsMstInstPol")
}

func (sm *ServiceManager) ReadRelationinfraRsMstInstPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsMstInstPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsMstInstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnStpInstPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsNetflowNodePol(parentDn, annotation, tnNetflowNodePolName string) error {
	dn := fmt.Sprintf("%s/rsnetflowNodePol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnNetflowNodePolName": "%s"
			}
		}
	}`, "infraRsNetflowNodePol", dn, annotation, tnNetflowNodePolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsNetflowNodePol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsnetflowNodePol", parentDn)
	return sm.DeleteByDn(dn, "infraRsNetflowNodePol")
}

func (sm *ServiceManager) ReadRelationinfraRsNetflowNodePol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsNetflowNodePol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsNetflowNodePol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnNetflowNodePolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsPoeInstPol(parentDn, annotation, tnPoeInstPolName string) error {
	dn := fmt.Sprintf("%s/rspoeInstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnPoeInstPolName": "%s"
			}
		}
	}`, "infraRsPoeInstPol", dn, annotation, tnPoeInstPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsPoeInstPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rspoeInstPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsPoeInstPol")
}

func (sm *ServiceManager) ReadRelationinfraRsPoeInstPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsPoeInstPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsPoeInstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnPoeInstPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsTopoctrlFastLinkFailoverInstPol(parentDn, annotation, tnTopoctrlFastLinkFailoverInstPolName string) error {
	dn := fmt.Sprintf("%s/rstopoctrlFastLinkFailoverInstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnTopoctrlFastLinkFailoverInstPolName": "%s"
			}
		}
	}`, "infraRsTopoctrlFastLinkFailoverInstPol", dn, annotation, tnTopoctrlFastLinkFailoverInstPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsTopoctrlFastLinkFailoverInstPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rstopoctrlFastLinkFailoverInstPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsTopoctrlFastLinkFailoverInstPol")
}

func (sm *ServiceManager) ReadRelationinfraRsTopoctrlFastLinkFailoverInstPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsTopoctrlFastLinkFailoverInstPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsTopoctrlFastLinkFailoverInstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnTopoctrlFastLinkFailoverInstPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsTopoctrlFwdScaleProfPol(parentDn, annotation, tnTopoctrlFwdScaleProfilePolName string) error {
	dn := fmt.Sprintf("%s/rstopoctrlFwdScaleProfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnTopoctrlFwdScaleProfilePolName": "%s"
			}
		}
	}`, "infraRsTopoctrlFwdScaleProfPol", dn, annotation, tnTopoctrlFwdScaleProfilePolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsTopoctrlFwdScaleProfPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rstopoctrlFwdScaleProfPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsTopoctrlFwdScaleProfPol")
}

func (sm *ServiceManager) ReadRelationinfraRsTopoctrlFwdScaleProfPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsTopoctrlFwdScaleProfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsTopoctrlFwdScaleProfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnTopoctrlFwdScaleProfilePolName")
		return dat, err
	} else {
		return nil, err
	}
}

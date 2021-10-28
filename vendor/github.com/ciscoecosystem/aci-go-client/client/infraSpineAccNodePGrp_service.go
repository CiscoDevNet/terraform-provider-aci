package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSpineSwitchPolicyGroup(name string, description string, nameAlias string, infraSpineAccNodePGrpAttr models.SpineSwitchPolicyGroupAttributes) (*models.SpineSwitchPolicyGroup, error) {
	rn := fmt.Sprintf(models.RninfraSpineAccNodePGrp, name)
	parentDn := fmt.Sprintf(models.ParentDninfraSpineAccNodePGrp)
	infraSpineAccNodePGrp := models.NewSpineSwitchPolicyGroup(rn, parentDn, description, nameAlias, infraSpineAccNodePGrpAttr)
	err := sm.Save(infraSpineAccNodePGrp)
	return infraSpineAccNodePGrp, err
}

func (sm *ServiceManager) ReadSpineSwitchPolicyGroup(name string) (*models.SpineSwitchPolicyGroup, error) {
	dn := fmt.Sprintf(models.DninfraSpineAccNodePGrp, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	infraSpineAccNodePGrp := models.SpineSwitchPolicyGroupFromContainer(cont)
	return infraSpineAccNodePGrp, nil
}

func (sm *ServiceManager) DeleteSpineSwitchPolicyGroup(name string) error {
	dn := fmt.Sprintf(models.DninfraSpineAccNodePGrp, name)
	return sm.DeleteByDn(dn, models.InfraspineaccnodepgrpClassName)
}

func (sm *ServiceManager) UpdateSpineSwitchPolicyGroup(name string, description string, nameAlias string, infraSpineAccNodePGrpAttr models.SpineSwitchPolicyGroupAttributes) (*models.SpineSwitchPolicyGroup, error) {
	rn := fmt.Sprintf(models.RninfraSpineAccNodePGrp, name)
	parentDn := fmt.Sprintf(models.ParentDninfraSpineAccNodePGrp)
	infraSpineAccNodePGrp := models.NewSpineSwitchPolicyGroup(rn, parentDn, description, nameAlias, infraSpineAccNodePGrpAttr)
	infraSpineAccNodePGrp.Status = "modified"
	err := sm.Save(infraSpineAccNodePGrp)
	return infraSpineAccNodePGrp, err
}

func (sm *ServiceManager) ListSpineSwitchPolicyGroup() ([]*models.SpineSwitchPolicyGroup, error) {
	dnUrl := fmt.Sprintf("%s/uni/infra/funcprof/infraSpineAccNodePGrp.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SpineSwitchPolicyGroupListFromContainer(cont)
	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsIaclSpineProfile(parentDn, annotation, tnIaclSpineProfileName string) error {
	dn := fmt.Sprintf("%s/rsiaclSpineProfile", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnIaclSpineProfileName": "%s"
			}
		}
	}`, "infraRsIaclSpineProfile", dn, annotation, tnIaclSpineProfileName))

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

func (sm *ServiceManager) DeleteRelationinfraRsIaclSpineProfile(parentDn string) error {
	dn := fmt.Sprintf("%s/rsiaclSpineProfile", parentDn)
	return sm.DeleteByDn(dn, "infraRsIaclSpineProfile")
}

func (sm *ServiceManager) ReadRelationinfraRsIaclSpineProfile(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsIaclSpineProfile")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsIaclSpineProfile")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnIaclSpineProfileName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsSpineBfdIpv4InstPol(parentDn, annotation, tnBfdIpv4InstPolName string) error {
	dn := fmt.Sprintf("%s/rsspineBfdIpv4InstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnBfdIpv4InstPolName": "%s"
			}
		}
	}`, "infraRsSpineBfdIpv4InstPol", dn, annotation, tnBfdIpv4InstPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsSpineBfdIpv4InstPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsspineBfdIpv4InstPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsSpineBfdIpv4InstPol")
}

func (sm *ServiceManager) ReadRelationinfraRsSpineBfdIpv4InstPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsSpineBfdIpv4InstPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsSpineBfdIpv4InstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnBfdIpv4InstPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsSpineBfdIpv6InstPol(parentDn, annotation, tnBfdIpv6InstPolName string) error {
	dn := fmt.Sprintf("%s/rsspineBfdIpv6InstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnBfdIpv6InstPolName": "%s"
			}
		}
	}`, "infraRsSpineBfdIpv6InstPol", dn, annotation, tnBfdIpv6InstPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsSpineBfdIpv6InstPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsspineBfdIpv6InstPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsSpineBfdIpv6InstPol")
}

func (sm *ServiceManager) ReadRelationinfraRsSpineBfdIpv6InstPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsSpineBfdIpv6InstPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsSpineBfdIpv6InstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnBfdIpv6InstPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsSpineCoppProfile(parentDn, annotation, tnCoppSpineProfileName string) error {
	dn := fmt.Sprintf("%s/rsspineCoppProfile", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnCoppSpineProfileName": "%s"
			}
		}
	}`, "infraRsSpineCoppProfile", dn, annotation, tnCoppSpineProfileName))

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

func (sm *ServiceManager) DeleteRelationinfraRsSpineCoppProfile(parentDn string) error {
	dn := fmt.Sprintf("%s/rsspineCoppProfile", parentDn)
	return sm.DeleteByDn(dn, "infraRsSpineCoppProfile")
}

func (sm *ServiceManager) ReadRelationinfraRsSpineCoppProfile(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsSpineCoppProfile")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsSpineCoppProfile")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnCoppSpineProfileName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsSpinePGrpToCdpIfPol(parentDn, annotation, tnCdpIfPolName string) error {
	dn := fmt.Sprintf("%s/rsspinePGrpToCdpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnCdpIfPolName": "%s"
			}
		}
	}`, "infraRsSpinePGrpToCdpIfPol", dn, annotation, tnCdpIfPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsSpinePGrpToCdpIfPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsspinePGrpToCdpIfPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsSpinePGrpToCdpIfPol")
}

func (sm *ServiceManager) ReadRelationinfraRsSpinePGrpToCdpIfPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsSpinePGrpToCdpIfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsSpinePGrpToCdpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnCdpIfPolName")
		return dat, err
	} else {
		return nil, err
	}
}

func (sm *ServiceManager) CreateRelationinfraRsSpinePGrpToLldpIfPol(parentDn, annotation, tnLldpIfPolName string) error {
	dn := fmt.Sprintf("%s/rsspinePGrpToLldpIfPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"annotation": "%s",
				"tnLldpIfPolName": "%s"
			}
		}
	}`, "infraRsSpinePGrpToLldpIfPol", dn, annotation, tnLldpIfPolName))

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

func (sm *ServiceManager) DeleteRelationinfraRsSpinePGrpToLldpIfPol(parentDn string) error {
	dn := fmt.Sprintf("%s/rsspinePGrpToLldpIfPol", parentDn)
	return sm.DeleteByDn(dn, "infraRsSpinePGrpToLldpIfPol")
}

func (sm *ServiceManager) ReadRelationinfraRsSpinePGrpToLldpIfPol(parentDn string) (interface{}, error) {
	dnUrl := fmt.Sprintf("%s/%s/%s.json", models.BaseurlStr, parentDn, "infraRsSpinePGrpToLldpIfPol")
	cont, err := sm.GetViaURL(dnUrl)
	contList := models.ListFromContainer(cont, "infraRsSpinePGrpToLldpIfPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnLldpIfPolName")
		return dat, err
	} else {
		return nil, err
	}
}

package client

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateVPCExplicitProtectionGroup(name string, description string, switch1 string, switch2 string, vpcDomainPolicy string, fabricExplicitGEpattr models.VPCExplicitProtectionGroupAttributes) (*models.VPCExplicitProtectionGroup, error) {
	rn := fmt.Sprintf("fabric/protpol/expgep-%s", name)
	parentDn := fmt.Sprintf("uni")
	fabricExplicitGEp := models.NewVPCExplicitProtectionGroup(rn, parentDn, description, fabricExplicitGEpattr)
	jsonPayload, _, err := sm.PrepareModel(fabricExplicitGEp)
	fabricNodePEp1 := []byte(fmt.Sprintf(`
	{
		"fabricNodePEp": {
			"attributes": {
				"id": "%s"
			}
		}
	}
	`, switch1))
	fabricNodePEp2 := []byte(fmt.Sprintf(`
	{
		"fabricNodePEp": {
			"attributes": {
				"id": "%s"
			}
		}
	}
	`, switch2))
	fabricRsVpcInstPol := []byte(fmt.Sprintf(`
	{
		"fabricRsVpcInstPol": {
			"attributes": {
				"tnVpcInstPolName": "%s"
			}
		}
	}
	`, vpcDomainPolicy))
	fabricNodePEp1Cont, err := container.ParseJSON(fabricNodePEp1)
	fabricNodePEp2Cont, err := container.ParseJSON(fabricNodePEp2)
	fabricRsVpcInstPolCont, err := container.ParseJSON(fabricRsVpcInstPol)
	if err != nil {
		return nil, err
	}
	jsonPayload.Array(fabricExplicitGEp.ClassName, "children")
	jsonPayload.ArrayAppend(fabricNodePEp1Cont.Data(), fabricExplicitGEp.ClassName, "children")
	jsonPayload.ArrayAppend(fabricNodePEp2Cont.Data(), fabricExplicitGEp.ClassName, "children")
	jsonPayload.ArrayAppend(fabricRsVpcInstPolCont.Data(), fabricExplicitGEp.ClassName, "children")
	jsonPayload.Set(name, fabricExplicitGEp.ClassName, "attributes", "name")
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("/api/node/mo/%s/%s.json", parentDn, rn), jsonPayload, true)
	if err != nil {
		return nil, err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return nil, err
	}

	return fabricExplicitGEp, CheckErrorCustom(cont)

}

func (sm *ServiceManager) ReadVPCExplicitProtectionGroup(name string) (*models.VPCExplicitProtectionGroup, error) {
	dn := fmt.Sprintf("uni/fabric/protpol/expgep-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricExplicitGEp := models.VPCExplicitProtectionGroupFromContainer(cont)
	return fabricExplicitGEp, nil
}

func (sm *ServiceManager) DeleteVPCExplicitProtectionGroup(name string) error {
	dn := fmt.Sprintf("uni/fabric/protpol/expgep-%s", name)
	return sm.DeleteByDn(dn, models.FabricexplicitgepClassName)
}

func (sm *ServiceManager) UpdateVPCExplicitProtectionGroup(name string, description string, switch1 string, switch2 string, vpcDomainPolicy string, fabricExplicitGEpattr models.VPCExplicitProtectionGroupAttributes) (*models.VPCExplicitProtectionGroup, error) {
	rn := fmt.Sprintf("fabric/protpol/expgep-%s", name)
	parentDn := fmt.Sprintf("uni")
	fabricExplicitGEp := models.NewVPCExplicitProtectionGroup(rn, parentDn, description, fabricExplicitGEpattr)

	fabricExplicitGEp.Status = "modified"
	jsonPayload, _, err := sm.PrepareModel(fabricExplicitGEp)
	fabricNodePEp1 := []byte(fmt.Sprintf(`
	{
		"fabricNodePEp": {
			"attributes": {
				"id": "%s"
			}
		}
	}
	`, switch1))
	fabricNodePEp2 := []byte(fmt.Sprintf(`
	{
		"fabricNodePEp": {
			"attributes": {
				"id": "%s"
			}
		}
	}
	`, switch2))
	fabricRsVpcInstPol := []byte(fmt.Sprintf(`
	{
		"fabricRsVpcInstPol": {
			"attributes": {
				"tnVpcInstPolName": "%s"
			}
		}
	}
	`, vpcDomainPolicy))
	fabricNodePEp1Cont, err := container.ParseJSON(fabricNodePEp1)
	fabricNodePEp2Cont, err := container.ParseJSON(fabricNodePEp2)
	fabricRsVpcInstPolCont, err := container.ParseJSON(fabricRsVpcInstPol)
	if err != nil {
		return nil, err
	}
	jsonPayload.Array(fabricExplicitGEp.ClassName, "children")
	jsonPayload.ArrayAppend(fabricNodePEp1Cont.Data(), fabricExplicitGEp.ClassName, "children")
	jsonPayload.ArrayAppend(fabricNodePEp2Cont.Data(), fabricExplicitGEp.ClassName, "children")
	jsonPayload.ArrayAppend(fabricRsVpcInstPolCont.Data(), fabricExplicitGEp.ClassName, "children")
	jsonPayload.Set(name, fabricExplicitGEp.ClassName, "attributes", "name")
	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("/api/node/mo/%s/%s.json", parentDn, rn), jsonPayload, true)
	if err != nil {
		return nil, err
	}
	cont, _, err := sm.client.Do(req)
	if err != nil {
		return nil, err
	}

	return fabricExplicitGEp, CheckErrorCustom(cont)

}

func CheckErrorCustom(cont *container.Container) error {
	number, err := strconv.Atoi(models.G(cont, "totalCount"))
	if err != nil {
		return err
	}
	imdata := cont.S("imdata").Index(0)
	if number > 0 {

		if imdata.Exists("error") {
			return errors.New(models.StripQuotes(imdata.Path("error.attributes.text").String()))
		}
	}
	return nil
}

func (sm *ServiceManager) ListVPCExplicitProtectionGroup() ([]*models.VPCExplicitProtectionGroup, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/fabricExplicitGEp.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.VPCExplicitProtectionGroupListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationfabricRsVpcInstPolFromVPCExplicitProtectionGroup(parentDn, tnVpcInstPolName string) error {
	dn := fmt.Sprintf("%s/rsvpcInstPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnVpcInstPolName": "%s"
								
			}
		}
	}`, "fabricRsVpcInstPol", dn, tnVpcInstPolName))

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

func (sm *ServiceManager) ReadRelationfabricRsVpcInstPolFromVPCExplicitProtectionGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fabricRsVpcInstPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fabricRsVpcInstPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tnVpcInstPolName")
		return dat, err
	} else {
		return nil, err
	}

}

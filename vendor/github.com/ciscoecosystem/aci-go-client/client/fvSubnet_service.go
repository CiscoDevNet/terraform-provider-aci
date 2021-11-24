package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateSubnet(ip string, bridge_domain string, tenant string, description string, fvSubnetattr models.SubnetAttributes) (*models.Subnet, error) {
	rn := fmt.Sprintf("subnet-[%s]", ip)
	parentDn := fmt.Sprintf("uni/tn-%s/BD-%s", tenant, bridge_domain)
	fvSubnet := models.NewSubnet(rn, parentDn, description, fvSubnetattr)
	err := sm.Save(fvSubnet)
	return fvSubnet, err
}

func (sm *ServiceManager) ReadSubnet(ip string, bridge_domain string, tenant string) (*models.Subnet, error) {
	dn := fmt.Sprintf("uni/tn-%s/BD-%s/subnet-[%s]", tenant, bridge_domain, ip)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	fvSubnet := models.SubnetFromContainer(cont)
	return fvSubnet, nil
}

func (sm *ServiceManager) DeleteSubnet(ip string, bridge_domain string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/BD-%s/subnet-[%s]", tenant, bridge_domain, ip)
	return sm.DeleteByDn(dn, models.FvsubnetClassName)
}

func (sm *ServiceManager) UpdateSubnet(ip string, bridge_domain string, tenant string, description string, fvSubnetattr models.SubnetAttributes) (*models.Subnet, error) {
	rn := fmt.Sprintf("subnet-[%s]", ip)
	parentDn := fmt.Sprintf("uni/tn-%s/BD-%s", tenant, bridge_domain)
	fvSubnet := models.NewSubnet(rn, parentDn, description, fvSubnetattr)

	fvSubnet.Status = "modified"
	err := sm.Save(fvSubnet)
	return fvSubnet, err

}

func (sm *ServiceManager) ListSubnet(bridge_domain string, tenant string) ([]*models.Subnet, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/BD-%s/fvSubnet.json", baseurlStr, tenant, bridge_domain)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.SubnetListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationfvRsBDSubnetToOutFromSubnet(parentDn, tnL3extOutName string) error {
	dn := fmt.Sprintf("%s/rsBDSubnetToOut-%s", parentDn, tnL3extOutName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsBDSubnetToOut", dn))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) DeleteRelationfvRsBDSubnetToOutFromSubnet(parentDn, tnL3extOutName string) error {
	dn := fmt.Sprintf("%s/rsBDSubnetToOut-%s", parentDn, tnL3extOutName)
	return sm.DeleteByDn(dn, "fvRsBDSubnetToOut")
}

func (sm *ServiceManager) ReadRelationfvRsBDSubnetToOutFromSubnet(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsBDSubnetToOut")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsBDSubnetToOut")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsNdPfxPolFromSubnet(parentDn, tnNdPfxPolName string) error {
	dn := fmt.Sprintf("%s/rsNdPfxPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnNdPfxPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsNdPfxPol", dn, tnNdPfxPolName))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) DeleteRelationfvRsNdPfxPolFromSubnet(parentDn string) error {
	dn := fmt.Sprintf("%s/rsNdPfxPol", parentDn)
	return sm.DeleteByDn(dn, "fvRsNdPfxPol")
}

func (sm *ServiceManager) ReadRelationfvRsNdPfxPolFromSubnet(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsNdPfxPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsNdPfxPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsBDSubnetToProfileFromSubnet(parentDn, tnRtctrlProfileName string) error {
	dn := fmt.Sprintf("%s/rsBDSubnetToProfile", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnRtctrlProfileName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsBDSubnetToProfile", dn, tnRtctrlProfileName))

	jsonPayload, err := container.ParseJSON(containerJSON)
	if err != nil {
		return err
	}

	req, err := sm.client.MakeRestRequest("POST", fmt.Sprintf("%s.json", sm.MOURL), jsonPayload, true)
	if err != nil {
		return err
	}

	_, _, err = sm.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (sm *ServiceManager) DeleteRelationfvRsBDSubnetToProfileFromSubnet(parentDn string) error {
	dn := fmt.Sprintf("%s/rsBDSubnetToProfile", parentDn)
	return sm.DeleteByDn(dn, "fvRsBDSubnetToProfile")
}

func (sm *ServiceManager) ReadRelationfvRsBDSubnetToProfileFromSubnet(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsBDSubnetToProfile")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsBDSubnetToProfile")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}

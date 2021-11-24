package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateFexBundleGroup(name string, fex_profile string, description string, infraFexBndlGrpattr models.FexBundleGroupAttributes) (*models.FexBundleGroup, error) {
	rn := fmt.Sprintf("fexbundle-%s", name)
	parentDn := fmt.Sprintf("uni/infra/fexprof-%s", fex_profile)
	infraFexBndlGrp := models.NewFexBundleGroup(rn, parentDn, description, infraFexBndlGrpattr)
	err := sm.Save(infraFexBndlGrp)
	return infraFexBndlGrp, err
}

func (sm *ServiceManager) ReadFexBundleGroup(name string, fex_profile string) (*models.FexBundleGroup, error) {
	dn := fmt.Sprintf("uni/infra/fexprof-%s/fexbundle-%s", fex_profile, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	infraFexBndlGrp := models.FexBundleGroupFromContainer(cont)
	return infraFexBndlGrp, nil
}

func (sm *ServiceManager) DeleteFexBundleGroup(name string, fex_profile string) error {
	dn := fmt.Sprintf("uni/infra/fexprof-%s/fexbundle-%s", fex_profile, name)
	return sm.DeleteByDn(dn, models.InfrafexbndlgrpClassName)
}

func (sm *ServiceManager) UpdateFexBundleGroup(name string, fex_profile string, description string, infraFexBndlGrpattr models.FexBundleGroupAttributes) (*models.FexBundleGroup, error) {
	rn := fmt.Sprintf("fexbundle-%s", name)
	parentDn := fmt.Sprintf("uni/infra/fexprof-%s", fex_profile)
	infraFexBndlGrp := models.NewFexBundleGroup(rn, parentDn, description, infraFexBndlGrpattr)

	infraFexBndlGrp.Status = "modified"
	err := sm.Save(infraFexBndlGrp)
	return infraFexBndlGrp, err

}

func (sm *ServiceManager) ListFexBundleGroup(fex_profile string) ([]*models.FexBundleGroup, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/infra/fexprof-%s/infraFexBndlGrp.json", baseurlStr, fex_profile)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.FexBundleGroupListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationinfraRsMonFexInfraPolFromFexBundleGroup(parentDn, tnMonInfraPolName string) error {
	dn := fmt.Sprintf("%s/rsmonFexInfraPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnMonInfraPolName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "infraRsMonFexInfraPol", dn, tnMonInfraPolName))

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

func (sm *ServiceManager) ReadRelationinfraRsMonFexInfraPolFromFexBundleGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsMonFexInfraPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsMonFexInfraPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationinfraRsFexBndlGrpToAggrIfFromFexBundleGroup(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsfexBndlGrpToAggrIf-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s"				
			}
		}
	}`, "infraRsFexBndlGrpToAggrIf", dn))

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

func (sm *ServiceManager) ReadRelationinfraRsFexBndlGrpToAggrIfFromFexBundleGroup(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "infraRsFexBndlGrpToAggrIf")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "infraRsFexBndlGrpToAggrIf")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}

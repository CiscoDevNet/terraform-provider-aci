package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sm *ServiceManager) CreateL2outExternalEpg(name string, l2_outside string, tenant string, description string, l2extInstPattr models.L2outExternalEpgAttributes) (*models.L2outExternalEpg, error) {
	rn := fmt.Sprintf("instP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/l2out-%s", tenant, l2_outside)
	l2extInstP := models.NewL2outExternalEpg(rn, parentDn, description, l2extInstPattr)
	err := sm.Save(l2extInstP)
	return l2extInstP, err
}

func (sm *ServiceManager) ReadL2outExternalEpg(name string, l2_outside string, tenant string) (*models.L2outExternalEpg, error) {
	dn := fmt.Sprintf("uni/tn-%s/l2out-%s/instP-%s", tenant, l2_outside, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l2extInstP := models.L2outExternalEpgFromContainer(cont)
	return l2extInstP, nil
}

func (sm *ServiceManager) DeleteL2outExternalEpg(name string, l2_outside string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/l2out-%s/instP-%s", tenant, l2_outside, name)
	return sm.DeleteByDn(dn, models.L2extinstpClassName)
}

func (sm *ServiceManager) UpdateL2outExternalEpg(name string, l2_outside string, tenant string, description string, l2extInstPattr models.L2outExternalEpgAttributes) (*models.L2outExternalEpg, error) {
	rn := fmt.Sprintf("instP-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/l2out-%s", tenant, l2_outside)
	l2extInstP := models.NewL2outExternalEpg(rn, parentDn, description, l2extInstPattr)

	l2extInstP.Status = "modified"
	err := sm.Save(l2extInstP)
	return l2extInstP, err

}

func (sm *ServiceManager) ListL2outExternalEpg(l2_outside string, tenant string) ([]*models.L2outExternalEpg, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/l2out-%s/l2extInstP.json", baseurlStr, tenant, l2_outside)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L2outExternalEpgListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationfvRsSecInheritedFromL2outExternalEpg(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tDn": "%s",
				"annotation":"orchestrator:terraform"
			}
		}
	}`, "fvRsSecInherited", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationfvRsSecInheritedFromL2outExternalEpg(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rssecInherited-[%s]", parentDn, tDn)
	return sm.DeleteByDn(dn, "fvRsSecInherited")
}

func (sm *ServiceManager) ReadRelationfvRsSecInheritedFromL2outExternalEpg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsSecInherited")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsSecInherited")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsProvFromL2outExternalEpg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnVzBrCPName": "%s",
				"annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsProv", dn, tnVzBrCPName))

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

func (sm *ServiceManager) DeleteRelationfvRsProvFromL2outExternalEpg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsprov-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsProv")
}

func (sm *ServiceManager) ReadRelationfvRsProvFromL2outExternalEpg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsProv")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsProv")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsConsIfFromL2outExternalEpg(parentDn, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnVzCPIfName": "%s",
				"annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsConsIf", dn, tnVzCPIfName))

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

func (sm *ServiceManager) DeleteRelationfvRsConsIfFromL2outExternalEpg(parentDn, tnVzCPIfName string) error {
	dn := fmt.Sprintf("%s/rsconsIf-%s", parentDn, tnVzCPIfName)
	return sm.DeleteByDn(dn, "fvRsConsIf")
}

func (sm *ServiceManager) ReadRelationfvRsConsIfFromL2outExternalEpg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsConsIf")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsConsIf")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsCustQosPolFromL2outExternalEpg(parentDn, tnQosCustomPolName string) error {
	dn := fmt.Sprintf("%s/rscustQosPol", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnQosCustomPolName": "%s",
				"annotation":"orchestrator:terraform"
								
			}
		}
	}`, "fvRsCustQosPol", dn, tnQosCustomPolName))

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

func (sm *ServiceManager) ReadRelationfvRsCustQosPolFromL2outExternalEpg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsCustQosPol")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsCustQosPol")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsConsFromL2outExternalEpg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnVzBrCPName": "%s",
				"annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsCons", dn, tnVzBrCPName))

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

func (sm *ServiceManager) DeleteRelationfvRsConsFromL2outExternalEpg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rscons-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsCons")
}

func (sm *ServiceManager) ReadRelationfvRsConsFromL2outExternalEpg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsCons")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsCons")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationl2extRsL2InstPToDomPFromL2outExternalEpg(parentDn, tnL2extDomPName string) error {
	dn := fmt.Sprintf("%s/rsl2InstPToDomP", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s"
								
			}
		}
	}`, "l2extRsL2InstPToDomP", dn, tnL2extDomPName))

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

func (sm *ServiceManager) ReadRelationl2extRsL2InstPToDomPFromL2outExternalEpg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l2extRsL2InstPToDomP")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l2extRsL2InstPToDomP")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationfvRsProtByFromL2outExternalEpg(parentDn, tnVzTabooName string) error {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnVzTabooName": "%s",
				"annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsProtBy", dn, tnVzTabooName))

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

func (sm *ServiceManager) DeleteRelationfvRsProtByFromL2outExternalEpg(parentDn, tnVzTabooName string) error {
	dn := fmt.Sprintf("%s/rsprotBy-%s", parentDn, tnVzTabooName)
	return sm.DeleteByDn(dn, "fvRsProtBy")
}

func (sm *ServiceManager) ReadRelationfvRsProtByFromL2outExternalEpg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsProtBy")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsProtBy")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}
func (sm *ServiceManager) CreateRelationfvRsIntraEpgFromL2outExternalEpg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s",
				"tnVzBrCPName": "%s",
				"annotation":"orchestrator:terraform"				
			}
		}
	}`, "fvRsIntraEpg", dn, tnVzBrCPName))

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

func (sm *ServiceManager) DeleteRelationfvRsIntraEpgFromL2outExternalEpg(parentDn, tnVzBrCPName string) error {
	dn := fmt.Sprintf("%s/rsintraEpg-%s", parentDn, tnVzBrCPName)
	return sm.DeleteByDn(dn, "fvRsIntraEpg")
}

func (sm *ServiceManager) ReadRelationfvRsIntraEpgFromL2outExternalEpg(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "fvRsIntraEpg")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "fvRsIntraEpg")

	st := &schema.Set{
		F: schema.HashString,
	}
	for _, contItem := range contList {
		dat := models.G(contItem, "tDn")
		st.Add(dat)
	}
	return st, err

}

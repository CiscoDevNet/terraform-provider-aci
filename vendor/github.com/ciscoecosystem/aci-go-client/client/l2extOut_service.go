package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateL2Outside(name string, tenant string, description string, l2extOutattr models.L2OutsideAttributes) (*models.L2Outside, error) {
	rn := fmt.Sprintf("l2out-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	l2extOut := models.NewL2Outside(rn, parentDn, description, l2extOutattr)
	err := sm.Save(l2extOut)
	return l2extOut, err
}

func (sm *ServiceManager) ReadL2Outside(name string, tenant string) (*models.L2Outside, error) {
	dn := fmt.Sprintf("uni/tn-%s/l2out-%s", tenant, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	l2extOut := models.L2OutsideFromContainer(cont)
	return l2extOut, nil
}

func (sm *ServiceManager) DeleteL2Outside(name string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/l2out-%s", tenant, name)
	return sm.DeleteByDn(dn, models.L2extoutClassName)
}

func (sm *ServiceManager) UpdateL2Outside(name string, tenant string, description string, l2extOutattr models.L2OutsideAttributes) (*models.L2Outside, error) {
	rn := fmt.Sprintf("l2out-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s", tenant)
	l2extOut := models.NewL2Outside(rn, parentDn, description, l2extOutattr)

	l2extOut.Status = "modified"
	err := sm.Save(l2extOut)
	return l2extOut, err

}

func (sm *ServiceManager) ListL2Outside(tenant string) ([]*models.L2Outside, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/l2extOut.json", baseurlStr, tenant)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.L2OutsideListFromContainer(cont)

	return list, err
}

func (sm *ServiceManager) CreateRelationl2extRsEBdFromL2Outside(parentDn, tnFvBDName string) error {
	dn := fmt.Sprintf("%s/rseBd", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tnFvBDName": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l2extRsEBd", dn, tnFvBDName))

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

func (sm *ServiceManager) ReadRelationl2extRsEBdFromL2Outside(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l2extRsEBd")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l2extRsEBd")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}
func (sm *ServiceManager) CreateRelationl2extRsL2DomAttFromL2Outside(parentDn, tDn string) error {
	dn := fmt.Sprintf("%s/rsl2DomAtt", parentDn)
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
				"dn": "%s","tDn": "%s","annotation":"orchestrator:terraform"
								
			}
		}
	}`, "l2extRsL2DomAtt", dn, tDn))

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

func (sm *ServiceManager) DeleteRelationl2extRsL2DomAttFromL2Outside(parentDn string) error {
	dn := fmt.Sprintf("%s/rsl2DomAtt", parentDn)
	return sm.DeleteByDn(dn, "l2extRsL2DomAtt")
}

func (sm *ServiceManager) ReadRelationl2extRsL2DomAttFromL2Outside(parentDn string) (interface{}, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, parentDn, "l2extRsL2DomAtt")
	cont, err := sm.GetViaURL(dnUrl)

	contList := models.ListFromContainer(cont, "l2extRsL2DomAtt")

	if len(contList) > 0 {
		dat := models.G(contList[0], "tDn")
		return dat, err
	} else {
		return nil, err
	}

}

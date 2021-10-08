package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateMgmtconnectivitypreference(description string, nameAlias string, mgmtConnectivityPrefsAttr models.MgmtconnectivitypreferenceAttributes) (*models.Mgmtconnectivitypreference, error) {
	rn := fmt.Sprintf(models.RnmgmtConnectivityPrefs)
	parentDn := fmt.Sprintf(models.ParentDnmgmtConnectivityPrefs)
	mgmtConnectivityPrefs := models.NewMgmtconnectivitypreference(rn, parentDn, description, nameAlias, mgmtConnectivityPrefsAttr)
	err := sm.Save(mgmtConnectivityPrefs)
	return mgmtConnectivityPrefs, err
}

func (sm *ServiceManager) ReadMgmtconnectivitypreference() (*models.Mgmtconnectivitypreference, error) {
	dn := fmt.Sprintf(models.DnmgmtConnectivityPrefs)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	mgmtConnectivityPrefs := models.MgmtconnectivitypreferenceFromContainer(cont)
	return mgmtConnectivityPrefs, nil
}

func (sm *ServiceManager) DeleteMgmtconnectivitypreference() error {
	dn := fmt.Sprintf(models.DnmgmtConnectivityPrefs)
	return sm.DeleteByDn(dn, models.MgmtconnectivityprefsClassName)
}

func (sm *ServiceManager) UpdateMgmtconnectivitypreference(description string, nameAlias string, mgmtConnectivityPrefsAttr models.MgmtconnectivitypreferenceAttributes) (*models.Mgmtconnectivitypreference, error) {
	rn := fmt.Sprintf(models.RnmgmtConnectivityPrefs)
	parentDn := fmt.Sprintf(models.ParentDnmgmtConnectivityPrefs)
	mgmtConnectivityPrefs := models.NewMgmtconnectivitypreference(rn, parentDn, description, nameAlias, mgmtConnectivityPrefsAttr)
	mgmtConnectivityPrefs.Status = "modified"
	err := sm.Save(mgmtConnectivityPrefs)
	return mgmtConnectivityPrefs, err
}

func (sm *ServiceManager) ListMgmtconnectivitypreference() ([]*models.Mgmtconnectivitypreference, error) {
	dnUrl := fmt.Sprintf("%s/uni/fabric/mgmtConnectivityPrefs.json", models.BaseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.MgmtconnectivitypreferenceListFromContainer(cont)
	return list, err
}

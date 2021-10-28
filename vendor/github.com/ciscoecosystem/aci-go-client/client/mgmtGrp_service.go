package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateManagedNodeConnectivityGroup(name string, description string, mgmtGrpattr models.ManagedNodeConnectivityGroupAttributes) (*models.ManagedNodeConnectivityGroup, error) {
	rn := fmt.Sprintf("infra/funcprof/grp-%s", name)
	parentDn := fmt.Sprintf("uni")
	mgmtGrp := models.NewManagedNodeConnectivityGroup(rn, parentDn, mgmtGrpattr)
	err := sm.Save(mgmtGrp)
	return mgmtGrp, err
}

func (sm *ServiceManager) ReadManagedNodeConnectivityGroup(name string) (*models.ManagedNodeConnectivityGroup, error) {
	dn := fmt.Sprintf("uni/infra/funcprof/grp-%s", name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	mgmtGrp := models.ManagedNodeConnectivityGroupFromContainer(cont)
	return mgmtGrp, nil
}

func (sm *ServiceManager) DeleteManagedNodeConnectivityGroup(name string) error {
	dn := fmt.Sprintf("uni/infra/funcprof/grp-%s", name)
	return sm.DeleteByDn(dn, models.MgmtgrpClassName)
}

func (sm *ServiceManager) UpdateManagedNodeConnectivityGroup(name string, description string, mgmtGrpattr models.ManagedNodeConnectivityGroupAttributes) (*models.ManagedNodeConnectivityGroup, error) {
	rn := fmt.Sprintf("infra/funcprof/grp-%s", name)
	parentDn := fmt.Sprintf("uni")
	mgmtGrp := models.NewManagedNodeConnectivityGroup(rn, parentDn, mgmtGrpattr)

	mgmtGrp.Status = "modified"
	err := sm.Save(mgmtGrp)
	return mgmtGrp, err

}

func (sm *ServiceManager) ListManagedNodeConnectivityGroup() ([]*models.ManagedNodeConnectivityGroup, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/mgmtGrp.json", baseurlStr)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.ManagedNodeConnectivityGroupListFromContainer(cont)

	return list, err
}

package client

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) ReadUpgJob(pod int, node int) (*models.UpgJob, error) {
	dn := fmt.Sprintf("topology/pod-%d/node-%d/sys/fwstatuscont/upgjob", pod, node)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	maintUpgJob := models.UpgJobFromContainer(cont)
	return maintUpgJob, nil
}

func (sm *ServiceManager) ListUpgJob() ([]*models.UpgJob, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/maintUpgJob.json", baseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.UpgJobListFromContainer(cont)
	return list, err
}

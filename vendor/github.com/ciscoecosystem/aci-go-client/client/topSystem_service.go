package client

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) ReadSystem(pod int, node int) (*models.System, error) {
	dn := fmt.Sprintf("topology/pod-%d/node-%d/sys", pod, node)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}
	topSystem := models.SystemFromContainer(cont)
	return topSystem, nil
}

func (sm *ServiceManager) ListSystem() ([]*models.System, error) {
	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/topSystem.json", baseurlStr)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SystemListFromContainer(cont)
	return list, err
}

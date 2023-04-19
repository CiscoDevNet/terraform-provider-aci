package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/models"
)

func (sm *ServiceManager) CreateSubnetPoolforIpSecTunnels(subnetpool string, infra_network_template string, tenant string, cloudtemplateIpSecTunnelSubnetPoolAttr models.SubnetPoolforIpSecTunnelsAttributes) (*models.SubnetPoolforIpSecTunnels, error) {
	rn := fmt.Sprintf(models.RncloudtemplateIpSecTunnelSubnetPool, subnetpool)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateIpSecTunnelSubnetPool, tenant, infra_network_template)
	cloudtemplateIpSecTunnelSubnetPool := models.NewSubnetPoolforIpSecTunnels(rn, parentDn, cloudtemplateIpSecTunnelSubnetPoolAttr)
	err := sm.Save(cloudtemplateIpSecTunnelSubnetPool)
	return cloudtemplateIpSecTunnelSubnetPool, err
}

func (sm *ServiceManager) ReadSubnetPoolforIpSecTunnels(subnetpool string, infra_network_template string, tenant string) (*models.SubnetPoolforIpSecTunnels, error) {
	dn := fmt.Sprintf(models.DncloudtemplateIpSecTunnelSubnetPool, tenant, infra_network_template, subnetpool)

	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudtemplateIpSecTunnelSubnetPool := models.SubnetPoolforIpSecTunnelsFromContainer(cont)
	return cloudtemplateIpSecTunnelSubnetPool, nil
}

func (sm *ServiceManager) DeleteSubnetPoolforIpSecTunnels(subnetpool string, infra_network_template string, tenant string) error {
	dn := fmt.Sprintf(models.DncloudtemplateIpSecTunnelSubnetPool, tenant, infra_network_template, subnetpool)
	return sm.DeleteByDn(dn, models.CloudtemplateipsectunnelsubnetpoolClassName)
}

func (sm *ServiceManager) UpdateSubnetPoolforIpSecTunnels(subnetpool string, infra_network_template string, tenant string, cloudtemplateIpSecTunnelSubnetPoolAttr models.SubnetPoolforIpSecTunnelsAttributes) (*models.SubnetPoolforIpSecTunnels, error) {
	rn := fmt.Sprintf(models.RncloudtemplateIpSecTunnelSubnetPool, subnetpool)
	parentDn := fmt.Sprintf(models.ParentDncloudtemplateIpSecTunnelSubnetPool, tenant, infra_network_template)
	cloudtemplateIpSecTunnelSubnetPool := models.NewSubnetPoolforIpSecTunnels(rn, parentDn, cloudtemplateIpSecTunnelSubnetPoolAttr)
	cloudtemplateIpSecTunnelSubnetPool.Status = "modified"
	err := sm.Save(cloudtemplateIpSecTunnelSubnetPool)
	return cloudtemplateIpSecTunnelSubnetPool, err
}

func (sm *ServiceManager) ListSubnetPoolforIpSecTunnels(infra_network_template string, tenant string) ([]*models.SubnetPoolforIpSecTunnels, error) {
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/infranetwork-%s/cloudtemplateIpSecTunnelSubnetPool.json", models.BaseurlStr, tenant, infra_network_template)
	cont, err := sm.GetViaURL(dnUrl)
	list := models.SubnetPoolforIpSecTunnelsListFromContainer(cont)
	return list, err
}

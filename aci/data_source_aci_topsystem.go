package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSystem() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciSystemRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{

			"system_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"pod_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"etep_addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"node_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"remote_network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"remote_node": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rldirect_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"server_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"boot_strap_state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"child_action": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"config_issues": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"control_plane_mtu": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"current_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enforce_subnet_check": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"fabric_domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"fabric_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"fabric_mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"inb_mgmt_addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"inb_mgmt_addr6": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"inb_mgmt_addr6_mask": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"inb_mgmt_addr_mask": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"inb_mgmt_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"inb_mgmt_gateway6": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_reboot_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"last_reset_reason": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"lc_own": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mod_ts": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mon_pol_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"oob_mgmt_addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"oob_mgmt_addr6": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"oob_mgmt_addr6_mask": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"oob_mgmt_addr_mask": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"oob_mgmt_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"oob_mgmt_gateway6": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rl_oper_pod_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rl_routable_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"serial": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"site_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"system_uptime": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tep_pool": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"unicast_xr_ep_learn_disable": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"virtual_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAciSystemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	pod := d.Get("pod_id").(string)
	node := d.Get("system_id").(string)
	dn := fmt.Sprintf("topology/pod-%s/node-%s/sys", pod, node)

	topSystem, err := getRemoteSystem(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setSystemAttributes(topSystem, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getRemoteSystem(client *client.Client, dn string) (*models.System, error) {
	topSystemCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	topSystem := models.SystemFromContainer(topSystemCont)

	if topSystem.DistinguishedName == "" {
		return nil, fmt.Errorf("System %s not found", topSystem.DistinguishedName)
	}

	return topSystem, nil
}

func setSystemAttributes(topSystem *models.System, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(topSystem.DistinguishedName)

	topSystemMap, err := topSystem.ToMap()

	if err != nil {
		return d, err
	}
	d.Set("address", topSystemMap["address"])
	d.Set("etep_addr", topSystemMap["etepAddr"])
	d.Set("system_id", topSystemMap["id"])
	d.Set("name_alias", topSystemMap["nameAlias"])
	d.Set("node_type", topSystemMap["nodeType"])
	d.Set("remote_network_id", topSystemMap["remoteNetworkId"])
	d.Set("remote_node", topSystemMap["remoteNode"])
	d.Set("rldirect_mode", topSystemMap["rldirectMode"])
	d.Set("role", topSystemMap["role"])
	d.Set("server_type", topSystemMap["serverType"])
	d.Set("boot_strap_state", topSystemMap["bootstrapState"])
	d.Set("child_action", topSystemMap["childAction"])
	d.Set("config_issues", topSystemMap["configIssues"])
	d.Set("control_plane_mtu", topSystemMap["controlPlaneMTU"])
	d.Set("current_time", topSystemMap["currentTime"])
	d.Set("enforce_subnet_check", topSystemMap["enforceSubnetCheck"])
	d.Set("fabric_domain", topSystemMap["fabricDomain"])
	d.Set("fabric_id", topSystemMap["fabricId"])
	d.Set("fabric_mac", topSystemMap["fabricMAC"])
	d.Set("inb_mgmt_addr", topSystemMap["inbMgmtAddr"])
	d.Set("inb_mgmt_addr6", topSystemMap["inbMgmtAddr6"])
	d.Set("inb_mgmt_addr6_mask", topSystemMap["inbMgmtAddr6Mask"])
	d.Set("inb_mgmt_addr_mask", topSystemMap["inbMgmtAddrMask"])
	d.Set("inb_mgmt_gateway", topSystemMap["inbMgmtGateway"])
	d.Set("inb_mgmt_gateway6", topSystemMap["inbMgmtGateway6"])
	d.Set("last_reboot_time", topSystemMap["lastRebootTime"])
	d.Set("last_reset_reason", topSystemMap["lastResetReason"])
	d.Set("lc_own", topSystemMap["lcOwn"])
	d.Set("mod_ts", topSystemMap["modTs"])
	d.Set("mode", topSystemMap["mode"])
	d.Set("mon_pol_dn", topSystemMap["monPolDn"])
	d.Set("name", topSystemMap["name"])
	d.Set("oob_mgmt_addr", topSystemMap["oobMgmtAddr"])
	d.Set("oob_mgmt_addr6", topSystemMap["oobMgmtAddr6"])
	d.Set("oob_mgmt_addr6_mask", topSystemMap["oobMgmtAddr6Mask"])
	d.Set("oob_mgmt_addr_mask", topSystemMap["oobMgmtAddrMask"])
	d.Set("oob_mgmt_gateway", topSystemMap["oobMgmtGateway"])
	d.Set("oob_mgmt_gateway6", topSystemMap["oobMgmtGateway6"])
	d.Set("pod_id", topSystemMap["podId"])
	d.Set("rl_oper_pod_id", topSystemMap["rlOperPodId"])
	d.Set("rl_routable_mode", topSystemMap["rlRoutableMode"])
	d.Set("serial", topSystemMap["serial"])
	d.Set("site_id", topSystemMap["siteId"])
	d.Set("state", topSystemMap["state"])
	d.Set("system_uptime", topSystemMap["systemUpTime"])
	d.Set("tep_pool", topSystemMap["tepPool"])
	d.Set("unicast_xr_ep_learn_disable", topSystemMap["unicastXrEpLearnDisable"])
	d.Set("version", topSystemMap["version"])
	d.Set("virtual_mode", topSystemMap["virtualMode"])

	return d, nil
}

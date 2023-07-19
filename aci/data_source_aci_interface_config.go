package aci

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciInterfaceConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciInterfaceConfigurationRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"node": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(101, 4000),
			},
			"interface": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateInterface,
			},
			"port_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"access",
					"fabric",
				}, false),
				Default: "access",
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"breakout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pc_member": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operational_associated_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operational_associated_sub_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_dn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pc_port_dn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciInterfaceConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	node := strconv.Itoa(d.Get("node").(int))
	portType := d.Get("port_type").(string)
	parsedInterface := strings.Split(getInterfaceVal(d.Get("interface")), "/")
	rn := fmt.Sprintf(models.RnInfraPortConfig, node, parsedInterface[0], parsedInterface[1], parsedInterface[2])
	var dn string
	if portType == "access" {
		dn = fmt.Sprintf("%s/%s", models.ParentDnInfraPortConfig, rn)
	} else if portType == "fabric" {
		dn = fmt.Sprintf("%s/%s", models.ParentDnFabricPortConfig, rn)
	}

	err := getAndSetRemoteInterfaceConfiguration(dn, aciClient, d)
	if err != nil {
		return nil
	}
	d.SetId(dn)
	return nil
}

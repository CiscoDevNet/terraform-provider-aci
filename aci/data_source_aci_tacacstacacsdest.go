package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciTACACSDestination() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciTACACSDestinationRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tacacs_accounting_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "49",
			},
		})),
	}
}

func dataSourceAciTACACSDestinationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	host := d.Get("host").(string)
	port := d.Get("port").(string)
	TACACSMonitoringDestinationGroupDn := d.Get("tacacs_accounting_dn").(string)
	rn := fmt.Sprintf("tacacsdest-%s-port-%s", host, port)
	dn := fmt.Sprintf("%s/%s", TACACSMonitoringDestinationGroupDn, rn)
	tacacsTacacsDest, err := getRemoteTACACSDestination(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setTACACSDestinationAttributes(tacacsTacacsDest, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

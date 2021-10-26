package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciTriggerScheduler() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciTriggerSchedulerRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciTriggerSchedulerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("fabric/schedp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	trigSchedP, err := getRemoteTriggerScheduler(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setTriggerSchedulerAttributes(trigSchedP, d)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

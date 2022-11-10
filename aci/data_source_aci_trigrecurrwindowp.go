package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciRecurringWindow() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciRecurringWindowRead,
		SchemaVersion: 1,
		Schema: AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"scheduler_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"concur_cap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"day": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hour": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"minute": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"node_upg_interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proc_break": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proc_cap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"time_cap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Default:  "orchestrator:terraform",
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
			},
		}),
	}
}

func dataSourceAciRecurringWindowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	SchedulerDn := d.Get("scheduler_dn").(string)
	rn := fmt.Sprintf("recurrwinp-%s", name)
	dn := fmt.Sprintf("%s/%s", SchedulerDn, rn)
	trigRecurrWindowP, err := getRemoteRecurringWindow(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setRecurringWindowAttributes(trigRecurrWindowP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

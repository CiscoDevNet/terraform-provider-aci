package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudExternalEPg() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciCloudExternalEPgRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_applicationcontainer_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"exception_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flood_on_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"route_reachability": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciCloudExternalEPgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("cloudextepg-%s", name)
	CloudApplicationcontainerDn := d.Get("cloud_applicationcontainer_dn").(string)

	dn := fmt.Sprintf("%s/%s", CloudApplicationcontainerDn, rn)

	cloudExtEPg, err := getRemoteCloudExternalEPg(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setCloudExternalEPgAttributes(cloudExtEPg, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

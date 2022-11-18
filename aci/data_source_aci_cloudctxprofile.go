package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudContextProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciCloudContextProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "component type",
			},
			"primary_cidr": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Primary CIDR block",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "region",
			},
			"cloud_vendor": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the vendor",
			},
			"relation_cloud_rs_ctx_to_flow_log": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to cloudAwsFlowLogPol",
			},
			"relation_cloud_rs_to_ctx": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to fvCtx",
			},
			"hub_network": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "hub network to enable transit gateway",
			},
			"cloud_brownfield": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Import Brownfield Virtual Network",
			},
			"access_policy_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cloud context access policy type",
			},
		}),
	}
}

func dataSourceAciCloudContextProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)
	rn := fmt.Sprintf(models.RncloudCtxProfile, name)
	TenantDn := d.Get("tenant_dn").(string)

	dn := fmt.Sprintf("%s/%s", TenantDn, rn)
	log.Printf("[DEBUG] %s: Data Source - Beginning Read", dn)

	cloudCtxProfile, err := getRemoteCloudContextProfile(aciClient, dn, d)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setCloudContextProfileAttributes(cloudCtxProfile, d)

	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Data Source - Read finished successfully", dn)
	return nil
}

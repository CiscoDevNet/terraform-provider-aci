package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudServiceEPg() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciCloudServiceEPgRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"cloud_application_container_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azure_private_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flood_on_encap": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"label_match_criteria": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"preferred_group_member": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prio": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_service_epg_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"relation_cloud_rs_cloud_epg_ctx": {
				Type: schema.TypeString,

				Computed:    true,
				Description: "Query fv:Ctx relationship object",
			},
			"relation_fv_rs_cons": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Query vzBrCP relationship object",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prio": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"target_dn": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"relation_fv_rs_cons_if": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Query vzCPIf relationship object",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prio": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"target_dn": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"relation_fv_rs_cust_qos_pol": {
				Type: schema.TypeString,

				Computed:    true,
				Description: "Query qos:CustomPol relationship object",
			},
			"relation_fv_rs_graph_def": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Query vz:GraphCont relationship object",
				Set:         schema.HashString,
			},
			"relation_fv_rs_intra_epg": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Query vz:BrCP relationship object",
				Set:         schema.HashString,
			},
			"relation_fv_rs_prot_by": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Query vz:Taboo relationship object",
				Set:         schema.HashString,
			},
			"relation_fv_rs_prov": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Query vzBrCP relationship object",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label_match_criteria": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"prio": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"target_dn": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"relation_fv_rs_prov_def": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Query vz:CtrctEPgCont relationship object",
				Set:         schema.HashString,
			},
			"relation_fv_rs_sec_inherited": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Query fv:EPg relationship object",
				Set:         schema.HashString,
			}})),
	}
}

func dataSourceAciCloudServiceEPgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	cloudApplicationcontainerDn := d.Get("cloud_application_container_dn").(string)
	rn := fmt.Sprintf(models.RnCloudSvcEPg, name)
	dn := fmt.Sprintf("%s/%s", cloudApplicationcontainerDn, rn)

	cloudSvcEPg, err := getRemoteCloudServiceEPg(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setCloudServiceEPgAttributes(cloudSvcEPg, d)
	if err != nil {
		return nil
	}

	// Get and Set Relational Attributes
	getAndSetCloudServiceEPgRelationalAttributes(aciClient, dn, d)
	return nil
}

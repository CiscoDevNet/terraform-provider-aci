package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciL4ToL7Devices() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciL4ToL7DevicesRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"context_aware": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"function_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_copy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"promiscuous_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"trunking": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_vns_rs_al_dev_to_dom_p": {
				Type:          schema.TypeSet,
				Optional:      true,
				Description:   "Create relation to vmmDomP",
				MaxItems:      1,
				ConflictsWith: []string{"relation_vns_rs_al_dev_to_phys_dom_p"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_dn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"switching_mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"AVE",
								"native",
							}, false),
						},
					},
				},
			},
			"relation_vns_rs_al_dev_to_phys_dom_p": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to phys:DomP",
			}})),
	}
}

func dataSourceAciL4ToL7DevicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(models.RnvnsLDevVip, name)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	vnsLDevVip, err := getRemoteL4ToL7Devices(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setL4ToL7DevicesAttributes(vnsLDevVip, d)
	if err != nil {
		return diag.FromErr(err)
	}

	vnsRsALDevToDomPData, err := aciClient.ReadRelationvnsRsALDevToDomP(dn)
	relParams := make([]map[string]string, 0, 1)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsALDevToDomP %v", err)
		d.Set("relation_vns_rs_al_dev_to_dom_p", relParams)
	} else {
		relParamsList := vnsRsALDevToDomPData.([]map[string]string)
		for _, obj := range relParamsList {
			relParams = append(relParams, map[string]string{
				"switching_mode": obj["switchingMode"],
				"domain_dn":      obj["tDn"],
			})
		}
		d.Set("relation_vns_rs_al_dev_to_dom_p", relParams)
	}

	vnsRsALDevToPhysDomPData, err := aciClient.ReadRelationvnsRsALDevToPhysDomP(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsALDevToPhysDomP %v", err)
		d.Set("relation_vns_rs_al_dev_to_phys_dom_p", "")
	} else {
		d.Set("relation_vns_rs_al_dev_to_phys_dom_p", vnsRsALDevToPhysDomPData.(string))
	}

	return nil
}

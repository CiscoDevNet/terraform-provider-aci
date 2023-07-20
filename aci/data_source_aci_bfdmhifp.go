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

func dataSourceAciBfdMultihopInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciBfdMultihopInterfaceProfileRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"key_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface_profile_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_bfd_rs_mh_if_pol": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create relation to bfd:MhIfPol",
			},
		})),
	}
}

func dataSourceAciBfdMultihopInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Begining data source.")
	aciClient := m.(*client.Client)
	logicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)
	dn := fmt.Sprintf("%s/%s", logicalInterfaceProfileDn, models.RnbfdMhIfP)

	bfdMhIfP, err := getRemoteAciBfdMultihopInterfaceProfile(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setAciBfdMultihopInterfaceProfileAttributes(bfdMhIfP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	getAndSetRelationAciBfdMultihopInterfacePolicy(aciClient, d.Id(), d)

	return nil
}

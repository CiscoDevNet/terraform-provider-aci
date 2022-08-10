package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAttachableAccessEntityProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciAttachableAccessEntityProfileRead,

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
			"relation_infra_rs_dom_p": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}

func dataSourceAciAttachableAccessEntityProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/attentp-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)
	log.Printf("[DEBUG] %s: Data Source - Beginning Read", dn)

	infraAttEntityP, err := getRemoteAttachableAccessEntityProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setAttachableAccessEntityProfileAttributes(infraAttEntityP, d)

	if err != nil {
		return diag.FromErr(err)
	}

	// infraRsDomP - Beginning Read
	log.Printf("[DEBUG] %s: infraRsDomP - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsDomPFromAttachableAccessEntityProfile(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsDomP - Read finished successfully", d.Get("relation_infra_rs_dom_p"))
	}
	// infraRsDomP - Read finished successfully

	log.Printf("[DEBUG] %s: Data Source - Read finished successfully", dn)
	return nil
}

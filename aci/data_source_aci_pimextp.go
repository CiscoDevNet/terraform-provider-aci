package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciExternalProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciExternalProfileRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"l3outside_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled_af": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciExternalProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	L3OutsideDn := d.Get("l3outside_dn").(string)
	rn := fmt.Sprintf(models.RnPimExtP)
	dn := fmt.Sprintf("%s/%s", L3OutsideDn, rn)

	pimExtP, err := getRemoteExternalProfile(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setExternalProfileAttributes(pimExtP, d)
	if err != nil {
		return nil
	}

	return nil
}

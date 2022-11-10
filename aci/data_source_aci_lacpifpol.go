package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLACPMemberPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciLACPMemberPolicyRead,
		SchemaVersion: 1,
		Schema: AppendAttrSchemas(map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"transmit_rate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}, GetBaseAttrSchema(), GetNameAliasAttrSchema()),
	}
}

func dataSourceAciLACPMemberPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	dn := fmt.Sprintf(models.DnlacpIfPol, name)

	lacpIfPol, err := getRemoteLACPMemberPolicy(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setLACPMemberPolicyAttributes(lacpIfPol, d)
	if err != nil {
		return nil
	}

	return nil
}

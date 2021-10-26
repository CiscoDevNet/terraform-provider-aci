package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciTACACSSource() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciTACACSSourceReadContext,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "uni/fabric/moncommon",
			},
			"incl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"min_sev": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciTACACSSourceReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	ParentDn := d.Get("parent_dn").(string)
	rn := fmt.Sprintf("tacacssrc-%s", name)
	dn := fmt.Sprintf("%s/%s", ParentDn, rn)
	tacacsSrc, err := getRemoteTACACSSource(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setTACACSSourceAttributes(tacacsSrc, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

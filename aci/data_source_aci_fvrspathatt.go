package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciStaticPath() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciStaticPathRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instr_imedcy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"primary_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciStaticPathRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	rn := fmt.Sprintf("rspathAtt-[%s]", tDn)
	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	dn := fmt.Sprintf("%s/%s", ApplicationEPGDn, rn)

	fvRsPathAtt, err := getRemoteStaticPath(aciClient, dn)

	if err != nil {
		return err
	}
	setStaticPathAttributes(fvRsPathAtt, d)
	return nil
}

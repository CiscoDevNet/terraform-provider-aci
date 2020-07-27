package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciEPGsUsingFunction() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciEPGsUsingFunctionRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"access_generic_dn": &schema.Schema{
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

func dataSourceAciEPGsUsingFunctionRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	rn := fmt.Sprintf("rsfuncToEpg-[%s]", tDn)
	AccessGenericDn := d.Get("access_generic_dn").(string)

	dn := fmt.Sprintf("%s/%s", AccessGenericDn, rn)

	infraRsFuncToEpg, err := getRemoteEPGsUsingFunction(aciClient, dn)

	if err != nil {
		return err
	}
	setEPGsUsingFunctionAttributes(infraRsFuncToEpg, d)
	return nil
}

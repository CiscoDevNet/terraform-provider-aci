package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAciNodeManagementEPg() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciNodeManagementEPgRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in_band",
					"out_of_band",
				}, false),
			},

			"management_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
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

			"exception_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flood_on_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciNodeManagementEPgRead(d *schema.ResourceData, m interface{}) error {
	if d.Get("type").(string) == "in_band" {
		return dataSourceInBandManagementEPgRead(d, m)
	}
	return dataSourceOutOfBandManagementEPgRead(d, m)
}

func dataSourceInBandManagementEPgRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("inb-%s", name)
	ManagementProfileDn := d.Get("management_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", ManagementProfileDn, rn)

	mgmtInB, err := getRemoteInBandManagementEPg(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setInBandManagementEPgAttributes(mgmtInB, d)
	return nil
}

func dataSourceOutOfBandManagementEPgRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("oob-%s", name)
	ManagementProfileDn := d.Get("management_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", ManagementProfileDn, rn)

	mgmtOoB, err := getRemoteOutOfBandManagementEPg(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setOutOfBandManagementEPgAttributes(mgmtOoB, d)
	return nil
}

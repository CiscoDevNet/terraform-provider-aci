package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciCloudEndpointSelector() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciCloudEndpointSelectorRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"match_expression": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciCloudEndpointSelectorRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("epselector-%s", name)
	CloudEPgDn := d.Get("cloud_epg_dn").(string)

	dn := fmt.Sprintf("%s/%s", CloudEPgDn, rn)

	cloudEPSelector, err := getRemoteCloudEndpointSelector(aciClient, dn)

	if err != nil {
		return err
	}
	setCloudEndpointSelectorAttributes(cloudEPSelector, d)
	return nil
}

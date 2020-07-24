package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciAccessGeneric() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciAccessGenericRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"attachable_access_entity_profile_dn": &schema.Schema{
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

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciAccessGenericRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("gen-%s", name)
	AttachableAccessEntityProfileDn := d.Get("attachable_access_entity_profile_dn").(string)

	dn := fmt.Sprintf("%s/%s", AttachableAccessEntityProfileDn, rn)

	infraGeneric, err := getRemoteAccessGeneric(aciClient, dn)

	if err != nil {
		return err
	}
	setAccessGenericAttributes(infraGeneric, d)
	return nil
}

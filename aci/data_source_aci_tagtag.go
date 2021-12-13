package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciTag() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciTagRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"parent_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciTagRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	key := d.Get("key").(string)
	parentDn := d.Get("parent_dn").(string)
	rn := fmt.Sprintf(models.RnTagTag, key)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)
	tagTag, err := getRemoteTag(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setTagAttributes(tagTag, d)
	return nil
}

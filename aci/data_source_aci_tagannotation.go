package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAnnotation() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciAnnotationRead,
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

func dataSourceAciAnnotationRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	key := d.Get("key").(string)
	parentDn := d.Get("parent_dn").(string)
	rn := fmt.Sprintf(models.RnTagAnnotation, key)
	dn := fmt.Sprintf("%s/%s", parentDn, rn)
	tagAnnotation, err := getRemoteAnnotation(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setAnnotationAttributes(tagAnnotation, d)
	return nil
}

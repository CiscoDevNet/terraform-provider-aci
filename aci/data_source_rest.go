package aci

import (
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceAciRest() *schema.Resource {
	return &schema.Resource{
		Read: datasourceAciRestRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"class_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"content": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			"payload": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func datasourceAciRestRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Rest data source: Beginning Read")

	aciClient := m.(*client.Client)

	path := d.Get("path").(string)

	cont, err := aciClient.GetViaURL(path)
	if err != nil {
		return err
	}

	payloadData := cont.S("imdata").Index(0)
	for k, _ := range payloadData.Data().(map[string]interface{}) {
		d.Set("class_name", k)
	}

	dn := stripQuotes(payloadData.S(d.Get("class_name").(string), "attributes", "dn").String())

	contentMap := payloadData.S(d.Get("class_name").(string), "attributes").Data().(map[string]interface{})
	d.Set("content", contentMap)

	d.Set("payload", payloadData.String())
	d.Set("dn", dn)
	d.SetId(dn)

	log.Println("[DEBUG] Rest data source: Ending Read ", d.Id())
	return nil
}

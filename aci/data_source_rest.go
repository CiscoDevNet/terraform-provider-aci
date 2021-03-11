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
				Required: true,
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

	aciClient := m.(*client.Client)

	path := d.Get("path").(string)

	cont, err := aciClient.GetViaURL(path)
	if err != nil {
		return err
	}

	payloadData := cont.S("imdata").Index(0)

	log.Println("check data : ", payloadData)

	return nil
}

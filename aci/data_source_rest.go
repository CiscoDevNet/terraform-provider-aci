package aci

import (
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/container"
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

			"children": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"child_class_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"child_content": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"content": &schema.Schema{
				Type:     schema.TypeMap,
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

	if payloadData.Exists(d.Get("class_name").(string), "children") {
		childrenSet, err := getChildrenAttrs(aciClient, payloadData.S(d.Get("class_name").(string), "children"))
		if err != nil {
			return err
		}

		d.Set("children", childrenSet)
	} else {
		d.Set("children", make([]interface{}, 0, 1))
	}

	d.Set("dn", dn)
	d.SetId(dn)

	log.Println("[DEBUG] Rest data source: Ending Read ", d.Id())
	return nil
}

func getChildrenAttrs(client *client.Client, cont *container.Container) ([]interface{}, error) {
	childCount, err := cont.ArrayCount()
	if err != nil {
		return nil, err
	}

	childrenSet := make([]interface{}, 0, 1)
	for i := 0; i < childCount; i++ {
		childCont := cont.Index(i)

		tpChildMap := make(map[string]interface{})
		for k, _ := range childCont.Data().(map[string]interface{}) {
			tpChildMap["child_class_name"] = k
		}

		childContentMap := childCont.S(tpChildMap["child_class_name"].(string), "attributes").Data().(map[string]interface{})
		tpChildMap["child_content"] = childContentMap

		childrenSet = append(childrenSet, tpChildMap)
	}

	return childrenSet, nil
}

package provider

import (
	"context"
	"log"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciRest() *schema.Resource {
	return &schema.Resource{
		Description: "This data source can read one ACI object and its children.",

		ReadContext: dataSourceAciRestRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The distinguished name of the object.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dn": {
				Type:        schema.TypeString,
				Description: "Distinguished name of object to be retrieved, e.g. uni/tn-EXAMPLE_TENANT.",
				Required:    true,
			},
			"class_name": {
				Type:        schema.TypeString,
				Description: "Class name of object being retrieved.",
				Computed:    true,
			},
			"content": {
				Type:        schema.TypeMap,
				Description: "Map of key-value pairs which represents the attributes of object being retrieved.",
				Computed:    true,
			},
			"child": {
				Type:        schema.TypeSet,
				Description: "Set of children of object being retrieved.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"class_name": {
							Type:        schema.TypeString,
							Description: "Class name of child object being retrieved.",
							Computed:    true,
						},
						"content": {
							Type:        schema.TypeMap,
							Description: "Map of key-value pairs which represents the attributes of child object being retrieved.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAciRestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	for attempts := 0; ; attempts++ {
		cont, diags := ApicRest(d, meta, "GET", true)
		if diags.HasError() {
			if ok := backoff(attempts, meta.(apiClient).Retries); !ok {
				return diags
			}
			log.Printf("[ERROR] Failed to read object: %s, retries: %v", diags[0].Summary, attempts)
			continue
		}

		// Check if we received an empty response without errors -> object has been deleted
		if cont == nil && diags == nil {
			d.SetId("")
			break
		}

		data, err := cont.Search("imdata").Index(0).ChildrenMap()
		if err != nil {
			return diag.FromErr(err)
		}

		var className string
		var obj *container.Container
		for k, o := range data {
			className = k
			obj = o
		}
		// Set class_name
		d.Set("class_name", className)

		// Set content
		contentMap := obj.Search("attributes").Data().(map[string]interface{})
		d.Set("content", contentMap)

		if obj.Exists("children") {
			children := obj.Search("children")
			childCount, err := children.ArrayCount()
			if err != nil {
				if ok := backoff(attempts, meta.(apiClient).Retries); !ok {
					return diag.FromErr(err)
				}
				log.Printf("[ERROR] Failed to decode response after reading object: %s, retries: %v", diags[0].Summary, attempts)
				continue
			}

			childrenSet := make([]interface{}, 0, 1)
			for i := 0; i < childCount; i++ {
				childMap := make(map[string]interface{})
				childData, err := children.Index(i).ChildrenMap()
				if err != nil {
					return diag.FromErr(err)
				}

				var childClassName string
				var childObj *container.Container
				for cn, o := range childData {
					childClassName = cn
					childObj = o
				}
				// Set child class_name
				childMap["class_name"] = childClassName

				// Set child content
				childContentMap := childObj.Search("attributes").Data().(map[string]interface{})
				childMap["content"] = childContentMap

				childrenSet = append(childrenSet, childMap)
			}

			d.Set("child", childrenSet)
		} else {
			d.Set("child", make([]interface{}, 0, 1))
		}

		// Set id
		d.SetId(d.Get("dn").(string))
		break
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

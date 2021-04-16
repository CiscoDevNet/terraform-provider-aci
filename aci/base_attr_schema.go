package aci

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func GetBaseAttrSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		"annotation": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			DefaultFunc: func() (interface{}, error) {
				return "orchestrator:terraform", nil
			},
		},
	}
}

// AppendBaseAttrSchema adds the BaseAttr to any schema
func AppendBaseAttrSchema(attrs map[string]*schema.Schema) map[string]*schema.Schema {
	for key, value := range GetBaseAttrSchema() {
		attrs[key] = value
	}

	return attrs
}

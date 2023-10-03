package aci

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetBaseAttrSchema() map[string]*schema.Schema {
	return AppendAttrSchemas(map[string]*schema.Schema{
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}, GetAnnotationAttrSchema())
}

func GetAllowEmptyAttrSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_empty_result": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
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

// AppendAttrSchemas adds a range of schemas to any schema
func AppendAttrSchemas(attrs map[string]*schema.Schema, mapsToAppend ...map[string]*schema.Schema) map[string]*schema.Schema {
	for _, m := range mapsToAppend {
		for key, value := range m {
			attrs[key] = value
		}
	}
	return attrs
}

func GetAnnotationAttrSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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

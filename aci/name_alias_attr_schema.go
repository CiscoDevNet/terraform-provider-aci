package aci

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetNameAliasAttrSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name_alias": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

// AppendNameAliasAttrSchema adds the NameAliasAttr to required schema
func AppendNameAliasAttrSchema(attrs map[string]*schema.Schema) map[string]*schema.Schema {
	for key, value := range GetNameAliasAttrSchema() {
		attrs[key] = value
	}
	return attrs
}

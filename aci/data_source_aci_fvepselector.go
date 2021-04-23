package aci


import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciEndpointSecurityGroupSelector() *schema.Resource {
	return &schema.Resource{

		Read:   dataSourceAciEndpointSecurityGroupSelectorRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"endpoint_security_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			
			
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			
            
			
			"match_expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			
            
			
			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			
            
			
			"userdom": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			
            
            }),
    }
}



func dataSourceAciEndpointSecurityGroupSelectorRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	
	matchExpression := d.Get("matchExpression").(string)
	

    rn := fmt.Sprintf("epselector-[%s]",matchExpression)
    EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)

    dn := fmt.Sprintf("%s/%s",EndpointSecurityGroupDn,rn)
    

	
	fvEPSelector, err := getRemoteEndpointSecurityGroupSelector(aciClient, dn)

	if err != nil {
		return err
	}
	setEndpointSecurityGroupSelectorAttributes(fvEPSelector, d)
	return nil
}
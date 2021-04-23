package aci


import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciEndpointSecurityGroup() *schema.Resource {
	return &schema.Resource{

		Read:   dataSourceAciEndpointSecurityGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			
			
			"name": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			
			
            
			
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			
            
			
			"exception_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			
            
			
			"flood_on_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			
            
			
			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			
            
			
			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			
            
			
			"pc_enf_pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			
            
			
			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			
            
			
			"prio": &schema.Schema{
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



func dataSourceAciEndpointSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	
	name := d.Get("name").(string)
	

    rn := fmt.Sprintf("esg-%s",name)
    ApplicationProfileDn := d.Get("application_profile_dn").(string)

    dn := fmt.Sprintf("%s/%s",ApplicationProfileDn,rn)
    

	
	fvESg, err := getRemoteEndpointSecurityGroup(aciClient, dn)

	if err != nil {
		return err
	}
	setEndpointSecurityGroupAttributes(fvESg, d)
	return nil
}
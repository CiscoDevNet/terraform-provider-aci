package aci


import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciEndpointSecurityGroupSelector() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciEndpointSecurityGroupSelectorCreate,
		Update: resourceAciEndpointSecurityGroupSelectorUpdate,
		Read:   resourceAciEndpointSecurityGroupSelectorRead,
		Delete: resourceAciEndpointSecurityGroupSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciEndpointSecurityGroupSelectorImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"endpoint_security_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			
			"matchExpression": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			


	     
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
			     
				Optional: true,
				Computed: true,
				Description: "Mo doc not defined in techpub!!!",
                
		},
             
			"match_expression": &schema.Schema{
				Type:     schema.TypeString,
			     
				Optional: true,
				Computed: true,
				Description: "Mo doc not defined in techpub!!!",
                
		},
             
			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
			     
				Optional: true,
				Computed: true,
				Description: "Mo doc not defined in techpub!!!",
                
		},
             
			"userdom": &schema.Schema{
				Type:     schema.TypeString,
			     
				Optional: true,
				Computed: true,
				Description: "Mo doc not defined in techpub!!!",
                
		},
            	    

			

		}),
	}
}

func getRemoteEndpointSecurityGroupSelector(client *client.Client, dn string) (*models.EndpointSecurityGroupSelector, error) {
	fvEPSelectorCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvEPSelector := models.EndpointSecurityGroupSelectorFromContainer(fvEPSelectorCont)

	if fvEPSelector.DistinguishedName == "" {
		return nil, fmt.Errorf("Bridge Domain %s not found", fvEPSelector.DistinguishedName)
	}

	return fvEPSelector, nil
}

func setEndpointSecurityGroupSelectorAttributes(fvEPSelector *models.EndpointSecurityGroupSelector, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvEPSelector.DistinguishedName)
	d.Set("description", fvEPSelector.Description)
	d.Set("endpoint_security_group_dn", GetParentDn(fvEPSelector.DistinguishedName))
	fvEPSelectorMap , _ := fvEPSelector.ToMap()
     
	d.Set("annotation", fvEPSelectorMap["annotation"]) 
	d.Set("match_expression", fvEPSelectorMap["matchExpression"]) 
	d.Set("name_alias", fvEPSelectorMap["nameAlias"]) 
	d.Set("userdom", fvEPSelectorMap["userdom"])
	return d
}

func resourceAciEndpointSecurityGroupSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	fvEPSelector, err := getRemoteEndpointSecurityGroupSelector(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setEndpointSecurityGroupSelectorAttributes(fvEPSelector, d)
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciEndpointSecurityGroupSelectorCreate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	matchExpression := d.Get("matchExpression").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)
	
	fvEPSelectorAttr := models.EndpointSecurityGroupSelectorAttributes{} 
    if Annotation, ok := d.GetOk("annotation"); ok {
        fvEPSelectorAttr.Annotation  = Annotation.(string)
    } 
    if MatchExpression, ok := d.GetOk("match_expression"); ok {
        fvEPSelectorAttr.MatchExpression  = MatchExpression.(string)
    } 
    if NameAlias, ok := d.GetOk("name_alias"); ok {
        fvEPSelectorAttr.NameAlias  = NameAlias.(string)
    } 
    if Userdom, ok := d.GetOk("userdom"); ok {
        fvEPSelectorAttr.Userdom  = Userdom.(string)
    }
	fvEPSelector := models.NewEndpointSecurityGroupSelector(fmt.Sprintf("epselector-[%s]",matchExpression),EndpointSecurityGroupDn, desc, fvEPSelectorAttr)  
	
	
	err := aciClient.Save(fvEPSelector)
	if err != nil {
		return err
	}
	

	d.SetId(fvEPSelector.DistinguishedName)
	return resourceAciEndpointSecurityGroupSelectorRead(d, m)
}

func resourceAciEndpointSecurityGroupSelectorUpdate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	
	matchExpression := d.Get("matchExpression").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)
	

    fvEPSelectorAttr := models.EndpointSecurityGroupSelectorAttributes{}     
    if Annotation, ok := d.GetOk("annotation"); ok {
        fvEPSelectorAttr.Annotation = Annotation.(string)
    }     
    if MatchExpression, ok := d.GetOk("match_expression"); ok {
        fvEPSelectorAttr.MatchExpression = MatchExpression.(string)
    }     
    if NameAlias, ok := d.GetOk("name_alias"); ok {
        fvEPSelectorAttr.NameAlias = NameAlias.(string)
    }     
    if Userdom, ok := d.GetOk("userdom"); ok {
        fvEPSelectorAttr.Userdom = Userdom.(string)
    }
	fvEPSelector := models.NewEndpointSecurityGroupSelector(fmt.Sprintf("epselector-[%s]",matchExpression),EndpointSecurityGroupDn, desc, fvEPSelectorAttr)  
		

	fvEPSelector.Status = "modified"

	err := aciClient.Save(fvEPSelector)
	
	if err != nil {
		return err
	}
	

	d.SetId(fvEPSelector.DistinguishedName)
	return resourceAciEndpointSecurityGroupSelectorRead(d, m)

}

func resourceAciEndpointSecurityGroupSelectorRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	dn := d.Id()
	fvEPSelector, err := getRemoteEndpointSecurityGroupSelector(aciClient, dn)

	if err != nil {
		return err
	}
	setEndpointSecurityGroupSelectorAttributes(fvEPSelector, d)
	return nil
}

func resourceAciEndpointSecurityGroupSelectorDelete(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvEPSelector")
	if err != nil {
		return err
	}

	d.SetId("")
	return err
}
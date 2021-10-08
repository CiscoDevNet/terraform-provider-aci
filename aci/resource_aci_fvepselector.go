package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciEndpointSecurityGroupSelector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciEndpointSecurityGroupSelectorCreate,
		UpdateContext: resourceAciEndpointSecurityGroupSelectorUpdate,
		ReadContext:   resourceAciEndpointSecurityGroupSelectorRead,
		DeleteContext: resourceAciEndpointSecurityGroupSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciEndpointSecurityGroupSelectorImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"endpoint_security_group_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"match_expression": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteEndpointSecurityGroupSelector(client *client.Client, dn string) (*models.EndpointSecurityGroupSelector, error) {
	fvEPSelectorCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvEPSelector := models.EndpointSecurityGroupSelectorFromContainer(fvEPSelectorCont)
	if fvEPSelector.DistinguishedName == "" {
		return nil, fmt.Errorf("EndpointSecurityGroupSelector %s not found", fvEPSelector.DistinguishedName)
	}
	return fvEPSelector, nil
}

func setEndpointSecurityGroupSelectorAttributes(fvEPSelector *models.EndpointSecurityGroupSelector, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fvEPSelector.DistinguishedName)
	d.Set("description", fvEPSelector.Description)
	fvEPSelectorMap, err := fvEPSelector.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("endpoint_security_group_dn", GetParentDn(fvEPSelector.DistinguishedName, fmt.Sprintf("/epselector-[%s]", fvEPSelectorMap["matchExpression"])))
	d.Set("annotation", fvEPSelectorMap["annotation"])
	d.Set("match_expression", fvEPSelectorMap["matchExpression"])
	d.Set("name", fvEPSelectorMap["name"])
	d.Set("name_alias", fvEPSelectorMap["nameAlias"])
	return d, nil
}

func resourceAciEndpointSecurityGroupSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvEPSelector, err := getRemoteEndpointSecurityGroupSelector(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setEndpointSecurityGroupSelectorAttributes(fvEPSelector, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciEndpointSecurityGroupSelectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EndpointSecurityGroupSelector: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	matchExpression := d.Get("match_expression").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)

	fvEPSelectorAttr := models.EndpointSecurityGroupSelectorAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvEPSelectorAttr.Annotation = Annotation.(string)
	} else {
		fvEPSelectorAttr.Annotation = "{}"
	}

	if MatchExpression, ok := d.GetOk("match_expression"); ok {
		fvEPSelectorAttr.MatchExpression = MatchExpression.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvEPSelectorAttr.Name = Name.(string)
	}
	fvEPSelector := models.NewEndpointSecurityGroupSelector(fmt.Sprintf("epselector-[%s]", matchExpression), EndpointSecurityGroupDn, desc, nameAlias, fvEPSelectorAttr)

	err := aciClient.Save(fvEPSelector)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvEPSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciEndpointSecurityGroupSelectorRead(ctx, d, m)
}

func resourceAciEndpointSecurityGroupSelectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EndpointSecurityGroupSelector: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	matchExpression := d.Get("match_expression").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)
	fvEPSelectorAttr := models.EndpointSecurityGroupSelectorAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvEPSelectorAttr.Annotation = Annotation.(string)
	} else {
		fvEPSelectorAttr.Annotation = "{}"
	}

	if MatchExpression, ok := d.GetOk("match_expression"); ok {
		fvEPSelectorAttr.MatchExpression = MatchExpression.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvEPSelectorAttr.Name = Name.(string)
	}
	fvEPSelector := models.NewEndpointSecurityGroupSelector(fmt.Sprintf("epselector-[%s]", matchExpression), EndpointSecurityGroupDn, desc, nameAlias, fvEPSelectorAttr)

	fvEPSelector.Status = "modified"
	err := aciClient.Save(fvEPSelector)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvEPSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciEndpointSecurityGroupSelectorRead(ctx, d, m)
}

func resourceAciEndpointSecurityGroupSelectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvEPSelector, err := getRemoteEndpointSecurityGroupSelector(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setEndpointSecurityGroupSelectorAttributes(fvEPSelector, d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciEndpointSecurityGroupSelectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvEPSelector")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

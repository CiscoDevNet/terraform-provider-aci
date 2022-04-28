package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciEndpointSecurityGroupTagSelector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciEndpointSecurityGroupTagSelectorCreate,
		UpdateContext: resourceAciEndpointSecurityGroupTagSelectorUpdate,
		ReadContext:   resourceAciEndpointSecurityGroupTagSelectorRead,
		DeleteContext: resourceAciEndpointSecurityGroupTagSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciEndpointSecurityGroupTagSelectorImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"endpoint_security_group_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"match_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"match_value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"value_operator": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"contains",
					"equals",
					"regex",
				}, false),
			},
		})),
	}
}

func getRemoteEndpointSecurityGroupTagSelector(client *client.Client, dn string) (*models.EndpointSecurityGroupTagSelector, error) {
	fvTagSelectorCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvTagSelector := models.EndpointSecurityGroupTagSelectorFromContainer(fvTagSelectorCont)
	if fvTagSelector.DistinguishedName == "" {
		return nil, fmt.Errorf("EndpointSecurityGroupTagSelector %s not found", fvTagSelector.DistinguishedName)
	}
	return fvTagSelector, nil
}

func setEndpointSecurityGroupTagSelectorAttributes(fvTagSelector *models.EndpointSecurityGroupTagSelector, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fvTagSelector.DistinguishedName)
	d.Set("description", fvTagSelector.Description)
	fvTagSelectorMap, err := fvTagSelector.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", fvTagSelectorMap["annotation"])
	d.Set("match_key", fvTagSelectorMap["matchKey"])
	d.Set("match_value", fvTagSelectorMap["matchValue"])
	d.Set("name", fvTagSelectorMap["name"])
	d.Set("value_operator", fvTagSelectorMap["valueOperator"])
	d.Set("name_alias", fvTagSelectorMap["nameAlias"])
	return d, nil
}

func resourceAciEndpointSecurityGroupTagSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvTagSelector, err := getRemoteEndpointSecurityGroupTagSelector(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setEndpointSecurityGroupTagSelectorAttributes(fvTagSelector, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciEndpointSecurityGroupTagSelectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EndpointSecurityGroupTagSelector: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	matchKey := d.Get("match_key").(string)
	matchValue := d.Get("match_value").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)

	fvTagSelectorAttr := models.EndpointSecurityGroupTagSelectorAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvTagSelectorAttr.Annotation = Annotation.(string)
	} else {
		fvTagSelectorAttr.Annotation = "{}"
	}

	if MatchKey, ok := d.GetOk("match_key"); ok {
		fvTagSelectorAttr.MatchKey = MatchKey.(string)
	}

	if MatchValue, ok := d.GetOk("match_value"); ok {
		fvTagSelectorAttr.MatchValue = MatchValue.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvTagSelectorAttr.Name = Name.(string)
	}

	if ValueOperator, ok := d.GetOk("value_operator"); ok {
		fvTagSelectorAttr.ValueOperator = ValueOperator.(string)
	}
	fvTagSelector := models.NewEndpointSecurityGroupTagSelector(fmt.Sprintf("tagselectorkey-[%s]-value-[%s]", matchKey, matchValue), EndpointSecurityGroupDn, desc, nameAlias, fvTagSelectorAttr)

	err := aciClient.Save(fvTagSelector)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvTagSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciEndpointSecurityGroupTagSelectorRead(ctx, d, m)
}

func resourceAciEndpointSecurityGroupTagSelectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EndpointSecurityGroupTagSelector: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	matchKey := d.Get("match_key").(string)
	matchValue := d.Get("match_value").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)
	fvTagSelectorAttr := models.EndpointSecurityGroupTagSelectorAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvTagSelectorAttr.Annotation = Annotation.(string)
	} else {
		fvTagSelectorAttr.Annotation = "{}"
	}

	if MatchKey, ok := d.GetOk("match_key"); ok {
		fvTagSelectorAttr.MatchKey = MatchKey.(string)
	}

	if MatchValue, ok := d.GetOk("match_value"); ok {
		fvTagSelectorAttr.MatchValue = MatchValue.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvTagSelectorAttr.Name = Name.(string)
	}

	if ValueOperator, ok := d.GetOk("value_operator"); ok {
		fvTagSelectorAttr.ValueOperator = ValueOperator.(string)
	}
	fvTagSelector := models.NewEndpointSecurityGroupTagSelector(fmt.Sprintf("tagselectorkey-[%s]-value-[%s]", matchKey, matchValue), EndpointSecurityGroupDn, desc, nameAlias, fvTagSelectorAttr)

	fvTagSelector.Status = "modified"
	err := aciClient.Save(fvTagSelector)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvTagSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciEndpointSecurityGroupTagSelectorRead(ctx, d, m)
}

func resourceAciEndpointSecurityGroupTagSelectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvTagSelector, err := getRemoteEndpointSecurityGroupTagSelector(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setEndpointSecurityGroupTagSelectorAttributes(fvTagSelector, d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciEndpointSecurityGroupTagSelectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvTagSelector")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

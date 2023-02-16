package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciEndpointSecurityGroupEPgSelector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciEndpointSecurityGroupEPgSelectorCreate,
		UpdateContext: resourceAciEndpointSecurityGroupEPgSelectorUpdate,
		ReadContext:   resourceAciEndpointSecurityGroupEPgSelectorRead,
		DeleteContext: resourceAciEndpointSecurityGroupEPgSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciEndpointSecurityGroupEPgSelectorImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"endpoint_security_group_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"match_epg_dn": {
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

func getRemoteEndpointSecurityGroupEPgSelector(client *client.Client, dn string) (*models.EndpointSecurityGroupEPgSelector, error) {
	fvEPgSelectorCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvEPgSelector := models.EndpointSecurityGroupEPgSelectorFromContainer(fvEPgSelectorCont)
	if fvEPgSelector.DistinguishedName == "" {
		return nil, fmt.Errorf("Endpoint Security Group EPG Selector %s not found", dn)
	}
	return fvEPgSelector, nil
}

func setEndpointSecurityGroupEPgSelectorAttributes(fvEPgSelector *models.EndpointSecurityGroupEPgSelector, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fvEPgSelector.DistinguishedName)
	d.Set("description", fvEPgSelector.Description)
	fvEPgSelectorMap, err := fvEPgSelector.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", fvEPgSelectorMap["annotation"])
	d.Set("match_epg_dn", fvEPgSelectorMap["matchEpgDn"])
	d.Set("name", fvEPgSelectorMap["name"])
	d.Set("name_alias", fvEPgSelectorMap["nameAlias"])
	return d, nil
}

func resourceAciEndpointSecurityGroupEPgSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvEPgSelector, err := getRemoteEndpointSecurityGroupEPgSelector(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setEndpointSecurityGroupEPgSelectorAttributes(fvEPgSelector, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciEndpointSecurityGroupEPgSelectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EndpointSecurityGroupEPgSelector: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	matchEpgDn := d.Get("match_epg_dn").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)

	fvEPgSelectorAttr := models.EndpointSecurityGroupEPgSelectorAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvEPgSelectorAttr.Annotation = Annotation.(string)
	} else {
		fvEPgSelectorAttr.Annotation = "{}"
	}

	if MatchEpgDn, ok := d.GetOk("match_epg_dn"); ok {
		fvEPgSelectorAttr.MatchEpgDn = MatchEpgDn.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvEPgSelectorAttr.Name = Name.(string)
	}
	fvEPgSelector := models.NewEndpointSecurityGroupEPgSelector(fmt.Sprintf("epgselector-[%s]", matchEpgDn), EndpointSecurityGroupDn, desc, nameAlias, fvEPgSelectorAttr)

	err := aciClient.Save(fvEPgSelector)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fvEPgSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciEndpointSecurityGroupEPgSelectorRead(ctx, d, m)
}

func resourceAciEndpointSecurityGroupEPgSelectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EndpointSecurityGroupEPgSelector: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	matchEpgDn := d.Get("match_epg_dn").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)
	fvEPgSelectorAttr := models.EndpointSecurityGroupEPgSelectorAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvEPgSelectorAttr.Annotation = Annotation.(string)
	} else {
		fvEPgSelectorAttr.Annotation = "{}"
	}

	if MatchEpgDn, ok := d.GetOk("match_epg_dn"); ok {
		fvEPgSelectorAttr.MatchEpgDn = MatchEpgDn.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvEPgSelectorAttr.Name = Name.(string)
	}
	fvEPgSelector := models.NewEndpointSecurityGroupEPgSelector(fmt.Sprintf("epgselector-[%s]", matchEpgDn), EndpointSecurityGroupDn, desc, nameAlias, fvEPgSelectorAttr)

	fvEPgSelector.Status = "modified"
	err := aciClient.Save(fvEPgSelector)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fvEPgSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciEndpointSecurityGroupEPgSelectorRead(ctx, d, m)
}

func resourceAciEndpointSecurityGroupEPgSelectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvEPgSelector, err := getRemoteEndpointSecurityGroupEPgSelector(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setEndpointSecurityGroupEPgSelectorAttributes(fvEPgSelector, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciEndpointSecurityGroupEPgSelectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvEPgSelector")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

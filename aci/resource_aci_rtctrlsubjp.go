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

func resourceAciMatchRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMatchRuleCreate,
		UpdateContext: resourceAciMatchRuleUpdate,
		ReadContext:   resourceAciMatchRuleRead,
		DeleteContext: resourceAciMatchRuleDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMatchRuleImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteMatchRule(client *client.Client, dn string) (*models.MatchRule, error) {
	rtctrlSubjPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSubjP := models.MatchRuleFromContainer(rtctrlSubjPCont)
	if rtctrlSubjP.DistinguishedName == "" {
		return nil, fmt.Errorf("MatchRule %s not found", dn)
	}
	return rtctrlSubjP, nil
}

func setMatchRuleAttributes(rtctrlSubjP *models.MatchRule, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(rtctrlSubjP.DistinguishedName)
	d.Set("description", rtctrlSubjP.Description)
	rtctrlSubjPMap, err := rtctrlSubjP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", rtctrlSubjPMap["annotation"])
	d.Set("name", rtctrlSubjPMap["name"])
	d.Set("name_alias", rtctrlSubjPMap["nameAlias"])
	return d, nil
}

func resourceAciMatchRuleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	rtctrlSubjP, err := getRemoteMatchRule(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setMatchRuleAttributes(rtctrlSubjP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMatchRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MatchRule: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	rtctrlSubjPAttr := models.MatchRuleAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlSubjPAttr.Annotation = Annotation.(string)
	} else {
		rtctrlSubjPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlSubjPAttr.Name = Name.(string)
	}
	rtctrlSubjP := models.NewMatchRule(fmt.Sprintf("subj-%s", name), TenantDn, desc, nameAlias, rtctrlSubjPAttr)

	err := aciClient.Save(rtctrlSubjP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlSubjP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciMatchRuleRead(ctx, d, m)
}

func resourceAciMatchRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MatchRule: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rtctrlSubjPAttr := models.MatchRuleAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlSubjPAttr.Annotation = Annotation.(string)
	} else {
		rtctrlSubjPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlSubjPAttr.Name = Name.(string)
	}
	rtctrlSubjP := models.NewMatchRule(fmt.Sprintf("subj-%s", name), TenantDn, desc, nameAlias, rtctrlSubjPAttr)

	rtctrlSubjP.Status = "modified"
	err := aciClient.Save(rtctrlSubjP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlSubjP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciMatchRuleRead(ctx, d, m)
}

func resourceAciMatchRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	rtctrlSubjP, err := getRemoteMatchRule(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setMatchRuleAttributes(rtctrlSubjP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciMatchRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "rtctrlSubjP")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

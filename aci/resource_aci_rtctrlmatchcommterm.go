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

func resourceAciMatchCommunityTerm() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMatchCommunityTermCreate,
		UpdateContext: resourceAciMatchCommunityTermUpdate,
		ReadContext:   resourceAciMatchCommunityTermRead,
		DeleteContext: resourceAciMatchCommunityTermDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMatchCommunityTermImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"match_rule_dn": {
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

func getRemoteMatchCommunityTerm(client *client.Client, dn string) (*models.MatchCommunityTerm, error) {
	rtctrlMatchCommTermCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlMatchCommTerm := models.MatchCommunityTermFromContainer(rtctrlMatchCommTermCont)
	if rtctrlMatchCommTerm.DistinguishedName == "" {
		return nil, fmt.Errorf("MatchCommunityTerm %s not found", dn)
	}
	return rtctrlMatchCommTerm, nil
}

func setMatchCommunityTermAttributes(rtctrlMatchCommTerm *models.MatchCommunityTerm, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(rtctrlMatchCommTerm.DistinguishedName)
	d.Set("description", rtctrlMatchCommTerm.Description)
	if dn != rtctrlMatchCommTerm.DistinguishedName {
		d.Set("match_rule_dn", "")
	}
	rtctrlMatchCommTermMap, err := rtctrlMatchCommTerm.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("match_rule_dn", GetParentDn(dn, fmt.Sprintf("/"+models.RnrtctrlMatchCommTerm, rtctrlMatchCommTermMap["match_rule_dn"])))
	d.Set("annotation", rtctrlMatchCommTermMap["annotation"])
	d.Set("name", rtctrlMatchCommTermMap["name"])
	d.Set("name_alias", rtctrlMatchCommTermMap["nameAlias"])
	return d, nil
}

func resourceAciMatchCommunityTermImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	rtctrlMatchCommTerm, err := getRemoteMatchCommunityTerm(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setMatchCommunityTermAttributes(rtctrlMatchCommTerm, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMatchCommunityTermCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MatchCommunityTerm: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	matchRuleDn := d.Get("match_rule_dn").(string)

	rtctrlMatchCommTermAttr := models.MatchCommunityTermAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlMatchCommTermAttr.Annotation = Annotation.(string)
	} else {
		rtctrlMatchCommTermAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlMatchCommTermAttr.Name = Name.(string)
	}
	rtctrlMatchCommTerm := models.NewMatchCommunityTerm(fmt.Sprintf(models.RnrtctrlMatchCommTerm, name), matchRuleDn, desc, nameAlias, rtctrlMatchCommTermAttr)

	err := aciClient.Save(rtctrlMatchCommTerm)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlMatchCommTerm.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciMatchCommunityTermRead(ctx, d, m)
}

func resourceAciMatchCommunityTermUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MatchCommunityTerm: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	matchRuleDn := d.Get("match_rule_dn").(string)

	rtctrlMatchCommTermAttr := models.MatchCommunityTermAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlMatchCommTermAttr.Annotation = Annotation.(string)
	} else {
		rtctrlMatchCommTermAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlMatchCommTermAttr.Name = Name.(string)
	}
	rtctrlMatchCommTerm := models.NewMatchCommunityTerm(fmt.Sprintf(models.RnrtctrlMatchCommTerm, name), matchRuleDn, desc, nameAlias, rtctrlMatchCommTermAttr)

	rtctrlMatchCommTerm.Status = "modified"

	err := aciClient.Save(rtctrlMatchCommTerm)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlMatchCommTerm.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciMatchCommunityTermRead(ctx, d, m)
}

func resourceAciMatchCommunityTermRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	rtctrlMatchCommTerm, err := getRemoteMatchCommunityTerm(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setMatchCommunityTermAttributes(rtctrlMatchCommTerm, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciMatchCommunityTermDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "rtctrlMatchCommTerm")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

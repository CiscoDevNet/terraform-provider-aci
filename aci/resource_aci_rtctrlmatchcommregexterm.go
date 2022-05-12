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

func resourceAciMatchRuleBasedonCommunityRegularExpression() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMatchRuleBasedonCommunityRegularExpressionCreate,
		UpdateContext: resourceAciMatchRuleBasedonCommunityRegularExpressionUpdate,
		ReadContext:   resourceAciMatchRuleBasedonCommunityRegularExpressionRead,
		DeleteContext: resourceAciMatchRuleBasedonCommunityRegularExpressionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMatchRuleBasedonCommunityRegularExpressionImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"match_rule_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"community_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"extended",
					"regular",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"regex": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteMatchRuleBasedonCommunityRegularExpression(client *client.Client, dn string) (*models.MatchRuleBasedonCommunityRegularExpression, error) {
	rtctrlMatchCommRegexTermCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlMatchCommRegexTerm := models.MatchRuleBasedonCommunityRegularExpressionFromContainer(rtctrlMatchCommRegexTermCont)
	if rtctrlMatchCommRegexTerm.DistinguishedName == "" {
		return nil, fmt.Errorf("MatchRuleBasedonCommunityRegularExpression %s not found", dn)
	}
	return rtctrlMatchCommRegexTerm, nil
}

func setMatchRuleBasedonCommunityRegularExpressionAttributes(rtctrlMatchCommRegexTerm *models.MatchRuleBasedonCommunityRegularExpression, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(rtctrlMatchCommRegexTerm.DistinguishedName)
	d.Set("description", rtctrlMatchCommRegexTerm.Description)
	if dn != rtctrlMatchCommRegexTerm.DistinguishedName {
		d.Set("match_rule_dn", "")
	}
	rtctrlMatchCommRegexTermMap, err := rtctrlMatchCommRegexTerm.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("match_rule_dn", GetParentDn(dn, fmt.Sprintf("/"+models.RnrtctrlMatchCommRegexTerm, rtctrlMatchCommRegexTermMap["match_rule_dn"])))
	d.Set("annotation", rtctrlMatchCommRegexTermMap["annotation"])
	d.Set("community_type", rtctrlMatchCommRegexTermMap["commType"])
	d.Set("name", rtctrlMatchCommRegexTermMap["name"])
	d.Set("regex", rtctrlMatchCommRegexTermMap["regex"])
	d.Set("name_alias", rtctrlMatchCommRegexTermMap["nameAlias"])
	return d, nil
}

func resourceAciMatchRuleBasedonCommunityRegularExpressionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	rtctrlMatchCommRegexTerm, err := getRemoteMatchRuleBasedonCommunityRegularExpression(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setMatchRuleBasedonCommunityRegularExpressionAttributes(rtctrlMatchCommRegexTerm, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMatchRuleBasedonCommunityRegularExpressionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MatchRuleBasedonCommunityRegularExpression: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	commType := d.Get("community_type").(string)
	matchRuleDn := d.Get("match_rule_dn").(string)

	rtctrlMatchCommRegexTermAttr := models.MatchRuleBasedonCommunityRegularExpressionAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlMatchCommRegexTermAttr.Annotation = Annotation.(string)
	} else {
		rtctrlMatchCommRegexTermAttr.Annotation = "{}"
	}

	if CommType, ok := d.GetOk("community_type"); ok {
		rtctrlMatchCommRegexTermAttr.CommType = CommType.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlMatchCommRegexTermAttr.Name = Name.(string)
	}

	if Regex, ok := d.GetOk("regex"); ok {
		rtctrlMatchCommRegexTermAttr.Regex = Regex.(string)
	}
	rtctrlMatchCommRegexTerm := models.NewMatchRuleBasedonCommunityRegularExpression(fmt.Sprintf(models.RnrtctrlMatchCommRegexTerm, commType), matchRuleDn, desc, nameAlias, rtctrlMatchCommRegexTermAttr)

	err := aciClient.Save(rtctrlMatchCommRegexTerm)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlMatchCommRegexTerm.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciMatchRuleBasedonCommunityRegularExpressionRead(ctx, d, m)
}

func resourceAciMatchRuleBasedonCommunityRegularExpressionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MatchRuleBasedonCommunityRegularExpression: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	commType := d.Get("community_type").(string)
	matchRuleDn := d.Get("match_rule_dn").(string)

	rtctrlMatchCommRegexTermAttr := models.MatchRuleBasedonCommunityRegularExpressionAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlMatchCommRegexTermAttr.Annotation = Annotation.(string)
	} else {
		rtctrlMatchCommRegexTermAttr.Annotation = "{}"
	}

	if CommType, ok := d.GetOk("community_type"); ok {
		rtctrlMatchCommRegexTermAttr.CommType = CommType.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlMatchCommRegexTermAttr.Name = Name.(string)
	}

	if Regex, ok := d.GetOk("regex"); ok {
		rtctrlMatchCommRegexTermAttr.Regex = Regex.(string)
	}
	rtctrlMatchCommRegexTerm := models.NewMatchRuleBasedonCommunityRegularExpression(fmt.Sprintf(models.RnrtctrlMatchCommRegexTerm, commType), matchRuleDn, desc, nameAlias, rtctrlMatchCommRegexTermAttr)

	rtctrlMatchCommRegexTerm.Status = "modified"

	err := aciClient.Save(rtctrlMatchCommRegexTerm)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlMatchCommRegexTerm.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciMatchRuleBasedonCommunityRegularExpressionRead(ctx, d, m)
}

func resourceAciMatchRuleBasedonCommunityRegularExpressionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	rtctrlMatchCommRegexTerm, err := getRemoteMatchRuleBasedonCommunityRegularExpression(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setMatchRuleBasedonCommunityRegularExpressionAttributes(rtctrlMatchCommRegexTerm, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciMatchRuleBasedonCommunityRegularExpressionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "rtctrlMatchCommRegexTerm")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

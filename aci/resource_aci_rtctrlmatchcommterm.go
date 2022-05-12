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
			"match_community_factors": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Create Community Factors",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"transitive",
								"non-transitive",
							}, false),
						},
						"community": {
							Required: true,
							Type:     schema.TypeString,
						},
						"description": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
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

func getRemoteMatchCommunityFactor(client *client.Client, dn string) (*models.MatchCommunityFactor, error) {
	rtctrlMatchCommFactorCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlMatchCommFactor := models.MatchCommunityFactorFromContainer(rtctrlMatchCommFactorCont)
	if rtctrlMatchCommFactor.DistinguishedName == "" {
		return nil, fmt.Errorf("MatchCommunityFactor %s not found", rtctrlMatchCommFactor.DistinguishedName)
	}
	return rtctrlMatchCommFactor, nil
}

func setMatchCommunityFactorAttributes(rtctrlMatchCommFactor *models.MatchCommunityFactor, d map[string]string) (map[string]string, error) {
	rtctrlMatchCommFactorMap, err := rtctrlMatchCommFactor.ToMap()
	if err != nil {
		return d, err
	}

	d = map[string]string{
		"community":   rtctrlMatchCommFactorMap["community"],
		"scope":       rtctrlMatchCommFactorMap["scope"],
		"description": rtctrlMatchCommFactor.Description,
	}

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

	rtctrlMatchCommFactors, err := aciClient.ListMatchCommFactorsFromCommunityTerm(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation rtctrlMatchCommFactor %v", err)
	}

	matchCommFactors := make([]map[string]string, 0, 1)

	for _, factor := range rtctrlMatchCommFactors {

		factorSet, err := setMatchCommunityFactorAttributes(factor, make(map[string]string))
		if err != nil {
			return nil, err
		}
		matchCommFactors = append(matchCommFactors, factorSet)
	}
	d.Set("match_community_factors", matchCommFactors)

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

	if matchCommunityFactors, ok := d.GetOk("match_community_factors"); ok {
		factors := matchCommunityFactors.(*schema.Set).List()
		for _, factor := range factors {
			factorMap := factor.(map[string]interface{})

			rtctrlMatchCommFactorAttr := models.MatchCommunityFactorAttributes{}
			rtctrlMatchCommFactorAttr.Scope = factorMap["scope"].(string)
			rtctrlMatchCommFactorAttr.Community = factorMap["community"].(string)
			rtctrlMatchCommFactorAttr.Annotation = rtctrlMatchCommTerm.Annotation

			rtctrlMatchCommFactor := models.NewMatchCommunityFactor(fmt.Sprintf(models.RnrtctrlMatchCommFactor, rtctrlMatchCommFactorAttr.Community), rtctrlMatchCommTerm.DistinguishedName, factorMap["description"].(string), "", rtctrlMatchCommFactorAttr)

			err := aciClient.Save(rtctrlMatchCommFactor)
			if err != nil {
				return diag.FromErr(err)
			}
		}
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
	if d.HasChange("match_community_factors") || d.HasChange("annotation") {
		previousMatchCommunityFactors, matchCommunityFactors := d.GetChange("match_community_factors")

		oldFactors := previousMatchCommunityFactors.(*schema.Set).List()
		factors := matchCommunityFactors.(*schema.Set).List()
		for _, oldFactor := range oldFactors {
			found := false
			oldFactorMap := oldFactor.(map[string]interface{})
			for _, factor := range factors {
				factorMap := factor.(map[string]interface{})
				if factorMap["community"].(string) == oldFactorMap["community"].(string) {
					found = true
					break
				}
			}
			if !found {
				dn := rtctrlMatchCommTerm.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlMatchCommFactor, oldFactorMap["community"].(string))

				err := aciClient.DeleteByDn(dn, "rtctrlMatchCommFactor")
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		for _, factor := range factors {
			factorMap := factor.(map[string]interface{})

			rtctrlMatchCommFactorAttr := models.MatchCommunityFactorAttributes{}
			rtctrlMatchCommFactorAttr.Scope = factorMap["scope"].(string)
			rtctrlMatchCommFactorAttr.Community = factorMap["community"].(string)
			rtctrlMatchCommFactorAttr.Annotation = rtctrlMatchCommTerm.Annotation

			found := false
			changed := false

			for _, oldFactor := range oldFactors {
				oldFactorMap := oldFactor.(map[string]interface{})
				if factorMap["community"].(string) == oldFactorMap["community"].(string) {
					found = true
					if factorMap["scope"].(string) != oldFactorMap["scope"].(string) || factorMap["description"].(string) != oldFactorMap["description"].(string) {
						changed = true
					}
					break
				}
			}
			if !found || changed {
				rtctrlMatchCommFactor := models.NewMatchCommunityFactor(fmt.Sprintf(models.RnrtctrlMatchCommFactor, rtctrlMatchCommFactorAttr.Community), rtctrlMatchCommTerm.DistinguishedName, factorMap["description"].(string), "", rtctrlMatchCommFactorAttr)
				err := aciClient.Save(rtctrlMatchCommFactor)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
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
		return nil
	}

	_, err = setMatchCommunityTermAttributes(rtctrlMatchCommTerm, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	if matchCommunityFactors, ok := d.GetOk("match_community_factors"); ok {
		factors := matchCommunityFactors.(*schema.Set).List()
		matchCommFactors := make([]map[string]string, 0, 1)

		for _, factor := range factors {
			factorMap := factor.(map[string]interface{})
			dn := rtctrlMatchCommTerm.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlMatchCommFactor, factorMap["community"].(string))
			rtctrlMatchCommFactor, err := getRemoteMatchCommunityFactor(aciClient, dn)
			if err != nil {
				d.SetId("")
				return diag.FromErr(err)
			}
			factorSet, err := setMatchCommunityFactorAttributes(rtctrlMatchCommFactor, make(map[string]string))
			if err != nil {
				return diag.FromErr(err)
			}
			matchCommFactors = append(matchCommFactors, factorSet)
		}
		d.Set("match_community_factors", matchCommFactors)
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

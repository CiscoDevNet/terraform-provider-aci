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
				Description: "Create relation to Community Terms",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope": {
							Required: true,
							Type:     schema.TypeString,
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

func setMatchCommunityFactorAttributes(rtctrlMatchCommFactor *models.MatchCommunityFactor, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(rtctrlMatchCommFactor.DistinguishedName)
	d.Set("description", rtctrlMatchCommFactor.Description)
	if dn != rtctrlMatchCommFactor.DistinguishedName {
		d.Set("match_community_term_dn", "")
	}
	rtctrlMatchCommFactorMap, err := rtctrlMatchCommFactor.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", rtctrlMatchCommFactorMap["annotation"])
	d.Set("community", rtctrlMatchCommFactorMap["community"])
	d.Set("name", rtctrlMatchCommFactorMap["name"])
	d.Set("scope", rtctrlMatchCommFactorMap["scope"])
	d.Set("name_alias", rtctrlMatchCommFactorMap["nameAlias"])
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

	if matchCommunityFactors, ok := d.GetOk("match_community_factors"); ok {
		factors := matchCommunityFactors.(*schema.Set).List()
		for _, factor := range factors {
			factorMap := factor.(map[string]interface{})

			rtctrlMatchCommFactorAttr := models.MatchCommunityFactorAttributes{}
			rtctrlMatchCommFactorAttr.Scope = factorMap["scope"].(string)
			rtctrlMatchCommFactorAttr.Community = factorMap["community"].(string)
			rtctrlMatchCommFactorAttr.Annotation = rtctrlMatchCommTerm.Annotation

			rtctrlMatchCommFactor := models.NewMatchCommunityFactor(fmt.Sprintf(models.RnrtctrlMatchCommFactor, rtctrlMatchCommFactorAttr.Community), rtctrlMatchCommTerm.DistinguishedName, factorMap["description"].(string), nil, rtctrlMatchCommFactorAttr)
			// err = aciClient.CreateMatchCommunityFactor(factorMap["community"].(string), name, rtctrlMatchCommTerm.Annotation, factorMap["scope"].(string))
			// CreateMatchCommunityFactor(community string, match_community_term string, match_rule string, tenant string, description string, nameAlias string, rtctrlMatchCommFactorAttr models.MatchCommunityFactorAttributes)

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

	// if d.HasChange("match_community_factors") {
	// 	oldRel, newRel := d.GetChange("match_community_factors")
	// 	oldRelList := oldRel.(*schema.Set).List()
	// 	newRelList := newRel.(*schema.Set).List()
	// 	for _, relationParam := range oldRelList {
	// 		paramMap := relationParam.(map[string]interface{})
	// 		err = aciClient.DeleteMatchCommunityFactor(rtctrlMatchCommTerm.DistinguishedName, paramMap["community"].(string), paramMap["scope"].(string))

	// 		if err != nil {
	// 			return diag.FromErr(err)
	// 		}
	// 	}
	// 	for _, relationParam := range newRelList {
	// 		paramMap := relationParam.(map[string]interface{})
	// 		err = aciClient.CreateMatchCommunityFactor(rtctrlMatchCommTerm.DistinguishedName, rtctrlMatchCommTerm.Annotation, paramMap["community"].(string), paramMap["scope"].(string))

	// 		if err != nil {
	// 			return diag.FromErr(err)
	// 		}
	// 	}
	// }
	// CreateMatchCommunityFactor(community string, match_community_term string, match_rule string, tenant string, description string, nameAlias string, rtctrlMatchCommFactorAttr models.MatchCommunityFactorAttributes)

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

	// ReadMatchCommunityFactor(community string, match_community_term string, match_rule string, tenant string) and delete
	// matchCommunityFactorData, err := aciClient.ReadMatchCommunityFactor(dn)

	// if err != nil {
	// 	log.Printf("[DEBUG] Error while reading relation rtctrlMatchCommFactor %v", err)
	// 	d.Set("match_community_factors", make([]map[string]string, 0))

	// } else {
	// 	matchCommunityFactorMap := matchCommunityFactorData.([]map[string]string)
	// 	st := make([]map[string]string, 0, 1)
	// 	for _, matchCommunityFactorObj := range matchCommunityFactorMap {
	// 		st = append(st, map[string]string{
	// 			"community": matchCommunityFactorObj["community"],
	// 			"scope":     matchCommunityFactorObj["scope"],
	// 		})
	// 	}
	// 	d.Set("match_community_factors", st)
	// }

	if matchCommunityFactors, ok := d.GetOk("match_community_factors"); ok {
		factors := matchCommunityFactors.(*schema.Set).List()
		for _, factor := range factors {
			factorMap := factor.(map[string]interface{})
			dn := rtctrlMatchCommTerm.DistinguishedName + fmt.Sprintf(models.RnrtctrlMatchCommFactor, factorMap["community"].(string))

			// replace by get and set

			rtctrlMatchCommFactor, err := getRemoteMatchCommunityFactor(aciClient, dn)
			if err != nil {
				d.SetId("")
				return diag.FromErr(err)
			}

			_, err = setMatchCommunityFactorAttributes(rtctrlMatchCommFactor, d)
			if err != nil {
				d.SetId("")
				return nil
			}

			// aciClient.ReadMatchCommunityFactor(dn)

			// rtctrlMatchCommFactorAttr := models.MatchCommunityFactorAttributes{}
			// rtctrlMatchCommFactorAttr.Scope = factorMap["scope"].(string)
			// rtctrlMatchCommFactorAttr.Community = factorMap["community"].(string)
			// rtctrlMatchCommFactorAttr.Annotation = rtctrlMatchCommTerm.Annotation

			// rtctrlMatchCommFactor := models.NewMatchCommunityFactor(fmt.Sprintf(models.RnrtctrlMatchCommFactor, rtctrlMatchCommFactorAttr.Community), rtctrlMatchCommTerm.DistinguishedName, factorMap["description"].(string), nil, rtctrlMatchCommFactorAttr)
			// // err = aciClient.CreateMatchCommunityFactor(factorMap["community"].(string), name, rtctrlMatchCommTerm.Annotation, factorMap["scope"].(string))
			// // CreateMatchCommunityFactor(community string, match_community_term string, match_rule string, tenant string, description string, nameAlias string, rtctrlMatchCommFactorAttr models.MatchCommunityFactorAttributes)

			// err := aciClient.Save(rtctrlMatchCommFactor)
			// if err != nil {
			// 	return diag.FromErr(err)
			// }
		}
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

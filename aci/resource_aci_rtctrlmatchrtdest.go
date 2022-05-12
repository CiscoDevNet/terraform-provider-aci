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

func resourceAciMatchRouteDestinationRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMatchRouteDestinationRuleCreate,
		UpdateContext: resourceAciMatchRouteDestinationRuleUpdate,
		ReadContext:   resourceAciMatchRouteDestinationRuleRead,
		DeleteContext: resourceAciMatchRouteDestinationRuleDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMatchRouteDestinationRuleImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"match_rule_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aggregate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"greater_than_mask": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"less_than_mask": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteMatchRouteDestinationRule(client *client.Client, dn string) (*models.MatchRouteDestinationRule, error) {
	rtctrlMatchRtDestCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlMatchRtDest := models.MatchRouteDestinationRuleFromContainer(rtctrlMatchRtDestCont)
	if rtctrlMatchRtDest.DistinguishedName == "" {
		return nil, fmt.Errorf("MatchRouteDestinationRule %s not found", rtctrlMatchRtDest.DistinguishedName)
	}
	return rtctrlMatchRtDest, nil
}

func setMatchRouteDestinationRuleAttributes(rtctrlMatchRtDest *models.MatchRouteDestinationRule, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(rtctrlMatchRtDest.DistinguishedName)
	d.Set("description", rtctrlMatchRtDest.Description)
	rtctrlMatchRtDestMap, err := rtctrlMatchRtDest.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("aggregate", rtctrlMatchRtDestMap["aggregate"])
	d.Set("annotation", rtctrlMatchRtDestMap["annotation"])
	d.Set("greater_than_mask", rtctrlMatchRtDestMap["fromPfxLen"])
	d.Set("ip", rtctrlMatchRtDestMap["ip"])
	d.Set("name", rtctrlMatchRtDestMap["name"])
	d.Set("less_than_mask", rtctrlMatchRtDestMap["toPfxLen"])
	d.Set("name_alias", rtctrlMatchRtDestMap["nameAlias"])
	return d, nil
}

func resourceAciMatchRouteDestinationRuleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	rtctrlMatchRtDest, err := getRemoteMatchRouteDestinationRule(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setMatchRouteDestinationRuleAttributes(rtctrlMatchRtDest, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMatchRouteDestinationRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MatchRouteDestinationRule: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	ip := d.Get("ip").(string)
	MatchRuleDn := d.Get("match_rule_dn").(string)

	rtctrlMatchRtDestAttr := models.MatchRouteDestinationRuleAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlMatchRtDestAttr.Annotation = Annotation.(string)
	} else {
		rtctrlMatchRtDestAttr.Annotation = "{}"
	}

	if Aggregate, ok := d.GetOk("aggregate"); ok {
		rtctrlMatchRtDestAttr.Aggregate = Aggregate.(string)
	}

	if FromPfxLen, ok := d.GetOk("greater_than_mask"); ok {
		rtctrlMatchRtDestAttr.FromPfxLen = FromPfxLen.(string)
	}

	if Ip, ok := d.GetOk("ip"); ok {
		rtctrlMatchRtDestAttr.Ip = Ip.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlMatchRtDestAttr.Name = Name.(string)
	}

	if ToPfxLen, ok := d.GetOk("less_than_mask"); ok {
		rtctrlMatchRtDestAttr.ToPfxLen = ToPfxLen.(string)
	}
	rtctrlMatchRtDest := models.NewMatchRouteDestinationRule(fmt.Sprintf("dest-[%s]", ip), MatchRuleDn, desc, nameAlias, rtctrlMatchRtDestAttr)

	err := aciClient.Save(rtctrlMatchRtDest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlMatchRtDest.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciMatchRouteDestinationRuleRead(ctx, d, m)
}

func resourceAciMatchRouteDestinationRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MatchRouteDestinationRule: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	ip := d.Get("ip").(string)
	MatchRuleDn := d.Get("match_rule_dn").(string)
	rtctrlMatchRtDestAttr := models.MatchRouteDestinationRuleAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlMatchRtDestAttr.Annotation = Annotation.(string)
	} else {
		rtctrlMatchRtDestAttr.Annotation = "{}"
	}

	if Aggregate, ok := d.GetOk("aggregate"); ok {
		rtctrlMatchRtDestAttr.Aggregate = Aggregate.(string)
	}

	if FromPfxLen, ok := d.GetOk("greater_than_mask"); ok {
		rtctrlMatchRtDestAttr.FromPfxLen = FromPfxLen.(string)
	}

	if Ip, ok := d.GetOk("ip"); ok {
		rtctrlMatchRtDestAttr.Ip = Ip.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlMatchRtDestAttr.Name = Name.(string)
	}

	if ToPfxLen, ok := d.GetOk("less_than_mask"); ok {
		rtctrlMatchRtDestAttr.ToPfxLen = ToPfxLen.(string)
	}
	rtctrlMatchRtDest := models.NewMatchRouteDestinationRule(fmt.Sprintf("dest-[%s]", ip), MatchRuleDn, desc, nameAlias, rtctrlMatchRtDestAttr)

	rtctrlMatchRtDest.Status = "modified"
	err := aciClient.Save(rtctrlMatchRtDest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlMatchRtDest.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciMatchRouteDestinationRuleRead(ctx, d, m)
}

func resourceAciMatchRouteDestinationRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	rtctrlMatchRtDest, err := getRemoteMatchRouteDestinationRule(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setMatchRouteDestinationRuleAttributes(rtctrlMatchRtDest, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciMatchRouteDestinationRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "rtctrlMatchRtDest")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

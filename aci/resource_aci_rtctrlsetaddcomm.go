package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciRtctrlSetAddComm() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciRtctrlSetAddCommCreate,
		UpdateContext: resourceAciRtctrlSetAddCommUpdate,
		ReadContext:   resourceAciRtctrlSetAddCommRead,
		DeleteContext: resourceAciRtctrlSetAddCommDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciRtctrlSetAddCommImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"action_rule_profile_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"community": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"set_criteria": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"append",
				}, false),
			},
		})),
	}
}

func getRemoteRtctrlSetAddComm(client *client.Client, dn string) (*models.RtctrlSetAddComm, error) {
	rtctrlSetAddCommCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetAddComm := models.RtctrlSetAddCommFromContainer(rtctrlSetAddCommCont)
	if rtctrlSetAddComm.DistinguishedName == "" {
		return nil, fmt.Errorf("RtctrlSetAddComm %s not found", dn)
	}
	return rtctrlSetAddComm, nil
}

func setRtctrlSetAddCommAttributes(rtctrlSetAddComm *models.RtctrlSetAddComm, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(rtctrlSetAddComm.DistinguishedName)
	d.Set("description", rtctrlSetAddComm.Description)
	if dn != rtctrlSetAddComm.DistinguishedName {
		d.Set("action_rule_profile_dn", "")
	}
	rtctrlSetAddCommMap, err := rtctrlSetAddComm.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("action_rule_profile_dn", GetParentDn(dn, fmt.Sprintf("/saddcomm-%s", rtctrlSetAddCommMap["community"])))
	d.Set("annotation", rtctrlSetAddCommMap["annotation"])
	d.Set("community", rtctrlSetAddCommMap["community"])
	d.Set("set_criteria", rtctrlSetAddCommMap["setCriteria"])
	d.Set("name_alias", rtctrlSetAddCommMap["nameAlias"])
	return d, nil
}

func resourceAciRtctrlSetAddCommImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	rtctrlSetAddComm, err := getRemoteRtctrlSetAddComm(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setRtctrlSetAddCommAttributes(rtctrlSetAddComm, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciRtctrlSetAddCommCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RtctrlSetAddComm: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	community := d.Get("community").(string)
	ActionRuleProfileDn := d.Get("action_rule_profile_dn").(string)

	rtctrlSetAddCommAttr := models.RtctrlSetAddCommAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlSetAddCommAttr.Annotation = Annotation.(string)
	} else {
		rtctrlSetAddCommAttr.Annotation = "{}"
	}

	if Community, ok := d.GetOk("community"); ok {
		rtctrlSetAddCommAttr.Community = Community.(string)
	}

	if SetCriteria, ok := d.GetOk("set_criteria"); ok {
		rtctrlSetAddCommAttr.SetCriteria = SetCriteria.(string)
	}

	rtctrlSetAddComm := models.NewRtctrlSetAddComm(fmt.Sprintf(models.RnrtctrlSetAddComm, community), ActionRuleProfileDn, desc, nameAlias, rtctrlSetAddCommAttr)

	err := aciClient.Save(rtctrlSetAddComm)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlSetAddComm.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciRtctrlSetAddCommRead(ctx, d, m)
}

func resourceAciRtctrlSetAddCommUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RtctrlSetAddComm: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	community := d.Get("community").(string)
	ActionRuleProfileDn := d.Get("action_rule_profile_dn").(string)

	rtctrlSetAddCommAttr := models.RtctrlSetAddCommAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlSetAddCommAttr.Annotation = Annotation.(string)
	} else {
		rtctrlSetAddCommAttr.Annotation = "{}"
	}

	if Community, ok := d.GetOk("community"); ok {
		rtctrlSetAddCommAttr.Community = Community.(string)
	}

	if SetCriteria, ok := d.GetOk("set_criteria"); ok {
		rtctrlSetAddCommAttr.SetCriteria = SetCriteria.(string)
	}

	rtctrlSetAddComm := models.NewRtctrlSetAddComm(fmt.Sprintf(models.RnrtctrlSetAddComm, community), ActionRuleProfileDn, desc, nameAlias, rtctrlSetAddCommAttr)

	rtctrlSetAddComm.Status = "modified"

	err := aciClient.Save(rtctrlSetAddComm)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlSetAddComm.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciRtctrlSetAddCommRead(ctx, d, m)
}

func resourceAciRtctrlSetAddCommRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	rtctrlSetAddComm, err := getRemoteRtctrlSetAddComm(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setRtctrlSetAddCommAttributes(rtctrlSetAddComm, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciRtctrlSetAddCommDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "rtctrlSetAddComm")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

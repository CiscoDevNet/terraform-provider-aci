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

func resourceAciActionRuleProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciActionRuleProfileCreate,
		UpdateContext: resourceAciActionRuleProfileUpdate,
		ReadContext:   resourceAciActionRuleProfileRead,
		DeleteContext: resourceAciActionRuleProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciActionRuleProfileImport,
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
			"set_route_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
		})),
	}
}
func getRemoteActionRuleProfile(client *client.Client, dn string) (*models.ActionRuleProfile, error) {
	rtctrlAttrPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlAttrP := models.ActionRuleProfileFromContainer(rtctrlAttrPCont)

	if rtctrlAttrP.DistinguishedName == "" {
		return nil, fmt.Errorf("ActionRuleProfile %s not found", dn)
	}

	return rtctrlAttrP, nil
}

func setActionRuleProfileAttributes(rtctrlAttrP *models.ActionRuleProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(rtctrlAttrP.DistinguishedName)
	d.Set("description", rtctrlAttrP.Description)
	if dn != rtctrlAttrP.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	rtctrlAttrPMap, err := rtctrlAttrP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", rtctrlAttrPMap["name"])
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/attr-%s", rtctrlAttrPMap["name"])))
	d.Set("annotation", rtctrlAttrPMap["annotation"])
	d.Set("name_alias", rtctrlAttrPMap["nameAlias"])
	return d, nil
}

func getRemoteRtctrlSetTag(client *client.Client, dn string) (*models.RtctrlSetTag, error) {
	rtctrlSetTagCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetTag := models.RtctrlSetTagFromContainer(rtctrlSetTagCont)
	if rtctrlSetTag.DistinguishedName == "" {
		return nil, fmt.Errorf("RtctrlSetTag %s not found", dn)
	}
	return rtctrlSetTag, nil
}

func setRtctrlSetTagAttributes(rtctrlSetTag *models.RtctrlSetTag, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetTagMap, err := rtctrlSetTag.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("set_route_tag", rtctrlSetTagMap["tag"])
	return d, nil
}

func resourceAciActionRuleProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	rtctrlAttrP, err := getRemoteActionRuleProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	rtctrlAttrPMap, err := rtctrlAttrP.ToMap()
	if err != nil {
		return nil, err
	}
	name := rtctrlAttrPMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/attr-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setActionRuleProfileAttributes(rtctrlAttrP, d)
	if err != nil {
		return nil, err
	}

	// rtctrlSetTag - Beginning Import

	setRouteTagDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetTag)
	rtctrlSetTag, err := getRemoteRtctrlSetTag(aciClient, setRouteTagDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetTag - Beginning Import", setRouteTagDn)
		_, err = setRtctrlSetTagAttributes(rtctrlSetTag, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetTag - Import finished successfully", setRouteTagDn)
	}

	// rtctrlSetTag - Import finished successfully

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciActionRuleProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ActionRuleProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	rtctrlAttrPAttr := models.ActionRuleProfileAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlAttrPAttr.Annotation = Annotation.(string)
	} else {
		rtctrlAttrPAttr.Annotation = "{}"
	}
	if Name, ok := d.GetOk("name"); ok {
		rtctrlAttrPAttr.Name = Name.(string)
	}
	rtctrlAttrP := models.NewActionRuleProfile(fmt.Sprintf(models.RnrtctrlAttrP, name), TenantDn, desc, nameAlias, rtctrlAttrPAttr)

	err := aciClient.Save(rtctrlAttrP)
	if err != nil {
		return diag.FromErr(err)
	}

	// rtctrlSetTag - Operations
	if setRouteTag, ok := d.GetOk("set_route_tag"); ok {

		log.Printf("[DEBUG] rtctrlSetTag: Beginning Creation")

		rtctrlSetTagAttr := models.RtctrlSetTagAttributes{}
		rtctrlSetTagAttr.Tag = setRouteTag.(string)
		rtctrlSetTag := models.NewRtctrlSetTag(fmt.Sprintf(models.RnrtctrlSetTag), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetTagAttr)

		creation_err := aciClient.Save(rtctrlSetTag)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetTag.DistinguishedName)
		resourceAciRtctrlSetTagRead(ctx, rtctrlSetTag.DistinguishedName, d, m)
	}

	d.SetId(rtctrlAttrP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciActionRuleProfileRead(ctx, d, m)
}

func resourceAciActionRuleProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ActionRuleProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	rtctrlAttrPAttr := models.ActionRuleProfileAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlAttrPAttr.Annotation = Annotation.(string)
	} else {
		rtctrlAttrPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlAttrPAttr.Name = Name.(string)
	}

	rtctrlAttrP := models.NewActionRuleProfile(fmt.Sprintf(models.RnrtctrlAttrP, name), TenantDn, desc, nameAlias, rtctrlAttrPAttr)

	rtctrlAttrP.Status = "modified"

	err := aciClient.Save(rtctrlAttrP)

	if err != nil {
		return diag.FromErr(err)
	}

	// rtctrlSetTag - Operations
	if d.HasChange("set_route_tag") {

		if setRouteTag, ok := d.GetOk("set_route_tag"); ok {

			log.Printf("[DEBUG] rtctrlSetTag - Beginning Creation")

			rtctrlSetTagAttr := models.RtctrlSetTagAttributes{}
			rtctrlSetTagAttr.Tag = setRouteTag.(string)

			setRouteTagDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetTag)

			deletion_err := aciClient.DeleteByDn(setRouteTagDn, "rtctrlSetTag")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetTag := models.NewRtctrlSetTag(fmt.Sprintf(models.RnrtctrlSetTag), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetTagAttr)

			err := aciClient.Save(rtctrlSetTag)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetTag - Creation finished successfully", rtctrlSetTag.DistinguishedName)
			resourceAciRtctrlSetTagRead(ctx, rtctrlSetTag.DistinguishedName, d, m)

		} else {
			setRouteTagDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetTag)
			log.Printf("[DEBUG] %s: rtctrlSetTag - Beginning Destroy", setRouteTagDn)

			err := aciClient.DeleteByDn(setRouteTagDn, "rtctrlSetTag")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetTag - Destroy finished successfully", setRouteTagDn)
		}
	}

	d.SetId(rtctrlAttrP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciActionRuleProfileRead(ctx, d, m)

}

func resourceAciRtctrlSetTagRead(ctx context.Context, dn string, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: rtctrlSetTag - Beginning Read", dn)
	aciClient := m.(*client.Client)

	rtctrlSetTag, err := getRemoteRtctrlSetTag(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = setRtctrlSetTagAttributes(rtctrlSetTag, d)
	if err != nil {
		return nil
	}

	log.Printf("[DEBUG] %s: rtctrlSetTag - Read finished successfully", dn)
	return nil
}

func resourceAciActionRuleProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	rtctrlAttrP, err := getRemoteActionRuleProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setActionRuleProfileAttributes(rtctrlAttrP, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciActionRuleProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "rtctrlAttrP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

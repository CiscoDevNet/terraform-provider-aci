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

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteActionRuleProfile(client *client.Client, dn string) (*models.ActionRuleProfile, error) {
	rtctrlAttrPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlAttrP := models.ActionRuleProfileFromContainer(rtctrlAttrPCont)

	if rtctrlAttrP.DistinguishedName == "" {
		return nil, fmt.Errorf("ActionRuleProfile %s not found", rtctrlAttrP.DistinguishedName)
	}

	return rtctrlAttrP, nil
}

func setActionRuleProfileAttributes(rtctrlAttrP *models.ActionRuleProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(rtctrlAttrP.DistinguishedName)
	d.Set("description", rtctrlAttrP.Description)
	// d.Set("tenant_dn", GetParentDn(rtctrlAttrP.DistinguishedName))
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
	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlAttrPAttr.Annotation = Annotation.(string)
	} else {
		rtctrlAttrPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		rtctrlAttrPAttr.NameAlias = NameAlias.(string)
	}
	rtctrlAttrP := models.NewActionRuleProfile(fmt.Sprintf("attr-%s", name), TenantDn, desc, rtctrlAttrPAttr)

	err := aciClient.Save(rtctrlAttrP)
	if err != nil {
		return diag.FromErr(err)
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
	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlAttrPAttr.Annotation = Annotation.(string)
	} else {
		rtctrlAttrPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		rtctrlAttrPAttr.NameAlias = NameAlias.(string)
	}
	rtctrlAttrP := models.NewActionRuleProfile(fmt.Sprintf("attr-%s", name), TenantDn, desc, rtctrlAttrPAttr)

	rtctrlAttrP.Status = "modified"

	err := aciClient.Save(rtctrlAttrP)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rtctrlAttrP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciActionRuleProfileRead(ctx, d, m)

}

func resourceAciActionRuleProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	rtctrlAttrP, err := getRemoteActionRuleProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
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

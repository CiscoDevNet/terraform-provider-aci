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

func resourceAciEndpointControlPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciEndpointControlPolicyCreate,
		UpdateContext: resourceAciEndpointControlPolicyUpdate,
		ReadContext:   resourceAciEndpointControlPolicyRead,
		DeleteContext: resourceAciEndpointControlPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciEndpointControlPolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},
			"hold_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rogue_ep_detect_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rogue_ep_detect_mult": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteEndpointControlPolicy(client *client.Client, dn string) (*models.EndpointControlPolicy, error) {
	epControlPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	epControlP := models.EndpointControlPolicyFromContainer(epControlPCont)
	if epControlP.DistinguishedName == "" {
		return nil, fmt.Errorf("EndpointControlPolicy %s not found", epControlP.DistinguishedName)
	}
	return epControlP, nil
}

func setEndpointControlPolicyAttributes(epControlP *models.EndpointControlPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(epControlP.DistinguishedName)
	d.Set("description", epControlP.Description)
	epControlPMap, err := epControlP.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("admin_st", epControlPMap["adminSt"])
	d.Set("annotation", epControlPMap["annotation"])
	d.Set("hold_intvl", epControlPMap["holdIntvl"])
	d.Set("rogue_ep_detect_intvl", epControlPMap["rogueEpDetectIntvl"])
	d.Set("rogue_ep_detect_mult", epControlPMap["rogueEpDetectMult"])
	d.Set("name_alias", epControlPMap["nameAlias"])
	return d, nil
}

func resourceAciEndpointControlPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	epControlP, err := getRemoteEndpointControlPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setEndpointControlPolicyAttributes(epControlP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciEndpointControlPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EndpointControlPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	epControlPAttr := models.EndpointControlPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		epControlPAttr.Annotation = Annotation.(string)
	} else {
		epControlPAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_st"); ok {
		epControlPAttr.AdminSt = AdminSt.(string)
	}

	if HoldIntvl, ok := d.GetOk("hold_intvl"); ok {
		epControlPAttr.HoldIntvl = HoldIntvl.(string)
	}

	epControlPAttr.Name = "default"

	if RogueEpDetectIntvl, ok := d.GetOk("rogue_ep_detect_intvl"); ok {
		epControlPAttr.RogueEpDetectIntvl = RogueEpDetectIntvl.(string)
	}

	if RogueEpDetectMult, ok := d.GetOk("rogue_ep_detect_mult"); ok {
		epControlPAttr.RogueEpDetectMult = RogueEpDetectMult.(string)
	}
	epControlP := models.NewEndpointControlPolicy(fmt.Sprintf("infra/epCtrlP-%s", name), "uni", desc, nameAlias, epControlPAttr)
	epControlP.Status = "modified"
	err := aciClient.Save(epControlP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(epControlP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciEndpointControlPolicyRead(ctx, d, m)
}

func resourceAciEndpointControlPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EndpointControlPolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	epControlPAttr := models.EndpointControlPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		epControlPAttr.Annotation = Annotation.(string)
	} else {
		epControlPAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_st"); ok {
		epControlPAttr.AdminSt = AdminSt.(string)
	}

	if HoldIntvl, ok := d.GetOk("hold_intvl"); ok {
		epControlPAttr.HoldIntvl = HoldIntvl.(string)
	}

	epControlPAttr.Name = "default"

	if RogueEpDetectIntvl, ok := d.GetOk("rogue_ep_detect_intvl"); ok {
		epControlPAttr.RogueEpDetectIntvl = RogueEpDetectIntvl.(string)
	}

	if RogueEpDetectMult, ok := d.GetOk("rogue_ep_detect_mult"); ok {
		epControlPAttr.RogueEpDetectMult = RogueEpDetectMult.(string)
	}
	epControlP := models.NewEndpointControlPolicy(fmt.Sprintf("infra/epCtrlP-%s", name), "uni", desc, nameAlias, epControlPAttr)
	epControlP.Status = "modified"
	err := aciClient.Save(epControlP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(epControlP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciEndpointControlPolicyRead(ctx, d, m)
}

func resourceAciEndpointControlPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	epControlP, err := getRemoteEndpointControlPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setEndpointControlPolicyAttributes(epControlP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciEndpointControlPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name epControlP cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}

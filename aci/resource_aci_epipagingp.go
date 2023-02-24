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

func resourceAciIPAgingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciIPAgingPolicyCreate,
		UpdateContext: resourceAciIPAgingPolicyUpdate,
		ReadContext:   resourceAciIPAgingPolicyRead,
		DeleteContext: resourceAciIPAgingPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciIPAgingPolicyImport,
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
		})),
	}
}

func getRemoteIPAgingPolicy(client *client.Client, dn string) (*models.IPAgingPolicy, error) {
	epIpAgingPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	epIpAgingP := models.IPAgingPolicyFromContainer(epIpAgingPCont)
	if epIpAgingP.DistinguishedName == "" {
		return nil, fmt.Errorf("IP Aging Policy %s not found", dn)
	}
	return epIpAgingP, nil
}

func setIPAgingPolicyAttributes(epIpAgingP *models.IPAgingPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(epIpAgingP.DistinguishedName)
	d.Set("description", epIpAgingP.Description)
	epIpAgingPMap, err := epIpAgingP.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("admin_st", epIpAgingPMap["adminSt"])
	d.Set("annotation", epIpAgingPMap["annotation"])
	d.Set("name_alias", epIpAgingPMap["nameAlias"])
	return d, nil
}

func resourceAciIPAgingPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	epIpAgingP, err := getRemoteIPAgingPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setIPAgingPolicyAttributes(epIpAgingP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciIPAgingPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] IPAgingPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	epIpAgingPAttr := models.IPAgingPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		epIpAgingPAttr.Annotation = Annotation.(string)
	} else {
		epIpAgingPAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_st"); ok {
		epIpAgingPAttr.AdminSt = AdminSt.(string)
	}

	epIpAgingPAttr.Name = "default"

	epIpAgingP := models.NewIPAgingPolicy(fmt.Sprintf("infra/ipAgingP-%s", name), "uni", desc, nameAlias, epIpAgingPAttr)
	epIpAgingP.Status = "modified"
	err := aciClient.Save(epIpAgingP)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(epIpAgingP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciIPAgingPolicyRead(ctx, d, m)
}

func resourceAciIPAgingPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] IPAgingPolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	epIpAgingPAttr := models.IPAgingPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		epIpAgingPAttr.Annotation = Annotation.(string)
	} else {
		epIpAgingPAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_st"); ok {
		epIpAgingPAttr.AdminSt = AdminSt.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		epIpAgingPAttr.Name = Name.(string)
	}
	epIpAgingP := models.NewIPAgingPolicy(fmt.Sprintf("infra/ipAgingP-%s", name), "uni", desc, nameAlias, epIpAgingPAttr)
	epIpAgingP.Status = "modified"
	err := aciClient.Save(epIpAgingP)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(epIpAgingP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciIPAgingPolicyRead(ctx, d, m)
}

func resourceAciIPAgingPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	epIpAgingP, err := getRemoteIPAgingPolicy(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setIPAgingPolicyAttributes(epIpAgingP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciIPAgingPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name epIpAgingP cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}

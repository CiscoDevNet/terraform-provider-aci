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

func resourceAciCOOPGroupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCOOPGroupPolicyCreate,
		UpdateContext: resourceAciCOOPGroupPolicyUpdate,
		ReadContext:   resourceAciCOOPGroupPolicyRead,
		DeleteContext: resourceAciCOOPGroupPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCOOPGroupPolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"compatible",
					"strict",
				}, false),
			},
		})),
	}
}

func getRemoteCOOPGroupPolicy(client *client.Client, dn string) (*models.COOPGroupPolicy, error) {
	coopPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	coopPol := models.COOPGroupPolicyFromContainer(coopPolCont)
	if coopPol.DistinguishedName == "" {
		return nil, fmt.Errorf("COOP Group Policy %s not found", dn)
	}
	return coopPol, nil
}

func setCOOPGroupPolicyAttributes(coopPol *models.COOPGroupPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(coopPol.DistinguishedName)
	d.Set("description", coopPol.Description)
	coopPolMap, err := coopPol.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", coopPolMap["annotation"])
	d.Set("type", coopPolMap["type"])
	d.Set("name_alias", coopPolMap["nameAlias"])
	return d, nil
}

func resourceAciCOOPGroupPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	coopPol, err := getRemoteCOOPGroupPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCOOPGroupPolicyAttributes(coopPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCOOPGroupPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] COOPGroupPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	coopPolAttr := models.COOPGroupPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		coopPolAttr.Annotation = Annotation.(string)
	} else {
		coopPolAttr.Annotation = "{}"
	}

	coopPolAttr.Name = "default"

	if Type, ok := d.GetOk("type"); ok {
		coopPolAttr.Type = Type.(string)
	}
	coopPol := models.NewCOOPGroupPolicy(fmt.Sprintf("fabric/pol-%s", name), "uni", desc, nameAlias, coopPolAttr)
	coopPol.Status = "modified"
	err := aciClient.Save(coopPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(coopPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCOOPGroupPolicyRead(ctx, d, m)
}

func resourceAciCOOPGroupPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] COOPGroupPolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	coopPolAttr := models.COOPGroupPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		coopPolAttr.Annotation = Annotation.(string)
	} else {
		coopPolAttr.Annotation = "{}"
	}

	coopPolAttr.Name = "default"

	if Type, ok := d.GetOk("type"); ok {
		coopPolAttr.Type = Type.(string)
	}
	coopPol := models.NewCOOPGroupPolicy(fmt.Sprintf("fabric/pol-%s", name), "uni", desc, nameAlias, coopPolAttr)
	coopPol.Status = "modified"
	err := aciClient.Save(coopPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(coopPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCOOPGroupPolicyRead(ctx, d, m)
}

func resourceAciCOOPGroupPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	coopPol, err := getRemoteCOOPGroupPolicy(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setCOOPGroupPolicyAttributes(coopPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCOOPGroupPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name coopPol cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}

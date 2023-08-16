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

func resourceAciPowerSupplyRedundancyPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciPowerSupplyRedundancyPolicyCreate,
		UpdateContext: resourceAciPowerSupplyRedundancyPolicyUpdate,
		ReadContext:   resourceAciPowerSupplyRedundancyPolicyRead,
		DeleteContext: resourceAciPowerSupplyRedundancyPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPowerSupplyRedundancyPolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"administrative_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"comb",
					"insrc-rdn",
					"n-rdn",
					"not-supp",
					"ps-rdn",
					"rdn",
					"sinin-rdn",
					"unknown",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemotePowerSupplyRedundancyPolicy(client *client.Client, dn string) (*models.PsuInstPol, error) {
	psuInstPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	psuInstPol := models.PsuInstPolFromContainer(psuInstPolCont)
	if psuInstPol.DistinguishedName == "" {
		return nil, fmt.Errorf("Power Supply Redundancy Policy %s not found", dn)
	}
	return psuInstPol, nil
}

func setPowerSupplyRedundancyPolicyAttributes(psuInstPol *models.PsuInstPol, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(psuInstPol.DistinguishedName)
	d.Set("description", psuInstPol.Description)
	psuInstPolMap, err := psuInstPol.ToMap()
	if err != nil {
		return nil, err
	}

	d.Set("administrative_state", psuInstPolMap["adminRdnM"])
	d.Set("annotation", psuInstPolMap["annotation"])
	d.Set("name", psuInstPolMap["name"])
	d.Set("name_alias", psuInstPolMap["nameAlias"])
	return d, nil
}

func resourceAciPowerSupplyRedundancyPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	psuInstPol, err := getRemotePowerSupplyRedundancyPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPowerSupplyRedundancyPolicyAttributes(psuInstPol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPowerSupplyRedundancyPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Power Supply Redundancy Policy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	psuInstPolAttr := models.PsuInstPolAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		psuInstPolAttr.Annotation = Annotation.(string)
	} else {
		psuInstPolAttr.Annotation = "{}"
	}

	if AdminRdnM, ok := d.GetOk("administrative_state"); ok {
		psuInstPolAttr.AdminRdnM = AdminRdnM.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		psuInstPolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		psuInstPolAttr.NameAlias = NameAlias.(string)
	}
	psuInstPol := models.NewPowerSupplyRedundancyPolicy(fmt.Sprintf(models.RnPsuInstPol, name), desc, psuInstPolAttr)
	err := aciClient.Save(psuInstPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(psuInstPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciPowerSupplyRedundancyPolicyRead(ctx, d, m)
}

func resourceAciPowerSupplyRedundancyPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Power Supply Redundancy Policy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)

	psuInstPolAttr := models.PsuInstPolAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		psuInstPolAttr.Annotation = Annotation.(string)
	} else {
		psuInstPolAttr.Annotation = "{}"
	}

	if AdminRdnM, ok := d.GetOk("administrative_state"); ok {
		psuInstPolAttr.AdminRdnM = AdminRdnM.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		psuInstPolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		psuInstPolAttr.NameAlias = NameAlias.(string)
	}
	psuInstPol := models.NewPowerSupplyRedundancyPolicy(fmt.Sprintf(models.RnPsuInstPol, name), desc, psuInstPolAttr)
	psuInstPol.Status = "modified"

	err := aciClient.Save(psuInstPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(psuInstPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciPowerSupplyRedundancyPolicyRead(ctx, d, m)
}

func resourceAciPowerSupplyRedundancyPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	psuInstPol, err := getRemotePowerSupplyRedundancyPolicy(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setPowerSupplyRedundancyPolicyAttributes(psuInstPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciPowerSupplyRedundancyPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, models.PsuInstPolClassName)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

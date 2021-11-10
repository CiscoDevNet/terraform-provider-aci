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

func resourceAciEPGsUsingFunction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciEPGsUsingFunctionCreate,
		UpdateContext: resourceAciEPGsUsingFunctionUpdate,
		ReadContext:   resourceAciEPGsUsingFunctionRead,
		DeleteContext: resourceAciEPGsUsingFunctionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciEPGsUsingFunctionImport,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"access_generic_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Default:  "orchestrator:terraform",
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"instr_imedcy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"immediate",
					"lazy",
				}, false),
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"regular",
					"native",
					"untagged",
				}, false),
			},

			"primary_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}
func getRemoteEPGsUsingFunction(client *client.Client, dn string) (*models.EPGsUsingFunction, error) {
	infraRsFuncToEpgCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraRsFuncToEpg := models.EPGsUsingFunctionFromContainer(infraRsFuncToEpgCont)

	if infraRsFuncToEpg.DistinguishedName == "" {
		return nil, fmt.Errorf("EPGsUsingFunction %s not found", infraRsFuncToEpg.DistinguishedName)
	}

	return infraRsFuncToEpg, nil
}

func setEPGsUsingFunctionAttributes(infraRsFuncToEpg *models.EPGsUsingFunction, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(infraRsFuncToEpg.DistinguishedName)
	if dn != infraRsFuncToEpg.DistinguishedName {
		d.Set("access_generic_dn", "")
	}
	infraRsFuncToEpgMap, err := infraRsFuncToEpg.ToMap()

	if err != nil {
		return d, err
	}
	d.Set("access_generic_dn", GetParentDn(infraRsFuncToEpg.DistinguishedName, fmt.Sprintf("/rsfuncToEpg-[%s]", infraRsFuncToEpgMap["tDn"])))
	d.Set("tdn", infraRsFuncToEpgMap["tDn"])
	d.Set("annotation", infraRsFuncToEpgMap["annotation"])
	d.Set("encap", infraRsFuncToEpgMap["encap"])
	d.Set("instr_imedcy", infraRsFuncToEpgMap["instrImedcy"])
	d.Set("mode", infraRsFuncToEpgMap["mode"])
	d.Set("primary_encap", infraRsFuncToEpgMap["primaryEncap"])
	return d, nil
}

func resourceAciEPGsUsingFunctionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraRsFuncToEpg, err := getRemoteEPGsUsingFunction(aciClient, dn)

	if err != nil {
		return nil, err
	}
	infraRsFuncToEpgMap, err := infraRsFuncToEpg.ToMap()

	if err != nil {
		return nil, err
	}

	tDn := infraRsFuncToEpgMap["tDn"]
	pDN := GetParentDn(dn, fmt.Sprintf("/rsfuncToEpg-[%s]", tDn))
	d.Set("access_generic_dn", pDN)
	schemaFilled, err := setEPGsUsingFunctionAttributes(infraRsFuncToEpg, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciEPGsUsingFunctionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EPGsUsingFunction: Beginning Creation")
	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	AccessGenericDn := d.Get("access_generic_dn").(string)

	infraRsFuncToEpgAttr := models.EPGsUsingFunctionAttributes{}
	infraRsFuncToEpgAttr.TDn = tDn
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraRsFuncToEpgAttr.Annotation = Annotation.(string)
	} else {
		infraRsFuncToEpgAttr.Annotation = "{}"
	}
	if Encap, ok := d.GetOk("encap"); ok {
		infraRsFuncToEpgAttr.Encap = Encap.(string)
	}
	if InstrImedcy, ok := d.GetOk("instr_imedcy"); ok {
		infraRsFuncToEpgAttr.InstrImedcy = InstrImedcy.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		infraRsFuncToEpgAttr.Mode = Mode.(string)
	}
	if PrimaryEncap, ok := d.GetOk("primary_encap"); ok {
		infraRsFuncToEpgAttr.PrimaryEncap = PrimaryEncap.(string)
	}
	infraRsFuncToEpg := models.NewEPGsUsingFunction(fmt.Sprintf("rsfuncToEpg-[%s]", tDn), AccessGenericDn, "", infraRsFuncToEpgAttr)

	err := aciClient.Save(infraRsFuncToEpg)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraRsFuncToEpg.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciEPGsUsingFunctionRead(ctx, d, m)
}

func resourceAciEPGsUsingFunctionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EPGsUsingFunction: Beginning Update")

	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	AccessGenericDn := d.Get("access_generic_dn").(string)

	infraRsFuncToEpgAttr := models.EPGsUsingFunctionAttributes{}
	infraRsFuncToEpgAttr.TDn = tDn
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraRsFuncToEpgAttr.Annotation = Annotation.(string)
	} else {
		infraRsFuncToEpgAttr.Annotation = "{}"
	}
	if Encap, ok := d.GetOk("encap"); ok {
		infraRsFuncToEpgAttr.Encap = Encap.(string)
	}
	if InstrImedcy, ok := d.GetOk("instr_imedcy"); ok {
		infraRsFuncToEpgAttr.InstrImedcy = InstrImedcy.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		infraRsFuncToEpgAttr.Mode = Mode.(string)
	}
	if PrimaryEncap, ok := d.GetOk("primary_encap"); ok {
		infraRsFuncToEpgAttr.PrimaryEncap = PrimaryEncap.(string)
	}
	infraRsFuncToEpg := models.NewEPGsUsingFunction(fmt.Sprintf("rsfuncToEpg-[%s]", tDn), AccessGenericDn, "", infraRsFuncToEpgAttr)

	infraRsFuncToEpg.Status = "modified"

	err := aciClient.Save(infraRsFuncToEpg)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraRsFuncToEpg.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciEPGsUsingFunctionRead(ctx, d, m)

}

func resourceAciEPGsUsingFunctionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraRsFuncToEpg, err := getRemoteEPGsUsingFunction(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setEPGsUsingFunctionAttributes(infraRsFuncToEpg, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciEPGsUsingFunctionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraRsFuncToEpg")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

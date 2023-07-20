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

func resourceAciBfdMultihopInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBfdMultihopInterfacePolicyCreate,
		UpdateContext: resourceAciBfdMultihopInterfacePolicyUpdate,
		ReadContext:   resourceAciBfdMultihopInterfacePolicyRead,
		DeleteContext: resourceAciBfdMultihopInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBfdMultihopInterfacePolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"admin_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},
			"detection_multiplier": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"min_receive_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"min_transmit_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getAciBfdMultihopInterfacePolicy(client *client.Client, dn string) (*models.AciBfdMultihopInterfacePolicy, error) {
	bfdMhIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	bfdMhIfPol := models.AciBfdMultihopInterfacePolicyFromContainer(bfdMhIfPolCont)
	if bfdMhIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("Aci BFD Multihop Interface Policy %s not found", dn)
	}
	return bfdMhIfPol, nil
}

func setAciBfdMultihopInterfacePolicyAttributes(bfdMhIfPol *models.AciBfdMultihopInterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(bfdMhIfPol.DistinguishedName)
	d.Set("description", bfdMhIfPol.Description)
	if dn != bfdMhIfPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	bfdMhIfPolMap, err := bfdMhIfPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("tenant_dn", GetParentDn(dn, "/"+fmt.Sprintf(models.RnbfdMhIfPol, bfdMhIfPolMap["name"])))
	d.Set("admin_state", bfdMhIfPolMap["adminSt"])
	d.Set("annotation", bfdMhIfPolMap["annotation"])
	d.Set("detection_multiplier", bfdMhIfPolMap["detectMult"])
	d.Set("min_receive_interval", bfdMhIfPolMap["minRxIntvl"])
	d.Set("min_transmit_interval", bfdMhIfPolMap["minTxIntvl"])
	d.Set("name", bfdMhIfPolMap["name"])
	d.Set("name_alias", bfdMhIfPolMap["nameAlias"])
	return d, nil
}

func resourceAciBfdMultihopInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	bfdMhIfPol, err := getAciBfdMultihopInterfacePolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setAciBfdMultihopInterfacePolicyAttributes(bfdMhIfPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBfdMultihopInterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Aci BFD Multihop Interface Policy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	nameAlias := d.Get("name_alias").(string)
	name := d.Get("name").(string)
	tenantDn := d.Get("tenant_dn").(string)

	bfdMhIfPolAttr := models.AciBfdMultihopInterfacePolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		bfdMhIfPolAttr.Annotation = Annotation.(string)
	} else {
		bfdMhIfPolAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_state"); ok {
		bfdMhIfPolAttr.AdminSt = AdminSt.(string)
	}

	if DetectMult, ok := d.GetOk("detection_multiplier"); ok {
		bfdMhIfPolAttr.DetectMult = DetectMult.(string)
	}

	if MinRxIntvl, ok := d.GetOk("min_receive_interval"); ok {
		bfdMhIfPolAttr.MinRxIntvl = MinRxIntvl.(string)
	}

	if MinTxIntvl, ok := d.GetOk("min_transmit_interval"); ok {
		bfdMhIfPolAttr.MinTxIntvl = MinTxIntvl.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		bfdMhIfPolAttr.Name = Name.(string)
	}
	bfdMhIfPol := models.NewAciBfdMultihopInterfacePolicy(fmt.Sprintf(models.RnbfdMhIfPol, name), tenantDn, desc, nameAlias, bfdMhIfPolAttr)

	err := aciClient.Save(bfdMhIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bfdMhIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciBfdMultihopInterfacePolicyRead(ctx, d, m)
}

func resourceAciBfdMultihopInterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Aci BFD Multihop Interface Policy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	nameAlias := d.Get("name_alias").(string)
	name := d.Get("name").(string)
	tenantDn := d.Get("tenant_dn").(string)

	bfdMhIfPolAttr := models.AciBfdMultihopInterfacePolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		bfdMhIfPolAttr.Annotation = Annotation.(string)
	} else {
		bfdMhIfPolAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_state"); ok {
		bfdMhIfPolAttr.AdminSt = AdminSt.(string)
	}

	if DetectMult, ok := d.GetOk("detection_multiplier"); ok {
		bfdMhIfPolAttr.DetectMult = DetectMult.(string)
	}

	if MinRxIntvl, ok := d.GetOk("min_receive_interval"); ok {
		bfdMhIfPolAttr.MinRxIntvl = MinRxIntvl.(string)
	}

	if MinTxIntvl, ok := d.GetOk("min_transmit_interval"); ok {
		bfdMhIfPolAttr.MinTxIntvl = MinTxIntvl.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		bfdMhIfPolAttr.Name = Name.(string)
	}
	bfdMhIfPol := models.NewAciBfdMultihopInterfacePolicy(fmt.Sprintf(models.RnbfdMhIfPol, name), tenantDn, desc, nameAlias, bfdMhIfPolAttr)

	bfdMhIfPol.Status = "modified"

	err := aciClient.Save(bfdMhIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bfdMhIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciBfdMultihopInterfacePolicyRead(ctx, d, m)
}

func resourceAciBfdMultihopInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	bfdMhIfPol, err := getAciBfdMultihopInterfacePolicy(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setAciBfdMultihopInterfacePolicyAttributes(bfdMhIfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciBfdMultihopInterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, models.BfdmhifpolClassName)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

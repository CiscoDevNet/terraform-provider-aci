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

func resourceAciBFDInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBFDInterfacePolicyCreate,
		UpdateContext: resourceAciBFDInterfacePolicyUpdate,
		ReadContext:   resourceAciBFDInterfacePolicyRead,
		DeleteContext: resourceAciBFDInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBFDInterfacePolicyImport,
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

			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"opt-subif", "none",
				}, false),
			},

			"detect_mult": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"echo_admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"echo_rx_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"min_rx_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"min_tx_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteBFDInterfacePolicy(client *client.Client, dn string) (*models.BFDInterfacePolicy, error) {
	bfdIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bfdIfPol := models.BFDInterfacePolicyFromContainer(bfdIfPolCont)

	if bfdIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("BFDInterfacePolicy %s not found", bfdIfPol.DistinguishedName)
	}

	return bfdIfPol, nil
}

func setBFDInterfacePolicyAttributes(bfdIfPol *models.BFDInterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(bfdIfPol.DistinguishedName)
	d.Set("description", bfdIfPol.Description)
	if dn != bfdIfPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	bfdIfPolMap, err := bfdIfPol.ToMap()

	if err != nil {
		return nil, err
	}

	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/bfdIfPol-%s", bfdIfPolMap["name"])))

	d.Set("name", bfdIfPolMap["name"])

	d.Set("admin_st", bfdIfPolMap["adminSt"])
	d.Set("annotation", bfdIfPolMap["annotation"])
	d.Set("ctrl", bfdIfPolMap["ctrl"])
	d.Set("detect_mult", bfdIfPolMap["detectMult"])
	d.Set("echo_admin_st", bfdIfPolMap["echoAdminSt"])
	d.Set("echo_rx_intvl", bfdIfPolMap["echoRxIntvl"])
	d.Set("min_rx_intvl", bfdIfPolMap["minRxIntvl"])
	d.Set("min_tx_intvl", bfdIfPolMap["minTxIntvl"])
	d.Set("name_alias", bfdIfPolMap["nameAlias"])
	return d, nil
}

func resourceAciBFDInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bfdIfPol, err := getRemoteBFDInterfacePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setBFDInterfacePolicyAttributes(bfdIfPol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBFDInterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BFDInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bfdIfPolAttr := models.BFDInterfacePolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		bfdIfPolAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bfdIfPolAttr.Annotation = Annotation.(string)
	} else {
		bfdIfPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		if Ctrl == "none" {
			bfdIfPolAttr.Ctrl = "0"
		} else {
			bfdIfPolAttr.Ctrl = Ctrl.(string)
		}
	}
	if DetectMult, ok := d.GetOk("detect_mult"); ok {
		bfdIfPolAttr.DetectMult = DetectMult.(string)
	}
	if EchoAdminSt, ok := d.GetOk("echo_admin_st"); ok {
		bfdIfPolAttr.EchoAdminSt = EchoAdminSt.(string)
	}
	if EchoRxIntvl, ok := d.GetOk("echo_rx_intvl"); ok {
		bfdIfPolAttr.EchoRxIntvl = EchoRxIntvl.(string)
	}
	if MinRxIntvl, ok := d.GetOk("min_rx_intvl"); ok {
		bfdIfPolAttr.MinRxIntvl = MinRxIntvl.(string)
	}
	if MinTxIntvl, ok := d.GetOk("min_tx_intvl"); ok {
		bfdIfPolAttr.MinTxIntvl = MinTxIntvl.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bfdIfPolAttr.NameAlias = NameAlias.(string)
	}
	bfdIfPol := models.NewBFDInterfacePolicy(fmt.Sprintf("bfdIfPol-%s", name), TenantDn, desc, bfdIfPolAttr)

	err := aciClient.Save(bfdIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bfdIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBFDInterfacePolicyRead(ctx, d, m)
}

func resourceAciBFDInterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BFDInterfacePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bfdIfPolAttr := models.BFDInterfacePolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		bfdIfPolAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bfdIfPolAttr.Annotation = Annotation.(string)
	} else {
		bfdIfPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		if Ctrl == "none" {
			bfdIfPolAttr.Ctrl = "0"
		} else {
			bfdIfPolAttr.Ctrl = Ctrl.(string)
		}
	}
	if DetectMult, ok := d.GetOk("detect_mult"); ok {
		bfdIfPolAttr.DetectMult = DetectMult.(string)
	}
	if EchoAdminSt, ok := d.GetOk("echo_admin_st"); ok {
		bfdIfPolAttr.EchoAdminSt = EchoAdminSt.(string)
	}
	if EchoRxIntvl, ok := d.GetOk("echo_rx_intvl"); ok {
		bfdIfPolAttr.EchoRxIntvl = EchoRxIntvl.(string)
	}
	if MinRxIntvl, ok := d.GetOk("min_rx_intvl"); ok {
		bfdIfPolAttr.MinRxIntvl = MinRxIntvl.(string)
	}
	if MinTxIntvl, ok := d.GetOk("min_tx_intvl"); ok {
		bfdIfPolAttr.MinTxIntvl = MinTxIntvl.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bfdIfPolAttr.NameAlias = NameAlias.(string)
	}
	bfdIfPol := models.NewBFDInterfacePolicy(fmt.Sprintf("bfdIfPol-%s", name), TenantDn, desc, bfdIfPolAttr)

	bfdIfPol.Status = "modified"

	err := aciClient.Save(bfdIfPol)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bfdIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBFDInterfacePolicyRead(ctx, d, m)

}

func resourceAciBFDInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bfdIfPol, err := getRemoteBFDInterfacePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setBFDInterfacePolicyAttributes(bfdIfPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBFDInterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bfdIfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

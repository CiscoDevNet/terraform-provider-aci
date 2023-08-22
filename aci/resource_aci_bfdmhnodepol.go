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

func resourceAciBFDMultihopNodePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBFDMultihopNodePolicyCreate,
		UpdateContext: resourceAciBFDMultihopNodePolicyUpdate,
		ReadContext:   resourceAciBFDMultihopNodePolicyRead,
		DeleteContext: resourceAciBFDMultihopNodePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBFDMultihopNodePolicyImport,
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
			"min_rx_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"min_tx_interval": {
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

func getRemoteBFDMultihopNodePolicy(client *client.Client, dn string) (*models.BFDMultihopNodePolicy, error) {
	bfdMhNodePolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	bfdMhNodePol := models.BFDMultihopNodePolicyFromContainer(bfdMhNodePolCont)
	if bfdMhNodePol.DistinguishedName == "" {
		return nil, fmt.Errorf("BFD Multihop Node Policy %s not found", dn)
	}
	return bfdMhNodePol, nil
}

func setBFDMultihopNodePolicyAttributes(bfdMhNodePol *models.BFDMultihopNodePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(bfdMhNodePol.DistinguishedName)
	d.Set("description", bfdMhNodePol.Description)
	bfdMhNodePolMap, err := bfdMhNodePol.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != bfdMhNodePol.DistinguishedName {
		d.Set("tenant_dn", "")
	} else {
		d.Set("tenant_dn", GetParentDn(bfdMhNodePol.DistinguishedName, fmt.Sprintf("/"+models.RnBfdMhNodePol, bfdMhNodePolMap["name"])))
	}
	d.Set("admin_state", bfdMhNodePolMap["adminSt"])
	d.Set("annotation", bfdMhNodePolMap["annotation"])
	d.Set("detection_multiplier", bfdMhNodePolMap["detectMult"])
	d.Set("min_rx_interval", bfdMhNodePolMap["minRxIntvl"])
	d.Set("min_tx_interval", bfdMhNodePolMap["minTxIntvl"])
	d.Set("name", bfdMhNodePolMap["name"])
	d.Set("name_alias", bfdMhNodePolMap["nameAlias"])
	return d, nil
}

func resourceAciBFDMultihopNodePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	bfdMhNodePol, err := getRemoteBFDMultihopNodePolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setBFDMultihopNodePolicyAttributes(bfdMhNodePol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBFDMultihopNodePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BFD Multihop Node Policy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	tenantDn := d.Get("tenant_dn").(string)

	bfdMhNodePolAttr := models.BFDMultihopNodePolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		bfdMhNodePolAttr.Annotation = Annotation.(string)
	} else {
		bfdMhNodePolAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_state"); ok {
		bfdMhNodePolAttr.AdminSt = AdminSt.(string)
	}

	if DetectMult, ok := d.GetOk("detection_multiplier"); ok {
		bfdMhNodePolAttr.DetectMult = DetectMult.(string)
	}

	if MinRxIntvl, ok := d.GetOk("min_rx_interval"); ok {
		bfdMhNodePolAttr.MinRxIntvl = MinRxIntvl.(string)
	}

	if MinTxIntvl, ok := d.GetOk("min_tx_interval"); ok {
		bfdMhNodePolAttr.MinTxIntvl = MinTxIntvl.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		bfdMhNodePolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bfdMhNodePolAttr.NameAlias = NameAlias.(string)
	}
	bfdMhNodePol := models.NewBFDMultihopNodePolicy(fmt.Sprintf(models.RnBfdMhNodePol, name), tenantDn, desc, bfdMhNodePolAttr)

	err := aciClient.Save(bfdMhNodePol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bfdMhNodePol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciBFDMultihopNodePolicyRead(ctx, d, m)
}
func resourceAciBFDMultihopNodePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BFD Multihop Node Policy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	tenantDn := d.Get("tenant_dn").(string)

	bfdMhNodePolAttr := models.BFDMultihopNodePolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		bfdMhNodePolAttr.Annotation = Annotation.(string)
	} else {
		bfdMhNodePolAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_state"); ok {
		bfdMhNodePolAttr.AdminSt = AdminSt.(string)
	}

	if DetectMult, ok := d.GetOk("detection_multiplier"); ok {
		bfdMhNodePolAttr.DetectMult = DetectMult.(string)
	}

	if MinRxIntvl, ok := d.GetOk("min_rx_interval"); ok {
		bfdMhNodePolAttr.MinRxIntvl = MinRxIntvl.(string)
	}

	if MinTxIntvl, ok := d.GetOk("min_tx_interval"); ok {
		bfdMhNodePolAttr.MinTxIntvl = MinTxIntvl.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		bfdMhNodePolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bfdMhNodePolAttr.NameAlias = NameAlias.(string)
	}
	bfdMhNodePol := models.NewBFDMultihopNodePolicy(fmt.Sprintf(models.RnBfdMhNodePol, name), tenantDn, desc, bfdMhNodePolAttr)

	bfdMhNodePol.Status = "modified"

	err := aciClient.Save(bfdMhNodePol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bfdMhNodePol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciBFDMultihopNodePolicyRead(ctx, d, m)
}

func resourceAciBFDMultihopNodePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	bfdMhNodePol, err := getRemoteBFDMultihopNodePolicy(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setBFDMultihopNodePolicyAttributes(bfdMhNodePol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciBFDMultihopNodePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, models.BfdMhNodePolClassName)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

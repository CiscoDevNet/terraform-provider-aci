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

func resourceAciBGPAddressFamilyContextPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBGPAddressFamilyContextPolicyCreate,
		UpdateContext: resourceAciBGPAddressFamilyContextPolicyUpdate,
		ReadContext:   resourceAciBGPAddressFamilyContextPolicyRead,
		DeleteContext: resourceAciBGPAddressFamilyContextPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBGPAddressFamilyContextPolicyImport,
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

			"ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"host-rt-leak",
				}, false),
			},

			"e_dist": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"i_dist": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"local_dist": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_ecmp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_ecmp_ibgp": &schema.Schema{
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

func getRemoteBGPAddressFamilyContextPolicy(client *client.Client, dn string) (*models.BGPAddressFamilyContextPolicy, error) {
	bgpCtxAfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpCtxAfPol := models.BGPAddressFamilyContextPolicyFromContainer(bgpCtxAfPolCont)

	if bgpCtxAfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("BGPAddressFamilyContextPolicy %s not found", bgpCtxAfPol.DistinguishedName)
	}

	return bgpCtxAfPol, nil
}

func setBGPAddressFamilyContextPolicyAttributes(bgpCtxAfPol *models.BGPAddressFamilyContextPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	d.SetId(bgpCtxAfPol.DistinguishedName)
	d.Set("description", bgpCtxAfPol.Description)
	if dn != bgpCtxAfPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	bgpCtxAfPolMap, err := bgpCtxAfPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/bgpCtxAfP-%s", bgpCtxAfPolMap["name"])))
	d.Set("name", bgpCtxAfPolMap["name"])

	d.Set("annotation", bgpCtxAfPolMap["annotation"])
	d.Set("ctrl", bgpCtxAfPolMap["ctrl"])
	d.Set("e_dist", bgpCtxAfPolMap["eDist"])
	d.Set("i_dist", bgpCtxAfPolMap["iDist"])
	d.Set("local_dist", bgpCtxAfPolMap["localDist"])
	d.Set("max_ecmp", bgpCtxAfPolMap["maxEcmp"])
	d.Set("max_ecmp_ibgp", bgpCtxAfPolMap["maxEcmpIbgp"])
	d.Set("name_alias", bgpCtxAfPolMap["nameAlias"])

	return d, nil
}

func resourceAciBGPAddressFamilyContextPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bgpCtxAfPol, err := getRemoteBGPAddressFamilyContextPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setBGPAddressFamilyContextPolicyAttributes(bgpCtxAfPol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBGPAddressFamilyContextPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BGPAddressFamilyContextPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpCtxAfPolAttr := models.BGPAddressFamilyContextPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpCtxAfPolAttr.Annotation = Annotation.(string)
	} else {
		bgpCtxAfPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		bgpCtxAfPolAttr.Ctrl = Ctrl.(string)
	}
	if EDist, ok := d.GetOk("e_dist"); ok {
		bgpCtxAfPolAttr.EDist = EDist.(string)
	}
	if IDist, ok := d.GetOk("i_dist"); ok {
		bgpCtxAfPolAttr.IDist = IDist.(string)
	}
	if LocalDist, ok := d.GetOk("local_dist"); ok {
		bgpCtxAfPolAttr.LocalDist = LocalDist.(string)
	}
	if MaxEcmp, ok := d.GetOk("max_ecmp"); ok {
		bgpCtxAfPolAttr.MaxEcmp = MaxEcmp.(string)
	}
	if MaxEcmpIbgp, ok := d.GetOk("max_ecmp_ibgp"); ok {
		bgpCtxAfPolAttr.MaxEcmpIbgp = MaxEcmpIbgp.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpCtxAfPolAttr.NameAlias = NameAlias.(string)
	}
	bgpCtxAfPol := models.NewBGPAddressFamilyContextPolicy(fmt.Sprintf("bgpCtxAfP-%s", name), TenantDn, desc, bgpCtxAfPolAttr)

	err := aciClient.Save(bgpCtxAfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bgpCtxAfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBGPAddressFamilyContextPolicyRead(ctx, d, m)
}

func resourceAciBGPAddressFamilyContextPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BGPAddressFamilyContextPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpCtxAfPolAttr := models.BGPAddressFamilyContextPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpCtxAfPolAttr.Annotation = Annotation.(string)
	} else {
		bgpCtxAfPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		bgpCtxAfPolAttr.Ctrl = Ctrl.(string)
	}
	if EDist, ok := d.GetOk("e_dist"); ok {
		bgpCtxAfPolAttr.EDist = EDist.(string)
	}
	if IDist, ok := d.GetOk("i_dist"); ok {
		bgpCtxAfPolAttr.IDist = IDist.(string)
	}
	if LocalDist, ok := d.GetOk("local_dist"); ok {
		bgpCtxAfPolAttr.LocalDist = LocalDist.(string)
	}
	if MaxEcmp, ok := d.GetOk("max_ecmp"); ok {
		bgpCtxAfPolAttr.MaxEcmp = MaxEcmp.(string)
	}
	if MaxEcmpIbgp, ok := d.GetOk("max_ecmp_ibgp"); ok {
		bgpCtxAfPolAttr.MaxEcmpIbgp = MaxEcmpIbgp.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpCtxAfPolAttr.NameAlias = NameAlias.(string)
	}
	bgpCtxAfPol := models.NewBGPAddressFamilyContextPolicy(fmt.Sprintf("bgpCtxAfP-%s", name), TenantDn, desc, bgpCtxAfPolAttr)

	bgpCtxAfPol.Status = "modified"

	err := aciClient.Save(bgpCtxAfPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bgpCtxAfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBGPAddressFamilyContextPolicyRead(ctx, d, m)

}

func resourceAciBGPAddressFamilyContextPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bgpCtxAfPol, err := getRemoteBGPAddressFamilyContextPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setBGPAddressFamilyContextPolicyAttributes(bgpCtxAfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBGPAddressFamilyContextPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bgpCtxAfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

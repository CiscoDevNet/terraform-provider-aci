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

func resourceAciLinkLevelPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLinkLevelPolicyCreate,
		UpdateContext: resourceAciLinkLevelPolicyUpdate,
		ReadContext:   resourceAciLinkLevelPolicyRead,
		DeleteContext: resourceAciLinkLevelPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLinkLevelPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"auto_neg": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"on",
					"off",
				}, false),
			},

			"fec_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"inherit",
					"cl91-rs-fec",
					"cl74-fc-fec",
					"ieee-rs-fec",
					"cons16-rs-fec",
					"kp-fec",
					"disable-fec",
				}, false),
			},

			"link_debounce": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"speed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unknown",
					"100M",
					"1G",
					"10G",
					"25G",
					"40G",
					"50G",
					"100G",
					"200G",
					"400G",
					"inherit",
				}, false),
			},
		}),
	}
}
func getRemoteLinkLevelPolicy(client *client.Client, dn string) (*models.LinkLevelPolicy, error) {
	fabricHIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricHIfPol := models.LinkLevelPolicyFromContainer(fabricHIfPolCont)

	if fabricHIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("LinkLevelPolicy %s not found", fabricHIfPol.DistinguishedName)
	}

	return fabricHIfPol, nil
}

func setLinkLevelPolicyAttributes(fabricHIfPol *models.LinkLevelPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fabricHIfPol.DistinguishedName)
	d.Set("description", fabricHIfPol.Description)
	fabricHIfPolMap, err := fabricHIfPol.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", fabricHIfPolMap["name"])

	d.Set("annotation", fabricHIfPolMap["annotation"])
	d.Set("auto_neg", fabricHIfPolMap["autoNeg"])
	d.Set("fec_mode", fabricHIfPolMap["fecMode"])
	d.Set("link_debounce", fabricHIfPolMap["linkDebounce"])
	d.Set("name_alias", fabricHIfPolMap["nameAlias"])
	d.Set("speed", fabricHIfPolMap["speed"])
	return d, nil
}

func resourceAciLinkLevelPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fabricHIfPol, err := getRemoteLinkLevelPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLinkLevelPolicyAttributes(fabricHIfPol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLinkLevelPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LinkLevelPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fabricHIfPolAttr := models.LinkLevelPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricHIfPolAttr.Annotation = Annotation.(string)
	} else {
		fabricHIfPolAttr.Annotation = "{}"
	}
	if AutoNeg, ok := d.GetOk("auto_neg"); ok {
		fabricHIfPolAttr.AutoNeg = AutoNeg.(string)
	}
	if FecMode, ok := d.GetOk("fec_mode"); ok {
		fabricHIfPolAttr.FecMode = FecMode.(string)
	}
	if LinkDebounce, ok := d.GetOk("link_debounce"); ok {
		fabricHIfPolAttr.LinkDebounce = LinkDebounce.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fabricHIfPolAttr.NameAlias = NameAlias.(string)
	}
	if Speed, ok := d.GetOk("speed"); ok {
		fabricHIfPolAttr.Speed = Speed.(string)
	}
	fabricHIfPol := models.NewLinkLevelPolicy(fmt.Sprintf("infra/hintfpol-%s", name), "uni", desc, fabricHIfPolAttr)

	err := aciClient.Save(fabricHIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricHIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLinkLevelPolicyRead(ctx, d, m)
}

func resourceAciLinkLevelPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LinkLevelPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fabricHIfPolAttr := models.LinkLevelPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricHIfPolAttr.Annotation = Annotation.(string)
	} else {
		fabricHIfPolAttr.Annotation = "{}"
	}
	if AutoNeg, ok := d.GetOk("auto_neg"); ok {
		fabricHIfPolAttr.AutoNeg = AutoNeg.(string)
	}
	if FecMode, ok := d.GetOk("fec_mode"); ok {
		fabricHIfPolAttr.FecMode = FecMode.(string)
	}
	if LinkDebounce, ok := d.GetOk("link_debounce"); ok {
		fabricHIfPolAttr.LinkDebounce = LinkDebounce.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fabricHIfPolAttr.NameAlias = NameAlias.(string)
	}
	if Speed, ok := d.GetOk("speed"); ok {
		fabricHIfPolAttr.Speed = Speed.(string)
	}
	fabricHIfPol := models.NewLinkLevelPolicy(fmt.Sprintf("infra/hintfpol-%s", name), "uni", desc, fabricHIfPolAttr)

	fabricHIfPol.Status = "modified"

	err := aciClient.Save(fabricHIfPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricHIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLinkLevelPolicyRead(ctx, d, m)

}

func resourceAciLinkLevelPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fabricHIfPol, err := getRemoteLinkLevelPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setLinkLevelPolicyAttributes(fabricHIfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLinkLevelPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fabricHIfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

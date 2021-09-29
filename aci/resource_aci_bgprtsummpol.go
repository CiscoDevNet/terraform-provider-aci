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

func resourceAciBgpRouteSummarization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBgpRouteSummarizationCreate,
		UpdateContext: resourceAciBgpRouteSummarizationUpdate,
		ReadContext:   resourceAciBgpRouteSummarizationRead,
		DeleteContext: resourceAciBgpRouteSummarizationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBgpRouteSummarizationImport,
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

			"attrmap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"as-set", "none",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteBgpRouteSummarization(client *client.Client, dn string) (*models.BgpRouteSummarization, error) {
	bgpRtSummPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpRtSummPol := models.BgpRouteSummarizationFromContainer(bgpRtSummPolCont)

	if bgpRtSummPol.DistinguishedName == "" {
		return nil, fmt.Errorf("BgpRouteSummarization %s not found", bgpRtSummPol.DistinguishedName)
	}

	return bgpRtSummPol, nil
}

func setBgpRouteSummarizationAttributes(bgpRtSummPol *models.BgpRouteSummarization, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(bgpRtSummPol.DistinguishedName)
	d.Set("description", bgpRtSummPol.Description)
	dn := d.Id()
	if dn != bgpRtSummPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	bgpRtSummPolMap, err := bgpRtSummPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/bgprtsum-%s", bgpRtSummPolMap["name"])))

	d.Set("name", bgpRtSummPolMap["name"])

	d.Set("annotation", bgpRtSummPolMap["annotation"])
	d.Set("attrmap", bgpRtSummPolMap["attrmap"])
	d.Set("ctrl", bgpRtSummPolMap["ctrl"])
	if bgpRtSummPolMap["ctrl"] == "" {
		d.Set("ctrl", "none")
	}
	d.Set("name_alias", bgpRtSummPolMap["nameAlias"])
	return d, nil
}

func resourceAciBgpRouteSummarizationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bgpRtSummPol, err := getRemoteBgpRouteSummarization(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setBgpRouteSummarizationAttributes(bgpRtSummPol, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBgpRouteSummarizationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BgpRouteSummarization: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpRtSummPolAttr := models.BgpRouteSummarizationAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpRtSummPolAttr.Annotation = Annotation.(string)
	} else {
		bgpRtSummPolAttr.Annotation = "{}"
	}
	if Attrmap, ok := d.GetOk("attrmap"); ok {
		bgpRtSummPolAttr.Attrmap = Attrmap.(string)
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		bgpRtSummPolAttr.Ctrl = Ctrl.(string)
		if Ctrl.(string) == "none" {
			bgpRtSummPolAttr.Ctrl = "{}"
		}
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpRtSummPolAttr.NameAlias = NameAlias.(string)
	}
	bgpRtSummPol := models.NewBgpRouteSummarization(fmt.Sprintf("bgprtsum-%s", name), TenantDn, desc, bgpRtSummPolAttr)

	err := aciClient.Save(bgpRtSummPol)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	d.SetId(bgpRtSummPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBgpRouteSummarizationRead(ctx, d, m)
}

func resourceAciBgpRouteSummarizationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BgpRouteSummarization: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpRtSummPolAttr := models.BgpRouteSummarizationAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpRtSummPolAttr.Annotation = Annotation.(string)
	} else {
		bgpRtSummPolAttr.Annotation = "{}"
	}
	if Attrmap, ok := d.GetOk("attrmap"); ok {
		bgpRtSummPolAttr.Attrmap = Attrmap.(string)
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		bgpRtSummPolAttr.Ctrl = Ctrl.(string)
		if Ctrl.(string) == "none" {
			bgpRtSummPolAttr.Ctrl = "{}"
		}
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpRtSummPolAttr.NameAlias = NameAlias.(string)
	}
	bgpRtSummPol := models.NewBgpRouteSummarization(fmt.Sprintf("bgprtsum-%s", name), TenantDn, desc, bgpRtSummPolAttr)

	bgpRtSummPol.Status = "modified"

	err := aciClient.Save(bgpRtSummPol)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	d.SetId(bgpRtSummPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBgpRouteSummarizationRead(ctx, d, m)

}

func resourceAciBgpRouteSummarizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bgpRtSummPol, err := getRemoteBgpRouteSummarization(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setBgpRouteSummarizationAttributes(bgpRtSummPol, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBgpRouteSummarizationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bgpRtSummPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

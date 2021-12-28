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

func resourceAciOspfRouteSummarization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciOspfRouteSummarizationCreate,
		UpdateContext: resourceAciOspfRouteSummarizationUpdate,
		ReadContext:   resourceAciOspfRouteSummarizationRead,
		DeleteContext: resourceAciOspfRouteSummarizationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciOspfRouteSummarizationImport,
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

			"cost": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"inter_area_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteOspfRouteSummarization(client *client.Client, dn string) (*models.OspfRouteSummarization, error) {
	ospfRtSummPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	ospfRtSummPol := models.OspfRouteSummarizationFromContainer(ospfRtSummPolCont)

	if ospfRtSummPol.DistinguishedName == "" {
		return nil, fmt.Errorf("OspfRouteSummarization %s not found", ospfRtSummPol.DistinguishedName)
	}

	return ospfRtSummPol, nil
}

func setOspfRouteSummarizationAttributes(ospfRtSummPol *models.OspfRouteSummarization, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(ospfRtSummPol.DistinguishedName)
	d.Set("description", ospfRtSummPol.Description)
	dn := d.Id()
	if dn != ospfRtSummPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	ospfRtSummPolMap, err := ospfRtSummPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/ospfrtsumm-%s", ospfRtSummPolMap["name"])))

	d.Set("name", ospfRtSummPolMap["name"])

	d.Set("annotation", ospfRtSummPolMap["annotation"])
	d.Set("cost", ospfRtSummPolMap["cost"])
	d.Set("inter_area_enabled", ospfRtSummPolMap["interAreaEnabled"])
	d.Set("name_alias", ospfRtSummPolMap["nameAlias"])
	d.Set("tag", ospfRtSummPolMap["tag"])
	return d, nil
}

func resourceAciOspfRouteSummarizationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	ospfRtSummPol, err := getRemoteOspfRouteSummarization(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setOspfRouteSummarizationAttributes(ospfRtSummPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciOspfRouteSummarizationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OspfRouteSummarization: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	ospfRtSummPolAttr := models.OspfRouteSummarizationAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ospfRtSummPolAttr.Annotation = Annotation.(string)
	} else {
		ospfRtSummPolAttr.Annotation = "{}"
	}
	if Cost, ok := d.GetOk("cost"); ok {
		ospfRtSummPolAttr.Cost = Cost.(string)
	}
	if InterAreaEnabled, ok := d.GetOk("inter_area_enabled"); ok {
		ospfRtSummPolAttr.InterAreaEnabled = InterAreaEnabled.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ospfRtSummPolAttr.NameAlias = NameAlias.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		ospfRtSummPolAttr.Tag = Tag.(string)
	}
	ospfRtSummPol := models.NewOspfRouteSummarization(fmt.Sprintf("ospfrtsumm-%s", name), TenantDn, desc, ospfRtSummPolAttr)

	err := aciClient.Save(ospfRtSummPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ospfRtSummPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciOspfRouteSummarizationRead(ctx, d, m)
}

func resourceAciOspfRouteSummarizationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OspfRouteSummarization: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	ospfRtSummPolAttr := models.OspfRouteSummarizationAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ospfRtSummPolAttr.Annotation = Annotation.(string)
	} else {
		ospfRtSummPolAttr.Annotation = "{}"
	}
	if Cost, ok := d.GetOk("cost"); ok {
		ospfRtSummPolAttr.Cost = Cost.(string)
	}
	if InterAreaEnabled, ok := d.GetOk("inter_area_enabled"); ok {
		ospfRtSummPolAttr.InterAreaEnabled = InterAreaEnabled.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ospfRtSummPolAttr.NameAlias = NameAlias.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		ospfRtSummPolAttr.Tag = Tag.(string)
	}
	ospfRtSummPol := models.NewOspfRouteSummarization(fmt.Sprintf("ospfrtsumm-%s", name), TenantDn, desc, ospfRtSummPolAttr)

	ospfRtSummPol.Status = "modified"

	err := aciClient.Save(ospfRtSummPol)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ospfRtSummPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciOspfRouteSummarizationRead(ctx, d, m)

}

func resourceAciOspfRouteSummarizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	ospfRtSummPol, err := getRemoteOspfRouteSummarization(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setOspfRouteSummarizationAttributes(ospfRtSummPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciOspfRouteSummarizationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "ospfRtSummPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

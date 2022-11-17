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

func resourceAciFabricNodeControl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciFabricNodeControlCreate,
		UpdateContext: resourceAciFabricNodeControlUpdate,
		ReadContext:   resourceAciFabricNodeControlRead,
		DeleteContext: resourceAciFabricNodeControlDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFabricNodeControlImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"control": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Dom",
					"None",
				}, false),
			},
			"feature_sel": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"analytics",
					"netflow",
					"telemetry",
				}, false),
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteFabricNodeControl(client *client.Client, dn string) (*models.FabricNodeControl, error) {
	fabricNodeControlCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fabricNodeControl := models.FabricNodeControlFromContainer(fabricNodeControlCont)
	if fabricNodeControl.DistinguishedName == "" {
		return nil, fmt.Errorf("FabricNodeControl %s not found", fabricNodeControl.DistinguishedName)
	}
	return fabricNodeControl, nil
}

func setFabricNodeControlAttributes(fabricNodeControl *models.FabricNodeControl, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fabricNodeControl.DistinguishedName)
	d.Set("description", fabricNodeControl.Description)
	fabricNodeControlMap, err := fabricNodeControl.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", fabricNodeControlMap["annotation"])
	if fabricNodeControlMap["control"] == "" {
		d.Set("control", "None")
	} else {
		d.Set("control", fabricNodeControlMap["control"])
	}
	d.Set("feature_sel", fabricNodeControlMap["featureSel"])
	d.Set("name", fabricNodeControlMap["name"])
	d.Set("name_alias", fabricNodeControlMap["nameAlias"])
	return d, nil
}

func resourceAciFabricNodeControlImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fabricNodeControl, err := getRemoteFabricNodeControl(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setFabricNodeControlAttributes(fabricNodeControl, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFabricNodeControlCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FabricNodeControl: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	fabricNodeControlAttr := models.FabricNodeControlAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricNodeControlAttr.Annotation = Annotation.(string)
	} else {
		fabricNodeControlAttr.Annotation = "{}"
	}

	if Control, ok := d.GetOk("control"); ok {
		fabricNodeControlAttr.Control = Control.(string)
	} else {
		fabricNodeControlAttr.Control = ""
	}

	if FeatureSel, ok := d.GetOk("feature_sel"); ok {
		fabricNodeControlAttr.FeatureSel = FeatureSel.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fabricNodeControlAttr.Name = Name.(string)
	}
	fabricNodeControl := models.NewFabricNodeControl(fmt.Sprintf("fabric/nodecontrol-%s", name), "uni", desc, nameAlias, fabricNodeControlAttr)
	err := aciClient.Save(fabricNodeControl)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricNodeControl.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciFabricNodeControlRead(ctx, d, m)
}

func resourceAciFabricNodeControlUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FabricNodeControl: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	fabricNodeControlAttr := models.FabricNodeControlAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricNodeControlAttr.Annotation = Annotation.(string)
	} else {
		fabricNodeControlAttr.Annotation = "{}"
	}

	if Control, ok := d.GetOk("control"); ok {
		fabricNodeControlAttr.Control = Control.(string)
	} else {
		fabricNodeControlAttr.Control = ""
	}

	if FeatureSel, ok := d.GetOk("feature_sel"); ok {
		fabricNodeControlAttr.FeatureSel = FeatureSel.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fabricNodeControlAttr.Name = Name.(string)
	}
	fabricNodeControl := models.NewFabricNodeControl(fmt.Sprintf("fabric/nodecontrol-%s", name), "uni", desc, nameAlias, fabricNodeControlAttr)
	fabricNodeControl.Status = "modified"
	err := aciClient.Save(fabricNodeControl)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricNodeControl.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciFabricNodeControlRead(ctx, d, m)
}

func resourceAciFabricNodeControlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fabricNodeControl, err := getRemoteFabricNodeControl(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setFabricNodeControlAttributes(fabricNodeControl, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciFabricNodeControlDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fabricNodeControl")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

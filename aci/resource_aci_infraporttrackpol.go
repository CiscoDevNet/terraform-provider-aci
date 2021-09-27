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

func resourceAciPortTracking() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciPortTrackingCreate,
		UpdateContext: resourceAciPortTrackingUpdate,
		ReadContext:   resourceAciPortTrackingRead,
		DeleteContext: resourceAciPortTrackingDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPortTrackingImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"off",
					"on",
				}, false),
			},
			"delay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"include_apic_ports": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"minlinks": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemotePortTracking(client *client.Client, dn string) (*models.PortTracking, error) {
	infraPortTrackPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	infraPortTrackPol := models.PortTrackingFromContainer(infraPortTrackPolCont)
	if infraPortTrackPol.DistinguishedName == "" {
		return nil, fmt.Errorf("PortTracking %s not found", infraPortTrackPol.DistinguishedName)
	}
	return infraPortTrackPol, nil
}

func setPortTrackingAttributes(infraPortTrackPol *models.PortTracking, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraPortTrackPol.DistinguishedName)
	d.Set("description", infraPortTrackPol.Description)
	infraPortTrackPolMap, err := infraPortTrackPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("admin_st", infraPortTrackPolMap["adminSt"])
	d.Set("annotation", infraPortTrackPolMap["annotation"])
	d.Set("delay", infraPortTrackPolMap["delay"])
	d.Set("include_apic_ports", infraPortTrackPolMap["includeApicPorts"])
	d.Set("minlinks", infraPortTrackPolMap["minlinks"])
	d.Set("name_alias", infraPortTrackPolMap["nameAlias"])
	return d, nil
}

func resourceAciPortTrackingImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraPortTrackPol, err := getRemotePortTracking(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPortTrackingAttributes(infraPortTrackPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPortTrackingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PortTracking: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	infraPortTrackPolAttr := models.PortTrackingAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraPortTrackPolAttr.Annotation = Annotation.(string)
	} else {
		infraPortTrackPolAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_st"); ok {
		infraPortTrackPolAttr.AdminSt = AdminSt.(string)
	}

	if Delay, ok := d.GetOk("delay"); ok {
		infraPortTrackPolAttr.Delay = Delay.(string)
	}

	if IncludeApicPorts, ok := d.GetOk("include_apic_ports"); ok {
		infraPortTrackPolAttr.IncludeApicPorts = IncludeApicPorts.(string)
	}

	if Minlinks, ok := d.GetOk("minlinks"); ok {
		infraPortTrackPolAttr.Minlinks = Minlinks.(string)
	}

	infraPortTrackPolAttr.Name = "default"

	infraPortTrackPol := models.NewPortTracking(fmt.Sprintf("infra/trackEqptFabP-%s", name), "uni", desc, nameAlias, infraPortTrackPolAttr)
	infraPortTrackPol.Status = "modified"
	err := aciClient.Save(infraPortTrackPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraPortTrackPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciPortTrackingRead(ctx, d, m)
}

func resourceAciPortTrackingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PortTracking: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	infraPortTrackPolAttr := models.PortTrackingAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraPortTrackPolAttr.Annotation = Annotation.(string)
	} else {
		infraPortTrackPolAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_st"); ok {
		infraPortTrackPolAttr.AdminSt = AdminSt.(string)
	}

	if Delay, ok := d.GetOk("delay"); ok {
		infraPortTrackPolAttr.Delay = Delay.(string)
	}

	if IncludeApicPorts, ok := d.GetOk("include_apic_ports"); ok {
		infraPortTrackPolAttr.IncludeApicPorts = IncludeApicPorts.(string)
	}

	if Minlinks, ok := d.GetOk("minlinks"); ok {
		infraPortTrackPolAttr.Minlinks = Minlinks.(string)
	}

	infraPortTrackPolAttr.Name = "default"

	infraPortTrackPol := models.NewPortTracking(fmt.Sprintf("infra/trackEqptFabP-%s", name), "uni", desc, nameAlias, infraPortTrackPolAttr)
	infraPortTrackPol.Status = "modified"
	err := aciClient.Save(infraPortTrackPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraPortTrackPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciPortTrackingRead(ctx, d, m)
}

func resourceAciPortTrackingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraPortTrackPol, err := getRemotePortTracking(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setPortTrackingAttributes(infraPortTrackPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciPortTrackingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name epLoopProtectP cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}

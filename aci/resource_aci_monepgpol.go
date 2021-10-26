package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciMonitoringPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMonitoringPolicyCreate,
		UpdateContext: resourceAciMonitoringPolicyUpdate,
		ReadContext:   resourceAciMonitoringPolicyRead,
		DeleteContext: resourceAciMonitoringPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMonitoringPolicyImport,
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

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteMonitoringPolicy(client *client.Client, dn string) (*models.MonitoringPolicy, error) {
	monEPGPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	monEPGPol := models.MonitoringPolicyFromContainer(monEPGPolCont)

	if monEPGPol.DistinguishedName == "" {
		return nil, fmt.Errorf("MonitoringPolicy %s not found", monEPGPol.DistinguishedName)
	}

	return monEPGPol, nil
}

func setMonitoringPolicyAttributes(monEPGPol *models.MonitoringPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(monEPGPol.DistinguishedName)
	d.Set("description", monEPGPol.Description)
	// d.Set("tenant_dn", GetParentDn(monEPGPol.DistinguishedName))
	if dn != monEPGPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	monEPGPolMap, err := monEPGPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", monEPGPolMap["name"])
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/monepg-%s", monEPGPolMap["name"])))
	d.Set("annotation", monEPGPolMap["annotation"])
	d.Set("name_alias", monEPGPolMap["nameAlias"])
	return d, nil
}

func resourceAciMonitoringPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	monEPGPol, err := getRemoteMonitoringPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	monEPGPolMap, err := monEPGPol.ToMap()
	if err != nil {
		return nil, err
	}
	name := monEPGPolMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/monepg-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setMonitoringPolicyAttributes(monEPGPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMonitoringPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MonitoringPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	monEPGPolAttr := models.MonitoringPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		monEPGPolAttr.Annotation = Annotation.(string)
	} else {
		monEPGPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		monEPGPolAttr.NameAlias = NameAlias.(string)
	}
	monEPGPol := models.NewMonitoringPolicy(fmt.Sprintf("monepg-%s", name), TenantDn, desc, monEPGPolAttr)

	err := aciClient.Save(monEPGPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(monEPGPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciMonitoringPolicyRead(ctx, d, m)
}

func resourceAciMonitoringPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MonitoringPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	monEPGPolAttr := models.MonitoringPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		monEPGPolAttr.Annotation = Annotation.(string)
	} else {
		monEPGPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		monEPGPolAttr.NameAlias = NameAlias.(string)
	}
	monEPGPol := models.NewMonitoringPolicy(fmt.Sprintf("monepg-%s", name), TenantDn, desc, monEPGPolAttr)

	monEPGPol.Status = "modified"

	err := aciClient.Save(monEPGPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(monEPGPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciMonitoringPolicyRead(ctx, d, m)

}

func resourceAciMonitoringPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	monEPGPol, err := getRemoteMonitoringPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setMonitoringPolicyAttributes(monEPGPol, d)

	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciMonitoringPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "monEPGPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

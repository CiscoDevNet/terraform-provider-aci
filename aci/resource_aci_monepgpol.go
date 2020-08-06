package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciMonitoringPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciMonitoringPolicyCreate,
		Update: resourceAciMonitoringPolicyUpdate,
		Read:   resourceAciMonitoringPolicyRead,
		Delete: resourceAciMonitoringPolicyDelete,

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

func setMonitoringPolicyAttributes(monEPGPol *models.MonitoringPolicy, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(monEPGPol.DistinguishedName)
	d.Set("description", monEPGPol.Description)
	// d.Set("tenant_dn", GetParentDn(monEPGPol.DistinguishedName))
	if dn != monEPGPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	monEPGPolMap, _ := monEPGPol.ToMap()

	d.Set("name", monEPGPolMap["name"])

	d.Set("annotation", monEPGPolMap["annotation"])
	d.Set("name_alias", monEPGPolMap["nameAlias"])
	return d
}

func resourceAciMonitoringPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	monEPGPol, err := getRemoteMonitoringPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setMonitoringPolicyAttributes(monEPGPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMonitoringPolicyCreate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(monEPGPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciMonitoringPolicyRead(d, m)
}

func resourceAciMonitoringPolicyUpdate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(monEPGPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciMonitoringPolicyRead(d, m)

}

func resourceAciMonitoringPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	monEPGPol, err := getRemoteMonitoringPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setMonitoringPolicyAttributes(monEPGPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciMonitoringPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "monEPGPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

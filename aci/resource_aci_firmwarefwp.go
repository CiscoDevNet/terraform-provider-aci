package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciFirmwarePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciFirmwarePolicyCreate,
		Update: resourceAciFirmwarePolicyUpdate,
		Read:   resourceAciFirmwarePolicyRead,
		Delete: resourceAciFirmwarePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFirmwarePolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"effective_on_reboot": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ignore_compat": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"internal_label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"version_check_override": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteFirmwarePolicy(client *client.Client, dn string) (*models.FirmwarePolicy, error) {
	firmwareFwPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	firmwareFwP := models.FirmwarePolicyFromContainer(firmwareFwPCont)

	if firmwareFwP.DistinguishedName == "" {
		return nil, fmt.Errorf("FirmwarePolicy %s not found", firmwareFwP.DistinguishedName)
	}

	return firmwareFwP, nil
}

func setFirmwarePolicyAttributes(firmwareFwP *models.FirmwarePolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(firmwareFwP.DistinguishedName)
	d.Set("description", firmwareFwP.Description)
	firmwareFwPMap, _ := firmwareFwP.ToMap()

	d.Set("name", firmwareFwPMap["name"])

	d.Set("annotation", firmwareFwPMap["annotation"])
	d.Set("effective_on_reboot", firmwareFwPMap["effectiveOnReboot"])
	d.Set("ignore_compat", firmwareFwPMap["ignoreCompat"])
	d.Set("internal_label", firmwareFwPMap["internalLabel"])
	d.Set("name_alias", firmwareFwPMap["nameAlias"])
	d.Set("version", firmwareFwPMap["version"])
	d.Set("version_check_override", firmwareFwPMap["versionCheckOverride"])
	return d
}

func resourceAciFirmwarePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	firmwareFwP, err := getRemoteFirmwarePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setFirmwarePolicyAttributes(firmwareFwP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFirmwarePolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FirmwarePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	firmwareFwPAttr := models.FirmwarePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		firmwareFwPAttr.Annotation = Annotation.(string)
	}
	if EffectiveOnReboot, ok := d.GetOk("effective_on_reboot"); ok {
		firmwareFwPAttr.EffectiveOnReboot = EffectiveOnReboot.(string)
	}
	if IgnoreCompat, ok := d.GetOk("ignore_compat"); ok {
		firmwareFwPAttr.IgnoreCompat = IgnoreCompat.(string)
	}
	if InternalLabel, ok := d.GetOk("internal_label"); ok {
		firmwareFwPAttr.InternalLabel = InternalLabel.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		firmwareFwPAttr.NameAlias = NameAlias.(string)
	}
	if Version, ok := d.GetOk("version"); ok {
		firmwareFwPAttr.Version = Version.(string)
	}
	if VersionCheckOverride, ok := d.GetOk("version_check_override"); ok {
		firmwareFwPAttr.VersionCheckOverride = VersionCheckOverride.(string)
	}
	firmwareFwP := models.NewFirmwarePolicy(fmt.Sprintf("fabric/fwpol-%s", name), "uni", desc, firmwareFwPAttr)

	err := aciClient.Save(firmwareFwP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(firmwareFwP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFirmwarePolicyRead(d, m)
}

func resourceAciFirmwarePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FirmwarePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	firmwareFwPAttr := models.FirmwarePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		firmwareFwPAttr.Annotation = Annotation.(string)
	}
	if EffectiveOnReboot, ok := d.GetOk("effective_on_reboot"); ok {
		firmwareFwPAttr.EffectiveOnReboot = EffectiveOnReboot.(string)
	}
	if IgnoreCompat, ok := d.GetOk("ignore_compat"); ok {
		firmwareFwPAttr.IgnoreCompat = IgnoreCompat.(string)
	}
	if InternalLabel, ok := d.GetOk("internal_label"); ok {
		firmwareFwPAttr.InternalLabel = InternalLabel.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		firmwareFwPAttr.NameAlias = NameAlias.(string)
	}
	if Version, ok := d.GetOk("version"); ok {
		firmwareFwPAttr.Version = Version.(string)
	}
	if VersionCheckOverride, ok := d.GetOk("version_check_override"); ok {
		firmwareFwPAttr.VersionCheckOverride = VersionCheckOverride.(string)
	}
	firmwareFwP := models.NewFirmwarePolicy(fmt.Sprintf("fabric/fwpol-%s", name), "uni", desc, firmwareFwPAttr)

	firmwareFwP.Status = "modified"

	err := aciClient.Save(firmwareFwP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(firmwareFwP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFirmwarePolicyRead(d, m)

}

func resourceAciFirmwarePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	firmwareFwP, err := getRemoteFirmwarePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setFirmwarePolicyAttributes(firmwareFwP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFirmwarePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "firmwareFwP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

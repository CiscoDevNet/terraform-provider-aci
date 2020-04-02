package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciCloudAvailabilityZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciCloudAvailabilityZoneCreate,
		Update: resourceAciCloudAvailabilityZoneUpdate,
		Read:   resourceAciCloudAvailabilityZoneRead,
		Delete: resourceAciCloudAvailabilityZoneDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudAvailabilityZoneImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_providers_region_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
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
func getRemoteCloudAvailabilityZone(client *client.Client, dn string) (*models.CloudAvailabilityZone, error) {
	cloudZoneCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudZone := models.CloudAvailabilityZoneFromContainer(cloudZoneCont)

	if cloudZone.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudAvailabilityZone %s not found", cloudZone.DistinguishedName)
	}

	return cloudZone, nil
}

func setCloudAvailabilityZoneAttributes(cloudZone *models.CloudAvailabilityZone, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(cloudZone.DistinguishedName)
	d.Set("description", cloudZone.Description)
	d.Set("cloud_providers_region_dn", GetParentDn(cloudZone.DistinguishedName))
	cloudZoneMap, _ := cloudZone.ToMap()

	d.Set("name", cloudZoneMap["name"])

	d.Set("annotation", cloudZoneMap["annotation"])
	d.Set("name_alias", cloudZoneMap["nameAlias"])
	return d
}

func resourceAciCloudAvailabilityZoneImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudZone, err := getRemoteCloudAvailabilityZone(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setCloudAvailabilityZoneAttributes(cloudZone, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudAvailabilityZoneCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudAvailabilityZone: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	CloudProvidersRegionDn := d.Get("cloud_providers_region_dn").(string)

	cloudZoneAttr := models.CloudAvailabilityZoneAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudZoneAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudZoneAttr.NameAlias = NameAlias.(string)
	}
	cloudZone := models.NewCloudAvailabilityZone(fmt.Sprintf("zone-%s", name), CloudProvidersRegionDn, desc, cloudZoneAttr)

	err := aciClient.Save(cloudZone)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(cloudZone.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCloudAvailabilityZoneRead(d, m)
}

func resourceAciCloudAvailabilityZoneUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudAvailabilityZone: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	CloudProvidersRegionDn := d.Get("cloud_providers_region_dn").(string)

	cloudZoneAttr := models.CloudAvailabilityZoneAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudZoneAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudZoneAttr.NameAlias = NameAlias.(string)
	}
	cloudZone := models.NewCloudAvailabilityZone(fmt.Sprintf("zone-%s", name), CloudProvidersRegionDn, desc, cloudZoneAttr)

	cloudZone.Status = "modified"

	err := aciClient.Save(cloudZone)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(cloudZone.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCloudAvailabilityZoneRead(d, m)

}

func resourceAciCloudAvailabilityZoneRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudZone, err := getRemoteCloudAvailabilityZone(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setCloudAvailabilityZoneAttributes(cloudZone, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCloudAvailabilityZoneDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudZone")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

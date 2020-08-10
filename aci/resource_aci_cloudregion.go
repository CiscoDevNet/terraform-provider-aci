package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciCloudProvidersRegion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciCloudProvidersRegionCreate,
		Update: resourceAciCloudProvidersRegionUpdate,
		Read:   resourceAciCloudProvidersRegionRead,
		Delete: resourceAciCloudProvidersRegionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudProvidersRegionImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_provider_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"admin_st": &schema.Schema{
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
func getRemoteCloudProvidersRegion(client *client.Client, dn string) (*models.CloudProvidersRegion, error) {
	cloudRegionCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudRegion := models.CloudProvidersRegionFromContainer(cloudRegionCont)

	if cloudRegion.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudProvidersRegion %s not found", cloudRegion.DistinguishedName)
	}

	return cloudRegion, nil
}

func setCloudProvidersRegionAttributes(cloudRegion *models.CloudProvidersRegion, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(cloudRegion.DistinguishedName)
	if dn != cloudRegion.DistinguishedName {
		d.Set("cloud_provider_profile_dn", "")
	}
	cloudRegionMap, _ := cloudRegion.ToMap()

	d.Set("name", cloudRegionMap["name"])

	d.Set("admin_st", cloudRegionMap["adminSt"])
	d.Set("name_alias", cloudRegionMap["nameAlias"])
	return d
}

func resourceAciCloudProvidersRegionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudRegion, err := getRemoteCloudProvidersRegion(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setCloudProvidersRegionAttributes(cloudRegion, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudProvidersRegionCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudProvidersRegion: Beginning Creation")
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	CloudProviderProfileDn := d.Get("cloud_provider_profile_dn").(string)

	cloudRegionAttr := models.CloudProvidersRegionAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		cloudRegionAttr.AdminSt = AdminSt.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudRegionAttr.NameAlias = NameAlias.(string)
	}
	cloudRegion := models.NewCloudProvidersRegion(fmt.Sprintf("region-%s", name), CloudProviderProfileDn, "", cloudRegionAttr)

	err := aciClient.Save(cloudRegion)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(cloudRegion.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCloudProvidersRegionRead(d, m)
}

func resourceAciCloudProvidersRegionUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudProvidersRegion: Beginning Update")

	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	CloudProviderProfileDn := d.Get("cloud_provider_profile_dn").(string)

	cloudRegionAttr := models.CloudProvidersRegionAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		cloudRegionAttr.AdminSt = AdminSt.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudRegionAttr.NameAlias = NameAlias.(string)
	}
	cloudRegion := models.NewCloudProvidersRegion(fmt.Sprintf("region-%s", name), CloudProviderProfileDn, "", cloudRegionAttr)

	cloudRegion.Status = "modified"

	err := aciClient.Save(cloudRegion)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(cloudRegion.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCloudProvidersRegionRead(d, m)

}

func resourceAciCloudProvidersRegionRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudRegion, err := getRemoteCloudProvidersRegion(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setCloudProvidersRegionAttributes(cloudRegion, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCloudProvidersRegionDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudRegion")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

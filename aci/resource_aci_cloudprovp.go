package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciCloudProviderProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciCloudProviderProfileCreate,
		Update: resourceAciCloudProviderProfileUpdate,
		Read:   resourceAciCloudProviderProfileRead,
		Delete: resourceAciCloudProviderProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudProviderProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"vendor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}
func getRemoteCloudProviderProfile(client *client.Client, dn string) (*models.CloudProviderProfile, error) {
	cloudProvPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudProvP := models.CloudProviderProfileFromContainer(cloudProvPCont)

	if cloudProvP.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudProviderProfile %s not found", cloudProvP.DistinguishedName)
	}

	return cloudProvP, nil
}

func setCloudProviderProfileAttributes(cloudProvP *models.CloudProviderProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(cloudProvP.DistinguishedName)
	cloudProvPMap, _ := cloudProvP.ToMap()

	d.Set("vendor", cloudProvPMap["vendor"])

	return d
}

func resourceAciCloudProviderProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudProvP, err := getRemoteCloudProviderProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setCloudProviderProfileAttributes(cloudProvP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudProviderProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudProviderProfile: Beginning Creation")
	aciClient := m.(*client.Client)

	vendor := d.Get("vendor").(string)

	cloudProvPAttr := models.CloudProviderProfileAttributes{}

	if Vendor, ok := d.GetOk("vendor"); ok {
		cloudProvPAttr.Vendor = Vendor.(string)
	}
	cloudProvP := models.NewCloudProviderProfile(fmt.Sprintf("clouddomp/provp-%s", vendor), "uni", "", cloudProvPAttr)

	err := aciClient.Save(cloudProvP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(cloudProvP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCloudProviderProfileRead(d, m)
}

func resourceAciCloudProviderProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudProviderProfile: Beginning Update")

	aciClient := m.(*client.Client)

	vendor := d.Get("vendor").(string)

	cloudProvPAttr := models.CloudProviderProfileAttributes{}

	if Vendor, ok := d.GetOk("vendor"); ok {
		cloudProvPAttr.Vendor = Vendor.(string)
	}
	cloudProvP := models.NewCloudProviderProfile(fmt.Sprintf("clouddomp/provp-%s", vendor), "uni", "", cloudProvPAttr)

	cloudProvP.Status = "modified"

	err := aciClient.Save(cloudProvP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(cloudProvP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCloudProviderProfileRead(d, m)

}

func resourceAciCloudProviderProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudProvP, err := getRemoteCloudProviderProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setCloudProviderProfileAttributes(cloudProvP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCloudProviderProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudProvP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciCloudApplicationcontainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciCloudApplicationcontainerCreate,
		Update: resourceAciCloudApplicationcontainerUpdate,
		Read:   resourceAciCloudApplicationcontainerRead,
		Delete: resourceAciCloudApplicationcontainerDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudApplicationcontainerImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
func getRemoteCloudApplicationcontainer(client *client.Client, dn string) (*models.CloudApplicationcontainer, error) {
	cloudAppCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudApp := models.CloudApplicationcontainerFromContainer(cloudAppCont)

	if cloudApp.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudApplicationcontainer %s not found", cloudApp.DistinguishedName)
	}

	return cloudApp, nil
}

func setCloudApplicationcontainerAttributes(cloudApp *models.CloudApplicationcontainer, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(cloudApp.DistinguishedName)
	d.Set("description", cloudApp.Description)
	d.Set("tenant_dn", GetParentDn(cloudApp.DistinguishedName))
	cloudAppMap, _ := cloudApp.ToMap()

	d.Set("name", cloudAppMap["name"])

	d.Set("annotation", cloudAppMap["annotation"])
	d.Set("name_alias", cloudAppMap["nameAlias"])
	return d
}

func resourceAciCloudApplicationcontainerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudApp, err := getRemoteCloudApplicationcontainer(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setCloudApplicationcontainerAttributes(cloudApp, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudApplicationcontainerCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudApplicationcontainer: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	cloudAppAttr := models.CloudApplicationcontainerAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudAppAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudAppAttr.NameAlias = NameAlias.(string)
	}
	cloudApp := models.NewCloudApplicationcontainer(fmt.Sprintf("cloudapp-%s", name), TenantDn, desc, cloudAppAttr)

	err := aciClient.Save(cloudApp)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(cloudApp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCloudApplicationcontainerRead(d, m)
}

func resourceAciCloudApplicationcontainerUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudApplicationcontainer: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	cloudAppAttr := models.CloudApplicationcontainerAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudAppAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudAppAttr.NameAlias = NameAlias.(string)
	}
	cloudApp := models.NewCloudApplicationcontainer(fmt.Sprintf("cloudapp-%s", name), TenantDn, desc, cloudAppAttr)

	cloudApp.Status = "modified"

	err := aciClient.Save(cloudApp)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(cloudApp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCloudApplicationcontainerRead(d, m)

}

func resourceAciCloudApplicationcontainerRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudApp, err := getRemoteCloudApplicationcontainer(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setCloudApplicationcontainerAttributes(cloudApp, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCloudApplicationcontainerDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudApp")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAciCloudEndpointSelector() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciCloudEndpointSelectorCreate,
		Update: resourceAciCloudEndpointSelectorUpdate,
		Read:   resourceAciCloudEndpointSelectorRead,
		Delete: resourceAciCloudEndpointSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudEndpointSelectorImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_e_pg_dn": &schema.Schema{
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

			"match_expression": &schema.Schema{
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
func getRemoteCloudEndpointSelector(client *client.Client, dn string) (*models.CloudEndpointSelector, error) {
	cloudEPSelectorCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudEPSelector := models.CloudEndpointSelectorFromContainer(cloudEPSelectorCont)

	if cloudEPSelector.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudEndpointSelector %s not found", cloudEPSelector.DistinguishedName)
	}

	return cloudEPSelector, nil
}

func setCloudEndpointSelectorAttributes(cloudEPSelector *models.CloudEndpointSelector, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(cloudEPSelector.DistinguishedName)
	d.Set("description", cloudEPSelector.Description)
	d.Set("cloud_e_pg_dn", GetParentDn(cloudEPSelector.DistinguishedName))
	cloudEPSelectorMap, _ := cloudEPSelector.ToMap()

	d.Set("name", cloudEPSelectorMap["name"])

	d.Set("annotation", cloudEPSelectorMap["annotation"])
	d.Set("match_expression", cloudEPSelectorMap["matchExpression"])
	d.Set("name_alias", cloudEPSelectorMap["nameAlias"])
	return d
}

func resourceAciCloudEndpointSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudEPSelector, err := getRemoteCloudEndpointSelector(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setCloudEndpointSelectorAttributes(cloudEPSelector, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudEndpointSelectorCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudEndpointSelector: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	CloudEPgDn := d.Get("cloud_e_pg_dn").(string)

	cloudEPSelectorAttr := models.CloudEndpointSelectorAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudEPSelectorAttr.Annotation = Annotation.(string)
	}
	if MatchExpression, ok := d.GetOk("match_expression"); ok {
		cloudEPSelectorAttr.MatchExpression = MatchExpression.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudEPSelectorAttr.NameAlias = NameAlias.(string)
	}
	cloudEPSelector := models.NewCloudEndpointSelector(fmt.Sprintf("epselector-%s", name), CloudEPgDn, desc, cloudEPSelectorAttr)

	err := aciClient.Save(cloudEPSelector)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(cloudEPSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCloudEndpointSelectorRead(d, m)
}

func resourceAciCloudEndpointSelectorUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudEndpointSelector: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	CloudEPgDn := d.Get("cloud_e_pg_dn").(string)

	cloudEPSelectorAttr := models.CloudEndpointSelectorAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudEPSelectorAttr.Annotation = Annotation.(string)
	}
	if MatchExpression, ok := d.GetOk("match_expression"); ok {
		cloudEPSelectorAttr.MatchExpression = MatchExpression.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudEPSelectorAttr.NameAlias = NameAlias.(string)
	}
	cloudEPSelector := models.NewCloudEndpointSelector(fmt.Sprintf("epselector-%s", name), CloudEPgDn, desc, cloudEPSelectorAttr)

	cloudEPSelector.Status = "modified"

	err := aciClient.Save(cloudEPSelector)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(cloudEPSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCloudEndpointSelectorRead(d, m)

}

func resourceAciCloudEndpointSelectorRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudEPSelector, err := getRemoteCloudEndpointSelector(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setCloudEndpointSelectorAttributes(cloudEPSelector, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCloudEndpointSelectorDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudEPSelector")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

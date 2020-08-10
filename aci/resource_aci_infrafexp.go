package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciFEXProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciFEXProfileCreate,
		Update: resourceAciFEXProfileUpdate,
		Read:   resourceAciFEXProfileRead,
		Delete: resourceAciFEXProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFEXProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

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
func getRemoteFEXProfile(client *client.Client, dn string) (*models.FEXProfile, error) {
	infraFexPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraFexP := models.FEXProfileFromContainer(infraFexPCont)

	if infraFexP.DistinguishedName == "" {
		return nil, fmt.Errorf("FEXProfile %s not found", infraFexP.DistinguishedName)
	}

	return infraFexP, nil
}

func setFEXProfileAttributes(infraFexP *models.FEXProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(infraFexP.DistinguishedName)
	d.Set("description", infraFexP.Description)
	infraFexPMap, _ := infraFexP.ToMap()

	d.Set("name", infraFexPMap["name"])

	d.Set("annotation", infraFexPMap["annotation"])
	d.Set("name_alias", infraFexPMap["nameAlias"])
	return d
}

func resourceAciFEXProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraFexP, err := getRemoteFEXProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setFEXProfileAttributes(infraFexP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFEXProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FEXProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraFexPAttr := models.FEXProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraFexPAttr.Annotation = Annotation.(string)
	} else {
		infraFexPAttr.Annotation = "{}"
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraFexPAttr.NameAlias = NameAlias.(string)
	}
	infraFexP := models.NewFEXProfile(fmt.Sprintf("infra/fexprof-%s", name), "uni", desc, infraFexPAttr)

	err := aciClient.Save(infraFexP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(infraFexP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFEXProfileRead(d, m)
}

func resourceAciFEXProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FEXProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraFexPAttr := models.FEXProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraFexPAttr.Annotation = Annotation.(string)
	} else {
		infraFexPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraFexPAttr.NameAlias = NameAlias.(string)
	}
	infraFexP := models.NewFEXProfile(fmt.Sprintf("infra/fexprof-%s", name), "uni", desc, infraFexPAttr)

	infraFexP.Status = "modified"

	err := aciClient.Save(infraFexP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(infraFexP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFEXProfileRead(d, m)

}

func resourceAciFEXProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraFexP, err := getRemoteFEXProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setFEXProfileAttributes(infraFexP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFEXProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraFexP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

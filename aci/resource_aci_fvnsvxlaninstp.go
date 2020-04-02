package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciVXLANPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciVXLANPoolCreate,
		Update: resourceAciVXLANPoolUpdate,
		Read:   resourceAciVXLANPoolRead,
		Delete: resourceAciVXLANPoolDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVXLANPoolImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

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
func getRemoteVXLANPool(client *client.Client, dn string) (*models.VXLANPool, error) {
	fvnsVxlanInstPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvnsVxlanInstP := models.VXLANPoolFromContainer(fvnsVxlanInstPCont)

	if fvnsVxlanInstP.DistinguishedName == "" {
		return nil, fmt.Errorf("VXLANPool %s not found", fvnsVxlanInstP.DistinguishedName)
	}

	return fvnsVxlanInstP, nil
}

func setVXLANPoolAttributes(fvnsVxlanInstP *models.VXLANPool, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvnsVxlanInstP.DistinguishedName)
	d.Set("description", fvnsVxlanInstP.Description)
	fvnsVxlanInstPMap, _ := fvnsVxlanInstP.ToMap()

	d.Set("name", fvnsVxlanInstPMap["name"])

	d.Set("annotation", fvnsVxlanInstPMap["annotation"])
	d.Set("name_alias", fvnsVxlanInstPMap["nameAlias"])
	return d
}

func resourceAciVXLANPoolImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvnsVxlanInstP, err := getRemoteVXLANPool(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setVXLANPoolAttributes(fvnsVxlanInstP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVXLANPoolCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] VXLANPool: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fvnsVxlanInstPAttr := models.VXLANPoolAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsVxlanInstPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsVxlanInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsVxlanInstP := models.NewVXLANPool(fmt.Sprintf("infra/vxlanns-%s", name), "uni", desc, fvnsVxlanInstPAttr)

	err := aciClient.Save(fvnsVxlanInstP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(fvnsVxlanInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciVXLANPoolRead(d, m)
}

func resourceAciVXLANPoolUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] VXLANPool: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fvnsVxlanInstPAttr := models.VXLANPoolAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsVxlanInstPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsVxlanInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsVxlanInstP := models.NewVXLANPool(fmt.Sprintf("infra/vxlanns-%s", name), "uni", desc, fvnsVxlanInstPAttr)

	fvnsVxlanInstP.Status = "modified"

	err := aciClient.Save(fvnsVxlanInstP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(fvnsVxlanInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciVXLANPoolRead(d, m)

}

func resourceAciVXLANPoolRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvnsVxlanInstP, err := getRemoteVXLANPool(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setVXLANPoolAttributes(fvnsVxlanInstP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciVXLANPoolDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvnsVxlanInstP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

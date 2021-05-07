package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciVSANPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciVSANPoolCreate,
		Update: resourceAciVSANPoolUpdate,
		Read:   resourceAciVSANPoolRead,
		Delete: resourceAciVSANPoolDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVSANPoolImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"alloc_mode": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"dynamic",
					"static",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteVSANPool(client *client.Client, dn string) (*models.VSANPool, error) {
	fvnsVsanInstPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvnsVsanInstP := models.VSANPoolFromContainer(fvnsVsanInstPCont)

	if fvnsVsanInstP.DistinguishedName == "" {
		return nil, fmt.Errorf("VSANPool %s not found", fvnsVsanInstP.DistinguishedName)
	}

	return fvnsVsanInstP, nil
}

func setVSANPoolAttributes(fvnsVsanInstP *models.VSANPool, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvnsVsanInstP.DistinguishedName)
	d.Set("description", fvnsVsanInstP.Description)
	fvnsVsanInstPMap, _ := fvnsVsanInstP.ToMap()

	d.Set("name", fvnsVsanInstPMap["name"])

	d.Set("alloc_mode", fvnsVsanInstPMap["allocMode"])
	d.Set("annotation", fvnsVsanInstPMap["annotation"])
	d.Set("name_alias", fvnsVsanInstPMap["nameAlias"])
	return d
}

func resourceAciVSANPoolImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvnsVsanInstP, err := getRemoteVSANPool(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setVSANPoolAttributes(fvnsVsanInstP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVSANPoolCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] VSANPool: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	allocMode := d.Get("alloc_mode").(string)

	fvnsVsanInstPAttr := models.VSANPoolAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsVsanInstPAttr.Annotation = Annotation.(string)
	} else {
		fvnsVsanInstPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsVsanInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsVsanInstP := models.NewVSANPool(fmt.Sprintf("infra/vsanns-[%s]-%s", name, allocMode), "uni", desc, fvnsVsanInstPAttr)

	err := aciClient.Save(fvnsVsanInstP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(fvnsVsanInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciVSANPoolRead(d, m)
}

func resourceAciVSANPoolUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] VSANPool: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	allocMode := d.Get("alloc_mode").(string)

	fvnsVsanInstPAttr := models.VSANPoolAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsVsanInstPAttr.Annotation = Annotation.(string)
	} else {
		fvnsVsanInstPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsVsanInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsVsanInstP := models.NewVSANPool(fmt.Sprintf("infra/vsanns-[%s]-%s", name, allocMode), "uni", desc, fvnsVsanInstPAttr)

	fvnsVsanInstP.Status = "modified"

	err := aciClient.Save(fvnsVsanInstP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(fvnsVsanInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciVSANPoolRead(d, m)

}

func resourceAciVSANPoolRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvnsVsanInstP, err := getRemoteVSANPool(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setVSANPoolAttributes(fvnsVsanInstP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciVSANPoolDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvnsVsanInstP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

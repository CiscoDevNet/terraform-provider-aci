package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciStaticPath() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciStaticPathCreate,
		Update: resourceAciStaticPathUpdate,
		Read:   resourceAciStaticPathRead,
		Delete: resourceAciStaticPathDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciStaticPathImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instr_imedcy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"primary_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteStaticPath(client *client.Client, dn string) (*models.StaticPath, error) {
	fvRsPathAttCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvRsPathAtt := models.StaticPathFromContainer(fvRsPathAttCont)

	if fvRsPathAtt.DistinguishedName == "" {
		return nil, fmt.Errorf("StaticPath %s not found", fvRsPathAtt.DistinguishedName)
	}

	return fvRsPathAtt, nil
}

func setStaticPathAttributes(fvRsPathAtt *models.StaticPath, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvRsPathAtt.DistinguishedName)
	d.Set("application_epg_dn", GetParentDn(fvRsPathAtt.DistinguishedName))
	fvRsPathAttMap, _ := fvRsPathAtt.ToMap()

	d.Set("tdn", fvRsPathAttMap["tDn"])

	d.Set("annotation", fvRsPathAttMap["annotation"])
	d.Set("encap", fvRsPathAttMap["encap"])
	d.Set("instr_imedcy", fvRsPathAttMap["instrImedcy"])
	d.Set("mode", fvRsPathAttMap["mode"])
	d.Set("primary_encap", fvRsPathAttMap["primaryEncap"])
	return d
}

func resourceAciStaticPathImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvRsPathAtt, err := getRemoteStaticPath(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setStaticPathAttributes(fvRsPathAtt, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciStaticPathCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] StaticPath: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	tDn := d.Get("tdn").(string)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	fvRsPathAttAttr := models.StaticPathAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsPathAttAttr.Annotation = Annotation.(string)
	}
	if Encap, ok := d.GetOk("encap"); ok {
		fvRsPathAttAttr.Encap = Encap.(string)
	}
	if InstrImedcy, ok := d.GetOk("instr_imedcy"); ok {
		fvRsPathAttAttr.InstrImedcy = InstrImedcy.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		fvRsPathAttAttr.Mode = Mode.(string)
	}
	if PrimaryEncap, ok := d.GetOk("primary_encap"); ok {
		fvRsPathAttAttr.PrimaryEncap = PrimaryEncap.(string)
	}

	fvRsPathAtt := models.NewStaticPath(fmt.Sprintf("rspathAtt-[%s]", tDn), ApplicationEPGDn, desc, fvRsPathAttAttr)

	err := aciClient.Save(fvRsPathAtt)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("tdn")

	d.Partial(false)

	d.SetId(fvRsPathAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciStaticPathRead(d, m)
}

func resourceAciStaticPathUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] StaticPath: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	tDn := d.Get("tdn").(string)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	fvRsPathAttAttr := models.StaticPathAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsPathAttAttr.Annotation = Annotation.(string)
	}
	if Encap, ok := d.GetOk("encap"); ok {
		fvRsPathAttAttr.Encap = Encap.(string)
	}
	if InstrImedcy, ok := d.GetOk("instr_imedcy"); ok {
		fvRsPathAttAttr.InstrImedcy = InstrImedcy.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		fvRsPathAttAttr.Mode = Mode.(string)
	}
	if PrimaryEncap, ok := d.GetOk("primary_encap"); ok {
		fvRsPathAttAttr.PrimaryEncap = PrimaryEncap.(string)
	}

	fvRsPathAtt := models.NewStaticPath(fmt.Sprintf("rspathAtt-[%s]", tDn), ApplicationEPGDn, desc, fvRsPathAttAttr)

	fvRsPathAtt.Status = "modified"

	err := aciClient.Save(fvRsPathAtt)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("tdn")

	d.Partial(false)

	d.SetId(fvRsPathAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciStaticPathRead(d, m)

}

func resourceAciStaticPathRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvRsPathAtt, err := getRemoteStaticPath(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setStaticPathAttributes(fvRsPathAtt, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciStaticPathDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvRsPathAtt")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

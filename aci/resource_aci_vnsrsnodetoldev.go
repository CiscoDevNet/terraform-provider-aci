package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciRelationfromaAbsNodetoanLDev() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciRelationfromaAbsNodetoanLDevCreate,
		Update: resourceAciRelationfromaAbsNodetoanLDevUpdate,
		Read:   resourceAciRelationfromaAbsNodetoanLDevRead,
		Delete: resourceAciRelationfromaAbsNodetoanLDevDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciRelationfromaAbsNodetoanLDevImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"function_node_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"t_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteRelationfromaAbsNodetoanLDev(client *client.Client, dn string) (*models.RelationfromaAbsNodetoanLDev, error) {
	vnsRsNodeToLDevCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsRsNodeToLDev := models.RelationfromaAbsNodetoanLDevFromContainer(vnsRsNodeToLDevCont)

	if vnsRsNodeToLDev.DistinguishedName == "" {
		return nil, fmt.Errorf("RelationfromaAbsNodetoanLDev %s not found", vnsRsNodeToLDev.DistinguishedName)
	}

	return vnsRsNodeToLDev, nil
}

func setRelationfromaAbsNodetoanLDevAttributes(vnsRsNodeToLDev *models.RelationfromaAbsNodetoanLDev, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(vnsRsNodeToLDev.DistinguishedName)
	d.Set("description", vnsRsNodeToLDev.Description)
	d.Set("function_node_dn", GetParentDn(vnsRsNodeToLDev.DistinguishedName))
	vnsRsNodeToLDevMap, _ := vnsRsNodeToLDev.ToMap()

	d.Set("annotation", vnsRsNodeToLDevMap["annotation"])
	d.Set("t_dn", vnsRsNodeToLDevMap["tDn"])
	return d
}

func resourceAciRelationfromaAbsNodetoanLDevImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vnsRsNodeToLDev, err := getRemoteRelationfromaAbsNodetoanLDev(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setRelationfromaAbsNodetoanLDevAttributes(vnsRsNodeToLDev, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciRelationfromaAbsNodetoanLDevCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] RelationfromaAbsNodetoanLDev: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	FunctionNodeDn := d.Get("function_node_dn").(string)

	vnsRsNodeToLDevAttr := models.RelationfromaAbsNodetoanLDevAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsRsNodeToLDevAttr.Annotation = Annotation.(string)
	}
	if TDn, ok := d.GetOk("t_dn"); ok {
		vnsRsNodeToLDevAttr.TDn = TDn.(string)
	}
	vnsRsNodeToLDev := models.NewRelationfromaAbsNodetoanLDev(fmt.Sprintf("rsNodeToLDev"), FunctionNodeDn, desc, vnsRsNodeToLDevAttr)

	err := aciClient.Save(vnsRsNodeToLDev)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	d.SetId(vnsRsNodeToLDev.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciRelationfromaAbsNodetoanLDevRead(d, m)
}

func resourceAciRelationfromaAbsNodetoanLDevUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] RelationfromaAbsNodetoanLDev: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	FunctionNodeDn := d.Get("function_node_dn").(string)

	vnsRsNodeToLDevAttr := models.RelationfromaAbsNodetoanLDevAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsRsNodeToLDevAttr.Annotation = Annotation.(string)
	}
	if TDn, ok := d.GetOk("t_dn"); ok {
		vnsRsNodeToLDevAttr.TDn = TDn.(string)
	}
	vnsRsNodeToLDev := models.NewRelationfromaAbsNodetoanLDev(fmt.Sprintf("rsNodeToLDev"), FunctionNodeDn, desc, vnsRsNodeToLDevAttr)

	vnsRsNodeToLDev.Status = "modified"

	err := aciClient.Save(vnsRsNodeToLDev)

	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	d.SetId(vnsRsNodeToLDev.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciRelationfromaAbsNodetoanLDevRead(d, m)

}

func resourceAciRelationfromaAbsNodetoanLDevRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vnsRsNodeToLDev, err := getRemoteRelationfromaAbsNodetoanLDev(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setRelationfromaAbsNodetoanLDevAttributes(vnsRsNodeToLDev, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciRelationfromaAbsNodetoanLDevDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vnsRsNodeToLDev")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

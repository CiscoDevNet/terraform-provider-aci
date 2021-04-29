package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciAutonomousSystemProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciAutonomousSystemProfileCreate,
		Update: resourceAciAutonomousSystemProfileUpdate,
		Read:   resourceAciAutonomousSystemProfileRead,
		Delete: resourceAciAutonomousSystemProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAutonomousSystemProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"asn": &schema.Schema{
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
func getRemoteAutonomousSystemProfile(client *client.Client, dn string) (*models.AutonomousSystemProfile, error) {
	cloudBgpAsPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudBgpAsP := models.AutonomousSystemProfileFromContainer(cloudBgpAsPCont)

	if cloudBgpAsP.DistinguishedName == "" {
		return nil, fmt.Errorf("AutonomousSystemProfile %s not found", cloudBgpAsP.DistinguishedName)
	}

	return cloudBgpAsP, nil
}

func setAutonomousSystemProfileAttributes(cloudBgpAsP *models.AutonomousSystemProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(cloudBgpAsP.DistinguishedName)
	d.Set("description", cloudBgpAsP.Description)
	cloudBgpAsPMap, _ := cloudBgpAsP.ToMap()

	d.Set("annotation", cloudBgpAsPMap["annotation"])
	d.Set("asn", cloudBgpAsPMap["asn"])
	d.Set("name_alias", cloudBgpAsPMap["nameAlias"])
	return d
}

func resourceAciAutonomousSystemProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudBgpAsP, err := getRemoteAutonomousSystemProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setAutonomousSystemProfileAttributes(cloudBgpAsP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAutonomousSystemProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AutonomousSystemProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	cloudBgpAsPAttr := models.AutonomousSystemProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudBgpAsPAttr.Annotation = Annotation.(string)
	} else {
		cloudBgpAsPAttr.Annotation = "{}"
	}
	if Asn, ok := d.GetOk("asn"); ok {
		cloudBgpAsPAttr.Asn = Asn.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudBgpAsPAttr.NameAlias = NameAlias.(string)
	}
	cloudBgpAsP := models.NewAutonomousSystemProfile(fmt.Sprintf("clouddomp/as"), "uni", desc, cloudBgpAsPAttr)

	err := aciClient.Save(cloudBgpAsP)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	d.SetId(cloudBgpAsP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciAutonomousSystemProfileRead(d, m)
}

func resourceAciAutonomousSystemProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AutonomousSystemProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	cloudBgpAsPAttr := models.AutonomousSystemProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudBgpAsPAttr.Annotation = Annotation.(string)
	} else {
		cloudBgpAsPAttr.Annotation = "{}"
	}
	if Asn, ok := d.GetOk("asn"); ok {
		cloudBgpAsPAttr.Asn = Asn.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudBgpAsPAttr.NameAlias = NameAlias.(string)
	}
	cloudBgpAsP := models.NewAutonomousSystemProfile(fmt.Sprintf("clouddomp/as"), "uni", desc, cloudBgpAsPAttr)

	cloudBgpAsP.Status = "modified"

	err := aciClient.Save(cloudBgpAsP)

	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	d.SetId(cloudBgpAsP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciAutonomousSystemProfileRead(d, m)

}

func resourceAciAutonomousSystemProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudBgpAsP, err := getRemoteAutonomousSystemProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setAutonomousSystemProfileAttributes(cloudBgpAsP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciAutonomousSystemProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudBgpAsP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

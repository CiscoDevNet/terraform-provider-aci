package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciSecurityDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciSecurityDomainCreate,
		Update: resourceAciSecurityDomainUpdate,
		Read:   resourceAciSecurityDomainRead,
		Delete: resourceAciSecurityDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSecurityDomainImport,
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
func getRemoteSecurityDomain(client *client.Client, dn string) (*models.SecurityDomain, error) {
	aaaDomainCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	aaaDomain := models.SecurityDomainFromContainer(aaaDomainCont)

	if aaaDomain.DistinguishedName == "" {
		return nil, fmt.Errorf("SecurityDomain %s not found", aaaDomain.DistinguishedName)
	}

	return aaaDomain, nil
}

func setSecurityDomainAttributes(aaaDomain *models.SecurityDomain, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(aaaDomain.DistinguishedName)
	d.Set("description", aaaDomain.Description)
	aaaDomainMap, _ := aaaDomain.ToMap()

	d.Set("name", aaaDomainMap["name"])

	d.Set("annotation", aaaDomainMap["annotation"])
	d.Set("name_alias", aaaDomainMap["nameAlias"])
	return d
}

func resourceAciSecurityDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	aaaDomain, err := getRemoteSecurityDomain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setSecurityDomainAttributes(aaaDomain, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSecurityDomainCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] SecurityDomain: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	aaaDomainAttr := models.SecurityDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaDomainAttr.Annotation = Annotation.(string)
	} else {
		aaaDomainAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		aaaDomainAttr.NameAlias = NameAlias.(string)
	}
	aaaDomain := models.NewSecurityDomain(fmt.Sprintf("userext/domain-%s", name), "uni", desc, aaaDomainAttr)

	err := aciClient.Save(aaaDomain)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(aaaDomain.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSecurityDomainRead(d, m)
}

func resourceAciSecurityDomainUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] SecurityDomain: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	aaaDomainAttr := models.SecurityDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaDomainAttr.Annotation = Annotation.(string)
	} else {
		aaaDomainAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		aaaDomainAttr.NameAlias = NameAlias.(string)
	}
	aaaDomain := models.NewSecurityDomain(fmt.Sprintf("userext/domain-%s", name), "uni", desc, aaaDomainAttr)

	aaaDomain.Status = "modified"

	err := aciClient.Save(aaaDomain)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(aaaDomain.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSecurityDomainRead(d, m)

}

func resourceAciSecurityDomainRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	aaaDomain, err := getRemoteSecurityDomain(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setSecurityDomainAttributes(aaaDomain, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSecurityDomainDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaDomain")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

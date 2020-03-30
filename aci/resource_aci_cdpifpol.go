package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciCDPInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciCDPInterfacePolicyCreate,
		Update: resourceAciCDPInterfacePolicyUpdate,
		Read:   resourceAciCDPInterfacePolicyRead,
		Delete: resourceAciCDPInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCDPInterfacePolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
func getRemoteCDPInterfacePolicy(client *client.Client, dn string) (*models.CDPInterfacePolicy, error) {
	cdpIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cdpIfPol := models.CDPInterfacePolicyFromContainer(cdpIfPolCont)

	if cdpIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("CDPInterfacePolicy %s not found", cdpIfPol.DistinguishedName)
	}

	return cdpIfPol, nil
}

func setCDPInterfacePolicyAttributes(cdpIfPol *models.CDPInterfacePolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(cdpIfPol.DistinguishedName)
	d.Set("description", cdpIfPol.Description)
	cdpIfPolMap, _ := cdpIfPol.ToMap()

	d.Set("name", cdpIfPolMap["name"])

	d.Set("admin_st", cdpIfPolMap["adminSt"])
	d.Set("annotation", cdpIfPolMap["annotation"])
	d.Set("name_alias", cdpIfPolMap["nameAlias"])
	return d
}

func resourceAciCDPInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cdpIfPol, err := getRemoteCDPInterfacePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setCDPInterfacePolicyAttributes(cdpIfPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCDPInterfacePolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CDPInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	cdpIfPolAttr := models.CDPInterfacePolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		cdpIfPolAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cdpIfPolAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cdpIfPolAttr.NameAlias = NameAlias.(string)
	}
	cdpIfPol := models.NewCDPInterfacePolicy(fmt.Sprintf("infra/cdpIfP-%s", name), "uni", desc, cdpIfPolAttr)

	err := aciClient.Save(cdpIfPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(cdpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCDPInterfacePolicyRead(d, m)
}

func resourceAciCDPInterfacePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CDPInterfacePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	cdpIfPolAttr := models.CDPInterfacePolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		cdpIfPolAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cdpIfPolAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cdpIfPolAttr.NameAlias = NameAlias.(string)
	}
	cdpIfPol := models.NewCDPInterfacePolicy(fmt.Sprintf("infra/cdpIfP-%s", name), "uni", desc, cdpIfPolAttr)

	cdpIfPol.Status = "modified"

	err := aciClient.Save(cdpIfPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(cdpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCDPInterfacePolicyRead(d, m)

}

func resourceAciCDPInterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cdpIfPol, err := getRemoteCDPInterfacePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setCDPInterfacePolicyAttributes(cdpIfPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCDPInterfacePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cdpIfPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

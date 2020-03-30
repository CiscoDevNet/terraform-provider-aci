package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciTabooContract() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciTabooContractCreate,
		Update: resourceAciTabooContractUpdate,
		Read:   resourceAciTabooContractRead,
		Delete: resourceAciTabooContractDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTabooContractImport,
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
func getRemoteTabooContract(client *client.Client, dn string) (*models.TabooContract, error) {
	vzTabooCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzTaboo := models.TabooContractFromContainer(vzTabooCont)

	if vzTaboo.DistinguishedName == "" {
		return nil, fmt.Errorf("TabooContract %s not found", vzTaboo.DistinguishedName)
	}

	return vzTaboo, nil
}

func setTabooContractAttributes(vzTaboo *models.TabooContract, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(vzTaboo.DistinguishedName)
	d.Set("description", vzTaboo.Description)
	d.Set("tenant_dn", GetParentDn(vzTaboo.DistinguishedName))
	vzTabooMap, _ := vzTaboo.ToMap()

	d.Set("name", vzTabooMap["name"])

	d.Set("annotation", vzTabooMap["annotation"])
	d.Set("name_alias", vzTabooMap["nameAlias"])
	return d
}

func resourceAciTabooContractImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzTaboo, err := getRemoteTabooContract(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setTabooContractAttributes(vzTaboo, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTabooContractCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] TabooContract: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzTabooAttr := models.TabooContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzTabooAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzTabooAttr.NameAlias = NameAlias.(string)
	}
	vzTaboo := models.NewTabooContract(fmt.Sprintf("taboo-%s", name), TenantDn, desc, vzTabooAttr)

	err := aciClient.Save(vzTaboo)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(vzTaboo.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciTabooContractRead(d, m)
}

func resourceAciTabooContractUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] TabooContract: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzTabooAttr := models.TabooContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzTabooAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzTabooAttr.NameAlias = NameAlias.(string)
	}
	vzTaboo := models.NewTabooContract(fmt.Sprintf("taboo-%s", name), TenantDn, desc, vzTabooAttr)

	vzTaboo.Status = "modified"

	err := aciClient.Save(vzTaboo)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(vzTaboo.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciTabooContractRead(d, m)

}

func resourceAciTabooContractRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzTaboo, err := getRemoteTabooContract(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setTabooContractAttributes(vzTaboo, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciTabooContractDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzTaboo")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

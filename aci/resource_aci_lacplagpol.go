package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciLACPPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciLACPPolicyCreate,
		Update: resourceAciLACPPolicyUpdate,
		Read:   resourceAciLACPPolicyRead,
		Delete: resourceAciLACPPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLACPPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_links": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"min_links": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
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
func getRemoteLACPPolicy(client *client.Client, dn string) (*models.LACPPolicy, error) {
	lacpLagPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	lacpLagPol := models.LACPPolicyFromContainer(lacpLagPolCont)

	if lacpLagPol.DistinguishedName == "" {
		return nil, fmt.Errorf("LACPPolicy %s not found", lacpLagPol.DistinguishedName)
	}

	return lacpLagPol, nil
}

func setLACPPolicyAttributes(lacpLagPol *models.LACPPolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(lacpLagPol.DistinguishedName)
	d.Set("description", lacpLagPol.Description)
	lacpLagPolMap, _ := lacpLagPol.ToMap()

	d.Set("name", lacpLagPolMap["name"])

	d.Set("annotation", lacpLagPolMap["annotation"])
	d.Set("ctrl", lacpLagPolMap["ctrl"])
	d.Set("max_links", lacpLagPolMap["maxLinks"])
	d.Set("min_links", lacpLagPolMap["minLinks"])
	d.Set("mode", lacpLagPolMap["mode"])
	d.Set("name_alias", lacpLagPolMap["nameAlias"])
	return d
}

func resourceAciLACPPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	lacpLagPol, err := getRemoteLACPPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setLACPPolicyAttributes(lacpLagPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLACPPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LACPPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	lacpLagPolAttr := models.LACPPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		lacpLagPolAttr.Annotation = Annotation.(string)
	} else {
		lacpLagPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		lacpLagPolAttr.Ctrl = Ctrl.(string)
	}
	if MaxLinks, ok := d.GetOk("max_links"); ok {
		lacpLagPolAttr.MaxLinks = MaxLinks.(string)
	}
	if MinLinks, ok := d.GetOk("min_links"); ok {
		lacpLagPolAttr.MinLinks = MinLinks.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		lacpLagPolAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		lacpLagPolAttr.NameAlias = NameAlias.(string)
	}
	lacpLagPol := models.NewLACPPolicy(fmt.Sprintf("infra/lacplagp-%s", name), "uni", desc, lacpLagPolAttr)

	err := aciClient.Save(lacpLagPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(lacpLagPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLACPPolicyRead(d, m)
}

func resourceAciLACPPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LACPPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	lacpLagPolAttr := models.LACPPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		lacpLagPolAttr.Annotation = Annotation.(string)
	} else {
		lacpLagPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		lacpLagPolAttr.Ctrl = Ctrl.(string)
	}
	if MaxLinks, ok := d.GetOk("max_links"); ok {
		lacpLagPolAttr.MaxLinks = MaxLinks.(string)
	}
	if MinLinks, ok := d.GetOk("min_links"); ok {
		lacpLagPolAttr.MinLinks = MinLinks.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		lacpLagPolAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		lacpLagPolAttr.NameAlias = NameAlias.(string)
	}
	lacpLagPol := models.NewLACPPolicy(fmt.Sprintf("infra/lacplagp-%s", name), "uni", desc, lacpLagPolAttr)

	lacpLagPol.Status = "modified"

	err := aciClient.Save(lacpLagPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(lacpLagPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLACPPolicyRead(d, m)

}

func resourceAciLACPPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	lacpLagPol, err := getRemoteLACPPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setLACPPolicyAttributes(lacpLagPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLACPPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "lacpLagPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

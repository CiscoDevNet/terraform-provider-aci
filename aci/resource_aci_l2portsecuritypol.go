package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAciPortSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciPortSecurityPolicyCreate,
		Update: resourceAciPortSecurityPolicyUpdate,
		Read:   resourceAciPortSecurityPolicyRead,
		Delete: resourceAciPortSecurityPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPortSecurityPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"maximum": &schema.Schema{
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

			"timeout": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"violation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemotePortSecurityPolicy(client *client.Client, dn string) (*models.PortSecurityPolicy, error) {
	l2PortSecurityPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l2PortSecurityPol := models.PortSecurityPolicyFromContainer(l2PortSecurityPolCont)

	if l2PortSecurityPol.DistinguishedName == "" {
		return nil, fmt.Errorf("PortSecurityPolicy %s not found", l2PortSecurityPol.DistinguishedName)
	}

	return l2PortSecurityPol, nil
}

func setPortSecurityPolicyAttributes(l2PortSecurityPol *models.PortSecurityPolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(l2PortSecurityPol.DistinguishedName)
	d.Set("description", l2PortSecurityPol.Description)
	l2PortSecurityPolMap, _ := l2PortSecurityPol.ToMap()

	d.Set("name", l2PortSecurityPolMap["name"])

	d.Set("annotation", l2PortSecurityPolMap["annotation"])
	d.Set("maximum", l2PortSecurityPolMap["maximum"])
	d.Set("mode", l2PortSecurityPolMap["mode"])
	d.Set("name_alias", l2PortSecurityPolMap["nameAlias"])
	d.Set("timeout", l2PortSecurityPolMap["timeout"])
	d.Set("violation", l2PortSecurityPolMap["violation"])
	return d
}

func resourceAciPortSecurityPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l2PortSecurityPol, err := getRemotePortSecurityPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setPortSecurityPolicyAttributes(l2PortSecurityPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPortSecurityPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] PortSecurityPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l2PortSecurityPolAttr := models.PortSecurityPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2PortSecurityPolAttr.Annotation = Annotation.(string)
	}
	if Maximum, ok := d.GetOk("maximum"); ok {
		l2PortSecurityPolAttr.Maximum = Maximum.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		l2PortSecurityPolAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2PortSecurityPolAttr.NameAlias = NameAlias.(string)
	}
	if Timeout, ok := d.GetOk("timeout"); ok {
		l2PortSecurityPolAttr.Timeout = Timeout.(string)
	}
	if Violation, ok := d.GetOk("violation"); ok {
		l2PortSecurityPolAttr.Violation = Violation.(string)
	}
	l2PortSecurityPol := models.NewPortSecurityPolicy(fmt.Sprintf("infra/portsecurityP-%s", name), "uni", desc, l2PortSecurityPolAttr)

	err := aciClient.Save(l2PortSecurityPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(l2PortSecurityPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciPortSecurityPolicyRead(d, m)
}

func resourceAciPortSecurityPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] PortSecurityPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l2PortSecurityPolAttr := models.PortSecurityPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2PortSecurityPolAttr.Annotation = Annotation.(string)
	}
	if Maximum, ok := d.GetOk("maximum"); ok {
		l2PortSecurityPolAttr.Maximum = Maximum.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		l2PortSecurityPolAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2PortSecurityPolAttr.NameAlias = NameAlias.(string)
	}
	if Timeout, ok := d.GetOk("timeout"); ok {
		l2PortSecurityPolAttr.Timeout = Timeout.(string)
	}
	if Violation, ok := d.GetOk("violation"); ok {
		l2PortSecurityPolAttr.Violation = Violation.(string)
	}
	l2PortSecurityPol := models.NewPortSecurityPolicy(fmt.Sprintf("infra/portsecurityP-%s", name), "uni", desc, l2PortSecurityPolAttr)

	l2PortSecurityPol.Status = "modified"

	err := aciClient.Save(l2PortSecurityPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(l2PortSecurityPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciPortSecurityPolicyRead(d, m)

}

func resourceAciPortSecurityPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l2PortSecurityPol, err := getRemotePortSecurityPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setPortSecurityPolicyAttributes(l2PortSecurityPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciPortSecurityPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l2PortSecurityPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

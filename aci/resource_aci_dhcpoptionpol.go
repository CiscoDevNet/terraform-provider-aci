package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciDHCPOptionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciDHCPOptionPolicyCreate,
		Update: resourceAciDHCPOptionPolicyUpdate,
		Read:   resourceAciDHCPOptionPolicyRead,
		Delete: resourceAciDHCPOptionPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciDHCPOptionPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

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
func getRemoteDHCPOptionPolicy(client *client.Client, dn string) (*models.DHCPOptionPolicy, error) {
	dhcpOptionPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpOptionPol := models.DHCPOptionPolicyFromContainer(dhcpOptionPolCont)

	if dhcpOptionPol.DistinguishedName == "" {
		return nil, fmt.Errorf("DHCPOptionPolicy %s not found", dhcpOptionPol.DistinguishedName)
	}

	return dhcpOptionPol, nil
}

func setDHCPOptionPolicyAttributes(dhcpOptionPol *models.DHCPOptionPolicy, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(dhcpOptionPol.DistinguishedName)
	d.Set("description", dhcpOptionPol.Description)

	if dn != dhcpOptionPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}

	dhcpOptionPolMap, _ := dhcpOptionPol.ToMap()

	d.Set("name", dhcpOptionPolMap["name"])

	d.Set("annotation", dhcpOptionPolMap["annotation"])
	d.Set("name_alias", dhcpOptionPolMap["nameAlias"])
	return d
}

func resourceAciDHCPOptionPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	dhcpOptionPol, err := getRemoteDHCPOptionPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setDHCPOptionPolicyAttributes(dhcpOptionPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciDHCPOptionPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] DHCPOptionPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	dhcpOptionPolAttr := models.DHCPOptionPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpOptionPolAttr.Annotation = Annotation.(string)
	} else {
		dhcpOptionPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpOptionPolAttr.NameAlias = NameAlias.(string)
	}
	dhcpOptionPol := models.NewDHCPOptionPolicy(fmt.Sprintf("dhcpoptpol-%s", name), TenantDn, desc, dhcpOptionPolAttr)

	err := aciClient.Save(dhcpOptionPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(dhcpOptionPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciDHCPOptionPolicyRead(d, m)
}

func resourceAciDHCPOptionPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] DHCPOptionPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	dhcpOptionPolAttr := models.DHCPOptionPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpOptionPolAttr.Annotation = Annotation.(string)
	} else {
		dhcpOptionPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpOptionPolAttr.NameAlias = NameAlias.(string)
	}
	dhcpOptionPol := models.NewDHCPOptionPolicy(fmt.Sprintf("dhcpoptpol-%s", name), TenantDn, desc, dhcpOptionPolAttr)

	dhcpOptionPol.Status = "modified"

	err := aciClient.Save(dhcpOptionPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(dhcpOptionPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciDHCPOptionPolicyRead(d, m)

}

func resourceAciDHCPOptionPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	dhcpOptionPol, err := getRemoteDHCPOptionPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setDHCPOptionPolicyAttributes(dhcpOptionPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciDHCPOptionPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "dhcpOptionPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

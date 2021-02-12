package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciDHCPOption() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciDHCPOptionCreate,
		Update: resourceAciDHCPOptionUpdate,
		Read:   resourceAciDHCPOptionRead,
		Delete: resourceAciDHCPOptionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciDHCPOptionImport,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"dhcp_option_policy_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "orchestrator:terraform",
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dhcp_option_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}
func getRemoteDHCPOption(client *client.Client, dn string) (*models.DHCPOption, error) {
	dhcpOptionCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpOption := models.DHCPOptionFromContainer(dhcpOptionCont)

	if dhcpOption.DistinguishedName == "" {
		return nil, fmt.Errorf("DHCPOption %s not found", dhcpOption.DistinguishedName)
	}

	return dhcpOption, nil
}

func setDHCPOptionAttributes(dhcpOption *models.DHCPOption, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(dhcpOption.DistinguishedName)
	//d.Set("description", dhcpOption.Description)

	if dn != dhcpOption.DistinguishedName {
		d.Set("dhcp_option_policy_dn", "")
	}
	dhcpOptionMap, _ := dhcpOption.ToMap()

	d.Set("name", dhcpOptionMap["name"])

	d.Set("annotation", dhcpOptionMap["annotation"])
	d.Set("data", dhcpOptionMap["data"])
	d.Set("dhcp_option_id", dhcpOptionMap["id"])
	d.Set("name_alias", dhcpOptionMap["nameAlias"])
	return d
}

func resourceAciDHCPOptionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	dhcpOption, err := getRemoteDHCPOption(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setDHCPOptionAttributes(dhcpOption, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciDHCPOptionCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] DHCPOption: Beginning Creation")
	aciClient := m.(*client.Client)
	//desc := d.Get("description").(string)

	name := d.Get("name").(string)

	DHCPOptionPolicyDn := d.Get("dhcp_option_policy_dn").(string)

	dhcpOptionAttr := models.DHCPOptionAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpOptionAttr.Annotation = Annotation.(string)
	} else {
		dhcpOptionAttr.Annotation = "{}"
	}
	if Data, ok := d.GetOk("data"); ok {
		dhcpOptionAttr.Data = Data.(string)
	}
	if DHCPOption_id, ok := d.GetOk("dhcp_option_id"); ok {
		dhcpOptionAttr.DHCPOption_id = DHCPOption_id.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpOptionAttr.NameAlias = NameAlias.(string)
	}
	dhcpOption := models.NewDHCPOption(fmt.Sprintf("opt-%s", name), DHCPOptionPolicyDn, dhcpOptionAttr)

	err := aciClient.Save(dhcpOption)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(dhcpOption.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciDHCPOptionRead(d, m)
}

func resourceAciDHCPOptionUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] DHCPOption: Beginning Update")

	aciClient := m.(*client.Client)
	//desc := d.Get("description").(string)

	name := d.Get("name").(string)

	DHCPOptionPolicyDn := d.Get("dhcp_option_policy_dn").(string)

	dhcpOptionAttr := models.DHCPOptionAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpOptionAttr.Annotation = Annotation.(string)
	} else {
		dhcpOptionAttr.Annotation = "{}"
	}
	if Data, ok := d.GetOk("data"); ok {
		dhcpOptionAttr.Data = Data.(string)
	}
	if DHCPOption_id, ok := d.GetOk("dhcp_option_id"); ok {
		dhcpOptionAttr.DHCPOption_id = DHCPOption_id.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpOptionAttr.NameAlias = NameAlias.(string)
	}
	dhcpOption := models.NewDHCPOption(fmt.Sprintf("opt-%s", name), DHCPOptionPolicyDn, dhcpOptionAttr)

	dhcpOption.Status = "modified"

	err := aciClient.Save(dhcpOption)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(dhcpOption.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciDHCPOptionRead(d, m)

}

func resourceAciDHCPOptionRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	dhcpOption, err := getRemoteDHCPOption(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setDHCPOptionAttributes(dhcpOption, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciDHCPOptionDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "dhcpOption")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

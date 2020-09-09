package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciL2InterfacePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciL2InterfacePolicyCreate,
		Update: resourceAciL2InterfacePolicyUpdate,
		Read:   resourceAciL2InterfacePolicyRead,
		Delete: resourceAciL2InterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL2InterfacePolicyImport,
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

			"qinq": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"edgePort",
					"corePort",
					"doubleQtagPort",
				}, false),
			},

			"vepa": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"vlan_scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"global",
					"portlocal",
				}, false),
			},
		}),
	}
}
func getRemoteL2InterfacePolicy(client *client.Client, dn string) (*models.L2InterfacePolicy, error) {
	l2IfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l2IfPol := models.L2InterfacePolicyFromContainer(l2IfPolCont)

	if l2IfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("L2InterfacePolicy %s not found", l2IfPol.DistinguishedName)
	}

	return l2IfPol, nil
}

func setL2InterfacePolicyAttributes(l2IfPol *models.L2InterfacePolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(l2IfPol.DistinguishedName)
	d.Set("description", l2IfPol.Description)
	l2IfPolMap, _ := l2IfPol.ToMap()

	d.Set("name", l2IfPolMap["name"])

	d.Set("annotation", l2IfPolMap["annotation"])
	d.Set("name_alias", l2IfPolMap["nameAlias"])
	d.Set("qinq", l2IfPolMap["qinq"])
	d.Set("vepa", l2IfPolMap["vepa"])
	d.Set("vlan_scope", l2IfPolMap["vlanScope"])
	return d
}

func resourceAciL2InterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l2IfPol, err := getRemoteL2InterfacePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setL2InterfacePolicyAttributes(l2IfPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL2InterfacePolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L2InterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l2IfPolAttr := models.L2InterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2IfPolAttr.Annotation = Annotation.(string)
	} else {
		l2IfPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2IfPolAttr.NameAlias = NameAlias.(string)
	}
	if Qinq, ok := d.GetOk("qinq"); ok {
		l2IfPolAttr.Qinq = Qinq.(string)
	}
	if Vepa, ok := d.GetOk("vepa"); ok {
		l2IfPolAttr.Vepa = Vepa.(string)
	}
	if VlanScope, ok := d.GetOk("vlan_scope"); ok {
		l2IfPolAttr.VlanScope = VlanScope.(string)
	}
	l2IfPol := models.NewL2InterfacePolicy(fmt.Sprintf("infra/l2IfP-%s", name), "uni", desc, l2IfPolAttr)

	err := aciClient.Save(l2IfPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(l2IfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL2InterfacePolicyRead(d, m)
}

func resourceAciL2InterfacePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L2InterfacePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l2IfPolAttr := models.L2InterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2IfPolAttr.Annotation = Annotation.(string)
	} else {
		l2IfPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2IfPolAttr.NameAlias = NameAlias.(string)
	}
	if Qinq, ok := d.GetOk("qinq"); ok {
		l2IfPolAttr.Qinq = Qinq.(string)
	}
	if Vepa, ok := d.GetOk("vepa"); ok {
		l2IfPolAttr.Vepa = Vepa.(string)
	}
	if VlanScope, ok := d.GetOk("vlan_scope"); ok {
		l2IfPolAttr.VlanScope = VlanScope.(string)
	}
	l2IfPol := models.NewL2InterfacePolicy(fmt.Sprintf("infra/l2IfP-%s", name), "uni", desc, l2IfPolAttr)

	l2IfPol.Status = "modified"

	err := aciClient.Save(l2IfPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(l2IfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL2InterfacePolicyRead(d, m)

}

func resourceAciL2InterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l2IfPol, err := getRemoteL2InterfacePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setL2InterfacePolicyAttributes(l2IfPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL2InterfacePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l2IfPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

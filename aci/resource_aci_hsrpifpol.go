package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciHSRPInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciHSRPInterfacePolicyCreate,
		Update: resourceAciHSRPInterfacePolicyUpdate,
		Read:   resourceAciHSRPInterfacePolicyRead,
		Delete: resourceAciHSRPInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciHSRPInterfacePolicyImport,
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

			"ctrl": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressBitMaskDiffFunc(),
				ValidateFunc: schema.SchemaValidateFunc(validateCommaSeparatedStringInSlice([]string{
					"bfd",
					"bia",
					"",
				}, false, "")),
			},

			"delay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"reload_delay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func getRemoteHSRPInterfacePolicy(client *client.Client, dn string) (*models.HSRPInterfacePolicy, error) {
	hsrpIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	hsrpIfPol := models.HSRPInterfacePolicyFromContainer(hsrpIfPolCont)

	if hsrpIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("HSRPInterfacePolicy %s not found", hsrpIfPol.DistinguishedName)
	}

	return hsrpIfPol, nil
}

func setHSRPInterfacePolicyAttributes(hsrpIfPol *models.HSRPInterfacePolicy, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(hsrpIfPol.DistinguishedName)
	d.Set("description", hsrpIfPol.Description)
	if dn != hsrpIfPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	hsrpIfPolMap, _ := hsrpIfPol.ToMap()

	d.Set("name", hsrpIfPolMap["name"])

	d.Set("annotation", hsrpIfPolMap["annotation"])
	d.Set("ctrl", hsrpIfPolMap["ctrl"])
	d.Set("delay", hsrpIfPolMap["delay"])
	d.Set("name_alias", hsrpIfPolMap["nameAlias"])
	d.Set("reload_delay", hsrpIfPolMap["reloadDelay"])

	return d
}

func resourceAciHSRPInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	hsrpIfPol, err := getRemoteHSRPInterfacePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setHSRPInterfacePolicyAttributes(hsrpIfPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciHSRPInterfacePolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] HSRPInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	hsrpIfPolAttr := models.HSRPInterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpIfPolAttr.Annotation = Annotation.(string)
	} else {
		hsrpIfPolAttr.Annotation = "{}"
	}
	if ctrl, ok := d.GetOk("ctrl"); ok {
		hsrpIfPolAttr.Ctrl = ctrl.(string)
	} else {
		hsrpIfPolAttr.Ctrl = "{}"
	}
	if Delay, ok := d.GetOk("delay"); ok {
		hsrpIfPolAttr.Delay = Delay.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpIfPolAttr.NameAlias = NameAlias.(string)
	}
	if ReloadDelay, ok := d.GetOk("reload_delay"); ok {
		hsrpIfPolAttr.ReloadDelay = ReloadDelay.(string)
	}
	hsrpIfPol := models.NewHSRPInterfacePolicy(fmt.Sprintf("hsrpIfPol-%s", name), TenantDn, desc, hsrpIfPolAttr)

	err := aciClient.Save(hsrpIfPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(hsrpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciHSRPInterfacePolicyRead(d, m)
}

func resourceAciHSRPInterfacePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] HSRPInterfacePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	hsrpIfPolAttr := models.HSRPInterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpIfPolAttr.Annotation = Annotation.(string)
	} else {
		hsrpIfPolAttr.Annotation = "{}"
	}
	if ctrl, ok := d.GetOk("ctrl"); ok {
		hsrpIfPolAttr.Ctrl = ctrl.(string)
	} else {
		hsrpIfPolAttr.Ctrl = "{}"
	}
	if Delay, ok := d.GetOk("delay"); ok {
		hsrpIfPolAttr.Delay = Delay.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpIfPolAttr.NameAlias = NameAlias.(string)
	}
	if ReloadDelay, ok := d.GetOk("reload_delay"); ok {
		hsrpIfPolAttr.ReloadDelay = ReloadDelay.(string)
	}
	hsrpIfPol := models.NewHSRPInterfacePolicy(fmt.Sprintf("hsrpIfPol-%s", name), TenantDn, desc, hsrpIfPolAttr)

	hsrpIfPol.Status = "modified"

	err := aciClient.Save(hsrpIfPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(hsrpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciHSRPInterfacePolicyRead(d, m)

}

func resourceAciHSRPInterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	hsrpIfPol, err := getRemoteHSRPInterfacePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setHSRPInterfacePolicyAttributes(hsrpIfPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciHSRPInterfacePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "hsrpIfPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

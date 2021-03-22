package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciHSRPGroupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciHSRPGroupPolicyCreate,
		Update: resourceAciHSRPGroupPolicyUpdate,
		Read:   resourceAciHSRPGroupPolicyRead,
		Delete: resourceAciHSRPGroupPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciHSRPGroupPolicyImport,
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"preempt", "0",
				}, false),
			},

			"hello_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"hold_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"preempt_delay_min": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"preempt_delay_reload": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"preempt_delay_sync": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"timeout": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"hsrp_group_policy_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"md5",
					"simple",
				}, false),
			},
		}),
	}
}
func getRemoteHSRPGroupPolicy(client *client.Client, dn string) (*models.HSRPGroupPolicy, error) {
	hsrpGroupPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	hsrpGroupPol := models.HSRPGroupPolicyFromContainer(hsrpGroupPolCont)

	if hsrpGroupPol.DistinguishedName == "" {
		return nil, fmt.Errorf("HSRPGroupPolicy %s not found", hsrpGroupPol.DistinguishedName)
	}

	return hsrpGroupPol, nil
}

func setHSRPGroupPolicyAttributes(hsrpGroupPol *models.HSRPGroupPolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(hsrpGroupPol.DistinguishedName)
	d.Set("description", hsrpGroupPol.Description)
	dn := d.Id()
	if dn != hsrpGroupPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	hsrpGroupPolMap, _ := hsrpGroupPol.ToMap()

	d.Set("name", hsrpGroupPolMap["name"])

	d.Set("annotation", hsrpGroupPolMap["annotation"])
	if hsrpGroupPolMap["ctrl"] == "" {
		d.Set("ctrl", "0")
	} else {
		d.Set("ctrl", hsrpGroupPolMap["ctrl"])
	}
	d.Set("hello_intvl", hsrpGroupPolMap["helloIntvl"])
	d.Set("hold_intvl", hsrpGroupPolMap["holdIntvl"])
	d.Set("key", hsrpGroupPolMap["key"])
	d.Set("name_alias", hsrpGroupPolMap["nameAlias"])
	d.Set("preempt_delay_min", hsrpGroupPolMap["preemptDelayMin"])
	d.Set("preempt_delay_reload", hsrpGroupPolMap["preemptDelayReload"])
	d.Set("preempt_delay_sync", hsrpGroupPolMap["preemptDelaySync"])
	d.Set("prio", hsrpGroupPolMap["prio"])
	d.Set("timeout", hsrpGroupPolMap["timeout"])
	d.Set("hsrp_group_policy_type", hsrpGroupPolMap["type"])
	return d
}

func resourceAciHSRPGroupPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	hsrpGroupPol, err := getRemoteHSRPGroupPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setHSRPGroupPolicyAttributes(hsrpGroupPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciHSRPGroupPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] HSRPGroupPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	hsrpGroupPolAttr := models.HSRPGroupPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpGroupPolAttr.Annotation = Annotation.(string)
	} else {
		hsrpGroupPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		hsrpGroupPolAttr.Ctrl = Ctrl.(string)
	}
	if HelloIntvl, ok := d.GetOk("hello_intvl"); ok {
		hsrpGroupPolAttr.HelloIntvl = HelloIntvl.(string)
	}
	if HoldIntvl, ok := d.GetOk("hold_intvl"); ok {
		hsrpGroupPolAttr.HoldIntvl = HoldIntvl.(string)
	}
	if Key, ok := d.GetOk("key"); ok {
		hsrpGroupPolAttr.Key = Key.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpGroupPolAttr.NameAlias = NameAlias.(string)
	}
	if PreemptDelayMin, ok := d.GetOk("preempt_delay_min"); ok {
		hsrpGroupPolAttr.PreemptDelayMin = PreemptDelayMin.(string)
	}
	if PreemptDelayReload, ok := d.GetOk("preempt_delay_reload"); ok {
		hsrpGroupPolAttr.PreemptDelayReload = PreemptDelayReload.(string)
	}
	if PreemptDelaySync, ok := d.GetOk("preempt_delay_sync"); ok {
		hsrpGroupPolAttr.PreemptDelaySync = PreemptDelaySync.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		hsrpGroupPolAttr.Prio = Prio.(string)
	}
	if Timeout, ok := d.GetOk("timeout"); ok {
		hsrpGroupPolAttr.Timeout = Timeout.(string)
	}
	if HSRPGroupPolicy_type, ok := d.GetOk("hsrp_group_policy_type"); ok {
		hsrpGroupPolAttr.HSRPGroupPolicy_type = HSRPGroupPolicy_type.(string)
	}
	hsrpGroupPol := models.NewHSRPGroupPolicy(fmt.Sprintf("hsrpGroupPol-%s", name), TenantDn, desc, hsrpGroupPolAttr)

	err := aciClient.Save(hsrpGroupPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(hsrpGroupPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciHSRPGroupPolicyRead(d, m)
}

func resourceAciHSRPGroupPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] HSRPGroupPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	hsrpGroupPolAttr := models.HSRPGroupPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpGroupPolAttr.Annotation = Annotation.(string)
	} else {
		hsrpGroupPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		hsrpGroupPolAttr.Ctrl = Ctrl.(string)
	}
	if HelloIntvl, ok := d.GetOk("hello_intvl"); ok {
		hsrpGroupPolAttr.HelloIntvl = HelloIntvl.(string)
	}
	if HoldIntvl, ok := d.GetOk("hold_intvl"); ok {
		hsrpGroupPolAttr.HoldIntvl = HoldIntvl.(string)
	}
	if Key, ok := d.GetOk("key"); ok {
		hsrpGroupPolAttr.Key = Key.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpGroupPolAttr.NameAlias = NameAlias.(string)
	}
	if PreemptDelayMin, ok := d.GetOk("preempt_delay_min"); ok {
		hsrpGroupPolAttr.PreemptDelayMin = PreemptDelayMin.(string)
	}
	if PreemptDelayReload, ok := d.GetOk("preempt_delay_reload"); ok {
		hsrpGroupPolAttr.PreemptDelayReload = PreemptDelayReload.(string)
	}
	if PreemptDelaySync, ok := d.GetOk("preempt_delay_sync"); ok {
		hsrpGroupPolAttr.PreemptDelaySync = PreemptDelaySync.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		hsrpGroupPolAttr.Prio = Prio.(string)
	}
	if Timeout, ok := d.GetOk("timeout"); ok {
		hsrpGroupPolAttr.Timeout = Timeout.(string)
	}
	if HSRPGroupPolicy_type, ok := d.GetOk("hsrp_group_policy_type"); ok {
		hsrpGroupPolAttr.HSRPGroupPolicy_type = HSRPGroupPolicy_type.(string)
	}
	hsrpGroupPol := models.NewHSRPGroupPolicy(fmt.Sprintf("hsrpGroupPol-%s", name), TenantDn, desc, hsrpGroupPolAttr)

	hsrpGroupPol.Status = "modified"

	err := aciClient.Save(hsrpGroupPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(hsrpGroupPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciHSRPGroupPolicyRead(d, m)

}

func resourceAciHSRPGroupPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	hsrpGroupPol, err := getRemoteHSRPGroupPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setHSRPGroupPolicyAttributes(hsrpGroupPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciHSRPGroupPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "hsrpGroupPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

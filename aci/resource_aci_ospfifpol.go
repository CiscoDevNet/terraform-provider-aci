package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAciOSPFInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciOSPFInterfacePolicyCreate,
		Update: resourceAciOSPFInterfacePolicyUpdate,
		Read:   resourceAciOSPFInterfacePolicyRead,
		Delete: resourceAciOSPFInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciOSPFInterfacePolicyImport,
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

			"cost": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dead_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"hello_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"nw_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pfx_suppress": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rexmit_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"xmit_delay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteOSPFInterfacePolicy(client *client.Client, dn string) (*models.OSPFInterfacePolicy, error) {
	ospfIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	ospfIfPol := models.OSPFInterfacePolicyFromContainer(ospfIfPolCont)

	if ospfIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("OSPFInterfacePolicy %s not found", ospfIfPol.DistinguishedName)
	}

	return ospfIfPol, nil
}

func setOSPFInterfacePolicyAttributes(ospfIfPol *models.OSPFInterfacePolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(ospfIfPol.DistinguishedName)
	d.Set("description", ospfIfPol.Description)
	d.Set("tenant_dn", GetParentDn(ospfIfPol.DistinguishedName))
	ospfIfPolMap, _ := ospfIfPol.ToMap()

	d.Set("name", ospfIfPolMap["name"])

	d.Set("annotation", ospfIfPolMap["annotation"])
	d.Set("cost", ospfIfPolMap["cost"])
	d.Set("ctrl", ospfIfPolMap["ctrl"])
	d.Set("dead_intvl", ospfIfPolMap["deadIntvl"])
	d.Set("hello_intvl", ospfIfPolMap["helloIntvl"])
	d.Set("name_alias", ospfIfPolMap["nameAlias"])
	d.Set("nw_t", ospfIfPolMap["nwT"])
	d.Set("pfx_suppress", ospfIfPolMap["pfxSuppress"])
	d.Set("prio", ospfIfPolMap["prio"])
	d.Set("rexmit_intvl", ospfIfPolMap["rexmitIntvl"])
	d.Set("xmit_delay", ospfIfPolMap["xmitDelay"])
	return d
}

func resourceAciOSPFInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	ospfIfPol, err := getRemoteOSPFInterfacePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setOSPFInterfacePolicyAttributes(ospfIfPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciOSPFInterfacePolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] OSPFInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	ospfIfPolAttr := models.OSPFInterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ospfIfPolAttr.Annotation = Annotation.(string)
	}
	if Cost, ok := d.GetOk("cost"); ok {
		ospfIfPolAttr.Cost = Cost.(string)
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		ospfIfPolAttr.Ctrl = Ctrl.(string)
	}
	if DeadIntvl, ok := d.GetOk("dead_intvl"); ok {
		ospfIfPolAttr.DeadIntvl = DeadIntvl.(string)
	}
	if HelloIntvl, ok := d.GetOk("hello_intvl"); ok {
		ospfIfPolAttr.HelloIntvl = HelloIntvl.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ospfIfPolAttr.NameAlias = NameAlias.(string)
	}
	if NwT, ok := d.GetOk("nw_t"); ok {
		ospfIfPolAttr.NwT = NwT.(string)
	}
	if PfxSuppress, ok := d.GetOk("pfx_suppress"); ok {
		ospfIfPolAttr.PfxSuppress = PfxSuppress.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		ospfIfPolAttr.Prio = Prio.(string)
	}
	if RexmitIntvl, ok := d.GetOk("rexmit_intvl"); ok {
		ospfIfPolAttr.RexmitIntvl = RexmitIntvl.(string)
	}
	if XmitDelay, ok := d.GetOk("xmit_delay"); ok {
		ospfIfPolAttr.XmitDelay = XmitDelay.(string)
	}
	ospfIfPol := models.NewOSPFInterfacePolicy(fmt.Sprintf("ospfIfPol-%s", name), TenantDn, desc, ospfIfPolAttr)

	err := aciClient.Save(ospfIfPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(ospfIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciOSPFInterfacePolicyRead(d, m)
}

func resourceAciOSPFInterfacePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] OSPFInterfacePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	ospfIfPolAttr := models.OSPFInterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ospfIfPolAttr.Annotation = Annotation.(string)
	}
	if Cost, ok := d.GetOk("cost"); ok {
		ospfIfPolAttr.Cost = Cost.(string)
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		ospfIfPolAttr.Ctrl = Ctrl.(string)
	}
	if DeadIntvl, ok := d.GetOk("dead_intvl"); ok {
		ospfIfPolAttr.DeadIntvl = DeadIntvl.(string)
	}
	if HelloIntvl, ok := d.GetOk("hello_intvl"); ok {
		ospfIfPolAttr.HelloIntvl = HelloIntvl.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ospfIfPolAttr.NameAlias = NameAlias.(string)
	}
	if NwT, ok := d.GetOk("nw_t"); ok {
		ospfIfPolAttr.NwT = NwT.(string)
	}
	if PfxSuppress, ok := d.GetOk("pfx_suppress"); ok {
		ospfIfPolAttr.PfxSuppress = PfxSuppress.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		ospfIfPolAttr.Prio = Prio.(string)
	}
	if RexmitIntvl, ok := d.GetOk("rexmit_intvl"); ok {
		ospfIfPolAttr.RexmitIntvl = RexmitIntvl.(string)
	}
	if XmitDelay, ok := d.GetOk("xmit_delay"); ok {
		ospfIfPolAttr.XmitDelay = XmitDelay.(string)
	}
	ospfIfPol := models.NewOSPFInterfacePolicy(fmt.Sprintf("ospfIfPol-%s", name), TenantDn, desc, ospfIfPolAttr)

	ospfIfPol.Status = "modified"

	err := aciClient.Save(ospfIfPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(ospfIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciOSPFInterfacePolicyRead(d, m)

}

func resourceAciOSPFInterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	ospfIfPol, err := getRemoteOSPFInterfacePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setOSPFInterfacePolicyAttributes(ospfIfPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciOSPFInterfacePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "ospfIfPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

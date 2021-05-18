package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciBGPTimersPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciBGPTimersPolicyCreate,
		Update: resourceAciBGPTimersPolicyUpdate,
		Read:   resourceAciBGPTimersPolicyRead,
		Delete: resourceAciBGPTimersPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBGPTimersPolicyImport,
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

			"gr_ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"hold_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ka_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_as_limit": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"stale_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func getRemoteBGPTimersPolicy(client *client.Client, dn string) (*models.BGPTimersPolicy, error) {
	bgpCtxPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpCtxPol := models.BGPTimersPolicyFromContainer(bgpCtxPolCont)

	if bgpCtxPol.DistinguishedName == "" {
		return nil, fmt.Errorf("BGPTimersPolicy %s not found", bgpCtxPol.DistinguishedName)
	}

	return bgpCtxPol, nil
}

func setBGPTimersPolicyAttributes(bgpCtxPol *models.BGPTimersPolicy, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()

	d.SetId(bgpCtxPol.DistinguishedName)
	d.Set("description", bgpCtxPol.Description)
	if dn != bgpCtxPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	bgpCtxPolMap, _ := bgpCtxPol.ToMap()

	d.Set("name", bgpCtxPolMap["name"])

	d.Set("annotation", bgpCtxPolMap["annotation"])
	d.Set("gr_ctrl", bgpCtxPolMap["grCtrl"])
	d.Set("hold_intvl", bgpCtxPolMap["holdIntvl"])
	d.Set("ka_intvl", bgpCtxPolMap["kaIntvl"])
	d.Set("max_as_limit", bgpCtxPolMap["maxAsLimit"])
	d.Set("name_alias", bgpCtxPolMap["nameAlias"])
	d.Set("stale_intvl", bgpCtxPolMap["staleIntvl"])

	return d
}

func resourceAciBGPTimersPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bgpCtxPol, err := getRemoteBGPTimersPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setBGPTimersPolicyAttributes(bgpCtxPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBGPTimersPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] BGPTimersPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpCtxPolAttr := models.BGPTimersPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpCtxPolAttr.Annotation = Annotation.(string)
	} else {
		bgpCtxPolAttr.Annotation = "{}"
	}
	if GrCtrl, ok := d.GetOk("gr_ctrl"); ok {
		bgpCtxPolAttr.GrCtrl = GrCtrl.(string)
	}
	if HoldIntvl, ok := d.GetOk("hold_intvl"); ok {
		bgpCtxPolAttr.HoldIntvl = HoldIntvl.(string)
	}
	if KaIntvl, ok := d.GetOk("ka_intvl"); ok {
		bgpCtxPolAttr.KaIntvl = KaIntvl.(string)
	}
	if MaxAsLimit, ok := d.GetOk("max_as_limit"); ok {
		bgpCtxPolAttr.MaxAsLimit = MaxAsLimit.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpCtxPolAttr.NameAlias = NameAlias.(string)
	}
	if StaleIntvl, ok := d.GetOk("stale_intvl"); ok {
		bgpCtxPolAttr.StaleIntvl = StaleIntvl.(string)
	}
	bgpCtxPol := models.NewBGPTimersPolicy(fmt.Sprintf("bgpCtxP-%s", name), TenantDn, desc, bgpCtxPolAttr)

	err := aciClient.Save(bgpCtxPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(bgpCtxPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBGPTimersPolicyRead(d, m)
}

func resourceAciBGPTimersPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] BGPTimersPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpCtxPolAttr := models.BGPTimersPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpCtxPolAttr.Annotation = Annotation.(string)
	} else {
		bgpCtxPolAttr.Annotation = "{}"
	}
	if GrCtrl, ok := d.GetOk("gr_ctrl"); ok {
		bgpCtxPolAttr.GrCtrl = GrCtrl.(string)
	}
	if HoldIntvl, ok := d.GetOk("hold_intvl"); ok {
		bgpCtxPolAttr.HoldIntvl = HoldIntvl.(string)
	}
	if KaIntvl, ok := d.GetOk("ka_intvl"); ok {
		bgpCtxPolAttr.KaIntvl = KaIntvl.(string)
	}
	if MaxAsLimit, ok := d.GetOk("max_as_limit"); ok {
		bgpCtxPolAttr.MaxAsLimit = MaxAsLimit.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpCtxPolAttr.NameAlias = NameAlias.(string)
	}
	if StaleIntvl, ok := d.GetOk("stale_intvl"); ok {
		bgpCtxPolAttr.StaleIntvl = StaleIntvl.(string)
	}
	bgpCtxPol := models.NewBGPTimersPolicy(fmt.Sprintf("bgpCtxP-%s", name), TenantDn, desc, bgpCtxPolAttr)

	bgpCtxPol.Status = "modified"

	err := aciClient.Save(bgpCtxPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(bgpCtxPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBGPTimersPolicyRead(d, m)
}

func resourceAciBGPTimersPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bgpCtxPol, err := getRemoteBGPTimersPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setBGPTimersPolicyAttributes(bgpCtxPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBGPTimersPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bgpCtxPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciInterfaceFCPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciInterfaceFCPolicyCreate,
		UpdateContext: resourceAciInterfaceFCPolicyUpdate,
		ReadContext:   resourceAciInterfaceFCPolicyRead,
		DeleteContext: resourceAciInterfaceFCPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciInterfaceFCPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"automaxspeed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"2G", "4G", "8G", "16G", "32G",
				}, false),
			},

			"fill_pattern": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ARBFF",
					"IDLE",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"port_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"f",
					"np",
				}, false),
			},

			"rx_bb_credit": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"speed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unknown",
					"auto",
					"4G",
					"8G",
					"16G",
					"32G",
				}, false),
			},

			"trunk_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"un-init",
					"trunk-off",
					"trunk-on",
					"auto",
				}, false),
			},
		}),
	}
}
func getRemoteInterfaceFCPolicy(client *client.Client, dn string) (*models.InterfaceFCPolicy, error) {
	fcIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fcIfPol := models.InterfaceFCPolicyFromContainer(fcIfPolCont)

	if fcIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("InterfaceFCPolicy %s not found", fcIfPol.DistinguishedName)
	}

	return fcIfPol, nil
}

func setInterfaceFCPolicyAttributes(fcIfPol *models.InterfaceFCPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fcIfPol.DistinguishedName)
	d.Set("description", fcIfPol.Description)
	fcIfPolMap, err := fcIfPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", fcIfPolMap["name"])

	d.Set("annotation", fcIfPolMap["annotation"])
	d.Set("automaxspeed", fcIfPolMap["automaxspeed"])
	d.Set("fill_pattern", fcIfPolMap["fillPattern"])
	d.Set("name_alias", fcIfPolMap["nameAlias"])
	d.Set("port_mode", fcIfPolMap["portMode"])
	d.Set("rx_bb_credit", fcIfPolMap["rxBBCredit"])
	d.Set("speed", fcIfPolMap["speed"])
	d.Set("trunk_mode", fcIfPolMap["trunkMode"])
	return d, nil
}

func resourceAciInterfaceFCPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fcIfPol, err := getRemoteInterfaceFCPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setInterfaceFCPolicyAttributes(fcIfPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciInterfaceFCPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InterfaceFCPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fcIfPolAttr := models.InterfaceFCPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fcIfPolAttr.Annotation = Annotation.(string)
	} else {
		fcIfPolAttr.Annotation = "{}"
	}
	if Automaxspeed, ok := d.GetOk("automaxspeed"); ok {
		fcIfPolAttr.Automaxspeed = Automaxspeed.(string)
	}
	if FillPattern, ok := d.GetOk("fill_pattern"); ok {
		fcIfPolAttr.FillPattern = FillPattern.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fcIfPolAttr.NameAlias = NameAlias.(string)
	}
	if PortMode, ok := d.GetOk("port_mode"); ok {
		fcIfPolAttr.PortMode = PortMode.(string)
	}
	if RxBBCredit, ok := d.GetOk("rx_bb_credit"); ok {
		fcIfPolAttr.RxBBCredit = RxBBCredit.(string)
	}
	if Speed, ok := d.GetOk("speed"); ok {
		fcIfPolAttr.Speed = Speed.(string)
	}
	if TrunkMode, ok := d.GetOk("trunk_mode"); ok {
		fcIfPolAttr.TrunkMode = TrunkMode.(string)
	}
	fcIfPol := models.NewInterfaceFCPolicy(fmt.Sprintf("infra/fcIfPol-%s", name), "uni", desc, fcIfPolAttr)

	err := aciClient.Save(fcIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fcIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciInterfaceFCPolicyRead(ctx, d, m)
}

func resourceAciInterfaceFCPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InterfaceFCPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fcIfPolAttr := models.InterfaceFCPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fcIfPolAttr.Annotation = Annotation.(string)
	} else {
		fcIfPolAttr.Annotation = "{}"
	}
	if Automaxspeed, ok := d.GetOk("automaxspeed"); ok {
		fcIfPolAttr.Automaxspeed = Automaxspeed.(string)
	}
	if FillPattern, ok := d.GetOk("fill_pattern"); ok {
		fcIfPolAttr.FillPattern = FillPattern.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fcIfPolAttr.NameAlias = NameAlias.(string)
	}
	if PortMode, ok := d.GetOk("port_mode"); ok {
		fcIfPolAttr.PortMode = PortMode.(string)
	}
	if RxBBCredit, ok := d.GetOk("rx_bb_credit"); ok {
		fcIfPolAttr.RxBBCredit = RxBBCredit.(string)
	}
	if Speed, ok := d.GetOk("speed"); ok {
		fcIfPolAttr.Speed = Speed.(string)
	}
	if TrunkMode, ok := d.GetOk("trunk_mode"); ok {
		fcIfPolAttr.TrunkMode = TrunkMode.(string)
	}
	fcIfPol := models.NewInterfaceFCPolicy(fmt.Sprintf("infra/fcIfPol-%s", name), "uni", desc, fcIfPolAttr)

	fcIfPol.Status = "modified"

	err := aciClient.Save(fcIfPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fcIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciInterfaceFCPolicyRead(ctx, d, m)

}

func resourceAciInterfaceFCPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fcIfPol, err := getRemoteInterfaceFCPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setInterfaceFCPolicyAttributes(fcIfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciInterfaceFCPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fcIfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

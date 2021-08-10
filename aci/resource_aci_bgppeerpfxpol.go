package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciBGPPeerPrefixPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBGPPeerPrefixPolicyCreate,
		UpdateContext: resourceAciBGPPeerPrefixPolicyUpdate,
		ReadContext:   resourceAciBGPPeerPrefixPolicyRead,
		DeleteContext: resourceAciBGPPeerPrefixPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBGPPeerPrefixPolicyImport,
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

			"action": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"log",
					"reject",
					"restart",
					"shut",
				}, false),
			},

			"max_pfx": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"restart_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"thresh": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func getRemoteBGPPeerPrefixPolicy(client *client.Client, dn string) (*models.BGPPeerPrefixPolicy, error) {
	bgpPeerPfxPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpPeerPfxPol := models.BGPPeerPrefixPolicyFromContainer(bgpPeerPfxPolCont)

	if bgpPeerPfxPol.DistinguishedName == "" {
		return nil, fmt.Errorf("BGPPeerPrefixPolicy %s not found", bgpPeerPfxPol.DistinguishedName)
	}

	return bgpPeerPfxPol, nil
}

func setBGPPeerPrefixPolicyAttributes(bgpPeerPfxPol *models.BGPPeerPrefixPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	d.SetId(bgpPeerPfxPol.DistinguishedName)
	d.Set("description", bgpPeerPfxPol.Description)
	if dn != bgpPeerPfxPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}

	bgpPeerPfxPolMap, err := bgpPeerPfxPol.ToMap()

	if err != nil {
		return d, err
	}
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/bgpPfxP-%s", bgpPeerPfxPolMap["name"])))

	d.Set("name", bgpPeerPfxPolMap["name"])
	d.Set("action", bgpPeerPfxPolMap["action"])
	d.Set("annotation", bgpPeerPfxPolMap["annotation"])
	d.Set("max_pfx", bgpPeerPfxPolMap["maxPfx"])
	d.Set("name_alias", bgpPeerPfxPolMap["nameAlias"])
	d.Set("restart_time", bgpPeerPfxPolMap["restartTime"])
	d.Set("thresh", bgpPeerPfxPolMap["thresh"])

	return d, nil
}

func resourceAciBGPPeerPrefixPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bgpPeerPfxPol, err := getRemoteBGPPeerPrefixPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setBGPPeerPrefixPolicyAttributes(bgpPeerPfxPol, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBGPPeerPrefixPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BGPPeerPrefixPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpPeerPfxPolAttr := models.BGPPeerPrefixPolicyAttributes{}
	if Action, ok := d.GetOk("action"); ok {
		bgpPeerPfxPolAttr.Action = Action.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpPeerPfxPolAttr.Annotation = Annotation.(string)
	} else {
		bgpPeerPfxPolAttr.Annotation = "{}"
	}
	if MaxPfx, ok := d.GetOk("max_pfx"); ok {
		bgpPeerPfxPolAttr.MaxPfx = MaxPfx.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpPeerPfxPolAttr.NameAlias = NameAlias.(string)
	}
	if RestartTime, ok := d.GetOk("restart_time"); ok {
		bgpPeerPfxPolAttr.RestartTime = RestartTime.(string)
	}
	if Thresh, ok := d.GetOk("thresh"); ok {
		bgpPeerPfxPolAttr.Thresh = Thresh.(string)
	}
	bgpPeerPfxPol := models.NewBGPPeerPrefixPolicy(fmt.Sprintf("bgpPfxP-%s", name), TenantDn, desc, bgpPeerPfxPolAttr)

	err := aciClient.Save(bgpPeerPfxPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bgpPeerPfxPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBGPPeerPrefixPolicyRead(ctx, d, m)
}

func resourceAciBGPPeerPrefixPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BGPPeerPrefixPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpPeerPfxPolAttr := models.BGPPeerPrefixPolicyAttributes{}
	if Action, ok := d.GetOk("action"); ok {
		bgpPeerPfxPolAttr.Action = Action.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpPeerPfxPolAttr.Annotation = Annotation.(string)
	} else {
		bgpPeerPfxPolAttr.Annotation = "{}"
	}
	if MaxPfx, ok := d.GetOk("max_pfx"); ok {
		bgpPeerPfxPolAttr.MaxPfx = MaxPfx.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpPeerPfxPolAttr.NameAlias = NameAlias.(string)
	}
	if RestartTime, ok := d.GetOk("restart_time"); ok {
		bgpPeerPfxPolAttr.RestartTime = RestartTime.(string)
	}
	if Thresh, ok := d.GetOk("thresh"); ok {
		bgpPeerPfxPolAttr.Thresh = Thresh.(string)
	}
	bgpPeerPfxPol := models.NewBGPPeerPrefixPolicy(fmt.Sprintf("bgpPfxP-%s", name), TenantDn, desc, bgpPeerPfxPolAttr)

	bgpPeerPfxPol.Status = "modified"

	err := aciClient.Save(bgpPeerPfxPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(bgpPeerPfxPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBGPPeerPrefixPolicyRead(ctx, d, m)

}

func resourceAciBGPPeerPrefixPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bgpPeerPfxPol, err := getRemoteBGPPeerPrefixPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setBGPPeerPrefixPolicyAttributes(bgpPeerPfxPol, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBGPPeerPrefixPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bgpPeerPfxPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

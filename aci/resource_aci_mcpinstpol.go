package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciMiscablingProtocolInstancePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMiscablingProtocolInstancePolicyCreate,
		UpdateContext: resourceAciMiscablingProtocolInstancePolicyUpdate,
		ReadContext:   resourceAciMiscablingProtocolInstancePolicyRead,
		DeleteContext: resourceAciMiscablingProtocolInstancePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMiscablingProtocolInstancePolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},
			"ctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"pdu-per-vlan",
						"stateful-ha",
					}, false),
				},
			},
			"init_delay_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"loop_detect_mult": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"loop_protect_act": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"port-disable",
					"none",
				}, false),
			},
			"tx_freq": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tx_freq_msec": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func GetRemoteMiscablingProtocolInstancePolicy(client *client.Client, dn string) (*models.MiscablingProtocolInstancePolicy, error) {
	mcpInstPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	mcpInstPol := models.MiscablingProtocolInstancePolicyFromContainer(mcpInstPolCont)
	if mcpInstPol.DistinguishedName == "" {
		return nil, fmt.Errorf("MiscablingProtocolInstancePolicy %s not found", mcpInstPol.DistinguishedName)
	}
	return mcpInstPol, nil
}

func setMiscablingProtocolInstancePolicyAttributes(mcpInstPol *models.MiscablingProtocolInstancePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(mcpInstPol.DistinguishedName)
	d.Set("description", mcpInstPol.Description)
	mcpInstPolMap, err := mcpInstPol.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("admin_st", mcpInstPolMap["adminSt"])
	d.Set("annotation", mcpInstPolMap["annotation"])
	ctrlGet := make([]string, 0, 1)
	if mcpInstPolMap["ctrl"] == "" {
		d.Set("ctrl", ctrlGet)
	} else {
		for _, val := range strings.Split(mcpInstPolMap["ctrl"], ",") {
			ctrlGet = append(ctrlGet, strings.Trim(val, " "))
		}
		sort.Strings(ctrlGet)
		if ctrlIntr, ok := d.GetOk("ctrl"); ok {
			ctrlAct := make([]string, 0, 1)
			for _, val := range ctrlIntr.([]interface{}) {
				ctrlAct = append(ctrlAct, val.(string))
			}
			sort.Strings(ctrlAct)
			if reflect.DeepEqual(ctrlAct, ctrlGet) {
				d.Set("ctrl", d.Get("ctrl").([]interface{}))
			} else {
				d.Set("ctrl", ctrlGet)
			}
		} else {
			d.Set("ctrl", ctrlGet)
		}
	}
	d.Set("init_delay_time", mcpInstPolMap["initDelayTime"])
	d.Set("loop_detect_mult", mcpInstPolMap["loopDetectMult"])
	if mcpInstPolMap["loopProtectAct"] == "" {
		d.Set("loop_protect_act", "none")
	} else {
		d.Set("loop_protect_act", "port-disable")
	}
	d.Set("tx_freq", mcpInstPolMap["txFreq"])
	d.Set("tx_freq_msec", mcpInstPolMap["txFreqMsec"])
	d.Set("name_alias", mcpInstPolMap["nameAlias"])
	return d, nil
}

func resourceAciMiscablingProtocolInstancePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	mcpInstPol, err := GetRemoteMiscablingProtocolInstancePolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setMiscablingProtocolInstancePolicyAttributes(mcpInstPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMiscablingProtocolInstancePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MiscablingProtocolInstancePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	mcpInstPolAttr := models.MiscablingProtocolInstancePolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		mcpInstPolAttr.Annotation = Annotation.(string)
	} else {
		mcpInstPolAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_st"); ok {
		mcpInstPolAttr.AdminSt = AdminSt.(string)
	}

	if Ctrl, ok := d.GetOk("ctrl"); ok {
		ctrlList := make([]string, 0, 1)
		for _, val := range Ctrl.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		err := checkDuplicate(ctrlList)
		if err != nil {
			return diag.FromErr(err)
		}
		Ctrl := strings.Join(ctrlList, ",")
		mcpInstPolAttr.Ctrl = Ctrl
	} else {
		mcpInstPolAttr.Ctrl = "{}"
	}

	if InitDelayTime, ok := d.GetOk("init_delay_time"); ok {
		mcpInstPolAttr.InitDelayTime = InitDelayTime.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		mcpInstPolAttr.Key = Key.(string)
	}

	if LoopDetectMult, ok := d.GetOk("loop_detect_mult"); ok {
		mcpInstPolAttr.LoopDetectMult = LoopDetectMult.(string)
	}

	if LoopProtectAct, ok := d.GetOk("loop_protect_act"); ok {
		loopAct := LoopProtectAct.(string)
		if loopAct == "none" {
			mcpInstPolAttr.LoopProtectAct = "{}"
		} else {
			mcpInstPolAttr.LoopProtectAct = LoopProtectAct.(string)
		}
	}

	mcpInstPolAttr.Name = "default"

	if TxFreq, ok := d.GetOk("tx_freq"); ok {
		mcpInstPolAttr.TxFreq = TxFreq.(string)
	}

	if TxFreqMsec, ok := d.GetOk("tx_freq_msec"); ok {
		mcpInstPolAttr.TxFreqMsec = TxFreqMsec.(string)
	}
	mcpInstPol := models.NewMiscablingProtocolInstancePolicy(fmt.Sprintf("infra/mcpInstP-%s", name), "uni", desc, nameAlias, mcpInstPolAttr)
	mcpInstPol.Status = "modified"
	err := aciClient.Save(mcpInstPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(mcpInstPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciMiscablingProtocolInstancePolicyRead(ctx, d, m)
}

func resourceAciMiscablingProtocolInstancePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MiscablingProtocolInstancePolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	mcpInstPolAttr := models.MiscablingProtocolInstancePolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		mcpInstPolAttr.Annotation = Annotation.(string)
	} else {
		mcpInstPolAttr.Annotation = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_st"); ok {
		mcpInstPolAttr.AdminSt = AdminSt.(string)
	}

	if Ctrl, ok := d.GetOk("ctrl"); ok {
		ctrlList := make([]string, 0, 1)
		for _, val := range Ctrl.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		err := checkDuplicate(ctrlList)
		if err != nil {
			return diag.FromErr(err)
		}
		Ctrl := strings.Join(ctrlList, ",")
		mcpInstPolAttr.Ctrl = Ctrl
	} else {
		mcpInstPolAttr.Ctrl = "{}"
	}

	if InitDelayTime, ok := d.GetOk("init_delay_time"); ok {
		mcpInstPolAttr.InitDelayTime = InitDelayTime.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		mcpInstPolAttr.Key = Key.(string)
	}

	if LoopDetectMult, ok := d.GetOk("loop_detect_mult"); ok {
		mcpInstPolAttr.LoopDetectMult = LoopDetectMult.(string)
	}

	if LoopProtectAct, ok := d.GetOk("loop_protect_act"); ok {
		loopAct := LoopProtectAct.(string)
		if loopAct == "none" {
			mcpInstPolAttr.LoopProtectAct = "{}"
		} else {
			mcpInstPolAttr.LoopProtectAct = LoopProtectAct.(string)
		}
	}
	mcpInstPolAttr.Name = "default"

	if TxFreq, ok := d.GetOk("tx_freq"); ok {
		mcpInstPolAttr.TxFreq = TxFreq.(string)
	}

	if TxFreqMsec, ok := d.GetOk("tx_freq_msec"); ok {
		mcpInstPolAttr.TxFreqMsec = TxFreqMsec.(string)
	}
	mcpInstPol := models.NewMiscablingProtocolInstancePolicy(fmt.Sprintf("infra/mcpInstP-%s", name), "uni", desc, nameAlias, mcpInstPolAttr)
	mcpInstPol.Status = "modified"
	err := aciClient.Save(mcpInstPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(mcpInstPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciMiscablingProtocolInstancePolicyRead(ctx, d, m)
}

func resourceAciMiscablingProtocolInstancePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	mcpInstPol, err := GetRemoteMiscablingProtocolInstancePolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setMiscablingProtocolInstancePolicyAttributes(mcpInstPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciMiscablingProtocolInstancePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name mcpInstPol cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}

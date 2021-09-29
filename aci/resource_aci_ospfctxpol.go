package aci

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciOSPFTimersPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciOSPFTimersPolicyCreate,
		UpdateContext: resourceAciOSPFTimersPolicyUpdate,
		ReadContext:   resourceAciOSPFTimersPolicyRead,
		DeleteContext: resourceAciOSPFTimersPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciOSPFTimersPolicyImport,
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

			"bw_ref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"name-lookup",
						"pfx-suppress",
					}, false),
				},
			},

			"dist": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"gr_ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"helper",
					"",
				}, false),
			},

			"lsa_arrival_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"lsa_gp_pacing_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"lsa_hold_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"lsa_max_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"lsa_start_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_ecmp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_lsa_action": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"reject",
					"restart",
					"log",
				}, false),
			},

			"max_lsa_num": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_lsa_reset_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_lsa_sleep_cnt": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_lsa_sleep_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_lsa_thresh": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"spf_hold_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"spf_init_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"spf_max_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func getRemoteOSPFTimersPolicy(client *client.Client, dn string) (*models.OSPFTimersPolicy, error) {
	ospfCtxPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	ospfCtxPol := models.OSPFTimersPolicyFromContainer(ospfCtxPolCont)

	if ospfCtxPol.DistinguishedName == "" {
		return nil, fmt.Errorf("OSPFTimersPolicy %s not found", ospfCtxPol.DistinguishedName)
	}

	return ospfCtxPol, nil
}

func setOSPFTimersPolicyAttributes(ospfCtxPol *models.OSPFTimersPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(ospfCtxPol.DistinguishedName)
	d.Set("description", ospfCtxPol.Description)
	if dn != ospfCtxPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	ospfCtxPolMap, err := ospfCtxPol.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/ospfCtxP-%s", ospfCtxPolMap["name"])))
	d.Set("name", ospfCtxPolMap["name"])

	d.Set("annotation", ospfCtxPolMap["annotation"])
	d.Set("bw_ref", ospfCtxPolMap["bwRef"])
	ctrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(ospfCtxPolMap["ctrl"], ",") {
		ctrlGet = append(ctrlGet, strings.Trim(val, " "))
	}
	sort.Strings(ctrlGet)
	if len(ctrlGet) == 1 && ctrlGet[0] == "" {
		d.Set("ctrl", make([]string, 0, 1))
	} else {
		d.Set("ctrl", ctrlGet)
	}
	d.Set("dist", ospfCtxPolMap["dist"])
	d.Set("gr_ctrl", ospfCtxPolMap["grCtrl"])
	d.Set("lsa_arrival_intvl", ospfCtxPolMap["lsaArrivalIntvl"])
	d.Set("lsa_gp_pacing_intvl", ospfCtxPolMap["lsaGpPacingIntvl"])
	d.Set("lsa_hold_intvl", ospfCtxPolMap["lsaHoldIntvl"])
	d.Set("lsa_max_intvl", ospfCtxPolMap["lsaMaxIntvl"])
	d.Set("lsa_start_intvl", ospfCtxPolMap["lsaStartIntvl"])
	d.Set("max_ecmp", ospfCtxPolMap["maxEcmp"])
	d.Set("max_lsa_action", ospfCtxPolMap["maxLsaAction"])
	d.Set("max_lsa_num", ospfCtxPolMap["maxLsaNum"])
	d.Set("max_lsa_reset_intvl", ospfCtxPolMap["maxLsaResetIntvl"])
	d.Set("max_lsa_sleep_cnt", ospfCtxPolMap["maxLsaSleepCnt"])
	d.Set("max_lsa_sleep_intvl", ospfCtxPolMap["maxLsaSleepIntvl"])
	d.Set("max_lsa_thresh", ospfCtxPolMap["maxLsaThresh"])
	d.Set("name_alias", ospfCtxPolMap["nameAlias"])
	d.Set("spf_hold_intvl", ospfCtxPolMap["spfHoldIntvl"])
	d.Set("spf_init_intvl", ospfCtxPolMap["spfInitIntvl"])
	d.Set("spf_max_intvl", ospfCtxPolMap["spfMaxIntvl"])

	return d, nil
}

func resourceAciOSPFTimersPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	ospfCtxPol, err := getRemoteOSPFTimersPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setOSPFTimersPolicyAttributes(ospfCtxPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciOSPFTimersPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OSPFTimersPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	ospfCtxPolAttr := models.OSPFTimersPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ospfCtxPolAttr.Annotation = Annotation.(string)
	} else {
		ospfCtxPolAttr.Annotation = "{}"
	}
	if BwRef, ok := d.GetOk("bw_ref"); ok {
		ospfCtxPolAttr.BwRef = BwRef.(string)
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		CtrlList := make([]string, 0, 1)
		for _, val := range Ctrl.([]interface{}) {
			CtrlList = append(CtrlList, val.(string))
		}
		Ctrl := strings.Join(CtrlList, ",")
		ospfCtxPolAttr.Ctrl = Ctrl
	}
	if Dist, ok := d.GetOk("dist"); ok {
		ospfCtxPolAttr.Dist = Dist.(string)
	}
	if GrCtrl, ok := d.GetOk("gr_ctrl"); ok {
		ospfCtxPolAttr.GrCtrl = GrCtrl.(string)
	} else {
		ospfCtxPolAttr.GrCtrl = "{}"
	}
	if LsaArrivalIntvl, ok := d.GetOk("lsa_arrival_intvl"); ok {
		ospfCtxPolAttr.LsaArrivalIntvl = LsaArrivalIntvl.(string)
	}
	if LsaGpPacingIntvl, ok := d.GetOk("lsa_gp_pacing_intvl"); ok {
		ospfCtxPolAttr.LsaGpPacingIntvl = LsaGpPacingIntvl.(string)
	}
	if LsaHoldIntvl, ok := d.GetOk("lsa_hold_intvl"); ok {
		ospfCtxPolAttr.LsaHoldIntvl = LsaHoldIntvl.(string)
	}
	if LsaMaxIntvl, ok := d.GetOk("lsa_max_intvl"); ok {
		ospfCtxPolAttr.LsaMaxIntvl = LsaMaxIntvl.(string)
	}
	if LsaStartIntvl, ok := d.GetOk("lsa_start_intvl"); ok {
		ospfCtxPolAttr.LsaStartIntvl = LsaStartIntvl.(string)
	}
	if MaxEcmp, ok := d.GetOk("max_ecmp"); ok {
		ospfCtxPolAttr.MaxEcmp = MaxEcmp.(string)
	}
	if MaxLsaAction, ok := d.GetOk("max_lsa_action"); ok {
		ospfCtxPolAttr.MaxLsaAction = MaxLsaAction.(string)
	}
	if MaxLsaNum, ok := d.GetOk("max_lsa_num"); ok {
		ospfCtxPolAttr.MaxLsaNum = MaxLsaNum.(string)
	}
	if MaxLsaResetIntvl, ok := d.GetOk("max_lsa_reset_intvl"); ok {
		ospfCtxPolAttr.MaxLsaResetIntvl = MaxLsaResetIntvl.(string)
	}
	if MaxLsaSleepCnt, ok := d.GetOk("max_lsa_sleep_cnt"); ok {
		ospfCtxPolAttr.MaxLsaSleepCnt = MaxLsaSleepCnt.(string)
	}
	if MaxLsaSleepIntvl, ok := d.GetOk("max_lsa_sleep_intvl"); ok {
		ospfCtxPolAttr.MaxLsaSleepIntvl = MaxLsaSleepIntvl.(string)
	}
	if MaxLsaThresh, ok := d.GetOk("max_lsa_thresh"); ok {
		ospfCtxPolAttr.MaxLsaThresh = MaxLsaThresh.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ospfCtxPolAttr.NameAlias = NameAlias.(string)
	}
	if SpfHoldIntvl, ok := d.GetOk("spf_hold_intvl"); ok {
		ospfCtxPolAttr.SpfHoldIntvl = SpfHoldIntvl.(string)
	}
	if SpfInitIntvl, ok := d.GetOk("spf_init_intvl"); ok {
		ospfCtxPolAttr.SpfInitIntvl = SpfInitIntvl.(string)
	}
	if SpfMaxIntvl, ok := d.GetOk("spf_max_intvl"); ok {
		ospfCtxPolAttr.SpfMaxIntvl = SpfMaxIntvl.(string)
	}
	ospfCtxPol := models.NewOSPFTimersPolicy(fmt.Sprintf("ospfCtxP-%s", name), TenantDn, desc, ospfCtxPolAttr)

	err := aciClient.Save(ospfCtxPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ospfCtxPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciOSPFTimersPolicyRead(ctx, d, m)
}

func resourceAciOSPFTimersPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OSPFTimersPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	ospfCtxPolAttr := models.OSPFTimersPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ospfCtxPolAttr.Annotation = Annotation.(string)
	} else {
		ospfCtxPolAttr.Annotation = "{}"
	}
	if BwRef, ok := d.GetOk("bw_ref"); ok {
		ospfCtxPolAttr.BwRef = BwRef.(string)
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		CtrlList := make([]string, 0, 1)
		for _, val := range Ctrl.([]interface{}) {
			CtrlList = append(CtrlList, val.(string))
		}
		Ctrl := strings.Join(CtrlList, ",")
		ospfCtxPolAttr.Ctrl = Ctrl
	}
	if Dist, ok := d.GetOk("dist"); ok {
		ospfCtxPolAttr.Dist = Dist.(string)
	}
	if GrCtrl, ok := d.GetOk("gr_ctrl"); ok {
		ospfCtxPolAttr.GrCtrl = GrCtrl.(string)
	} else {
		ospfCtxPolAttr.GrCtrl = "{}"
	}
	if LsaArrivalIntvl, ok := d.GetOk("lsa_arrival_intvl"); ok {
		ospfCtxPolAttr.LsaArrivalIntvl = LsaArrivalIntvl.(string)
	}
	if LsaGpPacingIntvl, ok := d.GetOk("lsa_gp_pacing_intvl"); ok {
		ospfCtxPolAttr.LsaGpPacingIntvl = LsaGpPacingIntvl.(string)
	}
	if LsaHoldIntvl, ok := d.GetOk("lsa_hold_intvl"); ok {
		ospfCtxPolAttr.LsaHoldIntvl = LsaHoldIntvl.(string)
	}
	if LsaMaxIntvl, ok := d.GetOk("lsa_max_intvl"); ok {
		ospfCtxPolAttr.LsaMaxIntvl = LsaMaxIntvl.(string)
	}
	if LsaStartIntvl, ok := d.GetOk("lsa_start_intvl"); ok {
		ospfCtxPolAttr.LsaStartIntvl = LsaStartIntvl.(string)
	}
	if MaxEcmp, ok := d.GetOk("max_ecmp"); ok {
		ospfCtxPolAttr.MaxEcmp = MaxEcmp.(string)
	}
	if MaxLsaAction, ok := d.GetOk("max_lsa_action"); ok {
		ospfCtxPolAttr.MaxLsaAction = MaxLsaAction.(string)
	}
	if MaxLsaNum, ok := d.GetOk("max_lsa_num"); ok {
		ospfCtxPolAttr.MaxLsaNum = MaxLsaNum.(string)
	}
	if MaxLsaResetIntvl, ok := d.GetOk("max_lsa_reset_intvl"); ok {
		ospfCtxPolAttr.MaxLsaResetIntvl = MaxLsaResetIntvl.(string)
	}
	if MaxLsaSleepCnt, ok := d.GetOk("max_lsa_sleep_cnt"); ok {
		ospfCtxPolAttr.MaxLsaSleepCnt = MaxLsaSleepCnt.(string)
	}
	if MaxLsaSleepIntvl, ok := d.GetOk("max_lsa_sleep_intvl"); ok {
		ospfCtxPolAttr.MaxLsaSleepIntvl = MaxLsaSleepIntvl.(string)
	}
	if MaxLsaThresh, ok := d.GetOk("max_lsa_thresh"); ok {
		ospfCtxPolAttr.MaxLsaThresh = MaxLsaThresh.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ospfCtxPolAttr.NameAlias = NameAlias.(string)
	}
	if SpfHoldIntvl, ok := d.GetOk("spf_hold_intvl"); ok {
		ospfCtxPolAttr.SpfHoldIntvl = SpfHoldIntvl.(string)
	}
	if SpfInitIntvl, ok := d.GetOk("spf_init_intvl"); ok {
		ospfCtxPolAttr.SpfInitIntvl = SpfInitIntvl.(string)
	}
	if SpfMaxIntvl, ok := d.GetOk("spf_max_intvl"); ok {
		ospfCtxPolAttr.SpfMaxIntvl = SpfMaxIntvl.(string)
	}
	ospfCtxPol := models.NewOSPFTimersPolicy(fmt.Sprintf("ospfCtxP-%s", name), TenantDn, desc, ospfCtxPolAttr)

	ospfCtxPol.Status = "modified"

	err := aciClient.Save(ospfCtxPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ospfCtxPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciOSPFTimersPolicyRead(ctx, d, m)

}

func resourceAciOSPFTimersPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	ospfCtxPol, err := getRemoteOSPFTimersPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setOSPFTimersPolicyAttributes(ospfCtxPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciOSPFTimersPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "ospfCtxPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

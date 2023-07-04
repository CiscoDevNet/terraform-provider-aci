package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciIGMPInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciIGMPInterfacePolicyCreate,
		UpdateContext: resourceAciIGMPInterfacePolicyUpdate,
		ReadContext:   resourceAciIGMPInterfacePolicyRead,
		DeleteContext: resourceAciIGMPInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciIGMPInterfacePolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"grp_timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"if_ctrl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"allow-v3-asm",
						"fast-leave",
						"rep-ll",
					}, false),
				},
			},
			"last_mbr_cnt": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_mbr_resp_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"querier_timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"query_intvl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"robust_fac": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rsp_intvl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"start_query_cnt": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"start_query_intvl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ver": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"v2",
					"v3",
				}, false),
			},
		})),
	}
}

func getRemoteIGMPInterfacePolicy(client *client.Client, dn string) (*models.IGMPInterfacePolicy, error) {
	igmpIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	igmpIfPol := models.IGMPInterfacePolicyFromContainer(igmpIfPolCont)
	if igmpIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("IGMPInterfacePolicy %s not found", dn)
	}
	return igmpIfPol, nil
}

func setIGMPInterfacePolicyAttributes(igmpIfPol *models.IGMPInterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(igmpIfPol.DistinguishedName)
	d.Set("description", igmpIfPol.Description)
	igmpIfPolMap, err := igmpIfPol.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != igmpIfPol.DistinguishedName {
		d.Set("tenant_dn", "")
	} else {
		d.Set("tenant_dn", GetParentDn(igmpIfPol.DistinguishedName, fmt.Sprintf("/"+models.RnIgmpIfPol, igmpIfPolMap["name"])))
	}
	d.Set("annotation", igmpIfPolMap["annotation"])
	d.Set("grp_timeout", igmpIfPolMap["grpTimeout"])
	ifCtrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(igmpIfPolMap["ifCtrl"], ",") {
		ifCtrlGet = append(ifCtrlGet, strings.Trim(val, " "))
	}
	sort.Strings(ifCtrlGet)
	if ifCtrlIntr, ok := d.GetOk("if_ctrl"); ok {
		ifCtrlAct := make([]string, 0, 1)
		for _, val := range ifCtrlIntr.([]interface{}) {
			ifCtrlAct = append(ifCtrlAct, val.(string))
		}
		sort.Strings(ifCtrlAct)
		if reflect.DeepEqual(ifCtrlAct, ifCtrlGet) {
			d.Set("if_ctrl", d.Get("if_ctrl").([]interface{}))
		} else {
			d.Set("if_ctrl", ifCtrlGet)
		}
	} else {
		d.Set("if_ctrl", ifCtrlGet)
	}
	d.Set("last_mbr_cnt", igmpIfPolMap["lastMbrCnt"])
	d.Set("last_mbr_resp_time", igmpIfPolMap["lastMbrRespTime"])
	d.Set("name", igmpIfPolMap["name"])
	d.Set("name_alias", igmpIfPolMap["nameAlias"])
	d.Set("querier_timeout", igmpIfPolMap["querierTimeout"])
	d.Set("query_intvl", igmpIfPolMap["queryIntvl"])
	d.Set("robust_fac", igmpIfPolMap["robustFac"])
	d.Set("rsp_intvl", igmpIfPolMap["rspIntvl"])
	d.Set("start_query_cnt", igmpIfPolMap["startQueryCnt"])
	d.Set("start_query_intvl", igmpIfPolMap["startQueryIntvl"])
	d.Set("ver", igmpIfPolMap["ver"])
	return d, nil
}

func resourceAciIGMPInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	igmpIfPol, err := getRemoteIGMPInterfacePolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setIGMPInterfacePolicyAttributes(igmpIfPol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciIGMPInterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] IGMPInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	igmpIfPolAttr := models.IGMPInterfacePolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		igmpIfPolAttr.Annotation = Annotation.(string)
	} else {
		igmpIfPolAttr.Annotation = "{}"
	}

	if GrpTimeout, ok := d.GetOk("grp_timeout"); ok {
		igmpIfPolAttr.GrpTimeout = GrpTimeout.(string)
	}

	if IfCtrl, ok := d.GetOk("if_ctrl"); ok {
		ifCtrlList := make([]string, 0, 1)
		for _, val := range IfCtrl.([]interface{}) {
			ifCtrlList = append(ifCtrlList, val.(string))
		}
		IfCtrl := strings.Join(ifCtrlList, ",")
		igmpIfPolAttr.IfCtrl = IfCtrl
	}

	if LastMbrCnt, ok := d.GetOk("last_mbr_cnt"); ok {
		igmpIfPolAttr.LastMbrCnt = LastMbrCnt.(string)
	}

	if LastMbrRespTime, ok := d.GetOk("last_mbr_resp_time"); ok {
		igmpIfPolAttr.LastMbrRespTime = LastMbrRespTime.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		igmpIfPolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		igmpIfPolAttr.NameAlias = NameAlias.(string)
	}

	if QuerierTimeout, ok := d.GetOk("querier_timeout"); ok {
		igmpIfPolAttr.QuerierTimeout = QuerierTimeout.(string)
	}

	if QueryIntvl, ok := d.GetOk("query_intvl"); ok {
		igmpIfPolAttr.QueryIntvl = QueryIntvl.(string)
	}

	if RobustFac, ok := d.GetOk("robust_fac"); ok {
		igmpIfPolAttr.RobustFac = RobustFac.(string)
	}

	if RspIntvl, ok := d.GetOk("rsp_intvl"); ok {
		igmpIfPolAttr.RspIntvl = RspIntvl.(string)
	}

	if StartQueryCnt, ok := d.GetOk("start_query_cnt"); ok {
		igmpIfPolAttr.StartQueryCnt = StartQueryCnt.(string)
	}

	if StartQueryIntvl, ok := d.GetOk("start_query_intvl"); ok {
		igmpIfPolAttr.StartQueryIntvl = StartQueryIntvl.(string)
	}

	if Ver, ok := d.GetOk("ver"); ok {
		igmpIfPolAttr.Ver = Ver.(string)
	}
	igmpIfPol := models.NewIGMPInterfacePolicy(fmt.Sprintf(models.RnIgmpIfPol, name), TenantDn, desc, igmpIfPolAttr)

	err := aciClient.Save(igmpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(igmpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciIGMPInterfacePolicyRead(ctx, d, m)
}
func resourceAciIGMPInterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] IGMPInterfacePolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	igmpIfPolAttr := models.IGMPInterfacePolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		igmpIfPolAttr.Annotation = Annotation.(string)
	} else {
		igmpIfPolAttr.Annotation = "{}"
	}

	if GrpTimeout, ok := d.GetOk("grp_timeout"); ok {
		igmpIfPolAttr.GrpTimeout = GrpTimeout.(string)
	}
	if IfCtrl, ok := d.GetOk("if_ctrl"); ok {
		ifCtrlList := make([]string, 0, 1)
		for _, val := range IfCtrl.([]interface{}) {
			ifCtrlList = append(ifCtrlList, val.(string))
		}
		IfCtrl := strings.Join(ifCtrlList, ",")
		igmpIfPolAttr.IfCtrl = IfCtrl
	}

	if LastMbrCnt, ok := d.GetOk("last_mbr_cnt"); ok {
		igmpIfPolAttr.LastMbrCnt = LastMbrCnt.(string)
	}

	if LastMbrRespTime, ok := d.GetOk("last_mbr_resp_time"); ok {
		igmpIfPolAttr.LastMbrRespTime = LastMbrRespTime.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		igmpIfPolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		igmpIfPolAttr.NameAlias = NameAlias.(string)
	}

	if QuerierTimeout, ok := d.GetOk("querier_timeout"); ok {
		igmpIfPolAttr.QuerierTimeout = QuerierTimeout.(string)
	}

	if QueryIntvl, ok := d.GetOk("query_intvl"); ok {
		igmpIfPolAttr.QueryIntvl = QueryIntvl.(string)
	}

	if RobustFac, ok := d.GetOk("robust_fac"); ok {
		igmpIfPolAttr.RobustFac = RobustFac.(string)
	}

	if RspIntvl, ok := d.GetOk("rsp_intvl"); ok {
		igmpIfPolAttr.RspIntvl = RspIntvl.(string)
	}

	if StartQueryCnt, ok := d.GetOk("start_query_cnt"); ok {
		igmpIfPolAttr.StartQueryCnt = StartQueryCnt.(string)
	}

	if StartQueryIntvl, ok := d.GetOk("start_query_intvl"); ok {
		igmpIfPolAttr.StartQueryIntvl = StartQueryIntvl.(string)
	}

	if Ver, ok := d.GetOk("ver"); ok {
		igmpIfPolAttr.Ver = Ver.(string)
	}
	igmpIfPol := models.NewIGMPInterfacePolicy(fmt.Sprintf(models.RnIgmpIfPol, name), TenantDn, desc, igmpIfPolAttr)

	igmpIfPol.Status = "modified"

	err := aciClient.Save(igmpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(igmpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciIGMPInterfacePolicyRead(ctx, d, m)
}

func resourceAciIGMPInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	igmpIfPol, err := getRemoteIGMPInterfacePolicy(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setIGMPInterfacePolicyAttributes(igmpIfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciIGMPInterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "igmpIfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

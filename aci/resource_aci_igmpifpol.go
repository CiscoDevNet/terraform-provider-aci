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
			"group_timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"control": {
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
				DiffSuppressFunc: suppressTypeListDiffFunc,
			},
			"last_member_count": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_member_response_time": {
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
			"query_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"robustness_variable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"response_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"startup_query_count": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"startup_query_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"v2",
					"v3",
				}, false),
			},
			"maximum_mulitcast_entries": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"reserved_mulitcast_entries": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"state_limit_route_map": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"report_policy_route_map": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"static_report_route_map": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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
		return nil, fmt.Errorf("IGMP Interface Policy %s not found", dn)
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
	d.Set("group_timeout", igmpIfPolMap["grpTimeout"])
	ifCtrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(igmpIfPolMap["ifCtrl"], ",") {
		ifCtrlGet = append(ifCtrlGet, strings.Trim(val, " "))
	}
	sort.Strings(ifCtrlGet)
	if ifCtrlIntr, ok := d.GetOk("control"); ok {
		ifCtrlAct := make([]string, 0, 1)
		for _, val := range ifCtrlIntr.([]interface{}) {
			ifCtrlAct = append(ifCtrlAct, val.(string))
		}
		sort.Strings(ifCtrlAct)
		if reflect.DeepEqual(ifCtrlAct, ifCtrlGet) {
			d.Set("control", d.Get("control").([]interface{}))
		} else {
			d.Set("control", ifCtrlGet)
		}
	} else {
		d.Set("control", ifCtrlGet)
	}
	d.Set("last_member_count", igmpIfPolMap["lastMbrCnt"])
	d.Set("last_member_response_time", igmpIfPolMap["lastMbrRespTime"])
	d.Set("name", igmpIfPolMap["name"])
	d.Set("name_alias", igmpIfPolMap["nameAlias"])
	d.Set("querier_timeout", igmpIfPolMap["querierTimeout"])
	d.Set("query_interval", igmpIfPolMap["queryIntvl"])
	d.Set("robustness_variable", igmpIfPolMap["robustFac"])
	d.Set("response_interval", igmpIfPolMap["rspIntvl"])
	d.Set("startup_query_count", igmpIfPolMap["startQueryCnt"])
	d.Set("startup_query_interval", igmpIfPolMap["startQueryIntvl"])
	d.Set("version", igmpIfPolMap["ver"])
	return d, nil
}

func getandSetIGMPIfPolRelationshipAttributes(aciClient *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	igmpRepPolData, err := aciClient.ReadRelationigmpRepPolrtdmcRsFilterToRtMapPol(fmt.Sprintf("%s/%s", dn, models.RnIgmpRepPol))
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation igmpRepPolData %v", err)
		d.Set("report_policy_route_map", "")
	} else {
		d.Set("report_policy_route_map", igmpRepPolData.(string))
	}

	igmpStateLPolData, err := aciClient.ReadRelationigmpStateLPolrtdmcRsFilterToRtMapPol(fmt.Sprintf("%s/%s", dn, models.RnIgmpStateLPol))
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation igmpStateLPol %v", err)
		d.Set("state_limit_route_map", "")
	} else {
		igmpStateLPolCont, err := aciClient.Get(fmt.Sprintf("%s/%s", dn, models.RnIgmpStateLPol))
		igmpStateLPolMap, err := models.IGMPStateLimitPolicyFromContainer(igmpStateLPolCont).ToMap()
		if err != nil {
			return d, err
		}
		d.Set("maximum_mulitcast_entries", igmpStateLPolMap["max"])
		d.Set("reserved_mulitcast_entries", igmpStateLPolMap["rsvd"])
		d.Set("state_limit_route_map", igmpStateLPolData.(string))
	}

	igmpStRepPolData, err := aciClient.ReadRelationigmpStRepPolrtdmcRsFilterToRtMapPol(fmt.Sprintf("%s/%s", dn, models.RnIgmpStRepPol))
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation igmpStRepPol %v", err)
		d.Set("static_report_route_map", "")
	} else {
		d.Set("static_report_route_map", igmpStRepPolData.(string))
	}
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

	_, err = getandSetIGMPIfPolRelationshipAttributes(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] IGMP Interface Policy Relationship Attributes - Read finished successfully")
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciIGMPInterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] IGMP Interface Policy: Beginning Creation")
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

	if GrpTimeout, ok := d.GetOk("group_timeout"); ok {
		igmpIfPolAttr.GrpTimeout = GrpTimeout.(string)
	}

	if IfCtrl, ok := d.GetOk("control"); ok {
		ifCtrlList := make([]string, 0, 1)
		for _, val := range IfCtrl.([]interface{}) {
			ifCtrlList = append(ifCtrlList, val.(string))
		}
		IfCtrl := strings.Join(ifCtrlList, ",")
		igmpIfPolAttr.IfCtrl = IfCtrl
	}

	if LastMbrCnt, ok := d.GetOk("last_member_count"); ok {
		igmpIfPolAttr.LastMbrCnt = LastMbrCnt.(string)
	}

	if LastMbrRespTime, ok := d.GetOk("last_member_response_time"); ok {
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

	if QueryIntvl, ok := d.GetOk("query_interval"); ok {
		igmpIfPolAttr.QueryIntvl = QueryIntvl.(string)
	}

	if RobustFac, ok := d.GetOk("robustness_variable"); ok {
		igmpIfPolAttr.RobustFac = RobustFac.(string)
	}

	if RspIntvl, ok := d.GetOk("response_interval"); ok {
		igmpIfPolAttr.RspIntvl = RspIntvl.(string)
	}

	if StartQueryCnt, ok := d.GetOk("startup_query_count"); ok {
		igmpIfPolAttr.StartQueryCnt = StartQueryCnt.(string)
	}

	if StartQueryIntvl, ok := d.GetOk("startup_query_interval"); ok {
		igmpIfPolAttr.StartQueryIntvl = StartQueryIntvl.(string)
	}

	if version, ok := d.GetOk("version"); ok {
		igmpIfPolAttr.Ver = version.(string)
	}
	igmpIfPol := models.NewIGMPInterfacePolicy(fmt.Sprintf(models.RnIgmpIfPol, name), TenantDn, desc, igmpIfPolAttr)

	err := aciClient.Save(igmpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)
	if stateLimitRouteMap, ok := d.GetOk("state_limit_route_map"); ok {
		relationParam := stateLimitRouteMap.(string)
		checkDns = append(checkDns, relationParam)
	}

	if reportPolicyRouteMap, ok := d.GetOk("report_policy_route_map"); ok {
		relationParam := reportPolicyRouteMap.(string)
		checkDns = append(checkDns, relationParam)
	}

	if staticReportRouteMap, ok := d.GetOk("static_report_route_map"); ok {
		relationParam := staticReportRouteMap.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if stateLimitRouteMap, ok := d.GetOk("state_limit_route_map"); ok {

		igmpStateLPolAttr := models.IGMPStateLimitPolicyAttributes{}

		if maxMulticastEntries, ok := d.GetOk("maximum_mulitcast_entries"); ok {
			igmpStateLPolAttr.Max = maxMulticastEntries.(string)
		}

		if reservedMulitcastEntries, ok := d.GetOk("reserved_mulitcast_entries"); ok {
			igmpStateLPolAttr.Rsvd = reservedMulitcastEntries.(string)
		}

		igmpStateLPol := models.NewIGMPStateLimitPolicy(igmpIfPol.DistinguishedName, "", igmpStateLPolAttr)
		err := aciClient.Save(igmpStateLPol)
		if err != nil {
			return diag.FromErr(err)
		}
		relationParam := stateLimitRouteMap.(string)
		err = aciClient.CreateRelationigmpStateLPolrtdmcRsFilterToRtMapPol(igmpStateLPol.DistinguishedName, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if reportPolicyRouteMap, ok := d.GetOk("report_policy_route_map"); ok {
		igmpRepPol := models.NewIGMPReportPolicy(igmpIfPol.DistinguishedName, "", models.IGMPReportPolicyAttributes{})
		err := aciClient.Save(igmpRepPol)
		if err != nil {
			return diag.FromErr(err)
		}
		relationParam := reportPolicyRouteMap.(string)
		err = aciClient.CreateRelationigmpRepPolrtdmcRsFilterToRtMapPol(igmpRepPol.DistinguishedName, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if staticReportRouteMap, ok := d.GetOk("static_report_route_map"); ok {
		igmpStRepPol := models.NewIGMPStaticReportPolicy(igmpIfPol.DistinguishedName, "", models.IGMPStaticReportPolicyAttributes{})
		err := aciClient.Save(igmpStRepPol)
		if err != nil {
			return diag.FromErr(err)
		}
		relationParam := staticReportRouteMap.(string)
		err = aciClient.CreateRelationigmpStRepPolrtdmcRsFilterToRtMapPol(igmpStRepPol.DistinguishedName, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(igmpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciIGMPInterfacePolicyRead(ctx, d, m)
}
func resourceAciIGMPInterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] IGMP Interface Policy: Beginning Update")
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

	if GrpTimeout, ok := d.GetOk("group_timeout"); ok {
		igmpIfPolAttr.GrpTimeout = GrpTimeout.(string)
	}
	if IfCtrl, ok := d.GetOk("control"); ok {
		ifCtrlList := make([]string, 0, 1)
		for _, val := range IfCtrl.([]interface{}) {
			ifCtrlList = append(ifCtrlList, val.(string))
		}
		IfCtrl := strings.Join(ifCtrlList, ",")
		igmpIfPolAttr.IfCtrl = IfCtrl
	}

	if LastMbrCnt, ok := d.GetOk("last_member_count"); ok {
		igmpIfPolAttr.LastMbrCnt = LastMbrCnt.(string)
	}

	if LastMbrRespTime, ok := d.GetOk("last_member_response_time"); ok {
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

	if QueryIntvl, ok := d.GetOk("query_interval"); ok {
		igmpIfPolAttr.QueryIntvl = QueryIntvl.(string)
	}

	if RobustFac, ok := d.GetOk("robustness_variable"); ok {
		igmpIfPolAttr.RobustFac = RobustFac.(string)
	}

	if RspIntvl, ok := d.GetOk("response_interval"); ok {
		igmpIfPolAttr.RspIntvl = RspIntvl.(string)
	}

	if StartQueryCnt, ok := d.GetOk("startup_query_count"); ok {
		igmpIfPolAttr.StartQueryCnt = StartQueryCnt.(string)
	}

	if StartQueryIntvl, ok := d.GetOk("startup_query_interval"); ok {
		igmpIfPolAttr.StartQueryIntvl = StartQueryIntvl.(string)
	}

	if version, ok := d.GetOk("version"); ok {
		igmpIfPolAttr.Ver = version.(string)
	}
	igmpIfPol := models.NewIGMPInterfacePolicy(fmt.Sprintf(models.RnIgmpIfPol, name), TenantDn, desc, igmpIfPolAttr)

	igmpIfPol.Status = "modified"

	err := aciClient.Save(igmpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)
	if d.HasChange("state_limit_route_map") {
		_, newRelParam := d.GetChange("state_limit_route_map")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("report_policy_route_map") {
		_, newRelParam := d.GetChange("report_policy_route_map")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("static_report_route_map") {
		_, newRelParam := d.GetChange("static_report_route_map")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("state_limit_route_map") || d.HasChange("maximum_mulitcast_entries") || d.HasChange("reserved_mulitcast_entries") {
		_, newRelParam := d.GetChange("state_limit_route_map")

		igmpStateLPolAttr := models.IGMPStateLimitPolicyAttributes{}

		if maxMulticastEntries, ok := d.GetOk("maximum_mulitcast_entries"); ok {
			igmpStateLPolAttr.Max = maxMulticastEntries.(string)
		}

		if reservedMulitcastEntries, ok := d.GetOk("reserved_mulitcast_entries"); ok {
			igmpStateLPolAttr.Rsvd = reservedMulitcastEntries.(string)
		}

		igmpStateLPol := models.NewIGMPStateLimitPolicy(igmpIfPol.DistinguishedName, "", igmpStateLPolAttr)

		igmpStateLPol.Status = "created,modified"

		err := aciClient.Save(igmpStateLPol)
		if err != nil {
			return diag.FromErr(err)
		}

		err = aciClient.CreateRelationigmpStateLPolrtdmcRsFilterToRtMapPol(igmpStateLPol.DistinguishedName, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if d.HasChange("report_policy_route_map") {
		_, newRelParam := d.GetChange("report_policy_route_map")

		igmpRepPol := models.NewIGMPReportPolicy(igmpIfPol.DistinguishedName, "", models.IGMPReportPolicyAttributes{})

		igmpRepPol.Status = "created,modified"

		err := aciClient.Save(igmpRepPol)
		if err != nil {
			return diag.FromErr(err)
		}

		err = aciClient.CreateRelationigmpRepPolrtdmcRsFilterToRtMapPol(igmpRepPol.DistinguishedName, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if d.HasChange("static_report_route_map") {
		_, newRelParam := d.GetChange("static_report_route_map")
		igmpStRepPol := models.NewIGMPStaticReportPolicy(igmpIfPol.DistinguishedName, "", models.IGMPStaticReportPolicyAttributes{})

		igmpStRepPol.Status = "created,modified"

		err := aciClient.Save(igmpStRepPol)
		if err != nil {
			return diag.FromErr(err)
		}

		err = aciClient.CreateRelationigmpStRepPolrtdmcRsFilterToRtMapPol(igmpStRepPol.DistinguishedName, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

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

	_, err = getandSetIGMPIfPolRelationshipAttributes(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] IGMP Interface Policy Relationship Attributes - Read finished successfully")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciIGMPInterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, models.IgmpIfPolClassName)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

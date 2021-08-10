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

func resourceAciVMMController() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciVMMControllerCreate,
		UpdateContext: resourceAciVMMControllerUpdate,
		ReadContext:   resourceAciVMMControllerRead,
		DeleteContext: resourceAciVMMControllerDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVMMControllerImport,
		},

		SchemaVersion: 1,
		Schema: AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Default:  "orchestrator:terraform",
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				}},
			"dvs_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"5.1",
					"5.5",
					"6.0",
					"6.5",
					"6.6",
					"7.0",
					"unmanaged",
				}, false),
			},
			"host_or_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"inventory_trig_st": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"autoTriggered",
					"triggered",
					"untriggered",
				}, false),
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"cf",
					"default",
					"k8s",
					"n1kv",
					"nsx",
					"openshift",
					"ovs",
					"rancher",
					"rhev",
					"unknown",
				}, false),
			},
			"msft_config_err_msg": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"msft_config_issues": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"aaacert-invalid",
						"duplicate-mac-in-inventory",
						"duplicate-rootContName",
						"invalid-object-in-inventory",
						"invalid-rootContName",
						"inventory-failed",
						"missing-hostGroup-in-cloud",
						"missing-rootContName",
						"not-applicable",
						"zero-mac-in-inventory",
					}, false),
				},
			},
			"n1kv_stats_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
					"unknown",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"root_cont_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MicrosoftSCVMM",
					"cloudfoundry",
					"iaas",
					"kubernetes",
					"network",
					"nsx",
					"openshift",
					"openstack",
					"rhev",
					"unmanaged",
					"vm",
				}, false),
			},
			"seq_num": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stats_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
					"unknown",
				}, false),
			},
			"vxlan_depl_pref": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"nsx",
					"vxlan",
				}, false),
			},

			"relation_vmm_rs_acc": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to vmm:UsrAccP",
			},
			"relation_vmm_rs_ctrlr_p_mon_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to mon:InfraPol",
			},
			"relation_vmm_rs_mcast_addr_ns": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fvns:McastAddrInstP",
			},
			"relation_vmm_rs_mgmt_e_pg": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fv:EPg",
			},
			"relation_vmm_rs_to_ext_dev_mgr": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to extdev:MgrP",
				Set:         schema.HashString,
			},
			"relation_vmm_rs_vmm_ctrlr_p": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to vmmCtrlrP",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"epg_depl_pref": {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"both",
								"local",
							}, false),
						},
						"target_dn": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"relation_vmm_rs_vxlan_ns": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fvns:VxlanInstP",
			},
			"relation_vmm_rs_vxlan_ns_def": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Create relation to fvns:AInstP",
			}}),
	}
}

func getRemoteVMMController(client *client.Client, dn string) (*models.VMMController, error) {
	vmmCtrlrPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vmmCtrlrP := models.VMMControllerFromContainer(vmmCtrlrPCont)
	if vmmCtrlrP.DistinguishedName == "" {
		return nil, fmt.Errorf("VMMController %s not found", vmmCtrlrP.DistinguishedName)
	}
	return vmmCtrlrP, nil
}

func setVMMControllerAttributes(vmmCtrlrP *models.VMMController, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(vmmCtrlrP.DistinguishedName)
	vmmCtrlrPMap, err := vmmCtrlrP.ToMap()

	if err != nil {
		return d, err
	}
	d.Set("annotation", vmmCtrlrPMap["annotation"])
	d.Set("dvs_version", vmmCtrlrPMap["dvsVersion"])
	d.Set("host_or_ip", vmmCtrlrPMap["hostOrIp"])
	// d.Set("inventory_trig_st", vmmCtrlrPMap["inventoryTrigSt"])
	d.Set("mode", vmmCtrlrPMap["mode"])
	d.Set("msft_config_err_msg", vmmCtrlrPMap["msftConfigErrMsg"])
	msftConfigIssuesGet := make([]string, 0, 1)
	for _, val := range strings.Split(vmmCtrlrPMap["msftConfigIssues"], ",") {
		msftConfigIssuesGet = append(msftConfigIssuesGet, strings.Trim(val, " "))
	}
	sort.Strings(msftConfigIssuesGet)
	if msftConfigIssuesIntr, ok := d.GetOk("msft_config_issues"); ok {
		msftConfigIssuesAct := make([]string, 0, 1)
		for _, val := range msftConfigIssuesIntr.([]interface{}) {
			msftConfigIssuesAct = append(msftConfigIssuesAct, val.(string))
		}
		sort.Strings(msftConfigIssuesAct)
		if reflect.DeepEqual(msftConfigIssuesAct, msftConfigIssuesGet) {
			d.Set("msft_config_issues", d.Get("msft_config_issues").([]interface{}))
		} else {
			d.Set("msft_config_issues", msftConfigIssuesGet)
		}
	} else {
		d.Set("msft_config_issues", msftConfigIssuesGet)
	}
	d.Set("vmm_domain_dn", GetParentDn(vmmCtrlrP.DistinguishedName, fmt.Sprintf("/ctrlr-%s", vmmCtrlrPMap["name"])))
	d.Set("n1kv_stats_mode", vmmCtrlrPMap["n1kvStatsMode"])
	d.Set("name", vmmCtrlrPMap["name"])
	d.Set("port", vmmCtrlrPMap["port"])
	d.Set("root_cont_name", vmmCtrlrPMap["rootContName"])
	d.Set("scope", vmmCtrlrPMap["scope"])
	d.Set("seq_num", vmmCtrlrPMap["seqNum"])
	d.Set("stats_mode", vmmCtrlrPMap["statsMode"])
	d.Set("vxlan_depl_pref", vmmCtrlrPMap["vxlanDeplPref"])
	d.Set("name_alias", vmmCtrlrPMap["nameAlias"])
	return d, nil
}

func resourceAciVMMControllerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vmmCtrlrP, err := getRemoteVMMController(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setVMMControllerAttributes(vmmCtrlrP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVMMControllerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VMMController: Beginning Creation")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)

	vmmCtrlrPAttr := models.VMMControllerAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmCtrlrPAttr.Annotation = Annotation.(string)
	} else {
		vmmCtrlrPAttr.Annotation = "{}"
	}

	if DvsVersion, ok := d.GetOk("dvs_version"); ok {
		vmmCtrlrPAttr.DvsVersion = DvsVersion.(string)
	}

	if HostOrIp, ok := d.GetOk("host_or_ip"); ok {
		vmmCtrlrPAttr.HostOrIp = HostOrIp.(string)
	}

	if InventoryTrigSt, ok := d.GetOk("inventory_trig_st"); ok {
		vmmCtrlrPAttr.InventoryTrigSt = InventoryTrigSt.(string)
	}

	if Mode, ok := d.GetOk("mode"); ok {
		vmmCtrlrPAttr.Mode = Mode.(string)
	}

	if MsftConfigErrMsg, ok := d.GetOk("msft_config_err_msg"); ok {
		vmmCtrlrPAttr.MsftConfigErrMsg = MsftConfigErrMsg.(string)
	}

	if MsftConfigIssues, ok := d.GetOk("msft_config_issues"); ok {
		msftConfigIssuesList := make([]string, 0, 1)
		for _, val := range MsftConfigIssues.([]interface{}) {
			msftConfigIssuesList = append(msftConfigIssuesList, val.(string))
		}
		MsftConfigIssues := strings.Join(msftConfigIssuesList, ",")
		vmmCtrlrPAttr.MsftConfigIssues = MsftConfigIssues
	}

	if N1kvStatsMode, ok := d.GetOk("n1kv_stats_mode"); ok {
		vmmCtrlrPAttr.N1kvStatsMode = N1kvStatsMode.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vmmCtrlrPAttr.Name = Name.(string)
	}

	if Port, ok := d.GetOk("port"); ok {
		vmmCtrlrPAttr.Port = Port.(string)
	}

	if RootContName, ok := d.GetOk("root_cont_name"); ok {
		vmmCtrlrPAttr.RootContName = RootContName.(string)
	}

	if Scope, ok := d.GetOk("scope"); ok {
		vmmCtrlrPAttr.Scope = Scope.(string)
	}

	if SeqNum, ok := d.GetOk("seq_num"); ok {
		vmmCtrlrPAttr.SeqNum = SeqNum.(string)
	}

	if StatsMode, ok := d.GetOk("stats_mode"); ok {
		vmmCtrlrPAttr.StatsMode = StatsMode.(string)
	}

	if VxlanDeplPref, ok := d.GetOk("vxlan_depl_pref"); ok {
		vmmCtrlrPAttr.VxlanDeplPref = VxlanDeplPref.(string)
	}
	vmmCtrlrP := models.NewVMMController(fmt.Sprintf("ctrlr-%s", name), VMMDomainDn, nameAlias, vmmCtrlrPAttr)

	err := aciClient.Save(vmmCtrlrP)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTovmmRsAcc, ok := d.GetOk("relation_vmm_rs_acc"); ok {
		relationParam := relationTovmmRsAcc.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovmmRsCtrlrPMonPol, ok := d.GetOk("relation_vmm_rs_ctrlr_p_mon_pol"); ok {
		relationParam := relationTovmmRsCtrlrPMonPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovmmRsMcastAddrNs, ok := d.GetOk("relation_vmm_rs_mcast_addr_ns"); ok {
		relationParam := relationTovmmRsMcastAddrNs.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovmmRsMgmtEPg, ok := d.GetOk("relation_vmm_rs_mgmt_e_pg"); ok {
		relationParam := relationTovmmRsMgmtEPg.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovmmRsToExtDevMgr, ok := d.GetOk("relation_vmm_rs_to_ext_dev_mgr"); ok {
		relationParamList := toStringList(relationTovmmRsToExtDevMgr.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTovmmRsVmmCtrlrP, ok := d.GetOk("relation_vmm_rs_vmm_ctrlr_p"); ok {
		relationParamList := toStringList(relationTovmmRsVmmCtrlrP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTovmmRsVxlanNs, ok := d.GetOk("relation_vmm_rs_vxlan_ns"); ok {
		relationParam := relationTovmmRsVxlanNs.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovmmRsVxlanNsDef, ok := d.GetOk("relation_vmm_rs_vxlan_ns_def"); ok {
		relationParam := relationTovmmRsVxlanNsDef.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovmmRsAcc, ok := d.GetOk("relation_vmm_rs_acc"); ok {
		relationParam := relationTovmmRsAcc.(string)
		err = aciClient.CreateRelationvmmRsAcc(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTovmmRsCtrlrPMonPol, ok := d.GetOk("relation_vmm_rs_ctrlr_p_mon_pol"); ok {
		relationParam := relationTovmmRsCtrlrPMonPol.(string)
		err = aciClient.CreateRelationvmmRsCtrlrPMonPol(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTovmmRsMcastAddrNs, ok := d.GetOk("relation_vmm_rs_mcast_addr_ns"); ok {
		relationParam := relationTovmmRsMcastAddrNs.(string)
		err = aciClient.CreateRelationvmmRsMcastAddrNs(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTovmmRsMgmtEPg, ok := d.GetOk("relation_vmm_rs_mgmt_e_pg"); ok {
		relationParam := relationTovmmRsMgmtEPg.(string)
		err = aciClient.CreateRelationvmmRsMgmtEPg(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTovmmRsToExtDevMgr, ok := d.GetOk("relation_vmm_rs_to_ext_dev_mgr"); ok {
		relationParamList := toStringList(relationTovmmRsToExtDevMgr.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationvmmRsToExtDevMgr(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTovmmRsVmmCtrlrP, ok := d.GetOk("relation_vmm_rs_vmm_ctrlr_p"); ok {
		relationParamList := relationTovmmRsVmmCtrlrP.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationvmmRsVmmCtrlrP(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, paramMap["epg_depl_pref"].(string), paramMap["target_dn"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTovmmRsVxlanNs, ok := d.GetOk("relation_vmm_rs_vxlan_ns"); ok {
		relationParam := relationTovmmRsVxlanNs.(string)
		err = aciClient.CreateRelationvmmRsVxlanNs(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTovmmRsVxlanNsDef, ok := d.GetOk("relation_vmm_rs_vxlan_ns_def"); ok {
		relationParam := relationTovmmRsVxlanNsDef.(string)
		err = aciClient.CreateRelationvmmRsVxlanNsDef(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vmmCtrlrP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciVMMControllerRead(ctx, d, m)
}

func resourceAciVMMControllerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VMMController: Beginning Update")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)
	vmmCtrlrPAttr := models.VMMControllerAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmCtrlrPAttr.Annotation = Annotation.(string)
	} else {
		vmmCtrlrPAttr.Annotation = "{}"
	}

	if DvsVersion, ok := d.GetOk("dvs_version"); ok {
		vmmCtrlrPAttr.DvsVersion = DvsVersion.(string)
	}

	if HostOrIp, ok := d.GetOk("host_or_ip"); ok {
		vmmCtrlrPAttr.HostOrIp = HostOrIp.(string)
	}

	if InventoryTrigSt, ok := d.GetOk("inventory_trig_st"); ok {
		vmmCtrlrPAttr.InventoryTrigSt = InventoryTrigSt.(string)
	}

	if Mode, ok := d.GetOk("mode"); ok {
		vmmCtrlrPAttr.Mode = Mode.(string)
	}

	if MsftConfigErrMsg, ok := d.GetOk("msft_config_err_msg"); ok {
		vmmCtrlrPAttr.MsftConfigErrMsg = MsftConfigErrMsg.(string)
	}
	if MsftConfigIssues, ok := d.GetOk("msft_config_issues"); ok {
		msftConfigIssuesList := make([]string, 0, 1)
		for _, val := range MsftConfigIssues.([]interface{}) {
			msftConfigIssuesList = append(msftConfigIssuesList, val.(string))
		}
		MsftConfigIssues := strings.Join(msftConfigIssuesList, ",")
		vmmCtrlrPAttr.MsftConfigIssues = MsftConfigIssues
	}

	if N1kvStatsMode, ok := d.GetOk("n1kv_stats_mode"); ok {
		vmmCtrlrPAttr.N1kvStatsMode = N1kvStatsMode.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vmmCtrlrPAttr.Name = Name.(string)
	}

	if Port, ok := d.GetOk("port"); ok {
		vmmCtrlrPAttr.Port = Port.(string)
	}

	if RootContName, ok := d.GetOk("root_cont_name"); ok {
		vmmCtrlrPAttr.RootContName = RootContName.(string)
	}

	if Scope, ok := d.GetOk("scope"); ok {
		vmmCtrlrPAttr.Scope = Scope.(string)
	}

	if SeqNum, ok := d.GetOk("seq_num"); ok {
		vmmCtrlrPAttr.SeqNum = SeqNum.(string)
	}

	if StatsMode, ok := d.GetOk("stats_mode"); ok {
		vmmCtrlrPAttr.StatsMode = StatsMode.(string)
	}

	if VxlanDeplPref, ok := d.GetOk("vxlan_depl_pref"); ok {
		vmmCtrlrPAttr.VxlanDeplPref = VxlanDeplPref.(string)
	}
	vmmCtrlrP := models.NewVMMController(fmt.Sprintf("ctrlr-%s", name), VMMDomainDn, nameAlias, vmmCtrlrPAttr)

	vmmCtrlrP.Status = "modified"
	err := aciClient.Save(vmmCtrlrP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vmm_rs_acc") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_acc")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vmm_rs_ctrlr_p_mon_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_ctrlr_p_mon_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vmm_rs_mcast_addr_ns") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_mcast_addr_ns")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vmm_rs_mgmt_e_pg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_mgmt_e_pg")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vmm_rs_to_ext_dev_mgr") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_vmm_rs_to_ext_dev_mgr")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_vmm_rs_vmm_ctrlr_p") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_vmm_rs_vmm_ctrlr_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_vmm_rs_vxlan_ns") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vxlan_ns")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vmm_rs_vxlan_ns_def") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vxlan_ns_def")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vmm_rs_acc") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_acc")
		err = aciClient.DeleteRelationvmmRsAcc(vmmCtrlrP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsAcc(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vmm_rs_ctrlr_p_mon_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_ctrlr_p_mon_pol")
		err = aciClient.DeleteRelationvmmRsCtrlrPMonPol(vmmCtrlrP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsCtrlrPMonPol(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vmm_rs_mcast_addr_ns") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_mcast_addr_ns")
		err = aciClient.DeleteRelationvmmRsMcastAddrNs(vmmCtrlrP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsMcastAddrNs(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vmm_rs_mgmt_e_pg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_mgmt_e_pg")
		err = aciClient.DeleteRelationvmmRsMgmtEPg(vmmCtrlrP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsMgmtEPg(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vmm_rs_to_ext_dev_mgr") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_vmm_rs_to_ext_dev_mgr")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationvmmRsToExtDevMgr(vmmCtrlrP.DistinguishedName, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationvmmRsToExtDevMgr(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_vmm_rs_vmm_ctrlr_p") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_vmm_rs_vmm_ctrlr_p")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationvmmRsVmmCtrlrP(vmmCtrlrP.DistinguishedName, paramMap["target_dn"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationvmmRsVmmCtrlrP(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, paramMap["epg_depl_pref"].(string), paramMap["target_dn"].(string))

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_vmm_rs_vxlan_ns") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vxlan_ns")
		err = aciClient.DeleteRelationvmmRsVxlanNs(vmmCtrlrP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsVxlanNs(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vmm_rs_vxlan_ns_def") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vxlan_ns_def")
		err = aciClient.DeleteRelationvmmRsVxlanNsDef(vmmCtrlrP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvmmRsVxlanNsDef(vmmCtrlrP.DistinguishedName, vmmCtrlrPAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vmmCtrlrP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciVMMControllerRead(ctx, d, m)
}

func resourceAciVMMControllerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vmmCtrlrP, err := getRemoteVMMController(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setVMMControllerAttributes(vmmCtrlrP, d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	vmmRsAccData, err := aciClient.ReadRelationvmmRsAcc(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsAcc %v", err)
		d.Set("relation_vmm_rs_acc", "")
	} else {
		d.Set("relation_vmm_rs_acc", vmmRsAccData.(string))
	}

	vmmRsCtrlrPMonPolData, err := aciClient.ReadRelationvmmRsCtrlrPMonPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsCtrlrPMonPol %v", err)
		d.Set("relation_vmm_rs_ctrlr_p_mon_pol", "")
	} else {
		d.Set("relation_vmm_rs_ctrlr_p_mon_pol", vmmRsCtrlrPMonPolData.(string))
	}

	vmmRsMcastAddrNsData, err := aciClient.ReadRelationvmmRsMcastAddrNs(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsMcastAddrNs %v", err)
		d.Set("relation_vmm_rs_mcast_addr_ns", "")
	} else {
		d.Set("relation_vmm_rs_mcast_addr_ns", vmmRsMcastAddrNsData.(string))
	}

	vmmRsMgmtEPgData, err := aciClient.ReadRelationvmmRsMgmtEPg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsMgmtEPg %v", err)
		d.Set("relation_vmm_rs_mgmt_e_pg", "")
	} else {
		d.Set("relation_vmm_rs_mgmt_e_pg", vmmRsMgmtEPgData.(string))
	}
	vmmRsToExtDevMgrData, err := aciClient.ReadRelationvmmRsToExtDevMgr(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsToExtDevMgr %v", err)
		d.Set("relation_vmm_rs_to_ext_dev_mgr", make([]string, 0, 1))
	} else {
		d.Set("relation_vmm_rs_to_ext_dev_mgr", toStringList(vmmRsToExtDevMgrData.(*schema.Set).List()))
	}

	vmmRsVmmCtrlrPData, err := aciClient.ReadRelationvmmRsVmmCtrlrP(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsVmmCtrlrP %v", err)
	} else {
		relParams := make([]map[string]string, 0, 1)
		relParamsList := vmmRsVmmCtrlrPData.([]map[string]string)
		for _, obj := range relParamsList {
			relParams = append(relParams, map[string]string{
				"epg_depl_pref": obj["epgDeplPref"],
				"target_dn":     obj["tDn"],
			})
		}
		d.Set("relation_vmm_rs_vmm_ctrlr_p", relParams)
	}

	vmmRsVxlanNsData, err := aciClient.ReadRelationvmmRsVxlanNs(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsVxlanNs %v", err)
		d.Set("relation_vmm_rs_vxlan_ns", "")
	} else {
		d.Set("relation_vmm_rs_vxlan_ns", vmmRsVxlanNsData.(string))
	}

	vmmRsVxlanNsDefData, err := aciClient.ReadRelationvmmRsVxlanNsDef(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vmmRsVxlanNsDef %v", err)
		d.Set("relation_vmm_rs_vxlan_ns_def", "")
	} else {
		d.Set("relation_vmm_rs_vxlan_ns_def", vmmRsVxlanNsDefData.(string))
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciVMMControllerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vmmCtrlrP")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

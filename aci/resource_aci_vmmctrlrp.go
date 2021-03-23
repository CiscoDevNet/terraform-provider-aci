package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciVMMController() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciVMMControllerCreate,
		Update: resourceAciVMMControllerUpdate,
		Read:   resourceAciVMMControllerRead,
		Delete: resourceAciVMMControllerDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVMMControllerImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"dvs_version": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "",

				ValidateFunc: validation.StringInSlice([]string{
					"5.1",
					"5.5",
					"6.0",
					"6.5",
					"6.6",
					"unmanaged",
				}, false),
			},

			"host_or_ip": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "host or ip",
			},

			"inventory_trig_st": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "",

				ValidateFunc: validation.StringInSlice([]string{
					"autoTriggered",
					"triggered",
					"untriggered",
				}, false),
			},

			"mode": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "mode of operation",

				ValidateFunc: validation.StringInSlice([]string{
					"cf",
					"default",
					"k8s",
					"n1kv",
					"openshift",
					"ovs",
					"rhev",
					"unknown",
				}, false),
			},

			"msft_config_err_msg": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"msft_config_issues": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "",

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

			"n1kv_stats_mode": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",

				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
					"unknown",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"port": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "service port number for LDAP service",
			},

			"root_cont_name": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "top level container name",
			},

			"scope": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "scope",

				ValidateFunc: validation.StringInSlice([]string{
					"MicrosoftSCVMM",
					"cloudfoundry",
					"iaas",
					"kubernetes",
					"network",
					"openshift",
					"openstack",
					"rhev",
					"unmanaged",
					"vm",
				}, false),
			},

			"seq_num": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "isis lsp sequence number",
			},

			"stats_mode": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "statistics mode",

				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
					"unknown",
				}, false),
			},

			"vxlan_depl_pref": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",

				ValidateFunc: validation.StringInSlice([]string{
					"nsx",
					"vxlan",
				}, false),
			},

			"relation_vmm_rs_mcast_addr_ns": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fvnsMcastAddrInstP",
			},
			"relation_vmm_rs_acc": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to vmmUsrAccP",
			},
			"relation_vmm_rs_vmm_ctrlr_p": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vmmCtrlrP",
				Set:         schema.HashString,
			},
			"relation_vmm_rs_vxlan_ns": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fvnsVxlanInstP",
			},
			"relation_vmm_rs_ctrlr_p_mon_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to monInfraPol",
			},
			"relation_vmm_rs_mgmt_e_pg": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fvEPg",
			},
		}),
	}
}

func getRemoteVMMController(client *client.Client, dn string) (*models.VMMController, error) {
	vmmCtrlrPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vmmCtrlrP := models.VMMControllerFromContainer(vmmCtrlrPCont)

	if vmmCtrlrP.DistinguishedName == "" {
		return nil, fmt.Errorf("VMM Controller %s not found", vmmCtrlrP.DistinguishedName)
	}

	return vmmCtrlrP, nil
}

func setVMMControllerAttributes(vmmCtrlrP *models.VMMController, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(vmmCtrlrP.DistinguishedName)
	d.Set("description", vmmCtrlrP.Description)
	// d.Set("vmm_domain_dn", GetParentDn(vmmCtrlrP.DistinguishedName))
	vmmCtrlrPMap, _ := vmmCtrlrP.ToMap()
	d.Set("name", vmmCtrlrPMap["name"])
	d.Set("vmm_domain_dn", GetParentDn(vmmCtrlrP.DistinguishedName, fmt.Sprintf("/ctrlr-%s", vmmCtrlrPMap["name"])))

	d.Set("annotation", vmmCtrlrPMap["annotation"])
	d.Set("dvs_version", vmmCtrlrPMap["dvsVersion"])
	d.Set("host_or_ip", vmmCtrlrPMap["hostOrIp"])
	d.Set("inventory_trig_st", vmmCtrlrPMap["inventoryTrigSt"])
	d.Set("mode", vmmCtrlrPMap["mode"])
	d.Set("msft_config_err_msg", vmmCtrlrPMap["msftConfigErrMsg"])
	d.Set("msft_config_issues", vmmCtrlrPMap["msftConfigIssues"])
	d.Set("n1kv_stats_mode", vmmCtrlrPMap["n1kvStatsMode"])
	d.Set("name_alias", vmmCtrlrPMap["nameAlias"])
	d.Set("port", vmmCtrlrPMap["port"])
	d.Set("root_cont_name", vmmCtrlrPMap["rootContName"])
	d.Set("scope", vmmCtrlrPMap["scope"])
	d.Set("seq_num", vmmCtrlrPMap["seqNum"])
	d.Set("stats_mode", vmmCtrlrPMap["statsMode"])
	d.Set("vxlan_depl_pref", vmmCtrlrPMap["vxlanDeplPref"])
	return d
}

func resourceAciVMMControllerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	vmmCtrlrP, err := getRemoteVMMController(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setVMMControllerAttributes(vmmCtrlrP, d)
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVMMControllerCreate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)

	vmmCtrlrPAttr := models.VMMControllerAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmCtrlrPAttr.Annotation = Annotation.(string)
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
		vmmCtrlrPAttr.MsftConfigIssues = MsftConfigIssues.(string)
	}
	if N1kvStatsMode, ok := d.GetOk("n1kv_stats_mode"); ok {
		vmmCtrlrPAttr.N1kvStatsMode = N1kvStatsMode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vmmCtrlrPAttr.NameAlias = NameAlias.(string)
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
	vmmCtrlrP := models.NewVMMController(fmt.Sprintf("ctrlr-%s", name), VMMDomainDn, desc, vmmCtrlrPAttr)

	err := aciClient.Save(vmmCtrlrP)
	if err != nil {
		return err
	}

	if relationTovmmRsMcastAddrNs, ok := d.GetOk("relation_vmm_rs_mcast_addr_ns"); ok {
		relationParam := relationTovmmRsMcastAddrNs.(string)
		err = aciClient.CreateRelationvmmRsMcastAddrNs(vmmCtrlrP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	if relationTovmmRsAcc, ok := d.GetOk("relation_vmm_rs_acc"); ok {
		relationParam := relationTovmmRsAcc.(string)
		err = aciClient.CreateRelationvmmRsAcc(vmmCtrlrP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	if relationTovmmRsVmmCtrlrP, ok := d.GetOk("relation_vmm_rs_vmm_ctrlr_p"); ok {
		relationParamList := toStringList(relationTovmmRsVmmCtrlrP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationvmmRsVmmCtrlrP(vmmCtrlrP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
		}
	}

	if relationTovmmRsVxlanNs, ok := d.GetOk("relation_vmm_rs_vxlan_ns"); ok {
		relationParam := relationTovmmRsVxlanNs.(string)
		err = aciClient.CreateRelationvmmRsVxlanNs(vmmCtrlrP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	if relationTovmmRsCtrlrPMonPol, ok := d.GetOk("relation_vmm_rs_ctrlr_p_mon_pol"); ok {
		relationParam := relationTovmmRsCtrlrPMonPol.(string)
		err = aciClient.CreateRelationvmmRsCtrlrPMonPol(vmmCtrlrP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	if relationTovmmRsMgmtEPg, ok := d.GetOk("relation_vmm_rs_mgmt_e_pg"); ok {
		relationParam := relationTovmmRsMgmtEPg.(string)
		err = aciClient.CreateRelationvmmRsMgmtEPg(vmmCtrlrP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	d.SetId(vmmCtrlrP.DistinguishedName)
	return resourceAciVMMControllerRead(d, m)
}

func resourceAciVMMControllerUpdate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)

	vmmCtrlrPAttr := models.VMMControllerAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmCtrlrPAttr.Annotation = Annotation.(string)
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
		vmmCtrlrPAttr.MsftConfigIssues = MsftConfigIssues.(string)
	}
	if N1kvStatsMode, ok := d.GetOk("n1kv_stats_mode"); ok {
		vmmCtrlrPAttr.N1kvStatsMode = N1kvStatsMode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vmmCtrlrPAttr.NameAlias = NameAlias.(string)
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
	vmmCtrlrP := models.NewVMMController(fmt.Sprintf("ctrlr-%s", name), VMMDomainDn, desc, vmmCtrlrPAttr)

	vmmCtrlrP.Status = "modified"

	err := aciClient.Save(vmmCtrlrP)

	if err != nil {
		return err
	}
	if d.HasChange("relation_vmm_rs_mcast_addr_ns") {
		_, newRelParam := d.GetChange("relation_vmm_rs_mcast_addr_ns")
		err = aciClient.DeleteRelationvmmRsMcastAddrNs(vmmCtrlrP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsMcastAddrNs(vmmCtrlrP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_vmm_rs_acc") {
		_, newRelParam := d.GetChange("relation_vmm_rs_acc")
		err = aciClient.DeleteRelationvmmRsAcc(vmmCtrlrP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsAcc(vmmCtrlrP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_vmm_rs_vmm_ctrlr_p") {
		oldRel, newRel := d.GetChange("relation_vmm_rs_vmm_ctrlr_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationvmmRsVmmCtrlrP(vmmCtrlrP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationvmmRsVmmCtrlrP(vmmCtrlrP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}
	}
	if d.HasChange("relation_vmm_rs_vxlan_ns") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vxlan_ns")
		err = aciClient.DeleteRelationvmmRsVxlanNs(vmmCtrlrP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsVxlanNs(vmmCtrlrP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_vmm_rs_ctrlr_p_mon_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_ctrlr_p_mon_pol")
		err = aciClient.DeleteRelationvmmRsCtrlrPMonPol(vmmCtrlrP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsCtrlrPMonPol(vmmCtrlrP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_vmm_rs_mgmt_e_pg") {
		_, newRelParam := d.GetChange("relation_vmm_rs_mgmt_e_pg")
		err = aciClient.DeleteRelationvmmRsMgmtEPg(vmmCtrlrP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsMgmtEPg(vmmCtrlrP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}

	d.SetId(vmmCtrlrP.DistinguishedName)
	return resourceAciVMMControllerRead(d, m)

}

func resourceAciVMMControllerRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	dn := d.Id()
	vmmCtrlrP, err := getRemoteVMMController(aciClient, dn)

	if err != nil {
		return err
	}
	setVMMControllerAttributes(vmmCtrlrP, d)
	return nil
}

func resourceAciVMMControllerDelete(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vmmCtrlrP")
	if err != nil {
		return err
	}

	d.SetId("")
	return err
}

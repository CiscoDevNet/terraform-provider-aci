package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciVSwitchPolicyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciVSwitchPolicyGroupCreate,
		Update: resourceAciVSwitchPolicyGroupUpdate,
		Read:   resourceAciVSwitchPolicyGroupRead,
		Delete: resourceAciVSwitchPolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVSwitchPolicyGroupImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"name_alias": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"relation_vmm_rs_vswitch_exporter_pol": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Create relation to netflowVmmExporterPol",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"exporter_pol_dn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"active_flow_timeout": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"idle_flow_timeout": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"sampling_rate": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"relation_vmm_rs_vswitch_override_fw_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to nwsFwPol",
			},
			"relation_vmm_rs_vswitch_override_stp_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to stpIfPol",
			},
			"relation_vmm_rs_vswitch_override_lldp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to lldpIfPol",
			},
			"relation_vmm_rs_vswitch_override_mcp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to mcpIfPol",
			},
			"relation_vmm_rs_vswitch_override_cdp_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to cdpIfPol",
			},
			"relation_vmm_rs_vswitch_override_lacp_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to lacpLagPol",
			},
			"relation_vmm_rs_vswitch_override_mtu_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to l2InstPol",
			},
		}),
	}
}

func getRemoteVSwitchPolicyGroup(client *client.Client, dn string) (*models.VSwitchPolicyGroup, error) {
	vmmVSwitchPolicyContCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vmmVSwitchPolicyCont := models.VSwitchPolicyGroupFromContainer(vmmVSwitchPolicyContCont)

	if vmmVSwitchPolicyCont.DistinguishedName == "" {
		return nil, fmt.Errorf("VMM vSwitch Policy %s not found", vmmVSwitchPolicyCont.DistinguishedName)
	}

	return vmmVSwitchPolicyCont, nil
}

func setVSwitchPolicyGroupAttributes(vmmVSwitchPolicyCont *models.VSwitchPolicyGroup, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(vmmVSwitchPolicyCont.DistinguishedName)
	d.Set("description", vmmVSwitchPolicyCont.Description)
	// d.Set("vmm_domain_dn", GetParentDn(vmmVSwitchPolicyCont.DistinguishedName))
	vmmVSwitchPolicyContMap, _ := vmmVSwitchPolicyCont.ToMap()
	// d.Set("name", vmmVSwitchPolicyContMap["name"])
	d.Set("vmm_domain_dn", GetParentDn(vmmVSwitchPolicyCont.DistinguishedName, "/vswitchpolcont"))

	d.Set("annotation", vmmVSwitchPolicyContMap["annotation"])
	d.Set("name_alias", vmmVSwitchPolicyContMap["nameAlias"])
	return d
}

func resourceAciVSwitchPolicyGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	vmmVSwitchPolicyCont, err := getRemoteVSwitchPolicyGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setVSwitchPolicyGroupAttributes(vmmVSwitchPolicyCont, d)
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVSwitchPolicyGroupCreate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)

	vmmVSwitchPolicyContAttr := models.VSwitchPolicyGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmVSwitchPolicyContAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vmmVSwitchPolicyContAttr.NameAlias = NameAlias.(string)
	}
	vmmVSwitchPolicyCont := models.NewVSwitchPolicyGroup(fmt.Sprintf("vswitchpolcont"), VMMDomainDn, desc, vmmVSwitchPolicyContAttr)

	err := aciClient.Save(vmmVSwitchPolicyCont)
	if err != nil {
		return err
	}

	exporterPolicyIDS := make([]string, 0, 1)
	if relationTovmmRsVswitchExporterPol, ok := d.GetOk("relation_vmm_rs_vswitch_exporter_pol"); ok {
		exporterPolicies := relationTovmmRsVswitchExporterPol.([]interface{})
		for _, relDn := range exporterPolicies {
			relation_vmm_rs_vswitch_exporter_pol := relDn.(map[string]interface{})
			var exporterPolicyDn string
			exporterPolicyDn, err = aciClient.CreateRelationvmmRsVswitchExporterPol(vmmVSwitchPolicyCont.DistinguishedName, relation_vmm_rs_vswitch_exporter_pol["exporter_pol_dn"].(string), relation_vmm_rs_vswitch_exporter_pol["active_flow_timeout"].(string), relation_vmm_rs_vswitch_exporter_pol["idle_flow_timeout"].(string), relation_vmm_rs_vswitch_exporter_pol["sampling_rate"].(string))

			if err != nil {
				return err
			}
			exporterPolicyIDS = append(exporterPolicyIDS, exporterPolicyDn)
		}
		log.Println("Check ... :", exporterPolicyIDS)
		d.Set("exporterPolicy_ids", exporterPolicyIDS)
	} else {
		d.Set("exporterPolicy_ids", exporterPolicyIDS)
	}

	if relationTovmmRsVswitchOverrideFwPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_fw_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideFwPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideFwPol(vmmVSwitchPolicyCont.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	if relationTovmmRsVswitchOverrideStpPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_stp_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideStpPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideStpPol(vmmVSwitchPolicyCont.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	if relationTovmmRsVswitchOverrideLldpIfPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_lldp_if_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideLldpIfPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideLldpIfPol(vmmVSwitchPolicyCont.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	if relationTovmmRsVswitchOverrideMcpIfPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_mcp_if_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideMcpIfPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideMcpIfPol(vmmVSwitchPolicyCont.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	if relationTovmmRsVswitchOverrideCdpIfPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_cdp_if_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideCdpIfPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideCdpIfPol(vmmVSwitchPolicyCont.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	if relationTovmmRsVswitchOverrideLacpPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_lacp_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideLacpPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideLacpPol(vmmVSwitchPolicyCont.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	if relationTovmmRsVswitchOverrideMtuPol, ok := d.GetOk("relation_vmm_rs_vswitch_override_mtu_pol"); ok {
		relationParam := relationTovmmRsVswitchOverrideMtuPol.(string)
		err = aciClient.CreateRelationvmmRsVswitchOverrideMtuPol(vmmVSwitchPolicyCont.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	d.SetId(vmmVSwitchPolicyCont.DistinguishedName)
	return resourceAciVSwitchPolicyGroupRead(d, m)
}

func resourceAciVSwitchPolicyGroupUpdate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	VMMDomainDn := d.Get("vmm_domain_dn").(string)

	vmmVSwitchPolicyContAttr := models.VSwitchPolicyGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmVSwitchPolicyContAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vmmVSwitchPolicyContAttr.NameAlias = NameAlias.(string)
	}
	vmmVSwitchPolicyCont := models.NewVSwitchPolicyGroup(fmt.Sprintf("vswitchpolcont"), VMMDomainDn, desc, vmmVSwitchPolicyContAttr)

	vmmVSwitchPolicyCont.Status = "modified"

	err := aciClient.Save(vmmVSwitchPolicyCont)

	if err != nil {
		return err
	}

	if d.HasChange("relation_vmm_rs_vswitch_exporter_pol") {
		relation_vmm_rs_vswitch_exporter_pol := d.Get("exporterPolicy_ids").([]interface{})
		for _, relDn := range relation_vmm_rs_vswitch_exporter_pol {
			err := aciClient.DeleteRelationvmmRsVswitchExporterPol(vmmVSwitchPolicyCont.DistinguishedName, relDn.(string))
			if err != nil {
				return err
			}
		}
		relationTovmmRsVswitchExporterPol := d.Get("relation_vmm_rs_vswitch_exporter_pol")
		exporterPolicyIDS := make([]string, 0, 1)
		exporterPolicies := relationTovmmRsVswitchExporterPol.([]interface{})
		for _, relDn := range exporterPolicies {
			relation_vmm_rs_vswitch_exporter_pol := relDn.(map[string]interface{})
			var exporterPolicyDn string
			exporterPolicyDn, err = aciClient.CreateRelationvmmRsVswitchExporterPol(vmmVSwitchPolicyCont.DistinguishedName, relation_vmm_rs_vswitch_exporter_pol["exporter_pol_dn"].(string), relation_vmm_rs_vswitch_exporter_pol["active_flow_timeout"].(string), relation_vmm_rs_vswitch_exporter_pol["idle_flow_timeout"].(string), relation_vmm_rs_vswitch_exporter_pol["sampling_rate"].(string))

			if err != nil {
				return err
			}

			exporterPolicyIDS = append(exporterPolicyIDS, exporterPolicyDn)
		}

		d.Set("exporterPolicy_ids", exporterPolicyIDS)

	}

	if d.HasChange("relation_vmm_rs_vswitch_override_fw_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_fw_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideFwPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideFwPol(vmmVSwitchPolicyCont.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_vmm_rs_vswitch_override_stp_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_stp_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideStpPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideStpPol(vmmVSwitchPolicyCont.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_vmm_rs_vswitch_override_lldp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_lldp_if_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideLldpIfPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideLldpIfPol(vmmVSwitchPolicyCont.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_vmm_rs_vswitch_override_mcp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_mcp_if_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideMcpIfPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideMcpIfPol(vmmVSwitchPolicyCont.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_vmm_rs_vswitch_override_cdp_if_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_cdp_if_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideCdpIfPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideCdpIfPol(vmmVSwitchPolicyCont.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_vmm_rs_vswitch_override_lacp_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_lacp_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideLacpPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideLacpPol(vmmVSwitchPolicyCont.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}

	if d.HasChange("relation_vmm_rs_vswitch_override_mtu_pol") {
		_, newRelParam := d.GetChange("relation_vmm_rs_vswitch_override_mtu_pol")
		err = aciClient.DeleteRelationvmmRsVswitchOverrideMtuPol(vmmVSwitchPolicyCont.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvmmRsVswitchOverrideMtuPol(vmmVSwitchPolicyCont.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}

	d.SetId(vmmVSwitchPolicyCont.DistinguishedName)
	return resourceAciVSwitchPolicyGroupRead(d, m)

}

func resourceAciVSwitchPolicyGroupRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	dn := d.Id()
	vmmVSwitchPolicyCont, err := getRemoteVSwitchPolicyGroup(aciClient, dn)

	if err != nil {
		return err
	}
	setVSwitchPolicyGroupAttributes(vmmVSwitchPolicyCont, d)
	return nil
}

func resourceAciVSwitchPolicyGroupDelete(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vmmVSwitchPolicyCont")
	if err != nil {
		return err
	}

	d.SetId("")
	return err
}

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

const (
	vlifplagpolattDnFromLogicalnterfaceProfileDn = "%s/rsdynPathAtt-[%s]/vlifplagpolatt"
)

func resourceAciVirtualLogicalInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciVirtualLogicalInterfaceProfileCreate,
		UpdateContext: resourceAciVirtualLogicalInterfaceProfileUpdate,
		ReadContext:   resourceAciVirtualLogicalInterfaceProfileRead,
		DeleteContext: resourceAciVirtualLogicalInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVirtualLogicalInterfaceProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"node_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"autostate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"encap_scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ctx",
					"local",
				}, false),
			},

			"if_inst_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ext-svi",
					"l3-port",
					"sub-interface",
					"unspecified",
				}, false),
			},

			"ipv6_dad": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"ll_addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"native",
					"regular",
					"untagged",
				}, false),
			},

			"mtu": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AF11",
					"AF12",
					"AF13",
					"AF21",
					"AF22",
					"AF23",
					"AF31",
					"AF32",
					"AF33",
					"AF41",
					"AF42",
					"AF43",
					"CS0",
					"CS1",
					"CS2",
					"CS3",
					"CS4",
					"CS5",
					"CS6",
					"CS7",
					"EF",
					"VA",
					"unspecified",
				}, false),
			},

			"relation_l3ext_rs_dyn_path_att": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tdn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"floating_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"forged_transmit": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Disabled",
								"Enabled",
							}, false),
							Default: "Disabled",
						},
						"mac_change": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Disabled",
								"Enabled",
							}, false),
							Default: "Disabled",
						},
						"promiscuous_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Disabled",
								"Enabled",
							}, false),
							Default: "Disabled",
						},
						"enhanced_lag_policy_dn": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"encap": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		}),
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
			configOld, configNew := diff.GetChange("relation_l3ext_rs_dyn_path_att")
			updatedConfig := compareConfigRsDynPathAttEncapBetweenUnknownEmpty(configOld.(*schema.Set).List(), configNew.(*schema.Set).List())
			diff.SetNew("relation_l3ext_rs_dyn_path_att", updatedConfig)
			return nil
		},
	}
}

func compareConfigRsDynPathAttEncapBetweenUnknownEmpty(old, new []interface{}) []interface{} {
	for _, oldItem := range old {
		oldMap := oldItem.(map[string]interface{})
		for _, newItem := range new {
			newMap := newItem.(map[string]interface{})
			if oldMap["tdn"] == newMap["tdn"] && oldMap["encap"] == "unknown" && newMap["encap"] == "" {
				newMap["encap"] = "unknown"
				break
			}
		}
	}
	return new
}

func getRemoteVirtualLogicalInterfaceProfile(client *client.Client, dn string) (*models.VirtualLogicalInterfaceProfile, error) {
	l3extVirtualLIfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extVirtualLIfP := models.VirtualLogicalInterfaceProfileFromContainer(l3extVirtualLIfPCont)

	if l3extVirtualLIfP.DistinguishedName == "" {
		return nil, fmt.Errorf("Logical Interface Profile %s not found", dn)
	}

	return l3extVirtualLIfP, nil
}

func setVirtualLogicalInterfaceProfileAttributes(l3extVirtualLIfP *models.VirtualLogicalInterfaceProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	d.SetId(l3extVirtualLIfP.DistinguishedName)
	d.Set("description", l3extVirtualLIfP.Description)
	if dn != l3extVirtualLIfP.DistinguishedName {
		d.Set("logical_interface_profile_dn", "")
	}
	l3extVirtualLIfPMap, err := l3extVirtualLIfP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("logical_interface_profile_dn", GetParentDn(dn, fmt.Sprintf("/vlifp-[%s]-[%s]", l3extVirtualLIfPMap["nodeDn"], l3extVirtualLIfPMap["encap"])))
	d.Set("node_dn", l3extVirtualLIfPMap["nodeDn"])
	d.Set("encap", l3extVirtualLIfPMap["encap"])
	d.Set("addr", l3extVirtualLIfPMap["addr"])
	d.Set("annotation", l3extVirtualLIfPMap["annotation"])
	d.Set("autostate", l3extVirtualLIfPMap["autostate"])
	d.Set("encap", l3extVirtualLIfPMap["encap"])
	d.Set("encap_scope", l3extVirtualLIfPMap["encapScope"])
	d.Set("if_inst_t", l3extVirtualLIfPMap["ifInstT"])
	d.Set("ipv6_dad", l3extVirtualLIfPMap["ipv6Dad"])
	d.Set("ll_addr", l3extVirtualLIfPMap["llAddr"])
	d.Set("mac", l3extVirtualLIfPMap["mac"])
	d.Set("mode", l3extVirtualLIfPMap["mode"])
	d.Set("mtu", l3extVirtualLIfPMap["mtu"])
	d.Set("target_dscp", l3extVirtualLIfPMap["targetDscp"])

	return d, nil
}

func getL3extRsVSwitchEnhancedLagPol(client *client.Client, dn, tDn string) string {
	l3extVirtualLIfPLagPolAttDn := fmt.Sprintf(vlifplagpolattDnFromLogicalnterfaceProfileDn, dn, tDn)
	l3extRsVSwitchEnhancedLagPol, err := client.ReadRelationl3extRsVSwitchEnhancedLagPol(l3extVirtualLIfPLagPolAttDn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsVSwitchEnhancedLagPol %v", l3extRsVSwitchEnhancedLagPol)
		return ""
	} else {
		return l3extRsVSwitchEnhancedLagPol.(string)
	}
}

func getAndSetL3extRsDynPathAttFromLogicalInterfaceProfile(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	l3extRsDynPathAttData, err := client.ReadRelationl3extRsDynPathAttFromLogicalInterfaceProfile(dn)
	l3extRsDynPaths := make([]map[string]string, 0)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsDynPathAtt %v", err)
		d.Set("relation_l3ext_rs_dyn_path_att", l3extRsDynPaths)
		return nil, err
	} else {
		l3extRsDynPathAttMap := l3extRsDynPathAttData.([]map[string]string)
		for _, l3extRsDynPathObj := range l3extRsDynPathAttMap {
			obj := make(map[string]string, 0)
			obj["tdn"] = l3extRsDynPathObj["tDn"]
			obj["floating_address"] = l3extRsDynPathObj["floatingAddr"]
			obj["forged_transmit"] = l3extRsDynPathObj["forgedTransmit"]
			obj["mac_change"] = l3extRsDynPathObj["macChange"]
			obj["promiscuous_mode"] = l3extRsDynPathObj["promMode"]
			obj["encap"] = l3extRsDynPathObj["encap"]
			obj["enhanced_lag_policy_dn"] = getL3extRsVSwitchEnhancedLagPol(client, dn, l3extRsDynPathObj["tDn"])
			l3extRsDynPaths = append(l3extRsDynPaths, obj)
		}
		d.Set("relation_l3ext_rs_dyn_path_att", l3extRsDynPaths)
		log.Printf("[DEBUG]: l3extRsDynPathAtt: %v finished successfully", l3extRsDynPaths)
	}
	return d, nil
}

func resourceAciVirtualLogicalInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extVirtualLIfP, err := getRemoteVirtualLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setVirtualLogicalInterfaceProfileAttributes(l3extVirtualLIfP, d)
	if err != nil {
		return nil, err
	}

	getAndSetL3extRsDynPathAttFromLogicalInterfaceProfile(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVirtualLogicalInterfaceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Logical Interface Profile: Beginning Creation")

	aciClient := m.(*client.Client)

	desc := d.Get("description").(string)

	nodeDn := d.Get("node_dn").(string)

	encap := d.Get("encap").(string)

	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	l3extVirtualLIfPAttr := models.VirtualLogicalInterfaceProfileAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extVirtualLIfPAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extVirtualLIfPAttr.Annotation = Annotation.(string)
	} else {
		l3extVirtualLIfPAttr.Annotation = "{}"
	}
	if Autostate, ok := d.GetOk("autostate"); ok {
		l3extVirtualLIfPAttr.Autostate = Autostate.(string)
	}
	if Encap, ok := d.GetOk("encap"); ok {
		l3extVirtualLIfPAttr.Encap = Encap.(string)
	}
	if EncapScope, ok := d.GetOk("encap_scope"); ok {
		l3extVirtualLIfPAttr.EncapScope = EncapScope.(string)
	}
	if IfInstT, ok := d.GetOk("if_inst_t"); ok {
		l3extVirtualLIfPAttr.IfInstT = IfInstT.(string)
	}
	if Ipv6Dad, ok := d.GetOk("ipv6_dad"); ok {
		l3extVirtualLIfPAttr.Ipv6Dad = Ipv6Dad.(string)
	}
	if LlAddr, ok := d.GetOk("ll_addr"); ok {
		l3extVirtualLIfPAttr.LlAddr = LlAddr.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		l3extVirtualLIfPAttr.Mac = Mac.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		l3extVirtualLIfPAttr.Mode = Mode.(string)
	}
	if Mtu, ok := d.GetOk("mtu"); ok {
		l3extVirtualLIfPAttr.Mtu = Mtu.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extVirtualLIfPAttr.TargetDscp = TargetDscp.(string)
	}
	l3extVirtualLIfP := models.NewVirtualLogicalInterfaceProfile(fmt.Sprintf("vlifp-[%s]-[%s]", nodeDn, encap), LogicalInterfaceProfileDn, desc, l3extVirtualLIfPAttr)

	err := aciClient.Save(l3extVirtualLIfP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTol3extRsDynPathAtt, ok := d.GetOk("relation_l3ext_rs_dyn_path_att"); ok {
		relationParamList := relationTol3extRsDynPathAtt.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			checkDns = append(checkDns, paramMap["tdn"].(string))
		}
	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationTol3extRsDynPathAtt, ok := d.GetOk("relation_l3ext_rs_dyn_path_att"); ok {
		relationParamList := relationTol3extRsDynPathAtt.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsDynPathAttFromLogicalInterfaceProfile(l3extVirtualLIfP.DistinguishedName, paramMap["tdn"].(string), paramMap["floating_address"].(string), paramMap["forged_transmit"].(string), paramMap["mac_change"].(string), paramMap["promiscuous_mode"].(string), paramMap["encap"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
			if paramMap["enhanced_lag_policy_dn"].(string) != "" {
				l3extVirtualLIfPLagPolAttDn := fmt.Sprintf(vlifplagpolattDnFromLogicalnterfaceProfileDn, l3extVirtualLIfP.DistinguishedName, paramMap["tdn"].(string))
				err = aciClient.CreateRelationl3extRsVSwitchEnhancedLagPol(l3extVirtualLIfPLagPolAttDn, paramMap["enhanced_lag_policy_dn"].(string))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	d.SetId(l3extVirtualLIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciVirtualLogicalInterfaceProfileRead(ctx, d, m)
}

func resourceAciVirtualLogicalInterfaceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Logical Interface Profile: Beginning Update")

	aciClient := m.(*client.Client)

	desc := d.Get("description").(string)

	nodeDn := d.Get("node_dn").(string)

	encap := d.Get("encap").(string)

	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	l3extVirtualLIfPAttr := models.VirtualLogicalInterfaceProfileAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extVirtualLIfPAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extVirtualLIfPAttr.Annotation = Annotation.(string)
	} else {
		l3extVirtualLIfPAttr.Annotation = "{}"
	}
	if Autostate, ok := d.GetOk("autostate"); ok {
		l3extVirtualLIfPAttr.Autostate = Autostate.(string)
	}
	if Encap, ok := d.GetOk("encap"); ok {
		l3extVirtualLIfPAttr.Encap = Encap.(string)
	}
	if EncapScope, ok := d.GetOk("encap_scope"); ok {
		l3extVirtualLIfPAttr.EncapScope = EncapScope.(string)
	}
	if IfInstT, ok := d.GetOk("if_inst_t"); ok {
		l3extVirtualLIfPAttr.IfInstT = IfInstT.(string)
	}
	if Ipv6Dad, ok := d.GetOk("ipv6_dad"); ok {
		l3extVirtualLIfPAttr.Ipv6Dad = Ipv6Dad.(string)
	}
	if LlAddr, ok := d.GetOk("ll_addr"); ok {
		l3extVirtualLIfPAttr.LlAddr = LlAddr.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		l3extVirtualLIfPAttr.Mac = Mac.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		l3extVirtualLIfPAttr.Mode = Mode.(string)
	}
	if Mtu, ok := d.GetOk("mtu"); ok {
		l3extVirtualLIfPAttr.Mtu = Mtu.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extVirtualLIfPAttr.TargetDscp = TargetDscp.(string)
	}
	l3extVirtualLIfP := models.NewVirtualLogicalInterfaceProfile(fmt.Sprintf("vlifp-[%s]-[%s]", nodeDn, encap), LogicalInterfaceProfileDn, desc, l3extVirtualLIfPAttr)

	l3extVirtualLIfP.Status = "modified"

	err := aciClient.Save(l3extVirtualLIfP)

	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTol3extRsDynPathAtt, ok := d.GetOk("relation_l3ext_rs_dyn_path_att"); ok {
		relationParamList := relationTol3extRsDynPathAtt.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			checkDns = append(checkDns, paramMap["tdn"].(string))
		}
	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("relation_l3ext_rs_dyn_path_att") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_dyn_path_att")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationl3extRsDynPathAttFromLogicalInterfaceProfile(l3extVirtualLIfP.DistinguishedName, paramMap["tdn"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationl3extRsDynPathAttFromLogicalInterfaceProfile(l3extVirtualLIfP.DistinguishedName, paramMap["tdn"].(string), paramMap["floating_address"].(string), paramMap["forged_transmit"].(string), paramMap["mac_change"].(string), paramMap["promiscuous_mode"].(string), paramMap["encap"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
			if paramMap["enhanced_lag_policy_dn"].(string) != "" {
				l3extVirtualLIfPLagPolAttDn := fmt.Sprintf(vlifplagpolattDnFromLogicalnterfaceProfileDn, l3extVirtualLIfP.DistinguishedName, paramMap["tdn"].(string))
				err = aciClient.CreateRelationl3extRsVSwitchEnhancedLagPol(l3extVirtualLIfPLagPolAttDn, paramMap["enhanced_lag_policy_dn"].(string))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

	}

	d.SetId(l3extVirtualLIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciVirtualLogicalInterfaceProfileRead(ctx, d, m)

}

func resourceAciVirtualLogicalInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extVirtualLIfP, err := getRemoteVirtualLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setVirtualLogicalInterfaceProfileAttributes(l3extVirtualLIfP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	getAndSetL3extRsDynPathAttFromLogicalInterfaceProfile(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciVirtualLogicalInterfaceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extVirtualLIfP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

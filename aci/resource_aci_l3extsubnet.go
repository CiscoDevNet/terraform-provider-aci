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

func resourceAciL3ExtSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL3ExtSubnetCreate,
		UpdateContext: resourceAciL3ExtSubnetUpdate,
		ReadContext:   resourceAciL3ExtSubnetRead,
		DeleteContext: resourceAciL3ExtSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3ExtSubnetImport,
		},

		SchemaVersion: 1,

		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
			// workaround because ValidateFunc is not allowed on TypeSet
			if diff.HasChange("relation_l3ext_rs_subnet_to_profile") {
				_, new := diff.GetChange("relation_l3ext_rs_subnet_to_profile")
				// validate that a direction type ( import/export ) is not defined more than once
				err := validateDirection(new.(*schema.Set).List())
				if err != nil {
					return err
				}
				// validate that dn and name are not both defined
				err2 := validateDnAndName(new.(*schema.Set).List())
				if err2 != nil {
					return err2
				}
			}
			return nil
		},

		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"external_network_instance_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"aggregate": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressBitMaskDiffFunc(),
				ValidateFunc: schema.SchemaValidateFunc(validateCommaSeparatedStringInSlice([]string{
					"import-rtctrl",
					"export-rtctrl",
					"shared-rtctrl",
					"none",
				}, false, "none")),
			},

			"scope": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"import-rtctrl",
						"export-rtctrl",
						"shared-rtctrl",
						"import-security",
						"shared-security",
					}, false),
				},
			},

			"relation_l3ext_rs_subnet_to_profile": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tn_rtctrl_profile_name": {
							Type:       schema.TypeString,
							Optional:   true,
							Default:    "", // workaround to mimic computed behaviour, which does not work in TypeSet
							Deprecated: "use tn_rtctrl_profile_dn instead",
						},
						"tn_rtctrl_profile_dn": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "", // workaround to mimic computed behaviour, which does not work in TypeSet
						},
						"direction": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"export",
								"import",
							}, false),
						},
					},
				},
			},
			"relation_l3ext_rs_subnet_to_rt_summ": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		})),
	}
}

func validateDirection(relationParamList []interface{}) error {
	err := fmt.Errorf("duplicate directions not allowed in relation_l3ext_rs_subnet_to_profile")
	if len(relationParamList) > 2 {
		return err
	}
	validateDirection := ""
	for _, relationParam := range relationParamList {
		paramMap := relationParam.(map[string]interface{})
		if paramMap["direction"].(string) != validateDirection {
			validateDirection = paramMap["direction"].(string)
		} else {
			return err
		}
	}
	return nil
}

func validateDnAndName(relationParamList []interface{}) error {
	for _, relationParam := range relationParamList {
		paramMap := relationParam.(map[string]interface{})
		if paramMap["tn_rtctrl_profile_dn"] != "" && paramMap["tn_rtctrl_profile_name"] != "" {
			return fmt.Errorf("Usage of both tn_rtctrl_profile_dn and tn_rtctrl_profile_name parameters is not supported. tn_rtctrl_profile_name parameter will be deprecated use tn_rtctrl_profile_dn instead.")
		}
		if paramMap["tn_rtctrl_profile_dn"] == "" && paramMap["tn_rtctrl_profile_name"] == "" {
			return fmt.Errorf("tn_rtctrl_profile_dn is required to generate Route Control Profile")
		}
	}
	return nil
}

func getTnRtctrlProfileName(paramMap map[string]interface{}) string {
	if paramMap["tn_rtctrl_profile_dn"] != "" {
		return GetMOName(paramMap["tn_rtctrl_profile_dn"].(string))
	} else {
		return paramMap["tn_rtctrl_profile_name"].(string)
	}
}

func getRemoteL3ExtSubnet(client *client.Client, dn string) (*models.L3ExtSubnet, error) {
	l3extSubnetCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extSubnet := models.L3ExtSubnetFromContainer(l3extSubnetCont)

	if l3extSubnet.DistinguishedName == "" {
		return nil, fmt.Errorf("Subnet %s not found", dn)
	}

	return l3extSubnet, nil
}

func setL3ExtSubnetAttributes(l3extSubnet *models.L3ExtSubnet, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(l3extSubnet.DistinguishedName)
	d.Set("description", l3extSubnet.Description)

	if dn != l3extSubnet.DistinguishedName {
		d.Set("external_network_instance_profile_dn", "")
	}
	l3extSubnetMap, err := l3extSubnet.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("external_network_instance_profile_dn", GetParentDn(dn, fmt.Sprintf("/extsubnet-[%s]", l3extSubnetMap["ip"])))
	d.Set("ip", l3extSubnetMap["ip"])

	if l3extSubnetMap["aggregate"] == "" {
		d.Set("aggregate", "none")
	} else {
		d.Set("aggregate", l3extSubnetMap["aggregate"])
	}

	d.Set("annotation", l3extSubnetMap["annotation"])
	d.Set("ip", l3extSubnetMap["ip"])
	d.Set("name_alias", l3extSubnetMap["nameAlias"])

	scpGet := make([]string, 0, 1)
	for _, val := range strings.Split(l3extSubnetMap["scope"], ",") {
		scpGet = append(scpGet, strings.Trim(val, " "))
	}
	sort.Strings(scpGet)
	if scpInp, ok := d.GetOk("scope"); ok {
		scpAct := make([]string, 0, 1)
		for _, val := range scpInp.([]interface{}) {
			scpAct = append(scpAct, val.(string))
		}
		sort.Strings(scpAct)
		if reflect.DeepEqual(scpAct, scpGet) {
			d.Set("scope", d.Get("scope").([]interface{}))
		} else {
			d.Set("scope", scpGet)
		}
	} else {
		d.Set("scope", scpGet)
	}

	return d, nil
}

func getAndSetl3extRsSubnetToProfileFromL3ExtSubnet(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	l3extRsSubnetToProfileData, err := client.ReadRelationl3extRsSubnetToProfileFromL3ExtSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsSubnetToProfile %v", err)
		d.Set("relation_l3ext_rs_subnet_to_profile", make([]map[string]string, 0, 1))
		return nil, err
	} else {
		relParamList := make([]map[string]string, 0, 1)
		relParams := l3extRsSubnetToProfileData.([]map[string]string)
		for _, obj := range relParams {
			relParamList = append(relParamList, map[string]string{
				// obj["tnRtctrlProfileName"] is set to tDN in aci-go-client thus name is assigned to tn_rtctrl_profile_dn
				"tn_rtctrl_profile_dn": obj["tnRtctrlProfileName"],
				"direction":            obj["direction"],
			})
		}
		d.Set("relation_l3ext_rs_subnet_to_profile", relParamList)
		log.Printf("[DEBUG]: l3extRsSubnetToProfileData: %v finished successfully", relParamList)
	}
	return d, nil
}

func getAndSetl3extRsSubnetToRtSummFromL3ExtSubnet(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	l3extRsSubnetToRtSummData, err := client.ReadRelationl3extRsSubnetToRtSummFromL3ExtSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsSubnetToRtSumm %v", err)
		d.Set("relation_l3ext_rs_subnet_to_rt_summ", "")
		return nil, err
	} else {
		d.Set("relation_l3ext_rs_subnet_to_rt_summ", l3extRsSubnetToRtSummData.(string))
		log.Printf("[DEBUG]: l3extRsSubnetToRtSummData: %v finished successfully", l3extRsSubnetToRtSummData)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return d, nil
}

func resourceAciL3ExtSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extSubnet, err := getRemoteL3ExtSubnet(aciClient, dn)

	if err != nil {
		return nil, err
	}
	l3extSubnetMap, err := l3extSubnet.ToMap()
	if err != nil {
		return nil, err
	}
	ip := l3extSubnetMap["ip"]
	pDN := GetParentDn(dn, fmt.Sprintf("/extsubnet-[%s]", ip))
	d.Set("external_network_instance_profile_dn", pDN)
	schemaFilled, err := setL3ExtSubnetAttributes(l3extSubnet, d)
	if err != nil {
		return nil, err
	}

	getAndSetl3extRsSubnetToProfileFromL3ExtSubnet(aciClient, dn, d)
	getAndSetl3extRsSubnetToRtSummFromL3ExtSubnet(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3ExtSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Subnet: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	ExternalNetworkInstanceProfileDn := d.Get("external_network_instance_profile_dn").(string)

	l3extSubnetAttr := models.L3ExtSubnetAttributes{}
	if Aggregate, ok := d.GetOk("aggregate"); ok {
		agg := Aggregate.(string)
		if agg == "none" {
			agg = ""
		}
		l3extSubnetAttr.Aggregate = agg
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extSubnetAttr.Annotation = Annotation.(string)
	} else {
		l3extSubnetAttr.Annotation = "{}"
	}
	if Ip, ok := d.GetOk("ip"); ok {
		l3extSubnetAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extSubnetAttr.NameAlias = NameAlias.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		scpList := make([]string, 0, 1)
		for _, val := range Scope.([]interface{}) {
			scpList = append(scpList, val.(string))
		}
		scp := strings.Join(scpList, ",")
		l3extSubnetAttr.Scope = scp
	}
	l3extSubnet := models.NewL3ExtSubnet(fmt.Sprintf("extsubnet-[%s]", ip), ExternalNetworkInstanceProfileDn, desc, l3extSubnetAttr)

	err := aciClient.Save(l3extSubnet)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTol3extRsSubnetToRtSumm, ok := d.GetOk("relation_l3ext_rs_subnet_to_rt_summ"); ok {
		relationParam := relationTol3extRsSubnetToRtSumm.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTol3extRsSubnetToProfile, ok := d.GetOk("relation_l3ext_rs_subnet_to_profile"); ok {

		relationParamList := relationTol3extRsSubnetToProfile.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			relationParamName := getTnRtctrlProfileName(paramMap)
			err = aciClient.CreateRelationl3extRsSubnetToProfileFromL3ExtSubnet(l3extSubnet.DistinguishedName, relationParamName, paramMap["direction"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTol3extRsSubnetToRtSumm, ok := d.GetOk("relation_l3ext_rs_subnet_to_rt_summ"); ok {
		relationParam := relationTol3extRsSubnetToRtSumm.(string)
		err = aciClient.CreateRelationl3extRsSubnetToRtSummFromL3ExtSubnet(l3extSubnet.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(l3extSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3ExtSubnetRead(ctx, d, m)
}

func resourceAciL3ExtSubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Subnet: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	ExternalNetworkInstanceProfileDn := d.Get("external_network_instance_profile_dn").(string)

	l3extSubnetAttr := models.L3ExtSubnetAttributes{}
	if Aggregate, ok := d.GetOk("aggregate"); ok {
		agg := Aggregate.(string)
		if agg == "none" {
			agg = ""
		}
		l3extSubnetAttr.Aggregate = agg
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extSubnetAttr.Annotation = Annotation.(string)
	} else {
		l3extSubnetAttr.Annotation = "{}"
	}
	if Ip, ok := d.GetOk("ip"); ok {
		l3extSubnetAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extSubnetAttr.NameAlias = NameAlias.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		scpList := make([]string, 0, 1)
		for _, val := range Scope.([]interface{}) {
			scpList = append(scpList, val.(string))
		}
		scp := strings.Join(scpList, ",")
		l3extSubnetAttr.Scope = scp
	}
	l3extSubnet := models.NewL3ExtSubnet(fmt.Sprintf("extsubnet-[%s]", ip), ExternalNetworkInstanceProfileDn, desc, l3extSubnetAttr)

	l3extSubnet.Status = "modified"

	err := aciClient.Save(l3extSubnet)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_l3ext_rs_subnet_to_rt_summ") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_subnet_to_rt_summ")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_l3ext_rs_subnet_to_profile") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_subnet_to_profile")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			relationParamName := getTnRtctrlProfileName(paramMap)
			err = aciClient.DeleteRelationl3extRsSubnetToProfileFromL3ExtSubnet(l3extSubnet.DistinguishedName, relationParamName, paramMap["direction"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			relationParamName := getTnRtctrlProfileName(paramMap)
			err = aciClient.CreateRelationl3extRsSubnetToProfileFromL3ExtSubnet(l3extSubnet.DistinguishedName, relationParamName, paramMap["direction"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_l3ext_rs_subnet_to_rt_summ") {
		_, newRelParam := d.GetChange("relation_l3ext_rs_subnet_to_rt_summ")
		err = aciClient.DeleteRelationl3extRsSubnetToRtSummFromL3ExtSubnet(l3extSubnet.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationl3extRsSubnetToRtSummFromL3ExtSubnet(l3extSubnet.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(l3extSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3ExtSubnetRead(ctx, d, m)

}

func resourceAciL3ExtSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extSubnet, err := getRemoteL3ExtSubnet(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setL3ExtSubnetAttributes(l3extSubnet, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	getAndSetl3extRsSubnetToProfileFromL3ExtSubnet(aciClient, dn, d)
	getAndSetl3extRsSubnetToRtSummFromL3ExtSubnet(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3ExtSubnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extSubnet")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

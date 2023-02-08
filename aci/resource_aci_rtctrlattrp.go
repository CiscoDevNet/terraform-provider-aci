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

func resourceAciActionRuleProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciActionRuleProfileCreate,
		UpdateContext: resourceAciActionRuleProfileUpdate,
		ReadContext:   resourceAciActionRuleProfileRead,
		DeleteContext: resourceAciActionRuleProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciActionRuleProfileImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"set_route_tag": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateIntBetweenFromString(0, 2147483647),
				ConflictsWith: []string{"multipath", "next_hop_propagation"},
				Computed:      true,
				// Set nexthop unchanged action cannot be configured along with set route tag action under the set action rule profile.
			},
			"set_preference": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIntBetweenFromString(0, 2147483647),
				Computed:     true,
			},
			"set_weight": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIntBetweenFromString(0, 2147483647),
				Computed:     true,
			},
			"set_metric": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIntBetweenFromString(0, 2147483647),
				Computed:     true,
			},
			"set_metric_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ospf-type1",
					"ospf-type2",
				}, false),
				Computed: true,
			},
			"set_next_hop": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"set_communities": {
				Optional: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"next_hop_propagation": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
				ConflictsWith: []string{"set_route_tag"},
				Computed:      true,
				// Set nexthop unchanged action cannot be configured along with set route tag action under the set action rule profile.
			},
			"multipath": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
				ConflictsWith: []string{"set_route_tag"},
				Computed:      true,
				// Set nexthop unchanged action cannot be configured along with set route tag action under the set action rule profile.
			},
			"set_as_path_prepend_last_as": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIntBetweenFromString(1, 10),
				Computed:     true,
			},
			"set_as_path_prepend_as": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIntBetweenFromString(0, 2147483647),
							Description:  "ASN must be between 0 and 2147483647",
						},
						"order": {
							Required:     true,
							Type:         schema.TypeString,
							ValidateFunc: validateIntBetweenFromString(0, 31),
							Description:  "Order must be between 0 and 31",
						},
					},
				},
				Computed: true,
			},
			"set_dampening": {
				Optional: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		})),
	}
}

func getRemoteActionRuleProfile(client *client.Client, dn string) (*models.ActionRuleProfile, error) {
	rtctrlAttrPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	rtctrlAttrP := models.ActionRuleProfileFromContainer(rtctrlAttrPCont)

	if rtctrlAttrP.DistinguishedName == "" {
		return nil, fmt.Errorf("Action Rule Profile %s not found", dn)
	}

	return rtctrlAttrP, nil
}

func setActionRuleProfileAttributes(rtctrlAttrP *models.ActionRuleProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(rtctrlAttrP.DistinguishedName)
	d.Set("description", rtctrlAttrP.Description)
	if dn != rtctrlAttrP.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	rtctrlAttrPMap, err := rtctrlAttrP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", rtctrlAttrPMap["name"])
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/attr-%s", rtctrlAttrPMap["name"])))
	d.Set("annotation", rtctrlAttrPMap["annotation"])
	d.Set("name_alias", rtctrlAttrPMap["nameAlias"])
	return d, nil
}

func getRemoteRtctrlSetTag(client *client.Client, dn string) (*models.RtctrlSetTag, error) {
	rtctrlSetTagCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetTag := models.RtctrlSetTagFromContainer(rtctrlSetTagCont)
	if rtctrlSetTag.DistinguishedName == "" {
		return nil, fmt.Errorf("RtctrlSetTag %s not found", dn)
	}
	return rtctrlSetTag, nil
}

func setRtctrlSetTagAttributes(rtctrlSetTag *models.RtctrlSetTag, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetTagMap, err := rtctrlSetTag.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("set_route_tag", rtctrlSetTagMap["tag"])
	return d, nil
}

func getRemoteRtctrlSetPref(client *client.Client, dn string) (*models.RtctrlSetPref, error) {
	rtctrlSetPrefCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetPref := models.RtctrlSetPrefFromContainer(rtctrlSetPrefCont)
	if rtctrlSetPref.DistinguishedName == "" {
		return nil, fmt.Errorf("rtctrlSetPref %s not found", dn)
	}
	return rtctrlSetPref, nil
}

func setRtctrlSetPrefAttributes(rtctrlSetPref *models.RtctrlSetPref, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetPrefMap, err := rtctrlSetPref.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("set_preference", rtctrlSetPrefMap["localPref"])
	return d, nil
}

func getRemoteRtctrlSetWeight(client *client.Client, dn string) (*models.RtctrlSetWeight, error) {
	rtctrlSetWeightCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetWeight := models.RtctrlSetWeightFromContainer(rtctrlSetWeightCont)
	if rtctrlSetWeight.DistinguishedName == "" {
		return nil, fmt.Errorf("rtctrlSetWeight %s not found", dn)
	}
	return rtctrlSetWeight, nil
}

func setRtctrlSetWeightAttributes(rtctrlSetWeight *models.RtctrlSetWeight, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetWeightMap, err := rtctrlSetWeight.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("set_weight", rtctrlSetWeightMap["weight"])
	return d, nil
}

func getRemoteRtctrlSetRtMetric(client *client.Client, dn string) (*models.RtctrlSetRtMetric, error) {
	rtctrlSetRtMetricCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetRtMetric := models.RtctrlSetRtMetricFromContainer(rtctrlSetRtMetricCont)
	if rtctrlSetRtMetric.DistinguishedName == "" {
		return nil, fmt.Errorf("Route Control Set Route Metric %s not found", dn)
	}
	return rtctrlSetRtMetric, nil
}

func setRtctrlSetRtMetricAttributes(rtctrlSetRtMetric *models.RtctrlSetRtMetric, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetRtMetricMap, err := rtctrlSetRtMetric.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("set_metric", rtctrlSetRtMetricMap["metric"])
	return d, nil
}

func getRemoteRtctrlSetRtMetricType(client *client.Client, dn string) (*models.RtctrlSetRtMetricType, error) {
	rtctrlSetRtMetricTypeCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetRtMetricType := models.RtctrlSetRtMetricTypeFromContainer(rtctrlSetRtMetricTypeCont)
	if rtctrlSetRtMetricType.DistinguishedName == "" {
		return nil, fmt.Errorf("Route Control Set Route Metric Type %s not found", dn)
	}
	return rtctrlSetRtMetricType, nil
}

func setRtctrlSetRtMetricTypeAttributes(rtctrlSetRtMetricType *models.RtctrlSetRtMetricType, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetRtMetricTypeMap, err := rtctrlSetRtMetricType.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("set_metric_type", rtctrlSetRtMetricTypeMap["metricType"])
	return d, nil
}

func getRemoteRtctrlSetNh(client *client.Client, dn string) (*models.RtctrlSetNh, error) {
	rtctrlSetNhCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetNh := models.RtctrlSetNhFromContainer(rtctrlSetNhCont)
	if rtctrlSetNh.DistinguishedName == "" {
		return nil, fmt.Errorf("Route Control Set Nexthop %s not found", dn)
	}
	return rtctrlSetNh, nil
}

func setRtctrlSetNhAttributes(rtctrlSetNh *models.RtctrlSetNh, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetNhMap, err := rtctrlSetNh.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("set_next_hop", rtctrlSetNhMap["addr"])
	return d, nil
}

func getRemoteRtctrlSetComm(client *client.Client, dn string) (*models.RtctrlSetComm, error) {
	rtctrlSetCommCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetComm := models.RtctrlSetCommFromContainer(rtctrlSetCommCont)
	if rtctrlSetComm.DistinguishedName == "" {
		return nil, fmt.Errorf("Route Control Set Community %s not found", dn)
	}
	return rtctrlSetComm, nil
}

func setRtctrlSetCommAttributes(rtctrlSetComm *models.RtctrlSetComm, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetCommMap, err := rtctrlSetComm.ToMap()
	if err != nil {
		return d, err
	}

	newContent := make(map[string]interface{})
	newContent["community"] = rtctrlSetCommMap["community"]
	newContent["criteria"] = rtctrlSetCommMap["setCriteria"]
	d.Set("set_communities", newContent)

	return d, nil
}

func getRemoteNexthopUnchangedAction(client *client.Client, dn string) (*models.NexthopUnchangedAction, error) {
	rtctrlSetNhUnchangedCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetNhUnchanged := models.NexthopUnchangedActionFromContainer(rtctrlSetNhUnchangedCont)
	if rtctrlSetNhUnchanged.DistinguishedName == "" {
		return nil, fmt.Errorf("Next hop Unchanged Action %s not found", dn)
	}
	return rtctrlSetNhUnchanged, nil
}

func setNexthopUnchangedActionAttributes(rtctrlSetNhUnchanged *models.NexthopUnchangedAction, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetNhUnchangedMap, err := rtctrlSetNhUnchanged.ToMap()
	if err != nil {
		return d, err
	}
	if rtctrlSetNhUnchangedMap["type"] != "" {
		d.Set("next_hop_propagation", "yes")
	}
	return d, nil
}

func getRemoteRtctrlSetRedistMultipath(client *client.Client, dn string) (*models.RedistributeMultipathAction, error) {
	rtctrlSetRedistMultipathCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetRedistMultipath := models.RedistributeMultipathActionFromContainer(rtctrlSetRedistMultipathCont)
	if rtctrlSetRedistMultipath.DistinguishedName == "" {
		return nil, fmt.Errorf("Route Control Set Redist Multipath %s not found", dn)
	}
	return rtctrlSetRedistMultipath, nil
}

func setRtctrlSetRedistMultipathAttributes(rtctrlSetRedistMultipath *models.RedistributeMultipathAction, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetRedistMultipathMap, err := rtctrlSetRedistMultipath.ToMap()
	if err != nil {
		return d, err
	}
	if rtctrlSetRedistMultipathMap["type"] != "" {
		d.Set("multipath", "yes")
	}
	return d, nil
}

func getRemoteRtctrlSetASPath(client *client.Client, dn string) (*models.SetASPath, error) {
	rtctrlSetASPathCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetASPath := models.SetASPathFromContainer(rtctrlSetASPathCont)
	if rtctrlSetASPath.DistinguishedName == "" {
		return nil, fmt.Errorf("Route Control Set AS Path %s not found", dn)
	}
	return rtctrlSetASPath, nil
}

func setRtctrlSetASPathAttributes(rtctrlSetASPath *models.SetASPath, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetASPathMap, err := rtctrlSetASPath.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("set_as_path_prepend_last_as", rtctrlSetASPathMap["lastnum"])
	return d, nil
}

func getAndSetRemoteSetASPathASNs(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	ReadRelationASNumberData, err := client.ListSetAsPathASNs(dn)
	if err == nil {
		ASNList := make([]interface{}, 0)
		for _, record := range ReadRelationASNumberData {
			ASNMap := make(map[string]interface{})
			ASNMap["asn"] = record.Asn
			ASNMap["order"] = record.Order
			ASNList = append(ASNList, ASNMap)
		}
		d.Set("set_as_path_prepend_as", ASNList)
	} else {
		d.Set("set_as_path_prepend_as", nil)
	}
	return d, nil
}

func getRemoteRtctrlSetDamp(client *client.Client, dn string) (*models.RtctrlSetDamp, error) {
	rtctrlSetDampCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	rtctrlSetDamp := models.RtctrlSetDampFromContainer(rtctrlSetDampCont)
	if rtctrlSetDamp.DistinguishedName == "" {
		return nil, fmt.Errorf("Route Control Set Damp %s not found", dn)
	}
	return rtctrlSetDamp, nil
}

func setRtctrlSetDampAttributes(rtctrlSetDamp *models.RtctrlSetDamp, d *schema.ResourceData) (*schema.ResourceData, error) {
	rtctrlSetDampMap, err := rtctrlSetDamp.ToMap()
	if err != nil {
		return d, err
	}

	newContent := make(map[string]interface{})
	newContent["half_life"] = rtctrlSetDampMap["halfLife"]
	newContent["reuse"] = rtctrlSetDampMap["reuse"]
	newContent["suppress"] = rtctrlSetDampMap["suppress"]
	newContent["max_suppress_time"] = rtctrlSetDampMap["maxSuppressTime"]
	d.Set("set_dampening", newContent)

	return d, nil
}

func resourceAciActionRuleProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	rtctrlAttrP, err := getRemoteActionRuleProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	rtctrlAttrPMap, err := rtctrlAttrP.ToMap()
	if err != nil {
		return nil, err
	}
	name := rtctrlAttrPMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/attr-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setActionRuleProfileAttributes(rtctrlAttrP, d)
	if err != nil {
		return nil, err
	}

	// rtctrlSetTag - Beginning Import
	setRouteTagDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetTag)
	rtctrlSetTag, err := getRemoteRtctrlSetTag(aciClient, setRouteTagDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetTag - Beginning Import", setRouteTagDn)
		_, err = setRtctrlSetTagAttributes(rtctrlSetTag, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetTag - Import finished successfully", setRouteTagDn)
	}
	// rtctrlSetTag - Import finished successfully

	// rtctrlSetPref - Beginning Import
	setPrefDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetPref)
	rtctrlSetPref, err := getRemoteRtctrlSetPref(aciClient, setPrefDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetPref - Beginning Import", setPrefDn)
		_, err = setRtctrlSetPrefAttributes(rtctrlSetPref, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetPref - Import finished successfully", setPrefDn)
	}
	// rtctrlSetPref - Import finished successfully

	// rtctrlSetWeight - Beginning Import
	setWeightDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetWeight)
	rtctrlSetWeight, err := getRemoteRtctrlSetWeight(aciClient, setWeightDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetWeight - Beginning Import", setWeightDn)
		_, err = setRtctrlSetWeightAttributes(rtctrlSetWeight, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetWeight - Import finished successfully", setWeightDn)
	}
	// rtctrlSetWeight - Import finished successfully

	// rtctrlSetRtMetric - Beginning Import
	setRtMetricDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRtMetric)
	rtctrlSetRtMetric, err := getRemoteRtctrlSetRtMetric(aciClient, setRtMetricDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetRtMetric - Beginning Import", setRtMetricDn)
		_, err = setRtctrlSetRtMetricAttributes(rtctrlSetRtMetric, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetRtMetric - Import finished successfully", setRtMetricDn)
	}
	// rtctrlSetRtMetric - Import finished successfully

	// rtctrlSetRtMetricType - Beginning Import
	setRtMetricTypeDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRtMetricType)
	rtctrlSetRtMetricType, err := getRemoteRtctrlSetRtMetricType(aciClient, setRtMetricTypeDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetRtMetricType - Beginning Import", setRtMetricTypeDn)
		_, err = setRtctrlSetRtMetricTypeAttributes(rtctrlSetRtMetricType, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetRtMetricType - Import finished successfully", setRtMetricTypeDn)
	}
	// rtctrlSetRtMetricType - Import finished successfully

	// rtctrlSetNh - Beginning Import
	setNhDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetNh)
	rtctrlSetNh, err := getRemoteRtctrlSetNh(aciClient, setNhDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetNh - Beginning Import", setNhDn)
		_, err = setRtctrlSetNhAttributes(rtctrlSetNh, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetNh - Import finished successfully", setNhDn)
	}
	// rtctrlSetNh - Import finished successfully

	// rtctrlSetComm - Beginning Import
	setCommDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetComm)
	rtctrlSetComm, err := getRemoteRtctrlSetComm(aciClient, setCommDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetComm - Beginning Import", setCommDn)
		_, err = setRtctrlSetCommAttributes(rtctrlSetComm, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetComm - Import finished successfully", setCommDn)
	}
	// rtctrlSetComm - Import finished successfully

	// rtctrlSetNhUnchanged - Beginning Import
	setNhUnchangedDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetNhUnchanged)
	rtctrlSetNhUnchanged, err := getRemoteNexthopUnchangedAction(aciClient, setNhUnchangedDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetNhUnchanged - Beginning Import", setNhUnchangedDn)
		_, err := setNexthopUnchangedActionAttributes(rtctrlSetNhUnchanged, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetNhUnchanged - Import finished successfully", setNhUnchangedDn)
	}
	// rtctrlSetNhUnchanged - Import finished successfully

	// rtctrlSetRedistMultipath - Beginning Import
	setRedistMultipathDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRedistMultipath)
	rtctrlSetRedistMultipath, err := getRemoteRtctrlSetRedistMultipath(aciClient, setRedistMultipathDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetRedistMultipath - Beginning Import", setRedistMultipathDn)
		_, err = setRtctrlSetRedistMultipathAttributes(rtctrlSetRedistMultipath, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetRedistMultipath - Import finished successfully", setRedistMultipathDn)
	}
	// rtctrlSetRedistMultipath - Import finished successfully

	// rtctrlSetASPath - Beginning Import

	setASPathPrependLastAsDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetASPath, "prepend-last-as")
	rtctrlSetASPathPrependLastAs, err := getRemoteRtctrlSetASPath(aciClient, setASPathPrependLastAsDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetASPath - Beginning Import", setASPathPrependLastAsDn)
		_, err = setRtctrlSetASPathAttributes(rtctrlSetASPathPrependLastAs, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetASPath - Import finished successfully", setASPathPrependLastAsDn)
	}

	// rtctrlSetASPath - Import finished successfully

	// rtctrlSetASPathASN - Beginning Import

	setASNumberDn := rtctrlAttrP.DistinguishedName + "/" + fmt.Sprintf(models.RnrtctrlSetASPath, "prepend")
	log.Printf("[DEBUG] %s: rtctrlSetASPathASN - Beginning Import", setASNumberDn)
	_, err = getAndSetRemoteSetASPathASNs(aciClient, setASNumberDn, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: rtctrlSetASPathASN - Import finished successfully", setASNumberDn)

	// rtctrlSetASPathASN - Import finished successfully

	// rtctrlSetDamp - Beginning Import
	setDampeningDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetDamp)
	rtctrlSetDamp, err := getRemoteRtctrlSetDamp(aciClient, setDampeningDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetDamp - Beginning Import", setDampeningDn)
		_, err = setRtctrlSetDampAttributes(rtctrlSetDamp, d)
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] %s: rtctrlSetDamp - Import finished successfully", setDampeningDn)
	}
	// rtctrlSetDamp - Import finished successfully

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciActionRuleProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ActionRuleProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	rtctrlAttrPAttr := models.ActionRuleProfileAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlAttrPAttr.Annotation = Annotation.(string)
	} else {
		rtctrlAttrPAttr.Annotation = "{}"
	}
	if Name, ok := d.GetOk("name"); ok {
		rtctrlAttrPAttr.Name = Name.(string)
	}
	rtctrlAttrP := models.NewActionRuleProfile(fmt.Sprintf(models.RnrtctrlAttrP, name), TenantDn, desc, nameAlias, rtctrlAttrPAttr)

	err := aciClient.Save(rtctrlAttrP)
	if err != nil {
		return diag.FromErr(err)
	}

	// rtctrlSetTag - Operations
	if setRouteTag, ok := d.GetOk("set_route_tag"); ok {
		log.Printf("[DEBUG] rtctrlSetTag: Beginning Creation")

		rtctrlSetTagAttr := models.RtctrlSetTagAttributes{}
		rtctrlSetTagAttr.Tag = setRouteTag.(string)
		rtctrlSetTag := models.NewRtctrlSetTag(fmt.Sprintf(models.RnrtctrlSetTag), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetTagAttr)

		creation_err := aciClient.Save(rtctrlSetTag)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetTag.DistinguishedName)
	}

	// rtctrlSetPref - Operations
	if setPref, ok := d.GetOk("set_preference"); ok {
		log.Printf("[DEBUG] rtctrlSetPref: Beginning Creation")

		rtctrlSetPrefAttr := models.RtctrlSetPrefAttributes{}
		rtctrlSetPrefAttr.LocalPref = setPref.(string)
		rtctrlSetPref := models.NewRtctrlSetPref(fmt.Sprintf(models.RnrtctrlSetPref), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetPrefAttr)

		creation_err := aciClient.Save(rtctrlSetPref)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetPref.DistinguishedName)
	}

	// rtctrlSetWeight - Operations
	if setWeight, ok := d.GetOk("set_weight"); ok {
		log.Printf("[DEBUG] rtctrlSetWeight: Beginning Creation")

		rtctrlSetWeightAttr := models.RtctrlSetWeightAttributes{}
		rtctrlSetWeightAttr.Weight = setWeight.(string)
		rtctrlSetWeight := models.NewRtctrlSetWeight(fmt.Sprintf(models.RnrtctrlSetWeight), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetWeightAttr)

		creation_err := aciClient.Save(rtctrlSetWeight)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetWeight.DistinguishedName)
	}

	// rtctrlSetRtMetric - Operations
	if setRtMetric, ok := d.GetOk("set_metric"); ok {
		log.Printf("[DEBUG] rtctrlSetRtMetric: Beginning Creation")

		rtctrlSetRtMetricAttr := models.RtctrlSetRtMetricAttributes{}
		rtctrlSetRtMetricAttr.Metric = setRtMetric.(string)
		rtctrlSetRtMetric := models.NewRtctrlSetRtMetric(fmt.Sprintf(models.RnrtctrlSetRtMetric), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetRtMetricAttr)

		creation_err := aciClient.Save(rtctrlSetRtMetric)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetRtMetric.DistinguishedName)
	}

	// rtctrlSetRtMetricType - Operations
	if setRtMetricType, ok := d.GetOk("set_metric_type"); ok {
		log.Printf("[DEBUG] rtctrlSetRtMetricType: Beginning Creation")

		rtctrlSetRtMetricTypeAttr := models.RtctrlSetRtMetricTypeAttributes{}
		rtctrlSetRtMetricTypeAttr.MetricType = setRtMetricType.(string)
		rtctrlSetRtMetricType := models.NewRtctrlSetRtMetricType(fmt.Sprintf(models.RnrtctrlSetRtMetricType), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetRtMetricTypeAttr)

		creation_err := aciClient.Save(rtctrlSetRtMetricType)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetRtMetricType.DistinguishedName)
	}

	// rtctrlSetNh - Operations
	if setNh, ok := d.GetOk("set_next_hop"); ok {
		log.Printf("[DEBUG] rtctrlSetNh: Beginning Creation")

		rtctrlSetNhAttr := models.RtctrlSetNhAttributes{}
		rtctrlSetNhAttr.Addr = setNh.(string)
		rtctrlSetNh := models.NewRtctrlSetNh(fmt.Sprintf(models.RnrtctrlSetNh), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetNhAttr)

		creation_err := aciClient.Save(rtctrlSetNh)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetNh.DistinguishedName)
	}

	// rtctrlSetComm - Operations
	if setComm, ok := d.GetOk("set_communities"); ok {
		log.Printf("[DEBUG] rtctrlSetComm: Beginning Creation")
		rtctrlSetCommAttr := models.RtctrlSetCommAttributes{}

		setCommMap := toStrMap(setComm.(map[string]interface{}))
		rtctrlSetCommAttr.Community = setCommMap["community"]
		rtctrlSetCommAttr.SetCriteria = setCommMap["criteria"]
		rtctrlSetComm := models.NewRtctrlSetComm(fmt.Sprintf(models.RnrtctrlSetComm), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetCommAttr)

		creation_err := aciClient.Save(rtctrlSetComm)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetComm.DistinguishedName)
	}

	// rtctrlSetNhUnchanged - Operations
	setNhUnchanged, ok := d.GetOk("next_hop_propagation")
	if ok && setNhUnchanged == "yes" {
		log.Printf("[DEBUG] rtctrlSetNhUnchanged: Beginning Creation")
		rtctrlSetNhUnchangedAttr := models.NexthopUnchangedActionAttributes{}
		rtctrlSetNhUnchanged := models.NewNexthopUnchangedAction(fmt.Sprintf(models.RnrtctrlSetNhUnchanged), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetNhUnchangedAttr)

		creation_err := aciClient.Save(rtctrlSetNhUnchanged)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetNhUnchanged.DistinguishedName)
	}

	// rtctrlSetRedistMultipath - Operations
	setRedistMultipath, ok := d.GetOk("multipath")
	if ok && setRedistMultipath == "yes" {
		log.Printf("[DEBUG] rtctrlSetRedistMultipath: Beginning Creation")

		rtctrlSetRedistMultipathAttr := models.RedistributeMultipathActionAttributes{}
		rtctrlSetRedistMultipath := models.NewRedistributeMultipathAction(fmt.Sprintf(models.RnrtctrlSetRedistMultipath), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetRedistMultipathAttr)

		creation_err := aciClient.Save(rtctrlSetRedistMultipath)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetRedistMultipath.DistinguishedName)
	}

	// rtctrlSetASPath - Operations
	if prependLastAS, ok := d.GetOk("set_as_path_prepend_last_as"); ok {
		log.Printf("[DEBUG] rtctrlSetASPath prepend-last-as: Beginning Creation")

		rtctrlSetASPathAttr := models.SetASPathAttributes{}
		rtctrlSetASPathAttr.Lastnum = prependLastAS.(string)
		rtctrlSetASPath := models.NewSetASPath(fmt.Sprintf(models.RnrtctrlSetASPath, "prepend-last-as"), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetASPathAttr)

		creation_err := aciClient.Save(rtctrlSetASPath)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetASPath.DistinguishedName)
	}

	// rtctrlSetASPathASN - Operations
	if SetAsPathASN, ok := d.GetOk("set_as_path_prepend_as"); ok {

		log.Printf("[DEBUG] rtctrlSetASPath prepend: Beginning Creation of Parent Object")

		// Parent Object creation - begins
		rtctrlSetASPathAttr := models.SetASPathAttributes{}

		rtctrlSetASPath := models.NewSetASPath(fmt.Sprintf(models.RnrtctrlSetASPath, "prepend"), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetASPathAttr)

		parent_creation_err := aciClient.Save(rtctrlSetASPath)
		if parent_creation_err != nil {
			return diag.FromErr(parent_creation_err)
		}

		log.Printf("[DEBUG] %s: Creation of parent object finished successfully", rtctrlSetASPath.DistinguishedName)
		log.Printf("[DEBUG] rtctrlSetASPathASN: Beginning Creation of Child Objects")

		SetAsPathASNList := SetAsPathASN.(*schema.Set).List()
		for _, ASN := range SetAsPathASNList {
			ASNMap := ASN.(map[string]interface{})

			rtctrlSetASPathASNAttr := models.ASNumberAttributes{}
			rtctrlSetASPathASNAttr.Asn = ASNMap["asn"].(string)
			rtctrlSetASPathASNAttr.Order = ASNMap["order"].(string)

			// Child Object creation
			rtctrlSetASPathASN := models.NewASNumber(fmt.Sprintf(models.RnrtctrlSetASPathASN, ASNMap["order"]), rtctrlSetASPath.DistinguishedName, "", "", rtctrlSetASPathASNAttr)

			err := aciClient.Save(rtctrlSetASPathASN)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetASPathASN.DistinguishedName)
		}

		log.Printf("[DEBUG] %s: Creation of all child objects finished successfully", rtctrlSetASPath.DistinguishedName)
	}

	// rtctrlSetDamp - Operations
	if setDampening, ok := d.GetOk("set_dampening"); ok {
		log.Printf("[DEBUG] rtctrlSetDamp: Beginning Creation")
		rtctrlSetDampAttr := models.RtctrlSetDampAttributes{}

		setDampeningMap := toStrMap(setDampening.(map[string]interface{}))
		rtctrlSetDampAttr.HalfLife = setDampeningMap["half_life"]
		rtctrlSetDampAttr.Reuse = setDampeningMap["reuse"]
		rtctrlSetDampAttr.Suppress = setDampeningMap["suppress"]
		rtctrlSetDampAttr.MaxSuppressTime = setDampeningMap["max_suppress_time"]
		rtctrlSetDamp := models.NewRtctrlSetDamp(fmt.Sprintf(models.RnrtctrlSetDamp), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetDampAttr)

		creation_err := aciClient.Save(rtctrlSetDamp)
		if creation_err != nil {
			return diag.FromErr(creation_err)
		}

		log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetDamp.DistinguishedName)
	}

	d.SetId(rtctrlAttrP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciActionRuleProfileRead(ctx, d, m)
}

func resourceAciActionRuleProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ActionRuleProfile: Beginning Update")
	next_hop_propagation_flag := true

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	rtctrlAttrPAttr := models.ActionRuleProfileAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		rtctrlAttrPAttr.Annotation = Annotation.(string)
	} else {
		rtctrlAttrPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		rtctrlAttrPAttr.Name = Name.(string)
	}

	rtctrlAttrP := models.NewActionRuleProfile(fmt.Sprintf(models.RnrtctrlAttrP, name), TenantDn, desc, nameAlias, rtctrlAttrPAttr)

	rtctrlAttrP.Status = "modified"

	err := aciClient.Save(rtctrlAttrP)

	if err != nil {
		return diag.FromErr(err)
	}

	// rtctrlSetTag - Operations
	if d.HasChange("set_route_tag") {
		if setRouteTag, ok := d.GetOk("set_route_tag"); ok {
			log.Printf("[DEBUG] rtctrlSetTag - Beginning Creation")

			// Removing multipath and next_hop_propagation
			setRedistMultipathDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRedistMultipath)
			remove_multipath := aciClient.DeleteByDn(setRedistMultipathDn, "rtctrlSetRedistMultipath")
			if remove_multipath != nil {
				return diag.FromErr(remove_multipath)
			}

			setNhUnchangedDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetNhUnchanged)
			remove_next_hop_propagation := aciClient.DeleteByDn(setNhUnchangedDn, "rtctrlSetNhUnchanged")
			if remove_next_hop_propagation != nil {
				return diag.FromErr(remove_next_hop_propagation)
			}

			//  Set Route Tag object creation
			rtctrlSetTagAttr := models.RtctrlSetTagAttributes{}
			rtctrlSetTagAttr.Tag = setRouteTag.(string)

			setRouteTagDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetTag)

			deletion_err := aciClient.DeleteByDn(setRouteTagDn, "rtctrlSetTag")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetTag := models.NewRtctrlSetTag(fmt.Sprintf(models.RnrtctrlSetTag), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetTagAttr)

			err := aciClient.Save(rtctrlSetTag)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetTag - Creation finished successfully", rtctrlSetTag.DistinguishedName)
		} else {
			setRouteTagDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetTag)
			log.Printf("[DEBUG] %s: rtctrlSetTag - Beginning Destroy", setRouteTagDn)

			err := aciClient.DeleteByDn(setRouteTagDn, "rtctrlSetTag")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetTag - Destroy finished successfully", setRouteTagDn)
		}
	}

	// rtctrlSetPref - Operations
	if d.HasChange("set_preference") {
		if setPref, ok := d.GetOk("set_preference"); ok {
			log.Printf("[DEBUG] rtctrlSetPref - Beginning Creation")

			rtctrlSetPrefAttr := models.RtctrlSetPrefAttributes{}
			rtctrlSetPrefAttr.LocalPref = setPref.(string)

			setPrefDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetPref)

			deletion_err := aciClient.DeleteByDn(setPrefDn, "rtctrlSetPref")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetPref := models.NewRtctrlSetPref(fmt.Sprintf(models.RnrtctrlSetPref), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetPrefAttr)

			err := aciClient.Save(rtctrlSetPref)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetPref - Creation finished successfully", rtctrlSetPref.DistinguishedName)
		} else {
			setPrefDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetPref)
			log.Printf("[DEBUG] %s: rtctrlSetPref - Beginning Destroy", setPrefDn)

			err := aciClient.DeleteByDn(setPrefDn, "rtctrlSetPref")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetPref - Destroy finished successfully", setPrefDn)
		}
	}

	// rtctrlSetWeight - Operations
	if d.HasChange("set_weight") {
		if setWeight, ok := d.GetOk("set_weight"); ok {
			log.Printf("[DEBUG] rtctrlSetWeight - Beginning Creation")

			rtctrlSetWeightAttr := models.RtctrlSetWeightAttributes{}
			rtctrlSetWeightAttr.Weight = setWeight.(string)

			setWeightDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetWeight)

			deletion_err := aciClient.DeleteByDn(setWeightDn, "rtctrlSetWeight")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetWeight := models.NewRtctrlSetWeight(fmt.Sprintf(models.RnrtctrlSetWeight), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetWeightAttr)

			err := aciClient.Save(rtctrlSetWeight)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetWeight - Creation finished successfully", rtctrlSetWeight.DistinguishedName)
		} else {
			setWeightDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetWeight)
			log.Printf("[DEBUG] %s: rtctrlSetWeight - Beginning Destroy", setWeightDn)

			err := aciClient.DeleteByDn(setWeightDn, "rtctrlSetWeight")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetWeight - Destroy finished successfully", setWeightDn)
		}
	}

	// rtctrlSetRtMetric - Operations
	if d.HasChange("set_metric") {
		if setRtMetric, ok := d.GetOk("set_metric"); ok {
			log.Printf("[DEBUG] rtctrlSetRtMetric - Beginning Creation")

			rtctrlSetRtMetricAttr := models.RtctrlSetRtMetricAttributes{}
			rtctrlSetRtMetricAttr.Metric = setRtMetric.(string)

			setRtMetricDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRtMetric)

			deletion_err := aciClient.DeleteByDn(setRtMetricDn, "rtctrlSetRtMetric")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetRtMetric := models.NewRtctrlSetRtMetric(fmt.Sprintf(models.RnrtctrlSetRtMetric), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetRtMetricAttr)

			err := aciClient.Save(rtctrlSetRtMetric)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetRtMetric - Creation finished successfully", rtctrlSetRtMetric.DistinguishedName)
		} else {
			setRtMetricDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRtMetric)
			log.Printf("[DEBUG] %s: rtctrlSetRtMetric - Beginning Destroy", setRtMetricDn)

			err := aciClient.DeleteByDn(setRtMetricDn, "rtctrlSetRtMetric")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetRtMetric - Destroy finished successfully", setRtMetricDn)
		}
	}

	// rtctrlSetRtMetricType - Operations
	if d.HasChange("set_metric_type") {
		if setRtMetricType, ok := d.GetOk("set_metric_type"); ok {
			log.Printf("[DEBUG] rtctrlSetRtMetricType - Beginning Creation")

			rtctrlSetRtMetricTypeAttr := models.RtctrlSetRtMetricTypeAttributes{}
			rtctrlSetRtMetricTypeAttr.MetricType = setRtMetricType.(string)

			setRtMetricTypeDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRtMetricType)

			deletion_err := aciClient.DeleteByDn(setRtMetricTypeDn, "rtctrlSetRtMetricType")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetRtMetricType := models.NewRtctrlSetRtMetricType(fmt.Sprintf(models.RnrtctrlSetRtMetricType), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetRtMetricTypeAttr)

			err := aciClient.Save(rtctrlSetRtMetricType)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetRtMetricType - Creation finished successfully", rtctrlSetRtMetricType.DistinguishedName)
		} else {
			setRtMetricTypeDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRtMetricType)
			log.Printf("[DEBUG] %s: rtctrlSetRtMetricType - Beginning Destroy", setRtMetricTypeDn)

			err := aciClient.DeleteByDn(setRtMetricTypeDn, "rtctrlSetRtMetricType")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetRtMetricType - Destroy finished successfully", setRtMetricTypeDn)
		}
	}

	// rtctrlSetNh - Operations
	if d.HasChange("set_next_hop") {
		if setNh, ok := d.GetOk("set_next_hop"); ok {
			log.Printf("[DEBUG] rtctrlSetNh - Beginning Creation")

			rtctrlSetNhAttr := models.RtctrlSetNhAttributes{}
			rtctrlSetNhAttr.Addr = setNh.(string)

			setNhDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetNh)

			deletion_err := aciClient.DeleteByDn(setNhDn, "rtctrlSetNh")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetNh := models.NewRtctrlSetNh(fmt.Sprintf(models.RnrtctrlSetNh), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetNhAttr)

			err := aciClient.Save(rtctrlSetNh)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetNh - Creation finished successfully", rtctrlSetNh.DistinguishedName)
		} else {
			setNhDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetNh)
			log.Printf("[DEBUG] %s: rtctrlSetNh - Beginning Destroy", setNhDn)

			err := aciClient.DeleteByDn(setNhDn, "rtctrlSetNh")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetNh - Destroy finished successfully", setNhDn)
		}
	}

	// rtctrlSetComm - Operations
	if d.HasChange("set_communities") {
		if setComm, ok := d.GetOk("set_communities"); ok {
			log.Printf("[DEBUG] rtctrlSetComm - Beginning Creation")

			rtctrlSetCommAttr := models.RtctrlSetCommAttributes{}
			setCommMap := toStrMap(setComm.(map[string]interface{}))
			rtctrlSetCommAttr.Community = setCommMap["community"]
			rtctrlSetCommAttr.SetCriteria = setCommMap["criteria"]

			setCommDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetComm)

			deletion_err := aciClient.DeleteByDn(setCommDn, "rtctrlSetComm")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetComm := models.NewRtctrlSetComm(fmt.Sprintf(models.RnrtctrlSetComm), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetCommAttr)

			err := aciClient.Save(rtctrlSetComm)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetComm - Creation finished successfully", rtctrlSetComm.DistinguishedName)
		} else {
			setCommDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetComm)
			log.Printf("[DEBUG] %s: rtctrlSetComm - Beginning Destroy", setCommDn)

			err := aciClient.DeleteByDn(setCommDn, "rtctrlSetComm")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetComm - Destroy finished successfully", setCommDn)
		}
	}

	// rtctrlSetNhUnchanged - Operations
	if d.HasChange("next_hop_propagation") {
		setNhUnchanged, ok := d.GetOk("next_hop_propagation")
		if ok && setNhUnchanged == "yes" {

			log.Printf("[DEBUG] rtctrlSetNhUnchanged - Beginning Creation")
			rtctrlSetNhUnchangedAttr := models.NexthopUnchangedActionAttributes{}

			setNhUnchangedDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetNhUnchanged)

			deletion_err := aciClient.DeleteByDn(setNhUnchangedDn, "rtctrlSetNhUnchanged")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetNhUnchanged := models.NewNexthopUnchangedAction(fmt.Sprintf(models.RnrtctrlSetNhUnchanged), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetNhUnchangedAttr)

			err := aciClient.Save(rtctrlSetNhUnchanged)
			if err != nil {
				return diag.FromErr(err)
			}

			next_hop_propagation_flag = false
			log.Printf("[DEBUG] %s: rtctrlSetNhUnchanged - Creation finished successfully", rtctrlSetNhUnchanged.DistinguishedName)
		}
	}

	// rtctrlSetRedistMultipath - Operations
	if d.HasChange("multipath") {
		setRedistMultipath, ok := d.GetOk("multipath")
		if ok && setRedistMultipath == "yes" {
			log.Printf("[DEBUG] rtctrlSetRedistMultipath - Beginning Creation")

			rtctrlSetRedistMultipathAttr := models.RedistributeMultipathActionAttributes{}

			setRedistMultipathDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRedistMultipath)

			deletion_err := aciClient.DeleteByDn(setRedistMultipathDn, "rtctrlSetRedistMultipath")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetRedistMultipath := models.NewRedistributeMultipathAction(fmt.Sprintf(models.RnrtctrlSetRedistMultipath), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetRedistMultipathAttr)

			err := aciClient.Save(rtctrlSetRedistMultipath)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetRedistMultipath - Creation finished successfully", rtctrlSetRedistMultipath.DistinguishedName)
		} else {
			setRedistMultipathDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRedistMultipath)
			log.Printf("[DEBUG] %s: rtctrlSetRedistMultipath - Beginning Destroy", setRedistMultipathDn)

			err := aciClient.DeleteByDn(setRedistMultipathDn, "rtctrlSetRedistMultipath")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetRedistMultipath - Destroy finished successfully", setRedistMultipathDn)
		}
	}

	// rtctrlSetNhUnchanged - Operations
	if d.HasChange("next_hop_propagation") || d.Get("next_hop_propagation") == "no" {
		setNhUnchanged, ok := d.GetOk("next_hop_propagation")
		if ok && setNhUnchanged == "yes" {
			if next_hop_propagation_flag {
				log.Printf("[DEBUG] rtctrlSetNhUnchanged - Beginning Creation")
				rtctrlSetNhUnchangedAttr := models.NexthopUnchangedActionAttributes{}

				setNhUnchangedDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetNhUnchanged)

				deletion_err := aciClient.DeleteByDn(setNhUnchangedDn, "rtctrlSetNhUnchanged")
				if deletion_err != nil {
					return diag.FromErr(err)
				}

				rtctrlSetNhUnchanged := models.NewNexthopUnchangedAction(fmt.Sprintf(models.RnrtctrlSetNhUnchanged), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetNhUnchangedAttr)

				err := aciClient.Save(rtctrlSetNhUnchanged)
				if err != nil {
					return diag.FromErr(err)
				}

				log.Printf("[DEBUG] %s: rtctrlSetNhUnchanged - Creation finished successfully", rtctrlSetNhUnchanged.DistinguishedName)
			}
		} else {
			setRedistMultipath, ok := d.GetOk("multipath")
			if !ok || setRedistMultipath == "no" {
				setNhUnchangedDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetNhUnchanged)
				log.Printf("[DEBUG] %s: rtctrlSetNhUnchanged - Beginning Destroy", setNhUnchangedDn)

				err := aciClient.DeleteByDn(setNhUnchangedDn, "rtctrlSetNhUnchanged")

				if err != nil {
					return diag.FromErr(err)
				}

				log.Printf("[DEBUG] %s: rtctrlSetNhUnchanged - Destroy finished successfully", setNhUnchangedDn)
			} else {
				return diag.FromErr(fmt.Errorf("Invalid Configuration Set Redistribute Multipath action cannot be configured without configuring the Next Hop Propagation"))
			}
		}
	}

	// rtctrlSetASPath - Operations
	if d.HasChange("set_as_path_prepend_last_as") {

		if prependLastAS, ok := d.GetOk("set_as_path_prepend_last_as"); ok {

			log.Printf("[DEBUG] rtctrlSetASPath prepend-last-as - Beginning Creation")

			rtctrlSetASPathAttr := models.SetASPathAttributes{}
			rtctrlSetASPathAttr.Lastnum = prependLastAS.(string)

			setASPathPrependLastASDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetASPath, "prepend-last-as")

			deletion_err := aciClient.DeleteByDn(setASPathPrependLastASDn, "rtctrlSetASPath")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetASPath := models.NewSetASPath(fmt.Sprintf(models.RnrtctrlSetASPath, "prepend-last-as"), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetASPathAttr)

			err := aciClient.Save(rtctrlSetASPath)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetASPath - Creation finished successfully", rtctrlSetASPath.DistinguishedName)

		} else {
			setASPathPrependLastASDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetASPath, "prepend-last-as")
			log.Printf("[DEBUG] %s: rtctrlSetASPath - Beginning Destroy", setASPathPrependLastASDn)

			err := aciClient.DeleteByDn(setASPathPrependLastASDn, "rtctrlSetASPath")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetASPath - Destroy finished successfully", setASPathPrependLastASDn)
		}
	}

	// rtctrlSetASPathASN - Operations
	if d.HasChange("set_as_path_prepend_as") {
		if SetAsPathASN, ok := d.GetOk("set_as_path_prepend_as"); ok {
			log.Printf("[DEBUG] rtctrlSetASPathASN prepend - Beginning Creation")

			// Parent Object deletion
			setASPathPrependDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetASPath, "prepend")

			parent_deletion_err := aciClient.DeleteByDn(setASPathPrependDn, "rtctrlSetASPath")
			if parent_deletion_err != nil {
				return diag.FromErr(parent_deletion_err)
			}

			// Parent Object creation
			rtctrlSetASPathAttr := models.SetASPathAttributes{}
			rtctrlSetASPath := models.NewSetASPath(fmt.Sprintf(models.RnrtctrlSetASPath, "prepend"), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetASPathAttr)

			parent_creation_err := aciClient.Save(rtctrlSetASPath)
			if parent_creation_err != nil {
				return diag.FromErr(parent_creation_err)
			}

			log.Printf("[DEBUG] %s: Creation of parent object finished successfully", rtctrlSetASPath.DistinguishedName)
			log.Printf("[DEBUG] rtctrlSetASPathASN: Beginning Creation of Child Objects")

			SetAsPathASNList := SetAsPathASN.(*schema.Set).List()

			for _, ASN := range SetAsPathASNList {
				ASNMap := ASN.(map[string]interface{})

				rtctrlSetASPathASNAttr := models.ASNumberAttributes{}
				rtctrlSetASPathASNAttr.Asn = ASNMap["asn"].(string)
				rtctrlSetASPathASNAttr.Order = ASNMap["order"].(string)

				// Child Object creation
				rtctrlSetASPathASN := models.NewASNumber(fmt.Sprintf(models.RnrtctrlSetASPathASN, ASNMap["order"]), rtctrlSetASPath.DistinguishedName, "", "", rtctrlSetASPathASNAttr)

				err := aciClient.Save(rtctrlSetASPathASN)
				if err != nil {
					return diag.FromErr(err)
				}
				log.Printf("[DEBUG] %s: Creation finished successfully", rtctrlSetASPathASN.DistinguishedName)
			}

			log.Printf("[DEBUG] %s: Creation of all child objects finished successfully", rtctrlSetASPath.DistinguishedName)

		} else {
			// Parent Object deletion
			setASPathPrependDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetASPath, "prepend")
			log.Printf("[DEBUG] %s: rtctrlSetASPath - Beginning Destroy", setASPathPrependDn)

			err := aciClient.DeleteByDn(setASPathPrependDn, "rtctrlSetASPath")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetASPath - Destroy finished successfully", setASPathPrependDn)
		}
	}

	// rtctrlSetDamp - Operations
	if d.HasChange("set_dampening") {
		if setDampening, ok := d.GetOk("set_dampening"); ok {
			log.Printf("[DEBUG] rtctrlSetDamp - Beginning Creation")

			rtctrlSetDampAttr := models.RtctrlSetDampAttributes{}
			setDampeningMap := toStrMap(setDampening.(map[string]interface{}))
			rtctrlSetDampAttr.HalfLife = setDampeningMap["half_life"]
			rtctrlSetDampAttr.Reuse = setDampeningMap["reuse"]
			rtctrlSetDampAttr.Suppress = setDampeningMap["suppress"]
			rtctrlSetDampAttr.MaxSuppressTime = setDampeningMap["max_suppress_time"]

			setDampeningDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetDamp)

			deletion_err := aciClient.DeleteByDn(setDampeningDn, "rtctrlSetDamp")
			if deletion_err != nil {
				return diag.FromErr(err)
			}

			rtctrlSetDamp := models.NewRtctrlSetDamp(fmt.Sprintf(models.RnrtctrlSetDamp), rtctrlAttrP.DistinguishedName, "", "", rtctrlSetDampAttr)

			err := aciClient.Save(rtctrlSetDamp)
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetDamp - Creation finished successfully", rtctrlSetDamp.DistinguishedName)
		} else {
			setDampeningDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetDamp)
			log.Printf("[DEBUG] %s: rtctrlSetDamp - Beginning Destroy", setDampeningDn)

			err := aciClient.DeleteByDn(setDampeningDn, "rtctrlSetDamp")
			if err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[DEBUG] %s: rtctrlSetDamp - Destroy finished successfully", setDampeningDn)
		}
	}

	d.SetId(rtctrlAttrP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciActionRuleProfileRead(ctx, d, m)
}

func resourceAciActionRuleProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	rtctrlAttrP, err := getRemoteActionRuleProfile(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setActionRuleProfileAttributes(rtctrlAttrP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// rtctrlSetTag - Beginning Read
	setRouteTagDn := dn + fmt.Sprintf("/"+models.RnrtctrlSetTag)
	rtctrlSetTag, err := getRemoteRtctrlSetTag(aciClient, setRouteTagDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetTag - Beginning Read", setRouteTagDn)
		_, err = setRtctrlSetTagAttributes(rtctrlSetTag, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: rtctrlSetTag - Read finished successfully", setRouteTagDn)
	} else {
		d.Set("set_route_tag", "")
	}
	// rtctrlSetTag - Read finished successfully

	// rtctrlSetPref - Beginning Read
	setPrefDn := dn + fmt.Sprintf("/"+models.RnrtctrlSetPref)
	rtctrlSetPref, err := getRemoteRtctrlSetPref(aciClient, setPrefDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetPref - Beginning Read", setPrefDn)
		_, err = setRtctrlSetPrefAttributes(rtctrlSetPref, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: rtctrlSetPref - Read finished successfully", setPrefDn)
	} else {
		d.Set("set_preference", "")
	}
	// rtctrlSetPref - Read finished successfully

	// rtctrlSetWeight - Beginning Read
	setWeightDn := dn + fmt.Sprintf("/"+models.RnrtctrlSetWeight)
	rtctrlSetWeight, err := getRemoteRtctrlSetWeight(aciClient, setWeightDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetWeight - Beginning Read", setWeightDn)
		_, err = setRtctrlSetWeightAttributes(rtctrlSetWeight, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: rtctrlSetWeight - Read finished successfully", setWeightDn)
	} else {
		d.Set("set_weight", "")
	}
	// rtctrlSetWeight - Read finished successfully

	// rtctrlSetRtMetric - Beginning Read
	setRtMetricDn := dn + fmt.Sprintf("/"+models.RnrtctrlSetRtMetric)
	rtctrlSetRtMetric, err := getRemoteRtctrlSetRtMetric(aciClient, setRtMetricDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetRtMetric - Beginning Read", setRtMetricDn)
		_, err = setRtctrlSetRtMetricAttributes(rtctrlSetRtMetric, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: rtctrlSetRtMetric - Read finished successfully", setRtMetricDn)
	} else {
		d.Set("set_metric", "")
	}
	// rtctrlSetRtMetric - Read finished successfully

	// rtctrlSetRtMetricType - Beginning Read
	setRtMetricTypeDn := dn + fmt.Sprintf("/"+models.RnrtctrlSetRtMetricType)
	rtctrlSetRtMetricType, err := getRemoteRtctrlSetRtMetricType(aciClient, setRtMetricTypeDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetRtMetricType - Beginning Read", setRtMetricTypeDn)
		_, err = setRtctrlSetRtMetricTypeAttributes(rtctrlSetRtMetricType, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: rtctrlSetRtMetricType - Read finished successfully", setRtMetricTypeDn)
	} else {
		d.Set("set_metric_type", "")
	}
	// rtctrlSetRtMetricType - Read finished successfully

	// rtctrlSetNh - Beginning Read
	setNhDn := dn + fmt.Sprintf("/"+models.RnrtctrlSetNh)
	rtctrlSetNh, err := getRemoteRtctrlSetNh(aciClient, setNhDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetNh - Beginning Read", setNhDn)
		_, err = setRtctrlSetNhAttributes(rtctrlSetNh, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: rtctrlSetNh - Read finished successfully", setNhDn)
	} else {
		d.Set("set_next_hop", "")
	}
	// rtctrlSetNh - Read finished successfully

	// rtctrlSetComm - Beginning Read
	setCommDn := dn + fmt.Sprintf("/"+models.RnrtctrlSetComm)
	rtctrlSetComm, err := getRemoteRtctrlSetComm(aciClient, setCommDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetComm - Beginning Read", setCommDn)
		_, err = setRtctrlSetCommAttributes(rtctrlSetComm, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: rtctrlSetComm - Read finished successfully", setCommDn)
	} else {
		d.Set("set_communities", nil)
	}
	// rtctrlSetComm - Read finished successfully

	// rtctrlSetNhUnchanged - Beginning Read
	setNhUnchangedDn := dn + fmt.Sprintf("/"+models.RnrtctrlSetNhUnchanged)
	rtctrlSetNhUnchanged, err := getRemoteNexthopUnchangedAction(aciClient, setNhUnchangedDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetNhUnchanged - Beginning Read", setNhUnchangedDn)
		_, err := setNexthopUnchangedActionAttributes(rtctrlSetNhUnchanged, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: rtctrlSetNhUnchanged - Read finished successfully", setNhUnchangedDn)
	} else {
		d.Set("next_hop_propagation", "")
	}
	// rtctrlSetNhUnchanged - Read finished successfully

	// rtctrlSetRedistMultipath - Beginning Read
	setRedistMultipathDn := dn + fmt.Sprintf("/"+models.RnrtctrlSetRedistMultipath)
	rtctrlSetRedistMultipath, err := getRemoteRtctrlSetRedistMultipath(aciClient, setRedistMultipathDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetRedistMultipath - Beginning Read", setRedistMultipathDn)
		_, err = setRtctrlSetRedistMultipathAttributes(rtctrlSetRedistMultipath, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: rtctrlSetRedistMultipath - Read finished successfully", setRedistMultipathDn)
	} else {
		d.Set("multipath", "")
	}
	// rtctrlSetRedistMultipath - Read finished successfully

	// rtctrlSetASPath - Beginning Read

	setASPathPrependLastAsDn := dn + fmt.Sprintf("/"+models.RnrtctrlSetASPath, "prepend-last-as")
	rtctrlSetASPathPrependLastAs, err := getRemoteRtctrlSetASPath(aciClient, setASPathPrependLastAsDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetASPath - Beginning Read", setASPathPrependLastAsDn)
		_, err = setRtctrlSetASPathAttributes(rtctrlSetASPathPrependLastAs, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: rtctrlSetASPath - Read finished successfully", setASPathPrependLastAsDn)
	} else {
		d.Set("set_as_path_prepend_last_as", "")
	}

	// rtctrlSetASPath - Read finished successfully

	// rtctrlSetASPathASN - Beginning Read

	setASNumberDn := dn + "/" + fmt.Sprintf(models.RnrtctrlSetASPath, "prepend")
	log.Printf("[DEBUG] %s: rtctrlSetASPathASN - Beginning Read", setASNumberDn)
	_, err = getAndSetRemoteSetASPathASNs(aciClient, setASNumberDn, d)
	if err != nil {
		return nil
	}
	log.Printf("[DEBUG] %s: rtctrlSetASPathASN - Read finished successfully", setASNumberDn)

	// rtctrlSetASPathASN - Read finished successfully

	// rtctrlSetDamp - Beginning Read
	setDampeningDn := dn + fmt.Sprintf("/"+models.RnrtctrlSetDamp)
	rtctrlSetDamp, err := getRemoteRtctrlSetDamp(aciClient, setDampeningDn)
	if err == nil {
		log.Printf("[DEBUG] %s: rtctrlSetDamp - Beginning Read", setDampeningDn)
		_, err = setRtctrlSetDampAttributes(rtctrlSetDamp, d)
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] %s: rtctrlSetDamp - Read finished successfully", setDampeningDn)
	} else {
		d.Set("set_dampening", nil)
	}
	// rtctrlSetDamp - Read finished successfully

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciActionRuleProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "rtctrlAttrP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

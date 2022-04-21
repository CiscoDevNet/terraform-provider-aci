package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"set_preference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"set_weight": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"set_metric": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"set_metric_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ospf-type1",
					"ospf-type2",
				}, false),
			},
			"set_next_hop": {
				Type:     schema.TypeString,
				Optional: true,
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
		return nil, fmt.Errorf("ActionRuleProfile %s not found", dn)
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
		return nil, fmt.Errorf("rtctrlSetRtMetric %s not found", dn)
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
		return nil, fmt.Errorf("rtctrlSetRtMetricType %s not found", dn)
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
		return nil, fmt.Errorf("rtctrlSetNh %s not found", dn)
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
		resourceAciRtctrlSetTagRead(ctx, rtctrlSetTag.DistinguishedName, d, m)
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
		resourceAciRtctrlSetPrefRead(ctx, rtctrlSetPref.DistinguishedName, d, m)
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
		resourceAciRtctrlSetWeightRead(ctx, rtctrlSetWeight.DistinguishedName, d, m)
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
		resourceAciRtctrlSetRtMetricRead(ctx, rtctrlSetRtMetric.DistinguishedName, d, m)
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
		resourceAciRtctrlSetRtMetricTypeRead(ctx, rtctrlSetRtMetricType.DistinguishedName, d, m)
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
		resourceAciRtctrlSetNhRead(ctx, rtctrlSetNh.DistinguishedName, d, m)
	}

	d.SetId(rtctrlAttrP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciActionRuleProfileRead(ctx, d, m)
}

func resourceAciActionRuleProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ActionRuleProfile: Beginning Update")

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
			resourceAciRtctrlSetTagRead(ctx, rtctrlSetTag.DistinguishedName, d, m)

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
			resourceAciRtctrlSetPrefRead(ctx, rtctrlSetPref.DistinguishedName, d, m)
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
			resourceAciRtctrlSetWeightRead(ctx, rtctrlSetWeight.DistinguishedName, d, m)
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
			resourceAciRtctrlSetRtMetricRead(ctx, rtctrlSetRtMetric.DistinguishedName, d, m)
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
			resourceAciRtctrlSetRtMetricTypeRead(ctx, rtctrlSetRtMetricType.DistinguishedName, d, m)
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
			resourceAciRtctrlSetNhRead(ctx, rtctrlSetNh.DistinguishedName, d, m)
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
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setActionRuleProfileAttributes(rtctrlAttrP, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciRtctrlSetTagRead(ctx context.Context, dn string, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: rtctrlSetTag - Beginning Read", dn)
	aciClient := m.(*client.Client)

	rtctrlSetTag, err := getRemoteRtctrlSetTag(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = setRtctrlSetTagAttributes(rtctrlSetTag, d)
	if err != nil {
		return nil
	}

	log.Printf("[DEBUG] %s: rtctrlSetTag - Read finished successfully", dn)
	return nil
}

func resourceAciRtctrlSetPrefRead(ctx context.Context, dn string, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: rtctrlSetPref - Beginning Read", dn)
	aciClient := m.(*client.Client)

	rtctrlSetPref, err := getRemoteRtctrlSetPref(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = setRtctrlSetPrefAttributes(rtctrlSetPref, d)
	if err != nil {
		return nil
	}

	log.Printf("[DEBUG] %s: rtctrlSetPref - Read finished successfully", dn)
	return nil
}

func resourceAciRtctrlSetWeightRead(ctx context.Context, dn string, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: rtctrlSetWeight - Beginning Read", dn)
	aciClient := m.(*client.Client)

	rtctrlSetWeight, err := getRemoteRtctrlSetWeight(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = setRtctrlSetWeightAttributes(rtctrlSetWeight, d)
	if err != nil {
		return nil
	}

	log.Printf("[DEBUG] %s: rtctrlSetWeight - Read finished successfully", dn)
	return nil
}

func resourceAciRtctrlSetRtMetricRead(ctx context.Context, dn string, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: rtctrlSetRtMetric - Beginning Read", dn)
	aciClient := m.(*client.Client)

	rtctrlSetRtMetric, err := getRemoteRtctrlSetRtMetric(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = setRtctrlSetRtMetricAttributes(rtctrlSetRtMetric, d)
	if err != nil {
		return nil
	}

	log.Printf("[DEBUG] %s: rtctrlSetRtMetric - Read finished successfully", dn)
	return nil
}

func resourceAciRtctrlSetRtMetricTypeRead(ctx context.Context, dn string, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: rtctrlSetRtMetricType - Beginning Read", dn)
	aciClient := m.(*client.Client)

	rtctrlSetRtMetricType, err := getRemoteRtctrlSetRtMetricType(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = setRtctrlSetRtMetricTypeAttributes(rtctrlSetRtMetricType, d)
	if err != nil {
		return nil
	}

	log.Printf("[DEBUG] %s: rtctrlSetRtMetricType - Read finished successfully", dn)
	return nil
}

func resourceAciRtctrlSetNhRead(ctx context.Context, dn string, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: rtctrlSetNh - Beginning Read", dn)
	aciClient := m.(*client.Client)

	rtctrlSetNh, err := getRemoteRtctrlSetNh(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = setRtctrlSetNhAttributes(rtctrlSetNh, d)
	if err != nil {
		return nil
	}

	log.Printf("[DEBUG] %s: rtctrlSetNh - Read finished successfully", dn)
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

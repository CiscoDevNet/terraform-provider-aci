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

func resourceAciMaintenancePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMaintenancePolicyCreate,
		UpdateContext: resourceAciMaintenancePolicyUpdate,
		ReadContext:   resourceAciMaintenancePolicyRead,
		DeleteContext: resourceAciMaintenancePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMaintenancePolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"untriggered",
					"triggered",
				}, false),
			},

			"graceful": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"ignore_compat": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"internal_label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"notif_cond": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"notifyOnlyOnFailures",
					"notifyAlwaysBetweenSets",
					"notifyNever",
				}, false),
			},

			"run_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"pauseOnlyOnFailures",
					"pauseAlwaysBetweenSets",
					"pauseNever",
				}, false),
			},

			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"version_check_override": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"trigger-immediate",
					"trigger",
					"triggered",
					"untriggered",
				}, false),
			},

			"relation_maint_rs_pol_scheduler": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_maint_rs_pol_notif": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_trig_rs_triggerable": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteMaintenancePolicy(client *client.Client, dn string) (*models.MaintenancePolicy, error) {
	maintMaintPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	maintMaintP := models.MaintenancePolicyFromContainer(maintMaintPCont)

	if maintMaintP.DistinguishedName == "" {
		return nil, fmt.Errorf("MaintenancePolicy %s not found", maintMaintP.DistinguishedName)
	}

	return maintMaintP, nil
}

func setMaintenancePolicyAttributes(maintMaintP *models.MaintenancePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(maintMaintP.DistinguishedName)
	d.Set("description", maintMaintP.Description)
	maintMaintPMap, err := maintMaintP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", maintMaintPMap["name"])

	d.Set("admin_st", maintMaintPMap["adminSt"])
	d.Set("annotation", maintMaintPMap["annotation"])
	d.Set("graceful", maintMaintPMap["graceful"])
	d.Set("ignore_compat", maintMaintPMap["ignoreCompat"])
	d.Set("internal_label", maintMaintPMap["internalLabel"])
	d.Set("name_alias", maintMaintPMap["nameAlias"])
	d.Set("notif_cond", maintMaintPMap["notifCond"])
	d.Set("run_mode", maintMaintPMap["runMode"])
	d.Set("version", maintMaintPMap["version"])
	d.Set("version_check_override", maintMaintPMap["versionCheckOverride"])
	return d, nil
}

func resourceAciMaintenancePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	maintMaintP, err := getRemoteMaintenancePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setMaintenancePolicyAttributes(maintMaintP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMaintenancePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MaintenancePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	maintMaintPAttr := models.MaintenancePolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		maintMaintPAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		maintMaintPAttr.Annotation = Annotation.(string)
	} else {
		maintMaintPAttr.Annotation = "{}"
	}
	if Graceful, ok := d.GetOk("graceful"); ok {
		maintMaintPAttr.Graceful = Graceful.(string)
	}
	if IgnoreCompat, ok := d.GetOk("ignore_compat"); ok {
		maintMaintPAttr.IgnoreCompat = IgnoreCompat.(string)
	}
	if InternalLabel, ok := d.GetOk("internal_label"); ok {
		maintMaintPAttr.InternalLabel = InternalLabel.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		maintMaintPAttr.NameAlias = NameAlias.(string)
	}
	if NotifCond, ok := d.GetOk("notif_cond"); ok {
		maintMaintPAttr.NotifCond = NotifCond.(string)
	}
	if RunMode, ok := d.GetOk("run_mode"); ok {
		maintMaintPAttr.RunMode = RunMode.(string)
	}
	if Version, ok := d.GetOk("version"); ok {
		maintMaintPAttr.Version = Version.(string)
	}
	if VersionCheckOverride, ok := d.GetOk("version_check_override"); ok {
		maintMaintPAttr.VersionCheckOverride = VersionCheckOverride.(string)
	}
	maintMaintP := models.NewMaintenancePolicy(fmt.Sprintf("fabric/maintpol-%s", name), "uni", desc, maintMaintPAttr)

	err := aciClient.Save(maintMaintP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTomaintRsPolScheduler, ok := d.GetOk("relation_maint_rs_pol_scheduler"); ok {
		relationParam := relationTomaintRsPolScheduler.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTomaintRsPolNotif, ok := d.GetOk("relation_maint_rs_pol_notif"); ok {
		relationParam := relationTomaintRsPolNotif.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTotrigRsTriggerable, ok := d.GetOk("relation_trig_rs_triggerable"); ok {
		relationParam := relationTotrigRsTriggerable.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTomaintRsPolScheduler, ok := d.GetOk("relation_maint_rs_pol_scheduler"); ok {
		relationParam := relationTomaintRsPolScheduler.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationmaintRsPolSchedulerFromMaintenancePolicy(maintMaintP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTomaintRsPolNotif, ok := d.GetOk("relation_maint_rs_pol_notif"); ok {
		relationParam := relationTomaintRsPolNotif.(string)
		err = aciClient.CreateRelationmaintRsPolNotifFromMaintenancePolicy(maintMaintP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTotrigRsTriggerable, ok := d.GetOk("relation_trig_rs_triggerable"); ok {
		relationParam := relationTotrigRsTriggerable.(string)
		err = aciClient.CreateRelationtrigRsTriggerableFromMaintenancePolicy(maintMaintP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(maintMaintP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciMaintenancePolicyRead(ctx, d, m)
}

func resourceAciMaintenancePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MaintenancePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	maintMaintPAttr := models.MaintenancePolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		maintMaintPAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		maintMaintPAttr.Annotation = Annotation.(string)
	} else {
		maintMaintPAttr.Annotation = "{}"
	}
	if Graceful, ok := d.GetOk("graceful"); ok {
		maintMaintPAttr.Graceful = Graceful.(string)
	}
	if IgnoreCompat, ok := d.GetOk("ignore_compat"); ok {
		maintMaintPAttr.IgnoreCompat = IgnoreCompat.(string)
	}
	if InternalLabel, ok := d.GetOk("internal_label"); ok {
		maintMaintPAttr.InternalLabel = InternalLabel.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		maintMaintPAttr.NameAlias = NameAlias.(string)
	}
	if NotifCond, ok := d.GetOk("notif_cond"); ok {
		maintMaintPAttr.NotifCond = NotifCond.(string)
	}
	if RunMode, ok := d.GetOk("run_mode"); ok {
		maintMaintPAttr.RunMode = RunMode.(string)
	}
	if Version, ok := d.GetOk("version"); ok {
		maintMaintPAttr.Version = Version.(string)
	}
	if VersionCheckOverride, ok := d.GetOk("version_check_override"); ok {
		maintMaintPAttr.VersionCheckOverride = VersionCheckOverride.(string)
	}
	maintMaintP := models.NewMaintenancePolicy(fmt.Sprintf("fabric/maintpol-%s", name), "uni", desc, maintMaintPAttr)

	maintMaintP.Status = "modified"

	err := aciClient.Save(maintMaintP)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_maint_rs_pol_scheduler") {
		_, newRelParam := d.GetChange("relation_maint_rs_pol_scheduler")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_maint_rs_pol_notif") {
		_, newRelParam := d.GetChange("relation_maint_rs_pol_notif")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_trig_rs_triggerable") {
		_, newRelParam := d.GetChange("relation_trig_rs_triggerable")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_maint_rs_pol_scheduler") {
		_, newRelParam := d.GetChange("relation_maint_rs_pol_scheduler")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationmaintRsPolSchedulerFromMaintenancePolicy(maintMaintP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_maint_rs_pol_notif") {
		_, newRelParam := d.GetChange("relation_maint_rs_pol_notif")
		err = aciClient.DeleteRelationmaintRsPolNotifFromMaintenancePolicy(maintMaintP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationmaintRsPolNotifFromMaintenancePolicy(maintMaintP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_trig_rs_triggerable") {
		_, newRelParam := d.GetChange("relation_trig_rs_triggerable")
		err = aciClient.CreateRelationtrigRsTriggerableFromMaintenancePolicy(maintMaintP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(maintMaintP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciMaintenancePolicyRead(ctx, d, m)

}

func resourceAciMaintenancePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	maintMaintP, err := getRemoteMaintenancePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setMaintenancePolicyAttributes(maintMaintP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	maintRsPolSchedulerData, err := aciClient.ReadRelationmaintRsPolSchedulerFromMaintenancePolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation maintRsPolScheduler %v", err)
		d.Set("relation_maint_rs_pol_scheduler", "")

	} else {
		d.Set("relation_maint_rs_pol_scheduler", maintRsPolSchedulerData.(string))
	}

	maintRsPolNotifData, err := aciClient.ReadRelationmaintRsPolNotifFromMaintenancePolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation maintRsPolNotif %v", err)
		d.Set("relation_maint_rs_pol_notif", "")

	} else {
		d.Set("relation_maint_rs_pol_notif", maintRsPolNotifData.(string))
	}

	trigRsTriggerableData, err := aciClient.ReadRelationtrigRsTriggerableFromMaintenancePolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation trigRsTriggerable %v", err)
		d.Set("relation_trig_rs_triggerable", "")

	} else {
		d.Set("relation_trig_rs_triggerable", trigRsTriggerableData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciMaintenancePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "maintMaintP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

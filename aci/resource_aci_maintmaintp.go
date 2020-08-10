package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciMaintenancePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciMaintenancePolicyCreate,
		Update: resourceAciMaintenancePolicyUpdate,
		Read:   resourceAciMaintenancePolicyRead,
		Delete: resourceAciMaintenancePolicyDelete,

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
			},

			"graceful": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ignore_compat": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			},

			"run_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func setMaintenancePolicyAttributes(maintMaintP *models.MaintenancePolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(maintMaintP.DistinguishedName)
	d.Set("description", maintMaintP.Description)
	maintMaintPMap, _ := maintMaintP.ToMap()

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
	return d
}

func resourceAciMaintenancePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	maintMaintP, err := getRemoteMaintenancePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setMaintenancePolicyAttributes(maintMaintP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMaintenancePolicyCreate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTomaintRsPolScheduler, ok := d.GetOk("relation_maint_rs_pol_scheduler"); ok {
		relationParam := relationTomaintRsPolScheduler.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationmaintRsPolSchedulerFromMaintenancePolicy(maintMaintP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_maint_rs_pol_scheduler")
		d.Partial(false)

	}
	if relationTomaintRsPolNotif, ok := d.GetOk("relation_maint_rs_pol_notif"); ok {
		relationParam := relationTomaintRsPolNotif.(string)
		err = aciClient.CreateRelationmaintRsPolNotifFromMaintenancePolicy(maintMaintP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_maint_rs_pol_notif")
		d.Partial(false)

	}
	if relationTotrigRsTriggerable, ok := d.GetOk("relation_trig_rs_triggerable"); ok {
		relationParam := relationTotrigRsTriggerable.(string)
		err = aciClient.CreateRelationtrigRsTriggerableFromMaintenancePolicy(maintMaintP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_trig_rs_triggerable")
		d.Partial(false)

	}

	d.SetId(maintMaintP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciMaintenancePolicyRead(d, m)
}

func resourceAciMaintenancePolicyUpdate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_maint_rs_pol_scheduler") {
		_, newRelParam := d.GetChange("relation_maint_rs_pol_scheduler")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationmaintRsPolSchedulerFromMaintenancePolicy(maintMaintP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_maint_rs_pol_scheduler")
		d.Partial(false)

	}
	if d.HasChange("relation_maint_rs_pol_notif") {
		_, newRelParam := d.GetChange("relation_maint_rs_pol_notif")
		err = aciClient.DeleteRelationmaintRsPolNotifFromMaintenancePolicy(maintMaintP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationmaintRsPolNotifFromMaintenancePolicy(maintMaintP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_maint_rs_pol_notif")
		d.Partial(false)

	}
	if d.HasChange("relation_trig_rs_triggerable") {
		_, newRelParam := d.GetChange("relation_trig_rs_triggerable")
		err = aciClient.CreateRelationtrigRsTriggerableFromMaintenancePolicy(maintMaintP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_trig_rs_triggerable")
		d.Partial(false)

	}

	d.SetId(maintMaintP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciMaintenancePolicyRead(d, m)

}

func resourceAciMaintenancePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	maintMaintP, err := getRemoteMaintenancePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setMaintenancePolicyAttributes(maintMaintP, d)

	maintRsPolSchedulerData, err := aciClient.ReadRelationmaintRsPolSchedulerFromMaintenancePolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation maintRsPolScheduler %v", err)

	} else {
		if _, ok := d.GetOk("relation_maint_rs_pol_scheduler"); ok {
			tfName := GetMOName(d.Get("relation_maint_rs_pol_scheduler").(string))
			if tfName != maintRsPolSchedulerData {
				d.Set("relation_maint_rs_pol_scheduler", "")
			}
		}
	}

	maintRsPolNotifData, err := aciClient.ReadRelationmaintRsPolNotifFromMaintenancePolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation maintRsPolNotif %v", err)

	} else {
		if _, ok := d.GetOk("relation_maint_rs_pol_notif"); ok {
			tfName := d.Get("relation_maint_rs_pol_notif").(string)
			if tfName != maintRsPolNotifData {
				d.Set("relation_maint_rs_pol_notif", "")
			}
		}
	}

	trigRsTriggerableData, err := aciClient.ReadRelationtrigRsTriggerableFromMaintenancePolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation trigRsTriggerable %v", err)

	} else {
		if _, ok := d.GetOk("relation_trig_rs_triggerable"); ok {
			tfName := d.Get("relation_trig_rs_triggerable").(string)
			if tfName != trigRsTriggerableData {
				d.Set("relation_trig_rs_triggerable", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciMaintenancePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "maintMaintP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

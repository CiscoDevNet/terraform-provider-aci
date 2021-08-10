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

func resourceAciConfigurationExportPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciConfigurationExportPolicyCreate,
		UpdateContext: resourceAciConfigurationExportPolicyUpdate,
		ReadContext:   resourceAciConfigurationExportPolicyRead,
		DeleteContext: resourceAciConfigurationExportPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciConfigurationExportPolicyImport,
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

			"format": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"xml",
					"json",
				}, false),
			},

			"include_secure_fields": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"max_snapshot_count": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"snapshot": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"target_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_config_rs_export_destination": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_trig_rs_triggerable": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_config_rs_remote_path": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_config_rs_export_scheduler": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteConfigurationExportPolicy(client *client.Client, dn string) (*models.ConfigurationExportPolicy, error) {
	configExportPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	configExportP := models.ConfigurationExportPolicyFromContainer(configExportPCont)

	if configExportP.DistinguishedName == "" {
		return nil, fmt.Errorf("ConfigurationExportPolicy %s not found", configExportP.DistinguishedName)
	}

	return configExportP, nil
}

func setConfigurationExportPolicyAttributes(configExportP *models.ConfigurationExportPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(configExportP.DistinguishedName)
	d.Set("description", configExportP.Description)
	configExportPMap, err := configExportP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", configExportPMap["name"])

	d.Set("admin_st", configExportPMap["adminSt"])
	d.Set("annotation", configExportPMap["annotation"])
	d.Set("format", configExportPMap["format"])
	d.Set("include_secure_fields", configExportPMap["includeSecureFields"])
	d.Set("max_snapshot_count", configExportPMap["maxSnapshotCount"])
	d.Set("name_alias", configExportPMap["nameAlias"])
	d.Set("snapshot", configExportPMap["snapshot"])
	d.Set("target_dn", configExportPMap["targetDn"])
	return d, nil
}

func resourceAciConfigurationExportPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	configExportP, err := getRemoteConfigurationExportPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setConfigurationExportPolicyAttributes(configExportP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciConfigurationExportPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ConfigurationExportPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	configExportPAttr := models.ConfigurationExportPolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		configExportPAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		configExportPAttr.Annotation = Annotation.(string)
	} else {
		configExportPAttr.Annotation = "{}"
	}
	if Format, ok := d.GetOk("format"); ok {
		configExportPAttr.Format = Format.(string)
	}
	if IncludeSecureFields, ok := d.GetOk("include_secure_fields"); ok {
		configExportPAttr.IncludeSecureFields = IncludeSecureFields.(string)
	}
	if MaxSnapshotCount, ok := d.GetOk("max_snapshot_count"); ok {
		configExportPAttr.MaxSnapshotCount = MaxSnapshotCount.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		configExportPAttr.NameAlias = NameAlias.(string)
	}
	if Snapshot, ok := d.GetOk("snapshot"); ok {
		configExportPAttr.Snapshot = Snapshot.(string)
	}
	if TargetDn, ok := d.GetOk("target_dn"); ok {
		configExportPAttr.TargetDn = TargetDn.(string)
	}
	configExportP := models.NewConfigurationExportPolicy(fmt.Sprintf("fabric/configexp-%s", name), "uni", desc, configExportPAttr)

	err := aciClient.Save(configExportP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToconfigRsExportDestination, ok := d.GetOk("relation_config_rs_export_destination"); ok {
		relationParam := relationToconfigRsExportDestination.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTotrigRsTriggerable, ok := d.GetOk("relation_trig_rs_triggerable"); ok {
		relationParam := relationTotrigRsTriggerable.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToconfigRsRemotePath, ok := d.GetOk("relation_config_rs_remote_path"); ok {
		relationParam := relationToconfigRsRemotePath.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToconfigRsExportScheduler, ok := d.GetOk("relation_config_rs_export_scheduler"); ok {
		relationParam := relationToconfigRsExportScheduler.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToconfigRsExportDestination, ok := d.GetOk("relation_config_rs_export_destination"); ok {
		relationParam := relationToconfigRsExportDestination.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationconfigRsExportDestinationFromConfigurationExportPolicy(configExportP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTotrigRsTriggerable, ok := d.GetOk("relation_trig_rs_triggerable"); ok {
		relationParam := relationTotrigRsTriggerable.(string)
		err = aciClient.CreateRelationtrigRsTriggerableFromConfigurationExportPolicy(configExportP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToconfigRsRemotePath, ok := d.GetOk("relation_config_rs_remote_path"); ok {
		relationParam := relationToconfigRsRemotePath.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationconfigRsRemotePathFromConfigurationExportPolicy(configExportP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToconfigRsExportScheduler, ok := d.GetOk("relation_config_rs_export_scheduler"); ok {
		relationParam := relationToconfigRsExportScheduler.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationconfigRsExportSchedulerFromConfigurationExportPolicy(configExportP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(configExportP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciConfigurationExportPolicyRead(ctx, d, m)
}

func resourceAciConfigurationExportPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ConfigurationExportPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	configExportPAttr := models.ConfigurationExportPolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		configExportPAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		configExportPAttr.Annotation = Annotation.(string)
	} else {
		configExportPAttr.Annotation = "{}"
	}
	if Format, ok := d.GetOk("format"); ok {
		configExportPAttr.Format = Format.(string)
	}
	if IncludeSecureFields, ok := d.GetOk("include_secure_fields"); ok {
		configExportPAttr.IncludeSecureFields = IncludeSecureFields.(string)
	}
	if MaxSnapshotCount, ok := d.GetOk("max_snapshot_count"); ok {
		configExportPAttr.MaxSnapshotCount = MaxSnapshotCount.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		configExportPAttr.NameAlias = NameAlias.(string)
	}
	if Snapshot, ok := d.GetOk("snapshot"); ok {
		configExportPAttr.Snapshot = Snapshot.(string)
	}
	if TargetDn, ok := d.GetOk("target_dn"); ok {
		configExportPAttr.TargetDn = TargetDn.(string)
	}
	configExportP := models.NewConfigurationExportPolicy(fmt.Sprintf("fabric/configexp-%s", name), "uni", desc, configExportPAttr)

	configExportP.Status = "modified"

	err := aciClient.Save(configExportP)

	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_config_rs_export_destination") {
		_, newRelParam := d.GetChange("relation_config_rs_export_destination")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_trig_rs_triggerable") {
		_, newRelParam := d.GetChange("relation_trig_rs_triggerable")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_config_rs_remote_path") {
		_, newRelParam := d.GetChange("relation_config_rs_remote_path")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_config_rs_export_scheduler") {
		_, newRelParam := d.GetChange("relation_config_rs_export_scheduler")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_config_rs_export_destination") {
		_, newRelParam := d.GetChange("relation_config_rs_export_destination")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationconfigRsExportDestinationFromConfigurationExportPolicy(configExportP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationconfigRsExportDestinationFromConfigurationExportPolicy(configExportP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_trig_rs_triggerable") {
		_, newRelParam := d.GetChange("relation_trig_rs_triggerable")
		err = aciClient.CreateRelationtrigRsTriggerableFromConfigurationExportPolicy(configExportP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_config_rs_remote_path") {
		_, newRelParam := d.GetChange("relation_config_rs_remote_path")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationconfigRsRemotePathFromConfigurationExportPolicy(configExportP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationconfigRsRemotePathFromConfigurationExportPolicy(configExportP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_config_rs_export_scheduler") {
		_, newRelParam := d.GetChange("relation_config_rs_export_scheduler")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationconfigRsExportSchedulerFromConfigurationExportPolicy(configExportP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationconfigRsExportSchedulerFromConfigurationExportPolicy(configExportP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(configExportP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciConfigurationExportPolicyRead(ctx, d, m)

}

func resourceAciConfigurationExportPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	configExportP, err := getRemoteConfigurationExportPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setConfigurationExportPolicyAttributes(configExportP, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	configRsExportDestinationData, err := aciClient.ReadRelationconfigRsExportDestinationFromConfigurationExportPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation configRsExportDestination %v", err)
		d.Set("relation_config_rs_export_destination", "")

	} else {
		d.Set("relation_config_rs_export_destination", configRsExportDestinationData.(string))
	}

	trigRsTriggerableData, err := aciClient.ReadRelationtrigRsTriggerableFromConfigurationExportPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation trigRsTriggerable %v", err)
		d.Set("relation_trig_rs_triggerable", "")

	} else {
		d.Set("relation_trig_rs_triggerable", trigRsTriggerableData.(string))
	}

	configRsRemotePathData, err := aciClient.ReadRelationconfigRsRemotePathFromConfigurationExportPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation configRsRemotePath %v", err)
		d.Set("relation_config_rs_remote_path", "")

	} else {
		d.Set("relation_config_rs_remote_path", configRsRemotePathData.(string))
	}

	configRsExportSchedulerData, err := aciClient.ReadRelationconfigRsExportSchedulerFromConfigurationExportPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation configRsExportScheduler %v", err)
		d.Set("relation_config_rs_export_scheduler", "")

	} else {
		d.Set("relation_config_rs_export_scheduler", configRsExportSchedulerData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciConfigurationExportPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "configExportP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

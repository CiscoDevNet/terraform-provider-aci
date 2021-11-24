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

func resourceAciConfigurationImportPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciConfigurationImportPolicyCreate,
		UpdateContext: resourceAciConfigurationImportPolicyUpdate,
		ReadContext:   resourceAciConfigurationImportPolicyRead,
		DeleteContext: resourceAciConfigurationImportPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciConfigurationImportPolicyImport,
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

			"fail_on_decrypt_errors": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"file_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"import_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"atomic",
					"best-effort",
				}, false),
			},

			"import_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"merge",
					"replace",
				}, false),
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

			"relation_config_rs_import_source": &schema.Schema{
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
		}),
	}
}
func getRemoteConfigurationImportPolicy(client *client.Client, dn string) (*models.ConfigurationImportPolicy, error) {
	configImportPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	configImportP := models.ConfigurationImportPolicyFromContainer(configImportPCont)

	if configImportP.DistinguishedName == "" {
		return nil, fmt.Errorf("ConfigurationImportPolicy %s not found", configImportP.DistinguishedName)
	}

	return configImportP, nil
}

func setConfigurationImportPolicyAttributes(configImportP *models.ConfigurationImportPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(configImportP.DistinguishedName)
	d.Set("description", configImportP.Description)
	configImportPMap, err := configImportP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", configImportPMap["name"])

	d.Set("admin_st", configImportPMap["adminSt"])
	d.Set("annotation", configImportPMap["annotation"])
	d.Set("fail_on_decrypt_errors", configImportPMap["failOnDecryptErrors"])
	d.Set("file_name", configImportPMap["fileName"])
	d.Set("import_mode", configImportPMap["importMode"])
	d.Set("import_type", configImportPMap["importType"])
	d.Set("name_alias", configImportPMap["nameAlias"])
	d.Set("snapshot", configImportPMap["snapshot"])
	return d, nil
}

func resourceAciConfigurationImportPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	configImportP, err := getRemoteConfigurationImportPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setConfigurationImportPolicyAttributes(configImportP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciConfigurationImportPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ConfigurationImportPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	configImportPAttr := models.ConfigurationImportPolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		configImportPAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		configImportPAttr.Annotation = Annotation.(string)
	} else {
		configImportPAttr.Annotation = "{}"
	}
	if FailOnDecryptErrors, ok := d.GetOk("fail_on_decrypt_errors"); ok {
		configImportPAttr.FailOnDecryptErrors = FailOnDecryptErrors.(string)
	}
	if FileName, ok := d.GetOk("file_name"); ok {
		configImportPAttr.FileName = FileName.(string)
	}
	if ImportMode, ok := d.GetOk("import_mode"); ok {
		configImportPAttr.ImportMode = ImportMode.(string)
	}
	if ImportType, ok := d.GetOk("import_type"); ok {
		configImportPAttr.ImportType = ImportType.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		configImportPAttr.NameAlias = NameAlias.(string)
	}
	if Snapshot, ok := d.GetOk("snapshot"); ok {
		configImportPAttr.Snapshot = Snapshot.(string)
	}
	configImportP := models.NewConfigurationImportPolicy(fmt.Sprintf("fabric/configimp-%s", name), "uni", desc, configImportPAttr)

	err := aciClient.Save(configImportP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToconfigRsImportSource, ok := d.GetOk("relation_config_rs_import_source"); ok {
		relationParam := relationToconfigRsImportSource.(string)
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

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToconfigRsImportSource, ok := d.GetOk("relation_config_rs_import_source"); ok {
		relationParam := relationToconfigRsImportSource.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationconfigRsImportSourceFromConfigurationImportPolicy(configImportP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTotrigRsTriggerable, ok := d.GetOk("relation_trig_rs_triggerable"); ok {
		relationParam := relationTotrigRsTriggerable.(string)
		err = aciClient.CreateRelationtrigRsTriggerableFromConfigurationImportPolicy(configImportP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToconfigRsRemotePath, ok := d.GetOk("relation_config_rs_remote_path"); ok {
		relationParam := relationToconfigRsRemotePath.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationconfigRsRemotePathFromConfigurationImportPolicy(configImportP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(configImportP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciConfigurationImportPolicyRead(ctx, d, m)
}

func resourceAciConfigurationImportPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ConfigurationImportPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	configImportPAttr := models.ConfigurationImportPolicyAttributes{}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		configImportPAttr.AdminSt = AdminSt.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		configImportPAttr.Annotation = Annotation.(string)
	} else {
		configImportPAttr.Annotation = "{}"
	}
	if FailOnDecryptErrors, ok := d.GetOk("fail_on_decrypt_errors"); ok {
		configImportPAttr.FailOnDecryptErrors = FailOnDecryptErrors.(string)
	}
	if FileName, ok := d.GetOk("file_name"); ok {
		configImportPAttr.FileName = FileName.(string)
	}
	if ImportMode, ok := d.GetOk("import_mode"); ok {
		configImportPAttr.ImportMode = ImportMode.(string)
	}
	if ImportType, ok := d.GetOk("import_type"); ok {
		configImportPAttr.ImportType = ImportType.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		configImportPAttr.NameAlias = NameAlias.(string)
	}
	if Snapshot, ok := d.GetOk("snapshot"); ok {
		configImportPAttr.Snapshot = Snapshot.(string)
	}
	configImportP := models.NewConfigurationImportPolicy(fmt.Sprintf("fabric/configimp-%s", name), "uni", desc, configImportPAttr)

	configImportP.Status = "modified"

	err := aciClient.Save(configImportP)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_config_rs_import_source") {
		_, newRelParam := d.GetChange("relation_config_rs_import_source")
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

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_config_rs_import_source") {
		_, newRelParam := d.GetChange("relation_config_rs_import_source")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationconfigRsImportSourceFromConfigurationImportPolicy(configImportP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationconfigRsImportSourceFromConfigurationImportPolicy(configImportP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_trig_rs_triggerable") {
		_, newRelParam := d.GetChange("relation_trig_rs_triggerable")
		err = aciClient.CreateRelationtrigRsTriggerableFromConfigurationImportPolicy(configImportP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_config_rs_remote_path") {
		_, newRelParam := d.GetChange("relation_config_rs_remote_path")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationconfigRsRemotePathFromConfigurationImportPolicy(configImportP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationconfigRsRemotePathFromConfigurationImportPolicy(configImportP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(configImportP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciConfigurationImportPolicyRead(ctx, d, m)

}

func resourceAciConfigurationImportPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	configImportP, err := getRemoteConfigurationImportPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setConfigurationImportPolicyAttributes(configImportP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	configRsImportSourceData, err := aciClient.ReadRelationconfigRsImportSourceFromConfigurationImportPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation configRsImportSource %v", err)
		d.Set("relation_config_rs_import_source", "")

	} else {
		d.Set("relation_config_rs_import_source", configRsImportSourceData.(string))
	}

	trigRsTriggerableData, err := aciClient.ReadRelationtrigRsTriggerableFromConfigurationImportPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation trigRsTriggerable %v", err)
		d.Set("relation_trig_rs_triggerable", "")

	} else {
		d.Set("relation_trig_rs_triggerable", trigRsTriggerableData.(string))
	}

	configRsRemotePathData, err := aciClient.ReadRelationconfigRsRemotePathFromConfigurationImportPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation configRsRemotePath %v", err)
		d.Set("relation_config_rs_remote_path", "")

	} else {
		d.Set("relation_config_rs_remote_path", configRsRemotePathData.(string))

	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciConfigurationImportPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "configImportP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

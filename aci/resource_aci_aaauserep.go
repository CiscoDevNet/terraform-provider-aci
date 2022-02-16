package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciUserManagement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciUserManagementCreate,
		UpdateContext: resourceAciUserManagementUpdate,
		ReadContext:   resourceAciUserManagementRead,
		DeleteContext: resourceAciUserManagementDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciUserManagementImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"pwd_strength_check": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"change_count": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"change_during_interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disable",
					"enable",
				}, false),
			},
			"change_interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expiration_warn_time": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"history_count": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"no_change_interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"block_duration": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_login_block": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disable",
					"enable",
				}, false),
			},
			"max_failed_attempts": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_failed_attempts_window": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maximum_validity_period": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"session_record_flags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"login",
						"logout",
						"refresh",
					}, false),
				},
			},
			"ui_idle_timeout_seconds": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"webtoken_timeout_seconds": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_aaa_rs_to_user_ep": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to aaa:UserEp",
			},
		})),
	}
}

func GetRemoteUserManagement(client *client.Client, dn string) (*models.UserManagement, error) {
	aaaUserEpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaUserEp := models.UserManagementFromContainer(aaaUserEpCont)
	if aaaUserEp.DistinguishedName == "" {
		return nil, fmt.Errorf("UserManagement %s not found", aaaUserEp.DistinguishedName)
	}
	return aaaUserEp, nil
}

func setUserManagementAttributes(aaaUserEp *models.UserManagement, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaUserEp.DistinguishedName)
	d.Set("description", aaaUserEp.Description)
	aaaUserEpMap, err := aaaUserEp.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", aaaUserEpMap["annotation"])
	d.Set("pwd_strength_check", aaaUserEpMap["pwdStrengthCheck"])
	d.Set("name_alias", aaaUserEpMap["nameAlias"])
	return d, nil
}

func GetRemotePasswordChangeExpirationPolicy(client *client.Client, dn string) (*models.PasswordChangeExpirationPolicy, error) {
	aaaPwdProfileCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaPwdProfile := models.PasswordChangeExpirationPolicyFromContainer(aaaPwdProfileCont)
	if aaaPwdProfile.DistinguishedName == "" {
		return nil, fmt.Errorf("PasswordChangeExpirationPolicy %s not found", aaaPwdProfile.DistinguishedName)
	}
	return aaaPwdProfile, nil
}

func setPasswordChangeExpirationPolicyAttributes(aaaPwdProfile *models.PasswordChangeExpirationPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	aaaPwdProfileMap, err := aaaPwdProfile.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("change_count", aaaPwdProfileMap["changeCount"])
	d.Set("change_during_interval", aaaPwdProfileMap["changeDuringInterval"])
	d.Set("change_interval", aaaPwdProfileMap["changeInterval"])
	d.Set("expiration_warn_time", aaaPwdProfileMap["expirationWarnTime"])
	d.Set("history_count", aaaPwdProfileMap["historyCount"])
	d.Set("no_change_interval", aaaPwdProfileMap["noChangeInterval"])
	return d, nil
}

func GetRemoteBlockUserLoginsPolicy(client *client.Client, dn string) (*models.BlockUserLoginsPolicy, error) {
	aaaBlockLoginProfileCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaBlockLoginProfile := models.BlockUserLoginsPolicyFromContainer(aaaBlockLoginProfileCont)
	if aaaBlockLoginProfile.DistinguishedName == "" {
		return nil, fmt.Errorf("BlockUserLoginsPolicy %s not found", aaaBlockLoginProfile.DistinguishedName)
	}
	return aaaBlockLoginProfile, nil
}

func setBlockUserLoginsPolicyAttributes(aaaBlockLoginProfile *models.BlockUserLoginsPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {

	aaaBlockLoginProfileMap, err := aaaBlockLoginProfile.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("block_duration", aaaBlockLoginProfileMap["blockDuration"])
	d.Set("enable_login_block", aaaBlockLoginProfileMap["enableLoginBlock"])
	d.Set("max_failed_attempts", aaaBlockLoginProfileMap["maxFailedAttempts"])
	d.Set("max_failed_attempts_window", aaaBlockLoginProfileMap["maxFailedAttemptsWindow"])
	return d, nil
}

func GetRemoteWebTokenData(client *client.Client, dn string) (*models.WebTokenData, error) {
	pkiWebTokenDataCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	pkiWebTokenData := models.WebTokenDataFromContainer(pkiWebTokenDataCont)
	if pkiWebTokenData.DistinguishedName == "" {
		return nil, fmt.Errorf("WebTokenData %s not found", pkiWebTokenData.DistinguishedName)
	}
	return pkiWebTokenData, nil
}

func setWebTokenDataAttributes(pkiWebTokenData *models.WebTokenData, d *schema.ResourceData) (*schema.ResourceData, error) {

	pkiWebTokenDataMap, err := pkiWebTokenData.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("maximum_validity_period", pkiWebTokenDataMap["maximumValidityPeriod"])
	sessionRecordFlagsGet := make([]string, 0, 1)
	for _, val := range strings.Split(pkiWebTokenDataMap["sessionRecordFlags"], ",") {
		sessionRecordFlagsGet = append(sessionRecordFlagsGet, strings.Trim(val, " "))
	}
	sort.Strings(sessionRecordFlagsGet)
	if sessionRecordFlagsIntr, ok := d.GetOk("session_record_flags"); ok {
		sessionRecordFlagsAct := make([]string, 0, 1)
		for _, val := range sessionRecordFlagsIntr.([]interface{}) {
			sessionRecordFlagsAct = append(sessionRecordFlagsAct, val.(string))
		}
		sort.Strings(sessionRecordFlagsAct)
		if reflect.DeepEqual(sessionRecordFlagsAct, sessionRecordFlagsGet) {
			d.Set("session_record_flags", d.Get("session_record_flags").([]interface{}))
		} else {
			d.Set("session_record_flags", sessionRecordFlagsGet)
		}
	} else {
		d.Set("session_record_flags", sessionRecordFlagsGet)
	}
	d.Set("ui_idle_timeout_seconds", pkiWebTokenDataMap["uiIdleTimeoutSeconds"])
	d.Set("webtoken_timeout_seconds", pkiWebTokenDataMap["webtokenTimeoutSeconds"])
	return d, nil
}

func resourceAciUserManagementImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaUserEp, err := GetRemoteUserManagement(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setUserManagementAttributes(aaaUserEp, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciUserManagementCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] UserManagement: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	aaaUserEpAttr := models.UserManagementAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaUserEpAttr.Annotation = Annotation.(string)
	} else {
		aaaUserEpAttr.Annotation = "{}"
	}

	if PwdStrengthCheck, ok := d.GetOk("pwd_strength_check"); ok {
		aaaUserEpAttr.PwdStrengthCheck = PwdStrengthCheck.(string)
	}
	aaaUserEp := models.NewUserManagement(fmt.Sprintf("userext"), "uni", desc, nameAlias, aaaUserEpAttr)
	aaaUserEp.Status = "modified"
	err := aciClient.Save(aaaUserEp)
	if err != nil {
		return diag.FromErr(err)
	}

	aaaPwdProfileAttr := models.PasswordChangeExpirationPolicyAttributes{}
	aaaPwdProfileOk := false
	if ChangeCount, ok := d.GetOk("change_count"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.ChangeCount = ChangeCount.(string)
	}

	if ChangeDuringInterval, ok := d.GetOk("change_during_interval"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.ChangeDuringInterval = ChangeDuringInterval.(string)
	}

	if ChangeInterval, ok := d.GetOk("change_interval"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.ChangeInterval = ChangeInterval.(string)
	}

	if ExpirationWarnTime, ok := d.GetOk("expiration_warn_time"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.ExpirationWarnTime = ExpirationWarnTime.(string)
	}

	if HistoryCount, ok := d.GetOk("history_count"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.HistoryCount = HistoryCount.(string)
	}

	if NoChangeInterval, ok := d.GetOk("no_change_interval"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.NoChangeInterval = NoChangeInterval.(string)
	}

	if aaaPwdProfileOk {
		aaaPwdProfile := models.NewPasswordChangeExpirationPolicy(fmt.Sprintf("userext/pwdprofile"), "uni", desc, nameAlias, aaaPwdProfileAttr)
		aaaPwdProfile.Status = "modified"
		err = aciClient.Save(aaaPwdProfile)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	aaaBlockLoginProfileAttr := models.BlockUserLoginsPolicyAttributes{}
	aaaBlockLoginProfileOk := false
	if BlockDuration, ok := d.GetOk("block_duration"); ok {
		aaaBlockLoginProfileOk = ok
		aaaBlockLoginProfileAttr.BlockDuration = BlockDuration.(string)
	}

	if EnableLoginBlock, ok := d.GetOk("enable_login_block"); ok {
		aaaBlockLoginProfileOk = ok
		aaaBlockLoginProfileAttr.EnableLoginBlock = EnableLoginBlock.(string)
	}

	if MaxFailedAttempts, ok := d.GetOk("max_failed_attempts"); ok {
		aaaBlockLoginProfileOk = ok
		aaaBlockLoginProfileAttr.MaxFailedAttempts = MaxFailedAttempts.(string)
	}

	if MaxFailedAttemptsWindow, ok := d.GetOk("max_failed_attempts_window"); ok {
		aaaBlockLoginProfileOk = ok
		aaaBlockLoginProfileAttr.MaxFailedAttemptsWindow = MaxFailedAttemptsWindow.(string)
	}

	if aaaBlockLoginProfileOk {
		aaaBlockLoginProfile := models.NewBlockUserLoginsPolicy(fmt.Sprintf("userext/blockloginp"), "uni", desc, nameAlias, aaaBlockLoginProfileAttr)
		aaaBlockLoginProfile.Status = "modified"
		err = aciClient.Save(aaaBlockLoginProfile)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	pkiEpAttr := models.PublicKeyManagementAttributes{}
	pkiEp := models.NewPublicKeyManagement(fmt.Sprintf("userext/pkiext"), "uni", desc, nameAlias, pkiEpAttr)
	err = aciClient.Save(pkiEp)
	if err != nil {
		return diag.FromErr(err)
	}

	pkiWebTokenDataAttr := models.WebTokenDataAttributes{}

	pkiWebTokenDataOk := false

	if MaximumValidityPeriod, ok := d.GetOk("maximum_validity_period"); ok {
		pkiWebTokenDataOk = ok
		pkiWebTokenDataAttr.MaximumValidityPeriod = MaximumValidityPeriod.(string)
	}

	if SessionRecordFlags, ok := d.GetOk("session_record_flags"); ok {
		pkiWebTokenDataOk = ok
		sessionRecordFlagsList := make([]string, 0, 1)
		for _, val := range SessionRecordFlags.([]interface{}) {
			sessionRecordFlagsList = append(sessionRecordFlagsList, val.(string))
		}
		err := checkDuplicate(sessionRecordFlagsList)
		if err != nil {
			return diag.FromErr(err)
		}
		SessionRecordFlags := strings.Join(sessionRecordFlagsList, ",")
		pkiWebTokenDataAttr.SessionRecordFlags = SessionRecordFlags
	}

	if UiIdleTimeoutSeconds, ok := d.GetOk("ui_idle_timeout_seconds"); ok {
		pkiWebTokenDataOk = ok
		pkiWebTokenDataAttr.UiIdleTimeoutSeconds = UiIdleTimeoutSeconds.(string)
	}

	if WebtokenTimeoutSeconds, ok := d.GetOk("webtoken_timeout_seconds"); ok {
		pkiWebTokenDataOk = ok
		pkiWebTokenDataAttr.WebtokenTimeoutSeconds = WebtokenTimeoutSeconds.(string)
	}
	if pkiWebTokenDataOk {

		pkiWebTokenData := models.NewWebTokenData(fmt.Sprintf("userext/pkiext/webtokendata"), "uni", desc, nameAlias, pkiWebTokenDataAttr)
		pkiWebTokenData.Status = "modified"
		err := aciClient.Save(pkiWebTokenData)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	checkDns := make([]string, 0, 1)

	if relationToaaaRsToUserEp, ok := d.GetOk("relation_aaa_rs_to_user_ep"); ok {
		relationParam := relationToaaaRsToUserEp.(string)
		checkDns = append(checkDns, relationParam)

	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationToaaaRsToUserEp, ok := d.GetOk("relation_aaa_rs_to_user_ep"); ok {
		relationParam := relationToaaaRsToUserEp.(string)
		err = aciClient.CreateRelationaaaRsToUserEp(aaaUserEp.DistinguishedName, aaaUserEpAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(aaaUserEp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciUserManagementRead(ctx, d, m)
}

func resourceAciUserManagementUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] UserManagement: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	aaaUserEpAttr := models.UserManagementAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaUserEpAttr.Annotation = Annotation.(string)
	} else {
		aaaUserEpAttr.Annotation = "{}"
	}

	if PwdStrengthCheck, ok := d.GetOk("pwd_strength_check"); ok {
		aaaUserEpAttr.PwdStrengthCheck = PwdStrengthCheck.(string)
	}
	aaaUserEp := models.NewUserManagement(fmt.Sprintf("userext"), "uni", desc, nameAlias, aaaUserEpAttr)
	aaaUserEp.Status = "modified"
	err := aciClient.Save(aaaUserEp)
	if err != nil {
		return diag.FromErr(err)
	}

	aaaPwdProfileAttr := models.PasswordChangeExpirationPolicyAttributes{}
	aaaPwdProfileOk := false
	if ChangeCount, ok := d.GetOk("change_count"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.ChangeCount = ChangeCount.(string)
	}

	if ChangeDuringInterval, ok := d.GetOk("change_during_interval"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.ChangeDuringInterval = ChangeDuringInterval.(string)
	}

	if ChangeInterval, ok := d.GetOk("change_interval"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.ChangeInterval = ChangeInterval.(string)
	}

	if ExpirationWarnTime, ok := d.GetOk("expiration_warn_time"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.ExpirationWarnTime = ExpirationWarnTime.(string)
	}

	if HistoryCount, ok := d.GetOk("history_count"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.HistoryCount = HistoryCount.(string)
	}

	if NoChangeInterval, ok := d.GetOk("no_change_interval"); ok {
		aaaPwdProfileOk = ok
		aaaPwdProfileAttr.NoChangeInterval = NoChangeInterval.(string)
	}

	if aaaPwdProfileOk {
		aaaPwdProfile := models.NewPasswordChangeExpirationPolicy(fmt.Sprintf("userext/pwdprofile"), "uni", desc, nameAlias, aaaPwdProfileAttr)
		aaaPwdProfile.Status = "modified"
		err = aciClient.Save(aaaPwdProfile)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	aaaBlockLoginProfileAttr := models.BlockUserLoginsPolicyAttributes{}
	aaaBlockLoginProfileOk := false
	if BlockDuration, ok := d.GetOk("block_duration"); ok {
		aaaBlockLoginProfileOk = ok
		aaaBlockLoginProfileAttr.BlockDuration = BlockDuration.(string)
	}

	if EnableLoginBlock, ok := d.GetOk("enable_login_block"); ok {
		aaaBlockLoginProfileOk = ok
		aaaBlockLoginProfileAttr.EnableLoginBlock = EnableLoginBlock.(string)
	}

	if MaxFailedAttempts, ok := d.GetOk("max_failed_attempts"); ok {
		aaaBlockLoginProfileOk = ok
		aaaBlockLoginProfileAttr.MaxFailedAttempts = MaxFailedAttempts.(string)
	}

	if MaxFailedAttemptsWindow, ok := d.GetOk("max_failed_attempts_window"); ok {
		aaaBlockLoginProfileOk = ok
		aaaBlockLoginProfileAttr.MaxFailedAttemptsWindow = MaxFailedAttemptsWindow.(string)
	}

	if aaaBlockLoginProfileOk {
		aaaBlockLoginProfile := models.NewBlockUserLoginsPolicy(fmt.Sprintf("userext/blockloginp"), "uni", desc, nameAlias, aaaBlockLoginProfileAttr)
		aaaBlockLoginProfile.Status = "modified"
		err = aciClient.Save(aaaBlockLoginProfile)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	pkiEpAttr := models.PublicKeyManagementAttributes{}
	pkiEp := models.NewPublicKeyManagement(fmt.Sprintf("userext/pkiext"), "uni", desc, nameAlias, pkiEpAttr)
	err = aciClient.Save(pkiEp)
	if err != nil {
		return diag.FromErr(err)
	}

	pkiWebTokenDataAttr := models.WebTokenDataAttributes{}

	pkiWebTokenDataOk := false

	if MaximumValidityPeriod, ok := d.GetOk("maximum_validity_period"); ok {
		pkiWebTokenDataOk = ok
		pkiWebTokenDataAttr.MaximumValidityPeriod = MaximumValidityPeriod.(string)
	}

	if SessionRecordFlags, ok := d.GetOk("session_record_flags"); ok {
		pkiWebTokenDataOk = ok
		sessionRecordFlagsList := make([]string, 0, 1)
		for _, val := range SessionRecordFlags.([]interface{}) {
			sessionRecordFlagsList = append(sessionRecordFlagsList, val.(string))
		}
		err := checkDuplicate(sessionRecordFlagsList)
		if err != nil {
			return diag.FromErr(err)
		}
		SessionRecordFlags := strings.Join(sessionRecordFlagsList, ",")
		pkiWebTokenDataAttr.SessionRecordFlags = SessionRecordFlags
	}

	if UiIdleTimeoutSeconds, ok := d.GetOk("ui_idle_timeout_seconds"); ok {
		pkiWebTokenDataOk = ok
		pkiWebTokenDataAttr.UiIdleTimeoutSeconds = UiIdleTimeoutSeconds.(string)
	}

	if WebtokenTimeoutSeconds, ok := d.GetOk("webtoken_timeout_seconds"); ok {
		pkiWebTokenDataOk = ok
		pkiWebTokenDataAttr.WebtokenTimeoutSeconds = WebtokenTimeoutSeconds.(string)
	}

	if pkiWebTokenDataOk {
		pkiWebTokenData := models.NewWebTokenData(fmt.Sprintf("userext/pkiext/webtokendata"), "uni", desc, nameAlias, pkiWebTokenDataAttr)
		pkiWebTokenData.Status = "modified"
		err := aciClient.Save(pkiWebTokenData)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_aaa_rs_to_user_ep") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_aaa_rs_to_user_ep")
		checkDns = append(checkDns, newRelParam.(string))

	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("relation_aaa_rs_to_user_ep") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_aaa_rs_to_user_ep")
		err = aciClient.DeleteRelationaaaRsToUserEp(aaaUserEp.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationaaaRsToUserEp(aaaUserEp.DistinguishedName, aaaUserEpAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(aaaUserEp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciUserManagementRead(ctx, d, m)
}

func resourceAciUserManagementRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaUserEp, err := GetRemoteUserManagement(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setUserManagementAttributes(aaaUserEp, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = aciClient.Get(dn + "/pwdprofile")
	if err == nil {
		aaaPwdProfileDn := dn + "/pwdprofile"
		aaaPwdProfile, err := GetRemotePasswordChangeExpirationPolicy(aciClient, aaaPwdProfileDn)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = setPasswordChangeExpirationPolicyAttributes(aaaPwdProfile, d)
		if err != nil {
			return nil
		}
	}

	_, err = aciClient.Get(dn + "/blockloginp")
	if err == nil {
		aaaBlockLoginProfileDn := dn + "/blockloginp"
		aaaBlockLoginProfile, err := GetRemoteBlockUserLoginsPolicy(aciClient, aaaBlockLoginProfileDn)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = setBlockUserLoginsPolicyAttributes(aaaBlockLoginProfile, d)
		if err != nil {
			return nil
		}
	}

	_, err = aciClient.Get(dn + "/pkiext/webtokendata")
	if err == nil {
		pkiWebTokenDn := dn + "/pkiext/webtokendata"
		pkiWebTokenData, err := GetRemoteWebTokenData(aciClient, pkiWebTokenDn)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = setWebTokenDataAttributes(pkiWebTokenData, d)
		if err != nil {
			return nil
		}
	}

	aaaRsToUserEpData, err := aciClient.ReadRelationaaaRsToUserEp(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation aaaRsToUserEp %v", err)
		d.Set("relation_aaa_rs_to_user_ep", "")
	} else {
		setRelationAttribute(d, "relation_aaa_rs_to_user_ep", aaaRsToUserEpData)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciUserManagementDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name aaaUserEp cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diags
}

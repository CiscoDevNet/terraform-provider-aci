package aci

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciCloudAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudAccountCreate,
		UpdateContext: resourceAciCloudAccountUpdate,
		ReadContext:   resourceAciCloudAccountRead,
		DeleteContext: resourceAciCloudAccountDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudAccountImport,
		},

		SchemaVersion: 1,

		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
			if diff.Get("access_type") == "credentials" {
				if diff.Get("cloud_credentials_dn") == "" {
					return errors.New(`"cloud_credentials_dn" is required when "access_type" is credentials`)
				}
			}

			return nil
		},

		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"credentials",
					"managed",
				}, false),
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"aws",
					"azure",
					"gcp",
					"unknown",
				}, false),
			},
			"relation_cloud_rs_account_to_access_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to cloud:AccessPolicy",
			},
			"cloud_credentials_dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to cloud:Credentials",
			}})),
	}
}

func getRemoteCloudAccount(client *client.Client, dn string) (*models.CloudAccount, error) {
	cloudAccountCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudAccount := models.CloudAccountFromContainer(cloudAccountCont)
	if cloudAccount.DistinguishedName == "" {
		return nil, fmt.Errorf("Cloud Account %s not found", cloudAccount.DistinguishedName)
	}
	return cloudAccount, nil
}

func setCloudAccountAttributes(cloudAccount *models.CloudAccount, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudAccount.DistinguishedName)
	if dn != cloudAccount.DistinguishedName {
		d.Set("tenant_dn", "")
	} else {
		d.Set("tenant_dn", GetParentDn(dn, "/ad"))
	}

	cloudAccountMap, err := cloudAccount.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("access_type", cloudAccountMap["accessType"])
	d.Set("account_id", cloudAccountMap["id"])
	d.Set("name", cloudAccountMap["name"])
	d.Set("vendor", cloudAccountMap["vendor"])
	d.Set("name_alias", cloudAccountMap["nameAlias"])
	d.Set("annotation", cloudAccountMap["annotation"])

	return d, nil
}

func resourceAciCloudAccountImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudAccount, err := getRemoteCloudAccount(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCloudAccountAttributes(cloudAccount, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Account: Beginning Creation")
	aciClient := m.(*client.Client)
	account_id := d.Get("account_id").(string)
	vendor := d.Get("vendor").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudAccountSet := make([]interface{}, 0, 1)
	cloudAccountMap := make(map[string]interface{})

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudAccountMap["nameAlias"] = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudAccountMap["annotation"] = Annotation.(string)
	} else {
		cloudAccountMap["annotation"] = "{}"
	}

	if AccessType, ok := d.GetOk("access_type"); ok {
		cloudAccountMap["accessType"] = AccessType.(string)
	}

	if Account_id, ok := d.GetOk("account_id"); ok {
		cloudAccountMap["id"] = Account_id.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudAccountMap["name"] = Name.(string)
	}

	if Vendor, ok := d.GetOk("vendor"); ok {
		cloudAccountMap["vendor"] = Vendor.(string)
	}

	dn := fmt.Sprintf("%s/%s", TenantDn, fmt.Sprintf(models.RncloudAccount, account_id, vendor))

	if relationTocloudRsCredentials, ok := d.GetOk("cloud_credentials_dn"); ok {

		cloudAccountCredentialsMap := make(map[string]interface{})
		cloudAccountCredentialsMap["class_name"] = "cloudRsCredentials"
		cloudAccountCredentialsContent := make(map[string]interface{})

		cloudAccountCredentialsContent["tDn"] = relationTocloudRsCredentials

		cloudAccountCredentialsMap["content"] = toStrMap(cloudAccountCredentialsContent)
		cloudAccountSet = append(cloudAccountSet, cloudAccountCredentialsMap)
	}

	cont, err := preparePayload(models.CloudaccountClassName, toStrMap(cloudAccountMap), cloudAccountSet)
	if err != nil {
		return diag.FromErr(err)
	}

	_, diags := cloudAccountRequest(aciClient, "POST", dn, cont)
	if diags.HasError() {
		return diags
	}

	checkDns := make([]string, 0, 1)

	if relationTocloudRsAccountToAccessPolicy, ok := d.GetOk("relation_cloud_rs_account_to_access_policy"); ok {
		relationParam := relationTocloudRsAccountToAccessPolicy.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTocloudRsAccountToAccessPolicy, ok := d.GetOk("relation_cloud_rs_account_to_access_policy"); ok {
		relationParam := relationTocloudRsAccountToAccessPolicy.(string)
		err = aciClient.CreateRelationcloudRsAccountToAccessPolicy(dn, "", relationParam)
		// err = aciClient.CreateRelationcloudRsAccountToAccessPolicy(cloudAccount.DistinguishedName, cloudAccountAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(dn)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudAccountRead(ctx, d, m)
}

func resourceAciCloudAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Account: Beginning Update")
	aciClient := m.(*client.Client)
	account_id := d.Get("account_id").(string)
	vendor := d.Get("vendor").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudAccountSet := make([]interface{}, 0, 1)
	cloudAccountMap := make(map[string]interface{})

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudAccountMap["nameAlias"] = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudAccountMap["annotation"] = Annotation.(string)
	} else {
		cloudAccountMap["annotation"] = "{}"
	}

	if AccessType, ok := d.GetOk("access_type"); ok {
		cloudAccountMap["accessType"] = AccessType.(string)
	}

	if Account_id, ok := d.GetOk("account_id"); ok {
		cloudAccountMap["id"] = Account_id.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudAccountMap["name"] = Name.(string)
	}

	if Vendor, ok := d.GetOk("vendor"); ok {
		cloudAccountMap["vendor"] = Vendor.(string)
	}

	dn := fmt.Sprintf("%s/%s", TenantDn, fmt.Sprintf(models.RncloudAccount, account_id, vendor))

	if relationTocloudRsCredentials, ok := d.GetOk("cloud_credentials_dn"); ok {

		cloudAccountCredentialsMap := make(map[string]interface{})
		cloudAccountCredentialsMap["class_name"] = "cloudRsCredentials"
		cloudAccountCredentialsContent := make(map[string]interface{})

		cloudAccountCredentialsContent["tDn"] = relationTocloudRsCredentials

		cloudAccountCredentialsMap["content"] = toStrMap(cloudAccountCredentialsContent)
		cloudAccountSet = append(cloudAccountSet, cloudAccountCredentialsMap)
	}

	cont, err := preparePayload(models.CloudaccountClassName, toStrMap(cloudAccountMap), cloudAccountSet)
	if err != nil {
		return diag.FromErr(err)
	}

	_, diags := cloudAccountRequest(aciClient, "POST", dn, cont)
	if diags.HasError() {
		return diags
	}
	// cloudAccount.Status = "modified"

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_cloud_rs_account_to_access_policy") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_cloud_rs_account_to_access_policy")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_cloud_rs_account_to_access_policy") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_cloud_rs_account_to_access_policy")
		err = aciClient.DeleteRelationcloudRsAccountToAccessPolicy(dn)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationcloudRsAccountToAccessPolicy(dn, "", newRelParam.(string))
		// err = aciClient.CreateRelationcloudRsAccountToAccessPolicy(cloudAccount.DistinguishedName, cloudAccountAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(dn)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudAccountRead(ctx, d, m)
}

func resourceAciCloudAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudAccount, err := getRemoteCloudAccount(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setCloudAccountAttributes(cloudAccount, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	cloudRsAccountToAccessPolicyData, err := aciClient.ReadRelationcloudRsAccountToAccessPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation cloudRsAccountToAccessPolicy %v", err)
		d.Set("relation_cloud_rs_account_to_access_policy", "")
	} else {
		if _, ok := d.GetOk("relation_cloud_rs_account_to_access_policy"); ok {
			tfName := d.Get("relation_cloud_rs_account_to_access_policy").(string)
			if tfName != cloudRsAccountToAccessPolicyData {
				d.Set("relation_cloud_rs_account_to_access_policy", "")
			}
		}
	}

	cloudRsCredentialsData, err := aciClient.ReadRelationcloudRsCredentials(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation cloudRsCredentials %v", err)
		d.Set("cloud_credentials_dn", "")
	} else {
		if _, ok := d.GetOk("cloud_credentials_dn"); ok {
			tfName := d.Get("cloud_credentials_dn").(string)
			if tfName != cloudRsCredentialsData {
				d.Set("cloud_credentials_dn", "")
			}
		}
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "cloudAccount")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

func cloudAccountRequest(aciClient *client.Client, method string, cloudAccountdn string, body *container.Container) (*container.Container, diag.Diagnostics) {
	url := "/api/mo/" + cloudAccountdn + ".json"
	req, err := aciClient.MakeRestRequest(method, url, body, true)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	respCont, _, err := aciClient.Do(req)
	if err != nil {
		return respCont, diag.FromErr(err)
	}
	err = client.CheckForErrors(respCont, method, false)
	if err != nil {
		return respCont, diag.FromErr(err)
	}
	if method == "POST" {
		return body, nil
	} else {
		return respCont, nil
	}
}

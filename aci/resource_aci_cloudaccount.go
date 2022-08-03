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
			"relation_cloud_rs_credentials": {
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
		d.Set("tenant_dn", GetParentDn(dn, "/act"))
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

	cloudAccountAttr := models.CloudAccountAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudAccountAttr.Annotation = Annotation.(string)
	} else {
		cloudAccountAttr.Annotation = "{}"
	}

	if AccessType, ok := d.GetOk("access_type"); ok {
		cloudAccountAttr.AccessType = AccessType.(string)
	}

	if Account_id, ok := d.GetOk("account_id"); ok {
		cloudAccountAttr.Account_id = Account_id.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudAccountAttr.Name = Name.(string)
	}

	if Vendor, ok := d.GetOk("vendor"); ok {
		cloudAccountAttr.Vendor = Vendor.(string)
	}
	cloudAccount := models.NewCloudAccount(fmt.Sprintf(models.RncloudAccount, account_id, vendor), TenantDn, nameAlias, cloudAccountAttr)

	err := aciClient.Save(cloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTocloudRsAccountToAccessPolicy, ok := d.GetOk("relation_cloud_rs_account_to_access_policy"); ok {
		relationParam := relationTocloudRsAccountToAccessPolicy.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTocloudRsCredentials, ok := d.GetOk("relation_cloud_rs_credentials"); ok {
		relationParam := relationTocloudRsCredentials.(string)
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
		err = aciClient.CreateRelationcloudRsAccountToAccessPolicy(cloudAccount.DistinguishedName, cloudAccountAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTocloudRsCredentials, ok := d.GetOk("relation_cloud_rs_credentials"); ok {
		relationParam := relationTocloudRsCredentials.(string)
		err = aciClient.CreateRelationcloudRsCredentials(cloudAccount.DistinguishedName, cloudAccountAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(cloudAccount.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudAccountRead(ctx, d, m)
}

func resourceAciCloudAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Account: Beginning Update")
	aciClient := m.(*client.Client)
	account_id := d.Get("account_id").(string)
	vendor := d.Get("vendor").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudAccountAttr := models.CloudAccountAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudAccountAttr.Annotation = Annotation.(string)
	} else {
		cloudAccountAttr.Annotation = "{}"
	}

	if AccessType, ok := d.GetOk("access_type"); ok {
		cloudAccountAttr.AccessType = AccessType.(string)
	}

	if Account_id, ok := d.GetOk("account_id"); ok {
		cloudAccountAttr.Account_id = Account_id.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudAccountAttr.Name = Name.(string)
	}

	if Vendor, ok := d.GetOk("vendor"); ok {
		cloudAccountAttr.Vendor = Vendor.(string)
	}
	cloudAccount := models.NewCloudAccount(fmt.Sprintf(models.RncloudAccount, account_id, vendor), TenantDn, nameAlias, cloudAccountAttr)

	cloudAccount.Status = "modified"

	err := aciClient.Save(cloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_cloud_rs_account_to_access_policy") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_cloud_rs_account_to_access_policy")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_cloud_rs_credentials") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_cloud_rs_credentials")
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
		err = aciClient.DeleteRelationcloudRsAccountToAccessPolicy(cloudAccount.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationcloudRsAccountToAccessPolicy(cloudAccount.DistinguishedName, cloudAccountAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_cloud_rs_credentials") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_cloud_rs_credentials")
		err = aciClient.DeleteRelationcloudRsCredentials(cloudAccount.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationcloudRsCredentials(cloudAccount.DistinguishedName, cloudAccountAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(cloudAccount.DistinguishedName)
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
		d.Set("relation_cloud_rs_credentials", "")
	} else {
		if _, ok := d.GetOk("relation_cloud_rs_credentials"); ok {
			tfName := d.Get("relation_cloud_rs_credentials").(string)
			if tfName != cloudRsCredentialsData {
				d.Set("relation_cloud_rs_credentials", "")
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

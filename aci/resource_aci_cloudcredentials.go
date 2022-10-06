package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciCloudCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudCredentialsCreate,
		UpdateContext: resourceAciCloudCredentialsUpdate,
		ReadContext:   resourceAciCloudCredentialsRead,
		DeleteContext: resourceAciCloudCredentialsDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudCredentialsImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email": { //only for gcp
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_proxy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_id": { //required for both (azure and gcp)
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rsa_private_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_cloud_rs_ad": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Create relation to cloud:AD",
			}})),
	}
}

func getRemoteCloudCredentials(client *client.Client, dn string) (*models.CloudCredentials, error) {
	cloudCredentialsCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudCredentials := models.CloudCredentialsFromContainer(cloudCredentialsCont)
	if cloudCredentials.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudCredentials %s not found", cloudCredentials.DistinguishedName)
	}
	return cloudCredentials, nil
}

func setCloudCredentialsAttributes(cloudCredentials *models.CloudCredentials, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudCredentials.DistinguishedName)
	d.Set("description", cloudCredentials.Description)
	if dn != cloudCredentials.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	cloudCredentialsMap, err := cloudCredentials.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("client_id", cloudCredentialsMap["clientId"])
	d.Set("email", cloudCredentialsMap["email"])
	d.Set("http_proxy", cloudCredentialsMap["httpProxy"])
	d.Set("key", cloudCredentialsMap["key"])
	d.Set("key_id", cloudCredentialsMap["keyId"])
	d.Set("name", cloudCredentialsMap["name"])
	d.Set("rsa_private_key", cloudCredentialsMap["rsaPrivateKey"])
	d.Set("name_alias", cloudCredentialsMap["nameAlias"])
	return d, nil
}

func resourceAciCloudCredentialsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudCredentials, err := getRemoteCloudCredentials(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCloudCredentialsAttributes(cloudCredentials, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudCredentialsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudCredentials: Beginning Creation")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudCredentialsAttr := models.CloudCredentialsAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudCredentialsAttr.Annotation = Annotation.(string)
	} else {
		cloudCredentialsAttr.Annotation = "{}"
	}

	if ClientId, ok := d.GetOk("client_id"); ok {
		cloudCredentialsAttr.ClientId = ClientId.(string)
	}

	if Email, ok := d.GetOk("email"); ok {
		cloudCredentialsAttr.Email = Email.(string)
	}

	if HttpProxy, ok := d.GetOk("http_proxy"); ok {
		cloudCredentialsAttr.HttpProxy = HttpProxy.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		cloudCredentialsAttr.Key = Key.(string)
	}

	if KeyId, ok := d.GetOk("key_id"); ok {
		cloudCredentialsAttr.KeyId = KeyId.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudCredentialsAttr.Name = Name.(string)
	}

	if RsaPrivateKey, ok := d.GetOk("rsa_private_key"); ok {
		cloudCredentialsAttr.RsaPrivateKey = RsaPrivateKey.(string)
	}
	cloudCredentials := models.NewCloudCredentials(fmt.Sprintf(models.RncloudCredentials, name), TenantDn, nameAlias, cloudCredentialsAttr)

	err := aciClient.Save(cloudCredentials)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTocloudRsAD, ok := d.GetOk("relation_cloud_rs_ad"); ok {
		relationParam := relationTocloudRsAD.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTocloudRsAD, ok := d.GetOk("relation_cloud_rs_ad"); ok {
		relationParam := relationTocloudRsAD.(string)
		err = aciClient.CreateRelationcloudRsAD(cloudCredentials.DistinguishedName, cloudCredentialsAttr.Annotation, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(cloudCredentials.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudCredentialsRead(ctx, d, m)
}

func resourceAciCloudCredentialsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudCredentials: Beginning Update")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudCredentialsAttr := models.CloudCredentialsAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudCredentialsAttr.Annotation = Annotation.(string)
	} else {
		cloudCredentialsAttr.Annotation = "{}"
	}

	if ClientId, ok := d.GetOk("client_id"); ok {
		cloudCredentialsAttr.ClientId = ClientId.(string)
	}

	if Email, ok := d.GetOk("email"); ok {
		cloudCredentialsAttr.Email = Email.(string)
	}

	if HttpProxy, ok := d.GetOk("http_proxy"); ok {
		cloudCredentialsAttr.HttpProxy = HttpProxy.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		cloudCredentialsAttr.Key = Key.(string)
	}

	if KeyId, ok := d.GetOk("key_id"); ok {
		cloudCredentialsAttr.KeyId = KeyId.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudCredentialsAttr.Name = Name.(string)
	}

	if RsaPrivateKey, ok := d.GetOk("rsa_private_key"); ok {
		cloudCredentialsAttr.RsaPrivateKey = RsaPrivateKey.(string)
	}
	cloudCredentials := models.NewCloudCredentials(fmt.Sprintf("credentials-%s", name), TenantDn, nameAlias, cloudCredentialsAttr)

	cloudCredentials.Status = "modified"

	err := aciClient.Save(cloudCredentials)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_cloud_rs_ad") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_cloud_rs_ad")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_cloud_rs_ad") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_cloud_rs_ad")
		err = aciClient.DeleteRelationcloudRsAD(cloudCredentials.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationcloudRsAD(cloudCredentials.DistinguishedName, cloudCredentialsAttr.Annotation, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(cloudCredentials.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudCredentialsRead(ctx, d, m)
}

func resourceAciCloudCredentialsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudCredentials, err := getRemoteCloudCredentials(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setCloudCredentialsAttributes(cloudCredentials, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	cloudRsADData, err := aciClient.ReadRelationcloudRsAD(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation cloudRsAD %v", err)
		d.Set("relation_cloud_rs_ad", "")
	} else {
		d.Set("relation_cloud_rs_ad", cloudRsADData.(string))
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudCredentialsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "cloudCredentials")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

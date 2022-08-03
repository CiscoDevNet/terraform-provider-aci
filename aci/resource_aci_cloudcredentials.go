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

func resourceAciAccessCredentialtomanagethecloudresources() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciAccessCredentialtomanagethecloudresourcesCreate,
		UpdateContext: resourceAciAccessCredentialtomanagethecloudresourcesUpdate,
		ReadContext:   resourceAciAccessCredentialtomanagethecloudresourcesRead,
		DeleteContext: resourceAciAccessCredentialtomanagethecloudresourcesDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAccessCredentialtomanagethecloudresourcesImport,
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
			"email": {
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
			"key_id": {
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
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to cloud:AD",
			}})),
	}
}

func getRemoteAccessCredentialtomanagethecloudresources(client *client.Client, dn string) (*models.AccessCredentialtomanagethecloudresources, error) {
	cloudCredentialsCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudCredentials := models.AccessCredentialtomanagethecloudresourcesFromContainer(cloudCredentialsCont)
	if cloudCredentials.DistinguishedName == "" {
		return nil, fmt.Errorf("AccessCredentialtomanagethecloudresources %s not found", cloudCredentials.DistinguishedName)
	}
	return cloudCredentials, nil
}

func setAccessCredentialtomanagethecloudresourcesAttributes(cloudCredentials *models.AccessCredentialtomanagethecloudresources, d *schema.ResourceData) (*schema.ResourceData, error) {
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

func resourceAciAccessCredentialtomanagethecloudresourcesImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudCredentials, err := getRemoteAccessCredentialtomanagethecloudresources(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setAccessCredentialtomanagethecloudresourcesAttributes(cloudCredentials, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAccessCredentialtomanagethecloudresourcesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AccessCredentialtomanagethecloudresources: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudCredentialsAttr := models.AccessCredentialtomanagethecloudresourcesAttributes{}

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
	cloudCredentials := models.NewAccessCredentialtomanagethecloudresources(fmt.Sprintf(models.RncloudCredentials, name), TenantDn, desc, nameAlias, cloudCredentialsAttr)

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
	return resourceAciAccessCredentialtomanagethecloudresourcesRead(ctx, d, m)
}

func resourceAciAccessCredentialtomanagethecloudresourcesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AccessCredentialtomanagethecloudresources: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudCredentialsAttr := models.AccessCredentialtomanagethecloudresourcesAttributes{}

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
	cloudCredentials := models.NewAccessCredentialtomanagethecloudresources(fmt.Sprintf("credentials-%s", name), TenantDn, desc, nameAlias, cloudCredentialsAttr)

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
	return resourceAciAccessCredentialtomanagethecloudresourcesRead(ctx, d, m)
}

func resourceAciAccessCredentialtomanagethecloudresourcesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudCredentials, err := getRemoteAccessCredentialtomanagethecloudresources(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setAccessCredentialtomanagethecloudresourcesAttributes(cloudCredentials, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	cloudRsADData, err := aciClient.ReadRelationcloudRsAD(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation cloudRsAD %v", err)
		d.Set("relation_cloud_rs_ad", "")
	} else {
		if _, ok := d.GetOk("relation_cloud_rs_ad"); ok {
			tfName := d.Get("relation_cloud_rs_ad").(string)
			if tfName != cloudRsADData {
				d.Set("relation_cloud_rs_ad", "")
			}
		}
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciAccessCredentialtomanagethecloudresourcesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciCloudAWSProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciCloudAWSProviderCreate,
		Update: resourceAciCloudAWSProviderUpdate,
		Read:   resourceAciCloudAWSProviderRead,
		Delete: resourceAciCloudAWSProviderDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudAWSProviderImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"access_key_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"http_proxy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_account_in_org": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_trusted": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"provider_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"secret_access_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteCloudAWSProvider(client *client.Client, dn string) (*models.CloudAWSProvider, error) {
	cloudAwsProviderCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudAwsProvider := models.CloudAWSProviderFromContainer(cloudAwsProviderCont)

	if cloudAwsProvider.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudAWSProvider %s not found", cloudAwsProvider.DistinguishedName)
	}

	return cloudAwsProvider, nil
}

func setCloudAWSProviderAttributes(cloudAwsProvider *models.CloudAWSProvider, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(cloudAwsProvider.DistinguishedName)
	d.Set("description", cloudAwsProvider.Description)
	d.Set("tenant_dn", GetParentDn(cloudAwsProvider.DistinguishedName))
	cloudAwsProviderMap, _ := cloudAwsProvider.ToMap()

	d.Set("access_key_id", cloudAwsProviderMap["accessKeyId"])
	d.Set("account_id", cloudAwsProviderMap["accountId"])
	d.Set("annotation", cloudAwsProviderMap["annotation"])
	d.Set("email", cloudAwsProviderMap["email"])
	d.Set("http_proxy", cloudAwsProviderMap["httpProxy"])
	d.Set("is_account_in_org", cloudAwsProviderMap["isAccountInOrg"])
	d.Set("is_trusted", cloudAwsProviderMap["isTrusted"])
	d.Set("name_alias", cloudAwsProviderMap["nameAlias"])
	d.Set("provider_id", cloudAwsProviderMap["providerId"])
	d.Set("region", cloudAwsProviderMap["region"])
	d.Set("secret_access_key", cloudAwsProviderMap["secretAccessKey"])
	return d
}

func resourceAciCloudAWSProviderImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudAwsProvider, err := getRemoteCloudAWSProvider(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setCloudAWSProviderAttributes(cloudAwsProvider, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudAWSProviderCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudAWSProvider: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudAwsProviderAttr := models.CloudAWSProviderAttributes{}
	if AccessKeyId, ok := d.GetOk("access_key_id"); ok {
		cloudAwsProviderAttr.AccessKeyId = AccessKeyId.(string)
	}
	if AccountId, ok := d.GetOk("account_id"); ok {
		cloudAwsProviderAttr.AccountId = AccountId.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudAwsProviderAttr.Annotation = Annotation.(string)
	}
	if Email, ok := d.GetOk("email"); ok {
		cloudAwsProviderAttr.Email = Email.(string)
	}
	if HttpProxy, ok := d.GetOk("http_proxy"); ok {
		cloudAwsProviderAttr.HttpProxy = HttpProxy.(string)
	}
	if IsAccountInOrg, ok := d.GetOk("is_account_in_org"); ok {
		cloudAwsProviderAttr.IsAccountInOrg = IsAccountInOrg.(string)
	}
	if IsTrusted, ok := d.GetOk("is_trusted"); ok {
		cloudAwsProviderAttr.IsTrusted = IsTrusted.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudAwsProviderAttr.NameAlias = NameAlias.(string)
	}
	if ProviderId, ok := d.GetOk("provider_id"); ok {
		cloudAwsProviderAttr.ProviderId = ProviderId.(string)
	}
	if Region, ok := d.GetOk("region"); ok {
		cloudAwsProviderAttr.Region = Region.(string)
	}
	if SecretAccessKey, ok := d.GetOk("secret_access_key"); ok {
		cloudAwsProviderAttr.SecretAccessKey = SecretAccessKey.(string)
	}
	cloudAwsProvider := models.NewCloudAWSProvider(fmt.Sprintf("awsprovider"), TenantDn, desc, cloudAwsProviderAttr)

	err := aciClient.Save(cloudAwsProvider)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	d.SetId(cloudAwsProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCloudAWSProviderRead(d, m)
}

func resourceAciCloudAWSProviderUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] CloudAWSProvider: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	TenantDn := d.Get("tenant_dn").(string)

	cloudAwsProviderAttr := models.CloudAWSProviderAttributes{}
	if AccessKeyId, ok := d.GetOk("access_key_id"); ok {
		cloudAwsProviderAttr.AccessKeyId = AccessKeyId.(string)
	}
	if AccountId, ok := d.GetOk("account_id"); ok {
		cloudAwsProviderAttr.AccountId = AccountId.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudAwsProviderAttr.Annotation = Annotation.(string)
	}
	if Email, ok := d.GetOk("email"); ok {
		cloudAwsProviderAttr.Email = Email.(string)
	}
	if HttpProxy, ok := d.GetOk("http_proxy"); ok {
		cloudAwsProviderAttr.HttpProxy = HttpProxy.(string)
	}
	if IsAccountInOrg, ok := d.GetOk("is_account_in_org"); ok {
		cloudAwsProviderAttr.IsAccountInOrg = IsAccountInOrg.(string)
	}
	if IsTrusted, ok := d.GetOk("is_trusted"); ok {
		cloudAwsProviderAttr.IsTrusted = IsTrusted.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudAwsProviderAttr.NameAlias = NameAlias.(string)
	}
	if ProviderId, ok := d.GetOk("provider_id"); ok {
		cloudAwsProviderAttr.ProviderId = ProviderId.(string)
	}
	if Region, ok := d.GetOk("region"); ok {
		cloudAwsProviderAttr.Region = Region.(string)
	}
	if SecretAccessKey, ok := d.GetOk("secret_access_key"); ok {
		cloudAwsProviderAttr.SecretAccessKey = SecretAccessKey.(string)
	}
	cloudAwsProvider := models.NewCloudAWSProvider(fmt.Sprintf("awsprovider"), TenantDn, desc, cloudAwsProviderAttr)

	cloudAwsProvider.Status = "modified"

	err := aciClient.Save(cloudAwsProvider)

	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	d.SetId(cloudAwsProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCloudAWSProviderRead(d, m)

}

func resourceAciCloudAWSProviderRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudAwsProvider, err := getRemoteCloudAWSProvider(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setCloudAWSProviderAttributes(cloudAwsProvider, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCloudAWSProviderDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudAwsProvider")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

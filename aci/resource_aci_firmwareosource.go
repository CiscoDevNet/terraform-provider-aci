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

func resourceAciFirmwareDownloadTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciFirmwareDownloadTaskCreate,
		UpdateContext: resourceAciFirmwareDownloadTaskUpdate,
		ReadContext:   resourceAciFirmwareDownloadTaskRead,
		DeleteContext: resourceAciFirmwareDownloadTaskDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFirmwareDownloadTaskImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"auth_pass": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"password",
					"key",
				}, false),
			},

			"auth_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"usePassword",
					"useSshKeyContents",
				}, false),
			},

			"dnld_task_flip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
			},

			"identity_private_key_contents": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"identity_private_key_passphrase": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"identity_public_key_contents": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"load_catalog_if_exists_and_newer": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"polling_interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"proto": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"scp",
					"http",
					"usbkey",
					"local",
				}, false),
			},

			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"user": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteFirmwareDownloadTask(client *client.Client, dn string) (*models.FirmwareDownloadTask, error) {
	firmwareOSourceCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	firmwareOSource := models.FirmwareDownloadTaskFromContainer(firmwareOSourceCont)

	if firmwareOSource.DistinguishedName == "" {
		return nil, fmt.Errorf("FirmwareDownloadTask %s not found", firmwareOSource.DistinguishedName)
	}

	return firmwareOSource, nil
}

func setFirmwareDownloadTaskAttributes(firmwareOSource *models.FirmwareDownloadTask, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(firmwareOSource.DistinguishedName)
	d.Set("description", firmwareOSource.Description)
	firmwareOSourceMap, err := firmwareOSource.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", firmwareOSourceMap["name"])
	d.Set("annotation", firmwareOSourceMap["annotation"])
	d.Set("auth_pass", firmwareOSourceMap["authPass"])
	d.Set("auth_type", firmwareOSourceMap["authType"])
	d.Set("dnld_task_flip", firmwareOSourceMap["dnldTaskFlip"])
	d.Set("load_catalog_if_exists_and_newer", firmwareOSourceMap["loadCatalogIfExistsAndNewer"])
	d.Set("name_alias", firmwareOSourceMap["nameAlias"])
	d.Set("polling_interval", firmwareOSourceMap["pollingInterval"])
	d.Set("proto", firmwareOSourceMap["proto"])
	d.Set("url", firmwareOSourceMap["url"])
	d.Set("user", firmwareOSourceMap["user"])
	return d, nil
}

func resourceAciFirmwareDownloadTaskImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	firmwareOSource, err := getRemoteFirmwareDownloadTask(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setFirmwareDownloadTaskAttributes(firmwareOSource, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFirmwareDownloadTaskCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FirmwareDownloadTask: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	firmwareOSourceAttr := models.FirmwareDownloadTaskAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		firmwareOSourceAttr.Annotation = Annotation.(string)
	} else {
		firmwareOSourceAttr.Annotation = "{}"
	}
	if AuthPass, ok := d.GetOk("auth_pass"); ok {
		firmwareOSourceAttr.AuthPass = AuthPass.(string)
	}
	if AuthType, ok := d.GetOk("auth_type"); ok {
		firmwareOSourceAttr.AuthType = AuthType.(string)
	}
	if DnldTaskFlip, ok := d.GetOk("dnld_task_flip"); ok {
		firmwareOSourceAttr.DnldTaskFlip = DnldTaskFlip.(string)
	}
	if IdentityPrivateKeyContents, ok := d.GetOk("identity_private_key_contents"); ok {
		firmwareOSourceAttr.IdentityPrivateKeyContents = IdentityPrivateKeyContents.(string)
	}
	if IdentityPrivateKeyPassphrase, ok := d.GetOk("identity_private_key_passphrase"); ok {
		firmwareOSourceAttr.IdentityPrivateKeyPassphrase = IdentityPrivateKeyPassphrase.(string)
	}
	if IdentityPublicKeyContents, ok := d.GetOk("identity_public_key_contents"); ok {
		firmwareOSourceAttr.IdentityPublicKeyContents = IdentityPublicKeyContents.(string)
	}
	if LoadCatalogIfExistsAndNewer, ok := d.GetOk("load_catalog_if_exists_and_newer"); ok {
		firmwareOSourceAttr.LoadCatalogIfExistsAndNewer = LoadCatalogIfExistsAndNewer.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		firmwareOSourceAttr.NameAlias = NameAlias.(string)
	}
	if Password, ok := d.GetOk("password"); ok {
		firmwareOSourceAttr.Password = Password.(string)
	}
	if PollingInterval, ok := d.GetOk("polling_interval"); ok {
		firmwareOSourceAttr.PollingInterval = PollingInterval.(string)
	}
	if Proto, ok := d.GetOk("proto"); ok {
		firmwareOSourceAttr.Proto = Proto.(string)
	}
	if Url, ok := d.GetOk("url"); ok {
		firmwareOSourceAttr.Url = Url.(string)
	}
	if User, ok := d.GetOk("user"); ok {
		firmwareOSourceAttr.User = User.(string)
	}
	firmwareOSource := models.NewFirmwareDownloadTask(fmt.Sprintf("fabric/fwrepop/osrc-%s", name), "uni", desc, firmwareOSourceAttr)

	err := aciClient.Save(firmwareOSource)
	if err != nil {
		return diag.FromErr(err)

	}

	d.SetId(firmwareOSource.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFirmwareDownloadTaskRead(ctx, d, m)
}

func resourceAciFirmwareDownloadTaskUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FirmwareDownloadTask: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	firmwareOSourceAttr := models.FirmwareDownloadTaskAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		firmwareOSourceAttr.Annotation = Annotation.(string)
	} else {
		firmwareOSourceAttr.Annotation = "{}"
	}
	if AuthPass, ok := d.GetOk("auth_pass"); ok {
		firmwareOSourceAttr.AuthPass = AuthPass.(string)
	}
	if AuthType, ok := d.GetOk("auth_type"); ok {
		firmwareOSourceAttr.AuthType = AuthType.(string)
	}
	if DnldTaskFlip, ok := d.GetOk("dnld_task_flip"); ok {
		firmwareOSourceAttr.DnldTaskFlip = DnldTaskFlip.(string)
	}
	if IdentityPrivateKeyContents, ok := d.GetOk("identity_private_key_contents"); ok {
		firmwareOSourceAttr.IdentityPrivateKeyContents = IdentityPrivateKeyContents.(string)
	}
	if IdentityPrivateKeyPassphrase, ok := d.GetOk("identity_private_key_passphrase"); ok {
		firmwareOSourceAttr.IdentityPrivateKeyPassphrase = IdentityPrivateKeyPassphrase.(string)
	}
	if IdentityPublicKeyContents, ok := d.GetOk("identity_public_key_contents"); ok {
		firmwareOSourceAttr.IdentityPublicKeyContents = IdentityPublicKeyContents.(string)
	}
	if LoadCatalogIfExistsAndNewer, ok := d.GetOk("load_catalog_if_exists_and_newer"); ok {
		firmwareOSourceAttr.LoadCatalogIfExistsAndNewer = LoadCatalogIfExistsAndNewer.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		firmwareOSourceAttr.NameAlias = NameAlias.(string)
	}
	if Password, ok := d.GetOk("password"); ok {
		firmwareOSourceAttr.Password = Password.(string)
	}
	if PollingInterval, ok := d.GetOk("polling_interval"); ok {
		firmwareOSourceAttr.PollingInterval = PollingInterval.(string)
	}
	if Proto, ok := d.GetOk("proto"); ok {
		firmwareOSourceAttr.Proto = Proto.(string)
	}
	if Url, ok := d.GetOk("url"); ok {
		firmwareOSourceAttr.Url = Url.(string)
	}
	if User, ok := d.GetOk("user"); ok {
		firmwareOSourceAttr.User = User.(string)
	}
	firmwareOSource := models.NewFirmwareDownloadTask(fmt.Sprintf("fabric/fwrepop/osrc-%s", name), "uni", desc, firmwareOSourceAttr)

	firmwareOSource.Status = "modified"

	err := aciClient.Save(firmwareOSource)

	if err != nil {
		return diag.FromErr(err)

	}

	d.SetId(firmwareOSource.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFirmwareDownloadTaskRead(ctx, d, m)

}

func resourceAciFirmwareDownloadTaskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	firmwareOSource, err := getRemoteFirmwareDownloadTask(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setFirmwareDownloadTaskAttributes(firmwareOSource, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFirmwareDownloadTaskDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "firmwareOSource")
	if err != nil {
		return diag.FromErr(err)

	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)

}

package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciCloudActiveDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudActiveDirectoryCreate,
		UpdateContext: resourceAciCloudActiveDirectoryUpdate,
		ReadContext:   resourceAciCloudActiveDirectoryRead,
		DeleteContext: resourceAciCloudActiveDirectoryDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudActiveDirectoryImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"active_directory_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteCloudActiveDirectory(client *client.Client, dn string) (*models.CloudActiveDirectory, error) {
	cloudADCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudAD := models.CloudActiveDirectoryFromContainer(cloudADCont)
	if cloudAD.DistinguishedName == "" {
		return nil, fmt.Errorf("Active Directory %s not found", dn)
	}
	return cloudAD, nil
}

func setCloudActiveDirectoryAttributes(cloudAD *models.CloudActiveDirectory, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudAD.DistinguishedName)
	if dn != cloudAD.DistinguishedName {
		d.Set("tenant_dn", "")
	} else {
		d.Set("tenant_dn", GetParentDn(dn, "/ad"))
	}

	cloudADMap, err := cloudAD.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("active_directory_id", cloudADMap["id"])
	d.Set("name", cloudADMap["name"])
	d.Set("name_alias", cloudADMap["nameAlias"])
	return d, nil
}

func resourceAciCloudActiveDirectoryImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudAD, err := getRemoteCloudActiveDirectory(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCloudActiveDirectoryAttributes(cloudAD, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudActiveDirectoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ActiveDirectory: Beginning Creation")
	aciClient := m.(*client.Client)
	active_directory_id := d.Get("active_directory_id").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudADAttr := models.CloudActiveDirectoryAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudADAttr.Annotation = Annotation.(string)
	} else {
		cloudADAttr.Annotation = "{}"
	}

	if ActiveDirectory_id, ok := d.GetOk("active_directory_id"); ok {
		cloudADAttr.ActiveDirectory_id = ActiveDirectory_id.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudADAttr.Name = Name.(string)
	}
	cloudAD := models.NewCloudActiveDirectory(fmt.Sprintf(models.RncloudAD, active_directory_id), TenantDn, nameAlias, cloudADAttr)

	err := aciClient.Save(cloudAD)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudAD.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudActiveDirectoryRead(ctx, d, m)
}

func resourceAciCloudActiveDirectoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ActiveDirectory: Beginning Update")
	aciClient := m.(*client.Client)
	active_directory_id := d.Get("active_directory_id").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudADAttr := models.CloudActiveDirectoryAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudADAttr.Annotation = Annotation.(string)
	} else {
		cloudADAttr.Annotation = "{}"
	}

	if ActiveDirectory_id, ok := d.GetOk("active_directory_id"); ok {
		cloudADAttr.ActiveDirectory_id = ActiveDirectory_id.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudADAttr.Name = Name.(string)
	}
	cloudAD := models.NewCloudActiveDirectory(fmt.Sprintf("ad-%s", active_directory_id), TenantDn, nameAlias, cloudADAttr)

	cloudAD.Status = "modified"

	err := aciClient.Save(cloudAD)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudAD.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudActiveDirectoryRead(ctx, d, m)
}

func resourceAciCloudActiveDirectoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudAD, err := getRemoteCloudActiveDirectory(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setCloudActiveDirectoryAttributes(cloudAD, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudActiveDirectoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "cloudAD")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

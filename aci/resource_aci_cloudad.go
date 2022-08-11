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

func resourceAciActiveDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciActiveDirectoryCreate,
		UpdateContext: resourceAciActiveDirectoryUpdate,
		ReadContext:   resourceAciActiveDirectoryRead,
		DeleteContext: resourceAciActiveDirectoryDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciActiveDirectoryImport,
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

func getRemoteActiveDirectory(client *client.Client, dn string) (*models.ActiveDirectory, error) {
	cloudADCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudAD := models.ActiveDirectoryFromContainer(cloudADCont)
	if cloudAD.DistinguishedName == "" {
		return nil, fmt.Errorf("ActiveDirectory %s not found", cloudAD.DistinguishedName)
	}
	return cloudAD, nil
}

func setActiveDirectoryAttributes(cloudAD *models.ActiveDirectory, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudAD.DistinguishedName)
	if dn != cloudAD.DistinguishedName {
		d.Set("tenant_dn", "")
	} else {
		d.Set("tenant_dn", GetParentDn(dn, "/act"))
	}

	cloudADMap, err := cloudAD.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("active_directory_id", cloudADMap["ActiveDirectory_id"])
	d.Set("name", cloudADMap["name"])
	d.Set("name_alias", cloudADMap["nameAlias"])
	return d, nil
}

func resourceAciActiveDirectoryImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudAD, err := getRemoteActiveDirectory(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setActiveDirectoryAttributes(cloudAD, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciActiveDirectoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ActiveDirectory: Beginning Creation")
	aciClient := m.(*client.Client)
	active_directory_id := d.Get("active_directory_id").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudADAttr := models.ActiveDirectoryAttributes{}

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
	cloudAD := models.NewActiveDirectory(fmt.Sprintf(models.RncloudAD, active_directory_id), TenantDn, nameAlias, cloudADAttr)

	err := aciClient.Save(cloudAD)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudAD.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciActiveDirectoryRead(ctx, d, m)
}

func resourceAciActiveDirectoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ActiveDirectory: Beginning Update")
	aciClient := m.(*client.Client)
	active_directory_id := d.Get("active_directory_id").(string)
	TenantDn := d.Get("tenant_dn").(string)

	cloudADAttr := models.ActiveDirectoryAttributes{}

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
	cloudAD := models.NewActiveDirectory(fmt.Sprintf("ad-%s", active_directory_id), TenantDn, nameAlias, cloudADAttr)

	cloudAD.Status = "modified"

	err := aciClient.Save(cloudAD)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudAD.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciActiveDirectoryRead(ctx, d, m)
}

func resourceAciActiveDirectoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudAD, err := getRemoteActiveDirectory(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setActiveDirectoryAttributes(cloudAD, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciActiveDirectoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

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

func resourceAciTenanttoCloudAccountAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciTenanttoCloudAccountAssociationCreate,
		UpdateContext: resourceAciTenanttoCloudAccountAssociationUpdate,
		ReadContext:   resourceAciTenanttoCloudAccountAssociationRead,
		DeleteContext: resourceAciTenanttoCloudAccountAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTenanttoCloudAccountAssociationImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cloud_account_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteTenanttoCloudAccountAssociation(client *client.Client, dn string) (*models.TenanttoCloudAccountAssociation, error) {
	fvRsCloudAccountCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvRsCloudAccount := models.TenanttoCloudAccountAssociationFromContainer(fvRsCloudAccountCont)
	if fvRsCloudAccount.DistinguishedName == "" {
		return nil, fmt.Errorf("TenanttoCloudAccountAssociation %s not found", fvRsCloudAccount.DistinguishedName)
	}
	return fvRsCloudAccount, nil
}

func setTenanttoCloudAccountAssociationAttributes(fvRsCloudAccount *models.TenanttoCloudAccountAssociation, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvRsCloudAccount.DistinguishedName)
	if dn != fvRsCloudAccount.DistinguishedName {
		d.Set("tenant_dn", "")
	} else {
		d.Set("tenant_dn", GetParentDn(dn, "/rsCloudAccount"))
	}
	fvRsCloudAccountMap, err := fvRsCloudAccount.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", fvRsCloudAccountMap["annotation"])
	d.Set("cloud_account_dn", fvRsCloudAccountMap["tDn"])
	d.Set("name_alias", fvRsCloudAccountMap["nameAlias"])
	return d, nil
}

func resourceAciTenanttoCloudAccountAssociationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvRsCloudAccount, err := getRemoteTenanttoCloudAccountAssociation(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setTenanttoCloudAccountAssociationAttributes(fvRsCloudAccount, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTenanttoCloudAccountAssociationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TenanttoCloudAccountAssociation: Beginning Creation")
	aciClient := m.(*client.Client)
	TenantDn := d.Get("tenant_dn").(string)

	fvRsCloudAccountAttr := models.TenanttoCloudAccountAssociationAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsCloudAccountAttr.Annotation = Annotation.(string)
	} else {
		fvRsCloudAccountAttr.Annotation = "{}"
	}

	if TDn, ok := d.GetOk("cloud_account_dn"); ok {
		fvRsCloudAccountAttr.TDn = TDn.(string)
	}
	fvRsCloudAccount := models.NewTenanttoCloudAccountAssociation(fmt.Sprintf(models.RnfvRsCloudAccount), TenantDn, nameAlias, fvRsCloudAccountAttr)

	err := aciClient.Save(fvRsCloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvRsCloudAccount.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciTenanttoCloudAccountAssociationRead(ctx, d, m)
}

func resourceAciTenanttoCloudAccountAssociationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TenanttoCloudAccountAssociation: Beginning Update")
	aciClient := m.(*client.Client)
	TenantDn := d.Get("tenant_dn").(string)

	fvRsCloudAccountAttr := models.TenanttoCloudAccountAssociationAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsCloudAccountAttr.Annotation = Annotation.(string)
	} else {
		fvRsCloudAccountAttr.Annotation = "{}"
	}

	if TDn, ok := d.GetOk("cloud_account_dn"); ok {
		fvRsCloudAccountAttr.TDn = TDn.(string)
	}
	fvRsCloudAccount := models.NewTenanttoCloudAccountAssociation(fmt.Sprintf("rsCloudAccount"), TenantDn, nameAlias, fvRsCloudAccountAttr)

	fvRsCloudAccount.Status = "modified"

	err := aciClient.Save(fvRsCloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvRsCloudAccount.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciTenanttoCloudAccountAssociationRead(ctx, d, m)
}

func resourceAciTenanttoCloudAccountAssociationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	fvRsCloudAccount, err := getRemoteTenanttoCloudAccountAssociation(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setTenanttoCloudAccountAssociationAttributes(fvRsCloudAccount, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciTenanttoCloudAccountAssociationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "fvRsCloudAccount")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

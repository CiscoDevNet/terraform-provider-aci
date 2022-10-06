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

func resourceAciTenantToCloudAccountAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciTenantToCloudAccountAssociationCreate,
		UpdateContext: resourceAciTenantToCloudAccountAssociationUpdate,
		ReadContext:   resourceAciTenantToCloudAccountAssociationRead,
		DeleteContext: resourceAciTenantToCloudAccountAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTenantToCloudAccountAssociationImport,
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

func getRemoteTenantToCloudAccountAssociation(client *client.Client, dn string) (*models.TenantToCloudAccountAssociation, error) {
	fvRsCloudAccountCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvRsCloudAccount := models.TenantToCloudAccountAssociationFromContainer(fvRsCloudAccountCont)
	if fvRsCloudAccount.DistinguishedName == "" {
		return nil, fmt.Errorf("Tenant To Cloud Account Association %s not found", fvRsCloudAccount.DistinguishedName)
	}
	return fvRsCloudAccount, nil
}

func setTenantToCloudAccountAssociationAttributes(fvRsCloudAccount *models.TenantToCloudAccountAssociation, d *schema.ResourceData) (*schema.ResourceData, error) {
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

func resourceAciTenantToCloudAccountAssociationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvRsCloudAccount, err := getRemoteTenantToCloudAccountAssociation(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setTenantToCloudAccountAssociationAttributes(fvRsCloudAccount, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTenantToCloudAccountAssociationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Tenant To Cloud Account Association: Beginning Creation")
	aciClient := m.(*client.Client)
	TenantDn := d.Get("tenant_dn").(string)

	fvRsCloudAccountAttr := models.TenantToCloudAccountAssociationAttributes{}

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
	fvRsCloudAccount := models.NewTenantToCloudAccountAssociation(fmt.Sprintf(models.RnfvRsCloudAccount), TenantDn, nameAlias, fvRsCloudAccountAttr)

	err := aciClient.Save(fvRsCloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvRsCloudAccount.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciTenantToCloudAccountAssociationRead(ctx, d, m)
}

func resourceAciTenantToCloudAccountAssociationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Tenant To Cloud Account Association: Beginning Update")
	aciClient := m.(*client.Client)
	TenantDn := d.Get("tenant_dn").(string)

	fvRsCloudAccountAttr := models.TenantToCloudAccountAssociationAttributes{}

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
	fvRsCloudAccount := models.NewTenantToCloudAccountAssociation(fmt.Sprintf("rsCloudAccount"), TenantDn, nameAlias, fvRsCloudAccountAttr)

	fvRsCloudAccount.Status = "modified"

	err := aciClient.Save(fvRsCloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvRsCloudAccount.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciTenantToCloudAccountAssociationRead(ctx, d, m)
}

func resourceAciTenantToCloudAccountAssociationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	fvRsCloudAccount, err := getRemoteTenantToCloudAccountAssociation(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setTenantToCloudAccountAssociationAttributes(fvRsCloudAccount, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciTenantToCloudAccountAssociationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

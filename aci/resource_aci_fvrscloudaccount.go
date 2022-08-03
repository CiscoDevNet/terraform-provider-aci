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

func resourceAciTenanttoaccountassociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciTenanttoaccountassociationCreate,
		UpdateContext: resourceAciTenanttoaccountassociationUpdate,
		ReadContext:   resourceAciTenanttoaccountassociationRead,
		DeleteContext: resourceAciTenanttoaccountassociationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTenanttoaccountassociationImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"t_dn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteTenanttoaccountassociation(client *client.Client, dn string) (*models.Tenanttoaccountassociation, error) {
	fvRsCloudAccountCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvRsCloudAccount := models.TenanttoaccountassociationFromContainer(fvRsCloudAccountCont)
	if fvRsCloudAccount.DistinguishedName == "" {
		return nil, fmt.Errorf("Tenanttoaccountassociation %s not found", fvRsCloudAccount.DistinguishedName)
	}
	return fvRsCloudAccount, nil
}

func setTenanttoaccountassociationAttributes(fvRsCloudAccount *models.Tenanttoaccountassociation, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvRsCloudAccount.DistinguishedName)
	d.Set("description", fvRsCloudAccount.Description)
	if dn != fvRsCloudAccount.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	fvRsCloudAccountMap, err := fvRsCloudAccount.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", fvRsCloudAccountMap["annotation"])
	d.Set("t_dn", fvRsCloudAccountMap["tDn"])
	d.Set("name_alias", fvRsCloudAccountMap["nameAlias"])
	return d, nil
}

func resourceAciTenanttoaccountassociationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvRsCloudAccount, err := getRemoteTenanttoaccountassociation(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setTenanttoaccountassociationAttributes(fvRsCloudAccount, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTenanttoaccountassociationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Tenanttoaccountassociation: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	TenantDn := d.Get("tenant_dn").(string)

	fvRsCloudAccountAttr := models.TenanttoaccountassociationAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsCloudAccountAttr.Annotation = Annotation.(string)
	} else {
		fvRsCloudAccountAttr.Annotation = "{}"
	}

	if TDn, ok := d.GetOk("t_dn"); ok {
		fvRsCloudAccountAttr.TDn = TDn.(string)
	}
	fvRsCloudAccount := models.NewTenanttoaccountassociation(fmt.Sprintf(models.RnfvRsCloudAccount), TenantDn, desc, nameAlias, fvRsCloudAccountAttr)

	err := aciClient.Save(fvRsCloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvRsCloudAccount.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciTenanttoaccountassociationRead(ctx, d, m)
}

func resourceAciTenanttoaccountassociationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Tenanttoaccountassociation: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	TenantDn := d.Get("tenant_dn").(string)

	fvRsCloudAccountAttr := models.TenanttoaccountassociationAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsCloudAccountAttr.Annotation = Annotation.(string)
	} else {
		fvRsCloudAccountAttr.Annotation = "{}"
	}

	if TDn, ok := d.GetOk("t_dn"); ok {
		fvRsCloudAccountAttr.TDn = TDn.(string)
	}
	fvRsCloudAccount := models.NewTenanttoaccountassociation(fmt.Sprintf("rsCloudAccount"), TenantDn, desc, nameAlias, fvRsCloudAccountAttr)

	fvRsCloudAccount.Status = "modified"

	err := aciClient.Save(fvRsCloudAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvRsCloudAccount.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciTenanttoaccountassociationRead(ctx, d, m)
}

func resourceAciTenanttoaccountassociationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	fvRsCloudAccount, err := getRemoteTenanttoaccountassociation(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setTenanttoaccountassociationAttributes(fvRsCloudAccount, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciTenanttoaccountassociationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

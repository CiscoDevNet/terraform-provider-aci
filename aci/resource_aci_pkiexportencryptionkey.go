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

func resourceAciAESEncryptionPassphraseandKeysforConfigExportImport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciAESEncryptionPassphraseandKeysforConfigExportImportCreate,
		UpdateContext: resourceAciAESEncryptionPassphraseandKeysforConfigExportImportUpdate,
		ReadContext:   resourceAciAESEncryptionPassphraseandKeysforConfigExportImportRead,
		DeleteContext: resourceAciAESEncryptionPassphraseandKeysforConfigExportImportDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAESEncryptionPassphraseandKeysforConfigExportImportImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"clear_encryption_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "no",
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"passphrase": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"passphrase_key_derivation_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"v1",
				}, false),
			},
			"strong_encryption_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
		})),
	}
}

func GetRemoteAESEncryptionPassphraseandKeysforConfigExportImport(client *client.Client, dn string) (*models.AESEncryptionPassphraseandKeysforConfigExportImport, error) {
	pkiExportEncryptionKeyCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	pkiExportEncryptionKey := models.AESEncryptionPassphraseandKeysforConfigExportImportFromContainer(pkiExportEncryptionKeyCont)
	if pkiExportEncryptionKey.DistinguishedName == "" {
		return nil, fmt.Errorf("AESEncryptionPassphraseandKeysforConfigExport(andImport) %s not found", pkiExportEncryptionKey.DistinguishedName)
	}
	return pkiExportEncryptionKey, nil
}

func setAESEncryptionPassphraseandKeysforConfigExportImportAttributes(pkiExportEncryptionKey *models.AESEncryptionPassphraseandKeysforConfigExportImport, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(pkiExportEncryptionKey.DistinguishedName)
	d.Set("description", pkiExportEncryptionKey.Description)
	pkiExportEncryptionKeyMap, err := pkiExportEncryptionKey.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", pkiExportEncryptionKeyMap["annotation"])
	d.Set("passphrase_key_derivation_version", pkiExportEncryptionKeyMap["passphraseKeyDerivationVersion"])
	d.Set("strong_encryption_enabled", pkiExportEncryptionKeyMap["strongEncryptionEnabled"])
	d.Set("name_alias", pkiExportEncryptionKeyMap["nameAlias"])
	return d, nil
}

func resourceAciAESEncryptionPassphraseandKeysforConfigExportImportImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	pkiExportEncryptionKey, err := GetRemoteAESEncryptionPassphraseandKeysforConfigExportImport(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setAESEncryptionPassphraseandKeysforConfigExportImportAttributes(pkiExportEncryptionKey, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAESEncryptionPassphraseandKeysforConfigExportImportCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AESEncryptionPassphraseandKeysforConfigExport(andImport): Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	pkiExportEncryptionKeyAttr := models.AESEncryptionPassphraseandKeysforConfigExportImportAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		pkiExportEncryptionKeyAttr.Annotation = Annotation.(string)
	} else {
		pkiExportEncryptionKeyAttr.Annotation = "{}"
	}

	if ClearEncryptionKey, ok := d.GetOk("clear_encryption_key"); ok {
		pkiExportEncryptionKeyAttr.ClearEncryptionKey = ClearEncryptionKey.(string)
	}

	if Passphrase, ok := d.GetOk("passphrase"); ok {
		pkiExportEncryptionKeyAttr.Passphrase = Passphrase.(string)
	}

	if PassphraseKeyDerivationVersion, ok := d.GetOk("passphrase_key_derivation_version"); ok {
		pkiExportEncryptionKeyAttr.PassphraseKeyDerivationVersion = PassphraseKeyDerivationVersion.(string)
	}

	if StrongEncryptionEnabled, ok := d.GetOk("strong_encryption_enabled"); ok {
		pkiExportEncryptionKeyAttr.StrongEncryptionEnabled = StrongEncryptionEnabled.(string)
	}

	pkiExportEncryptionKey := models.NewAESEncryptionPassphraseandKeysforConfigExportImport(fmt.Sprintf("exportcryptkey"), "uni", desc, nameAlias, pkiExportEncryptionKeyAttr)
	pkiExportEncryptionKey.Status = "modified"
	err := aciClient.Save(pkiExportEncryptionKey)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(pkiExportEncryptionKey.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciAESEncryptionPassphraseandKeysforConfigExportImportRead(ctx, d, m)
}

func resourceAciAESEncryptionPassphraseandKeysforConfigExportImportUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AESEncryptionPassphraseandKeysforConfigExport(andImport): Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	pkiExportEncryptionKeyAttr := models.AESEncryptionPassphraseandKeysforConfigExportImportAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		pkiExportEncryptionKeyAttr.Annotation = Annotation.(string)
	} else {
		pkiExportEncryptionKeyAttr.Annotation = "{}"
	}

	if ClearEncryptionKey, ok := d.GetOk("clear_encryption_key"); ok {
		pkiExportEncryptionKeyAttr.ClearEncryptionKey = ClearEncryptionKey.(string)
	}

	if Passphrase, ok := d.GetOk("passphrase"); ok {
		pkiExportEncryptionKeyAttr.Passphrase = Passphrase.(string)
	}

	if PassphraseKeyDerivationVersion, ok := d.GetOk("passphrase_key_derivation_version"); ok {
		pkiExportEncryptionKeyAttr.PassphraseKeyDerivationVersion = PassphraseKeyDerivationVersion.(string)
	}

	if StrongEncryptionEnabled, ok := d.GetOk("strong_encryption_enabled"); ok {
		pkiExportEncryptionKeyAttr.StrongEncryptionEnabled = StrongEncryptionEnabled.(string)
	}

	pkiExportEncryptionKey := models.NewAESEncryptionPassphraseandKeysforConfigExportImport(fmt.Sprintf("exportcryptkey"), "uni", desc, nameAlias, pkiExportEncryptionKeyAttr)
	pkiExportEncryptionKey.Status = "modified"
	err := aciClient.Save(pkiExportEncryptionKey)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(pkiExportEncryptionKey.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciAESEncryptionPassphraseandKeysforConfigExportImportRead(ctx, d, m)
}

func resourceAciAESEncryptionPassphraseandKeysforConfigExportImportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	pkiExportEncryptionKey, err := GetRemoteAESEncryptionPassphraseandKeysforConfigExportImport(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setAESEncryptionPassphraseandKeysforConfigExportImportAttributes(pkiExportEncryptionKey, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciAESEncryptionPassphraseandKeysforConfigExportImportDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name pkiExportEncryptionKey cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}

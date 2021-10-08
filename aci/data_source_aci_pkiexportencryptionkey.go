package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAESEncryptionPassphraseandKeysforConfigExportImport() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciAESEncryptionPassphraseandKeysforConfigExportImportRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"passphrase_key_derivation_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"strong_encryption_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciAESEncryptionPassphraseandKeysforConfigExportImportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("exportcryptkey")
	dn := fmt.Sprintf("uni/%s", rn)
	pkiExportEncryptionKey, err := getRemoteAESEncryptionPassphraseandKeysforConfigExportImport(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setAESEncryptionPassphraseandKeysforConfigExportImportAttributes(pkiExportEncryptionKey, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

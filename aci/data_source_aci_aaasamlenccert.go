package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciKeypairforSAMLEncryption() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciKeypairforSAMLEncryptionRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"regenerate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"certificate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expiry_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"certificate_validty": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteKeypairforSAMLEncryption(client *client.Client, dn string) (*models.KeypairforSAMLEncryption, error) {
	aaaSamlEncCertCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaSamlEncCert := models.KeypairforSAMLEncryptionFromContainer(aaaSamlEncCertCont)
	if aaaSamlEncCert.DistinguishedName == "" {
		return nil, fmt.Errorf("KeypairforSAMLEncryption %s not found", aaaSamlEncCert.DistinguishedName)
	}
	return aaaSamlEncCert, nil
}

func setKeypairforSAMLEncryptionAttributes(aaaSamlEncCert *models.KeypairforSAMLEncryption, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaSamlEncCert.DistinguishedName)
	d.Set("description", aaaSamlEncCert.Description)
	aaaSamlEncCertMap, err := aaaSamlEncCert.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaSamlEncCertMap["annotation"])
	d.Set("regenerate", aaaSamlEncCertMap["regenerate"])
	d.Set("name_alias", aaaSamlEncCertMap["nameAlias"])
	d.Set("certificate", aaaSamlEncCertMap["cert"])
	d.Set("expiry_status", aaaSamlEncCertMap["expState"])
	d.Set("certificate_validty", aaaSamlEncCertMap["certValidUntil"])
	return d, nil
}

func dataSourceAciKeypairforSAMLEncryptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := "default"

	rn := fmt.Sprintf("userext/samlext/samlenccert-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	aaaSamlEncCert, err := getRemoteKeypairforSAMLEncryption(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setKeypairforSAMLEncryptionAttributes(aaaSamlEncCert, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

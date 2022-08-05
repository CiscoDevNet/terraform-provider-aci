package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciContractProvider() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciContractProviderRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"application_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"contract_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"contract_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
			},
		},
	}
}

func dataSourceAciContractProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	contractType := d.Get("contract_type").(string)
	tnVzBrCPName := GetMOName(d.Get("contract_dn").(string))
	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	if contractType == "provider" {

		rn := fmt.Sprintf("rsprov-%s", tnVzBrCPName)
		dn := fmt.Sprintf("%s/%s", ApplicationEPGDn, rn)

		fvRsProv, err := getRemoteContractProvider(aciClient, dn)
		if err != nil {
			return nil
		}
		d.SetId(dn)
		_, err = setContractProviderDataAttributes(fvRsProv, d)
		if err != nil {
			return nil
		}

	} else if contractType == "consumer" {
		rn := fmt.Sprintf("rscons-%s", tnVzBrCPName)
		dn := fmt.Sprintf("%s/%s", ApplicationEPGDn, rn)

		fvRsCons, err := getRemoteContractConsumer(aciClient, dn)
		if err != nil {
			return nil
		}
		d.SetId(dn)
		_, err = setContractConsumerDataAttributes(fvRsCons, d)
		if err != nil {
			return nil
		}

	} else {
		return diag.Errorf("Contract Type: Value must be from [provider, consumer]")
	}

	return nil
}

func setContractConsumerDataAttributes(fvRsCons *models.ContractConsumer, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fvRsCons.DistinguishedName)

	// d.Set("application_epg_dn", GetParentDn(fvRsCons.DistinguishedName))
	fvRsConsMap, err := fvRsCons.ToMap()
	if err != nil {
		return d, err
	}

	// d.Set("contract_name", fvRsConsMap["tnVzBrCPName"])

	d.Set("annotation", fvRsConsMap["annotation"])
	d.Set("prio", fvRsConsMap["prio"])
	return d, nil
}

func setContractProviderDataAttributes(fvRsProv *models.ContractProvider, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fvRsProv.DistinguishedName)

	// d.Set("application_epg_dn", GetParentDn(fvRsProv.DistinguishedName))
	fvRsProvMap, err := fvRsProv.ToMap()
	if err != nil {
		return d, err
	}
	// d.Set("contract_name", fvRsProvMap["tnVzBrCPName"])

	d.Set("annotation", fvRsProvMap["annotation"])
	d.Set("match_t", fvRsProvMap["matchT"])
	d.Set("prio", fvRsProvMap["prio"])
	return d, nil
}

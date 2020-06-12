package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciContractProvider() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciContractProviderRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"contract_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"contract_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
		}),
	}
}

func dataSourceAciContractProviderRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	contractType := d.Get("contract_type").(string)
	tnVzBrCPName := d.Get("contract_name").(string)
	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	if contractType == "provider" {

		rn := fmt.Sprintf("rsprov-%s", tnVzBrCPName)
		dn := fmt.Sprintf("%s/%s", ApplicationEPGDn, rn)

		fvRsProv, err := getRemoteContractProvider(aciClient, dn)
		if err != nil {
			return nil
		}

		setContractProviderDataAttributes(fvRsProv, d)

	} else if contractType == "consumer" {
		rn := fmt.Sprintf("rscons-%s", tnVzBrCPName)
		dn := fmt.Sprintf("%s/%s", ApplicationEPGDn, rn)

		fvRsCons, err := getRemoteContractConsumer(aciClient, dn)
		if err != nil {
			return nil
		}

		setContractConsumerDataAttributes(fvRsCons, d)

	} else {
		return fmt.Errorf("Contract Type: Value must be from [provider, consumer]")
	}

	return nil
}

func setContractConsumerDataAttributes(fvRsCons *models.ContractConsumer, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvRsCons.DistinguishedName)

	// d.Set("application_epg_dn", GetParentDn(fvRsCons.DistinguishedName))
	fvRsConsMap, _ := fvRsCons.ToMap()

	d.Set("contract_name", fvRsConsMap["tnVzBrCPName"])

	d.Set("annotation", fvRsConsMap["annotation"])
	d.Set("prio", fvRsConsMap["prio"])
	return d
}

func setContractProviderDataAttributes(fvRsProv *models.ContractProvider, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvRsProv.DistinguishedName)

	// d.Set("application_epg_dn", GetParentDn(fvRsProv.DistinguishedName))
	fvRsProvMap, _ := fvRsProv.ToMap()
	d.Set("contract_name", fvRsProvMap["tnVzBrCPName"])

	d.Set("annotation", fvRsProvMap["annotation"])
	d.Set("match_t", fvRsProvMap["matchT"])
	d.Set("prio", fvRsProvMap["prio"])
	return d
}

package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciContractProvider() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciContractProviderRead,

		SchemaVersion: 1,

		Schema: AppendAttrSchemas(map[string]*schema.Schema{
			"application_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"contract_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"contract_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"consumer",
					"provider",
				}, false),
			},
			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciContractProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	contractType := d.Get("contract_type").(string)
	ApplicationEPGDn := d.Get("application_epg_dn").(string)
	tnVzBrCPName := GetMOName(d.Get("contract_dn").(string))

	if contractType == "provider" {

		rn := fmt.Sprintf("rsprov-%s", tnVzBrCPName)
		dn := fmt.Sprintf("%s/%s", ApplicationEPGDn, rn)

		fvRsProv, err := getRemoteContractProvider(aciClient, dn)
		if err != nil {
			return errorForObjectNotFound(err, dn, d)
		}
		fvRsProvMap, _ := fvRsProv.ToMap()
		name := fvRsProvMap["tnVzBrCPName"]
		pDN := GetParentDn(dn, fmt.Sprintf("/rsprov-%s", name))
		d.Set("application_epg_dn", pDN)
		_, err = setContractProviderAttributes(fvRsProv, d)
		if err != nil {
			d.SetId("")
			return nil
		}
		d.SetId(dn)

	} else if contractType == "consumer" {
		rn := fmt.Sprintf("rscons-%s", tnVzBrCPName)
		dn := fmt.Sprintf("%s/%s", ApplicationEPGDn, rn)

		fvRsCons, err := getRemoteContractConsumer(aciClient, dn)
		if err != nil {
			return errorForObjectNotFound(err, dn, d)
		}
		_, err = setContractConsumerAttributes(fvRsCons, d)
		if err != nil {
			d.SetId("")
			return nil
		}
		d.SetId(dn)
	}

	return nil
}

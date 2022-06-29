package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAbstractionofIPAddressBlock() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciAbstractionofIPAddressBlockRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"multicast_address_pool_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"from": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"to": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciAbstractionofIPAddressBlockRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	_from := d.Get("_from").(string)
	to := d.Get("to").(string)
	MulticastAddressPoolDn := d.Get("multicast_address_pool_dn").(string)
	rn := fmt.Sprintf(models.RnfvnsMcastAddrBlk, _from, to)
	dn := fmt.Sprintf("%s/%s", MulticastAddressPoolDn, rn)

	fvnsMcastAddrBlk, err := getRemoteAbstractionofIPAddressBlock(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setAbstractionofIPAddressBlockAttributes(fvnsMcastAddrBlk, d)
	if err != nil {
		return nil
	}

	return nil
}

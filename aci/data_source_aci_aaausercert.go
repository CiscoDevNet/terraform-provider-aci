package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciX509Certificate() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciX509CertificateRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"local_user_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciX509CertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("usercert-%s", name)
	LocalUserDn := d.Get("local_user_dn").(string)

	dn := fmt.Sprintf("%s/%s", LocalUserDn, rn)

	aaaUserCert, err := getRemoteX509Certificate(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)
	_, err = setX509CertificateAttributes(aaaUserCert, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

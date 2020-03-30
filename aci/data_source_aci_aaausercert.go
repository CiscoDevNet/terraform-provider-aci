package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciX509Certificate() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciX509CertificateRead,

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

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func dataSourceAciX509CertificateRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("usercert-%s", name)
	LocalUserDn := d.Get("local_user_dn").(string)

	dn := fmt.Sprintf("%s/%s", LocalUserDn, rn)

	aaaUserCert, err := getRemoteX509Certificate(aciClient, dn)

	if err != nil {
		return err
	}
	setX509CertificateAttributes(aaaUserCert, d)
	return nil
}

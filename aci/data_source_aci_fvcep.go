package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciClientEndPoint() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciClientEndPointRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"client_end_point_id": &schema.Schema{
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

func getRemoteClientEndPoint(client *client.Client, dn string) (*models.ClientEndPoint, error) {
	fvCEpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvCEp := models.ClientEndPointFromContainer(fvCEpCont)

	if fvCEp.DistinguishedName == "" {
		return nil, fmt.Errorf("ClientEndPoint %s not found", fvCEp.DistinguishedName)
	}

	return fvCEp, nil
}

func setClientEndPointAttributes(fvCEp *models.ClientEndPoint, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(fvCEp.DistinguishedName)
	d.Set("description", fvCEp.Description)
	if dn != fvCEp.DistinguishedName {
		d.Set("application_epg_dn", "")
	}
	fvCEpMap, _ := fvCEp.ToMap()

	d.Set("name", fvCEpMap["name"])
	d.Set("annotation", fvCEpMap["annotation"])
	d.Set("client_end_point_id", fvCEpMap["id"])
	d.Set("name_alias", fvCEpMap["nameAlias"])

	return d
}

func dataSourceAciClientEndPointRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("cep-%s", name)
	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	dn := fmt.Sprintf("%s/%s", ApplicationEPGDn, rn)

	fvCEp, err := getRemoteClientEndPoint(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setClientEndPointAttributes(fvCEp, d)
	return nil
}

package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAutonomousSystemProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciAutonomousSystemProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"asn": &schema.Schema{
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

func dataSourceAciAutonomousSystemProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("clouddomp/as")

	dn := fmt.Sprintf("uni/%s", rn)

	cloudBgpAsP, err := getRemoteAutonomousSystemProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setAutonomousSystemProfileAttributes(cloudBgpAsP, d)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getRemoteAutonomousSystemProfile(client *client.Client, dn string) (*models.AutonomousSystemProfile, error) {
	cloudBgpAsPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudBgpAsP := models.AutonomousSystemProfileFromContainer(cloudBgpAsPCont)

	if cloudBgpAsP.DistinguishedName == "" {
		return nil, fmt.Errorf("AutonomousSystemProfile %s not found", cloudBgpAsP.DistinguishedName)
	}

	return cloudBgpAsP, nil
}

func setAutonomousSystemProfileAttributes(cloudBgpAsP *models.AutonomousSystemProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(cloudBgpAsP.DistinguishedName)
	d.Set("description", cloudBgpAsP.Description)
	cloudBgpAsPMap, err := cloudBgpAsP.ToMap()

	if err != nil {
		return d, err
	}

	d.Set("annotation", cloudBgpAsPMap["annotation"])
	d.Set("asn", cloudBgpAsPMap["asn"])
	d.Set("name_alias", cloudBgpAsPMap["nameAlias"])
	return d, nil
}

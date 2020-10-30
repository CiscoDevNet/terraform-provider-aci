package aci

import (
	"fmt"
	"strings"

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
				Optional: true,
				Computed: true,
			},

			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vlan": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"object_dns": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}

func getRemoteClientEndPoint(client *client.Client, dn, query string) (oDn []string, err error) {
	baseURL := "/api/node/class"

	var duURL string
	if query == "" {
		duURL = fmt.Sprintf("%s/%s/fvCEp.json", baseURL, dn)
	} else {
		duURL = fmt.Sprintf("%s/%s/fvCEp.json?query-target-filter=and(%s)", baseURL, dn, query)
	}

	fvCEpCont, err := client.GetViaURL(duURL)
	if err != nil {
		return nil, err
	}

	objectIds := make([]string, 0, 1)

	count, err := fvCEpCont.ArrayCount("imdata")
	if err != nil {
		return nil, err
	}
	for i := 0; i < count; i++ {
		clientEndPointCont, err := fvCEpCont.ArrayElement(i, "imdata")
		if err != nil {
			return nil, err
		}
		objectIds = append(objectIds, models.G(clientEndPointCont.S("fvCEp", "attributes"), "dn"))
	}

	return objectIds, nil
}

func dataSourceAciClientEndPointRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	var queryString string
	if mac, ok := d.GetOk("mac"); ok {
		if queryString != "" {
			queryString = fmt.Sprintf("%s,eq(fvCEp.mac, \"%s\")", queryString, mac.(string))
		} else {
			queryString = fmt.Sprintf("eq(fvCEp.mac, \"%s\")", mac.(string))
		}
	}

	if ip, ok := d.GetOk("ip"); ok {
		if queryString != "" {
			queryString = fmt.Sprintf("%s,eq(fvCEp.ip, \"%s\")", queryString, ip.(string))
		} else {
			queryString = fmt.Sprintf("eq(fvCEp.ip, \"%s\")", ip.(string))
		}
	}

	if name, ok := d.GetOk("name"); ok {
		if queryString != "" {
			queryString = fmt.Sprintf("%s,eq(fvCEp.name, \"%s\")", queryString, name.(string))
		} else {
			queryString = fmt.Sprintf("eq(fvCEp.name, \"%s\")", name.(string))
		}
	}

	if vlan, ok := d.GetOk("vlan"); ok {
		if queryString != "" {
			queryString = fmt.Sprintf("%s,eq(fvCEp.encap, \"vlan-%s\")", queryString, vlan.(string))
		} else {
			queryString = fmt.Sprintf("eq(fvCEp.encap, \"vlan-%s\")", vlan.(string))
		}
	}

	objs, err := getRemoteClientEndPoint(aciClient, ApplicationEPGDn, queryString)
	if err != nil {
		return err
	}

	d.Set("object_dns", objs)
	d.SetId(strings.Join(objs, " "))
	return nil
}

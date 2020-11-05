package aci

import (
	"fmt"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciClientEndPoint() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciClientEndPointRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
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

			"fvcep_objects": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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

						"tenant_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"vrf_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"application_profile_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"epg_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"l2out_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"instance_profile_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		}),
	}
}

func extractInfo(con *container.Container) (obj map[string]interface{}, dn string) {
	infoMap := make(map[string]interface{})
	dnString := models.G(con, "dn")
	infoMap["name"] = models.G(con, "name")
	infoMap["mac"] = models.G(con, "mac")
	infoMap["ip"] = models.G(con, "ip")
	infoMap["vlan"] = models.G(con, "encap")

	dnInfoList := strings.Split(dnString, "/")
	tenantInfo := strings.Split(dnInfoList[1], "-")
	if tenantInfo[0] == "tn" {
		infoMap["tenant_name"] = tenantInfo[1]

		level2Info := strings.Split(dnInfoList[2], "-")
		if level2Info[0] == "ctx" {
			infoMap["vrf_name"] = level2Info[1]

		} else if level2Info[0] == "ap" {
			infoMap["application_profile_name"] = level2Info[1]

			level3Info := strings.Split(dnInfoList[3], "-")
			if level3Info[0] == "epg" {
				infoMap["epg_name"] = level3Info[1]
			} else {
				return nil, ""
			}

		} else if level2Info[0] == "l2out" {
			infoMap["l2out_name"] = level2Info[1]

			level3Info := strings.Split(dnInfoList[3], "-")
			if level3Info[0] == "instP" {
				infoMap["instance_profile_name"] = level3Info[1]
			} else {
				return nil, ""
			}

		} else {
			return nil, ""
		}

	} else {
		return nil, ""
	}

	return infoMap, dnString
}

func getRemoteClientEndPoint(client *client.Client, query string) (objMap []interface{}, objdns []string, err error) {
	baseURL := "/api/node/class"

	var duURL string
	if query == "" {
		duURL = fmt.Sprintf("%s/fvCEp.json", baseURL)
	} else {
		duURL = fmt.Sprintf("%s/fvCEp.json?query-target-filter=and(%s)", baseURL, query)
	}

	fvCEpCont, err := client.GetViaURL(duURL)
	if err != nil {
		return nil, nil, err
	}

	objects := make([]interface{}, 0, 1)
	dns := make([]string, 0, 1)

	count, err := fvCEpCont.ArrayCount("imdata")
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < count; i++ {
		clientEndPointCont, err := fvCEpCont.ArrayElement(i, "imdata")
		if err != nil {
			return nil, nil, err
		}

		objMap, dn := extractInfo(clientEndPointCont.S("fvCEp", "attributes"))
		if dn != "" {
			objects = append(objects, objMap)
			dns = append(dns, dn)
		}
	}

	return objects, dns, nil
}

func dataSourceAciClientEndPointRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

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

	objects, dns, err := getRemoteClientEndPoint(aciClient, queryString)
	if err != nil {
		return err
	}

	d.Set("fvcep_objects", objects)
	d.SetId(strings.Join(dns, " "))
	return nil
}

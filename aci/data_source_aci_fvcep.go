package aci

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciClientEndPoint() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciClientEndPointRead,

		SchemaVersion: 1,

		Schema: AppendAttrSchemas(map[string]*schema.Schema{
			"name": &schema.Schema{Type: schema.TypeString,
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
						"ips": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

						"esg_name": &schema.Schema{
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

						"endpoint_path": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"base_epg": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		}, GetBaseAttrSchema(), GetAllowEmptyAttrSchema()),
	}
}

func extractInfo(con *container.Container) (obj map[string]interface{}, dn string) {
	infoMap := make(map[string]interface{})
	dnString := models.G(con, "dn")
	infoMap["name"] = models.G(con, "name")
	infoMap["mac"] = models.G(con, "mac")
	infoMap["ip"] = models.G(con, "ip")
	if infoMap["ip"] == "{}" {
		infoMap["ip"] = ""
	}

	infoMap["vlan"] = models.G(con, "encap")

	dnInfoList := strings.Split(dnString, "/")
	tenantInfo := regexp.MustCompile("-").Split(dnInfoList[1], 2)
	if tenantInfo[0] == "tn" {
		infoMap["tenant_name"] = tenantInfo[1]

		level2Info := regexp.MustCompile("-").Split(dnInfoList[2], 2)
		if level2Info[0] == "ctx" {
			infoMap["vrf_name"] = level2Info[1]

		} else if level2Info[0] == "ap" {
			infoMap["application_profile_name"] = level2Info[1]

			level3Info := regexp.MustCompile("-").Split(dnInfoList[3], 2)
			if level3Info[0] == "epg" {
				infoMap["epg_name"] = level3Info[1]
			} else if level3Info[0] == "esg" {
				infoMap["esg_name"] = level3Info[1]
				baseEpgDnInfoList := strings.Split(models.G(con, "baseEpgDn"), "/")
				baseEpgTenantInfo := regexp.MustCompile("-").Split(baseEpgDnInfoList[1], 2)
				if baseEpgTenantInfo[0] == "tn" {
					baseEpgTenantName := baseEpgTenantInfo[1]
					baseEpgLevel2Info := regexp.MustCompile("-").Split(baseEpgDnInfoList[2], 2)
					if baseEpgLevel2Info[0] == "ap" {
						baseEpgApplicationProfileName := baseEpgLevel2Info[1]
						baseEpgLevel3Info := regexp.MustCompile("-").Split(baseEpgDnInfoList[3], 2)
						if baseEpgLevel3Info[0] == "epg" {
							baseEpgEpgName := baseEpgLevel3Info[1]
							infoMap["base_epg"] = map[string]interface{}{
								"tenant_name":              baseEpgTenantName,
								"application_profile_name": baseEpgApplicationProfileName,
								"epg_name":                 baseEpgEpgName,
							}
						}
					}
				}
			} else {
				return nil, ""
			}

		} else if level2Info[0] == "l2out" {
			infoMap["l2out_name"] = level2Info[1]

			level3Info := regexp.MustCompile("-").Split(dnInfoList[3], 2)
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

func extractEndpointPaths(cont *container.Container) ([]string, error) {
	paths := make([]string, 0, 1)

	count, err := cont.ArrayCount("imdata")
	if err != nil {
		return paths, err
	}

	for i := 0; i < count; i++ {
		pathEpCont, err := cont.ArrayElement(i, "imdata")
		if err != nil {
			return paths, err
		}

		tDN := models.StripQuotes(pathEpCont.S("fvRsCEpToPathEp", "attributes", "tDn").String())

		paths = append(paths, tDN)
	}
	return paths, nil
}

func getRemoteClientEndPoint(client *client.Client, query string, allowEmptyResult bool) (objMap []interface{}, objdns []string, err error) {
	baseURL := "/api/node/class"

	var duURL string
	if query == "" {
		duURL = fmt.Sprintf("%s/fvCEp.json", baseURL)
	} else {
		duURL = fmt.Sprintf("%s/fvCEp.json?query-target-filter=and(%s)&rsp-subtree=children&rsp-subtree-class=fvIp", baseURL, query)
	}

	fvCEpCont, err := client.GetViaURL(duURL)
	if err != nil {
		return nil, nil, allowEmpty(err, allowEmptyResult)
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

		// Reading fvIp object properties
		fvIpObjects, err := clientEndPointCont.S("fvCEp", "children").Children()
		if err == nil {
			ips := make([]string, 0, 1)

			if err == nil {
				for j := 0; j < len(fvIpObjects); j++ {
					ips = append(ips, models.G(fvIpObjects[j].S("fvIp", "attributes"), "addr"))
				}
			}

			objMap["ips"] = ips
			if len(ips) > 0 {
				objMap["ip"] = ips[0]
			}
		} else {
			if errors.Is(err, container.ErrNotObjOrArray) == false {
				return nil, nil, err
			}
		}

		if dn != "" {
			durl := fmt.Sprintf("%s/%s/fvRsCEpToPathEp.json", baseURL, dn)
			cepToPathEpCont, err := client.GetViaURL(durl)
			if err == nil {
				endpointPaths, err := extractEndpointPaths(cepToPathEpCont)
				if err == nil {
					objMap["endpoint_path"] = endpointPaths
				}
			}
			objects = append(objects, objMap)
			dns = append(dns, dn)
		}
	}

	return objects, dns, nil
}

func dataSourceAciClientEndPointRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	allowEmptyResult := d.Get("allow_empty_result").(bool)

	var queryString string
	if mac, ok := d.GetOk("mac"); ok {
		if queryString != "" {
			queryString = fmt.Sprintf("%s,eq(fvCEp.mac, \"%s\")", queryString, mac.(string))
		} else {
			queryString = fmt.Sprintf("eq(fvCEp.mac, \"%s\")", mac.(string))
		}
	}

	if ip, ok := d.GetOk("ip"); ok {
		fvIpDns, err := getRemotefvIpParentDn(aciClient, fmt.Sprintf("eq(fvIp.addr, \"%s\")", ip.(string)))
		if err != nil {
			return err
		}
		if len(fvIpDns) != 0 {
			macPattern := regexp.MustCompile("cep-(.+)/")
			tempQueryString := fmt.Sprintf("eq(fvCEp.mac, \"%s\")", macPattern.FindStringSubmatch(fvIpDns[0])[1])
			for i := 1; i < len(fvIpDns); i++ {
				tempQueryString = fmt.Sprintf("%s,eq(fvCEp.mac, \"%s\")", tempQueryString, macPattern.FindStringSubmatch(fvIpDns[i])[1])
			}
			if queryString != "" {
				queryString = fmt.Sprintf("%s,%s", queryString, tempQueryString)
			} else {
				queryString = tempQueryString
			}
		} else {
			d.SetId("")
			if allowEmptyResult {
				return nil
			} else {
				return errors.New("Error retrieving Object: Object may not exists")
			}
		}
		d.Set("ip", ip)
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

	objects, dns, err := getRemoteClientEndPoint(aciClient, queryString, allowEmptyResult)
	if err != nil {
		return err
	}

	if objects != nil {
		d.Set("fvcep_objects", objects)
	}

	if dns != nil {
		d.SetId(strings.Join(dns, " "))
	} else {
		d.SetId("")
	}

	return nil
}

func getRemotefvIpParentDn(client *client.Client, query string) (fvIpDns []string, err error) {

	if query == "" {
		return fvIpDns, errors.New("Failed to build fvIp query string.")
	}

	duURL := fmt.Sprintf("%s/fvIp.json?query-target-filter=and(%s)", models.BaseurlStr, query)
	fvIpsContainer, _ := client.GetViaURL(duURL)
	if err != nil {
		return fvIpDns, err
	}

	fvIpObjectsCount, err := fvIpsContainer.ArrayCount("imdata")
	if err != nil {
		return fvIpDns, err
	}

	for i := 0; i < fvIpObjectsCount; i++ {
		fvIpContainer, err := fvIpsContainer.ArrayElement(i, "imdata")
		if err != nil {
			return fvIpDns, err
		}
		fvIpDns = append(fvIpDns, models.G(fvIpContainer.S("fvIp", "attributes"), "dn"))
	}

	if err != nil {
		return fvIpDns, err
	}

	return fvIpDns, nil
}

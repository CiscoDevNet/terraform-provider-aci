package aci

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var esgParentDNPattern = regexp.MustCompile(`^uni/tn-[^/]+/ap-[^/]+/esg-[^/]+$`)

func dataSourceAciClientEndPoint() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciClientEndPointRead,

		SchemaVersion: 1,

		Schema: AppendAttrSchemas(map[string]*schema.Schema{
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

			"filter_dn": &schema.Schema{
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

type endpointView struct {
	object map[string]interface{}
	dn     string
}

func extractInfo(con *container.Container, fvIpObjects []*container.Container) []endpointView {
	endpointDN := models.G(con, "dn")

	// The primary view reflects the parent encoded in the endpoint DN.
	parentView := extractEndpointAttributes(con)
	if !parseACIPath(endpointDN, parentView) {
		return nil
	}
	if parentView["esg_name"] != nil {
		// Newer APIC releases omit duplicate fvCEp records under ESGs to reduce scale impact.
		// esgUsegDn retains the ESG association, so synthesize the ESG view from the EPG endpoint.
		return nil
	}

	views := []endpointView{{
		object: parentView,
		dn:     endpointDN,
	}}

	esgUsegDN := models.G(con, "esgUsegDn")
	endpointRN := lastDNComponent(endpointDN)
	if !strings.HasPrefix(endpointRN, "cep-") {
		return views
	}

	esgDNs := make([]string, 0)
	esgIPs := make(map[string][]string)
	if esgUsegDN != "" && esgUsegDN != "{}" {
		esgDNs = append(esgDNs, esgUsegDN)
		esgIPs[esgUsegDN] = nil
	}
	for _, fvIp := range fvIpObjects {
		attributes := fvIp.S("fvIp", "attributes")
		esgDN := models.G(attributes, "esgUsegDn")
		if esgDN == "" || esgDN == "{}" || esgDN == esgUsegDN {
			continue
		}
		if _, exists := esgIPs[esgDN]; !exists {
			esgDNs = append(esgDNs, esgDN)
		}
		esgIPs[esgDN] = append(esgIPs[esgDN], models.G(attributes, "addr"))
	}

	for _, esgDN := range esgDNs {
		// The ESG view uses the ESG parent path and the same endpoint RN.
		esgView := extractEndpointAttributes(con)
		if !parseACIPath(esgDN, esgView) || esgView["esg_name"] == nil {
			continue
		}
		if parentView["epg_name"] != nil {
			esgView["base_epg"] = map[string]interface{}{
				"tenant_name":              parentView["tenant_name"],
				"application_profile_name": parentView["application_profile_name"],
				"epg_name":                 parentView["epg_name"],
			}
		}
		if ips := esgIPs[esgDN]; ips != nil {
			esgView["ips"] = ips
			esgView["ip"] = ips[0]
		}
		views = append(views, endpointView{
			object: esgView,
			dn:     strings.TrimSuffix(esgDN, "/") + "/" + endpointRN,
		})
	}

	return views
}

func extractEndpointAttributes(con *container.Container) map[string]interface{} {
	info := map[string]interface{}{
		"name": models.G(con, "name"),
		"mac":  models.G(con, "mac"),
		"ip":   models.G(con, "ip"),
		"vlan": models.G(con, "encap"),
	}

	if info["ip"] == "{}" {
		info["ip"] = ""
	}

	return info
}

func parseACIPath(dn string, info map[string]interface{}) bool {
	parts := strings.Split(dn, "/")
	if len(parts) < 2 || parts[0] != "uni" {
		return false
	}

	for _, part := range parts[1:] {
		kind, name, ok := strings.Cut(part, "-")
		if !ok {
			continue
		}

		switch kind {
		case "tn":
			info["tenant_name"] = name
		case "ctx":
			info["vrf_name"] = name
		case "ap":
			info["application_profile_name"] = name
		case "epg":
			info["epg_name"] = name
		case "esg":
			info["esg_name"] = name
		case "l2out":
			info["l2out_name"] = name
		case "instP":
			info["instance_profile_name"] = name
		}
	}

	_, hasTenant := info["tenant_name"]
	return hasTenant
}

func lastDNComponent(dn string) string {
	parts := strings.Split(dn, "/")
	return parts[len(parts)-1]
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

func getRemoteClientEndPoints(client *client.Client, queryTargetFilter, rspSubtreeFilter string) ([]*container.Container, error) {
	const (
		pageSize                = 10000
		endpointURL             = models.BaseurlStr + "/fvCEp.json"
		ipSubtree               = "rsp-subtree=children&rsp-subtree-class=fvIp"
		filteredIPSubtreeFormat = ipSubtree + "&rsp-subtree-filter=%s&rsp-subtree-include=required"
	)

	var duURL string
	switch {
	case queryTargetFilter != "" && rspSubtreeFilter != "":
		duURL = fmt.Sprintf(endpointURL+"?query-target-filter=and(%s)&"+filteredIPSubtreeFormat, queryTargetFilter, rspSubtreeFilter)
	case queryTargetFilter != "":
		duURL = fmt.Sprintf(endpointURL+"?query-target-filter=and(%s)&"+ipSubtree, queryTargetFilter)
	case rspSubtreeFilter != "":
		// Select CEPs through their IP children when there is no CEP-level
		// filter. The required option excludes CEPs without a matching child.
		duURL = fmt.Sprintf(endpointURL+"?"+filteredIPSubtreeFormat, rspSubtreeFilter)
	default:
		// IP children also provide ESG associations when no filters are supplied.
		duURL = endpointURL + "?" + ipSubtree
	}

	// Page CEPs by DN to avoid requesting the entire class and its IP children
	// in one response. The page size limits parent CEPs, not their IP children.
	var endpoints []*container.Container
	for page := 0; ; page++ {
		pageURL := fmt.Sprintf("%s&order-by=fvCEp.dn|asc&page-size=%d&page=%d", duURL, pageSize, page)
		fvCEpCont, err := client.GetViaURL(pageURL)
		if err != nil {
			if fvCEpCont == nil {
				return nil, err
			}
			// The SDK reports an empty imdata array as a missing object. Accept
			// valid empty pages, but propagate errors returned as APIC objects.
			if count, countErr := fvCEpCont.ArrayCount("imdata"); countErr != nil || count != 0 {
				return nil, err
			}
		}

		totalCount, err := strconv.Atoi(models.G(fvCEpCont, "totalCount"))
		if err != nil || totalCount < 0 {
			return nil, fmt.Errorf("invalid client endpoint totalCount %q", models.G(fvCEpCont, "totalCount"))
		}
		count, err := fvCEpCont.ArrayCount("imdata")
		if err != nil {
			return nil, err
		}
		for i := 0; i < count; i++ {
			endpoint, err := fvCEpCont.ArrayElement(i, "imdata")
			if err != nil {
				return nil, err
			}
			endpoints = append(endpoints, endpoint)
		}

		// Do not end paging solely because a child-filtered response is short.
		if (page+1)*pageSize >= totalCount {
			break
		}
	}
	return endpoints, nil
}

// mergeEndpoints deduplicates the ESG selections by source CEP DN, preserving
// separate base EPGs for the same MAC. Parent results come first to retain all
// IP children allowed by the IP filter; child ESG results contain a subset.
func mergeEndpoints(parentEndpoints, childEndpoints []*container.Container) []*container.Container {
	var merged []*container.Container
	seen := make(map[string]struct{})
	for _, endpoints := range [][]*container.Container{parentEndpoints, childEndpoints} {
		for _, endpoint := range endpoints {
			dn := models.G(endpoint.S("fvCEp", "attributes"), "dn")
			if _, exists := seen[dn]; exists {
				continue
			}
			seen[dn] = struct{}{}
			merged = append(merged, endpoint)
		}
	}
	return merged
}

func extractClientEndPoints(client *client.Client, endpoints []*container.Container, filterDN string) ([]interface{}, []string, error) {
	objects := make([]interface{}, 0, 1)
	dns := make([]string, 0, 1)
	for _, clientEndPointCont := range endpoints {

		// Read fvIp children once for the endpoint and its ESG views.
		var ips []string
		fvIpObjects, err := clientEndPointCont.S("fvCEp", "children").Children()
		if err == nil {
			ips = make([]string, 0, len(fvIpObjects))
			for _, fvIpObject := range fvIpObjects {
				ips = append(ips, models.G(fvIpObject.S("fvIp", "attributes"), "addr"))
			}
		} else if !errors.Is(err, container.ErrNotObjOrArray) {
			return nil, nil, err
		}
		views := extractInfo(clientEndPointCont.S("fvCEp", "attributes"), fvIpObjects)
		if len(views) == 0 {
			continue
		}

		for _, view := range views {
			if filterDN != "" && !strings.Contains(view.dn, filterDN) {
				continue
			}

			if _, hasIPs := view.object["ips"]; ips != nil && !hasIPs {
				view.object["ips"] = append([]string(nil), ips...)
				if len(ips) > 0 {
					view.object["ip"] = ips[0]
				}
			}

			if view.dn == "" {
				continue
			}

			durl := fmt.Sprintf("%s/%s/fvRsCEpToPathEp.json", models.BaseurlStr, view.dn)
			cepToPathEpCont, err := client.GetViaURL(durl)
			if err == nil {
				endpointPaths, err := extractEndpointPaths(cepToPathEpCont)
				if err == nil {
					view.object["endpoint_path"] = endpointPaths
				}
			}
			objects = append(objects, view.object)
			dns = append(dns, view.dn)
		}
	}

	return objects, dns, nil
}

func dataSourceAciClientEndPointRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	filterDN := d.Get("filter_dn").(string)
	queryTargetFilter := filterBuilder("", "eq", "fvCEp.mac", d.Get("mac").(string))
	queryTargetFilter = filterBuilder(queryTargetFilter, "eq", "fvCEp.name", d.Get("name").(string))
	if vlan := d.Get("vlan").(string); vlan != "" {
		queryTargetFilter = filterBuilder(queryTargetFilter, "eq", "fvCEp.encap", "vlan-"+vlan)
	}

	filterAttribute := "fvCEp.dn"
	isESGFilter := esgParentDNPattern.MatchString(filterDN)
	if isESGFilter {
		filterAttribute = "fvCEp.esgUsegDn"
	}
	rspSubtreeFilter := filterBuilder("", "eq", "fvIp.addr", d.Get("ip").(string))

	// For ESG filters, APIC cannot OR a CEP property filter with an IP-child
	// subtree filter in one request. Select CEP-level matches first, then query
	// IP-level matches independently below so an empty fvCEp.esgUsegDn does not
	// hide a matching fvIp.esgUsegDn. Merging both selections implements OR.
	endpoints, err := getRemoteClientEndPoints(
		aciClient,
		filterBuilder(queryTargetFilter, "wcard", filterAttribute, filterDN),
		rspSubtreeFilter,
	)
	if err != nil {
		return err
	}
	if isESGFilter {
		childEndpoints, err := getRemoteClientEndPoints(
			aciClient,
			queryTargetFilter,
			fmt.Sprintf("and(%s)", filterBuilder(rspSubtreeFilter, "wcard", "fvIp.esgUsegDn", filterDN)),
		)
		if err != nil {
			return err
		}
		endpoints = mergeEndpoints(endpoints, childEndpoints)
	}
	if len(endpoints) == 0 {
		d.SetId("")
		if d.Get("allow_empty_result").(bool) {
			return nil
		}
		return errors.New("Error retrieving Object: Object may not exist")
	}

	objects, dns, err := extractClientEndPoints(aciClient, endpoints, filterDN)
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

func filterBuilder(current, operator, attribute, value string) string {
	if value == "" {
		return current
	}
	filter := fmt.Sprintf("%s(%s,%q)", operator, attribute, value)
	if current == "" {
		return filter
	}
	return current + "," + filter
}

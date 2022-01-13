package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciFilterEntry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciFilterEntryCreate,
		UpdateContext: resourceAciFilterEntryUpdate,
		ReadContext:   resourceAciFilterEntryRead,
		DeleteContext: resourceAciFilterEntryDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFilterEntryImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"filter_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"apply_to_frag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"arp_opc": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
					"req",
					"reply",
				}, false),
			},

			"d_from_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"d_to_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ether_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
					"ipv4",
					"trill",
					"arp",
					"ipv6",
					"mpls_ucast",
					"mac_security",
					"fcoe",
					"ip",
				}, false),
			},

			"icmpv4_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"echo-rep",
					"dst-unreach",
					"src-quench",
					"echo",
					"time-exceeded",
					"unspecified",
				}, false),
			},

			"icmpv6_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
					"dst-unreach",
					"time-exceeded",
					"echo-req",
					"echo-rep",
					"nbr-solicit",
					"nbr-advert",
					"redirect",
				}, false),
			},

			"match_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CS0",
					"CS1",
					"AF11",
					"AF12",
					"AF13",
					"CS2",
					"AF21",
					"AF22",
					"AF23",
					"CS3",
					"CS4",
					"CS5",
					"CS6",
					"CS7",
					"AF31",
					"AF32",
					"AF33",
					"AF41",
					"AF42",
					"AF43",
					"VA",
					"EF",
					"unspecified",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prot": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
					"icmp",
					"igmp",
					"tcp",
					"egp",
					"igp",
					"udp",
					"icmpv6",
					"eigrp",
					"ospfigp",
					"pim",
					"l2tp",
				}, false),
			},

			"s_from_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"s_to_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"stateful": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"tcp_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"ack",
						"est",
						"fin",
						"rst",
						"syn",
						"unspecified",
					}, false),
				},
			},
		}),
	}
}
func getRemoteFilterEntry(client *client.Client, dn string) (*models.FilterEntry, error) {
	vzEntryCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzEntry := models.FilterEntryFromContainer(vzEntryCont)

	if vzEntry.DistinguishedName == "" {
		return nil, fmt.Errorf("FilterEntry %s not found", vzEntry.DistinguishedName)
	}

	return vzEntry, nil
}

func setFilterEntryAttributes(vzEntry *models.FilterEntry, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(vzEntry.DistinguishedName)
	d.Set("description", vzEntry.Description)
	if dn != vzEntry.DistinguishedName {
		d.Set("filter_dn", "")
	}
	vzEntryMap, err := vzEntry.ToMap()
	if err != nil {
		return d, err
	}
	log.Println("Check .... :", d.Get("d_from_port"))

	d.Set("filter_dn", GetParentDn(dn, fmt.Sprintf("/e-%s", vzEntryMap["name"])))

	d.Set("name", vzEntryMap["name"])

	d.Set("annotation", vzEntryMap["annotation"])
	d.Set("apply_to_frag", vzEntryMap["applyToFrag"])
	d.Set("arp_opc", vzEntryMap["arpOpc"])
	d.Set("ether_t", vzEntryMap["etherT"])
	d.Set("icmpv4_t", vzEntryMap["icmpv4T"])
	d.Set("icmpv6_t", vzEntryMap["icmpv6T"])
	d.Set("match_dscp", vzEntryMap["matchDscp"])
	d.Set("name_alias", vzEntryMap["nameAlias"])
	d.Set("prot", vzEntryMap["prot"])
	d.Set("stateful", vzEntryMap["stateful"])
	tcpRulesGet := make([]string, 0, 1)
	for _, val := range strings.Split(vzEntryMap["tcpRules"], ",") {
		if val == "" {
			tcpRulesGet = append(tcpRulesGet, "unspecified")
		} else {
			tcpRulesGet = append(tcpRulesGet, strings.Trim(val, " "))
		}
	}
	sort.Strings(tcpRulesGet)
	if tcpRulesIntr, ok := d.GetOk("tcp_rules"); ok {
		tcpRulesAct := make([]string, 0, 1)
		for _, val := range tcpRulesIntr.([]interface{}) {
			tcpRulesAct = append(tcpRulesAct, val.(string))
		}
		sort.Strings(tcpRulesAct)
		if reflect.DeepEqual(tcpRulesAct, tcpRulesGet) {
			d.Set("tcp_rules", d.Get("tcp_rules").([]interface{}))
		} else {
			d.Set("tcp_rules", tcpRulesGet)
		}
	} else {
		d.Set("tcp_rules", tcpRulesGet)
	}
	return d, nil
}

func portConversionCheck(vzEntry *models.FilterEntry, d *schema.ResourceData) (*schema.ResourceData, error) {
	constantPortMapping := map[string]string{
		"smtp":        "25",
		"dns":         "53",
		"http":        "80",
		"https":       "443",
		"pop3":        "110",
		"rtsp":        "554",
		"ftpData":     "20",
		"ssh":         "22",
		"unspecified": "0",
	}
	vzEntryMap, err := vzEntry.ToMap()
	if err != nil {
		return d, err
	}
	if DFromPortTf, ok := d.GetOk("d_from_port"); ok {
		if DFromPortTf != vzEntryMap["dFromPort"] {
			if DFromPortTf != constantPortMapping[vzEntryMap["dFromPort"]] {
				d.Set("d_from_port", vzEntryMap["dFromPort"])
			} else {
				d.Set("d_from_port", DFromPortTf)
			}
		} else {
			d.Set("d_from_port", DFromPortTf)
		}
	} else {
		d.Set("d_from_port", vzEntryMap["dFromPort"])
	}

	if DToPortTf, ok := d.GetOk("d_to_port"); ok {
		if DToPortTf != vzEntryMap["dToPort"] {
			if DToPortTf != constantPortMapping[vzEntryMap["dToPort"]] {
				d.Set("d_to_port", vzEntryMap["dToPort"])
			} else {
				d.Set("d_to_port", DToPortTf)
			}
		} else {
			d.Set("d_to_port", DToPortTf)
		}
	} else {
		d.Set("d_to_port", vzEntryMap["dToPort"])
	}

	if SFromPortTf, ok := d.GetOk("s_from_port"); ok {
		if SFromPortTf != vzEntryMap["sFromPort"] {
			if SFromPortTf != constantPortMapping[vzEntryMap["sFromPort"]] {
				d.Set("s_from_port", vzEntryMap["sFromPort"])
			} else {
				d.Set("s_from_port", SFromPortTf)
			}
		} else {
			d.Set("s_from_port", SFromPortTf)
		}
	} else {
		d.Set("s_from_port", vzEntryMap["sFromPort"])
	}

	if SToPortTf, ok := d.GetOk("s_to_port"); ok {
		if SToPortTf != vzEntryMap["sToPort"] {
			if SToPortTf != constantPortMapping[vzEntryMap["sToPort"]] {
				d.Set("s_to_port", vzEntryMap["sToPort"])
			} else {
				d.Set("s_to_port", SToPortTf)
			}
		} else {
			d.Set("s_to_port", SToPortTf)
		}
	} else {
		d.Set("s_to_port", vzEntryMap["sToPort"])
	}
	return d, nil

}

func resourceAciFilterEntryImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzEntry, err := getRemoteFilterEntry(aciClient, dn)

	if err != nil {
		return nil, err
	}
	vzEntryMap, err := vzEntry.ToMap()
	if err != nil {
		return nil, err
	}
	name := vzEntryMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/e-%s", name))
	d.Set("filter_dn", pDN)
	schemaFilled, err := setFilterEntryAttributes(vzEntry, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFilterEntryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FilterEntry: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	FilterDn := d.Get("filter_dn").(string)

	vzEntryAttr := models.FilterEntryAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzEntryAttr.Annotation = Annotation.(string)
	} else {
		vzEntryAttr.Annotation = "{}"
	}
	if ApplyToFrag, ok := d.GetOk("apply_to_frag"); ok {
		vzEntryAttr.ApplyToFrag = ApplyToFrag.(string)
	}
	if ArpOpc, ok := d.GetOk("arp_opc"); ok {
		vzEntryAttr.ArpOpc = ArpOpc.(string)
	}
	if DFromPort, ok := d.GetOk("d_from_port"); ok {
		vzEntryAttr.DFromPort = DFromPort.(string)
	}
	if DToPort, ok := d.GetOk("d_to_port"); ok {
		vzEntryAttr.DToPort = DToPort.(string)
	}
	if EtherT, ok := d.GetOk("ether_t"); ok {
		vzEntryAttr.EtherT = EtherT.(string)
	}
	if Icmpv4T, ok := d.GetOk("icmpv4_t"); ok {
		vzEntryAttr.Icmpv4T = Icmpv4T.(string)
	}
	if Icmpv6T, ok := d.GetOk("icmpv6_t"); ok {
		vzEntryAttr.Icmpv6T = Icmpv6T.(string)
	}
	if MatchDscp, ok := d.GetOk("match_dscp"); ok {
		vzEntryAttr.MatchDscp = MatchDscp.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzEntryAttr.NameAlias = NameAlias.(string)
	}
	if Prot, ok := d.GetOk("prot"); ok {
		vzEntryAttr.Prot = Prot.(string)
	}
	if SFromPort, ok := d.GetOk("s_from_port"); ok {
		vzEntryAttr.SFromPort = SFromPort.(string)
	}
	if SToPort, ok := d.GetOk("s_to_port"); ok {
		vzEntryAttr.SToPort = SToPort.(string)
	}
	if Stateful, ok := d.GetOk("stateful"); ok {
		vzEntryAttr.Stateful = Stateful.(string)
	}
	if TcpRules, ok := d.GetOk("tcp_rules"); ok {
		tcpRulesList := make([]string, 0, 1)
		for _, val := range TcpRules.([]interface{}) {
			tcpRulesList = append(tcpRulesList, val.(string))
		}
		TcpRules := strings.Join(tcpRulesList, ",")
		vzEntryAttr.TcpRules = TcpRules
	}
	vzEntry := models.NewFilterEntry(fmt.Sprintf("e-%s", name), FilterDn, desc, vzEntryAttr)

	err := aciClient.Save(vzEntry)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vzEntry.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFilterEntryRead(ctx, d, m)
}

func resourceAciFilterEntryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FilterEntry: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	FilterDn := d.Get("filter_dn").(string)

	vzEntryAttr := models.FilterEntryAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzEntryAttr.Annotation = Annotation.(string)
	} else {
		vzEntryAttr.Annotation = "{}"
	}
	if ApplyToFrag, ok := d.GetOk("apply_to_frag"); ok {
		vzEntryAttr.ApplyToFrag = ApplyToFrag.(string)
	}
	if ArpOpc, ok := d.GetOk("arp_opc"); ok {
		vzEntryAttr.ArpOpc = ArpOpc.(string)
	}
	if DFromPort, ok := d.GetOk("d_from_port"); ok {
		vzEntryAttr.DFromPort = DFromPort.(string)
	}
	if DToPort, ok := d.GetOk("d_to_port"); ok {
		vzEntryAttr.DToPort = DToPort.(string)
	}
	if EtherT, ok := d.GetOk("ether_t"); ok {
		vzEntryAttr.EtherT = EtherT.(string)
	}
	if Icmpv4T, ok := d.GetOk("icmpv4_t"); ok {
		vzEntryAttr.Icmpv4T = Icmpv4T.(string)
	}
	if Icmpv6T, ok := d.GetOk("icmpv6_t"); ok {
		vzEntryAttr.Icmpv6T = Icmpv6T.(string)
	}
	if MatchDscp, ok := d.GetOk("match_dscp"); ok {
		vzEntryAttr.MatchDscp = MatchDscp.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzEntryAttr.NameAlias = NameAlias.(string)
	}
	if Prot, ok := d.GetOk("prot"); ok {
		vzEntryAttr.Prot = Prot.(string)
	}
	if SFromPort, ok := d.GetOk("s_from_port"); ok {
		vzEntryAttr.SFromPort = SFromPort.(string)
	}
	if SToPort, ok := d.GetOk("s_to_port"); ok {
		vzEntryAttr.SToPort = SToPort.(string)
	}
	if Stateful, ok := d.GetOk("stateful"); ok {
		vzEntryAttr.Stateful = Stateful.(string)
	}
	if TcpRules, ok := d.GetOk("tcp_rules"); ok {
		tcpRulesList := make([]string, 0, 1)
		for _, val := range TcpRules.([]interface{}) {
			tcpRulesList = append(tcpRulesList, val.(string))
		}
		TcpRules := strings.Join(tcpRulesList, ",")
		vzEntryAttr.TcpRules = TcpRules
	}
	vzEntry := models.NewFilterEntry(fmt.Sprintf("e-%s", name), FilterDn, desc, vzEntryAttr)

	vzEntry.Status = "modified"

	err := aciClient.Save(vzEntry)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vzEntry.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFilterEntryRead(ctx, d, m)
}

func resourceAciFilterEntryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzEntry, err := getRemoteFilterEntry(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	d, err = portConversionCheck(vzEntry, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setFilterEntryAttributes(vzEntry, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFilterEntryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzEntry")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

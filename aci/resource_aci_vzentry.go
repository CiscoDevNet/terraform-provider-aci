package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciFilterEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciFilterEntryCreate,
		Update: resourceAciFilterEntryUpdate,
		Read:   resourceAciFilterEntryRead,
		Delete: resourceAciFilterEntryDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFilterEntryImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"filter_dn": &schema.Schema{
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

			"apply_to_frag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"arp_opc": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			},

			"icmpv4_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"icmpv6_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"match_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			},

			"tcp_rules": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func setFilterEntryAttributes(vzEntry *models.FilterEntry, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(vzEntry.DistinguishedName)
	d.Set("description", vzEntry.Description)
	d.Set("filter_dn", GetParentDn(vzEntry.DistinguishedName))
	vzEntryMap, _ := vzEntry.ToMap()

	d.Set("name", vzEntryMap["name"])

	d.Set("annotation", vzEntryMap["annotation"])
	d.Set("apply_to_frag", vzEntryMap["applyToFrag"])
	d.Set("arp_opc", vzEntryMap["arpOpc"])
	d.Set("d_from_port", vzEntryMap["dFromPort"])
	d.Set("d_to_port", vzEntryMap["dToPort"])
	d.Set("ether_t", vzEntryMap["etherT"])
	d.Set("icmpv4_t", vzEntryMap["icmpv4T"])
	d.Set("icmpv6_t", vzEntryMap["icmpv6T"])
	d.Set("match_dscp", vzEntryMap["matchDscp"])
	d.Set("name_alias", vzEntryMap["nameAlias"])
	d.Set("prot", vzEntryMap["prot"])
	d.Set("s_from_port", vzEntryMap["sFromPort"])
	d.Set("s_to_port", vzEntryMap["sToPort"])
	d.Set("stateful", vzEntryMap["stateful"])
	d.Set("tcp_rules", vzEntryMap["tcpRules"])
	return d
}

func resourceAciFilterEntryImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzEntry, err := getRemoteFilterEntry(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setFilterEntryAttributes(vzEntry, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFilterEntryCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FilterEntry: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	FilterDn := d.Get("filter_dn").(string)

	vzEntryAttr := models.FilterEntryAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzEntryAttr.Annotation = Annotation.(string)
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
		vzEntryAttr.TcpRules = TcpRules.(string)
	}
	vzEntry := models.NewFilterEntry(fmt.Sprintf("e-%s", name), FilterDn, desc, vzEntryAttr)

	err := aciClient.Save(vzEntry)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(vzEntry.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFilterEntryRead(d, m)
}

func resourceAciFilterEntryUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FilterEntry: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	FilterDn := d.Get("filter_dn").(string)

	vzEntryAttr := models.FilterEntryAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzEntryAttr.Annotation = Annotation.(string)
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
		vzEntryAttr.TcpRules = TcpRules.(string)
	}
	vzEntry := models.NewFilterEntry(fmt.Sprintf("e-%s", name), FilterDn, desc, vzEntryAttr)

	vzEntry.Status = "modified"

	err := aciClient.Save(vzEntry)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(vzEntry.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFilterEntryRead(d, m)

}

func resourceAciFilterEntryRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzEntry, err := getRemoteFilterEntry(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setFilterEntryAttributes(vzEntry, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFilterEntryDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzEntry")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

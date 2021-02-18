package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciDHCPRelayPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciDHCPRelayPolicyCreate,
		Update: resourceAciDHCPRelayPolicyUpdate,
		Read:   resourceAciDHCPRelayPolicyRead,
		Delete: resourceAciDHCPRelayPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciDHCPRelayPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"visible",
					"not-visible",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"infra",
					"tenant",
				}, false),
			},

			"relation_dhcp_rs_prov": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteDHCPRelayPolicy(client *client.Client, dn string) (*models.DHCPRelayPolicy, error) {
	dhcpRelayPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpRelayP := models.DHCPRelayPolicyFromContainer(dhcpRelayPCont)

	if dhcpRelayP.DistinguishedName == "" {
		return nil, fmt.Errorf("DHCPRelayPolicy %s not found", dhcpRelayP.DistinguishedName)
	}

	return dhcpRelayP, nil
}

func setDHCPRelayPolicyAttributes(dhcpRelayP *models.DHCPRelayPolicy, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()

	d.SetId(dhcpRelayP.DistinguishedName)
	d.Set("description", dhcpRelayP.Description)

	if dn != dhcpRelayP.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	dhcpRelayPMap, _ := dhcpRelayP.ToMap()

	d.Set("name", dhcpRelayPMap["name"])

	d.Set("annotation", dhcpRelayPMap["annotation"])
	d.Set("mode", dhcpRelayPMap["mode"])
	d.Set("name_alias", dhcpRelayPMap["nameAlias"])
	d.Set("owner", dhcpRelayPMap["owner"])
	return d
}

func resourceAciDHCPRelayPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	dhcpRelayP, err := getRemoteDHCPRelayPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setDHCPRelayPolicyAttributes(dhcpRelayP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciDHCPRelayPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] DHCPRelayPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	dhcpRelayPAttr := models.DHCPRelayPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpRelayPAttr.Annotation = Annotation.(string)
	} else {
		dhcpRelayPAttr.Annotation = "{}"
	}
	if Mode, ok := d.GetOk("mode"); ok {
		dhcpRelayPAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpRelayPAttr.NameAlias = NameAlias.(string)
	}
	if Owner, ok := d.GetOk("owner"); ok {
		dhcpRelayPAttr.Owner = Owner.(string)
	}
	dhcpRelayP := models.NewDHCPRelayPolicy(fmt.Sprintf("relayp-%s", name), TenantDn, desc, dhcpRelayPAttr)

	err := aciClient.Save(dhcpRelayP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationTodhcpRsProv, ok := d.GetOk("relation_dhcp_rs_prov"); ok {
		relationParam := relationTodhcpRsProv.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if relationTodhcpRsProv, ok := d.GetOk("relation_dhcp_rs_prov"); ok {
		relationParamList := toStringList(relationTodhcpRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationdhcpRsProvFromDHCPRelayPolicy(dhcpRelayP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_dhcp_rs_prov")
			d.Partial(false)
		}
	}

	d.SetId(dhcpRelayP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciDHCPRelayPolicyRead(d, m)
}

func resourceAciDHCPRelayPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] DHCPRelayPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	dhcpRelayPAttr := models.DHCPRelayPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpRelayPAttr.Annotation = Annotation.(string)
	} else {
		dhcpRelayPAttr.Annotation = "{}"
	}
	if Mode, ok := d.GetOk("mode"); ok {
		dhcpRelayPAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpRelayPAttr.NameAlias = NameAlias.(string)
	}
	if Owner, ok := d.GetOk("owner"); ok {
		dhcpRelayPAttr.Owner = Owner.(string)
	}
	dhcpRelayP := models.NewDHCPRelayPolicy(fmt.Sprintf("relayp-%s", name), TenantDn, desc, dhcpRelayPAttr)

	dhcpRelayP.Status = "modified"

	err := aciClient.Save(dhcpRelayP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_dhcp_rs_prov") {
		_, newRelParam := d.GetChange("relation_dhcp_rs_prov")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if d.HasChange("relation_dhcp_rs_prov") {
		oldRel, newRel := d.GetChange("relation_dhcp_rs_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationdhcpRsProvFromDHCPRelayPolicy(dhcpRelayP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationdhcpRsProvFromDHCPRelayPolicy(dhcpRelayP.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_dhcp_rs_prov")
			d.Partial(false)

		}

	}

	d.SetId(dhcpRelayP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciDHCPRelayPolicyRead(d, m)

}

func resourceAciDHCPRelayPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	dhcpRelayP, err := getRemoteDHCPRelayPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setDHCPRelayPolicyAttributes(dhcpRelayP, d)

	dhcpRsProvData, err := aciClient.ReadRelationdhcpRsProvFromDHCPRelayPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation dhcpRsProv %v", err)
		d.Set("relation_dhcp_rs_prov", "")

	} else {
		if _, ok := d.GetOk("relation_dhcp_rs_prov"); ok {
			tfName := d.Get("relation_dhcp_rs_prov").(string)
			if tfName != dhcpRsProvData {
				d.Set("relation_dhcp_rs_prov", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciDHCPRelayPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "dhcpRelayP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

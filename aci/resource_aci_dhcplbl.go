package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciBDDHCPLabel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBDDHCPLabelCreate,
		UpdateContext: resourceAciBDDHCPLabelUpdate,
		ReadContext:   resourceAciBDDHCPLabelRead,
		DeleteContext: resourceAciBDDHCPLabelDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBDDHCPLabelImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"bridge_domain_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_dhcp_rs_dhcp_option_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		}),
	}
}
func getRemoteBDDHCPLabel(client *client.Client, dn string) (*models.BDDHCPLabel, error) {
	dhcpLblCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpLbl := models.BDDHCPLabelFromContainer(dhcpLblCont)

	if dhcpLbl.DistinguishedName == "" {
		return nil, fmt.Errorf("BDDHCPLabel %s not found", dhcpLbl.DistinguishedName)
	}

	return dhcpLbl, nil
}

func setBDDHCPLabelAttributes(dhcpLbl *models.BDDHCPLabel, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(dhcpLbl.DistinguishedName)
	d.Set("description", dhcpLbl.Description)

	if dn != dhcpLbl.DistinguishedName {
		d.Set("bridge_domain_dn", "")
	}
	dhcpLblMap, err := dhcpLbl.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("bridge_domain_dn", GetParentDn(dn, fmt.Sprintf("/dhcplbl-%s", dhcpLblMap["name"])))
	d.Set("name", dhcpLblMap["name"])

	d.Set("annotation", dhcpLblMap["annotation"])
	d.Set("name_alias", dhcpLblMap["nameAlias"])
	d.Set("owner", dhcpLblMap["owner"])
	d.Set("tag", dhcpLblMap["tag"])
	return d, nil
}

func resourceAciBDDHCPLabelImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	dhcpLbl, err := getRemoteBDDHCPLabel(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setBDDHCPLabelAttributes(dhcpLbl, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBDDHCPLabelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BDDHCPLabel: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	BridgeDomainDn := d.Get("bridge_domain_dn").(string)

	dhcpLblAttr := models.BDDHCPLabelAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpLblAttr.Annotation = Annotation.(string)
	} else {
		dhcpLblAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpLblAttr.NameAlias = NameAlias.(string)
	}
	if Owner, ok := d.GetOk("owner"); ok {
		dhcpLblAttr.Owner = Owner.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		dhcpLblAttr.Tag = Tag.(string)
	}
	dhcpLbl := models.NewBDDHCPLabel(fmt.Sprintf("dhcplbl-%s", name), BridgeDomainDn, desc, dhcpLblAttr)

	err := aciClient.Save(dhcpLbl)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTodhcpRsDhcpOptionPol, ok := d.GetOk("relation_dhcp_rs_dhcp_option_pol"); ok {
		relationParam := relationTodhcpRsDhcpOptionPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTodhcpRsDhcpOptionPol, ok := d.GetOk("relation_dhcp_rs_dhcp_option_pol"); ok {
		relationParam := relationTodhcpRsDhcpOptionPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationdhcpRsDhcpOptionPolFromBDDHCPLabel(dhcpLbl.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(dhcpLbl.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBDDHCPLabelRead(ctx, d, m)
}

func resourceAciBDDHCPLabelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BDDHCPLabel: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	BridgeDomainDn := d.Get("bridge_domain_dn").(string)

	dhcpLblAttr := models.BDDHCPLabelAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpLblAttr.Annotation = Annotation.(string)
	} else {
		dhcpLblAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpLblAttr.NameAlias = NameAlias.(string)
	}
	if Owner, ok := d.GetOk("owner"); ok {
		dhcpLblAttr.Owner = Owner.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		dhcpLblAttr.Tag = Tag.(string)
	}
	dhcpLbl := models.NewBDDHCPLabel(fmt.Sprintf("dhcplbl-%s", name), BridgeDomainDn, desc, dhcpLblAttr)

	dhcpLbl.Status = "modified"

	err := aciClient.Save(dhcpLbl)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_dhcp_rs_dhcp_option_pol") {
		_, newRelParam := d.GetChange("relation_dhcp_rs_dhcp_option_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_dhcp_rs_dhcp_option_pol") {
		_, newRelParam := d.GetChange("relation_dhcp_rs_dhcp_option_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationdhcpRsDhcpOptionPolFromBDDHCPLabel(dhcpLbl.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(dhcpLbl.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBDDHCPLabelRead(ctx, d, m)

}

func resourceAciBDDHCPLabelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	dhcpLbl, err := getRemoteBDDHCPLabel(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setBDDHCPLabelAttributes(dhcpLbl, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	dhcpRsDhcpOptionPolData, err := aciClient.ReadRelationdhcpRsDhcpOptionPolFromBDDHCPLabel(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation dhcpRsDhcpOptionPol %v", err)
		d.Set("relation_dhcp_rs_dhcp_option_pol", "")
	} else {
		d.Set("relation_dhcp_rs_dhcp_option_pol", dhcpRsDhcpOptionPolData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBDDHCPLabelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "dhcpLbl")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

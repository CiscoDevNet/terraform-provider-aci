package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciL3outPathAttachmentSecondaryIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL3outPathAttachmentSecondaryIpCreate,
		UpdateContext: resourceAciL3outPathAttachmentSecondaryIpUpdate,
		ReadContext:   resourceAciL3outPathAttachmentSecondaryIpRead,
		DeleteContext: resourceAciL3outPathAttachmentSecondaryIpDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outPathAttachmentSecondaryIpImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"l3out_path_attachment_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ipv6_dad": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},
			"dhcp_relay": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
				Default: "disabled",
			},
		})),
	}
}

func getRemoteL3outPathAttachmentSecondaryIp(client *client.Client, dn string) (*models.L3outPathAttachmentSecondaryIp, error) {
	l3extIpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extIp := models.L3outPathAttachmentSecondaryIpFromContainer(l3extIpCont)

	if l3extIp.DistinguishedName == "" {
		return nil, fmt.Errorf("L3Out Path Attachment Secondary IP %s not found", dn)
	}

	return l3extIp, nil
}

func setL3outPathAttachmentSecondaryIpAttributes(l3extIp *models.L3outPathAttachmentSecondaryIp, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(l3extIp.DistinguishedName)
	d.Set("description", l3extIp.Description)
	dn := d.Id()
	if dn != l3extIp.DistinguishedName {
		d.Set("l3out_path_attachment_dn", "")
	}
	l3extIpMap, err := l3extIp.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("l3out_path_attachment_dn", GetParentDn(dn, fmt.Sprintf("/addr-[%s]", l3extIpMap["addr"])))

	d.Set("addr", l3extIpMap["addr"])
	d.Set("annotation", l3extIpMap["annotation"])
	d.Set("ipv6_dad", l3extIpMap["ipv6Dad"])
	d.Set("name_alias", l3extIpMap["nameAlias"])
	return d, nil
}

func getAndSetReadDhcpRelayGwExtIp(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	_, err := client.ReadDhcpRelayGwExtIp(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading DHCP relay gateway %v", err)
		d.Set("dhcp_relay", "disabled")
		return nil, err
	} else {
		d.Set("dhcp_relay", "enabled")
		log.Printf("[DEBUG]: DHCP relay gateway: %s reading finished successfully", "enabled")
	}
	return d, nil
}

func resourceAciL3outPathAttachmentSecondaryIpImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extIp, err := getRemoteL3outPathAttachmentSecondaryIp(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL3outPathAttachmentSecondaryIpAttributes(l3extIp, d)
	if err != nil {
		return nil, err
	}

	// Importing dhcpRelayGwExtIp
	getAndSetReadDhcpRelayGwExtIp(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outPathAttachmentSecondaryIpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outPathAttachmentSecondaryIp: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	addr := d.Get("addr").(string)

	LeafPortDn := d.Get("l3out_path_attachment_dn").(string)

	l3extIpAttr := models.L3outPathAttachmentSecondaryIpAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extIpAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extIpAttr.Annotation = Annotation.(string)
	} else {
		l3extIpAttr.Annotation = "{}"
	}
	if Ipv6Dad, ok := d.GetOk("ipv6_dad"); ok {
		l3extIpAttr.Ipv6Dad = Ipv6Dad.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extIpAttr.NameAlias = NameAlias.(string)
	}
	l3extIp := models.NewL3outPathAttachmentSecondaryIp(fmt.Sprintf("addr-[%s]", addr), LeafPortDn, desc, l3extIpAttr)

	err := aciClient.Save(l3extIp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l3extIp.DistinguishedName)

	dhcpRelay := d.Get("dhcp_relay").(string)
	if dhcpRelay == "enabled" {
		dhcpRelayGwExtIpAttr := models.DhcpRelayGwExtIpAttributes{}
		dhcpRelayGwExtIp := models.NewDhcpRelayGwExtIp(models.RndhcpRelayGwExtIp, l3extIp.DistinguishedName, desc, dhcpRelayGwExtIpAttr)
		dhcpRelayErr := aciClient.Save(dhcpRelayGwExtIp)
		if dhcpRelayErr != nil {
			return diag.FromErr(dhcpRelayErr)
		}
		log.Printf("[DEBUG] %s: DHCP Relay enabled successfully", dhcpRelayGwExtIp.DistinguishedName)
	}

	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outPathAttachmentSecondaryIpRead(ctx, d, m)
}

func resourceAciL3outPathAttachmentSecondaryIpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outPathAttachmentSecondaryIp: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	addr := d.Get("addr").(string)

	LeafPortDn := d.Get("l3out_path_attachment_dn").(string)

	l3extIpAttr := models.L3outPathAttachmentSecondaryIpAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extIpAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extIpAttr.Annotation = Annotation.(string)
	} else {
		l3extIpAttr.Annotation = "{}"
	}
	if Ipv6Dad, ok := d.GetOk("ipv6_dad"); ok {
		l3extIpAttr.Ipv6Dad = Ipv6Dad.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extIpAttr.NameAlias = NameAlias.(string)
	}
	l3extIp := models.NewL3outPathAttachmentSecondaryIp(fmt.Sprintf("addr-[%s]", addr), LeafPortDn, desc, l3extIpAttr)

	l3extIp.Status = "modified"

	err := aciClient.Save(l3extIp)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l3extIp.DistinguishedName)

	if d.HasChange("dhcp_relay") {
		dhcpRelay := d.Get("dhcp_relay").(string)

		if dhcpRelay == "enabled" {
			dhcpRelayGwExtIpAttr := models.DhcpRelayGwExtIpAttributes{}
			dhcpRelayGwExtIp := models.NewDhcpRelayGwExtIp(models.RndhcpRelayGwExtIp, l3extIp.DistinguishedName, desc, dhcpRelayGwExtIpAttr)
			dhcpRelayErr := aciClient.Save(dhcpRelayGwExtIp)
			if dhcpRelayErr != nil {
				return diag.FromErr(dhcpRelayErr)
			}
			log.Printf("[DEBUG] %s: DHCP relay gateway enabled successfully", dhcpRelayGwExtIp.DistinguishedName)
		} else if dhcpRelay == "disabled" {
			aciClient.DeleteDhcpRelayGwExtIp(l3extIp.DistinguishedName)
			log.Printf("[DEBUG] %s: DHCP relay gateway disabled successfully", fmt.Sprintf("%s/%s", l3extIp.DistinguishedName, models.RndhcpRelayGwExtIp))
		}
	}

	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outPathAttachmentSecondaryIpRead(ctx, d, m)

}

func resourceAciL3outPathAttachmentSecondaryIpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extIp, err := getRemoteL3outPathAttachmentSecondaryIp(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setL3outPathAttachmentSecondaryIpAttributes(l3extIp, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// Importing dhcpRelayGwExtIp
	getAndSetReadDhcpRelayGwExtIp(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outPathAttachmentSecondaryIpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extIp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

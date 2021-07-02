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

func resourceAciL3outVPCMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL3outVPCMemberCreate,
		UpdateContext: resourceAciL3outVPCMemberUpdate,
		ReadContext:   resourceAciL3outVPCMemberRead,
		DeleteContext: resourceAciL3outVPCMemberDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outVPCMemberImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"leaf_port_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"side": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A",
					"B",
				}, false),
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

			"ll_addr": &schema.Schema{
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
func getRemoteL3outVPCMember(client *client.Client, dn string) (*models.L3outVPCMember, error) {
	l3extMemberCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extMember := models.L3outVPCMemberFromContainer(l3extMemberCont)

	if l3extMember.DistinguishedName == "" {
		return nil, fmt.Errorf("L3outVPCMember %s not found", l3extMember.DistinguishedName)
	}

	return l3extMember, nil
}

func setL3outVPCMemberAttributes(l3extMember *models.L3outVPCMember, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(l3extMember.DistinguishedName)
	d.Set("description", l3extMember.Description)
	dn := d.Id()
	if dn != l3extMember.DistinguishedName {
		d.Set("leaf_port_dn", "")
	}
	l3extMemberMap, err := l3extMember.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("side", l3extMemberMap["side"])

	d.Set("addr", l3extMemberMap["addr"])
	d.Set("annotation", l3extMemberMap["annotation"])
	d.Set("ipv6_dad", l3extMemberMap["ipv6Dad"])
	d.Set("ll_addr", l3extMemberMap["llAddr"])
	d.Set("name_alias", l3extMemberMap["nameAlias"])
	return d, nil
}

func resourceAciL3outVPCMemberImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extMember, err := getRemoteL3outVPCMember(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL3outVPCMemberAttributes(l3extMember, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outVPCMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outVPCMember: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	side := d.Get("side").(string)

	LeafPortDn := d.Get("leaf_port_dn").(string)

	l3extMemberAttr := models.L3outVPCMemberAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extMemberAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extMemberAttr.Annotation = Annotation.(string)
	} else {
		l3extMemberAttr.Annotation = "{}"
	}
	if Ipv6Dad, ok := d.GetOk("ipv6_dad"); ok {
		l3extMemberAttr.Ipv6Dad = Ipv6Dad.(string)
	}
	if LlAddr, ok := d.GetOk("ll_addr"); ok {
		l3extMemberAttr.LlAddr = LlAddr.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extMemberAttr.NameAlias = NameAlias.(string)
	}
	l3extMember := models.NewL3outVPCMember(fmt.Sprintf("mem-%s", side), LeafPortDn, desc, l3extMemberAttr)

	err := aciClient.Save(l3extMember)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l3extMember.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outVPCMemberRead(ctx, d, m)
}

func resourceAciL3outVPCMemberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outVPCMember: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	side := d.Get("side").(string)

	LeafPortDn := d.Get("leaf_port_dn").(string)

	l3extMemberAttr := models.L3outVPCMemberAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extMemberAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extMemberAttr.Annotation = Annotation.(string)
	} else {
		l3extMemberAttr.Annotation = "{}"
	}
	if Ipv6Dad, ok := d.GetOk("ipv6_dad"); ok {
		l3extMemberAttr.Ipv6Dad = Ipv6Dad.(string)
	}
	if LlAddr, ok := d.GetOk("ll_addr"); ok {
		l3extMemberAttr.LlAddr = LlAddr.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extMemberAttr.NameAlias = NameAlias.(string)
	}
	l3extMember := models.NewL3outVPCMember(fmt.Sprintf("mem-%s", side), LeafPortDn, desc, l3extMemberAttr)

	l3extMember.Status = "modified"

	err := aciClient.Save(l3extMember)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l3extMember.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outVPCMemberRead(ctx, d, m)

}

func resourceAciL3outVPCMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extMember, err := getRemoteL3outVPCMember(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setL3outVPCMemberAttributes(l3extMember, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outVPCMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extMember")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

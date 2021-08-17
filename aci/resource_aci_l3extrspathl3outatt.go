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

func resourceAciL3outPathAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL3outPathAttachmentCreate,
		UpdateContext: resourceAciL3outPathAttachmentUpdate,
		ReadContext:   resourceAciL3outPathAttachmentRead,
		DeleteContext: resourceAciL3outPathAttachmentDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outPathAttachmentImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"target_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"if_inst_t": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ext-svi",
					"l3-port",
					"sub-interface",
					"unspecified",
				}, false),
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"autostate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encap_scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ctx",
					"local",
				}, false),
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

			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"native",
					"regular",
					"untagged",
				}, false),
			},

			"mtu": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AF11",
					"AF12",
					"AF13",
					"AF21",
					"AF22",
					"AF23",
					"AF31",
					"AF32",
					"AF33",
					"AF41",
					"AF42",
					"AF43",
					"CS0",
					"CS1",
					"CS2",
					"CS3",
					"CS4",
					"CS5",
					"CS6",
					"CS7",
					"EF",
					"VA",
					"unspecified",
				}, false),
			},
		}),
	}
}
func getRemoteL3outPathAttachment(client *client.Client, dn string) (*models.L3outPathAttachment, error) {
	l3extRsPathL3OutAttCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extRsPathL3OutAtt := models.L3outPathAttachmentFromContainer(l3extRsPathL3OutAttCont)

	if l3extRsPathL3OutAtt.DistinguishedName == "" {
		return nil, fmt.Errorf("L3outPathAttachment %s not found", l3extRsPathL3OutAtt.DistinguishedName)
	}

	return l3extRsPathL3OutAtt, nil
}

func setL3outPathAttachmentAttributes(l3extRsPathL3OutAtt *models.L3outPathAttachment, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(l3extRsPathL3OutAtt.DistinguishedName)
	d.Set("description", l3extRsPathL3OutAtt.Description)
	dn := d.Id()
	if dn != l3extRsPathL3OutAtt.DistinguishedName {
		d.Set("logical_interface_profile_dn", "")
	}
	l3extRsPathL3OutAttMap, err := l3extRsPathL3OutAtt.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("logical_interface_profile_dn", GetParentDn(dn,fmt.Sprintf("/rspathL3OutAtt-[%s]",l3extRsPathL3OutAttMap["tDn"])))

	d.Set("target_dn", l3extRsPathL3OutAttMap["tDn"])

	d.Set("addr", l3extRsPathL3OutAttMap["addr"])
	d.Set("annotation", l3extRsPathL3OutAttMap["annotation"])
	d.Set("autostate", l3extRsPathL3OutAttMap["autostate"])
	d.Set("encap", l3extRsPathL3OutAttMap["encap"])
	d.Set("encap_scope", l3extRsPathL3OutAttMap["encapScope"])
	d.Set("if_inst_t", l3extRsPathL3OutAttMap["ifInstT"])
	d.Set("ipv6_dad", l3extRsPathL3OutAttMap["ipv6Dad"])
	d.Set("ll_addr", l3extRsPathL3OutAttMap["llAddr"])
	d.Set("mac", l3extRsPathL3OutAttMap["mac"])
	d.Set("mode", l3extRsPathL3OutAttMap["mode"])
	d.Set("mtu", l3extRsPathL3OutAttMap["mtu"])
	d.Set("target_dscp", l3extRsPathL3OutAttMap["targetDscp"])
	return d, nil
}

func resourceAciL3outPathAttachmentImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extRsPathL3OutAtt, err := getRemoteL3outPathAttachment(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL3outPathAttachmentAttributes(l3extRsPathL3OutAtt, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outPathAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outPathAttachment: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	tDn := d.Get("target_dn").(string)

	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	l3extRsPathL3OutAttAttr := models.L3outPathAttachmentAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extRsPathL3OutAttAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extRsPathL3OutAttAttr.Annotation = Annotation.(string)
	} else {
		l3extRsPathL3OutAttAttr.Annotation = "{}"
	}
	if Autostate, ok := d.GetOk("autostate"); ok {
		l3extRsPathL3OutAttAttr.Autostate = Autostate.(string)
	}
	if Encap, ok := d.GetOk("encap"); ok {
		l3extRsPathL3OutAttAttr.Encap = Encap.(string)
	}
	if EncapScope, ok := d.GetOk("encap_scope"); ok {
		l3extRsPathL3OutAttAttr.EncapScope = EncapScope.(string)
	}
	if IfInstT, ok := d.GetOk("if_inst_t"); ok {
		l3extRsPathL3OutAttAttr.IfInstT = IfInstT.(string)
	}
	if Ipv6Dad, ok := d.GetOk("ipv6_dad"); ok {
		l3extRsPathL3OutAttAttr.Ipv6Dad = Ipv6Dad.(string)
	}
	if LlAddr, ok := d.GetOk("ll_addr"); ok {
		l3extRsPathL3OutAttAttr.LlAddr = LlAddr.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		l3extRsPathL3OutAttAttr.Mac = Mac.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		l3extRsPathL3OutAttAttr.Mode = Mode.(string)
	}
	if Mtu, ok := d.GetOk("mtu"); ok {
		l3extRsPathL3OutAttAttr.Mtu = Mtu.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extRsPathL3OutAttAttr.TargetDscp = TargetDscp.(string)
	}
	l3extRsPathL3OutAtt := models.NewL3outPathAttachment(fmt.Sprintf("rspathL3OutAtt-[%s]", tDn), LogicalInterfaceProfileDn, desc, l3extRsPathL3OutAttAttr)

	err := aciClient.Save(l3extRsPathL3OutAtt)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l3extRsPathL3OutAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outPathAttachmentRead(ctx, d, m)
}

func resourceAciL3outPathAttachmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outPathAttachment: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	tDn := d.Get("target_dn").(string)

	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	l3extRsPathL3OutAttAttr := models.L3outPathAttachmentAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extRsPathL3OutAttAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extRsPathL3OutAttAttr.Annotation = Annotation.(string)
	} else {
		l3extRsPathL3OutAttAttr.Annotation = "{}"
	}
	if Autostate, ok := d.GetOk("autostate"); ok {
		l3extRsPathL3OutAttAttr.Autostate = Autostate.(string)
	}
	if Encap, ok := d.GetOk("encap"); ok {
		l3extRsPathL3OutAttAttr.Encap = Encap.(string)
	}
	if EncapScope, ok := d.GetOk("encap_scope"); ok {
		l3extRsPathL3OutAttAttr.EncapScope = EncapScope.(string)
	}
	if IfInstT, ok := d.GetOk("if_inst_t"); ok {
		l3extRsPathL3OutAttAttr.IfInstT = IfInstT.(string)
	}
	if Ipv6Dad, ok := d.GetOk("ipv6_dad"); ok {
		l3extRsPathL3OutAttAttr.Ipv6Dad = Ipv6Dad.(string)
	}
	if LlAddr, ok := d.GetOk("ll_addr"); ok {
		l3extRsPathL3OutAttAttr.LlAddr = LlAddr.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		l3extRsPathL3OutAttAttr.Mac = Mac.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		l3extRsPathL3OutAttAttr.Mode = Mode.(string)
	}
	if Mtu, ok := d.GetOk("mtu"); ok {
		l3extRsPathL3OutAttAttr.Mtu = Mtu.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extRsPathL3OutAttAttr.TargetDscp = TargetDscp.(string)
	}
	l3extRsPathL3OutAtt := models.NewL3outPathAttachment(fmt.Sprintf("rspathL3OutAtt-[%s]", tDn), LogicalInterfaceProfileDn, desc, l3extRsPathL3OutAttAttr)

	l3extRsPathL3OutAtt.Status = "modified"

	err := aciClient.Save(l3extRsPathL3OutAtt)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l3extRsPathL3OutAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outPathAttachmentRead(ctx, d, m)

}

func resourceAciL3outPathAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extRsPathL3OutAtt, err := getRemoteL3outPathAttachment(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setL3outPathAttachmentAttributes(l3extRsPathL3OutAtt, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outPathAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extRsPathL3OutAtt")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

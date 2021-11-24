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

func resourceAciVirtualLogicalInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciVirtualLogicalInterfaceProfileCreate,
		UpdateContext: resourceAciVirtualLogicalInterfaceProfileUpdate,
		ReadContext:   resourceAciVirtualLogicalInterfaceProfileRead,
		DeleteContext: resourceAciVirtualLogicalInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVirtualLogicalInterfaceProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"node_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"encap_scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ctx",
					"local",
				}, false),
			},

			"if_inst_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ext-svi",
					"l3-port",
					"sub-interface",
					"unspecified",
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

			"relation_l3ext_rs_dyn_path_att": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}

func getRemoteVirtualLogicalInterfaceProfile(client *client.Client, dn string) (*models.VirtualLogicalInterfaceProfile, error) {
	l3extVirtualLIfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extVirtualLIfP := models.VirtualLogicalInterfaceProfileFromContainer(l3extVirtualLIfPCont)

	if l3extVirtualLIfP.DistinguishedName == "" {
		return nil, fmt.Errorf("LogicalInterfaceProfile %s not found", l3extVirtualLIfP.DistinguishedName)
	}

	return l3extVirtualLIfP, nil
}

func setVirtualLogicalInterfaceProfileAttributes(l3extVirtualLIfP *models.VirtualLogicalInterfaceProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	d.SetId(l3extVirtualLIfP.DistinguishedName)
	d.Set("description", l3extVirtualLIfP.Description)
	if dn != l3extVirtualLIfP.DistinguishedName {
		d.Set("logical_interface_profile_dn", "")
	}
	l3extVirtualLIfPMap, err := l3extVirtualLIfP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("logical_interface_profile_dn", GetParentDn(dn, fmt.Sprintf("/vlifp-[%s]-[%s]", l3extVirtualLIfPMap["nodeDn"], l3extVirtualLIfPMap["encap"])))
	d.Set("node_dn", l3extVirtualLIfPMap["nodeDn"])
	d.Set("encap", l3extVirtualLIfPMap["encap"])
	d.Set("addr", l3extVirtualLIfPMap["addr"])
	d.Set("annotation", l3extVirtualLIfPMap["annotation"])
	d.Set("autostate", l3extVirtualLIfPMap["autostate"])
	d.Set("encap", l3extVirtualLIfPMap["encap"])
	d.Set("encap_scope", l3extVirtualLIfPMap["encapScope"])
	d.Set("if_inst_t", l3extVirtualLIfPMap["ifInstT"])
	d.Set("ipv6_dad", l3extVirtualLIfPMap["ipv6Dad"])
	d.Set("ll_addr", l3extVirtualLIfPMap["llAddr"])
	d.Set("mac", l3extVirtualLIfPMap["mac"])
	d.Set("mode", l3extVirtualLIfPMap["mode"])
	d.Set("mtu", l3extVirtualLIfPMap["mtu"])
	d.Set("target_dscp", l3extVirtualLIfPMap["targetDscp"])

	return d, nil
}

func resourceAciVirtualLogicalInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extVirtualLIfP, err := getRemoteVirtualLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setVirtualLogicalInterfaceProfileAttributes(l3extVirtualLIfP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVirtualLogicalInterfaceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalInterfaceProfile: Beginning Creation")

	aciClient := m.(*client.Client)

	desc := d.Get("description").(string)

	nodeDn := d.Get("node_dn").(string)

	encap := d.Get("encap").(string)

	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	l3extVirtualLIfPAttr := models.VirtualLogicalInterfaceProfileAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extVirtualLIfPAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extVirtualLIfPAttr.Annotation = Annotation.(string)
	} else {
		l3extVirtualLIfPAttr.Annotation = "{}"
	}
	if Autostate, ok := d.GetOk("autostate"); ok {
		l3extVirtualLIfPAttr.Autostate = Autostate.(string)
	}
	if Encap, ok := d.GetOk("encap"); ok {
		l3extVirtualLIfPAttr.Encap = Encap.(string)
	}
	if EncapScope, ok := d.GetOk("encap_scope"); ok {
		l3extVirtualLIfPAttr.EncapScope = EncapScope.(string)
	}
	if IfInstT, ok := d.GetOk("if_inst_t"); ok {
		l3extVirtualLIfPAttr.IfInstT = IfInstT.(string)
	}
	if Ipv6Dad, ok := d.GetOk("ipv6_dad"); ok {
		l3extVirtualLIfPAttr.Ipv6Dad = Ipv6Dad.(string)
	}
	if LlAddr, ok := d.GetOk("ll_addr"); ok {
		l3extVirtualLIfPAttr.LlAddr = LlAddr.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		l3extVirtualLIfPAttr.Mac = Mac.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		l3extVirtualLIfPAttr.Mode = Mode.(string)
	}
	if Mtu, ok := d.GetOk("mtu"); ok {
		l3extVirtualLIfPAttr.Mtu = Mtu.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extVirtualLIfPAttr.TargetDscp = TargetDscp.(string)
	}
	l3extVirtualLIfP := models.NewVirtualLogicalInterfaceProfile(fmt.Sprintf("vlifp-[%s]-[%s]", nodeDn, encap), LogicalInterfaceProfileDn, desc, l3extVirtualLIfPAttr)

	err := aciClient.Save(l3extVirtualLIfP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTol3extRsDynPathAtt, ok := d.GetOk("relation_l3ext_rs_dyn_path_att"); ok {
		relationParamList := toStringList(relationTol3extRsDynPathAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTol3extRsDynPathAtt, ok := d.GetOk("relation_l3ext_rs_dyn_path_att"); ok {
		relationParamList := toStringList(relationTol3extRsDynPathAtt.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationl3extRsDynPathAttFromLogicalInterfaceProfile(l3extVirtualLIfP.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}

	d.SetId(l3extVirtualLIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciVirtualLogicalInterfaceProfileRead(ctx, d, m)
}

func resourceAciVirtualLogicalInterfaceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalInterfaceProfile: Beginning Update")

	aciClient := m.(*client.Client)

	desc := d.Get("description").(string)

	nodeDn := d.Get("node_dn").(string)

	encap := d.Get("encap").(string)

	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	l3extVirtualLIfPAttr := models.VirtualLogicalInterfaceProfileAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extVirtualLIfPAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extVirtualLIfPAttr.Annotation = Annotation.(string)
	} else {
		l3extVirtualLIfPAttr.Annotation = "{}"
	}
	if Autostate, ok := d.GetOk("autostate"); ok {
		l3extVirtualLIfPAttr.Autostate = Autostate.(string)
	}
	if Encap, ok := d.GetOk("encap"); ok {
		l3extVirtualLIfPAttr.Encap = Encap.(string)
	}
	if EncapScope, ok := d.GetOk("encap_scope"); ok {
		l3extVirtualLIfPAttr.EncapScope = EncapScope.(string)
	}
	if IfInstT, ok := d.GetOk("if_inst_t"); ok {
		l3extVirtualLIfPAttr.IfInstT = IfInstT.(string)
	}
	if Ipv6Dad, ok := d.GetOk("ipv6_dad"); ok {
		l3extVirtualLIfPAttr.Ipv6Dad = Ipv6Dad.(string)
	}
	if LlAddr, ok := d.GetOk("ll_addr"); ok {
		l3extVirtualLIfPAttr.LlAddr = LlAddr.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		l3extVirtualLIfPAttr.Mac = Mac.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		l3extVirtualLIfPAttr.Mode = Mode.(string)
	}
	if Mtu, ok := d.GetOk("mtu"); ok {
		l3extVirtualLIfPAttr.Mtu = Mtu.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l3extVirtualLIfPAttr.TargetDscp = TargetDscp.(string)
	}
	l3extVirtualLIfP := models.NewVirtualLogicalInterfaceProfile(fmt.Sprintf("vlifp-[%s]-[%s]", nodeDn, encap), LogicalInterfaceProfileDn, desc, l3extVirtualLIfPAttr)

	l3extVirtualLIfP.Status = "modified"

	err := aciClient.Save(l3extVirtualLIfP)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_l3ext_rs_dyn_path_att") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_dyn_path_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_l3ext_rs_dyn_path_att") {
		oldRel, newRel := d.GetChange("relation_l3ext_rs_dyn_path_att")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationl3extRsDynPathAttFromLogicalInterfaceProfile(l3extVirtualLIfP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationl3extRsDynPathAttFromLogicalInterfaceProfile(l3extVirtualLIfP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}

	d.SetId(l3extVirtualLIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciVirtualLogicalInterfaceProfileRead(ctx, d, m)

}

func resourceAciVirtualLogicalInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extVirtualLIfP, err := getRemoteVirtualLogicalInterfaceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setVirtualLogicalInterfaceProfileAttributes(l3extVirtualLIfP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	l3extRsDynPathAttData, err := aciClient.ReadRelationl3extRsDynPathAttFromLogicalInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l3extRsDynPathAtt %v", err)
		d.Set("relation_l3ext_rs_dyn_path_att", make([]string, 0, 1))
	} else {
		d.Set("relation_l3ext_rs_dyn_path_att", toStringList(l3extRsDynPathAttData.(*schema.Set).List()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciVirtualLogicalInterfaceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extVirtualLIfP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

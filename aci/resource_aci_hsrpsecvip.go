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

func resourceAciL3outHSRPSecondaryVIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL3outHSRPSecondaryVIPCreate,
		UpdateContext: resourceAciL3outHSRPSecondaryVIPUpdate,
		ReadContext:   resourceAciL3outHSRPSecondaryVIPRead,
		DeleteContext: resourceAciL3outHSRPSecondaryVIPDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outHSRPSecondaryVIPImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l3out_hsrp_interface_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"config_issues": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GroupMac-Conflicts-Other-Group",
					"GroupName-Conflicts-Other-Group",
					"GroupVIP-Conflicts-Other-Group",
					"Multiple-Version-On-Interface",
					"Secondary-vip-conflicts-if-ip",
					"Secondary-vip-subnet-mismatch",
					"group-vip-conflicts-if-ip",
					"group-vip-subnet-mismatch",
					"none",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteL3outHSRPSecondaryVIP(client *client.Client, dn string) (*models.L3outHSRPSecondaryVIP, error) {
	hsrpSecVipCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	hsrpSecVip := models.L3outHSRPSecondaryVIPFromContainer(hsrpSecVipCont)

	if hsrpSecVip.DistinguishedName == "" {
		return nil, fmt.Errorf("L3outHSRPSecondaryVIP %s not found", hsrpSecVip.DistinguishedName)
	}

	return hsrpSecVip, nil
}

func setL3outHSRPSecondaryVIPAttributes(hsrpSecVip *models.L3outHSRPSecondaryVIP, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(hsrpSecVip.DistinguishedName)
	d.Set("description", hsrpSecVip.Description)
	dn := d.Id()
	if dn != hsrpSecVip.DistinguishedName {
		d.Set("l3out_hsrp_interface_group_dn", "")
	}
	hsrpSecVipMap, err := hsrpSecVip.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("l3out_hsrp_interface_group_dn", GetParentDn(dn, fmt.Sprintf("/hsrpSecVip-[%s]", hsrpSecVipMap["ip"])))
	d.Set("ip", hsrpSecVipMap["ip"])
	d.Set("annotation", hsrpSecVipMap["annotation"])
	if hsrpSecVipMap["configIssues"] == "" {
		d.Set("config_issues", "none")
	} else {
		d.Set("config_issues", hsrpSecVipMap["configIssues"])
	}
	d.Set("ip", hsrpSecVipMap["ip"])
	d.Set("name_alias", hsrpSecVipMap["nameAlias"])
	return d, nil
}

func resourceAciL3outHSRPSecondaryVIPImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	hsrpSecVip, err := getRemoteL3outHSRPSecondaryVIP(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL3outHSRPSecondaryVIPAttributes(hsrpSecVip, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outHSRPSecondaryVIPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outHSRPSecondaryVIP: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	HSRPGroupProfileDn := d.Get("l3out_hsrp_interface_group_dn").(string)

	hsrpSecVipAttr := models.L3outHSRPSecondaryVIPAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpSecVipAttr.Annotation = Annotation.(string)
	} else {
		hsrpSecVipAttr.Annotation = "{}"
	}
	if ConfigIssues, ok := d.GetOk("config_issues"); ok {
		hsrpSecVipAttr.ConfigIssues = ConfigIssues.(string)
	}
	if Ip, ok := d.GetOk("ip"); ok {
		hsrpSecVipAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpSecVipAttr.NameAlias = NameAlias.(string)
	}
	hsrpSecVip := models.NewL3outHSRPSecondaryVIP(fmt.Sprintf("hsrpSecVip-[%s]", ip), HSRPGroupProfileDn, desc, hsrpSecVipAttr)

	err := aciClient.Save(hsrpSecVip)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hsrpSecVip.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outHSRPSecondaryVIPRead(ctx, d, m)
}

func resourceAciL3outHSRPSecondaryVIPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outHSRPSecondaryVIP: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	HSRPGroupProfileDn := d.Get("l3out_hsrp_interface_group_dn").(string)

	hsrpSecVipAttr := models.L3outHSRPSecondaryVIPAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpSecVipAttr.Annotation = Annotation.(string)
	} else {
		hsrpSecVipAttr.Annotation = "{}"
	}
	if ConfigIssues, ok := d.GetOk("config_issues"); ok {
		hsrpSecVipAttr.ConfigIssues = ConfigIssues.(string)
	}
	if Ip, ok := d.GetOk("ip"); ok {
		hsrpSecVipAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpSecVipAttr.NameAlias = NameAlias.(string)
	}
	hsrpSecVip := models.NewL3outHSRPSecondaryVIP(fmt.Sprintf("hsrpSecVip-[%s]", ip), HSRPGroupProfileDn, desc, hsrpSecVipAttr)

	hsrpSecVip.Status = "modified"

	err := aciClient.Save(hsrpSecVip)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hsrpSecVip.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outHSRPSecondaryVIPRead(ctx, d, m)

}

func resourceAciL3outHSRPSecondaryVIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	hsrpSecVip, err := getRemoteL3outHSRPSecondaryVIP(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setL3outHSRPSecondaryVIPAttributes(hsrpSecVip, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outHSRPSecondaryVIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "hsrpSecVip")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

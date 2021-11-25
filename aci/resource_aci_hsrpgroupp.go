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

func resourceAciHSRPGroupProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciHSRPGroupProfileCreate,
		UpdateContext: resourceAciHSRPGroupProfileUpdate,
		ReadContext:   resourceAciHSRPGroupProfileRead,
		DeleteContext: resourceAciHSRPGroupProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciHSRPGroupProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l3out_hsrp_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
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

			"group_af": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ipv4",
					"ipv6",
				}, false),
			},

			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"group_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip_obtain_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"admin",
					"auto",
					"learn",
				}, false),
			},

			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_hsrp_rs_group_pol": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "uni/tn-common/hsrpGroupPol-default",
				Optional: true,
			},
		}),
	}
}

func getRemoteHSRPGroupProfile(client *client.Client, dn string) (*models.HSRPGroupProfile, error) {
	hsrpGroupPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	hsrpGroupP := models.HSRPGroupProfileFromContainer(hsrpGroupPCont)

	if hsrpGroupP.DistinguishedName == "" {
		return nil, fmt.Errorf("HSRPGroupProfile %s not found", hsrpGroupP.DistinguishedName)
	}

	return hsrpGroupP, nil
}

func setHSRPGroupProfileAttributes(hsrpGroupP *models.HSRPGroupProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	d.SetId(hsrpGroupP.DistinguishedName)
	d.Set("description", hsrpGroupP.Description)
	if dn != hsrpGroupP.DistinguishedName {
		d.Set("l3out_hsrp_interface_profile_dn", "")
	}
	hsrpGroupPMap, err := hsrpGroupP.ToMap()

	if err != nil {
		return d, err
	}
	d.Set("l3out_hsrp_interface_profile_dn", GetParentDn(hsrpGroupP.DistinguishedName, fmt.Sprintf("/hsrpGroupP-%s", hsrpGroupPMap["name"])))
	d.Set("name", hsrpGroupPMap["name"])
	d.Set("annotation", hsrpGroupPMap["annotation"])
	if hsrpGroupPMap["configIssues"] == "" {
		d.Set("config_issues", "none")
	} else {
		d.Set("config_issues", hsrpGroupPMap["configIssues"])
	}
	d.Set("group_af", hsrpGroupPMap["groupAf"])
	d.Set("group_id", hsrpGroupPMap["groupId"])
	d.Set("group_name", hsrpGroupPMap["groupName"])
	d.Set("ip", hsrpGroupPMap["ip"])
	d.Set("ip_obtain_mode", hsrpGroupPMap["ipObtainMode"])
	d.Set("mac", hsrpGroupPMap["mac"])
	d.Set("name_alias", hsrpGroupPMap["nameAlias"])

	return d, nil
}

func resourceAciHSRPGroupProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	hsrpGroupP, err := getRemoteHSRPGroupProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setHSRPGroupProfileAttributes(hsrpGroupP, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciHSRPGroupProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] HSRPGroupProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	parentDn := d.Get("l3out_hsrp_interface_profile_dn").(string)

	hsrpGroupPAttr := models.HSRPGroupProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpGroupPAttr.Annotation = Annotation.(string)
	} else {
		hsrpGroupPAttr.Annotation = "{}"
	}
	if ConfigIssues, ok := d.GetOk("config_issues"); ok {
		hsrpGroupPAttr.ConfigIssues = ConfigIssues.(string)
	}
	if GroupAf, ok := d.GetOk("group_af"); ok {
		hsrpGroupPAttr.GroupAf = GroupAf.(string)
	}
	if GroupId, ok := d.GetOk("group_id"); ok {
		hsrpGroupPAttr.GroupId = GroupId.(string)
	}
	if GroupName, ok := d.GetOk("group_name"); ok {
		hsrpGroupPAttr.GroupName = GroupName.(string)
	}
	if Ip, ok := d.GetOk("ip"); ok {
		hsrpGroupPAttr.Ip = Ip.(string)
	}
	if IpObtainMode, ok := d.GetOk("ip_obtain_mode"); ok {
		hsrpGroupPAttr.IpObtainMode = IpObtainMode.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		hsrpGroupPAttr.Mac = Mac.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpGroupPAttr.NameAlias = NameAlias.(string)
	}
	hsrpGroupP := models.NewHSRPGroupProfile(fmt.Sprintf("hsrpGroupP-%s", name), parentDn, desc, hsrpGroupPAttr)

	err := aciClient.Save(hsrpGroupP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTohsrpRsGroupPol, ok := d.GetOk("relation_hsrp_rs_group_pol"); ok {
		relationParam := relationTohsrpRsGroupPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTohsrpRsGroupPol, ok := d.GetOk("relation_hsrp_rs_group_pol"); ok {
		relationParam := relationTohsrpRsGroupPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationhsrpRsGroupPolFromHSRPGroupProfile(hsrpGroupP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(hsrpGroupP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciHSRPGroupProfileRead(ctx, d, m)
}

func resourceAciHSRPGroupProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] HSRPGroupProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	parentDn := d.Get("l3out_hsrp_interface_profile_dn").(string)

	hsrpGroupPAttr := models.HSRPGroupProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpGroupPAttr.Annotation = Annotation.(string)
	} else {
		hsrpGroupPAttr.Annotation = "{}"
	}
	if ConfigIssues, ok := d.GetOk("config_issues"); ok {
		hsrpGroupPAttr.ConfigIssues = ConfigIssues.(string)
	}
	if GroupAf, ok := d.GetOk("group_af"); ok {
		hsrpGroupPAttr.GroupAf = GroupAf.(string)
	}
	if GroupId, ok := d.GetOk("group_id"); ok {
		hsrpGroupPAttr.GroupId = GroupId.(string)
	}
	if GroupName, ok := d.GetOk("group_name"); ok {
		hsrpGroupPAttr.GroupName = GroupName.(string)
	}
	if Ip, ok := d.GetOk("ip"); ok {
		hsrpGroupPAttr.Ip = Ip.(string)
	}
	if IpObtainMode, ok := d.GetOk("ip_obtain_mode"); ok {
		hsrpGroupPAttr.IpObtainMode = IpObtainMode.(string)
	}
	if Mac, ok := d.GetOk("mac"); ok {
		hsrpGroupPAttr.Mac = Mac.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpGroupPAttr.NameAlias = NameAlias.(string)
	}
	hsrpGroupP := models.NewHSRPGroupProfile(fmt.Sprintf("hsrpGroupP-%s", name), parentDn, desc, hsrpGroupPAttr)

	hsrpGroupP.Status = "modified"

	err := aciClient.Save(hsrpGroupP)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_hsrp_rs_group_pol") {
		_, newRelParam := d.GetChange("relation_hsrp_rs_group_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_hsrp_rs_group_pol") {
		_, newRelParam := d.GetChange("relation_hsrp_rs_group_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationhsrpRsGroupPolFromHSRPGroupProfile(hsrpGroupP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(hsrpGroupP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciHSRPGroupProfileRead(ctx, d, m)

}

func resourceAciHSRPGroupProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	hsrpGroupP, err := getRemoteHSRPGroupProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setHSRPGroupProfileAttributes(hsrpGroupP, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	hsrpRsGroupPolData, err := aciClient.ReadRelationhsrpRsGroupPolFromHSRPGroupProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation hsrpRsGroupPol %v", err)
		d.Set("relation_hsrp_rs_group_pol", "")
	} else {
		d.Set("relation_hsrp_rs_group_pol", hsrpRsGroupPolData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciHSRPGroupProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "hsrpGroupP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

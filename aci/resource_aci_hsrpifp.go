package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciL3outHSRPInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciL3outHSRPInterfaceProfileCreate,
		Update: resourceAciL3outHSRPInterfaceProfileUpdate,
		Read:   resourceAciL3outHSRPInterfaceProfileRead,
		Delete: resourceAciL3outHSRPInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outHSRPInterfaceProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"v1",
					"v2",
				}, false),
			},

			"relation_hsrp_rs_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteL3outHSRPInterfaceProfile(client *client.Client, dn string) (*models.L3outHSRPInterfaceProfile, error) {
	hsrpIfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	hsrpIfP := models.L3outHSRPInterfaceProfileFromContainer(hsrpIfPCont)

	if hsrpIfP.DistinguishedName == "" {
		return nil, fmt.Errorf("L3outHSRPInterfaceProfile %s not found", hsrpIfP.DistinguishedName)
	}

	return hsrpIfP, nil
}

func setL3outHSRPInterfaceProfileAttributes(hsrpIfP *models.L3outHSRPInterfaceProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(hsrpIfP.DistinguishedName)
	d.Set("description", hsrpIfP.Description)
	dn := d.Id()
	if dn != hsrpIfP.DistinguishedName {
		d.Set("logical_interface_profile_dn", "")
	}
	hsrpIfPMap, _ := hsrpIfP.ToMap()

	d.Set("annotation", hsrpIfPMap["annotation"])
	d.Set("name_alias", hsrpIfPMap["nameAlias"])
	d.Set("version", hsrpIfPMap["version"])
	return d
}

func resourceAciL3outHSRPInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	hsrpIfP, err := getRemoteL3outHSRPInterfaceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setL3outHSRPInterfaceProfileAttributes(hsrpIfP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outHSRPInterfaceProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outHSRPInterfaceProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	hsrpIfPAttr := models.L3outHSRPInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpIfPAttr.Annotation = Annotation.(string)
	} else {
		hsrpIfPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpIfPAttr.NameAlias = NameAlias.(string)
	}
	if Version, ok := d.GetOk("version"); ok {
		hsrpIfPAttr.Version = Version.(string)
	}
	hsrpIfP := models.NewL3outHSRPInterfaceProfile(fmt.Sprintf("hsrpIfP"), LogicalInterfaceProfileDn, desc, hsrpIfPAttr)

	err := aciClient.Save(hsrpIfP)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationTohsrpRsIfPol, ok := d.GetOk("relation_hsrp_rs_if_pol"); ok {
		relationParam := relationTohsrpRsIfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if relationTohsrpRsIfPol, ok := d.GetOk("relation_hsrp_rs_if_pol"); ok {
		relationParam := GetMOName(relationTohsrpRsIfPol.(string))
		err = aciClient.CreateRelationhsrpRsIfPolFromL3outHSRPInterfaceProfile(hsrpIfP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_hsrp_rs_if_pol")
		d.Partial(false)

	}

	d.SetId(hsrpIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outHSRPInterfaceProfileRead(d, m)
}

func resourceAciL3outHSRPInterfaceProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outHSRPInterfaceProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	hsrpIfPAttr := models.L3outHSRPInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		hsrpIfPAttr.Annotation = Annotation.(string)
	} else {
		hsrpIfPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		hsrpIfPAttr.NameAlias = NameAlias.(string)
	}
	if Version, ok := d.GetOk("version"); ok {
		hsrpIfPAttr.Version = Version.(string)
	}
	hsrpIfP := models.NewL3outHSRPInterfaceProfile(fmt.Sprintf("hsrpIfP"), LogicalInterfaceProfileDn, desc, hsrpIfPAttr)

	hsrpIfP.Status = "modified"

	err := aciClient.Save(hsrpIfP)

	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_hsrp_rs_if_pol") {
		_, newRelParam := d.GetChange("relation_hsrp_rs_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if d.HasChange("relation_hsrp_rs_if_pol") {
		_, newRelParam := d.GetChange("relation_hsrp_rs_if_pol")
		err = aciClient.CreateRelationhsrpRsIfPolFromL3outHSRPInterfaceProfile(hsrpIfP.DistinguishedName, GetMOName(newRelParam.(string)))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_hsrp_rs_if_pol")
		d.Partial(false)

	}

	d.SetId(hsrpIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outHSRPInterfaceProfileRead(d, m)

}

func resourceAciL3outHSRPInterfaceProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	hsrpIfP, err := getRemoteL3outHSRPInterfaceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setL3outHSRPInterfaceProfileAttributes(hsrpIfP, d)

	hsrpRsIfPolData, err := aciClient.ReadRelationhsrpRsIfPolFromL3outHSRPInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation hsrpRsIfPol %v", err)
		d.Set("relation_hsrp_rs_if_pol", "")

	} else {
		if _, ok := d.GetOk("relation_hsrp_rs_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_hsrp_rs_if_pol").(string))
			if tfName != hsrpRsIfPolData {
				d.Set("relation_hsrp_rs_if_pol", "")
			}
		}

	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outHSRPInterfaceProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "hsrpIfP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

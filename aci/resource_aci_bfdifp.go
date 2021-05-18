package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciBFDInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciBFDInterfaceProfileCreate,
		Update: resourceAciBFDInterfaceProfileUpdate,
		Read:   resourceAciBFDInterfaceProfileRead,
		Delete: resourceAciBFDInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBFDInterfaceProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"key_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"interface_profile_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"sha1",
				}, false),
			},

			"relation_bfd_rs_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}

func getRemoteBFDInterfaceProfile(client *client.Client, dn string) (*models.BFDInterfaceProfile, error) {
	bfdIfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bfdIfP := models.BFDInterfaceProfileFromContainer(bfdIfPCont)

	if bfdIfP.DistinguishedName == "" {
		return nil, fmt.Errorf("InterfaceProfile %s not found", bfdIfP.DistinguishedName)
	}

	return bfdIfP, nil
}

func setBFDInterfaceProfileAttributes(bfdIfP *models.BFDInterfaceProfile, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()

	d.SetId(bfdIfP.DistinguishedName)
	d.Set("description", bfdIfP.Description)
	if dn != bfdIfP.DistinguishedName {
		d.Set("logical_interface_profile_dn", "")
	}
	bfdIfPMap, _ := bfdIfP.ToMap()

	d.Set("annotation", bfdIfPMap["annotation"])
	d.Set("key_id", bfdIfPMap["keyId"])
	d.Set("name_alias", bfdIfPMap["nameAlias"])
	d.Set("interface_profile_type", bfdIfPMap["type"])

	return d
}

func resourceAciBFDInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bfdIfP, err := getRemoteBFDInterfaceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setBFDInterfaceProfileAttributes(bfdIfP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBFDInterfaceProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] InterfaceProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	bfdIfPAttr := models.BFDInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bfdIfPAttr.Annotation = Annotation.(string)
	} else {
		bfdIfPAttr.Annotation = "{}"
	}
	if Key, ok := d.GetOk("key"); ok {
		bfdIfPAttr.Key = Key.(string)
	}
	if KeyId, ok := d.GetOk("key_id"); ok {
		bfdIfPAttr.KeyId = KeyId.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bfdIfPAttr.NameAlias = NameAlias.(string)
	}
	if InterfaceProfile_type, ok := d.GetOk("interface_profile_type"); ok {
		bfdIfPAttr.InterfaceProfileType = InterfaceProfile_type.(string)
	}
	bfdIfP := models.NewBFDInterfaceProfile(fmt.Sprintf("bfdIfP"), LogicalInterfaceProfileDn, desc, bfdIfPAttr)

	err := aciClient.Save(bfdIfP)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationTobfdRsIfPol, ok := d.GetOk("relation_bfd_rs_if_pol"); ok {
		relationParam := relationTobfdRsIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if relationTobfdRsIfPol, ok := d.GetOk("relation_bfd_rs_if_pol"); ok {
		relationParam := relationTobfdRsIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationbfdRsIfPolFromInterfaceProfile(bfdIfP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_bfd_rs_if_pol")
		d.Partial(false)

	}

	d.SetId(bfdIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBFDInterfaceProfileRead(d, m)
}

func resourceAciBFDInterfaceProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] InterfaceProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	bfdIfPAttr := models.BFDInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bfdIfPAttr.Annotation = Annotation.(string)
	} else {
		bfdIfPAttr.Annotation = "{}"
	}
	if Key, ok := d.GetOk("key"); ok {
		bfdIfPAttr.Key = Key.(string)
	}
	if KeyId, ok := d.GetOk("key_id"); ok {
		bfdIfPAttr.KeyId = KeyId.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bfdIfPAttr.NameAlias = NameAlias.(string)
	}
	if InterfaceProfile_type, ok := d.GetOk("interface_profile_type"); ok {
		bfdIfPAttr.InterfaceProfileType = InterfaceProfile_type.(string)
	}
	bfdIfP := models.NewBFDInterfaceProfile(fmt.Sprintf("bfdIfP"), LogicalInterfaceProfileDn, desc, bfdIfPAttr)

	bfdIfP.Status = "modified"

	err := aciClient.Save(bfdIfP)

	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_bfd_rs_if_pol") {
		_, newRelParam := d.GetChange("relation_bfd_rs_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if d.HasChange("relation_bfd_rs_if_pol") {
		_, newRelParam := d.GetChange("relation_bfd_rs_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationbfdRsIfPolFromInterfaceProfile(bfdIfP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_bfd_rs_if_pol")
		d.Partial(false)

	}

	d.SetId(bfdIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBFDInterfaceProfileRead(d, m)

}

func resourceAciBFDInterfaceProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bfdIfP, err := getRemoteBFDInterfaceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setBFDInterfaceProfileAttributes(bfdIfP, d)

	bfdRsIfPolData, err := aciClient.ReadRelationbfdRsIfPolFromInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation bfdRsIfPol %v", err)
		d.Set("relation_bfd_rs_if_pol", "")

	} else {
		if _, ok := d.GetOk("relation_bfd_rs_if_pol"); ok {
			tfName := GetMOName(d.Get("relation_bfd_rs_if_pol").(string))
			if tfName != bfdRsIfPolData {
				d.Set("relation_bfd_rs_if_pol", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBFDInterfaceProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bfdIfP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

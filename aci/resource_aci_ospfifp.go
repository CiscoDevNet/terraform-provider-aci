package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciOSPFInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciOSPFInterfaceProfileCreate,
		Update: resourceAciOSPFInterfaceProfileUpdate,
		Read:   resourceAciOSPFInterfaceProfileRead,
		Delete: resourceAciOSPFInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciInterfaceProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"auth_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"auth_key_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"auth_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_ospf_rs_if_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}

func getRemoteOSPFInterfaceProfile(client *client.Client, dn string) (*models.OSPFInterfaceProfile, error) {
	ospfIfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	ospfIfP := models.OSPFInterfaceProfileFromContainer(ospfIfPCont)

	if ospfIfP.DistinguishedName == "" {
		return nil, fmt.Errorf("InterfaceProfile %s not found", ospfIfP.DistinguishedName)
	}

	return ospfIfP, nil
}

func setOSPFInterfaceProfileAttributes(ospfIfP *models.OSPFInterfaceProfile, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()

	d.SetId(ospfIfP.DistinguishedName)
	d.Set("description", ospfIfP.Description)
	if dn != ospfIfP.DistinguishedName {
		d.Set("logical_interface_profile_dn", "")
	}
	ospfIfPMap, _ := ospfIfP.ToMap()

	d.Set("annotation", ospfIfPMap["annotation"])
	d.Set("auth_key", ospfIfPMap["authKey"])
	d.Set("auth_key_id", ospfIfPMap["authKeyId"])
	d.Set("auth_type", ospfIfPMap["authType"])
	d.Set("name_alias", ospfIfPMap["nameAlias"])
	return d
}

func resourceAciOSPFInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	ospfIfP, err := getRemoteOSPFInterfaceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setOSPFInterfaceProfileAttributes(ospfIfP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciOSPFInterfaceProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] InterfaceProfile: Beginning Creation")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	ospfIfPAttr := models.OSPFInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ospfIfPAttr.Annotation = Annotation.(string)
	} else {
		ospfIfPAttr.Annotation = "{}"
	}
	if AuthKey, ok := d.GetOk("auth_key"); ok {
		ospfIfPAttr.AuthKey = AuthKey.(string)
	}
	if AuthKeyId, ok := d.GetOk("auth_key_id"); ok {
		ospfIfPAttr.AuthKeyId = AuthKeyId.(string)
	}
	if AuthType, ok := d.GetOk("auth_type"); ok {
		ospfIfPAttr.AuthType = AuthType.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ospfIfPAttr.NameAlias = NameAlias.(string)
	}
	ospfIfP := models.NewOSPFInterfaceProfile(fmt.Sprintf("ospfIfP"), LogicalInterfaceProfileDn, desc, ospfIfPAttr)

	err := aciClient.Save(ospfIfP)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	if relationToospfRsIfPol, ok := d.GetOk("relation_ospf_rs_if_pol"); ok {
		relationParam := relationToospfRsIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationospfRsIfPolFromInterfaceProfile(ospfIfP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_ospf_rs_if_pol")
		d.Partial(false)

	}

	d.SetId(ospfIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciInterfaceProfileRead(d, m)
}

func resourceAciOSPFInterfaceProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] InterfaceProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	ospfIfPAttr := models.OSPFInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ospfIfPAttr.Annotation = Annotation.(string)
	} else {
		ospfIfPAttr.Annotation = "{}"
	}
	if AuthKey, ok := d.GetOk("auth_key"); ok {
		ospfIfPAttr.AuthKey = AuthKey.(string)
	}
	if AuthKeyId, ok := d.GetOk("auth_key_id"); ok {
		ospfIfPAttr.AuthKeyId = AuthKeyId.(string)
	}
	if AuthType, ok := d.GetOk("auth_type"); ok {
		ospfIfPAttr.AuthType = AuthType.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ospfIfPAttr.NameAlias = NameAlias.(string)
	}
	ospfIfP := models.NewOSPFInterfaceProfile(fmt.Sprintf("ospfIfP"), LogicalInterfaceProfileDn, desc, ospfIfPAttr)

	ospfIfP.Status = "modified"

	err := aciClient.Save(ospfIfP)

	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	if d.HasChange("relation_ospf_rs_if_pol") {
		_, newRelParam := d.GetChange("relation_ospf_rs_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationospfRsIfPolFromInterfaceProfile(ospfIfP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_ospf_rs_if_pol")
		d.Partial(false)

	}

	d.SetId(ospfIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciInterfaceProfileRead(d, m)

}

func resourceAciOSPFInterfaceProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	ospfIfP, err := getRemoteOSPFInterfaceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setOSPFInterfaceProfileAttributes(ospfIfP, d)

	ospfRsIfPolData, err := aciClient.ReadRelationospfRsIfPolFromInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation ospfRsIfPol %v", err)
		d.Set("relation_ospf_rs_if_pol", "")

	} else {
		if tf, ok := d.GetOk("relation_ospf_rs_if_pol"); ok {
			tfName := GetMOName(tf.(string))
			if tfName != ospfRsIfPolData {
				d.Set("relation_ospf_rs_if_pol", "")
			}
		} else {
			d.Set("relation_ospf_rs_if_pol", "")
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciOSPFInterfaceProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "ospfIfP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

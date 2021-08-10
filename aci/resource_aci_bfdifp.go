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

func resourceAciBFDInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBFDInterfaceProfileCreate,
		UpdateContext: resourceAciBFDInterfaceProfileUpdate,
		ReadContext:   resourceAciBFDInterfaceProfileRead,
		DeleteContext: resourceAciBFDInterfaceProfileDelete,

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
				Type:     schema.TypeString,
				Default:  "uni/tn-common/bfdIfPol-default",
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

func setBFDInterfaceProfileAttributes(bfdIfP *models.BFDInterfaceProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	d.SetId(bfdIfP.DistinguishedName)
	d.Set("description", bfdIfP.Description)
	if dn != bfdIfP.DistinguishedName {
		d.Set("logical_interface_profile_dn", "")
	}
	bfdIfPMap, err := bfdIfP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("logical_interface_profile_dn", GetParentDn(dn, fmt.Sprintf("/bfdIfP")))
	d.Set("annotation", bfdIfPMap["annotation"])
	d.Set("key_id", bfdIfPMap["keyId"])
	d.Set("name_alias", bfdIfPMap["nameAlias"])
	d.Set("interface_profile_type", bfdIfPMap["type"])

	return d, nil
}

func resourceAciBFDInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bfdIfP, err := getRemoteBFDInterfaceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setBFDInterfaceProfileAttributes(bfdIfP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBFDInterfaceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTobfdRsIfPol, ok := d.GetOk("relation_bfd_rs_if_pol"); ok {
		relationParam := relationTobfdRsIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTobfdRsIfPol, ok := d.GetOk("relation_bfd_rs_if_pol"); ok {
		relationParam := relationTobfdRsIfPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationbfdRsIfPolFromInterfaceProfile(bfdIfP.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(bfdIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBFDInterfaceProfileRead(ctx, d, m)
}

func resourceAciBFDInterfaceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_bfd_rs_if_pol") {
		_, newRelParam := d.GetChange("relation_bfd_rs_if_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_bfd_rs_if_pol") {
		_, newRelParam := d.GetChange("relation_bfd_rs_if_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationbfdRsIfPolFromInterfaceProfile(bfdIfP.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(bfdIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBFDInterfaceProfileRead(ctx, d, m)

}

func resourceAciBFDInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bfdIfP, err := getRemoteBFDInterfaceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setBFDInterfaceProfileAttributes(bfdIfP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	bfdRsIfPolData, err := aciClient.ReadRelationbfdRsIfPolFromInterfaceProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation bfdRsIfPol %v", err)
		d.Set("relation_bfd_rs_if_pol", "")

	} else {
		d.Set("relation_bfd_rs_if_pol", bfdRsIfPolData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBFDInterfaceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bfdIfP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

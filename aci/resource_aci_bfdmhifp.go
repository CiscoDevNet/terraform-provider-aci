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

func resourceAciBfdMultihopInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBfdMultihopInterfaceProfileCreate,
		UpdateContext: resourceAciBfdMultihopInterfaceProfileUpdate,
		ReadContext:   resourceAciBfdMultihopInterfaceProfileRead,
		DeleteContext: resourceAciBfdMultihopInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBfdMultihopInterfaceProfileImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"key_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
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
			"relation_bfd_rs_mh_if_pol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to bfd:MhIfPol",
			}})),
	}
}

func getRemoteAciBfdMultihopInterfaceProfile(client *client.Client, dn string) (*models.AciBfdMultihopInterfaceProfile, error) {
	bfdMhIfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	bfdMhIfP := models.AciBfdMultihopInterfaceProfileFromContainer(bfdMhIfPCont)
	if bfdMhIfP.DistinguishedName == "" {
		return nil, fmt.Errorf("Aci BFD Multihop Interface Profile %s not found", dn)
	}
	return bfdMhIfP, nil
}

func setAciBfdMultihopInterfaceProfileAttributes(bfdMhIfP *models.AciBfdMultihopInterfaceProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(bfdMhIfP.DistinguishedName)
	d.Set("description", bfdMhIfP.Description)
	if dn != bfdMhIfP.DistinguishedName {
		d.Set("logical_interface_profile_dn", "")
	}
	bfdMhIfPMap, err := bfdMhIfP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("logical_interface_profile_dn", GetParentDn(dn, fmt.Sprintf("/%s", models.RnbfdMhIfP)))
	d.Set("annotation", bfdMhIfPMap["annotation"])
	key := d.Get("key")
	if key != nil {
		keyValue := key.(string)
		if keyValue != "" {
			d.Set("key", keyValue)
		}
	}
	d.Set("key_id", bfdMhIfPMap["keyId"])
	d.Set("name", bfdMhIfPMap["name"])
	d.Set("interface_profile_type", bfdMhIfPMap["type"])
	d.Set("name_alias", bfdMhIfPMap["nameAlias"])
	return d, nil
}

func getAndSetRelationAciBfdMultihopInterfacePolicy(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	bfdRsMhIfPolData, err := client.ReadRelationbfdRsMhIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation bfdRsMhIfPol %v", err)
		d.Set("relation_bfd_rs_mh_if_pol", "")
		return nil, err
	} else {
		d.Set("relation_bfd_rs_mh_if_pol", bfdRsMhIfPolData.(string))
		log.Printf("[DEBUG] Reading relation bfdRsMhIfPol: %s finished successfully", bfdRsMhIfPolData.(string))
	}
	return d, nil
}

func resourceAciBfdMultihopInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Aci BFD Multihop Interface Profile %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	bfdMhIfP, err := getRemoteAciBfdMultihopInterfaceProfile(aciClient, d.Id())
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setAciBfdMultihopInterfaceProfileAttributes(bfdMhIfP, d)
	if err != nil {
		return nil, err
	}

	getAndSetRelationAciBfdMultihopInterfacePolicy(aciClient, d.Id(), d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBfdMultihopInterfaceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Aci BFD Multihop Interface Profile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	nameAlias := d.Get("name_alias").(string)
	logicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	bfdMhIfPAttr := models.AciBfdMultihopInterfaceProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		bfdMhIfPAttr.Annotation = Annotation.(string)
	} else {
		bfdMhIfPAttr.Annotation = "{}"
	}

	if Key, ok := d.GetOk("key"); ok {
		bfdMhIfPAttr.Key = Key.(string)
	}

	if KeyId, ok := d.GetOk("key_id"); ok {
		bfdMhIfPAttr.KeyId = KeyId.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		bfdMhIfPAttr.Name = Name.(string)
	}

	if InterfaceProfile_type, ok := d.GetOk("interface_profile_type"); ok {
		bfdMhIfPAttr.InterfaceProfile_type = InterfaceProfile_type.(string)
	}
	bfdMhIfP := models.NewAciBfdMultihopInterfaceProfile(fmt.Sprintf(models.RnbfdMhIfP), logicalInterfaceProfileDn, desc, nameAlias, bfdMhIfPAttr)

	err := aciClient.Save(bfdMhIfP)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTobfdRsMhIfPol, ok := d.GetOk("relation_bfd_rs_mh_if_pol"); ok {
		relationParam := relationTobfdRsMhIfPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTobfdRsMhIfPol, ok := d.GetOk("relation_bfd_rs_mh_if_pol"); ok {
		relationParam := relationTobfdRsMhIfPol.(string)
		err = aciClient.CreateRelationbfdRsMhIfPol(bfdMhIfP.DistinguishedName, bfdMhIfPAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(bfdMhIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciBfdMultihopInterfaceProfileRead(ctx, d, m)
}

func resourceAciBfdMultihopInterfaceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Aci BFD Multihop Interface Profile: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	nameAlias := d.Get("name_alias").(string)
	logicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	bfdMhIfPAttr := models.AciBfdMultihopInterfaceProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		bfdMhIfPAttr.Annotation = Annotation.(string)
	} else {
		bfdMhIfPAttr.Annotation = "{}"
	}

	if Key, ok := d.GetOk("key"); ok {
		bfdMhIfPAttr.Key = Key.(string)
	}

	if KeyId, ok := d.GetOk("key_id"); ok {
		bfdMhIfPAttr.KeyId = KeyId.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		bfdMhIfPAttr.Name = Name.(string)
	}

	if InterfaceProfile_type, ok := d.GetOk("interface_profile_type"); ok {
		bfdMhIfPAttr.InterfaceProfile_type = InterfaceProfile_type.(string)
	}
	bfdMhIfP := models.NewAciBfdMultihopInterfaceProfile(models.RnbfdMhIfP, logicalInterfaceProfileDn, desc, nameAlias, bfdMhIfPAttr)

	bfdMhIfP.Status = "modified"

	err := aciClient.Save(bfdMhIfP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_bfd_rs_mh_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_bfd_rs_mh_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_bfd_rs_mh_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_bfd_rs_mh_if_pol")
		err = aciClient.DeleteRelationbfdRsMhIfPol(bfdMhIfP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationbfdRsMhIfPol(bfdMhIfP.DistinguishedName, bfdMhIfPAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(bfdMhIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciBfdMultihopInterfaceProfileRead(ctx, d, m)
}

func resourceAciBfdMultihopInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Aci BFD Multihop Interface Profile %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)

	bfdMhIfP, err := getRemoteAciBfdMultihopInterfaceProfile(aciClient, d.Id())
	if err != nil {
		return errorForObjectNotFound(err, d.Id(), d)
	}

	_, err = setAciBfdMultihopInterfaceProfileAttributes(bfdMhIfP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	getAndSetRelationAciBfdMultihopInterfacePolicy(aciClient, d.Id(), d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciBfdMultihopInterfaceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Aci BFD Multihop Interface Profile %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)

	err := aciClient.DeleteByDn(d.Id(), models.RnbfdMhIfP)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

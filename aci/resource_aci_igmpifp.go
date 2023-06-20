package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciIGMPInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciIGMPInterfaceProfileCreate,
		UpdateContext: resourceAciIGMPInterfaceProfileUpdate,
		ReadContext:   resourceAciIGMPInterfaceProfileRead,
		DeleteContext: resourceAciIGMPInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciIGMPInterfaceProfileImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"logical_interface_profile_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_igmp_rs_if_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to igmp:IfPol",
			}})),
	}
}

func getRemoteIGMPInterfaceProfile(client *client.Client, dn string) (*models.IGMPInterfaceProfile, error) {
	igmpIfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	igmpIfP := models.IGMPInterfaceProfileFromContainer(igmpIfPCont)
	if igmpIfP.DistinguishedName == "" {
		return nil, fmt.Errorf("IGMP Interface Profile %s not found", dn)
	}
	return igmpIfP, nil
}

func setIGMPInterfaceProfileAttributes(igmpIfP *models.IGMPInterfaceProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(igmpIfP.DistinguishedName)
	d.Set("description", igmpIfP.Description)
	igmpIfPMap, err := igmpIfP.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != igmpIfP.DistinguishedName {
		d.Set("logical_interface_profile_dn", "")
	} else {
		d.Set("logical_interface_profile_dn", GetParentDn(igmpIfP.DistinguishedName, fmt.Sprintf("/"+models.RnIgmpIfP)))
	}
	d.Set("name", igmpIfPMap["name"])
	return d, nil
}

func getAndSetIGMPInterfaceProfileRelationalAttributes(client *client.Client, dn string, d *schema.ResourceData) {

	log.Printf("[DEBUG] igmpRsIfPol: Beginning Read")

	igmpRsIfPolData, err := client.ReadRelationigmpRsIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation igmpRsIfPol %v", err)
		d.Set("relation_igmp_rs_if_pol", "")
	} else {
		d.Set("relation_igmp_rs_if_pol", igmpRsIfPolData.(string))
		log.Printf("[DEBUG]: igmpRsIfPol: Reading finished successfully")
	}
}

func resourceAciIGMPInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	igmpIfP, err := getRemoteIGMPInterfaceProfile(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setIGMPInterfaceProfileAttributes(igmpIfP, d)
	if err != nil {
		return nil, err
	}

	// Get and Set Relational Attributes
	getAndSetIGMPInterfaceProfileRelationalAttributes(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciIGMPInterfaceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InterfaceProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	igmpIfPAttr := models.IGMPInterfaceProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		igmpIfPAttr.Annotation = Annotation.(string)
	} else {
		igmpIfPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		igmpIfPAttr.Name = Name.(string)
	}
	igmpIfP := models.NewIGMPInterfaceProfile(fmt.Sprintf(models.RnIgmpIfP), LogicalInterfaceProfileDn, desc, igmpIfPAttr)

	err := aciClient.Save(igmpIfP)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationToigmpRsIfPol, ok := d.GetOk("relation_igmp_rs_if_pol"); ok {
		relationParam := relationToigmpRsIfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToigmpRsIfPol, ok := d.GetOk("relation_igmp_rs_if_pol"); ok {
		relationParam := relationToigmpRsIfPol.(string)
		err = aciClient.CreateRelationigmpRsIfPol(igmpIfP.DistinguishedName, igmpIfPAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(igmpIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciIGMPInterfaceProfileRead(ctx, d, m)
}
func resourceAciIGMPInterfaceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Interface Profile: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	igmpIfPAttr := models.IGMPInterfaceProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		igmpIfPAttr.Annotation = Annotation.(string)
	} else {
		igmpIfPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		igmpIfPAttr.Name = Name.(string)
	}
	igmpIfP := models.NewIGMPInterfaceProfile(fmt.Sprintf(models.RnIgmpIfP), LogicalInterfaceProfileDn, desc, igmpIfPAttr)

	igmpIfP.Status = "modified"

	err := aciClient.Save(igmpIfP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_igmp_rs_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_igmp_rs_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_igmp_rs_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_igmp_rs_if_pol")
		err = aciClient.DeleteRelationigmpRsIfPol(igmpIfP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationigmpRsIfPol(igmpIfP.DistinguishedName, igmpIfPAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(igmpIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciIGMPInterfaceProfileRead(ctx, d, m)
}

func resourceAciIGMPInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	igmpIfP, err := getRemoteIGMPInterfaceProfile(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setIGMPInterfaceProfileAttributes(igmpIfP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// Get and Set Relational Attributes
	getAndSetIGMPInterfaceProfileRelationalAttributes(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciIGMPInterfaceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "igmpIfP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

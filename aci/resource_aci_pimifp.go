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

func resourceAciPimInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciPimIPInterfaceProfileCreate,
		UpdateContext: resourceAciPimIPInterfaceProfileUpdate,
		ReadContext:   resourceAciPimIPInterfaceProfileRead,
		DeleteContext: resourceAciPimIPInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPimIPInterfaceProfileImport,
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

			"relation_pim_rs_if_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to pim:IfPol",
			}})),
	}
}

func getRemotePimInterfaceProfile(client *client.Client, dn string) (*models.PimInterfaceProfile, error) {
	pimIfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	pimIfP := models.PimInterfaceProfileFromContainer(pimIfPCont)
	if pimIfP.DistinguishedName == "" {
		return nil, fmt.Errorf("Pim Interface Profile %s not found", dn)
	}
	return pimIfP, nil
}

func setPimInterfaceProfileAttributes(pimIfP *models.PimInterfaceProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(pimIfP.DistinguishedName)
	d.Set("description", pimIfP.Description)
	pimIfPMap, err := pimIfP.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != pimIfP.DistinguishedName {
		d.Set("logical_interface_profile_dn", "")
	} else {
		d.Set("logical_interface_profile_dn", GetParentDn(pimIfP.DistinguishedName, fmt.Sprintf("/"+models.RnPimIfP)))
	}
	d.Set("annotation", pimIfPMap["annotation"])
	d.Set("name", pimIfPMap["name"])
	d.Set("name_alias", pimIfPMap["nameAlias"])
	return d, nil
}

func getAndSetPimInterfaceProfileRelationalAttributes(client *client.Client, dn string, d *schema.ResourceData) {

	log.Printf("[DEBUG] pimRsIfPol: Beginning Read")

	pimRsIfPolData, err := client.ReadRelationPimRsIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation pimRsIfPol %v", err)
		d.Set("relation_pim_rs_if_pol", "")
	} else {
		d.Set("relation_pim_rs_if_pol", pimRsIfPolData.(string))
		log.Printf("[DEBUG]: pimRsIfPol: Reading finished successfully")
	}
}

func resourceAciPimIPInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	pimIfP, err := getRemotePimInterfaceProfile(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPimInterfaceProfileAttributes(pimIfP, d)
	if err != nil {
		return nil, err
	}

	// Get and Set Relational Attributes
	getAndSetPimInterfaceProfileRelationalAttributes(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPimIPInterfaceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InterfaceProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	pimIfPAttr := models.PimInterfaceProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		pimIfPAttr.Annotation = Annotation.(string)
	} else {
		pimIfPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		pimIfPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		pimIfPAttr.NameAlias = NameAlias.(string)
	}
	pimIfP := models.NewPimInterfaceProfile(fmt.Sprintf(models.RnPimIfP), LogicalInterfaceProfileDn, desc, pimIfPAttr)

	err := aciClient.Save(pimIfP)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTopimRsIfPol, ok := d.GetOk("relation_pim_rs_if_pol"); ok {
		relationParam := relationTopimRsIfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTopimRsIfPol, ok := d.GetOk("relation_pim_rs_if_pol"); ok {
		relationParam := relationTopimRsIfPol.(string)
		err = aciClient.CreateRelationPimRsIfPol(pimIfP.DistinguishedName, pimIfPAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(pimIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciPimIPInterfaceProfileRead(ctx, d, m)
}
func resourceAciPimIPInterfaceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InterfaceProfile: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	pimIfPAttr := models.PimInterfaceProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		pimIfPAttr.Annotation = Annotation.(string)
	} else {
		pimIfPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		pimIfPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		pimIfPAttr.NameAlias = NameAlias.(string)
	}
	pimIfP := models.NewPimInterfaceProfile(fmt.Sprintf(models.RnPimIfP), LogicalInterfaceProfileDn, desc, pimIfPAttr)

	pimIfP.Status = "modified"

	err := aciClient.Save(pimIfP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_pim_rs_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_pim_rs_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_pim_rs_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_pim_rs_if_pol")
		err = aciClient.DeleteRelationPimRsIfPol(pimIfP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationPimRsIfPol(pimIfP.DistinguishedName, pimIfPAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(pimIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciPimIPInterfaceProfileRead(ctx, d, m)
}

func resourceAciPimIPInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	pimIfP, err := getRemotePimInterfaceProfile(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setPimInterfaceProfileAttributes(pimIfP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// Get and Set Relational Attributes
	getAndSetPimInterfaceProfileRelationalAttributes(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciPimIPInterfaceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "pimIfP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

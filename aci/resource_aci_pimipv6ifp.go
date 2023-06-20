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

func resourceAciPimIPv6InterfaceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciPimIPv6InterfaceProfileCreate,
		UpdateContext: resourceAciPimIPv6InterfaceProfileUpdate,
		ReadContext:   resourceAciPimIPv6InterfaceProfileRead,
		DeleteContext: resourceAciPimIPv6InterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPimIPv6InterfaceProfileImport,
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

			"relation_pim_ipv6_rs_if_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to pim:IfPol",
			}})),
	}
}

func getRemotePimIPv6InterfaceProfile(client *client.Client, dn string) (*models.PimIPv6InterfaceProfile, error) {
	pimIPV6IfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	pimIPV6IfP := models.PimIPv6InterfaceProfileFromContainer(pimIPV6IfPCont)
	if pimIPV6IfP.DistinguishedName == "" {
		return nil, fmt.Errorf("Pim Interface Profile %s not found", dn)
	}
	return pimIPV6IfP, nil
}

func setPimIPv6InterfaceProfileAttributes(pimIPV6IfP *models.PimIPv6InterfaceProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(pimIPV6IfP.DistinguishedName)
	d.Set("description", pimIPV6IfP.Description)
	pimIPV6IfPMap, err := pimIPV6IfP.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != pimIPV6IfP.DistinguishedName {
		d.Set("logical_interface_profile_dn", "")
	} else {
		d.Set("logical_interface_profile_dn", GetParentDn(pimIPV6IfP.DistinguishedName, fmt.Sprintf("/"+models.RnPimIPV6IfP)))
	}
	d.Set("annotation", pimIPV6IfPMap["annotation"])
	d.Set("name", pimIPV6IfPMap["name"])
	d.Set("name_alias", pimIPV6IfPMap["nameAlias"])
	return d, nil
}

func getAndSetPimIPv6InterfaceProfileRelationalAttributes(client *client.Client, dn string, d *schema.ResourceData) {

	log.Printf("[DEBUG] pimRsV6IfPol: Beginning Read")

	pimRsV6IfPolData, err := client.ReadRelationPimIPv6RsIfPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation pimRsV6IfPol %v", err)
		d.Set("relation_pim_ipv6_rs_if_pol", "")
	} else {
		d.Set("relation_pim_ipv6_rs_if_pol", pimRsV6IfPolData.(string))
		log.Printf("[DEBUG]: pimRsV6IfPol: Read finished successfully")
	}
}

func resourceAciPimIPv6InterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	pimIPV6IfP, err := getRemotePimIPv6InterfaceProfile(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPimIPv6InterfaceProfileAttributes(pimIPV6IfP, d)
	if err != nil {
		return nil, err
	}

	// Get and Set Relational Attributes
	getAndSetPimIPv6InterfaceProfileRelationalAttributes(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPimIPv6InterfaceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InterfaceProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	pimIPV6IfPAttr := models.PimIPv6InterfaceProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		pimIPV6IfPAttr.Annotation = Annotation.(string)
	} else {
		pimIPV6IfPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		pimIPV6IfPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		pimIPV6IfPAttr.NameAlias = NameAlias.(string)
	}
	pimIPV6IfP := models.NewPimIPv6InterfaceProfile(fmt.Sprintf(models.RnPimIPV6IfP), LogicalInterfaceProfileDn, desc, pimIPV6IfPAttr)

	err := aciClient.Save(pimIPV6IfP)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTopimRsV6IfPol, ok := d.GetOk("relation_pim_ipv6_rs_if_pol"); ok {
		relationParam := relationTopimRsV6IfPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTopimRsV6IfPol, ok := d.GetOk("relation_pim_ipv6_rs_if_pol"); ok {
		relationParam := relationTopimRsV6IfPol.(string)
		err = aciClient.CreateRelationPimIPv6RsIfPol(pimIPV6IfP.DistinguishedName, pimIPV6IfPAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(pimIPV6IfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciPimIPv6InterfaceProfileRead(ctx, d, m)
}
func resourceAciPimIPv6InterfaceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InterfaceProfile: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	LogicalInterfaceProfileDn := d.Get("logical_interface_profile_dn").(string)

	pimIPV6IfPAttr := models.PimIPv6InterfaceProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		pimIPV6IfPAttr.Annotation = Annotation.(string)
	} else {
		pimIPV6IfPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		pimIPV6IfPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		pimIPV6IfPAttr.NameAlias = NameAlias.(string)
	}
	pimIPV6IfP := models.NewPimIPv6InterfaceProfile(fmt.Sprintf(models.RnPimIPV6IfP), LogicalInterfaceProfileDn, desc, pimIPV6IfPAttr)

	pimIPV6IfP.Status = "modified"

	err := aciClient.Save(pimIPV6IfP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_pim_ipv6_rs_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_pim_ipv6_rs_if_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_pim_ipv6_rs_if_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_pim_ipv6_rs_if_pol")
		err = aciClient.DeleteRelationPimIPv6RsIfPol(pimIPV6IfP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationPimIPv6RsIfPol(pimIPV6IfP.DistinguishedName, pimIPV6IfPAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(pimIPV6IfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciPimIPv6InterfaceProfileRead(ctx, d, m)
}

func resourceAciPimIPv6InterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	pimIPV6IfP, err := getRemotePimIPv6InterfaceProfile(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setPimIPv6InterfaceProfileAttributes(pimIPV6IfP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// Get and Set Relational Attributes
	getAndSetPimIPv6InterfaceProfileRelationalAttributes(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciPimIPv6InterfaceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "pimIPV6IfP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

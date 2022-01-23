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

func resourceAciSpineAccessPortSelector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSpineAccessPortSelectorCreate,
		UpdateContext: resourceAciSpineAccessPortSelectorUpdate,
		ReadContext:   resourceAciSpineAccessPortSelectorRead,
		DeleteContext: resourceAciSpineAccessPortSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSpineAccessPortSelectorImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"spine_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"spine_access_port_selector_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ALL",
					"range",
				}, false),
			},

			"relation_infra_rs_sp_acc_grp": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to infra:SpAccGrp",
			}})),
	}
}

func getRemoteSpineAccessPortSelector(client *client.Client, dn string) (*models.SpineAccessPortSelector, error) {
	infraSHPortSCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSHPortS := models.SpineAccessPortSelectorFromContainer(infraSHPortSCont)

	if infraSHPortS.DistinguishedName == "" {
		return nil, fmt.Errorf("SpineAccessPortSelector %s not found", infraSHPortS.DistinguishedName)
	}

	return infraSHPortS, nil
}

func setSpineAccessPortSelectorAttributes(infraSHPortS *models.SpineAccessPortSelector, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(infraSHPortS.DistinguishedName)
	d.Set("description", infraSHPortS.Description)

	if dn != infraSHPortS.DistinguishedName {
		d.Set("spine_interface_profile_dn", "")
	}
	infraSHPortSMap, err := infraSHPortS.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("spine_interface_profile_dn", GetParentDn(dn, fmt.Sprintf("/shports-%s-typ-%s", infraSHPortSMap["name"], infraSHPortSMap["type"])))
	d.Set("name", infraSHPortSMap["name"])
	d.Set("name_alias", infraSHPortSMap["nameAlias"])
	d.Set("annotation", infraSHPortSMap["annotation"])
	d.Set("spine_access_port_selector_type", infraSHPortSMap["type"])

	return d, nil
}

func resourceAciSpineAccessPortSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraSHPortS, err := getRemoteSpineAccessPortSelector(aciClient, dn)
	if err != nil {
		return nil, err
	}

	schemaFilled, err := setSpineAccessPortSelectorAttributes(infraSHPortS, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSpineAccessPortSelectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineAccessPortSelector: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	spine_access_port_selector_type := d.Get("spine_access_port_selector_type").(string)
	SpineInterfaceProfileDn := d.Get("spine_interface_profile_dn").(string)

	infraSHPortSAttr := models.SpineAccessPortSelectorAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSHPortSAttr.Annotation = Annotation.(string)
	} else {
		infraSHPortSAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		infraSHPortSAttr.Name = Name.(string)
	}

	if SpineAccessPortSelector_type, ok := d.GetOk("spine_access_port_selector_type"); ok {
		infraSHPortSAttr.SpineAccessPortSelector_type = SpineAccessPortSelector_type.(string)
	}
	infraSHPortS := models.NewSpineAccessPortSelector(fmt.Sprintf(models.RninfraSHPortS, name, spine_access_port_selector_type), SpineInterfaceProfileDn, desc, nameAlias, infraSHPortSAttr)

	err := aciClient.Save(infraSHPortS)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToinfraRsSpAccGrp, ok := d.GetOk("relation_infra_rs_sp_acc_grp"); ok {
		relationParam := relationToinfraRsSpAccGrp.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToinfraRsSpAccGrp, ok := d.GetOk("relation_infra_rs_sp_acc_grp"); ok {
		relationParam := relationToinfraRsSpAccGrp.(string)
		err = aciClient.CreateRelationinfraRsSpAccGrp(infraSHPortS.DistinguishedName, infraSHPortSAttr.Annotation, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(infraSHPortS.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSpineAccessPortSelectorRead(ctx, d, m)
}

func resourceAciSpineAccessPortSelectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineAccessPortSelector: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	spine_access_port_selector_type := d.Get("spine_access_port_selector_type").(string)
	SpineInterfaceProfileDn := d.Get("spine_interface_profile_dn").(string)
	infraSHPortSAttr := models.SpineAccessPortSelectorAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSHPortSAttr.Annotation = Annotation.(string)
	} else {
		infraSHPortSAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		infraSHPortSAttr.Name = Name.(string)
	}

	if SpineAccessPortSelector_type, ok := d.GetOk("spine_access_port_selector_type"); ok {
		infraSHPortSAttr.SpineAccessPortSelector_type = SpineAccessPortSelector_type.(string)
	}
	infraSHPortS := models.NewSpineAccessPortSelector(fmt.Sprintf(models.RninfraSHPortS, name, spine_access_port_selector_type), SpineInterfaceProfileDn, desc, nameAlias, infraSHPortSAttr)

	infraSHPortS.Status = "modified"
	err := aciClient.Save(infraSHPortS)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_sp_acc_grp") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_sp_acc_grp")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_sp_acc_grp") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_infra_rs_sp_acc_grp")
		err = aciClient.DeleteRelationinfraRsSpAccGrp(infraSHPortS.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsSpAccGrp(infraSHPortS.DistinguishedName, infraSHPortSAttr.Annotation, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(infraSHPortS.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSpineAccessPortSelectorRead(ctx, d, m)
}

func resourceAciSpineAccessPortSelectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraSHPortS, err := getRemoteSpineAccessPortSelector(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setSpineAccessPortSelectorAttributes(infraSHPortS, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsSpAccGrpData, err := aciClient.ReadRelationinfraRsSpAccGrp(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpAccGrp %v", err)
		d.Set("relation_infra_rs_sp_acc_grp", "")
	} else {
		if _, ok := d.GetOk("relation_infra_rs_sp_acc_grp"); ok {
			tfName := d.Get("relation_infra_rs_sp_acc_grp").(string)
			if tfName != infraRsSpAccGrpData {
				d.Set("relation_infra_rs_sp_acc_grp", "")
			}
		}
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciSpineAccessPortSelectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraSHPortS")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")

	return diag.FromErr(err)
}

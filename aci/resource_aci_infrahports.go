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

func resourceAciAccessPortSelector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciAccessPortSelectorCreate,
		UpdateContext: resourceAciAccessPortSelectorUpdate,
		ReadContext:   resourceAciAccessPortSelectorRead,
		DeleteContext: resourceAciAccessPortSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAccessPortSelectorImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"leaf_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"access_port_selector_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ALL",
					"range",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_infra_rs_acc_base_grp": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteAccessPortSelector(client *client.Client, dn string) (*models.AccessPortSelector, error) {
	infraHPortSCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraHPortS := models.AccessPortSelectorFromContainer(infraHPortSCont)

	if infraHPortS.DistinguishedName == "" {
		return nil, fmt.Errorf("AccessPortSelector %s not found", infraHPortS.DistinguishedName)
	}

	return infraHPortS, nil
}

func setAccessPortSelectorAttributes(infraHPortS *models.AccessPortSelector, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(infraHPortS.DistinguishedName)
	d.Set("description", infraHPortS.Description)

	if dn != infraHPortS.DistinguishedName {
		d.Set("leaf_interface_profile_dn", "")
	}
	infraHPortSMap, err := infraHPortS.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("leaf_interface_profile_dn", GetParentDn(dn, fmt.Sprintf("/hports-%s-typ-%s", infraHPortSMap["name"], infraHPortSMap["type"])))

	d.Set("name", infraHPortSMap["name"])

	d.Set("access_port_selector_type", infraHPortSMap["type"])

	d.Set("annotation", infraHPortSMap["annotation"])
	d.Set("name_alias", infraHPortSMap["nameAlias"])
	d.Set("access_port_selector_type", infraHPortSMap["type"])
	return d, nil
}

func resourceAciAccessPortSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraHPortS, err := getRemoteAccessPortSelector(aciClient, dn)

	if err != nil {
		return nil, err
	}
	infraHPortSMap, err := infraHPortS.ToMap()
	if err != nil {
		return nil, err
	}
	name := infraHPortSMap["name"]
	ptype := infraHPortSMap["type"]
	pDN := GetParentDn(dn, fmt.Sprintf("/hports-%s-typ-%s", name, ptype))
	d.Set("leaf_interface_profile_dn", pDN)
	schemaFilled, err := setAccessPortSelectorAttributes(infraHPortS, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAccessPortSelectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AccessPortSelector: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	access_port_selector_type := d.Get("access_port_selector_type").(string)

	LeafInterfaceProfileDn := d.Get("leaf_interface_profile_dn").(string)

	infraHPortSAttr := models.AccessPortSelectorAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraHPortSAttr.Annotation = Annotation.(string)
	} else {
		infraHPortSAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraHPortSAttr.NameAlias = NameAlias.(string)
	}
	if AccessPortSelector_type, ok := d.GetOk("access_port_selector_type"); ok {
		infraHPortSAttr.AccessPortSelector_type = AccessPortSelector_type.(string)
	}
	infraHPortS := models.NewAccessPortSelector(fmt.Sprintf("hports-%s-typ-%s", name, access_port_selector_type), LeafInterfaceProfileDn, desc, infraHPortSAttr)

	err := aciClient.Save(infraHPortS)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToinfraRsAccBaseGrp, ok := d.GetOk("relation_infra_rs_acc_base_grp"); ok {
		relationParam := relationToinfraRsAccBaseGrp.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToinfraRsAccBaseGrp, ok := d.GetOk("relation_infra_rs_acc_base_grp"); ok {
		relationParam := relationToinfraRsAccBaseGrp.(string)
		err = aciClient.CreateRelationinfraRsAccBaseGrpFromAccessPortSelector(infraHPortS.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraHPortS.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciAccessPortSelectorRead(ctx, d, m)
}

func resourceAciAccessPortSelectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AccessPortSelector: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	access_port_selector_type := d.Get("access_port_selector_type").(string)

	LeafInterfaceProfileDn := d.Get("leaf_interface_profile_dn").(string)

	infraHPortSAttr := models.AccessPortSelectorAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraHPortSAttr.Annotation = Annotation.(string)
	} else {
		infraHPortSAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraHPortSAttr.NameAlias = NameAlias.(string)
	}
	if AccessPortSelector_type, ok := d.GetOk("access_port_selector_type"); ok {
		infraHPortSAttr.AccessPortSelector_type = AccessPortSelector_type.(string)
	}
	infraHPortS := models.NewAccessPortSelector(fmt.Sprintf("hports-%s-typ-%s", name, access_port_selector_type), LeafInterfaceProfileDn, desc, infraHPortSAttr)

	infraHPortS.Status = "modified"

	err := aciClient.Save(infraHPortS)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_acc_base_grp") {
		_, newRelParam := d.GetChange("relation_infra_rs_acc_base_grp")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_acc_base_grp") {
		_, newRelParam := d.GetChange("relation_infra_rs_acc_base_grp")
		err = aciClient.DeleteRelationinfraRsAccBaseGrpFromAccessPortSelector(infraHPortS.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsAccBaseGrpFromAccessPortSelector(infraHPortS.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraHPortS.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciAccessPortSelectorRead(ctx, d, m)

}

func resourceAciAccessPortSelectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraHPortS, err := getRemoteAccessPortSelector(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setAccessPortSelectorAttributes(infraHPortS, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsAccBaseGrpData, err := aciClient.ReadRelationinfraRsAccBaseGrpFromAccessPortSelector(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAccBaseGrp %v", err)
		d.Set("relation_infra_rs_acc_base_grp", "")

	} else {
		if _, ok := d.GetOk("relation_infra_rs_acc_base_grp"); ok {
			tfName := d.Get("relation_infra_rs_acc_base_grp").(string)
			if tfName != infraRsAccBaseGrpData {
				d.Set("relation_infra_rs_acc_base_grp", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciAccessPortSelectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraHPortS")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

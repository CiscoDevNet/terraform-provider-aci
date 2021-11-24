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

func resourceAciL2Outside() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL2OutsideCreate,
		UpdateContext: resourceAciL2OutsideUpdate,
		ReadContext:   resourceAciL2OutsideRead,
		DeleteContext: resourceAciL2OutsideDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL2OutsideImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"target_dscp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AF11",
					"AF12",
					"AF13",
					"AF21",
					"AF22",
					"AF23",
					"AF31",
					"AF32",
					"AF33",
					"AF41",
					"AF42",
					"AF43",
					"CS0",
					"CS1",
					"CS2",
					"CS3",
					"CS4",
					"CS5",
					"CS6",
					"CS7",
					"EF",
					"VA",
					"unspecified",
				}, false),
			},

			"relation_l2ext_rs_e_bd": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_l2ext_rs_l2_dom_att": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteL2Outside(client *client.Client, dn string) (*models.L2Outside, error) {
	l2extOutCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l2extOut := models.L2OutsideFromContainer(l2extOutCont)

	if l2extOut.DistinguishedName == "" {
		return nil, fmt.Errorf("L2Outside %s not found", l2extOut.DistinguishedName)
	}

	return l2extOut, nil
}

func setL2OutsideAttributes(l2extOut *models.L2Outside, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(l2extOut.DistinguishedName)
	d.Set("description", l2extOut.Description)
	if dn != l2extOut.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	l2extOutMap, err := l2extOut.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", l2extOutMap["name"])
	d.Set("tenant_dn", GetParentDn(l2extOut.DistinguishedName, fmt.Sprintf("/l2out-%s", l2extOutMap["name"])))

	d.Set("annotation", l2extOutMap["annotation"])
	d.Set("name_alias", l2extOutMap["nameAlias"])
	d.Set("target_dscp", l2extOutMap["targetDscp"])
	return d, nil
}

func resourceAciL2OutsideImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l2extOut, err := getRemoteL2Outside(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL2OutsideAttributes(l2extOut, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL2OutsideCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L2Outside: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	l2extOutAttr := models.L2OutsideAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2extOutAttr.Annotation = Annotation.(string)
	} else {
		l2extOutAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2extOutAttr.NameAlias = NameAlias.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l2extOutAttr.TargetDscp = TargetDscp.(string)
	}
	l2extOut := models.NewL2Outside(fmt.Sprintf("l2out-%s", name), TenantDn, desc, l2extOutAttr)

	err := aciClient.Save(l2extOut)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTol2extRsEBd, ok := d.GetOk("relation_l2ext_rs_e_bd"); ok {
		relationParam := relationTol2extRsEBd.(string)
		checkDns = append(checkDns, relationParam)

	}
	if relationTol2extRsL2DomAtt, ok := d.GetOk("relation_l2ext_rs_l2_dom_att"); ok {
		relationParam := relationTol2extRsL2DomAtt.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTol2extRsEBd, ok := d.GetOk("relation_l2ext_rs_e_bd"); ok {
		relationParam := relationTol2extRsEBd.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationl2extRsEBdFromL2Outside(l2extOut.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTol2extRsL2DomAtt, ok := d.GetOk("relation_l2ext_rs_l2_dom_att"); ok {
		relationParam := relationTol2extRsL2DomAtt.(string)
		err = aciClient.CreateRelationl2extRsL2DomAttFromL2Outside(l2extOut.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(l2extOut.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL2OutsideRead(ctx, d, m)
}

func resourceAciL2OutsideUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L2Outside: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	l2extOutAttr := models.L2OutsideAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2extOutAttr.Annotation = Annotation.(string)
	} else {
		l2extOutAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2extOutAttr.NameAlias = NameAlias.(string)
	}
	if TargetDscp, ok := d.GetOk("target_dscp"); ok {
		l2extOutAttr.TargetDscp = TargetDscp.(string)
	}
	l2extOut := models.NewL2Outside(fmt.Sprintf("l2out-%s", name), TenantDn, desc, l2extOutAttr)

	l2extOut.Status = "modified"

	err := aciClient.Save(l2extOut)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_l2ext_rs_e_bd") {
		_, newRelParam := d.GetChange("relation_l2ext_rs_e_bd")
		checkDns = append(checkDns, newRelParam.(string))

	}
	if d.HasChange("relation_l2ext_rs_l2_dom_att") {
		_, newRelParam := d.GetChange("relation_l2ext_rs_l2_dom_att")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_l2ext_rs_e_bd") {
		_, newRelParam := d.GetChange("relation_l2ext_rs_e_bd")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationl2extRsEBdFromL2Outside(l2extOut.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_l2ext_rs_l2_dom_att") {
		_, newRelParam := d.GetChange("relation_l2ext_rs_l2_dom_att")
		err = aciClient.DeleteRelationl2extRsL2DomAttFromL2Outside(l2extOut.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationl2extRsL2DomAttFromL2Outside(l2extOut.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(l2extOut.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL2OutsideRead(ctx, d, m)

}

func resourceAciL2OutsideRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l2extOut, err := getRemoteL2Outside(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setL2OutsideAttributes(l2extOut, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	l2extRsEBdData, err := aciClient.ReadRelationl2extRsEBdFromL2Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l2extRsEBd %v", err)
		d.Set("relation_l2ext_rs_e_bd", "")

	} else {
		d.Set("relation_l2ext_rs_e_bd", l2extRsEBdData.(string))
	}

	l2extRsL2DomAttData, err := aciClient.ReadRelationl2extRsL2DomAttFromL2Outside(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation l2extRsL2DomAtt %v", err)
		d.Set("relation_l2ext_rs_l2_dom_att", "")

	} else {
		d.Set("relation_l2ext_rs_l2_dom_att", l2extRsL2DomAttData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL2OutsideDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l2extOut")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

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

func resourceAciL3outBGPProtocolProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL3outBGPProtocolProfileCreate,
		UpdateContext: resourceAciL3outBGPProtocolProfileUpdate,
		ReadContext:   resourceAciL3outBGPProtocolProfileRead,
		DeleteContext: resourceAciL3outBGPProtocolProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outBGPProtocolProfileImport,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"logical_node_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "orchestrator:terraform",
			},
			"relation_bgp_rs_best_path_ctrl_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_bgp_rs_bgp_node_ctx_pol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
func getRemoteL3outBGPProtocolProfile(client *client.Client, dn string) (*models.L3outBGPProtocolProfile, error) {
	bgpProtPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpProtP := models.L3outBGPProtocolProfileFromContainer(bgpProtPCont)

	if bgpProtP.DistinguishedName == "" {
		return nil, fmt.Errorf("L3Out BGP Protocol Profile %s not found", dn)
	}

	return bgpProtP, nil
}

func setL3outBGPProtocolProfileAttributes(bgpProtP *models.L3outBGPProtocolProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(bgpProtP.DistinguishedName)
	dn := d.Id()
	if dn != bgpProtP.DistinguishedName {
		d.Set("logical_node_profile_dn", "")
	}
	bgpProtPMap, err := bgpProtP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("logical_node_profile_dn", GetParentDn(dn, fmt.Sprintf("/%s", models.RnbgpProtP)))
	d.Set("annotation", bgpProtPMap["annotation"])
	d.Set("name", bgpProtPMap["name"])
	d.Set("name_alias", bgpProtPMap["nameAlias"])

	return d, nil
}

func resourceAciL3outBGPProtocolProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bgpProtP, err := getRemoteL3outBGPProtocolProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL3outBGPProtocolProfileAttributes(bgpProtP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outBGPProtocolProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outBGPProtocolProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	bgpProtPAttr := models.L3outBGPProtocolProfileAttributes{}
	if Name, ok := d.GetOk("name"); ok {
		bgpProtPAttr.Name = Name.(string)
	}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpProtPAttr.Annotation = Annotation.(string)
	} else {
		bgpProtPAttr.Annotation = "{}"
	}

	bgpProtP := models.NewL3outBGPProtocolProfile(fmt.Sprintf(models.RnbgpProtP), LogicalNodeProfileDn, nameAlias, bgpProtPAttr)

	err := aciClient.Save(bgpProtP)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTobgpRsBestPathCtrlPol, ok := d.GetOk("relation_bgp_rs_best_path_ctrl_pol"); ok {
		relationParam := relationTobgpRsBestPathCtrlPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTobgpRsBgpNodeCtxPol, ok := d.GetOk("relation_bgp_rs_bgp_node_ctx_pol"); ok {
		relationParam := relationTobgpRsBgpNodeCtxPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTobgpRsBestPathCtrlPol, ok := d.GetOk("relation_bgp_rs_best_path_ctrl_pol"); ok {
		relationParam := GetMOName(relationTobgpRsBestPathCtrlPol.(string))
		err = aciClient.CreateRelationbgpRsBestPathCtrlPol(bgpProtP.DistinguishedName, bgpProtPAttr.Annotation, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTobgpRsBgpNodeCtxPol, ok := d.GetOk("relation_bgp_rs_bgp_node_ctx_pol"); ok {
		relationParam := GetMOName(relationTobgpRsBgpNodeCtxPol.(string))
		err = aciClient.CreateRelationbgpRsBgpNodeCtxPolFromL3outBGPProtocolProfile(bgpProtP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(bgpProtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciL3outBGPProtocolProfileRead(ctx, d, m)
}

func resourceAciL3outBGPProtocolProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3outBGPProtocolProfile: Beginning Update")
	aciClient := m.(*client.Client)
	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	bgpProtPAttr := models.L3outBGPProtocolProfileAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpProtPAttr.Annotation = Annotation.(string)
	} else {
		bgpProtPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		bgpProtPAttr.Name = Name.(string)
	}
	bgpProtP := models.NewL3outBGPProtocolProfile(fmt.Sprintf(models.RnbgpProtP), LogicalNodeProfileDn, nameAlias, bgpProtPAttr)

	bgpProtP.Status = "modified"

	err := aciClient.Save(bgpProtP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_bgp_rs_best_path_ctrl_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_bgp_rs_best_path_ctrl_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_bgp_rs_bgp_node_ctx_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_bgp_rs_bgp_node_ctx_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_bgp_rs_best_path_ctrl_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_bgp_rs_best_path_ctrl_pol")
		err = aciClient.DeleteRelationbgpRsBestPathCtrlPol(bgpProtP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationbgpRsBestPathCtrlPol(bgpProtP.DistinguishedName, bgpProtPAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_bgp_rs_bgp_node_ctx_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_bgp_rs_bgp_node_ctx_pol")
		err = aciClient.CreateRelationbgpRsBgpNodeCtxPolFromL3outBGPProtocolProfile(bgpProtP.DistinguishedName, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(bgpProtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciL3outBGPProtocolProfileRead(ctx, d, m)
}

func resourceAciL3outBGPProtocolProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	bgpProtP, err := getRemoteL3outBGPProtocolProfile(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setL3outBGPProtocolProfileAttributes(bgpProtP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	bgpRsBestPathCtrlPolData, err := aciClient.ReadRelationbgpRsBestPathCtrlPol(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation bgpRsBestPathCtrlPol %v", err)
		setRelationAttribute(d, "relation_bgp_rs_best_path_ctrl_pol", "")
	} else {
		setRelationAttribute(d, "relation_bgp_rs_best_path_ctrl_pol", bgpRsBestPathCtrlPolData)
	}

	bgpRsBgpNodeCtxPolData, err := aciClient.ReadRelationbgpRsBgpNodeCtxPolFromL3outBGPProtocolProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation bgpRsBgpNodeCtxPol %v", err)
		setRelationAttribute(d, "relation_bgp_rs_bgp_node_ctx_pol", "")
	} else {
		setRelationAttribute(d, "relation_bgp_rs_bgp_node_ctx_pol", bgpRsBgpNodeCtxPolData)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciL3outBGPProtocolProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "bgpProtP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

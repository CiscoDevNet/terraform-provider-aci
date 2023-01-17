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

func resourceAciL1L2RedirectDestTraffic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL1L2RedirectDestTrafficCreate,
		UpdateContext: resourceAciL1L2RedirectDestTrafficUpdate,
		ReadContext:   resourceAciL1L2RedirectDestTrafficRead,
		DeleteContext: resourceAciL1L2RedirectDestTrafficDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL1L2RedirectDestTrafficImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"policy_based_redirect_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pod_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_vns_rs_l1_l2_redirect_health_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to vns:RedirectHealthGroup",
			},
			"relation_vns_rs_to_c_if": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Create relation to vns:CIf",
			}})),
	}
}

func getRemoteL1L2RedirectDestTraffic(client *client.Client, dn string) (*models.L1L2RedirectDestTraffic, error) {
	vnsL1L2RedirectDestCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vnsL1L2RedirectDest := models.L1L2RedirectDestTrafficFromContainer(vnsL1L2RedirectDestCont)
	if vnsL1L2RedirectDest.DistinguishedName == "" {
		return nil, fmt.Errorf("L1/L2 Redirect Destination Traffic %s not found", dn)
	}
	return vnsL1L2RedirectDest, nil
}

func setL1L2RedirectDestTrafficAttributes(vnsL1L2RedirectDest *models.L1L2RedirectDestTraffic, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(vnsL1L2RedirectDest.DistinguishedName)
	d.Set("description", vnsL1L2RedirectDest.Description)

	vnsL1L2RedirectDestMap, err := vnsL1L2RedirectDest.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("policy_based_redirect_dn", GetParentDn(vnsL1L2RedirectDest.DistinguishedName, fmt.Sprintf("/%s", fmt.Sprintf(models.RnvnsL1L2RedirectDest, vnsL1L2RedirectDestMap["destName"]))))
	d.Set("annotation", vnsL1L2RedirectDestMap["annotation"])
	d.Set("destination_name", vnsL1L2RedirectDestMap["destName"])
	d.Set("mac", vnsL1L2RedirectDestMap["mac"])
	d.Set("name", vnsL1L2RedirectDestMap["name"])
	d.Set("pod_id", vnsL1L2RedirectDestMap["podId"])
	d.Set("name_alias", vnsL1L2RedirectDestMap["nameAlias"])
	return d, nil
}

func getAndSetRemoteReadRelationvnsRsL1L2RedirectHealthGroup(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	vnsRsL1L2RedirectHealthGroupData, err := client.ReadRelationvnsRsL1L2RedirectHealthGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsL1L2RedirectHealthGroup %v", err)
		d.Set("relation_vns_rs_l1_l2_redirect_health_group", nil)
		return nil, err
	} else {
		d.Set("relation_vns_rs_l1_l2_redirect_health_group", vnsRsL1L2RedirectHealthGroupData.(string))
		log.Printf("[DEBUG]: vnsRsRedirectHealthGroup: %s finished successfully", vnsRsL1L2RedirectHealthGroupData.(string))
	}
	return d, nil
}

func getAndSetRemoteReadRelationvnsRsToCIf(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	vnsRsToCIfData, err := client.ReadRelationvnsRsToCIf(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsToCIf %v", err)
		d.Set("relation_vns_rs_to_c_if", nil)
		return nil, err
	} else {
		d.Set("relation_vns_rs_to_c_if", vnsRsToCIfData.(string))
		log.Printf("[DEBUG]: vnsRsToCIf: %s finished successfully", vnsRsToCIfData.(string))
	}
	return d, nil
}

func resourceAciL1L2RedirectDestTrafficImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vnsL1L2RedirectDest, err := getRemoteL1L2RedirectDestTraffic(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL1L2RedirectDestTrafficAttributes(vnsL1L2RedirectDest, d)
	if err != nil {
		return nil, err
	}

	// Importing vnsRsL1L2RedirectHealthGroup object
	getAndSetRemoteReadRelationvnsRsL1L2RedirectHealthGroup(aciClient, dn, d)

	// Importing vnsRsToCIf object
	getAndSetRemoteReadRelationvnsRsToCIf(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL1L2RedirectDestTrafficCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L1L2RedirectDestTraffic: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	policyBasedRedirectDn := d.Get("policy_based_redirect_dn").(string)

	vnsL1L2RedirectDestAttr := models.L1L2RedirectDestTrafficAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsL1L2RedirectDestAttr.Annotation = Annotation.(string)
	} else {
		vnsL1L2RedirectDestAttr.Annotation = "{}"
	}

	if DestName, ok := d.GetOk("destination_name"); ok {
		vnsL1L2RedirectDestAttr.DestName = DestName.(string)
	}

	if Mac, ok := d.GetOk("mac"); ok {
		vnsL1L2RedirectDestAttr.Mac = Mac.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsL1L2RedirectDestAttr.Name = Name.(string)
	}

	if PodId, ok := d.GetOk("pod_id"); ok {
		vnsL1L2RedirectDestAttr.PodId = PodId.(string)
	}
	vnsL1L2RedirectDest := models.NewL1L2RedirectDestTraffic(fmt.Sprintf(models.RnvnsL1L2RedirectDest, vnsL1L2RedirectDestAttr.DestName), policyBasedRedirectDn, desc, nameAlias, vnsL1L2RedirectDestAttr)

	err := aciClient.Save(vnsL1L2RedirectDest)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTovnsRsL1L2RedirectHealthGroup, ok := d.GetOk("relation_vns_rs_l1_l2_redirect_health_group"); ok {
		relationParam := relationTovnsRsL1L2RedirectHealthGroup.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTovnsRsToCIf, ok := d.GetOk("relation_vns_rs_to_c_if"); ok {
		relationParam := relationTovnsRsToCIf.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovnsRsL1L2RedirectHealthGroup, ok := d.GetOk("relation_vns_rs_l1_l2_redirect_health_group"); ok {
		relationParam := relationTovnsRsL1L2RedirectHealthGroup.(string)
		err = aciClient.CreateRelationvnsRsL1L2RedirectHealthGroup(vnsL1L2RedirectDest.DistinguishedName, vnsL1L2RedirectDestAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTovnsRsToCIf, ok := d.GetOk("relation_vns_rs_to_c_if"); ok {
		relationParam := relationTovnsRsToCIf.(string)
		err = aciClient.CreateRelationvnsRsToCIf(vnsL1L2RedirectDest.DistinguishedName, vnsL1L2RedirectDestAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vnsL1L2RedirectDest.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciL1L2RedirectDestTrafficRead(ctx, d, m)
}

func resourceAciL1L2RedirectDestTrafficUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L1L2RedirectDestTraffic: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	policyBasedRedirectDn := d.Get("policy_based_redirect_dn").(string)

	vnsL1L2RedirectDestAttr := models.L1L2RedirectDestTrafficAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsL1L2RedirectDestAttr.Annotation = Annotation.(string)
	} else {
		vnsL1L2RedirectDestAttr.Annotation = "{}"
	}

	if DestName, ok := d.GetOk("destination_name"); ok {
		vnsL1L2RedirectDestAttr.DestName = DestName.(string)
	}

	if Mac, ok := d.GetOk("mac"); ok {
		vnsL1L2RedirectDestAttr.Mac = Mac.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsL1L2RedirectDestAttr.Name = Name.(string)
	}

	if PodId, ok := d.GetOk("pod_id"); ok {
		vnsL1L2RedirectDestAttr.PodId = PodId.(string)
	}
	vnsL1L2RedirectDest := models.NewL1L2RedirectDestTraffic(fmt.Sprintf(models.RnvnsL1L2RedirectDest, vnsL1L2RedirectDestAttr.DestName), policyBasedRedirectDn, desc, nameAlias, vnsL1L2RedirectDestAttr)

	err := aciClient.Save(vnsL1L2RedirectDest)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vns_rs_l1_l2_redirect_health_group") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vns_rs_l1_l2_redirect_health_group")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_vns_rs_to_c_if") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vns_rs_to_c_if")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vns_rs_l1_l2_redirect_health_group") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vns_rs_l1_l2_redirect_health_group")
		err = aciClient.DeleteRelationvnsRsL1L2RedirectHealthGroup(vnsL1L2RedirectDest.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsL1L2RedirectHealthGroup(vnsL1L2RedirectDest.DistinguishedName, vnsL1L2RedirectDestAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vns_rs_to_c_if") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vns_rs_to_c_if")
		err = aciClient.DeleteRelationvnsRsToCIf(vnsL1L2RedirectDest.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsToCIf(vnsL1L2RedirectDest.DistinguishedName, vnsL1L2RedirectDestAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vnsL1L2RedirectDest.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciL1L2RedirectDestTrafficRead(ctx, d, m)
}

func resourceAciL1L2RedirectDestTrafficRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vnsL1L2RedirectDest, err := getRemoteL1L2RedirectDestTraffic(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setL1L2RedirectDestTrafficAttributes(vnsL1L2RedirectDest, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// Reading vnsRsL1L2RedirectHealthGroup object
	getAndSetRemoteReadRelationvnsRsL1L2RedirectHealthGroup(aciClient, dn, d)

	// Reading vnsRsToCIf object
	getAndSetRemoteReadRelationvnsRsToCIf(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciL1L2RedirectDestTrafficDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "vnsL1L2RedirectDest")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

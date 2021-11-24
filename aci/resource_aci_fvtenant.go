package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciTenant() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciTenantCreate,
		UpdateContext: resourceAciTenantUpdate,
		ReadContext:   resourceAciTenantRead,
		DeleteContext: resourceAciTenantDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTenantImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

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

			"relation_fv_rs_tn_deny_rule": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_tenant_mon_pol": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		}),
	}
}
func getRemoteTenant(client *client.Client, dn string) (*models.Tenant, error) {
	fvTenantCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvTenant := models.TenantFromContainer(fvTenantCont)

	if fvTenant.DistinguishedName == "" {
		return nil, fmt.Errorf("Tenant %s not found", fvTenant.DistinguishedName)
	}

	return fvTenant, nil
}

func setTenantAttributes(fvTenant *models.Tenant, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fvTenant.DistinguishedName)
	d.Set("description", fvTenant.Description)
	fvTenantMap, err := fvTenant.ToMap()

	if err != nil {
		return d, err
	}

	d.Set("name", fvTenantMap["name"])

	d.Set("annotation", fvTenantMap["annotation"])
	d.Set("name_alias", fvTenantMap["nameAlias"])
	return d, nil
}

func resourceAciTenantImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvTenant, err := getRemoteTenant(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setTenantAttributes(fvTenant, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTenantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Tenant: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fvTenantAttr := models.TenantAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvTenantAttr.Annotation = Annotation.(string)
	} else {
		fvTenantAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvTenantAttr.NameAlias = NameAlias.(string)
	}
	fvTenant := models.NewTenant(fmt.Sprintf("tn-%s", name), "uni", desc, fvTenantAttr)

	err := aciClient.Save(fvTenant)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTofvRsTnDenyRule, ok := d.GetOk("relation_fv_rs_tn_deny_rule"); ok {
		relationParamList := toStringList(relationTofvRsTnDenyRule.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsTenantMonPol, ok := d.GetOk("relation_fv_rs_tenant_mon_pol"); ok {
		relationParam := relationTofvRsTenantMonPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTofvRsTnDenyRule, ok := d.GetOk("relation_fv_rs_tn_deny_rule"); ok {
		relationParamList := toStringList(relationTofvRsTnDenyRule.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsTnDenyRuleFromTenant(fvTenant.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationTofvRsTenantMonPol, ok := d.GetOk("relation_fv_rs_tenant_mon_pol"); ok {
		relationParam := relationTofvRsTenantMonPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsTenantMonPolFromTenant(fvTenant.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(fvTenant.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciTenantRead(ctx, d, m)
}

func resourceAciTenantUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Tenant: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fvTenantAttr := models.TenantAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvTenantAttr.Annotation = Annotation.(string)
	} else {
		fvTenantAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvTenantAttr.NameAlias = NameAlias.(string)
	}
	fvTenant := models.NewTenant(fmt.Sprintf("tn-%s", name), "uni", desc, fvTenantAttr)

	fvTenant.Status = "modified"

	err := aciClient.Save(fvTenant)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_fv_rs_tn_deny_rule") {
		oldRel, newRel := d.GetChange("relation_fv_rs_tn_deny_rule")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fv_rs_tenant_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_tenant_mon_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_fv_rs_tn_deny_rule") {
		oldRel, newRel := d.GetChange("relation_fv_rs_tn_deny_rule")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsTnDenyRuleFromTenant(fvTenant.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsTnDenyRuleFromTenant(fvTenant.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if d.HasChange("relation_fv_rs_tenant_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_tenant_mon_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsTenantMonPolFromTenant(fvTenant.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fvTenant.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciTenantRead(ctx, d, m)

}

func resourceAciTenantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvTenant, err := getRemoteTenant(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setTenantAttributes(fvTenant, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	fvRsTnDenyRuleData, err := aciClient.ReadRelationfvRsTnDenyRuleFromTenant(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsTnDenyRule %v", err)
		d.Set("relation_fv_rs_tn_deny_rule", make([]string, 0, 1))

	} else {
		d.Set("relation_fv_rs_tn_deny_rule", toStringList(fvRsTnDenyRuleData.(*schema.Set).List()))
	}

	fvRsTenantMonPolData, err := aciClient.ReadRelationfvRsTenantMonPolFromTenant(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsTenantMonPol %v", err)
		d.Set("relation_fv_rs_tenant_mon_pol", "")

	} else {
		d.Set("relation_fv_rs_tenant_mon_pol", fvRsTenantMonPolData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciTenantDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvTenant")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciTenant() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciTenantCreate,
		Update: resourceAciTenantUpdate,
		Read:   resourceAciTenantRead,
		Delete: resourceAciTenantDelete,

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
				Type: schema.TypeString,

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

func setTenantAttributes(fvTenant *models.Tenant, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvTenant.DistinguishedName)
	d.Set("description", fvTenant.Description)
	fvTenantMap, _ := fvTenant.ToMap()

	d.Set("name", fvTenantMap["name"])

	d.Set("annotation", fvTenantMap["annotation"])
	d.Set("name_alias", fvTenantMap["nameAlias"])
	return d
}

func resourceAciTenantImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvTenant, err := getRemoteTenant(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setTenantAttributes(fvTenant, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTenantCreate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTofvRsTnDenyRule, ok := d.GetOk("relation_fv_rs_tn_deny_rule"); ok {
		relationParamList := toStringList(relationTofvRsTnDenyRule.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsTnDenyRuleFromTenant(fvTenant.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_tn_deny_rule")
			d.Partial(false)
		}
	}
	if relationTofvRsTenantMonPol, ok := d.GetOk("relation_fv_rs_tenant_mon_pol"); ok {
		relationParam := relationTofvRsTenantMonPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfvRsTenantMonPolFromTenant(fvTenant.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_tenant_mon_pol")
		d.Partial(false)

	}

	d.SetId(fvTenant.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciTenantRead(d, m)
}

func resourceAciTenantUpdate(d *schema.ResourceData, m interface{}) error {
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
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

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
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsTnDenyRuleFromTenant(fvTenant.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_tn_deny_rule")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_tenant_mon_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_tenant_mon_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfvRsTenantMonPolFromTenant(fvTenant.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_tenant_mon_pol")
		d.Partial(false)

	}

	d.SetId(fvTenant.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciTenantRead(d, m)

}

func resourceAciTenantRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvTenant, err := getRemoteTenant(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setTenantAttributes(fvTenant, d)

	fvRsTnDenyRuleData, err := aciClient.ReadRelationfvRsTnDenyRuleFromTenant(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsTnDenyRule %v", err)

	} else {
		d.Set("relation_fv_rs_tn_deny_rule", fvRsTnDenyRuleData)
	}

	fvRsTenantMonPolData, err := aciClient.ReadRelationfvRsTenantMonPolFromTenant(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsTenantMonPol %v", err)

	} else {
		if _, ok := d.GetOk("relation_fv_rs_tenant_mon_pol"); ok {
			tfName := GetMOName(d.Get("relation_fv_rs_tenant_mon_pol").(string))
			if tfName != fvRsTenantMonPolData {
				d.Set("relation_fv_rs_tenant_mon_pol", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciTenantDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvTenant")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

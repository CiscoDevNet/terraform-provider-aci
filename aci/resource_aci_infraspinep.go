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

func resourceAciSpineProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSpineProfileCreate,
		UpdateContext: resourceAciSpineProfileUpdate,
		ReadContext:   resourceAciSpineProfileRead,
		DeleteContext: resourceAciSpineProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSpineProfileImport,
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

			"relation_infra_rs_sp_acc_port_p": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteSpineProfile(client *client.Client, dn string) (*models.SpineProfile, error) {
	infraSpinePCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSpineP := models.SpineProfileFromContainer(infraSpinePCont)

	if infraSpineP.DistinguishedName == "" {
		return nil, fmt.Errorf("SpineProfile %s not found", infraSpineP.DistinguishedName)
	}

	return infraSpineP, nil
}

func setSpineProfileAttributes(infraSpineP *models.SpineProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraSpineP.DistinguishedName)
	d.Set("description", infraSpineP.Description)
	infraSpinePMap, err := infraSpineP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", infraSpinePMap["name"])
	d.Set("annotation", infraSpinePMap["annotation"])
	d.Set("name_alias", infraSpinePMap["nameAlias"])
	return d, nil
}

func resourceAciSpineProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraSpineP, err := getRemoteSpineProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setSpineProfileAttributes(infraSpineP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSpineProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraSpinePAttr := models.SpineProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpinePAttr.Annotation = Annotation.(string)
	} else {
		infraSpinePAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpinePAttr.NameAlias = NameAlias.(string)
	}
	infraSpineP := models.NewSpineProfile(fmt.Sprintf("infra/spprof-%s", name), "uni", desc, infraSpinePAttr)

	err := aciClient.Save(infraSpineP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToinfraRsSpAccPortP, ok := d.GetOk("relation_infra_rs_sp_acc_port_p"); ok {
		relationParamList := toStringList(relationToinfraRsSpAccPortP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToinfraRsSpAccPortP, ok := d.GetOk("relation_infra_rs_sp_acc_port_p"); ok {
		relationParamList := toStringList(relationToinfraRsSpAccPortP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsSpAccPortPFromSpineProfile(infraSpineP.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(infraSpineP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSpineProfileRead(ctx, d, m)
}

func resourceAciSpineProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraSpinePAttr := models.SpineProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpinePAttr.Annotation = Annotation.(string)
	} else {
		infraSpinePAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpinePAttr.NameAlias = NameAlias.(string)
	}
	infraSpineP := models.NewSpineProfile(fmt.Sprintf("infra/spprof-%s", name), "uni", desc, infraSpinePAttr)

	infraSpineP.Status = "modified"

	err := aciClient.Save(infraSpineP)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_sp_acc_port_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_sp_acc_port_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_sp_acc_port_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_sp_acc_port_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationinfraRsSpAccPortPFromSpineProfile(infraSpineP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsSpAccPortPFromSpineProfile(infraSpineP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(infraSpineP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSpineProfileRead(ctx, d, m)

}

func resourceAciSpineProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraSpineP, err := getRemoteSpineProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setSpineProfileAttributes(infraSpineP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsSpAccPortPData, err := aciClient.ReadRelationinfraRsSpAccPortPFromSpineProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpAccPortP %v", err)
		d.Set("relation_infra_rs_sp_acc_port_p", make([]interface{}, 0, 1))

	} else {
		d.Set("relation_infra_rs_sp_acc_port_p", toStringList(infraRsSpAccPortPData.(*schema.Set).List()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSpineProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraSpineP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

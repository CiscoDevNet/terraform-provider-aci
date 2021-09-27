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

func resourceAciImportedContract() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciImportedContractCreate,
		UpdateContext: resourceAciImportedContractUpdate,
		ReadContext:   resourceAciImportedContractRead,
		DeleteContext: resourceAciImportedContractDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciImportedContractImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_vz_rs_if": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteImportedContract(client *client.Client, dn string) (*models.ImportedContract, error) {
	vzCPIfCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzCPIf := models.ImportedContractFromContainer(vzCPIfCont)

	if vzCPIf.DistinguishedName == "" {
		return nil, fmt.Errorf("ImportedContract %s not found", vzCPIf.DistinguishedName)
	}

	return vzCPIf, nil
}

func setImportedContractAttributes(vzCPIf *models.ImportedContract, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(vzCPIf.DistinguishedName)
	d.Set("description", vzCPIf.Description)

	if dn != vzCPIf.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	vzCPIfMap, err := vzCPIf.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/cif-%s", vzCPIfMap["name"])))

	d.Set("name", vzCPIfMap["name"])

	d.Set("annotation", vzCPIfMap["annotation"])
	d.Set("name_alias", vzCPIfMap["nameAlias"])
	return d, nil
}

func resourceAciImportedContractImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzCPIf, err := getRemoteImportedContract(aciClient, dn)

	if err != nil {
		return nil, err
	}
	vzCPIfMap, err := vzCPIf.ToMap()
	if err != nil {
		return nil, err
	}
	name := vzCPIfMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/cif-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setImportedContractAttributes(vzCPIf, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciImportedContractCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ImportedContract: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzCPIfAttr := models.ImportedContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzCPIfAttr.Annotation = Annotation.(string)
	} else {
		vzCPIfAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzCPIfAttr.NameAlias = NameAlias.(string)
	}
	vzCPIf := models.NewImportedContract(fmt.Sprintf("cif-%s", name), TenantDn, desc, vzCPIfAttr)

	err := aciClient.Save(vzCPIf)
	if err != nil {
		return diag.FromErr(err)

	}
	d.Partial(true)

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationTovzRsIf, ok := d.GetOk("relation_vz_rs_if"); ok {
		relationParam := relationTovzRsIf.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)

	}
	d.Partial(false)

	if relationTovzRsIf, ok := d.GetOk("relation_vz_rs_if"); ok {
		relationParam := relationTovzRsIf.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationvzRsIfFromImportedContract(vzCPIf.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)

		}
		d.Partial(true)
		d.Partial(false)

	}

	d.SetId(vzCPIf.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciImportedContractRead(ctx, d, m)
}

func resourceAciImportedContractUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ImportedContract: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vzCPIfAttr := models.ImportedContractAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzCPIfAttr.Annotation = Annotation.(string)
	} else {
		vzCPIfAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzCPIfAttr.NameAlias = NameAlias.(string)
	}
	vzCPIf := models.NewImportedContract(fmt.Sprintf("cif-%s", name), TenantDn, desc, vzCPIfAttr)

	vzCPIf.Status = "modified"

	err := aciClient.Save(vzCPIf)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(true)

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vz_rs_if") {
		_, newRelParam := d.GetChange("relation_vz_rs_if")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vz_rs_if") {
		_, newRelParam := d.GetChange("relation_vz_rs_if")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationvzRsIfFromImportedContract(vzCPIf.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvzRsIfFromImportedContract(vzCPIf.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Partial(true)
		d.Partial(false)

	}

	d.SetId(vzCPIf.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciImportedContractRead(ctx, d, m)

}

func resourceAciImportedContractRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzCPIf, err := getRemoteImportedContract(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setImportedContractAttributes(vzCPIf, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	vzRsIfData, err := aciClient.ReadRelationvzRsIfFromImportedContract(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsIf %v", err)
		d.Set("relation_vz_rs_if", "")

	} else {
		if _, ok := d.GetOk("relation_vz_rs_if"); ok {
			tfName := GetMOName(d.Get("relation_vz_rs_if").(string))
			if tfName != vzRsIfData {
				d.Set("relation_vz_rs_if", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciImportedContractDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzCPIf")
	if err != nil {
		return diag.FromErr(err)

	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)

}

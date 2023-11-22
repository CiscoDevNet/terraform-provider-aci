package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciContractProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciContractProviderCreate,
		UpdateContext: resourceAciContractProviderUpdate,
		ReadContext:   resourceAciContractProviderRead,
		DeleteContext: resourceAciContractProviderDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciContractProviderImport,
		},

		SchemaVersion: 1,

		Schema: AppendAttrSchemas(map[string]*schema.Schema{
			"application_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"contract_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"contract_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"consumer",
					"provider",
				}, false),
			},

			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"All",
					"AtleastOne",
					"AtmostOne",
					"None",
				}, false),
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"unspecified",
					"level6",
					"level5",
					"level4",
					"level3",
					"level2",
					"level1",
				}, false),
			},
		}, GetAnnotationAttrSchema()),
		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
			contractType := diff.Get("contract_type").(string)
			matchT := diff.Get("match_t").(string)
			if contractType == "consumer" && matchT != "" {
				return fmt.Errorf("MatchT is not supported for consumer contracts")
			}
			return nil
		},
	}
}

func getRemoteContractConsumer(client *client.Client, dn string) (*models.ContractConsumer, error) {
	fvRsConsCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvRsCons := models.ContractConsumerFromContainer(fvRsConsCont)
	if fvRsCons.DistinguishedName == "" {
		return nil, fmt.Errorf("Contract Consumer %s not found", dn)
	}

	return fvRsCons, nil
}

func getRemoteContractProvider(client *client.Client, dn string) (*models.ContractProvider, error) {
	fvRsProvCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvRsProv := models.ContractProviderFromContainer(fvRsProvCont)

	if fvRsProv.DistinguishedName == "" {
		return nil, fmt.Errorf("Contract Provider %s not found", dn)
	}

	return fvRsProv, nil
}

func setContractConsumerAttributes(fvRsCons *models.ContractConsumer, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvRsCons.DistinguishedName)

	if dn != fvRsCons.DistinguishedName {
		d.Set("application_epg_dn", "")
	}

	fvRsConsMap, err := fvRsCons.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("application_epg_dn", GetParentDn(fvRsCons.DistinguishedName, fmt.Sprintf("/rscons-%s", fvRsConsMap["tnVzBrCPName"])))

	tnVzBrCPName := GetMOName(d.Get("contract_dn").(string))
	if tnVzBrCPName != fvRsConsMap["tnVzBrCPName"] {
		d.Set("contract_dn", "")
	}
	d.Set("contract_dn", fvRsConsMap["tDn"])
	d.Set("annotation", fvRsConsMap["annotation"])
	d.Set("prio", fvRsConsMap["prio"])
	d.Set("contract_type", "consumer")
	return d, nil
}

func setContractProviderAttributes(fvRsProv *models.ContractProvider, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvRsProv.DistinguishedName)

	if dn != fvRsProv.DistinguishedName {
		d.Set("application_epg_dn", "")
	}
	// d.Set("application_epg_dn", GetParentDn(fvRsProv.DistinguishedName))
	fvRsProvMap, err := fvRsProv.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("application_epg_dn", GetParentDn(fvRsProv.DistinguishedName, fmt.Sprintf("/rsprov-%s", fvRsProvMap["tnVzBrCPName"])))
	tnVzBrCPName := GetMOName(d.Get("contract_dn").(string))
	if tnVzBrCPName != fvRsProvMap["tnVzBrCPName"] {
		d.Set("contract_dn", "")
	}
	d.Set("contract_dn", fvRsProvMap["tDn"])
	d.Set("annotation", fvRsProvMap["annotation"])
	d.Set("match_t", fvRsProvMap["matchT"])
	d.Set("prio", fvRsProvMap["prio"])
	d.Set("contract_type", "provider")
	return d, nil
}

func resourceAciContractProviderImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()
	rnValue := models.GetMORnPrefix(dn)

	var schemaFilled *schema.ResourceData

	if rnValue == "rsprov" {
		fvRsProv, err := getRemoteContractProvider(aciClient, dn)

		if err != nil {
			return nil, err
		}
		fvRsProvMap, _ := fvRsProv.ToMap()
		name := fvRsProvMap["tnVzBrCPName"]
		pDN := GetParentDn(dn, fmt.Sprintf("/rsprov-%s", name))
		d.Set("application_epg_dn", pDN)
		schemaFilled, err = setContractProviderAttributes(fvRsProv, d)
		if err != nil {
			return nil, err
		}
	} else if rnValue == "rscons" {
		fvRsCons, err := getRemoteContractConsumer(aciClient, dn)

		if err != nil {
			return nil, err
		}
		fvRsConsMap, _ := fvRsCons.ToMap()
		name := fvRsConsMap["tnVzBrCPName"]
		pDN := GetParentDn(dn, fmt.Sprintf("/rscons-%s", name))
		d.Set("application_epg_dn", pDN)
		schemaFilled, err = setContractConsumerAttributes(fvRsCons, d)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Failed to import, invalid DN: %s", dn)
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciContractProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return postContractConfig(ctx, "", d, m)
}

func resourceAciContractProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return postContractConfig(ctx, "modified", d, m)
}

func resourceAciContractProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	contractType := d.Get("contract_type").(string)

	if contractType == "provider" {
		fvRsProv, err := getRemoteContractProvider(aciClient, dn)
		if err != nil {
			return errorForObjectNotFound(err, dn, d)
		}
		_, err = setContractProviderAttributes(fvRsProv, d)
		if err != nil {
			d.SetId("")
			return nil
		}
	} else if contractType == "consumer" {
		fvRsCons, err := getRemoteContractConsumer(aciClient, dn)
		if err != nil {
			return errorForObjectNotFound(err, dn, d)
		}
		_, err = setContractConsumerAttributes(fvRsCons, d)
		if err != nil {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciContractProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	contractType := d.Get("contract_type").(string)

	if contractType == "provider" {
		err := aciClient.DeleteByDn(dn, "fvRsProv")
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
		d.SetId("")
		return diag.FromErr(err)

	} else if contractType == "consumer" {
		err := aciClient.DeleteByDn(dn, "fvRsCons")
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
		d.SetId("")
		return diag.FromErr(err)
	}
	return nil
}

func postContractConfig(ctx context.Context, status string, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	contractOperation := ""
	if status == "modified" {
		contractOperation = "Update"
	} else {
		contractOperation = "Create"
	}

	aciClient := m.(*client.Client)

	tnVzBrCPName := GetMOName(d.Get("contract_dn").(string))

	contractType := d.Get("contract_type").(string)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	if contractType == "provider" {
		log.Printf("[DEBUG] ContractProvider: Beginning %s", contractOperation)

		MatchT := d.Get("match_t").(string)
		fvRsProvAttr := models.ContractProviderAttributes{}
		if Annotation, ok := d.GetOk("annotation"); ok {
			fvRsProvAttr.Annotation = Annotation.(string)
		} else {
			fvRsProvAttr.Annotation = "{}"
		}
		if MatchT != "" {
			fvRsProvAttr.MatchT = MatchT
		}
		if Prio, ok := d.GetOk("prio"); ok {
			fvRsProvAttr.Prio = Prio.(string)
		}
		fvRsProvAttr.TnVzBrCPName = tnVzBrCPName
		fvRsProv := models.NewContractProvider(fmt.Sprintf("rsprov-%s", tnVzBrCPName), ApplicationEPGDn, fvRsProvAttr)

		fvRsProv.Status = status

		err := aciClient.Save(fvRsProv)

		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fvRsProv.DistinguishedName)
	} else if contractType == "consumer" {
		log.Printf("[DEBUG] ContractConsumer: Beginning %s", contractOperation)

		fvRsConsAttr := models.ContractConsumerAttributes{}
		if Annotation, ok := d.GetOk("annotation"); ok {
			fvRsConsAttr.Annotation = Annotation.(string)
		} else {
			fvRsConsAttr.Annotation = "{}"
		}
		if Prio, ok := d.GetOk("prio"); ok {
			fvRsConsAttr.Prio = Prio.(string)
		}
		fvRsConsAttr.TnVzBrCPName = tnVzBrCPName
		fvRsCons := models.NewContractConsumer(fmt.Sprintf("rscons-%s", tnVzBrCPName), ApplicationEPGDn, fvRsConsAttr)

		fvRsCons.Status = status

		err := aciClient.Save(fvRsCons)

		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fvRsCons.DistinguishedName)
	}
	if status == "modified" {
		log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	} else {
		log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	}
	return resourceAciContractProviderRead(ctx, d, m)
}

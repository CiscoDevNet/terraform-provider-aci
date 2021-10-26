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

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
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
		}),
	}
}

func getRemoteContractConsumer(client *client.Client, dn string) (*models.ContractConsumer, error) {
	fvRsConsCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvRsCons := models.ContractConsumerFromContainer(fvRsConsCont)
	if fvRsCons.DistinguishedName == "" {
		return nil, fmt.Errorf("ContractConsumer %s not found", fvRsCons.DistinguishedName)
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
		return nil, fmt.Errorf("ContractProvider %s not found", fvRsProv.DistinguishedName)
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
	return d, nil
}

func resourceAciContractProviderImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()
	contractType := d.Get("contract_type").(string)
	var schemaFilled *schema.ResourceData

	if contractType == "provider" {
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

	} else if contractType == "consumer" {
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
		return nil, fmt.Errorf("Contract Type: Value must be from [provider, consumer]")
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciContractProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	aciClient := m.(*client.Client)

	tnVzBrCPName := GetMOName(d.Get("contract_dn").(string))

	contractType := d.Get("contract_type").(string)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	if contractType == "provider" {
		log.Printf("[DEBUG] ContractProvider: Beginning Creation")

		fvRsProvAttr := models.ContractProviderAttributes{}
		if Annotation, ok := d.GetOk("annotation"); ok {
			fvRsProvAttr.Annotation = Annotation.(string)
		} else {
			fvRsProvAttr.Annotation = "{}"
		}
		if MatchT, ok := d.GetOk("match_t"); ok {
			fvRsProvAttr.MatchT = MatchT.(string)
		}
		if Prio, ok := d.GetOk("prio"); ok {
			fvRsProvAttr.Prio = Prio.(string)
		}
		fvRsProvAttr.TnVzBrCPName = tnVzBrCPName
		fvRsProv := models.NewContractProvider(fmt.Sprintf("rsprov-%s", tnVzBrCPName), ApplicationEPGDn, fvRsProvAttr)

		err := aciClient.Save(fvRsProv)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fvRsProv.DistinguishedName)

	} else if contractType == "consumer" {
		log.Printf("[DEBUG] ContractConsumer: Beginning Creation")

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

		err := aciClient.Save(fvRsCons)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fvRsCons.DistinguishedName)

	} else {
		return diag.Errorf(fmt.Sprintf("Contract Type: Value must be from [provider, consumer] = %s", contractType))
	}

	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciContractProviderRead(ctx, d, m)
}

func resourceAciContractProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	aciClient := m.(*client.Client)

	tnVzBrCPName := GetMOName(d.Get("contract_dn").(string))

	contractType := d.Get("contract_type").(string)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	if contractType == "provider" {
		log.Printf("[DEBUG] ContractProvider: Beginning Update")

		fvRsProvAttr := models.ContractProviderAttributes{}
		if Annotation, ok := d.GetOk("annotation"); ok {
			fvRsProvAttr.Annotation = Annotation.(string)
		} else {
			fvRsProvAttr.Annotation = "{}"
		}
		if MatchT, ok := d.GetOk("match_t"); ok {
			fvRsProvAttr.MatchT = MatchT.(string)
		}
		if Prio, ok := d.GetOk("prio"); ok {
			fvRsProvAttr.Prio = Prio.(string)
		}
		fvRsProvAttr.TnVzBrCPName = tnVzBrCPName
		fvRsProv := models.NewContractProvider(fmt.Sprintf("rsprov-%s", tnVzBrCPName), ApplicationEPGDn, fvRsProvAttr)

		fvRsProv.Status = "modified"

		err := aciClient.Save(fvRsProv)

		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fvRsProv.DistinguishedName)

	} else if contractType == "consumer" {
		log.Printf("[DEBUG] ContractConsumer: Beginning Update")

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

		fvRsCons.Status = "modified"

		err := aciClient.Save(fvRsCons)

		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fvRsCons.DistinguishedName)

	} else {
		return diag.Errorf("Contract Type: Value must be from [provider, consumer]")
	}

	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciContractProviderRead(ctx, d, m)

}

func resourceAciContractProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)
	contractType := d.Get("contract_type").(string)
	dn := d.Id()

	if contractType == "provider" {
		fvRsProv, err := getRemoteContractProvider(aciClient, dn)
		if err != nil {
			d.SetId("")
			return nil
		}
		_, err = setContractProviderAttributes(fvRsProv, d)
		if err != nil {
			d.SetId("")
			return nil
		}

	} else if contractType == "consumer" {
		fvRsCons, err := getRemoteContractConsumer(aciClient, dn)
		if err != nil {
			d.SetId("")
			return nil
		}
		_, err = setContractConsumerAttributes(fvRsCons, d)
		if err != nil {
			d.SetId("")
			return nil
		}

	} else {
		return diag.Errorf(fmt.Sprintf("Contract Type: Value must be from [provider, consumer] = %s", contractType))
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

	} else {
		return diag.Errorf("Contract Type: Value must be from [provider, consumer]")
	}

}

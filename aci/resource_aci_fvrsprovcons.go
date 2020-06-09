package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciContractProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciContractProviderCreate,
		Update: resourceAciContractProviderUpdate,
		Read:   resourceAciContractProviderRead,
		Delete: resourceAciContractProviderDelete,

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

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func setContractConsumerAttributes(fvRsCons *models.ContractConsumer, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(fvRsCons.DistinguishedName)

	if dn != fvRsCons.DistinguishedName {
		d.Set("application_epg_dn", "")
	}

	// d.Set("application_epg_dn", GetParentDn(fvRsCons.DistinguishedName))
	fvRsConsMap, _ := fvRsCons.ToMap()

	tnVzBrCPName := GetMOName(d.Get("contract_dn").(string))
	if tnVzBrCPName != fvRsConsMap["tnVzBrCPName"] {
		d.Set("contract_dn", "")
	}

	d.Set("annotation", fvRsConsMap["annotation"])
	d.Set("prio", fvRsConsMap["prio"])
	return d
}

func setContractProviderAttributes(fvRsProv *models.ContractProvider, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(fvRsProv.DistinguishedName)

	if dn != fvRsProv.DistinguishedName {
		d.Set("application_epg_dn", "")
	}
	// d.Set("application_epg_dn", GetParentDn(fvRsProv.DistinguishedName))
	fvRsProvMap, _ := fvRsProv.ToMap()
	tnVzBrCPName := GetMOName(d.Get("contract_dn").(string))
	if tnVzBrCPName != fvRsProvMap["tnVzBrCPName"] {
		d.Set("contract_dn", "")
	}

	d.Set("annotation", fvRsProvMap["annotation"])
	d.Set("match_t", fvRsProvMap["matchT"])
	d.Set("prio", fvRsProvMap["prio"])
	return d
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
		schemaFilled = setContractProviderAttributes(fvRsProv, d)

	} else if contractType == "consumer" {
		fvRsCons, err := getRemoteContractConsumer(aciClient, dn)

		if err != nil {
			return nil, err
		}
		schemaFilled = setContractConsumerAttributes(fvRsCons, d)

	} else {
		return nil, fmt.Errorf("Contract Type: Value must be from [provider, consumer]")
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciContractProviderCreate(d *schema.ResourceData, m interface{}) error {

	aciClient := m.(*client.Client)

	tnVzBrCPName := GetMOName(d.Get("contract_dn").(string))

	contractType := d.Get("contract_type").(string)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	if contractType == "provider" {
		log.Printf("[DEBUG] ContractProvider: Beginning Creation")

		fvRsProvAttr := models.ContractProviderAttributes{}
		if Annotation, ok := d.GetOk("annotation"); ok {
			fvRsProvAttr.Annotation = Annotation.(string)
		}
		if MatchT, ok := d.GetOk("match_t"); ok {
			fvRsProvAttr.MatchT = MatchT.(string)
		}
		if Prio, ok := d.GetOk("prio"); ok {
			fvRsProvAttr.Prio = Prio.(string)
		}

		fvRsProv := models.NewContractProvider(fmt.Sprintf("rsprov-%s", tnVzBrCPName), ApplicationEPGDn, fvRsProvAttr)

		err := aciClient.Save(fvRsProv)
		if err != nil {
			return err
		}
		d.Partial(true)

		d.SetPartial("tnVzBrCPName")

		d.Partial(false)

		d.SetId(fvRsProv.DistinguishedName)

	} else if contractType == "consumer" {
		log.Printf("[DEBUG] ContractConsumer: Beginning Creation")

		fvRsConsAttr := models.ContractConsumerAttributes{}
		if Annotation, ok := d.GetOk("annotation"); ok {
			fvRsConsAttr.Annotation = Annotation.(string)
		}
		if Prio, ok := d.GetOk("prio"); ok {
			fvRsConsAttr.Prio = Prio.(string)
		}
		if TnVzBrCPName, ok := d.GetOk("contract_name"); ok {
			fvRsConsAttr.TnVzBrCPName = TnVzBrCPName.(string)
		}
		fvRsCons := models.NewContractConsumer(fmt.Sprintf("rscons-%s", tnVzBrCPName), ApplicationEPGDn, fvRsConsAttr)

		err := aciClient.Save(fvRsCons)
		if err != nil {
			return err
		}
		d.Partial(true)

		d.SetPartial("tnVzBrCPName")

		d.Partial(false)

		d.SetId(fvRsCons.DistinguishedName)

	} else {
		return fmt.Errorf("Contract Type: Value must be from [provider, consumer]")
	}

	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciContractProviderRead(d, m)
}

func resourceAciContractProviderUpdate(d *schema.ResourceData, m interface{}) error {

	aciClient := m.(*client.Client)

	tnVzBrCPName := GetMOName(d.Get("contract_name").(string))

	contractType := d.Get("contract_type").(string)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	if contractType == "provider" {
		log.Printf("[DEBUG] ContractProvider: Beginning Update")

		fvRsProvAttr := models.ContractProviderAttributes{}
		if Annotation, ok := d.GetOk("annotation"); ok {
			fvRsProvAttr.Annotation = Annotation.(string)
		}
		if MatchT, ok := d.GetOk("match_t"); ok {
			fvRsProvAttr.MatchT = MatchT.(string)
		}
		if Prio, ok := d.GetOk("prio"); ok {
			fvRsProvAttr.Prio = Prio.(string)
		}

		fvRsProv := models.NewContractProvider(fmt.Sprintf("rsprov-%s", tnVzBrCPName), ApplicationEPGDn, fvRsProvAttr)

		fvRsProv.Status = "modified"

		err := aciClient.Save(fvRsProv)

		if err != nil {
			return err
		}
		d.Partial(true)

		d.SetPartial("tnVzBrCPName")

		d.Partial(false)

		d.SetId(fvRsProv.DistinguishedName)

	} else if contractType == "consumer" {
		log.Printf("[DEBUG] ContractConsumer: Beginning Update")

		fvRsConsAttr := models.ContractConsumerAttributes{}
		if Annotation, ok := d.GetOk("annotation"); ok {
			fvRsConsAttr.Annotation = Annotation.(string)
		}
		if Prio, ok := d.GetOk("prio"); ok {
			fvRsConsAttr.Prio = Prio.(string)
		}
		if TnVzBrCPName, ok := d.GetOk("contract_name"); ok {
			fvRsConsAttr.TnVzBrCPName = TnVzBrCPName.(string)
		}
		fvRsCons := models.NewContractConsumer(fmt.Sprintf("rscons-%s", tnVzBrCPName), ApplicationEPGDn, fvRsConsAttr)

		fvRsCons.Status = "modified"

		err := aciClient.Save(fvRsCons)

		if err != nil {
			return err
		}
		d.Partial(true)

		d.SetPartial("tnVzBrCPName")

		d.Partial(false)

		d.SetId(fvRsCons.DistinguishedName)

	} else {
		return fmt.Errorf("Contract Type: Value must be from [provider, consumer]")
	}

	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciContractProviderRead(d, m)

}

func resourceAciContractProviderRead(d *schema.ResourceData, m interface{}) error {
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
		setContractProviderAttributes(fvRsProv, d)

	} else if contractType == "consumer" {
		fvRsCons, err := getRemoteContractConsumer(aciClient, dn)
		if err != nil {
			d.SetId("")
			return nil
		}
		setContractConsumerAttributes(fvRsCons, d)

	} else {
		return fmt.Errorf("Contract Type: Value must be from [provider, consumer]")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciContractProviderDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	contractType := d.Get("contract_type").(string)

	if contractType == "provider" {
		err := aciClient.DeleteByDn(dn, "fvRsProv")
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
		d.SetId("")
		return err

	} else if contractType == "consumer" {
		err := aciClient.DeleteByDn(dn, "fvRsCons")
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
		d.SetId("")
		return err

	} else {
		return fmt.Errorf("Contract Type: Value must be from [provider, consumer]")
	}

}

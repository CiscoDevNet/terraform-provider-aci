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

func resourceAciMulticastAddressPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMulticastAddressPoolCreate,
		UpdateContext: resourceAciMulticastAddressPoolUpdate,
		ReadContext:   resourceAciMulticastAddressPoolRead,
		DeleteContext: resourceAciMulticastAddressPoolDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMulticastAddressPoolImport,
		},

		SchemaVersion: 1,
		Schema: AppendAttrSchemas(map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"multicast_address_block": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: AppendAttrSchemas(map[string]*schema.Schema{
						// Internally used for updating (delete/re-create) of multicast_address_block
						"dn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"to": {
							Type:     schema.TypeString,
							Required: true,
						},
					}, GetBaseAttrSchema(), GetNameAliasAttrSchema()),
				},
			},
		}, GetBaseAttrSchema(), GetNameAliasAttrSchema()),
	}
}

func getRemoteMulticastAddressPool(client *client.Client, dn string) (*models.MulticastAddressPool, error) {
	fvnsMcastAddrInstPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvnsMcastAddrInstP := models.MulticastAddressPoolFromContainer(fvnsMcastAddrInstPCont)
	if fvnsMcastAddrInstP.DistinguishedName == "" {
		return nil, fmt.Errorf("Multicast Address Pool %s not found", dn)
	}
	return fvnsMcastAddrInstP, nil
}

func setMulticastAddressPoolAttributes(fvnsMcastAddrInstP *models.MulticastAddressPool, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fvnsMcastAddrInstP.DistinguishedName)
	d.Set("description", fvnsMcastAddrInstP.Description)
	fvnsMcastAddrInstPMap, err := fvnsMcastAddrInstP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("annotation", fvnsMcastAddrInstPMap["annotation"])
	d.Set("name", fvnsMcastAddrInstPMap["name"])
	d.Set("name_alias", fvnsMcastAddrInstPMap["nameAlias"])
	return d, nil
}

func getMulticastAddressBlocks(callType, multicastAddressPool string, client *client.Client, d *schema.ResourceData) []*models.MulticastAddressBlock {
	log.Printf("[DEBUG] Beginning GET called by %s function for address pool name %s", callType, multicastAddressPool)
	readMulticastAddressBlockData, err := client.ListMulticastAddressBlock(multicastAddressPool)
	if err == nil {
		log.Printf("[DEBUG] Finished GET called by %s successfully with result: %v", callType, d.Get("multicast_address_block"))
	} else {
		log.Printf("[DEBUG] Error during GET operation of multicast address blocks: %v", err)
		readMulticastAddressBlockData = nil
	}
	return readMulticastAddressBlockData
}

func setMulticastAddressBlocks(callType, multicastAddressPool string, readMulticastAddressBlockData []*models.MulticastAddressBlock, client *client.Client, d *schema.ResourceData) {
	log.Printf("[DEBUG] Beginning SET called by %s function for address blocks %v", callType, readMulticastAddressBlockData)
	multicastAddressBlockList := make([]interface{}, 0)
	for _, record := range readMulticastAddressBlockData {
		multicastAddressBlockMap := make(map[string]interface{})
		multicastAddressBlockMap["annotation"] = record.Annotation
		multicastAddressBlockMap["from"] = record.From
		multicastAddressBlockMap["name"] = record.Name
		multicastAddressBlockMap["name_alias"] = record.NameAlias
		multicastAddressBlockMap["to"] = record.To
		multicastAddressBlockMap["dn"] = fmt.Sprintf(models.DnfvnsMcastAddrBlk, multicastAddressPool, record.From, record.To)
		multicastAddressBlockList = append(multicastAddressBlockList, multicastAddressBlockMap)
	}
	d.Set("multicast_address_block", multicastAddressBlockList)
	log.Printf("[DEBUG] Finished SET called by %s successfully with result: %v", callType, multicastAddressBlockList)
}

func resourceAciMulticastAddressPoolImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvnsMcastAddrInstP, err := getRemoteMulticastAddressPool(aciClient, dn)
	if err != nil {
		return nil, err
	}
	_, err = setMulticastAddressPoolAttributes(fvnsMcastAddrInstP, d)
	if err != nil {
		return nil, err
	}

	multicastAddressBlocks := getMulticastAddressBlocks("Import", fvnsMcastAddrInstP.Name, aciClient, d)
	if multicastAddressBlocks != nil {
		setMulticastAddressBlocks("Import", fvnsMcastAddrInstP.Name, multicastAddressBlocks, aciClient, d)
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceAciMulticastAddressPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MulticastAddressPool: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	fvnsMcastAddrInstPAttr := models.MulticastAddressPoolAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsMcastAddrInstPAttr.Annotation = Annotation.(string)
	} else {
		fvnsMcastAddrInstPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		fvnsMcastAddrInstPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsMcastAddrInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsMcastAddrInstP := models.NewMulticastAddressPool(fmt.Sprintf(models.RnfvnsMcastAddrInstP, name), models.ParentDnfvnsMcastAddrInstP, desc, fvnsMcastAddrInstPAttr)
	err := aciClient.Save(fvnsMcastAddrInstP)
	if err != nil {
		return diag.FromErr(err)
	}

	if multicastAddressBlock, ok := d.GetOk("multicast_address_block"); ok {
		multicastAddressBlockList := multicastAddressBlock.(*schema.Set).List()
		multicastAddressBlocks := make([]*models.MulticastAddressBlock, 0)
		for _, block := range multicastAddressBlockList {
			blockMap := block.(map[string]interface{})
			fvnsMcastAddrBlkAttr := models.MulticastAddressBlockAttributes{}
			fvnsMcastAddrBlkAttr.Annotation = blockMap["annotation"].(string)
			fvnsMcastAddrBlkAttr.From = blockMap["from"].(string)
			fvnsMcastAddrBlkAttr.Name = blockMap["name"].(string)
			fvnsMcastAddrBlkAttr.To = blockMap["to"].(string)
			if blockMap["nameAlias"] != nil {
				fvnsMcastAddrBlkAttr.NameAlias = blockMap["nameAlias"].(string)
			}
			desc := ""
			if blockMap["description"] != nil {
				desc = blockMap["description"].(string)
			}
			fvnsMcastAddrBlk := models.NewMulticastAddressBlock(fmt.Sprintf(models.RnfvnsMcastAddrBlk, fvnsMcastAddrBlkAttr.From, fvnsMcastAddrBlkAttr.To), fvnsMcastAddrInstP.DistinguishedName, desc, fvnsMcastAddrBlkAttr)
			err := aciClient.Save(fvnsMcastAddrBlk)
			if err != nil {
				return diag.FromErr(err)
			}
			multicastAddressBlocks = append(multicastAddressBlocks, fvnsMcastAddrBlk)
		}
		setMulticastAddressBlocks("Create", fvnsMcastAddrInstP.Name, multicastAddressBlocks, aciClient, d)
	}

	d.SetId(fvnsMcastAddrInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciMulticastAddressPoolRead(ctx, d, m)
}

func resourceAciMulticastAddressPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MulticastAddressPool: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)

	fvnsMcastAddrInstPAttr := models.MulticastAddressPoolAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsMcastAddrInstPAttr.Annotation = Annotation.(string)
	} else {
		fvnsMcastAddrInstPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		fvnsMcastAddrInstPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsMcastAddrInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsMcastAddrInstP := models.NewMulticastAddressPool(fmt.Sprintf(models.RnfvnsMcastAddrInstP, name), models.ParentDnfvnsMcastAddrInstP, desc, fvnsMcastAddrInstPAttr)
	err := aciClient.Save(fvnsMcastAddrInstP)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("multicast_address_block") {
		oldSchemaObjs, newSchemaObjs := d.GetChange("multicast_address_block")
		missingOldObjects := getOldObjectsNotInNew("dn", oldSchemaObjs.(*schema.Set), newSchemaObjs.(*schema.Set))

		for _, missingOldObject := range missingOldObjects {
			err := aciClient.DeleteByDn(missingOldObject.(map[string]interface{})["dn"].(string), fmt.Sprintf(models.FvnsmcastaddrblkClassName))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		multicastAddressBlockList := newSchemaObjs.(*schema.Set).List()
		multicastAddressBlocks := make([]*models.MulticastAddressBlock, 0)
		for _, block := range multicastAddressBlockList {
			blockMap := block.(map[string]interface{})
			log.Printf("[DEBUG] blockMap: %v", blockMap)
			fvnsMcastAddrBlkAttr := models.MulticastAddressBlockAttributes{}
			fvnsMcastAddrBlkAttr.Annotation = blockMap["annotation"].(string)
			fvnsMcastAddrBlkAttr.From = blockMap["from"].(string)
			fvnsMcastAddrBlkAttr.Name = blockMap["name"].(string)
			fvnsMcastAddrBlkAttr.To = blockMap["to"].(string)
			if blockMap["nameAlias"] != nil {
				fvnsMcastAddrBlkAttr.NameAlias = blockMap["nameAlias"].(string)
			}
			desc := ""
			if blockMap["description"] != nil {
				desc = blockMap["description"].(string)
			}
			fvnsMcastAddrBlk := models.NewMulticastAddressBlock(fmt.Sprintf(models.RnfvnsMcastAddrBlk, fvnsMcastAddrBlkAttr.From, fvnsMcastAddrBlkAttr.To), fvnsMcastAddrInstP.DistinguishedName, desc, fvnsMcastAddrBlkAttr)
			err := aciClient.Save(fvnsMcastAddrBlk)
			if err != nil {
				return diag.FromErr(err)
			}
			blockMap["dn"] = fvnsMcastAddrBlk.DistinguishedName
			multicastAddressBlocks = append(multicastAddressBlocks, fvnsMcastAddrBlk)
		}
		setMulticastAddressBlocks("Update", fvnsMcastAddrInstP.Name, multicastAddressBlocks, aciClient, d)
	}

	d.SetId(fvnsMcastAddrInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciMulticastAddressPoolRead(ctx, d, m)
}

func resourceAciMulticastAddressPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	fvnsMcastAddrInstP, err := getRemoteMulticastAddressPool(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setMulticastAddressPoolAttributes(fvnsMcastAddrInstP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	if _, ok := d.GetOk("multicast_address_block"); ok {
		multicastAddressBlocks := getMulticastAddressBlocks("Read", fvnsMcastAddrInstP.Name, aciClient, d)
		if multicastAddressBlocks != nil {
			setMulticastAddressBlocks("Read", fvnsMcastAddrInstP.Name, multicastAddressBlocks, aciClient, d)
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciMulticastAddressPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "fvnsMcastAddrInstP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

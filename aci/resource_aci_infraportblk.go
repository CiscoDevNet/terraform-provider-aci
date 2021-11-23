package aci

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciAccessPortBlock() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciAccessPortBlockCreate,
		UpdateContext: resourceAciAccessPortBlockUpdate,
		ReadContext:   resourceAciAccessPortBlockRead,
		DeleteContext: resourceAciAccessPortBlockDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAccessPortBlockImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"access_port_selector_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"from_card": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"from_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"to_card": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"to_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_infra_rs_acc_bndl_subgrp": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteAccessPortBlock(client *client.Client, dn string) (*models.AccessPortBlock, error) {
	infraPortBlkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraPortBlk := models.AccessPortBlockFromContainer(infraPortBlkCont)

	if infraPortBlk.DistinguishedName == "" {
		return nil, fmt.Errorf("AccessPortBlock %s not found", infraPortBlk.DistinguishedName)
	}

	return infraPortBlk, nil
}

func setAccessPortBlockAttributes(infraPortBlk *models.AccessPortBlock, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(infraPortBlk.DistinguishedName)
	d.Set("description", infraPortBlk.Description)

	if dn != infraPortBlk.DistinguishedName {
		d.Set("access_port_selector_dn", "")
	}
	infraPortBlkMap, err := infraPortBlk.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("access_port_selector_dn", GetParentDn(dn, fmt.Sprintf("/portblk-%s", infraPortBlkMap["name"])))
	d.Set("name", infraPortBlkMap["name"])

	d.Set("annotation", infraPortBlkMap["annotation"])
	d.Set("from_card", infraPortBlkMap["fromCard"])
	d.Set("from_port", infraPortBlkMap["fromPort"])
	d.Set("name_alias", infraPortBlkMap["nameAlias"])
	d.Set("to_card", infraPortBlkMap["toCard"])
	d.Set("to_port", infraPortBlkMap["toPort"])
	return d, nil
}

func resourceAciAccessPortBlockImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraPortBlk, err := getRemoteAccessPortBlock(aciClient, dn)

	if err != nil {
		return nil, err
	}
	infraPortBlkMap, err := infraPortBlk.ToMap()
	if err != nil {
		return nil, err
	}

	name := infraPortBlkMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/portblk-%s", name))
	d.Set("access_port_selector_dn", pDN)
	schemaFilled, err := setAccessPortBlockAttributes(infraPortBlk, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAccessPortBlockCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AccessPortBlock: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	AccessPortSelectorDn := d.Get("access_port_selector_dn").(string)

	var name string
	if _, ok := d.GetOk("name"); !ok {
		baseurlStr := "/api/node/class"
		dnUrl := fmt.Sprintf("%s/%s/%s.json", baseurlStr, AccessPortSelectorDn, "infraPortBlk")

		cont, err := aciClient.GetViaURL(dnUrl)
		if err != nil {
			if models.G(cont, "totalCount") != "0" {
				return diag.FromErr(err)
			}
		}
		contList := models.ListFromContainer(cont, "infraPortBlk")
		contListLen := len(contList)

		blkNames := make([]string, 0, 1)
		for i := 0; i < contListLen; i++ {
			tp := models.G(contList[i], "name")
			blkNames = append(blkNames, tp)
		}
		log.Println("check .. : ", blkNames)

		cnt := contListLen + 1
		for true {
			flag := false
			tpName := fmt.Sprintf("Block%s", strconv.Itoa(cnt))
			for _, val := range blkNames {
				if val == tpName {
					flag = true
					cnt = cnt + 1
					break
				}
			}
			if !flag {
				name = tpName
				break
			}
		}
	} else {
		name = d.Get("name").(string)
	}

	infraPortBlkAttr := models.AccessPortBlockAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraPortBlkAttr.Annotation = Annotation.(string)
	} else {
		infraPortBlkAttr.Annotation = "{}"
	}
	if FromCard, ok := d.GetOk("from_card"); ok {
		infraPortBlkAttr.FromCard = FromCard.(string)
	}
	if FromPort, ok := d.GetOk("from_port"); ok {
		infraPortBlkAttr.FromPort = FromPort.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraPortBlkAttr.NameAlias = NameAlias.(string)
	}
	if ToCard, ok := d.GetOk("to_card"); ok {
		infraPortBlkAttr.ToCard = ToCard.(string)
	}
	if ToPort, ok := d.GetOk("to_port"); ok {
		infraPortBlkAttr.ToPort = ToPort.(string)
	}
	infraPortBlk := models.NewAccessPortBlock(fmt.Sprintf("portblk-%s", name), AccessPortSelectorDn, desc, infraPortBlkAttr)

	err := aciClient.Save(infraPortBlk)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToinfraRsAccBndlSubgrp, ok := d.GetOk("relation_infra_rs_acc_bndl_subgrp"); ok {
		relationParam := relationToinfraRsAccBndlSubgrp.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToinfraRsAccBndlSubgrp, ok := d.GetOk("relation_infra_rs_acc_bndl_subgrp"); ok {
		relationParam := relationToinfraRsAccBndlSubgrp.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsAccBndlSubgrpFromAccessPortBlock(infraPortBlk.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraPortBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciAccessPortBlockRead(ctx, d, m)
}

func resourceAciAccessPortBlockUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AccessPortBlock: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	AccessPortSelectorDn := d.Get("access_port_selector_dn").(string)

	infraPortBlkAttr := models.AccessPortBlockAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraPortBlkAttr.Annotation = Annotation.(string)
	} else {
		infraPortBlkAttr.Annotation = "{}"
	}
	if FromCard, ok := d.GetOk("from_card"); ok {
		infraPortBlkAttr.FromCard = FromCard.(string)
	}
	if FromPort, ok := d.GetOk("from_port"); ok {
		infraPortBlkAttr.FromPort = FromPort.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraPortBlkAttr.NameAlias = NameAlias.(string)
	}
	if ToCard, ok := d.GetOk("to_card"); ok {
		infraPortBlkAttr.ToCard = ToCard.(string)
	}
	if ToPort, ok := d.GetOk("to_port"); ok {
		infraPortBlkAttr.ToPort = ToPort.(string)
	}
	infraPortBlk := models.NewAccessPortBlock(fmt.Sprintf("portblk-%s", name), AccessPortSelectorDn, desc, infraPortBlkAttr)

	infraPortBlk.Status = "modified"

	err := aciClient.Save(infraPortBlk)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_acc_bndl_subgrp") {
		_, newRelParam := d.GetChange("relation_infra_rs_acc_bndl_subgrp")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_acc_bndl_subgrp") {
		_, newRelParam := d.GetChange("relation_infra_rs_acc_bndl_subgrp")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationinfraRsAccBndlSubgrpFromAccessPortBlock(infraPortBlk.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsAccBndlSubgrpFromAccessPortBlock(infraPortBlk.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraPortBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciAccessPortBlockRead(ctx, d, m)

}

func resourceAciAccessPortBlockRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraPortBlk, err := getRemoteAccessPortBlock(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setAccessPortBlockAttributes(infraPortBlk, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsAccBndlSubgrpData, err := aciClient.ReadRelationinfraRsAccBndlSubgrpFromAccessPortBlock(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAccBndlSubgrp %v", err)
		d.Set("relation_infra_rs_acc_bndl_subgrp", "")

	} else {
		d.Set("relation_infra_rs_acc_bndl_subgrp", infraRsAccBndlSubgrpData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciAccessPortBlockDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraPortBlk")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

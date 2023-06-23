package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAccessPortBlock() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciAccessPortBlockRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"access_port_selector_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func dataSourceAciAccessPortBlockRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("portblk-%s", name)
	AccessPortSelectorDn := d.Get("access_port_selector_dn").(string)

	dn := fmt.Sprintf("%s/%s", AccessPortSelectorDn, rn)

	infraPortBlk, err := getRemoteAccessPortBlock(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setAccessPortBlockAttributes(infraPortBlk, d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: infraRsAccBndlSubgrp - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsAccBndlSubgrp(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsAccBndlSubgrp - Read finished successfully", d.Get("relation_infra_rs_acc_bndl_subgrp"))
	}

	return nil
}

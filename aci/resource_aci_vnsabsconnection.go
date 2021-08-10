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

func resourceAciConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciConnectionCreate,
		UpdateContext: resourceAciConnectionUpdate,
		ReadContext:   resourceAciConnectionRead,
		DeleteContext: resourceAciConnectionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciConnectionImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l4_l7_service_graph_template_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"adj_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"L2",
					"L3",
				}, false),
			},

			"conn_dir": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"consumer",
					"provider",
				}, false),
			},

			"conn_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"external",
					"internal",
				}, false),
			},

			"direct_connect": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"unicast_route": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
			},

			"relation_vns_rs_abs_copy_connection": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_vns_rs_abs_connection_conns": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteConnection(client *client.Client, dn string) (*models.Connection, error) {
	vnsAbsConnectionCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsConnection := models.ConnectionFromContainer(vnsAbsConnectionCont)

	if vnsAbsConnection.DistinguishedName == "" {
		return nil, fmt.Errorf("Connection %s not found", vnsAbsConnection.DistinguishedName)
	}

	return vnsAbsConnection, nil
}

func setConnectionAttributes(vnsAbsConnection *models.Connection, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	d.SetId(vnsAbsConnection.DistinguishedName)
	d.Set("description", vnsAbsConnection.Description)
	if dn != vnsAbsConnection.DistinguishedName {
		d.Set("l4_l7_service_graph_template_dn", "")
	}
	vnsAbsConnectionMap, err := vnsAbsConnection.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", vnsAbsConnectionMap["name"])

	d.Set("adj_type", vnsAbsConnectionMap["adjType"])
	d.Set("annotation", vnsAbsConnectionMap["annotation"])
	d.Set("conn_dir", vnsAbsConnectionMap["connDir"])
	d.Set("conn_type", vnsAbsConnectionMap["connType"])
	d.Set("direct_connect", vnsAbsConnectionMap["directConnect"])
	d.Set("name_alias", vnsAbsConnectionMap["nameAlias"])
	d.Set("unicast_route", vnsAbsConnectionMap["unicastRoute"])
	return d, nil
}

func resourceAciConnectionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vnsAbsConnection, err := getRemoteConnection(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setConnectionAttributes(vnsAbsConnection, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Connection: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L4L7ServiceGraphTemplateDn := d.Get("l4_l7_service_graph_template_dn").(string)

	vnsAbsConnectionAttr := models.ConnectionAttributes{}
	if AdjType, ok := d.GetOk("adj_type"); ok {
		vnsAbsConnectionAttr.AdjType = AdjType.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsAbsConnectionAttr.Annotation = Annotation.(string)
	} else {
		vnsAbsConnectionAttr.Annotation = "{}"
	}
	if ConnDir, ok := d.GetOk("conn_dir"); ok {
		vnsAbsConnectionAttr.ConnDir = ConnDir.(string)
	}
	if ConnType, ok := d.GetOk("conn_type"); ok {
		vnsAbsConnectionAttr.ConnType = ConnType.(string)
	}
	if DirectConnect, ok := d.GetOk("direct_connect"); ok {
		vnsAbsConnectionAttr.DirectConnect = DirectConnect.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsAbsConnectionAttr.NameAlias = NameAlias.(string)
	}
	if UnicastRoute, ok := d.GetOk("unicast_route"); ok {
		vnsAbsConnectionAttr.UnicastRoute = UnicastRoute.(string)
	}
	vnsAbsConnection := models.NewConnection(fmt.Sprintf("AbsConnection-%s", name), L4L7ServiceGraphTemplateDn, desc, vnsAbsConnectionAttr)

	err := aciClient.Save(vnsAbsConnection)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTovnsRsAbsCopyConnection, ok := d.GetOk("relation_vns_rs_abs_copy_connection"); ok {
		relationParamList := toStringList(relationTovnsRsAbsCopyConnection.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}
	if relationTovnsRsAbsConnectionConns, ok := d.GetOk("relation_vns_rs_abs_connection_conns"); ok {
		relationParamList := toStringList(relationTovnsRsAbsConnectionConns.(*schema.Set).List())
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

	if relationTovnsRsAbsCopyConnection, ok := d.GetOk("relation_vns_rs_abs_copy_connection"); ok {
		relationParamList := toStringList(relationTovnsRsAbsCopyConnection.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationvnsRsAbsCopyConnectionFromConnection(vnsAbsConnection.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}
	if relationTovnsRsAbsConnectionConns, ok := d.GetOk("relation_vns_rs_abs_connection_conns"); ok {
		relationParamList := toStringList(relationTovnsRsAbsConnectionConns.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationvnsRsAbsConnectionConnsFromConnection(vnsAbsConnection.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}

		}
	}

	d.SetId(vnsAbsConnection.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciConnectionRead(ctx, d, m)
}

func resourceAciConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Connection: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	L4L7ServiceGraphTemplateDn := d.Get("l4_l7_service_graph_template_dn").(string)

	vnsAbsConnectionAttr := models.ConnectionAttributes{}
	if AdjType, ok := d.GetOk("adj_type"); ok {
		vnsAbsConnectionAttr.AdjType = AdjType.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsAbsConnectionAttr.Annotation = Annotation.(string)
	} else {
		vnsAbsConnectionAttr.Annotation = "{}"
	}
	if ConnDir, ok := d.GetOk("conn_dir"); ok {
		vnsAbsConnectionAttr.ConnDir = ConnDir.(string)
	}
	if ConnType, ok := d.GetOk("conn_type"); ok {
		vnsAbsConnectionAttr.ConnType = ConnType.(string)
	}
	if DirectConnect, ok := d.GetOk("direct_connect"); ok {
		vnsAbsConnectionAttr.DirectConnect = DirectConnect.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsAbsConnectionAttr.NameAlias = NameAlias.(string)
	}
	if UnicastRoute, ok := d.GetOk("unicast_route"); ok {
		vnsAbsConnectionAttr.UnicastRoute = UnicastRoute.(string)
	}
	vnsAbsConnection := models.NewConnection(fmt.Sprintf("AbsConnection-%s", name), L4L7ServiceGraphTemplateDn, desc, vnsAbsConnectionAttr)

	vnsAbsConnection.Status = "modified"

	err := aciClient.Save(vnsAbsConnection)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vns_rs_abs_copy_connection") {
		oldRel, newRel := d.GetChange("relation_vns_rs_abs_copy_connection")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_vns_rs_abs_connection_conns") {
		oldRel, newRel := d.GetChange("relation_vns_rs_abs_connection_conns")
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

	if d.HasChange("relation_vns_rs_abs_copy_connection") {
		oldRel, newRel := d.GetChange("relation_vns_rs_abs_copy_connection")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationvnsRsAbsCopyConnectionFromConnection(vnsAbsConnection.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationvnsRsAbsCopyConnectionFromConnection(vnsAbsConnection.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_vns_rs_abs_connection_conns") {
		oldRel, newRel := d.GetChange("relation_vns_rs_abs_connection_conns")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationvnsRsAbsConnectionConnsFromConnection(vnsAbsConnection.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationvnsRsAbsConnectionConnsFromConnection(vnsAbsConnection.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}

	d.SetId(vnsAbsConnection.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciConnectionRead(ctx, d, m)

}

func resourceAciConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vnsAbsConnection, err := getRemoteConnection(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setConnectionAttributes(vnsAbsConnection, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	vnsRsAbsCopyConnectionData, err := aciClient.ReadRelationvnsRsAbsCopyConnectionFromConnection(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsAbsCopyConnection %v", err)
		d.Set("relation_vns_rs_abs_copy_connection", make([]string, 0, 1))

	} else {
		d.Set("relation_vns_rs_abs_copy_connection", toStringList(vnsRsAbsCopyConnectionData.(*schema.Set).List()))
	}

	vnsRsAbsConnectionConnsData, err := aciClient.ReadRelationvnsRsAbsConnectionConnsFromConnection(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsAbsConnectionConns %v", err)
		d.Set("relation_vns_rs_abs_connection_conns", make([]string, 0, 1))

	} else {
		d.Set("relation_vns_rs_abs_connection_conns", toStringList(vnsRsAbsConnectionConnsData.(*schema.Set).List()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vnsAbsConnection")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}

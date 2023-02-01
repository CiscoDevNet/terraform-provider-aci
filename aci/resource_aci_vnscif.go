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

func resourceAciConcreteInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciConcreteInterfaceCreate,
		UpdateContext: resourceAciConcreteInterfaceUpdate,
		ReadContext:   resourceAciConcreteInterfaceRead,
		DeleteContext: resourceAciConcreteInterfaceDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciConcreteInterfaceImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"concrete_device_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encap": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vnic_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_vns_rs_c_if_path_att": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to fabric:PathEp",
			},
		})),
	}
}

func getRemoteConcreteInterface(client *client.Client, dn string) (*models.ConcreteInterface, error) {
	vnsCIfCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vnsCIf := models.ConcreteInterfaceFromContainer(vnsCIfCont)
	if vnsCIf.DistinguishedName == "" {
		return nil, fmt.Errorf("Concrete Interface %s not found", dn)
	}
	return vnsCIf, nil
}

func setConcreteInterfaceAttributes(vnsCIf *models.ConcreteInterface, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(vnsCIf.DistinguishedName)
	vnsCIfMap, err := vnsCIf.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != vnsCIf.DistinguishedName {
		d.Set("concrete_device_dn", "")
	} else {
		d.Set("concrete_device_dn", GetParentDn(vnsCIf.DistinguishedName, fmt.Sprintf("/"+models.RnvnsCIf, vnsCIfMap["name"])))
	}
	d.Set("annotation", vnsCIfMap["annotation"])
	d.Set("encap", vnsCIfMap["encap"])
	d.Set("name", vnsCIfMap["name"])
	d.Set("vnic_name", vnsCIfMap["vnicName"])
	d.Set("name_alias", vnsCIfMap["nameAlias"])
	return d, nil
}

func resourceAciConcreteInterfaceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vnsCIf, err := getRemoteConcreteInterface(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setConcreteInterfaceAttributes(vnsCIf, d)
	if err != nil {
		return nil, err
	}
	vnsRsCIfPathAttData, err := aciClient.ReadRelationvnsRsCIfPathAtt(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsCIfPathAtt %v", err)
		d.Set("relation_vns_rs_c_if_path_att", "")
	} else {
		d.Set("relation_vns_rs_c_if_path_att", vnsRsCIfPathAttData.(string))
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciConcreteInterfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Concrete Interface: Beginning Creation")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	ConcreteDeviceDn := d.Get("concrete_device_dn").(string)

	vnsCIfAttr := models.ConcreteInterfaceAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsCIfAttr.Annotation = Annotation.(string)
	} else {
		vnsCIfAttr.Annotation = "{}"
	}

	if Encap, ok := d.GetOk("encap"); ok {
		vnsCIfAttr.Encap = Encap.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsCIfAttr.Name = Name.(string)
	}

	if VnicName, ok := d.GetOk("vnic_name"); ok {
		vnsCIfAttr.VnicName = VnicName.(string)
	}
	vnsCIf := models.NewConcreteInterface(fmt.Sprintf(models.RnvnsCIf, name), ConcreteDeviceDn, nameAlias, vnsCIfAttr)

	err := aciClient.Save(vnsCIf)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTovnsRsCIfPathAtt, ok := d.GetOk("relation_vns_rs_c_if_path_att"); ok {
		relationParam := relationTovnsRsCIfPathAtt.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovnsRsCIfPathAtt, ok := d.GetOk("relation_vns_rs_c_if_path_att"); ok {
		relationParam := relationTovnsRsCIfPathAtt.(string)
		err = aciClient.CreateRelationvnsRsCIfPathAtt(vnsCIf.DistinguishedName, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vnsCIf.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciConcreteInterfaceRead(ctx, d, m)
}
func resourceAciConcreteInterfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ConcreteInterface: Beginning Update")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	ConcreteDeviceDn := d.Get("concrete_device_dn").(string)

	vnsCIfAttr := models.ConcreteInterfaceAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsCIfAttr.Annotation = Annotation.(string)
	} else {
		vnsCIfAttr.Annotation = "{}"
	}

	if Encap, ok := d.GetOk("encap"); ok {
		vnsCIfAttr.Encap = Encap.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsCIfAttr.Name = Name.(string)
	}

	if VnicName, ok := d.GetOk("vnic_name"); ok {
		vnsCIfAttr.VnicName = VnicName.(string)
	}
	vnsCIf := models.NewConcreteInterface(fmt.Sprintf(models.RnvnsCIf, name), ConcreteDeviceDn, nameAlias, vnsCIfAttr)

	vnsCIf.Status = "modified"

	err := aciClient.Save(vnsCIf)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vns_rs_c_if_path_att") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vns_rs_c_if_path_att")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vns_rs_c_if_path_att") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vns_rs_c_if_path_att")
		err = aciClient.DeleteRelationvnsRsCIfPathAtt(vnsCIf.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsCIfPathAtt(vnsCIf.DistinguishedName, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(vnsCIf.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciConcreteInterfaceRead(ctx, d, m)
}

func resourceAciConcreteInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vnsCIf, err := getRemoteConcreteInterface(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setConcreteInterfaceAttributes(vnsCIf, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	vnsRsCIfPathAttData, err := aciClient.ReadRelationvnsRsCIfPathAtt(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsCIfPathAtt %v", err)
		setRelationAttribute(d, "relation_vns_rs_c_if_path_att", "")
	} else {
		setRelationAttribute(d, "relation_vns_rs_c_if_path_att", vnsRsCIfPathAttData.(string))
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciConcreteInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "vnsCIf")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

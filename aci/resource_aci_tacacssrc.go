package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciTACACSSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciTACACSSourceCreate,
		UpdateContext: resourceAciTACACSSourceUpdate,
		ReadContext:   resourceAciTACACSSourceRead,
		DeleteContext: resourceAciTACACSSourceDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTACACSSourceImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "uni/fabric/moncommon",
			},
			"incl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"audit",
						"events",
						"faults",
						"session",
					}, false),
				},
			},
			"min_sev": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"cleared",
					"critical",
					"info",
					"major",
					"minor",
					"warning",
				}, false),
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"relation_tacacs_rs_dest_group": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to tacacs:Group",
			}})),
	}
}

func getRemoteTACACSSource(client *client.Client, dn string) (*models.TACACSSource, error) {
	tacacsSrcCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	tacacsSrc := models.TACACSSourceFromContainer(tacacsSrcCont)
	if tacacsSrc.DistinguishedName == "" {
		return nil, fmt.Errorf("TACACSSource %s not found", tacacsSrc.DistinguishedName)
	}
	return tacacsSrc, nil
}

func setTACACSSourceAttributes(tacacsSrc *models.TACACSSource, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(tacacsSrc.DistinguishedName)
	d.Set("description", tacacsSrc.Description)
	tacacsSrcMap, err := tacacsSrc.ToMap()
	if err != nil {
		return d, err
	}
	inclGet := make([]string, 0, 1)
	for _, val := range strings.Split(tacacsSrcMap["incl"], ",") {
		inclGet = append(inclGet, strings.Trim(val, " "))
	}
	if len(inclGet) == 5 {
		inclGet = inclGet[1:]
	}
	sort.Strings(inclGet)
	if inclIntr, ok := d.GetOk("incl"); ok {
		inclAct := make([]string, 0, 1)
		for _, val := range inclIntr.([]interface{}) {
			inclAct = append(inclAct, val.(string))
		}
		sort.Strings(inclAct)
		if reflect.DeepEqual(inclAct, inclGet) {
			d.Set("incl", d.Get("incl").([]interface{}))
		} else {
			d.Set("incl", inclGet)
		}

	} else {
		d.Set("incl", inclGet)
	}
	d.Set("min_sev", tacacsSrcMap["minSev"])
	d.Set("name", tacacsSrcMap["name"])
	d.Set("name_alias", tacacsSrcMap["nameAlias"])
	d.Set("annotation", tacacsSrcMap["annotation"])
	return d, nil
}

func resourceAciTACACSSourceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	tacacsSrc, err := getRemoteTACACSSource(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setTACACSSourceAttributes(tacacsSrc, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTACACSSourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TACACSSource: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	ParentDn := d.Get("parent_dn").(string)

	tacacsSrcAttr := models.TACACSSourceAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		tacacsSrcAttr.Annotation = Annotation.(string)
	} else {
		tacacsSrcAttr.Annotation = "{}"
	}

	if Incl, ok := d.GetOk("incl"); ok {
		inclList := make([]string, 0, 1)
		for _, val := range Incl.([]interface{}) {
			inclList = append(inclList, val.(string))
		}
		if len(inclList) == 0 {
			tacacsSrcAttr.Incl = "none"
		} else {
			Incl := strings.Join(inclList, ",")
			tacacsSrcAttr.Incl = Incl
		}
	}

	if MinSev, ok := d.GetOk("min_sev"); ok {
		tacacsSrcAttr.MinSev = MinSev.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		tacacsSrcAttr.Name = Name.(string)
	}
	tacacsSrc := models.NewTACACSSource(fmt.Sprintf("tacacssrc-%s", name), ParentDn, desc, nameAlias, tacacsSrcAttr)
	err := aciClient.Save(tacacsSrc)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTotacacsRsDestGroup, ok := d.GetOk("relation_tacacs_rs_dest_group"); ok {
		relationParam := relationTotacacsRsDestGroup.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTotacacsRsDestGroup, ok := d.GetOk("relation_tacacs_rs_dest_group"); ok {
		relationParam := relationTotacacsRsDestGroup.(string)
		err = aciClient.CreateRelationtacacsRsDestGroup(tacacsSrc.DistinguishedName, tacacsSrcAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(tacacsSrc.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciTACACSSourceRead(ctx, d, m)
}

func resourceAciTACACSSourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TACACSSource: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	ParentDn := d.Get("parent_dn").(string)
	tacacsSrcAttr := models.TACACSSourceAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		tacacsSrcAttr.Annotation = Annotation.(string)
	} else {
		tacacsSrcAttr.Annotation = "{}"
	}
	if Incl, ok := d.GetOk("incl"); ok {
		inclList := make([]string, 0, 1)
		for _, val := range Incl.([]interface{}) {
			inclList = append(inclList, val.(string))
		}
		if len(inclList) == 0 {
			tacacsSrcAttr.Incl = "none"
		} else {
			Incl := strings.Join(inclList, ",")
			tacacsSrcAttr.Incl = Incl
		}
	} else {
		tacacsSrcAttr.Incl = "none"
	}

	if MinSev, ok := d.GetOk("min_sev"); ok {
		tacacsSrcAttr.MinSev = MinSev.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		tacacsSrcAttr.Name = Name.(string)
	}
	tacacsSrc := models.NewTACACSSource(fmt.Sprintf("tacacssrc-%s", name), ParentDn, desc, nameAlias, tacacsSrcAttr)
	tacacsSrc.Status = "modified"
	err := aciClient.Save(tacacsSrc)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_tacacs_rs_dest_group") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_tacacs_rs_dest_group")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_tacacs_rs_dest_group") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_tacacs_rs_dest_group")
		err = aciClient.DeleteRelationtacacsRsDestGroup(tacacsSrc.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationtacacsRsDestGroup(tacacsSrc.DistinguishedName, tacacsSrcAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(tacacsSrc.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciTACACSSourceRead(ctx, d, m)
}

func resourceAciTACACSSourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	tacacsSrc, err := getRemoteTACACSSource(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setTACACSSourceAttributes(tacacsSrc, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	tacacsRsDestGroupData, err := aciClient.ReadRelationtacacsRsDestGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation tacacsRsDestGroup %v", err)
		d.Set("relation_tacacs_rs_dest_group", "")
	} else {
		setRelationAttribute(d, "relation_tacacs_rs_dest_group", tacacsRsDestGroupData)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciTACACSSourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "tacacsSrc")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

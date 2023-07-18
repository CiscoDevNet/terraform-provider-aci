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

func resourceAciSnmpUserProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSnmpUserProfileCreate,
		UpdateContext: resourceAciSnmpUserProfileUpdate,
		ReadContext:   resourceAciSnmpUserProfileRead,
		DeleteContext: resourceAciSnmpUserProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSnmpUserProfileImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"snmp_policy_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"authorization_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"authorization_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"hmac-md5-96",
					"hmac-sha1-96",
					"hmac-sha2-224",
					"hmac-sha2-256",
					"hmac-sha2-384",
					"hmac-sha2-512",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"privacy_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"privacy_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"aes-128",
					"des",
					"none",
				}, false),
			},
		})),
	}
}

func getRemoteUserProfile(client *client.Client, dn string) (*models.SnmpUserProfile, error) {
	snmpUserPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	snmpUserP := models.SnmpUserProfileFromContainer(snmpUserPCont)
	if snmpUserP.DistinguishedName == "" {
		return nil, fmt.Errorf("SNMP User Profile %s not found", dn)
	}
	return snmpUserP, nil
}

func setUserProfileAttributes(snmpUserP *models.SnmpUserProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(snmpUserP.DistinguishedName)
	d.Set("description", snmpUserP.Description)
	snmpUserPMap, err := snmpUserP.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != snmpUserP.DistinguishedName {
		d.Set("snmp_policy_dn", "")
	} else {
		d.Set("snmp_policy_dn", GetParentDn(snmpUserP.DistinguishedName, fmt.Sprintf("/"+models.RnSnmpUserP, snmpUserPMap["name"])))
	}
	d.Set("annotation", snmpUserPMap["annotation"])
	authKey := d.Get("authorization_key").(string)
	if authKey != "" {
		d.Set("authorization_key", snmpUserPMap["authKey"])
	}
	d.Set("authorization_type", snmpUserPMap["authType"])
	d.Set("name", snmpUserPMap["name"])
	d.Set("name_alias", snmpUserPMap["nameAlias"])
	privKey := d.Get("privacy_key").(string)
	if privKey != "" {
		d.Set("privacy_key", snmpUserPMap["privKey"])
	}
	d.Set("privacy_type", snmpUserPMap["privType"])
	return d, nil
}

func resourceAciSnmpUserProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	snmpUserP, err := getRemoteUserProfile(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setUserProfileAttributes(snmpUserP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSnmpUserProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SNMP User Profile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	SNMPPolicyDn := d.Get("snmp_policy_dn").(string)

	snmpUserPAttr := models.SnmpUserProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		snmpUserPAttr.Annotation = Annotation.(string)
	} else {
		snmpUserPAttr.Annotation = "{}"
	}

	if AuthKey, ok := d.GetOk("authorization_key"); ok {
		snmpUserPAttr.AuthKey = AuthKey.(string)
	}

	if AuthType, ok := d.GetOk("authorization_type"); ok {
		snmpUserPAttr.AuthType = AuthType.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		snmpUserPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		snmpUserPAttr.NameAlias = NameAlias.(string)
	}

	if PrivKey, ok := d.GetOk("privacy_key"); ok {
		snmpUserPAttr.PrivKey = PrivKey.(string)
	}

	if PrivType, ok := d.GetOk("privacy_type"); ok {
		snmpUserPAttr.PrivType = PrivType.(string)
	}
	snmpUserP := models.NewSnmpUserProfile(fmt.Sprintf(models.RnSnmpUserP, name), SNMPPolicyDn, desc, snmpUserPAttr)

	err := aciClient.Save(snmpUserP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(snmpUserP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciSnmpUserProfileRead(ctx, d, m)
}
func resourceAciSnmpUserProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SNMP User Profile: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	SNMPPolicyDn := d.Get("snmp_policy_dn").(string)

	snmpUserPAttr := models.SnmpUserProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		snmpUserPAttr.Annotation = Annotation.(string)
	} else {
		snmpUserPAttr.Annotation = "{}"
	}

	if AuthKey, ok := d.GetOk("authorization_key"); ok {
		snmpUserPAttr.AuthKey = AuthKey.(string)
	}

	if AuthType, ok := d.GetOk("authorization_type"); ok {
		snmpUserPAttr.AuthType = AuthType.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		snmpUserPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		snmpUserPAttr.NameAlias = NameAlias.(string)
	}

	if PrivKey, ok := d.GetOk("privacy_key"); ok {
		snmpUserPAttr.PrivKey = PrivKey.(string)
	}

	if PrivType, ok := d.GetOk("privacy_type"); ok {
		snmpUserPAttr.PrivType = PrivType.(string)
	}
	snmpUserP := models.NewSnmpUserProfile(fmt.Sprintf(models.RnSnmpUserP, name), SNMPPolicyDn, desc, snmpUserPAttr)

	snmpUserP.Status = "modified"

	err := aciClient.Save(snmpUserP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(snmpUserP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciSnmpUserProfileRead(ctx, d, m)
}

func resourceAciSnmpUserProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	snmpUserP, err := getRemoteUserProfile(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setUserProfileAttributes(snmpUserP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciSnmpUserProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, models.SnmpUserPClassName)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}

package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccProviderVersionName = "tf_acc_apic_version"

func TestAccProviderAPICVersion(t *testing.T) {
	var version string
	versionVariables := config.Variables{}
	tenantDN := fmt.Sprintf("uni/tn-%s", testAccProviderVersionName)
	applicationDN := fmt.Sprintf("%s/ap-%s", tenantDN, testAccProviderVersionName)
	epgDN := fmt.Sprintf("%s/epg-%s", applicationDN, testAccProviderVersionName)
	restManagedEPGDN := fmt.Sprintf("%s_rest", epgDN)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			version = testAccProviderVersionPreCheck(t)
			versionVariables["apic_version"] = config.StringVariable(version)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					fmt.Println("Reject an invalid APIC version")
				},
				Config:      testAccProviderVersionConfig("invalid-version"),
				ExpectError: regexp.MustCompile("Invalid APIC version"),
			},
			{
				PreConfig: func() {
					fmt.Println("Reject extra text around an APIC version")
				},
				// Check that surrounding text is rejected as well.
				Config:      testAccProviderVersionConfig("release 6.0(1h)"),
				ExpectError: regexp.MustCompile("Invalid APIC version"),
			},
			{
				PreConfig: func() {
					fmt.Println("Fail version discovery for the tenant-scoped user without an override")
				},
				// Discovery must fail for this tenant-scoped user.
				Config: testAccProviderVersionConfig(""),
				ExpectError: regexp.MustCompile(
					"Data for topSystem class could not be retrieved",
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Read the tenant with an explicit APIC version")
				},
				Config:          testAccProviderVersionVariableConfig(),
				ConfigVariables: versionVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.aci_tenant.test", "name", testAccProviderVersionName),
					resource.TestCheckResourceAttr("data.aci_tenant.test", "id", tenantDN),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Read the tenant using ACI_APIC_VERSION")
					t.Setenv("ACI_APIC_VERSION", version)
				},
				// Omit the attribute to exercise the environment fallback.
				Config: testAccProviderVersionEnvironmentConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.aci_tenant.test", "name", testAccProviderVersionName),
					resource.TestCheckResourceAttr("data.aci_tenant.test", "id", tenantDN),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Read rest_managed EPG children with the supported APIC version")
				},
				Config:          testAccProviderVersionRestManagedConfig(),
				ConfigVariables: versionVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.aci_tenant.test", "name", testAccProviderVersionName),
					resource.TestCheckResourceAttr("data.aci_tenant.test", "id", tenantDN),
					resource.TestCheckResourceAttr("aci_application_profile.test", "name", testAccProviderVersionName),
					resource.TestCheckResourceAttr("aci_application_profile.test", "id", applicationDN),
					resource.TestCheckResourceAttr("aci_application_profile.test", "parent_dn", tenantDN),
					resource.TestCheckResourceAttr("aci_rest_managed.epg", "id", restManagedEPGDN),
					resource.TestCheckResourceAttr("aci_rest_managed.epg", "dn", restManagedEPGDN),
					resource.TestCheckResourceAttr("aci_rest_managed.epg", "class_name", "fvAEPg"),
					resource.TestCheckResourceAttr("aci_rest_managed.epg", "content.name", fmt.Sprintf("%s_rest", testAccProviderVersionName)),
					resource.TestCheckResourceAttr("aci_rest_managed.epg", "child.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("aci_rest_managed.epg", "child.*", map[string]string{
						"class_name":    "tagAnnotation",
						"rn":            "annotationKey-[version-test]",
						"content.key":   "version-test",
						"content.value": "test-value",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("aci_rest_managed.epg", "child.*", map[string]string{
						"class_name":    "tagTag",
						"rn":            "tagKey-version-test",
						"content.key":   "version-test",
						"content.value": "test-value",
					}),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Read rest_managed EPG children with the older APIC version")
				},
				// rest_managed reads configured children without version filtering.
				Config: testAccProviderVersionRestManagedConfig(),
				ConfigVariables: config.Variables{
					"apic_version": config.StringVariable("3.1(1i)"),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.aci_tenant.test", "name", testAccProviderVersionName),
					resource.TestCheckResourceAttr("data.aci_tenant.test", "id", tenantDN),
					resource.TestCheckResourceAttr("aci_application_profile.test", "name", testAccProviderVersionName),
					resource.TestCheckResourceAttr("aci_application_profile.test", "id", applicationDN),
					resource.TestCheckResourceAttr("aci_application_profile.test", "parent_dn", tenantDN),
					resource.TestCheckResourceAttr("aci_rest_managed.epg", "id", restManagedEPGDN),
					resource.TestCheckResourceAttr("aci_rest_managed.epg", "dn", restManagedEPGDN),
					resource.TestCheckResourceAttr("aci_rest_managed.epg", "class_name", "fvAEPg"),
					resource.TestCheckResourceAttr("aci_rest_managed.epg", "content.name", fmt.Sprintf("%s_rest", testAccProviderVersionName)),
					resource.TestCheckResourceAttr("aci_rest_managed.epg", "child.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("aci_rest_managed.epg", "child.*", map[string]string{
						"class_name":    "tagAnnotation",
						"rn":            "annotationKey-[version-test]",
						"content.key":   "version-test",
						"content.value": "test-value",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("aci_rest_managed.epg", "child.*", map[string]string{
						"class_name":    "tagTag",
						"rn":            "tagKey-version-test",
						"content.key":   "version-test",
						"content.value": "test-value",
					}),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Read generated EPG resource children with the supported APIC version")
				},
				Config:          testAccProviderVersionResourcesConfig(),
				ConfigVariables: versionVariables,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.aci_tenant.test", "name", testAccProviderVersionName),
					resource.TestCheckResourceAttr("data.aci_tenant.test", "id", tenantDN),
					resource.TestCheckResourceAttr("aci_application_profile.test", "name", testAccProviderVersionName),
					resource.TestCheckResourceAttr("aci_application_profile.test", "id", applicationDN),
					resource.TestCheckResourceAttr("aci_application_profile.test", "parent_dn", tenantDN),
					resource.TestCheckResourceAttr("aci_application_epg.test", "name", testAccProviderVersionName),
					resource.TestCheckResourceAttr("aci_application_epg.test", "id", epgDN),
					resource.TestCheckResourceAttr("aci_application_epg.test", "parent_dn", applicationDN),
					resource.TestCheckResourceAttr("aci_application_epg.test", "annotations.#", "1"),
					resource.TestCheckResourceAttr("aci_application_epg.test", "annotations.0.key", "version-test"),
					resource.TestCheckResourceAttr("aci_application_epg.test", "annotations.0.value", "test-value"),
					resource.TestCheckResourceAttr("aci_application_epg.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("aci_application_epg.test", "tags.0.key", "version-test"),
					resource.TestCheckResourceAttr("aci_application_epg.test", "tags.0.value", "test-value"),
				),
			},
			{
				PreConfig: func() {
					fmt.Println("Expect inconsistent EPG state when the older version filters out its children")
				},
				Config: testAccProviderVersionResourcesConfig(),
				ConfigVariables: config.Variables{
					"apic_version": config.StringVariable("3.1(1i)"),
				},
				// APIC rejects an unknown rsp-subtree-class with HTTP 400/code 12
				// ("Unknown class ..."). Here the real APIC supports annotations and
				// tags (introduced in 3.2(1l)), but the older override omits them from
				// the GET, so Terraform reports inconsistent state after apply.
				ExpectError: regexp.MustCompile(
					`(?s)Provider produced inconsistent result after apply.*(annotations|tags)`,
				),
			},
		},
	})
}

func testAccProviderVersionPreCheck(t *testing.T) string {
	t.Helper()
	admin := getACIClientTest(t)
	var diagnostics diag.Diagnostics
	version := getVersionAPIC(context.Background(), &diagnostics, admin)
	if diagnostics.HasError() {
		diagnostic := diagnostics.Errors()[0]
		t.Fatalf("retrieve APIC version: %s: %s", diagnostic.Summary(), diagnostic.Detail())
	}
	if version == "" {
		t.Fatal("APIC returned no controller version")
	}
	supported, err := CompareVersionsRange(version, "3.2(1l)-", "inside")
	if err != nil {
		t.Fatalf("check annotation and tag support: %v", err)
	}
	if !supported {
		t.Skipf("annotations and tags require APIC 3.2(1l) or later; got %s", version)
	}

	password := "V7!qR2#nL9@xT4"
	testAccCreateProviderVersionFixtures(t, admin, password)

	previousVersion := apicVersion
	t.Cleanup(func() { apicVersion = previousVersion })
	t.Setenv("ACI_USERNAME", testAccProviderVersionName)
	t.Setenv("ACI_PASSWORD", password)
	t.Setenv("ACI_APIC_VERSION", "")

	tenantDN := fmt.Sprintf("uni/tn-%s", testAccProviderVersionName)
	restricted := testAccProviderVersionClient(testAccProviderVersionName, password)
	if _, err := restricted.Get(tenantDN); err != nil {
		t.Fatalf("restricted user cannot read the test tenant: %v", err)
	}

	return version
}

func testAccProviderVersionClient(username, password string) *client.Client {
	return client.NewClient(
		os.Getenv("ACI_URL"), username,
		client.Password(password),
		client.Insecure(true),
		client.ValidateRelationDn(false),
		client.ProxyUrl(os.Getenv("ACI_PROXY_URL")),
		client.ProxyCreds(os.Getenv("ACI_PROXY_CREDS")),
		client.ReqTimeout(10),
		client.MaxRetries(0),
		client.SkipLoggingPayload(true),
	)
}

func testAccProviderVersionConfig(version string) string {
	return fmt.Sprintf(`
provider "aci" {
  apic_version = %q
}

data "aci_tenant" "test" {
  name = %q
}
`, version, testAccProviderVersionName)
}

func testAccProviderVersionVariableConfig() string {
	return fmt.Sprintf(`
variable "apic_version" {
  type = string
}

provider "aci" {
  apic_version = var.apic_version
}

data "aci_tenant" "test" {
  name = %q
}
`, testAccProviderVersionName)
}

func testAccProviderVersionEnvironmentConfig() string {
	return fmt.Sprintf(`
provider "aci" {}

data "aci_tenant" "test" {
  name = %q
}
`, testAccProviderVersionName)
}

func testAccProviderVersionResourcesConfig() string {
	return fmt.Sprintf(`
variable "apic_version" {
  type = string
}

provider "aci" {
  apic_version = var.apic_version
}

data "aci_tenant" "test" {
  name = %[1]q
}

resource "aci_application_profile" "test" {
  parent_dn = data.aci_tenant.test.id
  name      = %[1]q
}

resource "aci_application_epg" "test" {
  parent_dn = aci_application_profile.test.id
  name      = %[1]q
  annotations = [{
    key   = "version-test"
    value = "test-value"
  }]
  tags = [{
    key   = "version-test"
    value = "test-value"
  }]
}
`, testAccProviderVersionName)
}

func testAccProviderVersionRestManagedConfig() string {
	return fmt.Sprintf(`
variable "apic_version" {
  type = string
}

provider "aci" {
  apic_version = var.apic_version
}

data "aci_tenant" "test" {
  name = %[1]q
}

resource "aci_application_profile" "test" {
  parent_dn = data.aci_tenant.test.id
  name      = %[1]q
}

resource "aci_rest_managed" "epg" {
  dn         = "${aci_application_profile.test.id}/epg-%[1]s_rest"
  class_name = "fvAEPg"
  content = {
    name = "%[1]s_rest"
  }

  child {
    rn         = "annotationKey-[version-test]"
    class_name = "tagAnnotation"
    content = {
      key   = "version-test"
      value = "test-value"
    }
  }

  child {
    rn         = "tagKey-version-test"
    class_name = "tagTag"
    content = {
      key   = "version-test"
      value = "test-value"
    }
  }
}
`, testAccProviderVersionName)
}

func testAccCreateProviderVersionFixtures(t *testing.T, admin *client.Client, password string) {
	t.Helper()
	name := testAccProviderVersionName
	domainDN := fmt.Sprintf("uni/userext/domain-%s", name)
	tenantDN := fmt.Sprintf("uni/tn-%s", name)
	userDN := fmt.Sprintf("uni/userext/user-%s", name)

	// Fixed fixture names must not overwrite existing objects.
	for _, dn := range []string{domainDN, tenantDN, userDN} {
		response := testAccProviderVersionGet(t, admin, dn)
		count, err := response.ArrayCount("imdata")
		if err != nil {
			t.Fatalf("inspect fixture %s: %v", dn, err)
		}
		if count != 0 {
			t.Fatalf("fixture %s already exists", dn)
		}
	}

	_, err := admin.CreateSecurityDomain(name, "APIC version acceptance test", models.SecurityDomainAttributes{Name: name})
	if err != nil {
		t.Fatalf("create security domain: %v", err)
	}

	testAccProviderVersionCleanup(t, admin, domainDN, func() error {
		return admin.DeleteSecurityDomain(name)
	})

	_, err = admin.CreateTenant(name, "APIC version acceptance test", models.TenantAttributes{Name: name})
	if err != nil {
		t.Fatalf("create tenant: %v", err)
	}

	testAccProviderVersionCleanup(t, admin, tenantDN, func() error {
		return admin.DeleteTenant(name)
	})

	_, err = admin.CreateAaaDomainRef(name, tenantDN, "", models.AaaDomainRefAttributes{Name: name})
	if err != nil {
		t.Fatalf("assign tenant security domain: %v", err)
	}

	_, err = admin.CreateLocalUser(name, "APIC version acceptance test",
		models.LocalUserAttributes{
			Name:              name,
			Pwd:               password,
			AccountStatus:     "active",
			Expires:           "no",
			PwdUpdateRequired: "no",
		})
	if err != nil {
		t.Fatalf("create local user: %v", err)
	}

	testAccProviderVersionCleanup(t, admin, userDN, func() error {
		return admin.DeleteLocalUser(name)
	})

	_, err = admin.CreateUserDomain(name, name, "", "", models.UserDomainAttributes{Name: name})
	if err != nil {
		t.Fatalf("assign user security domain: %v", err)
	}

	_, err = admin.CreateUserRole("tenant-admin", name, name, "", "", models.UserRoleAttributes{Name: "tenant-admin", PrivType: "writePriv"})
	if err != nil {
		t.Fatalf("assign tenant role: %v", err)
	}
}

func testAccProviderVersionGet(t *testing.T, aciClient *client.Client, dn string) *container.Container {
	t.Helper()
	var diagnostics diag.Diagnostics
	response := DoRestRequest(context.Background(), &diagnostics, aciClient, fmt.Sprintf("api/mo/%s.json", dn), "GET", nil)
	if diagnostics.HasError() {
		diagnostic := diagnostics.Errors()[0]
		t.Fatalf("read %s: %s: %s", dn,
			diagnostic.Summary(), diagnostic.Detail())
	}
	if response == nil {
		t.Fatalf("read %s returned no response", dn)
	}
	return response
}

func testAccProviderVersionCleanup(t *testing.T, admin *client.Client, dn string, remove func() error) {
	t.Helper()
	t.Cleanup(func() {
		if err := remove(); err != nil {
			t.Errorf("delete fixture %s: %v", dn, err)
			return
		}
		var diagnostics diag.Diagnostics
		response := DoRestRequest(context.Background(), &diagnostics, admin, fmt.Sprintf("api/mo/%s.json", dn), "GET", nil)
		if diagnostics.HasError() {
			diagnostic := diagnostics.Errors()[0]
			t.Errorf("verify cleanup of %s: %s: %s", dn,
				diagnostic.Summary(), diagnostic.Detail())
			return
		}
		if response == nil {
			t.Errorf("verify cleanup of %s: no response", dn)
			return
		}
		count, err := response.ArrayCount("imdata")
		if err != nil || count != 0 {
			t.Errorf("fixture %s was not confirmed deleted", dn)
		}
	})
}

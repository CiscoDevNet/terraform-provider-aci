package acctest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciContractSubject_Basic(t *testing.T) {
	var contract_subject_default models.ContractSubject
	var contract_subject_updated models.ContractSubject
	resourceName := "aci_contract_subject.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOtherName := makeTestVariable(acctest.RandString(5))
	parentOtherName := makeTestVariable(acctest.RandString(5))
	longerName := acctest.RandString(65)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccContractSubjectWithoutContract(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccContractSubjectWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccContractSubjectConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_default),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "cons_match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "prov_match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "rev_flt_ports", "yes"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "contract_dn", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_sdwan_pol", ""),
					// resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_filt_att.#", "0"), giving null on initial apply
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_graph_att", ""),
				),
			},
			{
				Config: CreateAccContractSubjectConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "test_annotation"),
					resource.TestCheckResourceAttr(resourceName, "cons_match_t", "All"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_name_alias"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "prov_match_t", "All"),
					resource.TestCheckResourceAttr(resourceName, "rev_flt_ports", "no"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS0"),
					resource.TestCheckResourceAttr(resourceName, "contract_dn", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_sdwan_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_filt_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_graph_att", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccContractSubjectUpdatedName(rName, longerName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of subj-%s failed validation for value '%s'", longerName, longerName)),
			},
			{
				Config: CreateAccContractSubjectConfigWithParentAndName(rName, rOtherName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "contract_dn", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rOtherName),
					testAccCheckAciContractSubjectIdNotEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectConfig(rName),
			},
			{
				Config: CreateAccContractSubjectConfigWithParentAndName(parentOtherName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "contract_dn", fmt.Sprintf("uni/tn-%s/brc-%s", parentOtherName, parentOtherName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciContractSubjectIdNotEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
		},
	})
}

func TestAccAciContractSubject_Update(t *testing.T) {
	var contract_subject_default models.ContractSubject
	var contract_subject_updated models.ContractSubject
	resourceName := "aci_contract_subject.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccContractSubjectConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_default),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "cons_match_t", "None"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "cons_match_t", "None"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "cons_match_t", "AtmostOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "cons_match_t", "AtmostOne"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "prov_match_t", "None"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "prov_match_t", "None"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "prov_match_t", "AtmostOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "prov_match_t", "AtmostOne"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "CS1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS1"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF11"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF11"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF12"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF12"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF13"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF13"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "CS2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS2"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF21"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF21"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF22"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF22"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF23"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF23"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "CS3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS3"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF31"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF31"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF32"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF33"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF33"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "CS4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS4"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF41"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF41"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF42"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF42"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "AF43"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF43"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "CS5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS5"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "VA"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "VA"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "EF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "EF"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "CS6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS6"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
			{
				Config: CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", "CS7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS7"),
					testAccCheckAciContractSubjectIdEqual(&contract_subject_default, &contract_subject_updated),
				),
			},
		},
	})
}

func TestAccAciContractSubject_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longAnnotationDesc := acctest.RandString(129)
	longNameAlias := acctest.RandString(65)
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccContractSubjectConfig(rName),
			},
			{
				Config:      CreateAccContractSubjectWithInvalidContract(rName),
				ExpectError: regexp.MustCompile(`Invalid request. dn '(.)+' is not valid for class vzSubj`),
			},
			{
				Config:      CreateAccContractSubjectUpdatedAttr(rName, "description", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccContractSubjectUpdatedAttr(rName, "annotation", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccContractSubjectUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccContractSubjectUpdatedAttr(rName, "cons_match_t", randomValue),
				ExpectError: regexp.MustCompile(`expected cons_match_t to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractSubjectUpdatedAttr(rName, "prov_match_t", randomValue),
				ExpectError: regexp.MustCompile(`expected prov_match_t to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractSubjectUpdatedAttr(rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`expected prio to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractSubjectUpdatedAttr(rName, "rev_flt_ports", randomValue),
				ExpectError: regexp.MustCompile(`expected rev_flt_ports to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractSubjectUpdatedAttr(rName, "target_dscp", randomValue),
				ExpectError: regexp.MustCompile(`expected target_dscp to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractSubjectUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccContractSubjectConfig(rName),
			},
		},
	})
}

func TestAccAciContractSubject_RelationParameters(t *testing.T) {
	var contract_subject_default models.ContractSubject
	var contract_subject_rel1 models.ContractSubject
	var contract_subject_rel2 models.ContractSubject
	resourceName := "aci_contract_subject.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomName1 := makeTestVariable(acctest.RandString(5))
	randomName2 := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccContractSubjectConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_default),
				),
			},
			{
				Config: CreateAccContractSubjectConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_default),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_sdwan_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_graph_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_filt_att.#", "0"),
				),
			},
			{
				Config: CreateAccContractSubjectRelConfigInitial(rName, randomName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_rel1),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_graph_att", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_filt_att.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_subj_filt_att.*", fmt.Sprintf("uni/tn-%s/flt-%s", rName, randomName1)),
				),
			},
			{
				Config: CreateAccContractSubjectRelConfigFinal(rName, randomName1, randomName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_rel2),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_graph_att", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_filt_att.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_subj_filt_att.*", fmt.Sprintf("uni/tn-%s/flt-%s", rName, randomName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_vz_rs_subj_filt_att.*", fmt.Sprintf("uni/tn-%s/flt-%s", rName, randomName2)),
				),
			},
			{
				Config: CreateAccContractSubjectConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractSubjectExists(resourceName, &contract_subject_default),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_sdwan_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_graph_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_subj_filt_att.#", "0"),
				),
			},
		},
	})
}

func TestAccContractSubject_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccContractSubjectsConfig(rName),
			},
		},
	})
}

func CreateAccContractSubjectWithoutContract(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract subject without creating contract")
	resource := fmt.Sprintf(`
	resource "aci_contract_subject" "test" {
		name = "%s"
	}
	`, rName)
	return resource

}

func CreateAccContractSubjectWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract subject without passing name attribute")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
	}
	`, rName, rName)
	return resource
}

func CreateAccContractSubjectConfig(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract subject creation with required paramters only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccContractSubjectConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract subject creation with optional paramters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
		description = "test_description"
		annotation = "test_annotation"
		cons_match_t = "All"
		name_alias = "test_name_alias"
		prio = "level1"
		prov_match_t = "All"
		rev_flt_ports = "no"
		target_dscp = "CS0"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccContractSubjectUpdatedName(rName, longerName string) string {
	fmt.Println("=== STEP  Basic: testing contract subject creation with invalid name with long lenght")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
	}
	`, rName, rName, longerName)
	return resource
}

func CreateAccContractSubjectConfigWithParentAndName(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing contract subject creation with contract name %s and contract subject name %s\n", parentName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
	}
	`, parentName, parentName, rName)
	return resource
}

func CreateAccContractSubjectUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing contract subject %s = %s\n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
		%s = "%s"
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccContractSubjectWithInvalidContract(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract subject updation with Invalid contract_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccContractSubjectRelConfigInitial(rName, relName1 string) string {
	fmt.Printf("=== STEP  Basic: testing contract subject relation")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
		relation_vz_rs_subj_graph_att = aci_l4_l7_service_graph_template.test.id
		relation_vz_rs_subj_filt_att = [aci_filter.test.id]
	}

	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName, rName, relName1, relName1)
	return resource
}

func CreateAccContractSubjectRelConfigFinal(rName, relName1, relName2 string) string {
	fmt.Printf("=== STEP  Basic: testing contract subject relations")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
		relation_vz_rs_subj_graph_att = aci_l4_l7_service_graph_template.test.id
		relation_vz_rs_subj_filt_att = [aci_filter.test.id,aci_filter.test1.id]
	}

	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter" "test1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName, rName, relName2, relName1, relName2)
	return resource
}

func CreateAccContractSubjectsConfig(rName string) string {
	fmt.Println("=== STEP  creating multiple Contract Subjects")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test1"{
		contract_dn = aci_contract.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test2"{
		contract_dn = aci_contract.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test3"{
		contract_dn = aci_contract.test.id
		name = "%s"
	}

	`, rName, rName, rName, rName+"1", rName+"2", rName+"3")
	return resource

}

func testAccCheckAciContractSubjectDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing Contract Subject destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_contract_subject" {
			cont, err := client.Get(rs.Primary.ID)
			contract_subject := models.ContractSubjectFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Contract Subject %s still exists", contract_subject.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciContractSubjectExists(name string, contract_subject *models.ContractSubject) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Contract Subject %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Contract Subject Dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		contract_subjectFound := models.ContractSubjectFromContainer(cont)
		if contract_subjectFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Contract Subject %s not found", rs.Primary.ID)
		}
		*contract_subject = *contract_subjectFound
		return nil
	}
}

func testAccCheckAciContractSubjectIdNotEqual(cs1, cs2 *models.ContractSubject) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cs1.DistinguishedName == cs2.DistinguishedName {
			return fmt.Errorf("Contract Subject DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciContractSubjectIdEqual(cs1, cs2 *models.ContractSubject) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cs1.DistinguishedName != cs2.DistinguishedName {
			return fmt.Errorf("Contract Subject DNs are no equal")
		}
		return nil
	}
}

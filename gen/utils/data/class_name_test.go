package data

import (
	"fmt"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

type splitClassNameExpected struct {
	PackageName string
	ShortName   string
	Error       bool
	ErrorMsg    string
}

func TestSplitClassNameToPackageNameAndShortName(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:  "test_single_word",
			Input: "fvTenant",
			Expected: splitClassNameExpected{
				PackageName: "fv",
				ShortName:   "Tenant",
			},
		},
		{
			Name:  "test_single_word_meta_style_class_name",
			Input: "fv:Tenant",
			Expected: splitClassNameExpected{
				PackageName: "fv",
				ShortName:   "Tenant",
			},
		},
		{
			Name:  "test_multiple_words",
			Input: "fvRsIpslaMonPol",
			Expected: splitClassNameExpected{
				PackageName: "fv",
				ShortName:   "RsIpslaMonPol",
			},
		},
		{
			Name:  "test_multiple_words_meta_style_class_name",
			Input: "fv:RsIpslaMonPol",
			Expected: splitClassNameExpected{
				PackageName: "fv",
				ShortName:   "RsIpslaMonPol",
			},
		},
		{
			Name:  "test_longer_package_name",
			Input: "netflowAExporterPol",
			Expected: splitClassNameExpected{
				PackageName: "netflow",
				ShortName:   "AExporterPol",
			},
		},
		{
			Name:  "test_longer_package_meta_style_class_name",
			Input: "netflow:AExporterPol",
			Expected: splitClassNameExpected{
				PackageName: "netflow",
				ShortName:   "AExporterPol",
			},
		},
		{
			Name:  "test_error_no_uppercase",
			Input: "error",
			Expected: splitClassNameExpected{
				Error:    true,
				ErrorMsg: "failed to split class name 'error' for name space separation",
			},
		},
		{
			Name:  "test_empty_string",
			Input: "",
			Expected: splitClassNameExpected{
				Error:    true,
				ErrorMsg: "failed to split class name '' for name space separation",
			},
		},
		{
			Name:  "test_error_multiple_colons",
			Input: "fv:Rs:Bd",
			Expected: splitClassNameExpected{
				Error:    true,
				ErrorMsg: "invalid class name 'fv:Rs:Bd': multiple colons detected",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			expected := testCase.Expected.(splitClassNameExpected)
			packageName, shortName, err := splitClassNameToPackageNameAndShortName(testCase.Input.(string))

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			}
			assert.Equal(t, expected.PackageName, packageName, test.MessageEqual(expected.PackageName, packageName, testCase.Name))
			assert.Equal(t, expected.ShortName, shortName, test.MessageEqual(expected.ShortName, shortName, testCase.Name))
		})
	}
}

type sanitizeClassNameExpected struct {
	Result   string
	Error    bool
	ErrorMsg string
}

func TestSanitizeClassName(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:  "test_standard_class_name",
			Input: "fv:Tenant",
			Expected: sanitizeClassNameExpected{
				Result: "fvTenant",
			},
		},
		{
			Name:  "test_relation_class_name",
			Input: "fv:RsBdToOut",
			Expected: sanitizeClassNameExpected{
				Result: "fvRsBdToOut",
			},
		},
		{
			Name:  "test_longer_package_name",
			Input: "netflow:AExporterPol",
			Expected: sanitizeClassNameExpected{
				Result: "netflowAExporterPol",
			},
		},
		{
			Name:  "test_already_sanitized",
			Input: "fvTenant",
			Expected: sanitizeClassNameExpected{
				Result: "fvTenant",
			},
		},
		{
			Name:  "test_error_multiple_colons",
			Input: "fv:Rs:Bd",
			Expected: sanitizeClassNameExpected{
				Error:    true,
				ErrorMsg: "invalid class name 'fv:Rs:Bd': multiple colons detected",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			expected := testCase.Expected.(sanitizeClassNameExpected)
			result, err := sanitizeClassName(testCase.Input.(string))

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else {
				assert.Equal(t, expected.Result, result, test.MessageEqual(expected.Result, result, testCase.Name))
			}
		})
	}
}

type newClassNameExpected struct {
	String      string
	Capitalized string
	Package     string
	Short       string
	MetaStyle   string
	Error       bool
	ErrorMsg    string
}

func TestNewClassName(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:  "test_full_class_name",
			Input: "fvTenant",
			Expected: newClassNameExpected{
				String:      "fvTenant",
				Capitalized: "FvTenant",
				Package:     "fv",
				Short:       "Tenant",
				MetaStyle:   "fv:Tenant",
			},
		},
		{
			Name:  "test_meta_style_class_name",
			Input: "fv:Tenant",
			Expected: newClassNameExpected{
				String:      "fvTenant",
				Capitalized: "FvTenant",
				Package:     "fv",
				Short:       "Tenant",
				MetaStyle:   "fv:Tenant",
			},
		},
		{
			Name:  "test_longer_package_name",
			Input: "netflowAExporterPol",
			Expected: newClassNameExpected{
				String:      "netflowAExporterPol",
				Capitalized: "NetflowAExporterPol",
				Package:     "netflow",
				Short:       "AExporterPol",
				MetaStyle:   "netflow:AExporterPol",
			},
		},
		{
			Name:  "test_meta_style_longer_package",
			Input: "netflow:AExporterPol",
			Expected: newClassNameExpected{
				String:      "netflowAExporterPol",
				Capitalized: "NetflowAExporterPol",
				Package:     "netflow",
				Short:       "AExporterPol",
				MetaStyle:   "netflow:AExporterPol",
			},
		},
		{
			Name:  "test_relation_class_name",
			Input: "fvRsBdToOut",
			Expected: newClassNameExpected{
				String:      "fvRsBdToOut",
				Capitalized: "FvRsBdToOut",
				Package:     "fv",
				Short:       "RsBdToOut",
				MetaStyle:   "fv:RsBdToOut",
			},
		},
		{
			Name:  "test_error_no_uppercase",
			Input: "error",
			Expected: newClassNameExpected{
				Error:    true,
				ErrorMsg: "failed to split class name 'error' for name space separation",
			},
		},
		{
			Name:  "test_error_empty_string",
			Input: "",
			Expected: newClassNameExpected{
				Error:    true,
				ErrorMsg: "failed to split class name '' for name space separation",
			},
		},
		{
			Name:  "test_error_multiple_colons",
			Input: "fv:Rs:Bd",
			Expected: newClassNameExpected{
				Error:    true,
				ErrorMsg: "invalid class name 'fv:Rs:Bd': multiple colons detected",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			expected := testCase.Expected.(newClassNameExpected)
			className, err := NewClassName(testCase.Input.(string))

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else {
				assert.Equal(t, expected.String, className.String(), test.MessageEqual(expected.String, className.String(), testCase.Name))
				assert.Equal(t, expected.Capitalized, className.Capitalized(), test.MessageEqual(expected.Capitalized, className.Capitalized(), testCase.Name))
				assert.Equal(t, expected.Package, className.Package(), test.MessageEqual(expected.Package, className.Package(), testCase.Name))
				assert.Equal(t, expected.Short, className.Short(), test.MessageEqual(expected.Short, className.Short(), testCase.Name))
				assert.Equal(t, expected.MetaStyle, className.MetaStyle(), test.MessageEqual(expected.MetaStyle, className.MetaStyle(), testCase.Name))
				assert.Equal(t, expected.String, fmt.Sprintf("%s", className), test.MessageEqual(expected.String, fmt.Sprintf("%s", className), testCase.Name))
			}
		})
	}
}

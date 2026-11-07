package acceptancetest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

func TestAcceptancePrefixOptions(t *testing.T) {
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_default", Input: PrefixOptions{}, Expected: 8},
		{Name: "test_letters_only_override", Input: PrefixOptions{Characters: "abcdef", Length: 12}, Expected: 12},
		{Name: "test_uppercase_override", Input: PrefixOptions{Characters: "ABCDEF", Length: 6}, Expected: 6},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			options := testCase.Input.(PrefixOptions)
			prefix := Prefix(t, options)
			assert.Len(t, []rune(prefix), testCase.Expected.(int))
			if options.Characters != "" {
				for _, character := range prefix {
					assert.Contains(t, options.Characters, string(character))
				}
			}
		})
	}
}

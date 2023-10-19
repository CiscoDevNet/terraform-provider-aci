package models

import (
	"strconv"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
)

func toStringMap(intf interface{}) map[string]string {

	result := make(map[string]string)
	temp := intf.(map[string]interface{})

	for key, value := range temp {
		A(result, key, value.(string))

	}

	return result
}

func StripQuotes(word string) string {
	if strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "\""), "\"")
	}
	return word
}

func StripSquareBrackets(word string) string {
	if strings.HasPrefix(word, "[") && strings.HasSuffix(word, "]") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "["), "]")
	}
	return word
}

func BoolToString(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func StringToBool(value string) bool {
	if value == "yes" {
		return true
	}
	return false
}

func A(data map[string]string, key, value string) {

	if value != "" {
		data[key] = value
	}

	if value == "{}" {
		data[key] = ""
	}
}

func G(cont *container.Container, key string) string {
	return StripQuotes(cont.S(key).String())
}

func GetDnFromContainer(cont *container.Container) string {
	return CurlyBraces(StripQuotes(cont.S("dn").String()))
}

func GetMOName(dn string) string {
	arr := strings.Split(dn, "/")
	hashedName := arr[len(arr)-1]
	nameArr := strings.Split(hashedName, "-")
	name := strings.Join(nameArr[1:], "-")
	return name

}

func ListFromContainer(cont *container.Container, klass string) []*container.Container {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*container.Container, length)
	for i := 0; i < length; i++ {

		arr[i] = cont.S("imdata").Index(i).S(klass, "attributes")
	}
	return arr

}

func CurlyBraces(value string) string {
	if value == "{}" {
		return ""
	} else {
		return value
	}
}

func GetMORnPrefix(DistinguishedName string) string {
	if DistinguishedName != "" {
		bracketIndex := 0
		rnIndex := 0
		reversedDistinguishedName := reverseString(DistinguishedName)
		for i, v := range reversedDistinguishedName {
			stringValue := string(v)
			if stringValue == "]" {
				bracketIndex = bracketIndex + 1
			} else if stringValue == "[" {
				bracketIndex = bracketIndex - 1
			} else if stringValue == "/" && bracketIndex == 0 {
				rnIndex = i
				break
			}
		}
		return strings.Split(reverseString(reversedDistinguishedName[0:rnIndex]), "-")[0]
	}
	return DistinguishedName
}

func reverseString(s string) string {
	reverse := ""
	for i := len(s) - 1; i >= 0; i-- {
		reverse += string(s[i])
	}
	return reverse
}

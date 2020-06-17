package models

import (
	"strconv"
	"strings"

	"github.com/ciscoecosystem/mso-go-client/container"
)

func toStringMap(intf interface{}) map[string]interface{} {

	result := make(map[string]interface{})
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

func A(data map[string]interface{}, key string, value interface{}) {

	if value != "" {
		data[key] = value
	}

	if value == "{}" {
		data[key] = ""
	}

	if value == nil {
		data[key] = ""
	}
}

func G(cont *container.Container, key string) string {
	return StripQuotes(cont.S(key).String())
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

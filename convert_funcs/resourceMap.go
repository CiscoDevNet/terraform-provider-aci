package convert_funcs

type createFunc func(map[string]interface{}, string) map[string]interface{}

var ResourceMap = map[string]createFunc{}

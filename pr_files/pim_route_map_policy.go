






func extractRelations(relationType string, relations []interface{}) []map[string]Item {
	var children []map[string]Item
	for _, relation := range relations {
		if relationMap, ok := relation.(map[string]interface{}); ok {
			child := map[string]Item{
				relationType: {
					Attributes: relationMap,
				},
			}
			children = append(children, child)
		}
	}
	return children
}

func createChildrenFromAttributes(attributes map[string]interface{}) []map[string]Item {
	var children []map[string]Item

	for key, value := range attributes {
		if nestedObjects, exists := value.([]interface{}); exists {
			children = append(children, extractRelations(key, nestedObjects)...)
		}
	}
	return children
}

func processResource(resourceName string, resourceValues map[string]interface{}) map[string]Item {
	attributes := make(map[string]interface{})
	var children []map[string]Item

	switch resourceName {

	}
	return nil
}


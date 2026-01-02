package tables

/*
import "encoding/json"

// JsonSchemaEncoder is a type responsible for encoding structs into JSON schemas by processing field definitions recursively.
type JsonSchemaEncoder struct{}

// Encode serializes a list of IStructField objects into formatted JSON representing their default values or prototypes.
func (e *JsonSchemaEncoder) Encode(fields []IStructField) ([]byte, error) {
	defaults := make(map[string]interface{})
	for _, field := range fields {
		defaults[field.FieldName()] = e.buildPrototype(field)
	}
	//return json.MarshalIndent(defaults, "", " ")
	return json.Marshal(defaults)
}

// buildPrototype recursively constructs a default prototype value based on the container and kind of the given field.
func (e *JsonSchemaEncoder) buildPrototype(field IStructField) interface{} {
	_, container, kind := field.Options()
	if container == "array" {
		return []interface{}{}
	}
	if container == "map" {
		return map[string]interface{}{}
	}
	switch kind {
	case "struct":
		kk := field.(*Struct)
		subMap := make(map[string]interface{})
		for _, sf := range kk.fields {
			subMap[sf.FieldName()] = e.buildPrototype(sf) // <--- Ricorsione
		}
		return subMap
	case "interface":
		return nil
	case "string":
		return ""
	case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64":
		return 0.0
	case "float32", "float64":
		return 0.0
	case "bool":
		return false
	}
	return nil
}
*/

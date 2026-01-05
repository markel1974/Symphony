package tables

import "encoding/json"

// StructAnalyzer provides utilities for analyzing and processing structured types and their associated metadata.
type StructAnalyzer struct{}

// NewStructAnalyzer creates and returns a new instance of StructAnalyzer.
func NewStructAnalyzer() *StructAnalyzer {
	return &StructAnalyzer{}
}

// Dump generates a JSON representation of the provided fields' default values, structured by their configurations.
func (e *StructAnalyzer) Dump(fields []IStructField) ([]byte, error) {
	defaults := make(map[string]interface{})
	for _, field := range fields {
		defaults[field.FieldName()] = e.buildPrototype(field)
	}
	return json.MarshalIndent(defaults, "", " ")
	//return json.Marshal(defaults)
}

// buildPrototype generates a prototype value for a given IStructField based on its options and kind recursively.
func (e *StructAnalyzer) buildPrototype(field IStructField) interface{} {
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

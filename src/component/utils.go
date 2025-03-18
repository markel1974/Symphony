package component

import "fmt"

func GetSegmentKeys(s interface{}) ([]string, error) {
	if s == nil {
		return nil, fmt.Errorf("error getting segment keys: %s", "nil interface")
	}
	data, ok := s.(map[string]interface{})
	if !ok || data == nil {
		return nil, fmt.Errorf("error getting segment keys: %s", "invalid object")
	}
	var out []string
	for k := range data {
		out = append(out, k)
	}
	return out, nil
}

func GetSegment(id string, s interface{}) (map[string]interface{}, error) {
	if s == nil {
		return nil, fmt.Errorf("error getting segment %s: %s", id, "nil interface")
	}
	data, ok := s.(map[string]interface{})
	if !ok || data == nil {
		return nil, fmt.Errorf("error getting segment %s: %s", id, "invalid interface")
	}
	segmentI, ok := data[id]
	if !ok {
		return nil, fmt.Errorf("error getting segment %s: %s", id, "missing property")
	}
	segment, ok := segmentI.(map[string]interface{})
	if !ok || segment == nil {
		return nil, fmt.Errorf("error getting segment %s: %s", id, "invalid segment interface")
	}
	return segment, nil
}

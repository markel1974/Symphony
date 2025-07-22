package component

import (
	"fmt"
	"strconv"
	"strings"
)

// GetSegmentKeys retrieves all the keys from a map if the provided interface is a map[string]interface{}.
// Returns an error if the input is nil or not a valid map.
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

// GetSegment retrieves a map segment identified by the given id from the provided interface.
// Returns an error if the interface is nil, invalid, or the segment does not exist.
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

// ComponentData parses a colon-separated string into label, id, instance, and validates the input format.
func ComponentData(data string) (string, string, int, error) {
	p := strings.Split(data, ":")
	if len(p) < 3 {
		return "", "", 0, fmt.Errorf("error restoring component %s: %s", data, "invalid component id")
	}
	label := p[0]
	id := p[1]
	instance, err := strconv.Atoi(p[2])
	if err != nil {
		return "", "", 0, fmt.Errorf("error restoring component %s: %s", id, err.Error())
	}
	return label, id, instance, nil
}

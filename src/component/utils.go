package component

import (
	"fmt"
	"strconv"
	"strings"
)

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

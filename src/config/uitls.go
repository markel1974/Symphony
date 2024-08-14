package config

import "strings"

type KV struct {
	K string
	V string
}

func KeyVal(data string) []KV {
	var kvs []KV
	for _, c := range strings.Split(data, ";") {
		k := ""
		v := c
		if opts := strings.Split(c, ":"); len(opts) > 1 {
			k = strings.ToUpper(strings.TrimSpace(opts[0]))
			v = strings.TrimSpace(opts[1])
		}
		kvs = append(kvs, KV{k, v})
	}
	return kvs
}

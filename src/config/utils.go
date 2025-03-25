package config

import (
	"io"
	"os"
	"strings"
)

// KV represents a key-value pair with both key and value as strings.
type KV struct {
	K string
	V string
}

// KeyVal parses a semicolon-separated string into a slice of KV structs, splitting each entry into key and value pairs.
func KeyVal(data string) []KV {
	var kvs []KV
	for _, c := range strings.Split(data, ";") {
		k := ""
		v := c
		if opts := strings.Split(c, ":"); len(opts) > 1 {
			k = strings.TrimSpace(opts[0])
			v = strings.TrimSpace(opts[1])
		}
		kvs = append(kvs, KV{k, v})
	}
	return kvs
}

func ImageFromFile(path string) ([]byte, bool, error) {
	if len(path) == 0 {
		return nil, false, nil
	}
	wp := false
	fd, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		if fd, err = os.OpenFile(path, os.O_RDONLY, 0); err != nil {
			return nil, false, err
		}
		wp = true
	}
	defer fd.Close()
	image, err := io.ReadAll(fd)
	if err != nil {
		return nil, false, err
	}
	return image, wp, nil
}

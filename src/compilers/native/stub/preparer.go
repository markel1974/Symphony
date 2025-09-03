package stub

import (
	"log"
	"os"
	"sort"
	"strings"
)

func Prepare(baseDir string, prefix string) []string {
	data, err := os.ReadDir(baseDir)
	var files []string
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	for _, v := range data {
		if v.IsDir() {
			continue
		}
		if !strings.HasPrefix(v.Name(), prefix) {
			continue
		}
		files = append(files, v.Name())
	}
	sort.Strings(files)
	return files
}

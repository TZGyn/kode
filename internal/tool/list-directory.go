package tool

import (
	"os"
	"sort"
	"strings"
)

func ListDirectory(directory string) ([]string, error) {
	result := []string{}

	entires, err := os.ReadDir("./" + directory)
	if err != nil {
		return result, err
	}

	sort.Slice(entires, func(i, j int) bool {
		if entires[i].IsDir() == entires[j].IsDir() {
			return strings.ToLower(entires[i].Name()) < strings.ToLower(entires[j].Name())
		}
		if !entires[i].IsDir() {
			return true
		}
		return false
	})

	for _, e := range entires {
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		result = append(result, name)
	}

	return result, nil
}

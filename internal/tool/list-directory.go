package tool

import "os"

func ListDirectory(directory string) ([]string, error) {
	result := []string{}

	entires, err := os.ReadDir("./" + directory)
	if err != nil {
		return result, err
	}

	for _, e := range entires {
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}
		result = append(result, name)
	}

	return result, nil
}

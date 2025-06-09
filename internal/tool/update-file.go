package tool

import "os"

func UpdateFile(path string, new_content string) (string, error) {
	err := os.WriteFile("./"+path, []byte(new_content), 0644)
	if err != nil {
		return "", err
	}

	return "File updated", nil
}

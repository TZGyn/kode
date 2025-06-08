package tool

import "os"

func CreateFile(path string) error {
	_, err := os.Create(path)

	if err != nil {
		return err
	}

	return nil
}

package tool

import (
	"os"
)

func CatFile(filePath string) (string, error) {
	result := ""

	// stat, err := os.Stat("./" + filePath)
	// if err != nil {
	// 	return result, err
	// }

	file, err := os.ReadFile("./" + filePath)
	if err != nil {
		return result, err
	}

	result += "```\n"
	result += string(file) + "\n"
	result += "```\n"

	return result, nil
}

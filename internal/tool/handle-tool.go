package tool

import (
	"errors"
	"os"
	"strings"

	"github.com/aymanbagabas/go-udiff"
)

func HandleTool(toolName string, args map[string]any, response *string) (string, error) {
	if toolName == "list_directory" {
		result := []string{}

		directory, ok := args["directory"].(string)
		if ok {
			entires, err := ListDirectory(directory)
			if err == nil {
				result = entires
			}

			toolResult := ""
			toolResult += "## Files Start\n"
			for _, entry := range result {
				toolResult += "- " + entry + "\n"
			}
			toolResult += "## Files End\n"

			*response = *response + toolResult
		}
		return strings.Join(result, "\n"), nil
	}
	if toolName == "cat_file" {
		result := ""

		filePath, ok := args["filePath"].(string)
		if ok {
			content, err := CatFile(filePath)
			if err == nil {
				result = content
			}

			toolResult := ""
			toolResult += "## File content " + filePath + "\n"
			toolResult += result + "\n"
			toolResult += "## File content\n"

			*response = *response + toolResult
		}

		return result, nil
	}
	if toolName == "create_file" {
		result := ""
		path, ok := args["filePath"].(string)
		if ok {
			err := CreateFile(path)
			if err == nil {
				result = "File Created Successfully"
			} else {
				result = err.Error()
			}

			toolResult := ""
			toolResult += "## File create\n"
			toolResult += path + "\n"
			toolResult += "## File create\n"

			*response = *response + toolResult
		}

		return result, nil
	}

	if toolName == "update_file" {
		result := ""
		path, pathOk := args["path"].(string)
		new_content, ok := args["new_content"].(string)
		if ok && pathOk {
			file, err := os.ReadFile("./" + path)
			if err != nil {
				result = err.Error()
				return result, err
			}

			output, err := UpdateFile(path, new_content)
			result = output
			if err != nil {
				result = err.Error()
				return result, err
			}

			edits := udiff.Strings(string(file), new_content)

			unified, _ := udiff.ToUnified("a/"+path, "b/"+path, string(file), edits, 8)

			toolResult := ""
			toolResult += "## File update\n"
			toolResult += "```diff\n"
			toolResult += unified + "\n"
			toolResult += "```\n"
			toolResult += "## File update\n"

			*response = *response + toolResult
		}
		return result, nil
	}

	return "", errors.New("invalid tool")
}

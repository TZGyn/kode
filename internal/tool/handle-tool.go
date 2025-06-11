package tool

import (
	"errors"
	"strings"
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
		return strings.Join(result, ","), nil
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

	if toolName == "apply_patch" {
		result := ""
		patch, ok := args["patch"].(string)
		if ok {
			output, _ := ApplyPatch(patch)
			result = output

			toolResult := ""
			toolResult += "## File patch\n"
			toolResult += patch + "\n"
			toolResult += "## File patch\n"

			*response = *response + toolResult
		}
		return result, nil
	}

	if toolName == "update_file" {
		result := ""
		path, pathOk := args["path"].(string)
		new_content, ok := args["new_content"].(string)
		if ok && pathOk {
			output, err := UpdateFile(path, new_content)
			result = output
			if err != nil {
				result = err.Error()
			}

			toolResult := ""
			toolResult += "## File update " + path + "\n"
			toolResult += new_content + "\n"
			toolResult += "## File update\n"

			*response = *response + toolResult
		}
		return result, nil
	}

	return "", errors.New("Invalid Tool")
}

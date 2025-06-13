package openai

import (
	"github.com/openai/openai-go"
)

var tools = []openai.ChatCompletionToolParam{
	{
		Function: openai.FunctionDefinitionParam{
			Name:        "list_directory",
			Description: openai.String("Given a directory, return all the children of it"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"directory": map[string]string{
						"type":        "string",
						"description": "the directory to output, use . for root",
					},
				},
				"required": []string{"directory"},
			},
		},
	},
	{
		Function: openai.FunctionDefinitionParam{
			Name:        "cat_file",
			Description: openai.String("Given a file path, return all its content as string"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"filePath": map[string]string{
						"type":        "string",
						"description": "the file path to output relative to root",
					},
				},
			},
		},
	},
	{
		Function: openai.FunctionDefinitionParam{
			Name:        "create_file",
			Description: openai.String("Given a file path, create a empty file in the path"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"filePath": map[string]string{
						"type":        "string",
						"description": "the file path to create relative to root",
					},
				},
			},
		},
	},
	{
		Function: openai.FunctionDefinitionParam{
			Name:        "update_file",
			Description: openai.String("Update file, given the complete new file content"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"path": map[string]string{
						"type":        "string",
						"description": "Path of the file",
					},
					"new_content": map[string]string{
						"type":        "string",
						"description": "New file content",
					},
				},
			},
		},
	},
}

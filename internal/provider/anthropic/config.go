package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go"
)

var tools = []anthropic.ToolUnionParam{
	{
		OfTool: &anthropic.ToolParam{
			Name:        "list_directory",
			Description: anthropic.String("Given a directory, return all the children of it"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]any{
					"directory": map[string]string{
						"type":        "string",
						"description": "the directory to output, use . for root",
					},
				},
			},
		},
	},
	{
		OfTool: &anthropic.ToolParam{

			Name:        "cat_file",
			Description: anthropic.String("Given a file path, return all its content as string"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]any{
					"filePath": map[string]string{
						"type":        "string",
						"description": "the file path to output relative to root",
					},
				},
			},
		},
	},
	{
		OfTool: &anthropic.ToolParam{
			Name:        "create_file",
			Description: anthropic.String("Given a file path, create a empty file in the path"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]any{
					"filePath": map[string]string{
						"type":        "string",
						"description": "the file path to create relative to root",
					},
				},
			},
		},
	},
	{
		OfTool: &anthropic.ToolParam{
			Name:        "update_file",
			Description: anthropic.String("Update file, given the complete new file content"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]any{
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

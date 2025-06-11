package google

import (
	"fmt"
	"time"

	"google.golang.org/genai"
)

var googleConfig = &genai.GenerateContentConfig{
	SystemInstruction: &genai.Content{
		Role: "system",
		Parts: []*genai.Part{
			{
				Text: fmt.Sprintf(`
							You are a cli code assistant named kode
							Today's Date: %s

							It is a must to generate some text, letting the user knows your thinking process before using a tool.
							Thus providing better user experience, rather than immediately jump to using the tool and generate a conclusion

							Common Order: Tool, Text
							Better order you must follow: Text, Tool, Text

							You have been given tools to fulfill user request, make sure to keep using them until the user request is fulfilled
							Always check the progress to make sure you dont infinite loop
						`, time.Now().Format("2006-01-02 15:04:05")),
			},
		},
	},
	Tools: tools,
}

var tools = []*genai.Tool{
	{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "list_directory",
				Description: "Given a directory, return all the children of it",
				Parameters: &genai.Schema{
					Type: "object",
					Properties: map[string]*genai.Schema{
						"directory": {
							Type:        "string",
							Description: "the directory to output, use . for root",
						},
					},
				},
				Response: &genai.Schema{
					Type: "object",
					Properties: map[string]*genai.Schema{
						"result": {
							Type:        "string",
							Description: "children as list",
							// Items: &genai.Schema{
							// 	Type:        "string",
							// 	Description: "children, folder or file",
							// },
						},
					},
				},
			},
			{
				Name:        "cat_file",
				Description: "Given a file path, return all its content as string",
				Parameters: &genai.Schema{
					Type: "object",
					Properties: map[string]*genai.Schema{
						"filePath": {
							Type:        "string",
							Description: "the file path to output relative to root",
						},
					},
				},
				Response: &genai.Schema{
					Type: "object",
					Properties: map[string]*genai.Schema{
						"result": {
							Type:        "string",
							Description: "file content",
						},
					},
				},
			},
			{
				Name:        "create_file",
				Description: "Given a file path, create a empty file in the path",
				Parameters: &genai.Schema{
					Type: "object",
					Properties: map[string]*genai.Schema{
						"filePath": {
							Type:        "string",
							Description: "the file path to create relative to root",
						},
					},
				},
				Response: &genai.Schema{
					Type: "object",
					Properties: map[string]*genai.Schema{
						"result": {
							Type:        "string",
							Description: "file create result, either file created successfully or an error message",
						},
					},
				},
			},
			// {
			// 	Name:        "apply_patch",
			// 	Description: "Given a patch file content, apply the patch",
			// 	Parameters: &genai.Schema{
			// 		Type: "object",
			// 		Properties: map[string]*genai.Schema{
			// 			"patch": {
			// 				Type:        "string",
			// 				Description: "the patch file content",
			// 			},
			// 		},
			// 	},
			// 	Response: &genai.Schema{
			// 		Type: "object",
			// 		Properties: map[string]*genai.Schema{
			// 			"result": {
			// 				Type:        "string",
			// 				Description: "file patch result, either file patch successfully or an error message",
			// 			},
			// 		},
			// 	},
			// },
			{
				Name:        "update_file",
				Description: "Update file, given the complete new file content",
				Parameters: &genai.Schema{
					Type: "object",
					Properties: map[string]*genai.Schema{
						"path": {
							Type:        "string",
							Description: "Path of the file",
						},
						"new_content": {
							Type:        "string",
							Description: "New file content",
						},
					},
				},
				Response: &genai.Schema{
					Type: "object",
					Properties: map[string]*genai.Schema{
						"result": {
							Type:        "string",
							Description: "file update result, either file updated successfully or an error message",
						},
					},
				},
			},
		},
	},
}

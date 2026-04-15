package tools

import (
	"context"
	_ "embed"
	"os"
	"strings"

	"charm.land/fantasy"
)

type LSParams struct {
	Path string `json:"path,omitempty" description:"The path to the directory to list (defaults to current working directory)"`
	// Ignore []string `json:"ignore,omitempty" description:"List of glob patterns to ignore"`
	Depth int `json:"depth,omitempty" description:"The maximum depth to traverse"`
}

const LSToolName = "ls"

//go:embed ls.md
var lsDescription []byte

func NewLsTool() fantasy.AgentTool {
	// fmt.Println(string(lsDescription))
	return fantasy.NewAgentTool(
		LSToolName,
		string(lsDescription),
		func(ctx context.Context, input LSParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			files, err := os.ReadDir("./" + input.Path)
			if err != nil {
				return fantasy.ToolResponse{}, err
			}

			result := []string{}

			for _, f := range files {
				name := f.Name()
				result = append(result, name)
			}

			return fantasy.NewTextResponse(printLs(result, input.Path)), nil
		},
	)
}

func printLs(files []string, root string) string {
	var str strings.Builder

	for _, file := range files {
		str.WriteString(file + "\n")
	}

	return str.String()
}

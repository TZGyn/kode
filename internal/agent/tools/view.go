package tools

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"charm.land/fantasy"
)

//go:embed view.md
var viewDescription []byte

type ViewParams struct {
	FilePath string `json:"file_path" description:"The path to the file to read"`
	Offset   int    `json:"offset,omitempty" description:"The line number to start reading from (0-based)"`
	Limit    int    `json:"limit,omitempty" description:"The number of lines to read (defaults to 2000)"`
}

type ViewResourceType string

const (
	ViewResourceUnset ViewResourceType = ""
	ViewResourceSkill ViewResourceType = "skill"
)

type ViewResponseMetadata struct {
	FilePath            string           `json:"file_path"`
	Content             string           `json:"content"`
	ResourceType        ViewResourceType `json:"resource_type,omitempty"`
	ResourceName        string           `json:"resource_name,omitempty"`
	ResourceDescription string           `json:"resource_description,omitempty"`
}

const (
	ViewToolName     = "view"
	MaxViewSize      = 1 * 1024 * 1024 // 1MB
	DefaultReadLimit = 2000
	MaxLineLength    = 2000
)

func NewViewTool() fantasy.AgentTool {
	return fantasy.NewAgentTool(
		ViewToolName,
		string(viewDescription),
		func(ctx context.Context, params ViewParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			file, err := os.Open(params.FilePath)
			if err != nil {
				return fantasy.ToolResponse{}, err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			// Increase buffer size to handle large lines (e.g., minified JSON, HTML)
			// Default is 64KB, set to 1MB
			buf := make([]byte, 0, 64*1024)
			scanner.Buffer(buf, 1024*1024)

			offset := params.Offset

			if offset > 0 {
				skipped := 0
				for skipped < offset && scanner.Scan() {
					skipped++
				}
				if err = scanner.Err(); err != nil {
					return fantasy.ToolResponse{}, err
				}
			}

			limit := params.Limit

			// Pre-allocate slice with expected capacity.
			lines := make([]string, 0, limit)

			for len(lines) < limit && scanner.Scan() {
				lineText := scanner.Text()
				if len(lineText) > MaxLineLength {
					lineText = lineText[:MaxLineLength] + "..."
				}
				lines = append(lines, lineText)
			}

			// Peek one more line only when we filled the limit.
			hasMore := len(lines) == limit && scanner.Scan()

			if err := scanner.Err(); err != nil {
				return fantasy.ToolResponse{}, err
			}

			content := strings.Join(lines, "\n")

			if !utf8.ValidString(content) {
				return fantasy.NewTextErrorResponse("File content is not valid UTF-8"), nil
			}
			output := "<file>\n"
			output += addLineNumbers(content, params.Offset+1)

			if hasMore {
				output += fmt.Sprintf("\n\n(File has more lines. Use 'offset' parameter to read beyond line %d)",
					params.Offset+len(strings.Split(content, "\n")))
			}
			output += "\n</file>\n"

			return fantasy.NewTextResponse(output), nil
		})
}

func addLineNumbers(content string, startLine int) string {
	if content == "" {
		return ""
	}

	lines := strings.Split(content, "\n")

	var result []string
	for i, line := range lines {
		line = strings.TrimSuffix(line, "\r")

		lineNum := i + startLine
		numStr := fmt.Sprintf("%d", lineNum)

		if len(numStr) >= 6 {
			result = append(result, fmt.Sprintf("%s|%s", numStr, line))
		} else {
			paddedNum := fmt.Sprintf("%6s", numStr)
			result = append(result, fmt.Sprintf("%s|%s", paddedNum, line))
		}
	}

	return strings.Join(result, "\n")
}

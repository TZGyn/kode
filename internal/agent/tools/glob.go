package tools

import (
	"context"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"charm.land/fantasy"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/charlievieth/fastwalk"
)

type GlobParams struct {
	Pattern string `json:"pattern,omitempty" description:"Pattern to match, for example **/*.go"`
}

const GlobToolName = "glob"

type GlobResponseMetadata struct {
	NumberOfFiles int  `json:"number_of_files"`
	Truncated     bool `json:"truncated"`
}

//go:embed glob.md
var globDescription []byte

func NewGlobTool() fantasy.AgentTool {
	// fmt.Println(string(lsDescription))
	return fantasy.NewAgentTool(
		GlobToolName,
		string(globDescription),
		func(ctx context.Context, input GlobParams, call fantasy.ToolCall) (fantasy.ToolResponse, error) {

			root := "."

			conf := fastwalk.Config{
				Follow:  false,
				ToSlash: fastwalk.DefaultToSlash(),
			}

			pattern := input.Pattern
			pattern = filepath.ToSlash(pattern)

			files := []string{}

			walkFn := func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
					os.Exit(1)
					return nil // returning the error stops iteration
				}

				if skipHidden(path) {
					return nil
				}

				if match, err := doublestar.PathMatch(pattern, filepath.ToSlash(path)); err == nil {
					if match {
						files = append(files, path)
					}
				} else {
				}
				// _, err = fmt.Println(path)
				return err
			}

			if err := fastwalk.Walk(&conf, root, walkFn); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", root, err)
				os.Exit(1)
				return fantasy.ToolResponse{}, err
			}

			tree := createFileTree(files, root)

			return fantasy.NewTextResponse(printTree(tree, root)), nil
		},
	)
}

type NodeType string

const (
	NodeTypeFile      NodeType = "file"
	NodeTypeDirectory NodeType = "directory"
)

type TreeNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Type     NodeType    `json:"type"`
	Children []*TreeNode `json:"children,omitempty"`
}

func createFileTree(sortedPaths []string, rootPath string) []*TreeNode {
	root := []*TreeNode{}
	pathMap := make(map[string]*TreeNode)

	for _, path := range sortedPaths {
		relativePath := strings.TrimPrefix(path, rootPath)
		parts := strings.Split(relativePath, string(filepath.Separator))
		currentPath := ""
		var parentPath string

		var cleanParts []string
		for _, part := range parts {
			if part != "" {
				cleanParts = append(cleanParts, part)
			}
		}
		parts = cleanParts

		if len(parts) == 0 {
			continue
		}

		for i, part := range parts {
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath = filepath.Join(currentPath, part)
			}

			if _, exists := pathMap[currentPath]; exists {
				parentPath = currentPath
				continue
			}

			isLastPart := i == len(parts)-1
			isDir := !isLastPart || strings.HasSuffix(relativePath, string(filepath.Separator))
			nodeType := NodeTypeFile
			if isDir {
				nodeType = NodeTypeDirectory
			}
			newNode := &TreeNode{
				Name:     part,
				Path:     currentPath,
				Type:     nodeType,
				Children: []*TreeNode{},
			}

			pathMap[currentPath] = newNode

			if i > 0 && parentPath != "" {
				if parent, ok := pathMap[parentPath]; ok {
					parent.Children = append(parent.Children, newNode)
				}
			} else {
				root = append(root, newNode)
			}

			parentPath = currentPath
		}
	}

	return root
}

func printTree(tree []*TreeNode, rootPath string) string {
	var result strings.Builder

	result.WriteString("- ")
	result.WriteString(filepath.ToSlash(rootPath))
	if rootPath[len(rootPath)-1] != '/' {
		result.WriteByte('/')
	}
	result.WriteByte('\n')

	for _, node := range tree {
		printNode(&result, node, 1)
	}

	return result.String()
}

func printNode(builder *strings.Builder, node *TreeNode, level int) {
	indent := strings.Repeat("  ", level)

	nodeName := node.Name
	if node.Type == NodeTypeDirectory {
		nodeName = nodeName + "/"
	}

	fmt.Fprintf(builder, "%s- %s\n", indent, nodeName)

	if node.Type == NodeTypeDirectory && len(node.Children) > 0 {
		for _, child := range node.Children {
			printNode(builder, child, level+1)
		}
	}
}

func skipHidden(path string) bool {
	// Check for hidden files (starting with a dot)
	base := filepath.Base(path)
	if base != "." && strings.HasPrefix(base, ".") {
		return true
	}

	commonIgnoredDirs := map[string]bool{
		".crush":           true,
		"node_modules":     true,
		"vendor":           true,
		"dist":             true,
		"build":            true,
		"target":           true,
		".git":             true,
		".idea":            true,
		".vscode":          true,
		"__pycache__":      true,
		"bin":              true,
		"obj":              true,
		"out":              true,
		"coverage":         true,
		"logs":             true,
		"generated":        true,
		"bower_components": true,
		"jspm_packages":    true,
	}

	parts := strings.SplitSeq(path, string(os.PathSeparator))
	for part := range parts {
		if commonIgnoredDirs[part] {
			return true
		}
	}
	return false
}

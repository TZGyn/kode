package markdown

import (
	"image/color"

	"charm.land/glamour/v2"
	"charm.land/glamour/v2/ansi"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/TZGyn/kode/internal/theme"
	"github.com/lucasb-eyer/go-colorful"
)

var RootPath string
var CwdPath string

const defaultMargin = 1

func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }
func uintPtr(u uint) *uint       { return &u }

func GetMarkdownRenderer(width int) *glamour.TermRenderer {
	r, _ := glamour.NewTermRenderer(
		glamour.WithStyles(generateMarkdownStyleConfig()),
		glamour.WithWordWrap(width),
		// glamour.WithAutoStyle(),
		glamour.WithChromaFormatter("terminal16m"),
	)
	return r
}

func generateMarkdownStyleConfig() ansi.StyleConfig {
	t := theme.GetTheme()
	background := AdaptiveColorToString(t.Background())
	// background := AdaptiveColorToString(compat.AdaptiveColor{
	// 	Light: lipgloss.Color("#ffffff"),
	// 	Dark:  lipgloss.Color("#000000"),
	// })

	return ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix:     "",
				BlockSuffix:     "",
				BackgroundColor: background,
				Color:           AdaptiveColorToString(t.MarkdownText()),
			},
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color:  AdaptiveColorToString(t.MarkdownBlockQuote()),
				Italic: boolPtr(true),
				Prefix: "┃ ",
			},
			Indent:      uintPtr(1),
			IndentToken: stringPtr(" "),
		},
		List: ansi.StyleList{
			LevelIndent: defaultMargin,
			StyleBlock: ansi.StyleBlock{
				IndentToken: stringPtr(" "),
				StylePrimitive: ansi.StylePrimitive{
					Color: AdaptiveColorToString(t.MarkdownText()),
				},
			},
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       AdaptiveColorToString(t.MarkdownHeading()),
				Bold:        boolPtr(true),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "# ",
				Color:  AdaptiveColorToString(t.MarkdownHeading()),
				Bold:   boolPtr(true),
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
				Color:  AdaptiveColorToString(t.MarkdownHeading()),
				Bold:   boolPtr(true),
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
				Color:  AdaptiveColorToString(t.MarkdownHeading()),
				Bold:   boolPtr(true),
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
				Color:  AdaptiveColorToString(t.MarkdownHeading()),
				Bold:   boolPtr(true),
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "##### ",
				Color:  AdaptiveColorToString(t.MarkdownHeading()),
				Bold:   boolPtr(true),
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
				Color:  AdaptiveColorToString(t.MarkdownHeading()),
				Bold:   boolPtr(true),
			},
		},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: boolPtr(true),
			Color:      AdaptiveColorToString(t.TextMuted()),
		},
		Emph: ansi.StylePrimitive{
			Color:  AdaptiveColorToString(t.MarkdownEmph()),
			Italic: boolPtr(true),
		},
		Strong: ansi.StylePrimitive{
			Bold:  boolPtr(true),
			Color: AdaptiveColorToString(t.MarkdownStrong()),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  AdaptiveColorToString(t.MarkdownHorizontalRule()),
			Format: "\n─────────────────────────────────────────\n",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
			Color:       AdaptiveColorToString(t.MarkdownListItem()),
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
			Color:       AdaptiveColorToString(t.MarkdownListEnumeration()),
		},
		Task: ansi.StyleTask{
			Ticked:   "[✓] ",
			Unticked: "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color:     AdaptiveColorToString(t.MarkdownLink()),
			Underline: boolPtr(true),
		},
		LinkText: ansi.StylePrimitive{
			Color: AdaptiveColorToString(t.MarkdownLinkText()),
			Bold:  boolPtr(true),
		},
		Image: ansi.StylePrimitive{
			Color:     AdaptiveColorToString(t.MarkdownImage()),
			Underline: boolPtr(true),
			Format:    "🖼 {{.text}}",
		},
		ImageText: ansi.StylePrimitive{
			Color:  AdaptiveColorToString(t.MarkdownImageText()),
			Format: "{{.text}}",
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BackgroundColor: background,
				Color:           AdaptiveColorToString(t.MarkdownCode()),
				Prefix:          "",
				Suffix:          "",
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					BackgroundColor: background,
					Prefix:          " ",
					Color:           AdaptiveColorToString(t.MarkdownCodeBlock()),
				},
			},
			Chroma: &ansi.Chroma{
				Background: ansi.StylePrimitive{
					BackgroundColor: background,
				},
				Text: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.MarkdownText()),
				},
				Error: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.Error()),
				},
				Comment: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxComment()),
				},
				CommentPreproc: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxKeyword()),
				},
				Keyword: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxKeyword()),
				},
				KeywordReserved: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxKeyword()),
				},
				KeywordNamespace: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxKeyword()),
				},
				KeywordType: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxType()),
				},
				Operator: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxOperator()),
				},
				Punctuation: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxPunctuation()),
				},
				Name: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxVariable()),
				},
				NameBuiltin: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxVariable()),
				},
				NameTag: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxKeyword()),
				},
				NameAttribute: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxFunction()),
				},
				NameClass: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxType()),
				},
				NameConstant: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxVariable()),
				},
				NameDecorator: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxFunction()),
				},
				NameFunction: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxFunction()),
				},
				LiteralNumber: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxNumber()),
				},
				LiteralString: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxString()),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.SyntaxKeyword()),
				},
				GenericDeleted: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.DiffRemoved()),
				},
				GenericEmph: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.MarkdownEmph()),
					Italic:          boolPtr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.DiffAdded()),
				},
				GenericStrong: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.MarkdownStrong()),
					Bold:            boolPtr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					BackgroundColor: background,
					Color:           AdaptiveColorToString(t.MarkdownHeading()),
				},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					BlockSuffix: "\n",
					// TODO: find better way to fix markdown table renders
					BackgroundColor: stringPtr(""),
				},
			},
			CenterSeparator: stringPtr("┼"),
			ColumnSeparator: stringPtr("│"),
			RowSeparator:    stringPtr("─"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: "\n ❯ ",
			Color:       AdaptiveColorToString(t.MarkdownLinkText()),
		},
		Text: ansi.StylePrimitive{
			Color: AdaptiveColorToString(t.MarkdownText()),
		},
		Paragraph: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: AdaptiveColorToString(t.MarkdownText()),
			},
		},
	}
}

type TerminalInfo struct {
	Background       color.Color
	BackgroundIsDark bool
}

var Terminal *TerminalInfo

func init() {
	Terminal = &TerminalInfo{
		Background:       color.Black,
		BackgroundIsDark: true,
	}
}

// AdaptiveColorToString converts a compat.AdaptiveColor to the appropriate
// hex color string based on the current terminal background
func AdaptiveColorToString(color compat.AdaptiveColor) *string {
	if Terminal.BackgroundIsDark {
		if _, ok := color.Dark.(lipgloss.NoColor); ok {
			return nil
		}
		c1, _ := colorful.MakeColor(color.Dark)
		return stringPtr(c1.Hex())
	}
	if _, ok := color.Light.(lipgloss.NoColor); ok {
		return nil
	}
	c1, _ := colorful.MakeColor(color.Light)
	return stringPtr(c1.Hex())
}

func ToMarkdown(content string, width int) string {
	r := GetMarkdownRenderer(width)

	rendered, _ := r.Render(content)
	return rendered
}

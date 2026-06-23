// glow — Markdown в терминале.
//
//	go run . README.md
//	go run . --theme dark README.md
//	cat FILE.md | go run .
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// ──────────────────────────────────────────────
// Темы
// ──────────────────────────────────────────────

type Theme struct {
	H1, H2, H3   lipgloss.Style
	Bold, Italic  lipgloss.Style
	Code          lipgloss.Style
	CodeBlock     lipgloss.Style
	Blockquote    lipgloss.Style
	Link          lipgloss.Style
	HRule         lipgloss.Style
	Normal        lipgloss.Style
}

func darkTheme() Theme {
	return Theme{
		H1:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1),
		H2:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("117")),
		H3:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81")),
		Bold:       lipgloss.NewStyle().Bold(true),
		Italic:     lipgloss.NewStyle().Italic(true),
		Code:       lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Background(lipgloss.Color("236")),
		CodeBlock:  lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Background(lipgloss.Color("235")).Padding(1, 2).MarginTop(1).MarginBottom(1),
		Blockquote: lipgloss.NewStyle().Foreground(lipgloss.Color("244")).BorderLeft(true).BorderStyle(lipgloss.NormalBorder()).PaddingLeft(1),
		Link:       lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Underline(true),
		HRule:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Normal:     lipgloss.NewStyle(),
	}
}

func lightTheme() Theme {
	return Theme{
		H1:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("125")).MarginBottom(1),
		H2:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("25")),
		H3:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("24")),
		Bold:       lipgloss.NewStyle().Bold(true),
		Italic:     lipgloss.NewStyle().Italic(true),
		Code:       lipgloss.NewStyle().Foreground(lipgloss.Color("124")).Background(lipgloss.Color("254")),
		CodeBlock:  lipgloss.NewStyle().Foreground(lipgloss.Color("232")).Background(lipgloss.Color("253")).Padding(1, 2).MarginTop(1).MarginBottom(1),
		Blockquote: lipgloss.NewStyle().Foreground(lipgloss.Color("240")).BorderLeft(true).BorderStyle(lipgloss.NormalBorder()).PaddingLeft(1),
		Link:       lipgloss.NewStyle().Foreground(lipgloss.Color("27")).Underline(true),
		HRule:      lipgloss.NewStyle().Foreground(lipgloss.Color("250")),
		Normal:     lipgloss.NewStyle(),
	}
}

var themes = map[string]func() Theme{
	"dark":  darkTheme,
	"light": lightTheme,
}

// ──────────────────────────────────────────────
// Рендерер
// ──────────────────────────────────────────────

type Renderer struct {
	theme Theme
	width int
	out   strings.Builder
}

func (r *Renderer) renderNode(node ast.Node, src []byte, entering bool) (ast.WalkStatus, error) {
	switch n := node.(type) {

	case *ast.Heading:
		if entering {
			return ast.WalkContinue, nil
		}
		text := r.flushInline()
		level := n.Level
		switch level {
		case 1:
			r.out.WriteString(r.theme.H1.Width(r.width).Render(text) + "\n")
		case 2:
			r.out.WriteString(r.theme.H2.Render("## "+text) + "\n")
		default:
			r.out.WriteString(r.theme.H3.Render(strings.Repeat("#", level)+" "+text) + "\n")
		}

	case *ast.Paragraph:
		if !entering {
			text := r.flushInline()
			r.out.WriteString(r.theme.Normal.Width(r.width).Render(text) + "\n\n")
		}

	case *ast.Text:
		if entering {
			r.inlineAppend(string(n.Segment.Value(src)))
		}

	case *ast.CodeSpan:
		if !entering {
			text := r.flushInline()
			r.inlineAppend(r.theme.Code.Render("`" + text + "`"))
		}

	case *ast.FencedCodeBlock:
		if entering {
			var buf bytes.Buffer
			for i := 0; i < n.Lines().Len(); i++ {
				line := n.Lines().At(i)
				buf.Write(line.Value(src))
			}
			r.out.WriteString(r.theme.CodeBlock.Width(r.width).Render(buf.String()) + "\n")
		}

	case *ast.Blockquote:
		// handled by children

	case *ast.ThematicBreak:
		if entering {
			r.out.WriteString(r.theme.HRule.Render(strings.Repeat("─", r.width)) + "\n\n")
		}

	case *ast.Link:
		if !entering {
			text := r.flushInline()
			dest := string(n.Destination)
			r.inlineAppend(r.theme.Link.Render(text) + " (" + dest + ")")
		}

	case *ast.Emphasis:
		if !entering {
			text := r.flushInline()
			if n.Level == 2 {
				r.inlineAppend(r.theme.Bold.Render(text))
			} else {
				r.inlineAppend(r.theme.Italic.Render(text))
			}
		}

	case *ast.ListItem:
		if !entering {
			text := strings.TrimSpace(r.flushInline())
			r.out.WriteString("  • " + text + "\n")
		}

	case *ast.List:
		if !entering {
			r.out.WriteString("\n")
		}
	}

	return ast.WalkContinue, nil
}

// inline буфер для текста внутри параграфов/заголовков
var inlineBuf strings.Builder

func (r *Renderer) inlineAppend(s string) { inlineBuf.WriteString(s) }
func (r *Renderer) flushInline() string {
	s := inlineBuf.String()
	inlineBuf.Reset()
	return s
}

func render(md []byte, theme Theme, width int) string {
	renderer := &Renderer{theme: theme, width: width}

	parser := goldmark.DefaultParser()
	reader := text.NewReader(md)
	doc := parser.Parse(reader)

	ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		return renderer.renderNode(node, md, entering)
	})

	return renderer.out.String()
}

// ──────────────────────────────────────────────
// main
// ──────────────────────────────────────────────

func main() {
	themeName := flag.String("theme", "dark", "тема: dark, light")
	width     := flag.Int("width", 100, "ширина вывода")
	flag.Parse()

	themeFn, ok := themes[*themeName]
	if !ok {
		fmt.Fprintf(os.Stderr, "неизвестная тема %q. Доступны: dark, light\n", *themeName)
		os.Exit(1)
	}

	var src []byte
	var err error

	if flag.NArg() > 0 {
		src, err = os.ReadFile(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "ошибка чтения файла:", err)
			os.Exit(1)
		}
	} else {
		src, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ошибка чтения stdin:", err)
			os.Exit(1)
		}
	}

	fmt.Print(render(src, themeFn(), *width))
}

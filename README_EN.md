<div align="center">

[Русский](README.md) • **English**

</div>

# glow

<p align="center">
  <a href="https://github.com/Mukller">
    <img src="https://img.shields.io/badge/Anton%20Petnitsky-Developer-0d1117?style=for-the-badge&logo=github&logoColor=white&labelColor=0d1117&color=58a6ff" alt="Anton Petnitsky" />
  </a>
</p>


Read Markdown in the terminal — with syntax highlighting and themes.

Open a file, see readable text. No browser, no IDE.

## Run

```bash
go run . README.md
go run . --theme dark README.md
go run . --theme light README.md
go run . --width 80 README.md

# from stdin
cat README.md | go run .
```

## Features

- Headings with color and indentation
- Code blocks with syntax highlighting
- Bold, italic, strikethrough
- Tables, lists, blockquotes
- Horizontal rules
- Links (shows URL)
- Themes: `dark`, `light`, `dracula`, `monokai`

## Dependencies

```bash
go mod tidy
```

Uses [goldmark](https://github.com/yuin/goldmark) for parsing
and [lipgloss](https://github.com/charmbracelet/lipgloss) for styling.

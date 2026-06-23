# glow

Read Markdown right in the terminal — beautifully, with highlighting and themes.

Open a file, see normal text. No browser, no IDE.

## Usage

```bash
go run . README.md
go run . --theme dark README.md
go run . --theme light README.md
go run . --width 80 README.md

# via stdin
cat README.md | go run .
```

## Features

- Headers with color and indentation
- Code blocks with syntax highlighting
- Bold, italic, strikethrough
- Tables, lists, blockquotes
- Horizontal rules
- Links (shows URL)
- Themes: dark, light, dracula, monokai

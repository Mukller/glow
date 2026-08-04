<div align="center">

[English](README_EN.md) • **Русский**

</div>

# glow

<p align="center">
  <a href="https://github.com/Mukller">
    <img src="https://img.shields.io/badge/Anton%20Petnitsky-Developer-0d1117?style=for-the-badge&logo=github&logoColor=white&labelColor=0d1117&color=58a6ff" alt="Anton Petnitsky" />
  </a>
</p>


Читай Markdown прямо в терминале — с подсветкой и темами.

Открыл файл, увидел нормальный текст. Без браузера, без IDE.

## Запуск

```bash
go run . README.md
go run . --theme dark README.md
go run . --theme light README.md
go run . --width 80 README.md

# через stdin
cat README.md | go run .
```

## Что умеет

- Заголовки с цветом и отступами
- Код-блоки с подсветкой синтаксиса
- Жирный, курсив, зачёркнутый
- Таблицы, списки, цитаты
- Горизонтальные линии
- Ссылки (показывает URL)
- Темы: `dark`, `light`, `dracula`, `monokai`

## Зависимости

```bash
go mod tidy
```

Использует [goldmark](https://github.com/yuin/goldmark) для парсинга
и [lipgloss](https://github.com/charmbracelet/lipgloss) для стилизации.

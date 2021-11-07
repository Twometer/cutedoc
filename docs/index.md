# :ribbon: Cutedoc

Cutedoc is a fast and clean static site generator for beautiful documentation pages, written in Go.

## Features

- Write documentation using Markdown
- Easy to setup, with only a single `.ini` file
- Outputs static HTML that can be hosted anywhere
- Supports syntax highlighting and emojis
- Easily extensible with custom themes using [Go templates](https://pkg.go.dev/text/template)
- Built-in development server

## Getting started

To get started with using cutedoc as your documentation site generator, you need a `cutedoc.ini` file at your project root:

```ini
[Page]
Name = My documentation
Theme = default

[Files]
Input = ./docs
Output = ./docs_gen
```

Only the `Name` key is required. All the other settings will default to the values shown in the example above if omitted.

In the input directory, you need at minimum a `index.md` file which contains the main page of your documentation. The final documentation site will be generated from the Markdown files in the `Input` directory and written into the `Output` directory.

If you want to create your own theme, please see [Creating themes](creating-themes).

## Default theme

The default theme for cutedoc was once a [theme for MkDocs](https://github.com/Twometer/cutedoc/tree/mkdocs), which is where the project started. It is inspired by the GitBook theme, but written completely from scratch. You can find more details about the theme in the [Widget gallery](design/widget-gallery).


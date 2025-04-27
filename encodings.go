package tg

import (
	"strings"
	"sync"
)

const (
	ParseModeHTML       = "HTML"
	ParseModeMarkdown   = "Markdown"
	ParseModeMarkdownV2 = "MarkdownV2"
)

var (
	encodingsHTML = sync.OnceValue(func() *strings.Replacer {
		return strings.NewReplacer(
			"<", "&lt;",
			">", "&gt;",
			"&", "&amp;",
		)
	})

	encodingsMarkdown = sync.OnceValue(func() *strings.Replacer {
		return strings.NewReplacer(
			"_", "\\_",
			"*", "\\*",
			"`", "\\`",
			"[", "\\[",
			"]", "\\]",
		)
	})

	encodingsMarkdownV2 = sync.OnceValue(func() *strings.Replacer {
		return strings.NewReplacer(
			"_", "\\_",
			"*", "\\*",
			"[", "\\[",
			"]", "\\]",
			"(", "\\(",
			")", "\\)",
			"~", "\\~",
			"`", "\\`",
			">", "\\>",
			"#", "\\#",
			"+", "\\+",
			"-", "\\-",
			"=", "\\=",
			"|", "\\|",
			"{", "\\{",
			"}", "\\}",
			".", "\\.",
			"!", "\\!",
		)
	})
)

// EscapeParseMode helps with formatted messages.
// Use example: fmt.Sprintf("```python\n%s\n```", tg.EscapeParseMode(tg.ParseModeMarkdownV2, "print('hello world')"))
// Check https://core.telegram.org/bots/api#formatting-options for Telegram's documentation.
func EscapeParseMode(encoding string, text string) string {
	switch encoding {
	case ParseModeHTML:
		return HTML(text)
	case ParseModeMarkdown:
		return encodingsMarkdown().Replace(text)
	case ParseModeMarkdownV2:
		return Md(text)
	default:
		panic("EscapeParseMode: unknown encoding: " + encoding)
	}
}

func Md(text string) string {
	return encodingsMarkdownV2().Replace(text)
}

func HTML(text string) string {
	return encodingsHTML().Replace(text)
}

package tg

import "strings"

const (
	ParseModeHTML       = "HTML"
	ParseModeMarkdown   = "Markdown"
	ParseModeMarkdownV2 = "MarkdownV2"
)

// EscapeParseMode helps with formatted messages.
// Use example: fmt.Sprintf("```python\n%s\n```", tg.EscapeParseMode(tg.ParseModeMarkdownV2, "print('hello world')"))
// Check https://core.telegram.org/bots/api#formatting-options for Telegram's documentation.
func EscapeParseMode(encoding string, text string) string {
	switch encoding {
	case ParseModeHTML:
		replacer := strings.NewReplacer(
			"<", "&lt;",
			">", "&gt;",
			"&", "&amp;",
		)
		return replacer.Replace(text)

	case ParseModeMarkdown:
		replacer := strings.NewReplacer(
			"_", "\\_",
			"*", "\\*",
			"`", "\\`",
			"[", "\\[",
			"]", "\\]",
		)
		return replacer.Replace(text)

	case ParseModeMarkdownV2:
		replacer := strings.NewReplacer(
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
		return replacer.Replace(text)

	default:
		panic("EscapeParseMode: unknown encoding: " + encoding)
	}
}

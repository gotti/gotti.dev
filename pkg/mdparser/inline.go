package mdparser

import (
	"regexp"
)

// ReplaceInline replaces inline elements
func ReplaceInline(text string) string {
	text = replaceInlineBold(text)
	text = replaceInlineItalic(text)
	text = replaceInlineCode(text)
	text = replaceInlineImage(text)
	text = replaceInlineLink(text)
	return text
}

func replaceInlineBold(text string) string {
	re := regexp.MustCompile(`\*\*([^\*]+)\*\*`)
	return re.ReplaceAllString(text, "<b>$1</b>")
}

func replaceInlineItalic(text string) string {
	re := regexp.MustCompile(`\*([^\*]+)\*`)
	return re.ReplaceAllString(text, "<i>$1</i>")
}

func replaceInlineCode(text string) string {
	re := regexp.MustCompile("`(.+?)`")
	return re.ReplaceAllString(text, "<code>$1</code>")
}

func replaceInlineLink(text string) string {
	exre := regexp.MustCompile(`\[(.*?)\]\((.+?)\)`)
	return exre.ReplaceAllString(text, "<a href=\"$2\">$1</a>")
}

func replaceInlineImage(text string) string {
	re := regexp.MustCompile(`!\[(.*?)\]\((.+?)\)`)
	return re.ReplaceAllString(text, "<div><div>$1</div><div><img src=\"$2\" alt=\"$1\" /></div></div>")
}

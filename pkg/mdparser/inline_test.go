package mdparser

import "testing"

func TestReplaceInlineBold(t *testing.T) {
	text := "**bold**"
	expected := "<b>bold</b>"
	actual := replaceInlineBold(text)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestReplaceInlineItalic(t *testing.T) {
	text := "*italic*"
	expected := "<i>italic</i>"
	actual := replaceInlineItalic(text)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestReplaceInlineCode(t *testing.T) {
	text := "`code`"
	expected := "<code>code</code>"
	actual := replaceInlineCode(text)
	if actual != expected {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}

func TestReplaceInlineLink(t *testing.T) {
	data := []struct {
		input    string
		expected string
	}{
		{
			input:    "[link](http://example.com)",
			expected: "<a href=\"http://example.com\">link</a>",
		},
		{
			input:    "aieuo[link](http://example.com)",
			expected: "aieuo<a href=\"http://example.com\">link</a>",
		},
		{
			input:    "[link](http://example.com)aieuo",
			expected: "<a href=\"http://example.com\">link</a>aieuo",
		},
		{
			input:    "aieuo[link](http://example.com)aieuo",
			expected: "aieuo<a href=\"http://example.com\">link</a>aieuo",
		},
	}
	for _, d := range data {
		actual := replaceInlineLink(d.input)
		if actual != d.expected {
			t.Errorf("Expected %s, but got %s", d.expected, actual)
		}
	}
}

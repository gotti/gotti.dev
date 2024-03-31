package mdparser

import (
	"fmt"
	"regexp"
)

type InlineType string

const (
	InlineTypeBlocks InlineType = "blocks"
	InlineTypeText   InlineType = "text"
	InlineTypeBold   InlineType = "bold"
	InlineTypeItalic InlineType = "italic"
	InlineTypeCode   InlineType = "code"
	InlineTypeLink   InlineType = "link"
	InlineTypeImage  InlineType = "image"
)

type InlineBlock interface {
	GetType() InlineType
	GetMDText() []rune
	ToHTML() string
}

type InlineBlocks struct {
	Children []InlineBlock
	MDText   []rune
}

func (i InlineBlocks) ToHTML() string {
	html := ""
	for _, child := range i.Children {
		html += child.ToHTML()
	}
	return html
}

func (i InlineBlocks) String() string {
	return string(i.MDText)
}

type InlineRichImpl struct {
	Type   InlineType
	MDText []rune
}

func (i InlineRichImpl) GetType() InlineType {
	return i.Type
}

func (i InlineRichImpl) GetMDText() []rune {
	return i.MDText
}

type inlinePlainImpl struct {
	Type   InlineType
	MDText []rune
}

func (i inlinePlainImpl) GetType() InlineType {
	return i.Type
}

func (i inlinePlainImpl) GetMDText() []rune {
	return i.MDText
}

func (i inlinePlainImpl) String() string {
	return string(i.MDText)
}

type InlineBold struct {
	InlineRichImpl
	Text InlineBlocks
}

func (i InlineBold) ToHTML() string {
	return "<strong>" + i.Text.ToHTML() + "</strong>"
}

func (i InlineBold) String() string {
	return string(i.MDText)
}

type ParserFunc func([]rune) InlineBlocks

type ParserInlineBold struct {
}

func (pi *ParserInlineBold) Parse(text []rune, p ParserFunc) (block InlineBlock, consumed int) {
	if len(text) < 2 {
		return nil, 0
	}
	if string(text[0:2]) != "**" {
		return nil, 0
	}
	// find the end of the bold
	for i := 2; i < len(text); i++ {
		if string(text[i:i+2]) == "**" {
			return InlineBold{
				InlineRichImpl: InlineRichImpl{
					Type:   InlineTypeBold,
					MDText: text[0 : i+2],
				},
				Text: p(text[2:i]),
			}, i + 2
		}
	}
	return nil, 0
}

type InlineItalic struct {
	InlineRichImpl
	Text InlineBlocks
}

func (i InlineItalic) ToHTML() string {
	return "<em>" + i.Text.ToHTML() + "</em>"
}

func (i InlineItalic) String() string {
	return string(i.MDText)
}

type ParserInlineItalic struct {
}

func (pi *ParserInlineItalic) Parse(text []rune, p ParserFunc) (block InlineBlock, consumed int) {
	if len(text) < 1 {
		return nil, 0
	}
	if text[0] != '*' {
		return nil, 0
	}
	// find the end of the italic
	for i := 1; i < len(text); i++ {
		if text[i] == '*' {
			return InlineItalic{
				InlineRichImpl{
					Type:   InlineTypeItalic,
					MDText: text[0 : i+1],
				},
				p(text[1:i]),
			}, i + 1
		}
	}
	return nil, 0
}

type InlineCode struct {
	InlineRichImpl
	Text []rune
}

func (i InlineCode) ToHTML() string {
	return "<code>" + string(i.Text) + "</code>"
}

func (i InlineCode) String() string {
	return string(i.MDText)
}

type ParserInlineCode struct {
}

func (pi *ParserInlineCode) Parse(text []rune, p ParserFunc) (block InlineBlock, consumed int) {
	if len(text) < 1 {
		return nil, 0
	}
	if text[0] != '`' {
		return nil, 0
	}
	// find the end of the code
	for i := 1; i < len(text); i++ {
		if text[i] == '`' {
			return InlineCode{
				InlineRichImpl{
					Type:   InlineTypeCode,
					MDText: text[0 : i+1],
				},
				text[1:i],
			}, i + 1
		}
	}
	return nil, 0
}

type InlineLink struct {
	InlineRichImpl
	URL  []rune
	Text InlineBlocks
}

func (i InlineLink) ToHTML() string {
	return "<a href=\"" + string(i.URL) + "\">" + i.Text.ToHTML() + "</a>"
}

func (i InlineLink) String() string {
	return string(i.MDText)
}

type ParserInlineLink struct {
}

func (pi *ParserInlineLink) Parse(text []rune, p ParserFunc) (block InlineBlock, consumed int) {
	if len(text) < 1 {
		return nil, 0
	}
	if text[0] != '[' {
		return nil, 0
	}
	fmt.Println("parseLink", text)
	// find the end of the link
	for i := 1; i < len(text); i++ {
		if text[i] == ']' {
			if i+1 < len(text) && text[i+1] == '(' {
				for j := i + 2; j < len(text); j++ {
					if text[j] == ')' {
						return InlineLink{
							InlineRichImpl{
								Type:   InlineTypeLink,
								MDText: text[0 : j+1],
							},
							text[i+2 : j],
							p(text[1:i]),
						}, j + 1
					}
				}
			}
		}
	}
	return nil, 0
}

type InlineImage struct {
	InlineRichImpl
	URL []rune
	Alt InlineBlocks
}

func (i InlineImage) ToHTML() string {
	return "<img src=\"" + string(i.URL) + "\" alt=\"" + i.Alt.ToHTML() + "\" />"
}

func (i InlineImage) String() string {
	return string(i.MDText)
}

type ParserInlineImage struct {
}

func (pi *ParserInlineImage) Parse(text []rune, p ParserFunc) (block InlineBlock, consumed int) {
	if len(text) < 1 {
		return nil, 0
	}
	if text[0] != '!' {
		return nil, 0
	}
	// find the end of the image
	for i := 1; i < len(text); i++ {
		if text[i] == '[' {
			for j := i + 1; j < len(text); j++ {
				if text[j] == ']' {
					if j+1 < len(text) && text[j+1] == '(' {
						for k := j + 2; k < len(text); k++ {
							if text[k] == ')' {
								return InlineImage{
									InlineRichImpl{
										Type:   InlineTypeImage,
										MDText: text[0 : k+1],
									},
									text[j+2 : k],
									p(text[i+1 : j]),
								}, k + 1
							}
						}
					}
				}
			}
		}
	}
	return nil, 0
}

type InlineImplicitLink struct {
	InlineRichImpl
	URL []rune
}

func (i InlineImplicitLink) ToHTML() string {
	return "<a href=\"" + string(i.URL) + "\">" + string(i.URL) + "</a>"
}

func (i InlineImplicitLink) String() string {
	return string(i.MDText)
}

type ParserImplicitLink struct {
}

func (pi *ParserImplicitLink) Parse(text []rune, p ParserFunc) (block InlineBlock, consumed int) {
	re := regexp.MustCompile(`^https?://`)
	if !re.MatchString(string(text)) {
		return nil, 0
	}
	// find the end of the implicit link
	for i := 0; i < len(text); i++ {
		if text[i] == ' ' || text[i] == '\n' {
			return InlineImplicitLink{
				InlineRichImpl{
					Type:   InlineTypeLink,
					MDText: text[0:i],
				},
				text[0:i],
			}, i
		}
	}
	return InlineImplicitLink{
		InlineRichImpl{
			Type:   InlineTypeLink,
			MDText: text,
		},
		text,
	}, len(text)
}

type InlineText struct {
	InlineRichImpl
	Text []rune
}

func (i InlineText) ToHTML() string {
	return string(i.Text)
}

func (i InlineText) String() string {
	return string(i.MDText)
}

type InlineSingleParser interface {
	Parse(text []rune, p ParserFunc) (block InlineBlock, consumed int)
}

type InlineParser struct {
	parsers []InlineSingleParser
}

func NewInlineParser(addon []InlineSingleParser) *InlineParser {
	parsers := []InlineSingleParser{
		&ParserInlineBold{},
		&ParserInlineItalic{},
		&ParserInlineCode{},
		&ParserInlineImage{},
		&ParserInlineLink{},
		&ParserImplicitLink{},
	}
	parsers = append(parsers, addon...)
	return &InlineParser{
		parsers: parsers,
	}
}

func (p *InlineParser) parse(text []rune) InlineBlocks {
	inlineBlocks := InlineBlocks{
		MDText:   text,
		Children: []InlineBlock{},
	}
	lastedit := 0
	for i := 0; i < len([]rune(text)); i++ {
		for _, parser := range p.parsers {
			if block, consumed := parser.Parse(text[i:], p.parse); block != nil {
				skipped := []rune(text)[lastedit:i]
				if string(skipped) != "" {
					inlineBlocks.Children = append(inlineBlocks.Children, InlineText{InlineRichImpl: InlineRichImpl{
						Type:   InlineTypeText,
						MDText: skipped,
					}, Text: skipped})
				}
				inlineBlocks.Children = append(inlineBlocks.Children, block)
				i += consumed
				lastedit = i
				break
			}
		}
	}
	if lastedit < len(text) {
		inlineBlocks.Children = append(inlineBlocks.Children, InlineText{InlineRichImpl: InlineRichImpl{
			Type:   InlineTypeText,
			MDText: text[lastedit:],
		}, Text: text[lastedit:]})
	}
	return inlineBlocks
}

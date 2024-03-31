package mdparser

import (
	"errors"
	"fmt"
	"regexp"
)

// BlockType is a type of block
type BlockType string

const (
	// LineBlockTypeHeading is a heading
	LineBlockTypeHeading = BlockType("LineBlockHeading")
	// LineBlockTypeListItem is a list item
	LineBlockTypeListItem = BlockType("LineBlockListItem")
	// LineBlockTypeOrderedListItem is an ordered list item
	LineBlockTypeOrderedListItem = BlockType("LineBlockOrderedListItem")
	// LineBlockTypeBlockQuote is a quotation
	LineBlockTypeBlockQuote = BlockType("LineBlockQuotation")
	// LineBlockTypePagination is a pagination
	LineBlockTypePagination = BlockType("LineBlockPagination")
	// LineBlockTypeCodeStartOrEnd is a code block
	LineBlockTypeCodeStartOrEnd = BlockType("LineBlockCode")
	// LineBlockTypeIndented is an indented block
	LineBlockTypeIndented = BlockType("LineBlockIndented")
	// LineBlockTypeDivider is a divider
	LineBlockTypeDivider = BlockType("LineBlockDivider")
	// LineBlockTypeSimple is a normal block
	LineBlockTypeSimple = BlockType("LineBlockSimple")
)

// LineBlock is a block
type LineBlock interface {
	Type() BlockType
	TokenText() string
	InnerText() string
}

type lineBlockImpl struct {
	btype     BlockType
	tokenText string
	innerText string
}

func (l *lineBlockImpl) Type() BlockType {
	return l.btype
}

func (l *lineBlockImpl) TokenText() string {
	return l.tokenText
}

func (l *lineBlockImpl) InnerText() string {
	return l.innerText
}

// LineBlockTokenizer is a block, such as a list or headings
type LineBlockTokenizer struct {
}

var matchers = []LineBlockMatcher{
	LineBlockHeadingMatcher{},
	LineBlockListItemMatcher{},
	LineBlockOrderedListItemMatcher{},
	LineBlockIndentedMatcher{},
	LineBlockDividerMatcher{},
	LineBlockPaginationMatcher{},
	LineBlockCodeStartOrEndMatcher{},
	LineBlockSimpleMatcher{},
}

// Tokenize a line
func (l LineBlockTokenizer) Tokenize(line string) (LineBlock, error) {
	errs := []error{}
	for _, m := range matchers {
		obj, err := m.ParseOnce(line)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		return obj, nil
	}
	return nil, fmt.Errorf("error parsing line: %v, errors: %w", line, errors.Join(errs...))
}

// LineBlockMatcher is a matcher
type LineBlockMatcher interface {
	ParseOnce(line string) (LineBlock, error)
}

// LineBlockHeading is a heading
type LineBlockHeading struct {
	lineBlockImpl
	Level int
}

// LineBlockHeadingMatcher is a heading matcher
type LineBlockHeadingMatcher struct {
}

// ParseOnce a line
func (l LineBlockHeadingMatcher) ParseOnce(line string) (LineBlock, error) {
	r := regexp.MustCompile(`^(#+)\s([^\s]+)`)
	if r.MatchString(line) {
		f := r.FindStringSubmatch(line)

		return &LineBlockHeading{Level: len(f[1]), lineBlockImpl: lineBlockImpl{btype: LineBlockTypeHeading, tokenText: f[1], innerText: f[2]}}, nil
	}
	return nil, fmt.Errorf("error parsing line: %v", line)
}

// LineBlockListItem is a list item
type LineBlockListItem struct {
	lineBlockImpl
}

// LineBlockListItemMatcher is a list item
type LineBlockListItemMatcher struct {
}

// ParseOnce a line
func (l LineBlockListItemMatcher) ParseOnce(line string) (LineBlock, error) {
	r := regexp.MustCompile(`^(\*|\-) `)
	if r.MatchString(line) {
		f := r.FindString(line)
		return &LineBlockListItem{lineBlockImpl: lineBlockImpl{btype: LineBlockTypeListItem, tokenText: f, innerText: line[len(f):]}}, nil
	}
	return nil, fmt.Errorf("error parsing line: %v", line)
}

// LineBlockOrderedListItem is a list item
type LineBlockOrderedListItem struct {
	lineBlockImpl
}

// LineBlockOrderedListItemMatcher is a list item
type LineBlockOrderedListItemMatcher struct {
}

// ParseOnce a line
func (l LineBlockOrderedListItemMatcher) ParseOnce(line string) (LineBlock, error) {
	r := regexp.MustCompile(`^([0-9])+\. `)
	if r.MatchString(line) {
		f := r.FindString(line)
		return &LineBlockOrderedListItem{lineBlockImpl: lineBlockImpl{btype: LineBlockTypeOrderedListItem, tokenText: f, innerText: line[len(f):]}}, nil
	}
	return nil, fmt.Errorf("error parsing line: %v", line)
}

// LineBlockIndented is an indented block
type LineBlockIndented struct {
	lineBlockImpl
	Level int
}

// RemoveIndent remove specified spaces and concatenates the rest and inner
func (l *LineBlockIndented) RemoveIndent(level int) LineBlock {
	if level > l.Level {
		panic("cannot set indent to a higher level")
	}
	fmt.Printf("parsing: %v:%v, tokenLen: %v\n", l.tokenText, l.innerText, len(l.tokenText))
	l.Level -= level
	l.tokenText = l.tokenText[:l.Level]
	l.innerText = l.tokenText[l.Level:] + l.innerText
	fmt.Printf("parsing: %v:%v, tokenLen: %v\n", l.tokenText, l.innerText, len(l.tokenText))
	return l
}

// LineBlockIndentedMatcher is an indented block
type LineBlockIndentedMatcher struct {
}

// ParseOnce a line
func (l LineBlockIndentedMatcher) ParseOnce(line string) (LineBlock, error) {
	r := regexp.MustCompile(`^\s+`)
	if r.MatchString(line) {
		f := r.FindString(line)
		return &LineBlockIndented{Level: len(f), lineBlockImpl: lineBlockImpl{btype: LineBlockTypeIndented, tokenText: f, innerText: line[len(f):]}}, nil
	}
	return nil, fmt.Errorf("error parsing line: %v", line)
}

// LineBlockDivider is a divider, an empty line
type LineBlockDivider struct {
	lineBlockImpl
}

// LineBlockDividerMatcher is a divider, an empty line
type LineBlockDividerMatcher struct {
}

// ParseOnce a line
func (l LineBlockDividerMatcher) ParseOnce(line string) (LineBlock, error) {
	if line == "" {
		return &LineBlockDivider{lineBlockImpl: lineBlockImpl{btype: LineBlockTypeDivider, tokenText: "", innerText: ""}}, nil
	}
	return nil, fmt.Errorf("error parsing line: %v", line)
}

// LineBlockSimple is a normal block
type LineBlockSimple struct {
	lineBlockImpl
}

// LineBlockSimpleMatcher is a normal block
type LineBlockSimpleMatcher struct {
}

// ParseOnce a line
func (l LineBlockSimpleMatcher) ParseOnce(line string) (LineBlock, error) {
	return &LineBlockSimple{lineBlockImpl: lineBlockImpl{btype: LineBlockTypeSimple, tokenText: "", innerText: line}}, nil
}

// LineBlockPagination is a pagination block
type LineBlockPagination struct {
	lineBlockImpl
}

// LineBlockPaginationMatcher is a pagination block
type LineBlockPaginationMatcher struct {
}

// ParseOnce a line
func (l LineBlockPaginationMatcher) ParseOnce(line string) (LineBlock, error) {
	r := regexp.MustCompile(`^---$`)
	if r.MatchString(line) {
		return &LineBlockPagination{lineBlockImpl: lineBlockImpl{btype: LineBlockTypePagination, tokenText: "---", innerText: ""}}, nil
	}
	return nil, fmt.Errorf("error parsing line: %v", line)
}

// LineBlockCode is a code block
type LineBlockCode struct {
	lineBlockImpl
	File string
}

// LineBlockCodeStartOrEndMatcher is a code block
type LineBlockCodeStartOrEndMatcher struct {
}

// ParseOnce a line
func (l LineBlockCodeStartOrEndMatcher) ParseOnce(line string) (LineBlock, error) {
	r := regexp.MustCompile("^```(.*)$")
	if r.MatchString(line) {
		// extract "```<file>" from line
		f := r.FindStringSubmatch(line)
		return &LineBlockCode{lineBlockImpl: lineBlockImpl{btype: LineBlockTypeCodeStartOrEnd, tokenText: "```", innerText: ""}, File: f[1]}, nil
	}
	return nil, fmt.Errorf("error parsing line: %v", line)
}

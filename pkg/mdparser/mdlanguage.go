package mdparser

import (
	"fmt"
	"html"
	"strings"
	"time"
)

// MultiLineType is a type of processing multiline blocks
type MultiLineType = string

const (
	// MultiLineTypeNormal is a normal process
	MultiLineTypeNormal = MultiLineType("MultiLineNormal")
	// MultiLineList is a list process
	MultiLineList = MultiLineType("MultiLineList")
	// MultiLineOrderedList is an ordered list process
	MultiLineOrderedList = MultiLineType("MultiLineOrderedList")
	// MultiLineTypeCode is a code block
	MultiLineTypeCode = MultiLineType("MultiLineCode")
)

// Object is a type of object
type Object interface {
	ToHTML() string
}

// Container can contain objects
type Container interface {
	GetChildren() Objects
}

type EndObject interface {
	GetContents() InlineBlocks
}

type InlineContainerObjectImpl struct {
	Contents InlineBlocks
}

// PlainObjectImpl is a text object implementation
type PlainObjectImpl struct {
	InlineContainerObjectImpl
}

// String returns the text
func (o PlainObjectImpl) String() string {
	return string(o.InlineContainerObjectImpl.Contents.MDText)
}

func (o PlainObjectImpl) GetContents() InlineBlocks {
	return o.InlineContainerObjectImpl.Contents
}

// Objects is a list of objects
type Objects []Object

func (o Objects) GetChildren() Objects {
	return o
}

// ToHTML returns the objects as HTML
func (o Objects) ToHTML() string {
	s := ""
	for _, obj := range o {
		s += obj.ToHTML()
	}
	return s
}

// List is a list
type List struct {
	Objects
}

func (l List) String() string {
	return fmt.Sprintf("List{Items: %v}", l.Objects)
}

// ToHTML returns the list as HTML
func (l List) ToHTML() string {
	s := "<ul>"
	for _, i := range l.Objects {
		s += "<li>" + i.ToHTML() + "</li>"
	}
	s += "</ul>"
	return s
}

// OrderedList is an ordered list
type OrderedList struct {
	Objects
}

func (l OrderedList) String() string {
	return fmt.Sprintf("OrderedList{Items: %v}", l.Objects)
}

// ToHTML returns the ordered list as HTML
func (l *OrderedList) ToHTML() string {
	s := "<ol>"
	for _, i := range l.Objects {
		s += "<li>" + i.ToHTML() + "</li>"
	}
	s += "</ol>"
	return s
}

// Heading is a heading object
type Heading struct {
	Level int
	PlainObjectImpl
}

func (h Heading) String() string {
	return fmt.Sprintf("Heading{Level: %v, Text: %v}", h.Level, h.ToHTML())
}

// ToHTML returns the heading as HTML
func (h Heading) ToHTML() string {
	return fmt.Sprintf("<h%d>%v</h%d>", h.Level, h.PlainObjectImpl.InlineContainerObjectImpl.Contents.ToHTML(), h.Level)
}

// BlockQuote is a quotation object
type BlockQuote struct {
	Objects
}

func (b BlockQuote) String() string {
	return fmt.Sprintf("BlockQuote{Items: %v}", b.Objects)
}

// ToHTML returns the block quote as HTML
func (b BlockQuote) ToHTML() string {
	s := "<blockquote>"
	for _, i := range b.Objects {
		s += i.ToHTML()
	}
	s += "</blockquote>"
	return s
}

// CodeBlock is a code block object
type CodeBlock struct {
	PlainObjectImpl
	FileName string
}

// ToHTML returns the code block as HTML
func (c CodeBlock) ToHTML() string {
	return fmt.Sprintf("<pre>%s\n<code>%v</code></pre>", c.FileName, c.PlainObjectImpl.InlineContainerObjectImpl.Contents.ToHTML())
}

// Divider is a divider object
type Divider struct {
	Objects
}

func (d Divider) String() string {
	return fmt.Sprintf("Divider{Items: %v}", d.Objects)
}

// ToHTML returns the divider as HTML
func (d Divider) ToHTML() string {
	s := "<div>"
	for _, i := range d.Objects {
		s += i.ToHTML()
	}
	s += "</div>\n"
	return s
}

// FrontMatter is a pagination object
type FrontMatter struct {
	MetaData map[string]string
}

func (f FrontMatter) String() string {
	return fmt.Sprintf("FrontMatter{MetaData: %v}", f.MetaData)
}

// Text is a text object
type Text struct {
	PlainObjectImpl
}

// ToHTML returns the text as HTML
func (t Text) ToHTML() string {
	return t.PlainObjectImpl.Contents.ToHTML()
}

func captureIndentedOrAbove(lines []LineBlock, indent int) []LineBlock {
	indented := []LineBlock{}
	for _, l := range lines {
		if l.Type() == LineBlockTypeIndented {
			ind := l.(*LineBlockIndented)
			if ind.Level < indent {
				return indented
			}
			in := ind.RemoveIndent(indent)
			indented = append(indented, in)
		} else {
			return indented
		}
	}
	return indented
}

func reTokenizeString(lines []string) ([]LineBlock, error) {
	t := LineBlockTokenizer{}
	reTokenized := []LineBlock{}
	for _, l := range lines {
		reTokenizedLine, err := t.Tokenize(l)
		if err != nil {
			return nil, fmt.Errorf("reTokenize: %v", err)
		}
		reTokenized = append(reTokenized, reTokenizedLine)
	}
	return reTokenized, nil
}

func reTokenize(lines []LineBlock) ([]LineBlock, error) {
	linesStr := []string{}
	for _, l := range lines {
		linesStr = append(linesStr, l.TokenText()+l.InnerText())
	}
	return reTokenizeString(linesStr)
}

type ParserAddon interface {
	GetSingleParsers() []InlineSingleParser
}

// LineParser is a function to parse lines
type LineParser struct {
	inlineParser *InlineParser
	addons       []ParserAddon
}

func NewLineParser(addon []ParserAddon) *LineParser {
	singleParsers := []InlineSingleParser{}
	for _, a := range addon {
		if p := a.GetSingleParsers(); p != nil {
			singleParsers = append(singleParsers, p...)
		}
	}
	p := NewInlineParser(singleParsers)
	return &LineParser{inlineParser: p, addons: addon}
}

func (p LineParser) parseListItems(lines []LineBlock) (Object, int, error) {
	list := &List{}
	lastConsumed := len(lines) - 1
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		fmt.Printf("Parse List Items: %v\n", l)
		if l.Type() == LineBlockTypeListItem {
			h := p.inlineParser.parse([]rune(l.InnerText()))
			list.Objects = append(list.Objects, &Text{PlainObjectImpl: PlainObjectImpl{InlineContainerObjectImpl: InlineContainerObjectImpl{Contents: h}}})
		} else if l.Type() == LineBlockTypeIndented {
			captured := captureIndentedOrAbove(lines[i:], l.(*LineBlockIndented).Level)
			reTokenized, err := reTokenize(captured)
			if err != nil {
				return nil, 0, fmt.Errorf("Parse List Items: %v", err)
			}
			fmt.Println("retokenized, parsing")
			parsed, err := p.parseLineBlocks(reTokenized)
			if err != nil {
				return nil, 0, fmt.Errorf("Parse List Items: %v", err)
			}
			fmt.Printf("parsed: %v\n", parsed)

			objs := &Objects{}
			*objs = append(*objs, list.Objects[len(list.Objects)-1])
			*objs = append(*objs, parsed.Objects...)
			list.Objects[len(list.Objects)-1] = objs
			i += len(captured) - 1
		} else {
			lastConsumed = i
			break
		}
	}
	return list, lastConsumed, nil
}

func (p LineParser) parseOrderedListItems(lines []LineBlock) (Object, int, error) {
	list := &OrderedList{}
	lastConsumed := len(lines) - 1
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		fmt.Printf("Parse List Items: %v\n", l)
		if l.Type() == LineBlockTypeOrderedListItem {
			h := p.inlineParser.parse([]rune(l.InnerText()))
			list.Objects = append(list.Objects, &Text{PlainObjectImpl: PlainObjectImpl{InlineContainerObjectImpl: InlineContainerObjectImpl{Contents: h}}})
		} else if l.Type() == LineBlockTypeIndented {
			captured := captureIndentedOrAbove(lines[i:], l.(*LineBlockIndented).Level)
			reTokenized, err := reTokenize(captured)
			if err != nil {
				return nil, 0, fmt.Errorf("Parse List Items: %v", err)
			}
			fmt.Println("retokenized, parsing")
			parsed, err := p.parseLineBlocks(reTokenized)
			if err != nil {
				return nil, 0, fmt.Errorf("Parse List Items: %v", err)
			}
			fmt.Printf("parsed: %v\n", parsed)

			objs := &Objects{}
			*objs = append(*objs, list.Objects[len(list.Objects)-1])
			*objs = append(*objs, parsed.Objects...)
			list.Objects[len(list.Objects)-1] = objs
			i += len(captured) - 1
		} else {
			lastConsumed = i
			break
		}
	}
	return list, lastConsumed, nil
}

func (p LineParser) parseBlockQuote(lines []LineBlock) (Object, int, error) {
	innerTexts := []string{}
	consumed := 0
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if l.Type() == LineBlockTypeBlockQuote {
			innerTexts = append(innerTexts, l.InnerText())
		} else {
			consumed = i
			break
		}
	}
	reTokenized, err := reTokenizeString(innerTexts)
	if err != nil {
		return nil, 0, fmt.Errorf("Parse BlockQuote: %v", err)
	}
	parsed, err := p.parseLineBlocks(reTokenized)
	if err != nil {
		return nil, 0, fmt.Errorf("Parse BlockQuote: %v", err)
	}

	objs := parsed

	return &BlockQuote{
		Objects: objs.Objects,
	}, consumed, err
}

func parseCodeBlock(lines []LineBlock) (Object, int, error) {
	innerTexts := []string{}
	count := 0
	filename := ""
	fmt.Printf("Parse CodeBlock start: %v\n", lines[1])
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if l.Type() == LineBlockTypeCodeStartOrEnd {
			if count == 0 {
				filename = l.(*LineBlockCode).File
				count++
			} else {
				fmt.Printf("Parse CodeBlock end: %v\n", lines[i-1])
				return &CodeBlock{
					FileName: filename,
					PlainObjectImpl: PlainObjectImpl{
						InlineContainerObjectImpl: InlineContainerObjectImpl{
							Contents: InlineBlocks{
								Children: []InlineBlock{
									InlineText{
										Text: []rune(html.EscapeString(strings.Join(innerTexts, "\n"))),
									},
								},
							},
						},
					},
				}, i, nil
			}
		} else {
			fmt.Printf("Parse CodeBlock: %v, inner: %v\n", l, l.TokenText()+l.InnerText())
			innerTexts = append(innerTexts, l.InnerText())
		}
	}
	return nil, 0, fmt.Errorf("Parse CodeBlock: no end")
}

func (p LineParser) parseDivider(lines []LineBlock) (Object, int, error) {
	initial := false
	lineblocks := []LineBlock{}
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		fmt.Printf("Parse Divider: %v\n", l)
		if l.Type() == LineBlockTypeDivider {
			if initial == false {
				initial = true
			} else {
				break
			}
		} else if l.Type() != LineBlockTypeHeading {
			lineblocks = append(lineblocks, l)
		} else {
			break
		}
	}
	if len(lineblocks) == 0 {
		fmt.Println("no lineblocks")
		return nil, 0, nil
	}
	fmt.Println("lineblocks", lineblocks)
	obj, err := p.parseLineBlocks(lineblocks)
	if err != nil {
		return nil, 0, fmt.Errorf("Parse Divider: %v", err)
	}
	return &Divider{Objects: obj.Objects}, len(lineblocks) + 2, nil
}

func parseFrontMatter(lines []LineBlock) (map[string]string, int, error) {
	initial := false
	ret := map[string]string{}
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if l.Type() == LineBlockTypePagination {
			if initial == false {
				initial = true
			} else {
				return ret, i, nil
			}
		} else {
			s := strings.Split(l.InnerText(), ": ")
			if len(s) != 2 {
				return nil, 0, fmt.Errorf("Parse FrontMatter: invalid format, %v", l.InnerText())
			}
			ret[s[0]] = s[1]
		}
	}
	return nil, 0, fmt.Errorf("Parse FrontMatter: no end")
}

// Root is root object
type Root struct {
	Objects
	MetaData MetaData
}

// MetaData is meta data
type MetaData struct {
	Tags      []string
	Title     *string
	Thumbnail *string
	Date      *time.Time
}

// Parse parses lines
func (p LineParser) Parse(lines string) (*Root, error) {
	t := LineBlockTokenizer{}
	lineBlocks := []LineBlock{}
	l := strings.ReplaceAll(lines, "\r\n", "\n")
	d := strings.Split(l, "\n")
	for _, li := range d {
		l, err := t.Tokenize(li)
		if err != nil {
			return nil, fmt.Errorf("ParseRaw: %v", err)
		}
		if l == nil {
			return nil, fmt.Errorf("ParseRaw: nil lineblock for line=%v", li)
		}
		lineBlocks = append(lineBlocks, l)
	}
	return p.parseLineBlocks(lineBlocks)
}

// parseLineBlocks parses lines
func (p LineParser) parseLineBlocks(lines []LineBlock) (*Root, error) {
	root := Root{}
	tmp := Objects{}
	fmt.Printf("%v\n", lines)
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		fmt.Printf("Parse: %v\n", l)
		switch l.Type() {
		case LineBlockTypeHeading:
			tmp = append(tmp, &Heading{
				Level: l.(*LineBlockHeading).Level,
				PlainObjectImpl: PlainObjectImpl{
					InlineContainerObjectImpl: InlineContainerObjectImpl{
						Contents: InlineBlocks{
							Children: []InlineBlock{
								InlineText{
									Text: []rune(html.EscapeString(l.InnerText())),
								},
							},
						},
					},
				},
			})
			root.Objects = append(root.Objects, tmp)
			//fmt.Printf("Parse: %v\n", root.Objects)
			tmp = Objects{}
		case LineBlockTypeListItem:
			list, lastConsumed, err := p.parseListItems(lines[i:])
			if err != nil {
				return nil, fmt.Errorf("Parse List: %v", err)
			}
			tmp = append(tmp, list)
			i += lastConsumed
		case LineBlockTypeOrderedListItem:
			list, lastConsumed, err := p.parseOrderedListItems(lines[i:])
			if err != nil {
				return nil, fmt.Errorf("Parse List: %v", err)
			}
			tmp = append(tmp, list)
			i += lastConsumed + 1
		case LineBlockTypeBlockQuote:
			p.parseBlockQuote(lines[i:])
		case LineBlockTypePagination:
			if i == 0 { // ParseFrontMatter
				f, lastConsumed, err := parseFrontMatter(lines[i:])
				if err != nil {
					return nil, fmt.Errorf("Parse FrontMatter: %v", err)
				}
				for k, v := range f {
					fmt.Printf("Parse FrontMatter: %v, %v\n", k, v)
					v = strings.TrimPrefix(v, "\"")
					v = strings.TrimSuffix(v, "\"")
					switch k {
					case "title":
						tmp := v
						root.MetaData.Title = &tmp
					case "date":
						t, err := time.Parse("2006-01-02", v)
						if err != nil {
							t2, err := time.Parse("2006-1-2", v)
							if err != nil {
								return nil, fmt.Errorf("Parse FrontMatter datetime=%v, %w", v, err)
							}
							t = t2
						}
						root.MetaData.Date = &t
					case "thumbnail":
						tmp := v
						root.MetaData.Thumbnail = &tmp
					}
				}
				i += lastConsumed
			}
		case LineBlockTypeCodeStartOrEnd:
			line, lastConsumed, err := parseCodeBlock(lines[i:])
			if err != nil {
				return nil, fmt.Errorf("Parse CodeBlock: %v", err)
			}
			tmp = append(tmp, line)
			i += lastConsumed
		case LineBlockTypeDivider:
			if len(tmp) > 0 {
				root.Objects = append(root.Objects, &Divider{Objects: tmp})
				tmp = Objects{}
			}
		case LineBlockTypeIndented:
			fallthrough
		case LineBlockTypeSimple:
			fallthrough
		default:
			if len(tmp) > 1 {
				if _, ok := tmp[len(tmp)-1].(*Text); ok {
					obs := p.inlineParser.parse([]rune(l.InnerText()))
					tmp[len(tmp)-1].(*Text).PlainObjectImpl.InlineContainerObjectImpl.Contents.Children = append(tmp[len(tmp)-1].(*Text).PlainObjectImpl.InlineContainerObjectImpl.Contents.Children, obs.Children...)
				}
			} else {
				obs := p.inlineParser.parse([]rune(l.InnerText()))
				tmp = append(tmp, &Text{PlainObjectImpl: PlainObjectImpl{
					InlineContainerObjectImpl: InlineContainerObjectImpl{
						Contents: InlineBlocks{
							Children: obs.Children,
						},
					},
				}})
			}
		}
	}
	if len(tmp) > 0 {
		root.Objects = append(root.Objects, &tmp)
	}
	return &root, nil
}

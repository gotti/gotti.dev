package mdparser

import (
	"reflect"
	"testing"
)

func TestParseList(t *testing.T) {
	l := LineBlockTokenizer{}
	data := []struct {
		inputLine string
		expected  LineBlock
	}{
		{"- item1", &LineBlockListItem{lineBlockImpl: lineBlockImpl{btype: LineBlockTypeListItem, tokenText: "- ", innerText: "item1"}}},
		{"    hoge", &LineBlockIndented{lineBlockImpl: lineBlockImpl{btype: LineBlockTypeIndented, tokenText: "    ", innerText: "hoge"}}},
	}
	for _, d := range data {
		actual, err := l.Tokenize(d.inputLine)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if actual.Type() != d.expected.Type() {
			t.Errorf("Type: expected %v, but got %v", d.expected.Type(), actual.Type())
		}
		if actual.TokenText() != d.expected.TokenText() {
			t.Errorf("TokenText: expected %v, but got %v", d.expected.TokenText(), actual.TokenText())
		}
		if actual.InnerText() != d.expected.InnerText() {
			t.Errorf("InnerText: expected %v, but got %v", d.expected.InnerText(), actual.InnerText())
		}
	}
}

func TestParseTags(t *testing.T) {
	l := LineBlockTokenizer{}
	data := []struct {
		inputLine string
		expected  LineBlock
	}{
		{"#hoge #huga", &LineBlockTags{lineBlockImpl: lineBlockImpl{btype: LineBlockTypeTags, tokenText: "hoge huga"}}},
	}
	for _, d := range data {
		actual, err := l.Tokenize(d.inputLine)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if actual.Type() != d.expected.Type() {
			t.Errorf("Type: expected %v, but got %v", d.expected.Type(), actual.Type())
		}
		if actual.TokenText() != d.expected.TokenText() {
			t.Errorf("TokenText: expected %v, but got %v", d.expected.TokenText(), actual.TokenText())
		}
	}

}

func TestParseHeading(t *testing.T) {
	l := LineBlockTokenizer{}
	data := []struct {
		inputLine string
		expected  LineBlock
	}{
		{"# hoge", &LineBlockHeading{Level: 1, lineBlockImpl: lineBlockImpl{btype: LineBlockTypeHeading, tokenText: "#", innerText: "hoge"}}},
		{"## hoge", &LineBlockHeading{Level: 2, lineBlockImpl: lineBlockImpl{btype: LineBlockTypeHeading, tokenText: "##", innerText: "hoge"}}},
	}
	for _, d := range data {
		actual, err := l.Tokenize(d.inputLine)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !reflect.DeepEqual(actual, d.expected) {
			t.Errorf("InnerText: expected %v, but got %v", d.expected, actual)
		}
	}
}

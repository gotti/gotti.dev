package mdparser

import (
	"testing"
)

func Test1(t *testing.T) {
  parser := NewSingleParser()
  md := ",[光クロスではPPPoEを提供するISPはほとんど無い](https://flets.com/cross/pppoe/isp.html)．"
  actual := parser.parse([]rune(md))
  expected := "hello <strong>bold</strong>"
  if actual.ToHTML() != expected {
    t.Errorf("HTML: Expected %s, but got %s, %+v", expected, actual.ToHTML(), actual)
  }
}


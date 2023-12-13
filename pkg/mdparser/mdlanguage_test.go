package mdparser

import (
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	data := []struct {
		input    string
		expected Objects
	}{
		{
			input:	
`- item1
- item2`,
			expected: Objects{
						&List{
							Objects: Objects{&Text{TextObjectImpl: TextObjectImpl{"item1"}}, &Text{TextObjectImpl: TextObjectImpl{"item2"}}},
						},
					},
		},
		{
			input:	
`- item1
    - hogehoge
- item2`,
			expected: Objects{
					&List{
						Objects: Objects{
							Objects{
								&Text{TextObjectImpl: TextObjectImpl{"item1"}},
								&List{
									Objects: Objects{&Text{TextObjectImpl: TextObjectImpl{"hogehoge"}}},
								},
							},
						&Text{TextObjectImpl: TextObjectImpl{"item2"}},
						},
					},
			},
		},
		{
			input:	
`- item1
    - hogehoge
    - hogehuga
- item2`,
			expected: Objects{
					&List{
						Objects: Objects{
							&Objects{
								&Text{TextObjectImpl: TextObjectImpl{"item1"}},
								&List{
									Objects: Objects{
										&Text{TextObjectImpl: TextObjectImpl{"hogehoge"}},
										&Text{TextObjectImpl: TextObjectImpl{"hogehuga"}},
									},
								},
							},
							&Text{TextObjectImpl: TextObjectImpl{"item2"}},
						},
					},
				},
		},
		{
			input:	
`- item1
    - hogehoge
        - hogehuga
- item2`,
			expected: Objects{
					&List{
						Objects: []Object{&Objects{
							&Text{TextObjectImpl: TextObjectImpl{"item1"}},
							&List{
								Objects: Objects{&Objects{
									&Text{TextObjectImpl:  TextObjectImpl{"hogehoge"}},
									&List{
										Objects: []Object{&Text{TextObjectImpl: TextObjectImpl{"hogehuga"}}},
									},
								}}},
						},
						&Text{TextObjectImpl: TextObjectImpl{"item2"}}},
					},
				},
		},
		{
			input:	
`- item1
    1. hogehoge
    2. hogehuga
- item2`,
			expected: Objects{
					&List{
						Objects: []Object{
							&Objects{
								&Text{TextObjectImpl: TextObjectImpl{"item1"}},
								&OrderedList{
									Objects: Objects{
										&Text{TextObjectImpl: TextObjectImpl{"hogehoge"}},
										&Text{TextObjectImpl: TextObjectImpl{"hogehuga"}},
									},
								},
							},
							&Text{TextObjectImpl: TextObjectImpl{"item2"}}},
					},
				},
			},
	}
	for _, d := range data {
		actual, err := Parse(d.input)
		t.Logf("\n%v", d.input)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if !reflect.DeepEqual(actual.Objects.String(), d.expected.String()){
			t.Errorf("In expected: \n%v\n but got: \n%v", d.expected, actual.Objects)
			t.Errorf("In expected: \n%T\n but got: \n%T", d.expected, actual.Objects)
		}
	}
}

package orgmode

import (
	"reflect"
	"strings"
	"testing"
)

func TestHeadingParser(t *testing.T) {
	text := `* level 1
line 2`
	n, c := headingParser(text)
	if n.Type != NodeTypeHeading {
		t.Fatal("unexpected type:", n.Type)
	}
	if n.Heading != "level 1" {
		t.Fatal("unexpected heading:", n.Heading)
	}
	if c != strings.Index(text, "line 2") {
		t.Fatal("unexpected consumed length:", n)
	}
}

func TestDescItemParser(t *testing.T) {
	text := `- search engine :: a fast document finder
Line 2
`
	n, c := descriptionItemParser(text)
	expected := NewNode(NodeTypeDescItem, "", DescItem{Term: "search engine", Desc: "a fast document finder"})
	if !reflect.DeepEqual(expected, n) {
		t.Fatal("unexpected node:", n, expected)
	}
	if c != strings.Index(text, "Line 2") {
		t.Fatal("unexpected consumed:", c)
	}
}

func TestTextParser(t *testing.T) {
	text := `Line 1
Line 2

Line 4
- foo :: bar
Line 5
`
	n, c := textParser(text)
	expected := NewNode(NodeTypeText, "", `Line 1
Line 2

Line 4
`)
	if !reflect.DeepEqual(expected, n) {
		t.Fatal("unexpected node:", n, expected)
	}
	if c != strings.Index(text, "- foo :: bar") {
		t.Fatal("unexpected consumed:", c)
	}
}

func TestOrgDocParser(t *testing.T) {
	testData := []struct {
		text     string
		expected *OrgDoc
	}{
		{
			`before
* heading 1
`,
			&OrgDoc{
				[]*Node{
					{NodeTypeText, "_", "before\n"},
					{NodeTypeHeading, "heading 1", ""},
				},
			},
		},
		{
			`* heading 1
* heading 2
`,
			&OrgDoc{
				[]*Node{
					{NodeTypeHeading, "heading 1", ""},
					{NodeTypeHeading, "heading 2", ""},
				},
			},
		},
		{
			`* heading 1
text
- term :: desc
- home :: where the heart is
`,
			&OrgDoc{
				[]*Node{
					{NodeTypeHeading, "heading 1", ""},
					{NodeTypeText, "heading 1", "text\n"},
					{NodeTypeDescItem, "heading 1", DescItem{"term", "desc"}},
					{NodeTypeDescItem, "heading 1", DescItem{"home", "where the heart is"}},
				},
			},
		},
	}
	for idx, test := range testData {
		orgDoc, _ := ParseOrgDoc(test.text)
		if !reflect.DeepEqual(orgDoc, test.expected) {
			t.Error("failed test case:", idx+1)
		}
	}
}

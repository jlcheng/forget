package orgmode

import (
	"bytes"
	"fmt"
	"github.com/jlcheng/forget/trace"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type Node struct {
	Type    string
	Heading string
	Value   interface{}
}

type OrgParser struct {
	heading string
	text    string
	index   int
	nodes   []*Node
}

type OrgDoc struct {
	Nodes []*Node
}

type DescItem struct {
	Term string
	Desc string
}

type subparser func(text string) (*Node, int)

var NodeTypeText string = "TEXT"
var NodeTypeHeading string = "HEADING"
var NodeTypeDescItem string = "DESC_ITEM"

func NewNode(nodeType string, heading string, value interface{}) *Node {
	return &Node{
		Type:    nodeType,
		Heading: heading,
		Value:   value,
	}
}

func NewParser(text string) OrgParser {
	return OrgParser{
		heading: "_",
		text:    text,
		index:   0,
		nodes:   make([]*Node, 0),
	}
}

func textParser(text string) (*Node, int) {
	idx := 0
	contents := ""
	for idx < len(text) {
		// Stop consuming contents of text if it starts with a heading or a description item
		_, consumed := headingParser(text[idx:])
		if consumed != 0 {
			break
		}
		_, consumed = descriptionItemParser(text[idx:])
		if consumed != 0 {
			break
		}

		// Read in one line of the substring. We will have to check for heading/description on every
		// line
		ni := strings.IndexRune(text[idx:], '\n')
		if ni == -1 {
			// No more newlines, read in the rest of text
			ni = len(text)
		} else {
			// Read in the line, including the newline
			ni += 1
		}
		contents += text[idx : idx+ni]
		idx += ni // Move the start of the next substring forward
	}

	// Return a text node if we found contents, otherwise, return nil
	if contents != "" {
		return NewNode(NodeTypeText, "", contents), idx
	}
	return nil, 0
}

func headingParser(text string) (*Node, int) {
	headingP := regexp.MustCompile(`^\*+\s+?(?P<text>.+)\n?`)
	matches := headingP.FindStringSubmatch(text)
	if matches == nil {
		return nil, 0
	}
	return NewNode(NodeTypeHeading, matches[1], ""), len(matches[0])
}

func descriptionItemParser(text string) (*Node, int) {
	descItemP := regexp.MustCompile(`^-\s+?(?P<term>.+?)\s+::\s+(?P<desc>.+)\n?`)
	matches := descItemP.FindStringSubmatch(text)
	if matches == nil {
		return nil, 0
	}

	return NewNode(NodeTypeDescItem, "", DescItem{matches[1], matches[2]}), len(matches[0])
}

func (parser OrgParser) parseNext() (*Node, int) {
	parsers := []subparser{
		headingParser,
		descriptionItemParser,
		textParser,
	}
	for _, p := range parsers {
		node, consumed := p(parser.text[parser.index:])
		if consumed != 0 {
			return node, consumed
		}
	}
	return nil, 0
}

func (parser OrgParser) finished() bool {
	return parser.index >= len(parser.text)
}

func ParseOrgDoc(text string) (*OrgDoc, error) {
	doc := &OrgDoc{
		Nodes: make([]*Node, 0),
	}
	orgParser := NewParser(text)
	for !orgParser.finished() {
		node, consumed := orgParser.parseNext()
		if node == nil {
			break
		}
		orgParser.index += consumed
		doc.Nodes = append(doc.Nodes, node)

		if node.Type == NodeTypeHeading {
			// Store the last heading
			orgParser.heading = node.Heading
		} else {
			// Non-heading nodes should be tagged with the last seen heading
			node.Heading = orgParser.heading
		}
	}
	if !orgParser.finished() {
		trace.Warn("unparsable org-mode file")
		trace.Warn(text)
		return nil, errors.New("unparsable org-mode file")
	}
	return doc, nil
}

func (doc *OrgDoc) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for _, node := range doc.Nodes {
		buf.WriteString(fmt.Sprintf("%v", node))
		buf.WriteString(",\n")
	}
	buf.WriteString("]")
	return buf.String()
}

func (node Node) String() string {
	if _, ok := node.Value.(string); ok {
		return fmt.Sprintf("{Type: %v, Heading: \"%v\", value: \"%v\"}", node.Type, node.Heading, node.Value)
	}
	return fmt.Sprintf("{Type: %v, Heading: \"%v\", value: %v}", node.Type, node.Heading, node.Value)
}

func (node Node) TextValue() string {
	textContent, ok := node.Value.(string)
	if !ok {
		return ""
	}
	return textContent
}

func (node Node) DescItemValue() DescItem {
	descItem, ok := node.Value.(DescItem)
	if !ok {
		return DescItem{}
	}
	return descItem
}

func (descItem DescItem) String() string {
	return fmt.Sprintf("{%v=%v}", descItem.Term, descItem.Desc)
}

package xhtml

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Node struct {
	*html.Node
}

func NewNode(node *html.Node) Node {
	return Node{
		node,
	}
}

func NilNode() Node {
	return Node{
		nil,
	}
}

func Parse(r io.Reader) (Node, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return NilNode(), err
	}
	return NewNode(doc), nil
}

func (node Node) String() string {
	return fmt.Sprintf("%+v", node.Node)
}

func (node Node) IsNil() bool {
	return node.Node == nil
}

func (node Node) GetData() string {
	if node.IsNil() {
		return ""
	}
	return node.Data
}

func (node Node) IsElement() bool {
	if node.IsNil() {
		return false
	}
	return node.Type == html.ElementNode
}

func (node Node) IsText() bool {
	if node.IsNil() {
		return false
	}
	return node.Type == html.TextNode
}

func (node Node) IsLeaf() bool {
	if node.IsNil() {
		return false
	}
	return node.FirstChild == nil
}

func (node Node) PrevElement() Node {
	if node.IsNil() {
		return NilNode()
	}
	for prev := node.PrevSibling; prev != nil; prev = prev.PrevSibling {
		if prev.Type == html.ElementNode {
			return NewNode(prev)
		}
	}

	return NilNode()
}

func (node Node) NextElement() Node {
	if node.IsNil() {
		return NilNode()
	}
	for next := node.NextSibling; next != nil; next = next.NextSibling {
		if next.Type == html.ElementNode {
			return NewNode(next)
		}
	}

	return NilNode()
}

func (node Node) NextChildElement() Node {
	if node.IsNil() {
		return NilNode()
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			return NewNode(child)
		}
	}

	return NilNode()
}

func (node Node) GetAllChildElements() []Node {
	children := []Node{}
	if node.IsNil() {
		return children
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			children = append(children, NewNode(child))
		}
	}

	return children
}

func (node Node) GetAttributes() []string {
	keys := []string{}

	if node.IsNil() {
		return keys
	}
	for _, a := range node.Attr {
		keys = append(keys, a.Key)
	}

	return keys
}

func (node Node) GetAttributeValue(key string) (string, bool) {
	if node.IsNil() {
		return "", false
	}
	for _, a := range node.Attr {
		if strings.EqualFold(a.Key, key) {
			return a.Val, true
		}
	}

	return "", false
}

func (node Node) GetNodeLink() (string, bool) {
	if node.IsNil() {
		return "", false
	}
	if node.Type == html.ElementNode &&
		strings.EqualFold(node.Data, A) {
		return node.GetAttributeValue(HREF)
	}
	if node.Type == html.ElementNode &&
		strings.EqualFold(node.Data, IMG) {
		return node.GetAttributeValue(SRC)
	}

	return "", false
}

func (node Node) ConcatAllText(sep string) string {
	if node.IsNil() {
		return ""
	}
	if node.IsLeaf() && node.IsText() {
		return node.GetData()
	}
	texts := []string{}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		txt := NewNode(child).ConcatAllText(sep)
		if strings.TrimSpace(txt) != "" {
			texts = append(texts, txt)
		}
	}

	return strings.Join(texts, sep)
}

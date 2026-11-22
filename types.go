package xhtml

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Node wraps an *html.Node and adds navigation, search and extraction helpers.
// The zero value, also returned by NilNode, represents the absence of a node:
// every method can be safely called on it.
type Node struct {
	*html.Node
}

// NewNode wraps node into a Node. A nil node gives the same result as NilNode.
func NewNode(node *html.Node) Node {
	return Node{
		node,
	}
}

// NilNode returns a Node that holds no underlying node. See Node.IsNil.
func NilNode() Node {
	return Node{
		nil,
	}
}

// Parse parses the HTML document read from r, as html.Parse does, and returns
// its root node.
func Parse(r io.Reader) (Node, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return NilNode(), err
	}
	return NewNode(doc), nil
}

// String returns a debug representation of the underlying html.Node.
func (node Node) String() string {
	return fmt.Sprintf("%+v", node.Node)
}

// IsNil reports whether node holds no underlying node. The Find* methods and
// the navigation methods return a nil Node when nothing is found.
func (node Node) IsNil() bool {
	return node.Node == nil
}

// GetData returns the data of the node: the tag name for an element, the
// content for a text node. It returns an empty string for a nil node.
func (node Node) GetData() string {
	if node.IsNil() {
		return ""
	}
	return node.Data
}

// IsElement reports whether node is an element node.
func (node Node) IsElement() bool {
	if node.IsNil() {
		return false
	}
	return node.Type == html.ElementNode
}

// IsText reports whether node is a text node.
func (node Node) IsText() bool {
	if node.IsNil() {
		return false
	}
	return node.Type == html.TextNode
}

// IsLeaf reports whether node has no child. It is false for a nil node.
func (node Node) IsLeaf() bool {
	if node.IsNil() {
		return false
	}
	return node.FirstChild == nil
}

// PrevElement returns the closest previous sibling that is an element, or a
// nil Node if there is none.
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

// NextElement returns the closest next sibling that is an element, or a nil
// Node if there is none.
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

// NextChildElement returns the first child that is an element, or a nil Node
// if there is none.
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

// GetAllChildElements returns the direct children of node that are elements,
// in document order.
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

// GetAttributes returns the names of the attributes of node, in document
// order.
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

// GetAttributeValue returns the value of the attribute named key (compared
// case-insensitively) and whether the attribute exists.
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

// GetNodeLink returns the href attribute of an <a> element or the src
// attribute of an <img> element, and whether it exists. Any other node gives
// an empty string and false.
func (node Node) GetNodeLink() (string, bool) {
	if node.IsNil() {
		return "", false
	}
	if node.Type == html.ElementNode &&
		strings.EqualFold(node.Data, TagA) {
		return node.GetAttributeValue(AttrHref)
	}
	if node.Type == html.ElementNode &&
		strings.EqualFold(node.Data, TagImg) {
		return node.GetAttributeValue(AttrSrc)
	}

	return "", false
}

// ConcatAllText joins with sep the text of all the descendants of node, in
// document order, skipping the blank ones. The content of <script> and
// <style> elements is regular text for the parser, so it is included.
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

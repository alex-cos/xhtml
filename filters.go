package xhtml

import (
	"regexp"
	"strings"
)

type FilterFunc func(f *Node) bool

func (filter FilterFunc) Not() FilterFunc {
	return func(node *Node) bool {
		return !filter(node)
	}
}

func (filter FilterFunc) And(right FilterFunc) FilterFunc {
	return func(node *Node) bool {
		return filter(node) && right(node)
	}
}

func (filter FilterFunc) Or(right FilterFunc) FilterFunc {
	return func(node *Node) bool {
		return filter(node) || right(node)
	}
}

func FilterAnyElement() FilterFunc {
	return func(node *Node) bool {
		return node.IsElement()
	}
}

func FilterAnyText() FilterFunc {
	return func(node *Node) bool {
		return node.IsText()
	}
}

func FilterAnyLeaf() FilterFunc {
	return func(node *Node) bool {
		return node.IsLeaf()
	}
}

func FilterByElement(name string) FilterFunc {
	return func(node *Node) bool {
		return node.IsElement() &&
			strings.EqualFold(node.Node.Data, name)
	}
}

func FilterTextByRegEx(regex *regexp.Regexp) FilterFunc {
	return func(node *Node) bool {
		return node.IsText() &&
			regex.MatchString(node.Node.Data)
	}
}

func FilterByAttribute(name, value string) FilterFunc {
	return func(node *Node) bool {
		if !node.IsElement() {
			return false
		}
		for _, a := range node.Attr {
			if strings.EqualFold(a.Key, name) &&
				a.Val == value {
				return true
			}
		}
		return false
	}
}

func FilterByClassName(className string) FilterFunc {
	return func(node *Node) bool {
		if !node.IsElement() {
			return false
		}
		for _, a := range node.Attr {
			if strings.EqualFold(a.Key, "class") {
				classes := strings.Split(a.Val, " ")
				for _, class := range classes {
					if strings.EqualFold(class, className) {
						return true
					}
				}
			}
		}
		return false
	}
}

func FilterByID(id string) FilterFunc {
	return func(node *Node) bool {
		if !node.IsElement() {
			return false
		}
		for _, a := range node.Attr {
			if strings.EqualFold(a.Key, "id") && strings.EqualFold(a.Val, id) {
				return true
			}
		}
		return false
	}
}

func FilterByRef(ref string) FilterFunc {
	return func(node *Node) bool {
		if !node.IsElement() {
			return false
		}
		for _, a := range node.Attr {
			if strings.EqualFold(a.Key, "ref") && strings.EqualFold(a.Val, ref) {
				return true
			}
		}
		return false
	}
}

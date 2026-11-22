package xhtml

import (
	"regexp"
	"strings"
)

// FilterFunc reports whether a node matches a criterion. Filters are passed to
// the Find* methods and can be combined with Not, And and Or.
type FilterFunc func(f Node) bool

// Not returns a filter that matches the nodes filter does not match.
func (filter FilterFunc) Not() FilterFunc {
	return func(node Node) bool {
		return !filter(node)
	}
}

// And returns a filter that matches the nodes both filter and right match.
func (filter FilterFunc) And(right FilterFunc) FilterFunc {
	return func(node Node) bool {
		return filter(node) && right(node)
	}
}

// Or returns a filter that matches the nodes filter or right matches.
func (filter FilterFunc) Or(right FilterFunc) FilterFunc {
	return func(node Node) bool {
		return filter(node) || right(node)
	}
}

// FilterAnyElement matches any element node.
func FilterAnyElement() FilterFunc {
	return func(node Node) bool {
		return node.IsElement()
	}
}

// FilterAnyText matches any text node.
func FilterAnyText() FilterFunc {
	return func(node Node) bool {
		return node.IsText()
	}
}

// FilterAnyLeaf matches any node that has no child.
func FilterAnyLeaf() FilterFunc {
	return func(node Node) bool {
		return node.IsLeaf()
	}
}

// FilterByElement matches elements whose tag name is name (compared
// case-insensitively).
func FilterByElement(name string) FilterFunc {
	return func(node Node) bool {
		return node.IsElement() &&
			strings.EqualFold(node.Data, name)
	}
}

// FilterTextByRegEx matches text nodes whose content matches regex.
func FilterTextByRegEx(regex *regexp.Regexp) FilterFunc {
	return func(node Node) bool {
		return node.IsText() &&
			regex.MatchString(node.Data)
	}
}

// FilterByAttribute matches elements that have the attribute name (compared
// case-insensitively) with exactly the given value.
func FilterByAttribute(name, value string) FilterFunc {
	return func(node Node) bool {
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

// FilterByClassName matches elements whose class attribute contains className
// as one of its white space separated classes (compared case-sensitively).
func FilterByClassName(className string) FilterFunc {
	return func(node Node) bool {
		if !node.IsElement() {
			return false
		}
		for _, a := range node.Attr {
			if strings.EqualFold(a.Key, AttrClass) {
				classes := strings.FieldsSeq(a.Val)
				for class := range classes {
					if class == className {
						return true
					}
				}
			}
		}
		return false
	}
}

// FilterByID matches elements whose id attribute is exactly id.
func FilterByID(id string) FilterFunc {
	return func(node Node) bool {
		if !node.IsElement() {
			return false
		}
		for _, a := range node.Attr {
			if strings.EqualFold(a.Key, AttrID) && a.Val == id {
				return true
			}
		}
		return false
	}
}

// FilterByRef matches elements whose ref attribute is ref (compared
// case-insensitively). ref is not a standard HTML attribute.
func FilterByRef(ref string) FilterFunc {
	return func(node Node) bool {
		if !node.IsElement() {
			return false
		}
		for _, a := range node.Attr {
			if strings.EqualFold(a.Key, AttrRef) && strings.EqualFold(a.Val, ref) {
				return true
			}
		}
		return false
	}
}

// FilterByAttributeExists matches elements that have the attribute name,
// whatever its value (the name is compared case-insensitively).
func FilterByAttributeExists(name string) FilterFunc {
	return filterByAttributeValue(name, func(string) bool { return true })
}

// FilterByAttributePrefix matches elements whose attribute name has a value
// starting with prefix (the value is compared case-sensitively).
func FilterByAttributePrefix(name, prefix string) FilterFunc {
	return filterByAttributeValue(name, func(val string) bool {
		return strings.HasPrefix(val, prefix)
	})
}

// FilterByAttributeSuffix matches elements whose attribute name has a value
// ending with suffix (the value is compared case-sensitively).
func FilterByAttributeSuffix(name, suffix string) FilterFunc {
	return filterByAttributeValue(name, func(val string) bool {
		return strings.HasSuffix(val, suffix)
	})
}

// FilterByAttributeContains matches elements whose attribute name has a value
// containing substr (the value is compared case-sensitively).
func FilterByAttributeContains(name, substr string) FilterFunc {
	return filterByAttributeValue(name, func(val string) bool {
		return strings.Contains(val, substr)
	})
}

// FilterByClassNames matches elements that carry all the given classes,
// in any order. Without any class, it matches any element.
func FilterByClassNames(classes ...string) FilterFunc {
	filter := FilterAnyElement()
	for _, class := range classes {
		filter = filter.And(FilterByClassName(class))
	}

	return filter
}

// FilterByText matches text nodes whose content, once surrounding white
// space is trimmed, is exactly text.
func FilterByText(text string) FilterFunc {
	return func(node Node) bool {
		return node.IsText() &&
			strings.TrimSpace(node.Data) == text
	}
}

// FilterByTextContains matches text nodes whose content contains substr.
func FilterByTextContains(substr string) FilterFunc {
	return func(node Node) bool {
		return node.IsText() &&
			strings.Contains(node.Data, substr)
	}
}

func filterByAttributeValue(name string, match func(string) bool) FilterFunc {
	return func(node Node) bool {
		if !node.IsElement() {
			return false
		}
		for _, a := range node.Attr {
			if strings.EqualFold(a.Key, name) && match(a.Val) {
				return true
			}
		}

		return false
	}
}

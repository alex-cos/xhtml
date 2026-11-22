// Package xhtml is a small wrapper around golang.org/x/net/html that makes it
// easier to traverse and query an HTML document.
//
// A Node wraps an *html.Node and provides search methods (FindFirstNode,
// FindLastNode, FindAllNodes, FindNthNode), navigation helpers and text and
// attribute accessors. The search methods take a FilterFunc, and filters can
// be combined with Not, And and Or:
//
//	root, err := xhtml.Parse(r)
//	if err != nil {
//		return err
//	}
//
//	filter := xhtml.FilterByElement(xhtml.TagA).And(xhtml.FilterByClassName("external"))
//	for _, link := range root.FindAllNodes(filter) {
//		href, _ := link.GetAttributeValue(xhtml.AttrHref)
//		fmt.Println(href, link.ConcatAllText(" "))
//	}
package xhtml

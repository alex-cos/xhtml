package xhtml_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/alex-cos/xhtml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterAnyElement(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterAnyElement())
	assert.NotEmpty(t, nodes)

	for _, n := range nodes {
		assert.True(t, n.IsElement())
	}
}

func TestFilterAnyText(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterAnyText())
	assert.NotEmpty(t, nodes)

	for _, n := range nodes {
		assert.True(t, n.IsText())
	}
}

func TestFilterAnyLeaf(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterAnyLeaf())
	assert.NotEmpty(t, nodes)

	for _, n := range nodes {
		assert.True(t, n.IsLeaf())
	}
}

func TestFilterByElement(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterByElement(xhtml.TagH2))
	assert.Len(t, nodes, 2)

	nodes = body.FindAllNodes(xhtml.FilterByElement(xhtml.TagLi))
	assert.Len(t, nodes, 6)

	nodes = body.FindAllNodes(xhtml.FilterByElement("nonexistent"))
	assert.Empty(t, nodes)
}

func TestFilterByAttribute(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterByAttribute("href", "#"))
	assert.NotEmpty(t, nodes)

	nodes = body.FindAllNodes(xhtml.FilterByAttribute("data-role", "container"))
	assert.Len(t, nodes, 1)

	nodes = body.FindAllNodes(xhtml.FilterByAttribute("nonexistent", "value"))
	assert.Empty(t, nodes)
}

func TestFilterByClassName(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterByClassName("fas"))
	assert.Len(t, nodes, 3)

	nodes = body.FindAllNodes(xhtml.FilterByClassName("image-list"))
	assert.Len(t, nodes, 1)

	nodes = body.FindAllNodes(xhtml.FilterByClassName("highlight"))
	assert.Len(t, nodes, 3)

	nodes = body.FindAllNodes(xhtml.FilterByClassName("spaced"))
	assert.Len(t, nodes, 1)

	nodes = body.FindAllNodes(xhtml.FilterByClassName("nonexistent"))
	assert.Empty(t, nodes)
}

func TestFilterByID(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	node := body.FindFirstNode(xhtml.FilterByID("images"))
	assert.False(t, node.IsNil())
	assert.Equal(t, xhtml.TagUl, node.GetData())

	node = body.FindFirstNode(xhtml.FilterByID("nonexistent"))
	assert.True(t, node.IsNil())
}

func TestFilterByRef(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	node := body.FindFirstNode(xhtml.FilterByRef("mainTitle"))
	assert.False(t, node.IsNil())
	assert.Equal(t, xhtml.TagH1, node.GetData())

	node = body.FindFirstNode(xhtml.FilterByRef("nonexistent"))
	assert.True(t, node.IsNil())
}

func TestFilterTextByRegEx(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	regex := regexp.MustCompile(`^Galerie.*$`)
	nodes := body.FindAllNodes(xhtml.FilterTextByRegEx(regex))
	assert.NotEmpty(t, nodes)

	regex = regexp.MustCompile(`^Bienvenue.*$`)
	nodes = body.FindAllNodes(xhtml.FilterTextByRegEx(regex))
	assert.NotEmpty(t, nodes)

	regex = regexp.MustCompile(`^zzzznotfoundzzzz$`)
	nodes = body.FindAllNodes(xhtml.FilterTextByRegEx(regex))
	assert.Empty(t, nodes)
}

func TestFilterNot(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	notBr := xhtml.FilterByElement(xhtml.TagBr).Not()
	nodes := body.FindAllNodes(notBr)

	for _, n := range nodes {
		if n.IsElement() {
			assert.NotEqual(t, xhtml.TagBr, n.GetData())
		}
	}
}

func TestFilterAnd(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	filter := xhtml.FilterByElement(xhtml.TagSpan).And(xhtml.FilterByClassName("highlight"))
	nodes := body.FindAllNodes(filter)
	assert.Len(t, nodes, 3)

	filter = xhtml.FilterByElement(xhtml.TagSpan).And(xhtml.FilterByClassName("nonexistent"))
	nodes = body.FindAllNodes(filter)
	assert.Empty(t, nodes)
}

func TestFilterOr(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	filter := xhtml.FilterByElement(xhtml.TagH1).Or(xhtml.FilterByElement(xhtml.TagH2))
	nodes := body.FindAllNodes(filter)
	assert.Len(t, nodes, 3)

	for _, n := range nodes {
		assert.True(t, n.IsElement())
		data := n.GetData()
		assert.True(t, data == "h1" || data == "h2")
	}
}

func TestFilterCompositionChained(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	filter := xhtml.FilterByElement(xhtml.TagSpan).
		And(xhtml.FilterByClassName("highlight")).
		And(xhtml.FilterByAttribute("data-id", "1"))

	nodes := body.FindAllNodes(filter)
	assert.Len(t, nodes, 1)
}

func TestFilterOnNilNode(t *testing.T) {
	t.Parallel()

	nilNode := xhtml.NilNode()

	assert.False(t, xhtml.FilterAnyElement()(nilNode))
	assert.False(t, xhtml.FilterAnyText()(nilNode))
	assert.False(t, xhtml.FilterAnyLeaf()(nilNode))
	assert.False(t, xhtml.FilterByElement(xhtml.TagDiv)(nilNode))
	assert.False(t, xhtml.FilterByClassName("test")(nilNode))
	assert.False(t, xhtml.FilterByID("test")(nilNode))
	assert.False(t, xhtml.FilterByRef("test")(nilNode))
	assert.False(t, xhtml.FilterByAttribute("key", "val")(nilNode))
	assert.False(t, xhtml.FilterByAttributeExists("key")(nilNode))
	assert.False(t, xhtml.FilterByAttributePrefix("key", "")(nilNode))
	assert.False(t, xhtml.FilterByAttributeSuffix("key", "")(nilNode))
	assert.False(t, xhtml.FilterByAttributeContains("key", "")(nilNode))
	assert.False(t, xhtml.FilterByClassNames("test")(nilNode))
	assert.False(t, xhtml.FilterByClassNames()(nilNode))
	assert.False(t, xhtml.FilterByText("test")(nilNode))
	assert.False(t, xhtml.FilterByTextContains("test")(nilNode))

	regex := regexp.MustCompile(`.*`)
	assert.False(t, xhtml.FilterTextByRegEx(regex)(nilNode))
}

func TestFilterByAttributeExists(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterByAttributeExists("data-id"))
	assert.Len(t, nodes, 3)

	nodes = body.FindAllNodes(xhtml.FilterByAttributeExists("DATA-ROLE"))
	assert.Len(t, nodes, 1)

	nodes = body.FindAllNodes(xhtml.FilterByAttributeExists("nonexistent"))
	assert.Empty(t, nodes)
}

func TestFilterByAttributePrefix(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterByAttributePrefix("src", "https://"))
	assert.Len(t, nodes, 3)

	nodes = body.FindAllNodes(xhtml.FilterByAttributePrefix("src", "HTTPS://"))
	assert.Empty(t, nodes)

	nodes = body.FindAllNodes(xhtml.FilterByAttributePrefix("nonexistent", ""))
	assert.Empty(t, nodes)
}

func TestFilterByAttributeSuffix(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterByAttributeSuffix("src", "/150"))
	assert.Len(t, nodes, 3)

	nodes = body.FindAllNodes(xhtml.FilterByAttributeSuffix("src", "/151"))
	assert.Empty(t, nodes)
}

func TestFilterByAttributeContains(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterByAttributeContains("src", "placeholder"))
	assert.Len(t, nodes, 3)

	nodes = body.FindAllNodes(xhtml.FilterByAttributeContains("alt", "Image"))
	assert.Len(t, nodes, 3)

	nodes = body.FindAllNodes(xhtml.FilterByAttributeContains("alt", "zzz"))
	assert.Empty(t, nodes)

	// Combined with an element filter
	filter := xhtml.FilterByElement(xhtml.TagA).And(xhtml.FilterByAttributeContains("href", "#"))
	nodes = body.FindAllNodes(filter)
	assert.NotEmpty(t, nodes)
}

func TestFilterByClassNames(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterByClassNames("multi", "classes"))
	assert.Len(t, nodes, 1)

	nodes = body.FindAllNodes(xhtml.FilterByClassNames("classes", "spaced", "multi"))
	assert.Len(t, nodes, 1)

	nodes = body.FindAllNodes(xhtml.FilterByClassNames("fas", "fa-home"))
	assert.Len(t, nodes, 1)

	nodes = body.FindAllNodes(xhtml.FilterByClassNames("multi", "nonexistent"))
	assert.Empty(t, nodes)

	nodes = body.FindAllNodes(xhtml.FilterByClassNames())
	assert.Equal(t, body.FindAllNodes(xhtml.FilterAnyElement()), nodes)
}

func TestFilterByText(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterByText("Text 1"))
	require.Len(t, nodes, 1)
	assert.Equal(t, xhtml.TagSpan, nodes[0].Parent.Data)

	nodes = body.FindAllNodes(xhtml.FilterByText("Text"))
	assert.Empty(t, nodes)

	// Surrounding white space is ignored
	doc, err := xhtml.Parse(strings.NewReader("<p>\n  hello  \n</p>"))
	require.NoError(t, err)
	assert.Len(t, doc.FindAllNodes(xhtml.FilterByText("hello")), 1)
}

func TestFilterByTextContains(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterByTextContains("Text"))
	assert.Len(t, nodes, 3)

	nodes = body.FindAllNodes(xhtml.FilterByTextContains("zzzznotfoundzzzz"))
	assert.Empty(t, nodes)
}

package xhtml_test

import (
	"regexp"
	"testing"

	"github.com/alex-cos/xhtml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterAnyElement(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
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
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
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
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
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
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	nodes := body.FindAllNodes(xhtml.FilterByElement(xhtml.H2))
	assert.Len(t, nodes, 2)

	nodes = body.FindAllNodes(xhtml.FilterByElement(xhtml.LI))
	assert.Len(t, nodes, 6)

	nodes = body.FindAllNodes(xhtml.FilterByElement("nonexistent"))
	assert.Empty(t, nodes)
}

func TestFilterByAttribute(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
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
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
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
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	node := body.FindFirstNode(xhtml.FilterByID("images"))
	assert.False(t, node.IsNil())
	assert.Equal(t, xhtml.UL, node.GetData())

	node = body.FindFirstNode(xhtml.FilterByID("nonexistent"))
	assert.True(t, node.IsNil())
}

func TestFilterByRef(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	node := body.FindFirstNode(xhtml.FilterByRef("mainTitle"))
	assert.False(t, node.IsNil())
	assert.Equal(t, xhtml.H1, node.GetData())

	node = body.FindFirstNode(xhtml.FilterByRef("nonexistent"))
	assert.True(t, node.IsNil())
}

func TestFilterTextByRegEx(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
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
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	notBr := xhtml.FilterByElement(xhtml.BR).Not()
	nodes := body.FindAllNodes(notBr)

	for _, n := range nodes {
		if n.IsElement() {
			assert.NotEqual(t, xhtml.BR, n.GetData())
		}
	}
}

func TestFilterAnd(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	filter := xhtml.FilterByElement(xhtml.SPAN).And(xhtml.FilterByClassName("highlight"))
	nodes := body.FindAllNodes(filter)
	assert.Len(t, nodes, 3)

	filter = xhtml.FilterByElement(xhtml.SPAN).And(xhtml.FilterByClassName("nonexistent"))
	nodes = body.FindAllNodes(filter)
	assert.Empty(t, nodes)
}

func TestFilterOr(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	filter := xhtml.FilterByElement(xhtml.H1).Or(xhtml.FilterByElement(xhtml.H2))
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
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	filter := xhtml.FilterByElement(xhtml.SPAN).
		And(xhtml.FilterByClassName("highlight")).
		And(xhtml.FilterByAttribute("data-id", "1"))

	nodes := body.FindAllNodes(filter)
	assert.Len(t, nodes, 1)
}

func TestFilterOnNilNode(t *testing.T) {
	t.Parallel()

	nilNode := xhtml.NilNode()

	assert.False(t, xhtml.FilterAnyElement()(&nilNode))
	assert.False(t, xhtml.FilterAnyText()(&nilNode))
	assert.False(t, xhtml.FilterAnyLeaf()(&nilNode))
	assert.False(t, xhtml.FilterByElement(xhtml.DIV)(&nilNode))
	assert.False(t, xhtml.FilterByClassName("test")(&nilNode))
	assert.False(t, xhtml.FilterByID("test")(&nilNode))
	assert.False(t, xhtml.FilterByRef("test")(&nilNode))
	assert.False(t, xhtml.FilterByAttribute("key", "val")(&nilNode))

	regex := regexp.MustCompile(`.*`)
	assert.False(t, xhtml.FilterTextByRegEx(regex)(&nilNode))
}

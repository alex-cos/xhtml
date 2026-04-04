package xhtml_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/alex-cos/xhtml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func parseTestFile(t *testing.T) xhtml.Node {
	t.Helper()
	testfile := filepath.Join("testdata", "test1.html")
	data, err := os.ReadFile(testfile)
	require.NoError(t, err)
	doc, err := html.Parse(bytes.NewReader(data))
	require.NoError(t, err)
	return xhtml.NewNode(doc)
}

func TestNilNode(t *testing.T) {
	t.Parallel()

	node := xhtml.NilNode()
	assert.True(t, node.IsNil())
}

func TestNewNode(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	assert.False(t, root.IsNil())
}

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("valid HTML", func(t *testing.T) {
		t.Parallel()
		raw := `<html><body><p>Hello</p></body></html>`
		node, err := xhtml.Parse(bytes.NewReader([]byte(raw)))
		assert.NoError(t, err)
		assert.False(t, node.IsNil())
	})

	t.Run("invalid HTML returns error", func(t *testing.T) {
		t.Parallel()
		_, err := xhtml.Parse(&brokenReader{})
		assert.Error(t, err)
	})
}

type brokenReader struct{}

func (b *brokenReader) Read(p []byte) (int, error) {
	return 0, os.ErrInvalid
}

func TestNodeIsNil(t *testing.T) {
	t.Parallel()

	assert.True(t, xhtml.NilNode().IsNil())

	root := parseTestFile(t)
	assert.False(t, root.IsNil())
}

func TestNodeGetData(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	assert.Equal(t, "body", body.GetData())

	assert.Empty(t, xhtml.NilNode().GetData())
}

func TestNodeIsElement(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	assert.True(t, body.IsElement())

	textNodes := root.FindAllNodes(xhtml.FilterAnyText())
	require.NotEmpty(t, textNodes)
	assert.False(t, textNodes[0].IsElement())

	assert.False(t, xhtml.NilNode().IsElement())
}

func TestNodeIsText(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	textNodes := root.FindAllNodes(xhtml.FilterAnyText())
	require.NotEmpty(t, textNodes)
	assert.True(t, textNodes[0].IsText())

	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	assert.False(t, body.IsText())

	assert.False(t, xhtml.NilNode().IsText())
}

func TestNodeIsLeaf(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	br := root.FindFirstNode(xhtml.FilterByElement(xhtml.BR))
	assert.True(t, br.IsLeaf())

	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	assert.False(t, body.IsLeaf())

	assert.False(t, xhtml.NilNode().IsLeaf())
}

func TestNodePrevElement(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	secondSpan := root.FindNthNode(2, xhtml.FilterByElement(xhtml.SPAN))
	require.False(t, secondSpan.IsNil())

	prev := secondSpan.PrevElement()
	assert.False(t, prev.IsNil())
	assert.Equal(t, xhtml.SPAN, prev.GetData())

	thirdSpan := root.FindNthNode(3, xhtml.FilterByElement(xhtml.SPAN))
	require.False(t, thirdSpan.IsNil())

	prev = thirdSpan.PrevElement()
	assert.False(t, prev.IsNil())
	assert.Equal(t, xhtml.HR, prev.GetData())

	firstSpan := root.FindNthNode(1, xhtml.FilterByElement(xhtml.SPAN))
	require.False(t, firstSpan.IsNil())

	prev = firstSpan.PrevElement()
	assert.True(t, prev.IsNil())

	assert.True(t, xhtml.NilNode().PrevElement().IsNil())
}

func TestNodeNextElement(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	firstSpan := root.FindNthNode(1, xhtml.FilterByElement(xhtml.SPAN))
	require.False(t, firstSpan.IsNil())

	next := firstSpan.NextElement()
	assert.False(t, next.IsNil())
	assert.Equal(t, xhtml.SPAN, next.GetData())

	thirdSpan := root.FindNthNode(3, xhtml.FilterByElement(xhtml.SPAN))
	next = thirdSpan.NextElement()
	assert.True(t, next.IsNil())

	assert.True(t, xhtml.NilNode().NextElement().IsNil())
}

func TestNodeNextChildElement(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	firstChild := body.NextChildElement()
	assert.False(t, firstChild.IsNil())
	assert.Equal(t, xhtml.HEADER, firstChild.GetData())

	assert.True(t, xhtml.NilNode().NextChildElement().IsNil())

	br := root.FindFirstNode(xhtml.FilterByElement(xhtml.BR))
	assert.True(t, br.NextChildElement().IsNil())
}

func TestNodeGetAllChildElements(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	children := body.GetAllChildElements()
	assert.Len(t, children, 4)
	assert.Equal(t, xhtml.HEADER, children[0].GetData())
	assert.Equal(t, xhtml.MAIN, children[1].GetData())
	assert.Equal(t, xhtml.DIV, children[2].GetData())
	assert.Equal(t, xhtml.FOOTER, children[3].GetData())

	assert.Empty(t, xhtml.NilNode().GetAllChildElements())
}

func TestNodeGetAttributes(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	img := root.FindFirstNode(xhtml.FilterByElement(xhtml.IMG))
	require.False(t, img.IsNil())

	attrs := img.GetAttributes()
	assert.Contains(t, attrs, "src")
	assert.Contains(t, attrs, "alt")

	assert.Empty(t, xhtml.NilNode().GetAttributes())
}

func TestNodeGetAttributeValue(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	img := root.FindFirstNode(xhtml.FilterByElement(xhtml.IMG))
	require.False(t, img.IsNil())

	val, ok := img.GetAttributeValue("src")
	assert.True(t, ok)
	assert.Equal(t, "https://via.placeholder.com/150", val)

	val, ok = img.GetAttributeValue("Src")
	assert.True(t, ok)
	assert.Equal(t, "https://via.placeholder.com/150", val)

	val, ok = img.GetAttributeValue("nonexistent")
	assert.False(t, ok)
	assert.Empty(t, val)

	val, ok = xhtml.NilNode().GetAttributeValue("src")
	assert.False(t, ok)
	assert.Empty(t, val)
}

func TestNodeGetNodeLink(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)

	link := root.FindFirstNode(xhtml.FilterByElement(xhtml.A))
	require.False(t, link.IsNil())

	href, ok := link.GetNodeLink()
	assert.True(t, ok)
	assert.Equal(t, "#", href)

	img := root.FindFirstNode(xhtml.FilterByElement(xhtml.IMG))
	require.False(t, img.IsNil())

	src, ok := img.GetNodeLink()
	assert.True(t, ok)
	assert.Equal(t, "https://via.placeholder.com/150", src)

	div := root.FindFirstNode(xhtml.FilterByElement(xhtml.DIV))
	require.False(t, div.IsNil())

	linkVal, ok := div.GetNodeLink()
	assert.False(t, ok)
	assert.Empty(t, linkVal)

	_, ok = xhtml.NilNode().GetNodeLink()
	assert.False(t, ok)
}

func TestNodeConcatAllText(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	text := body.ConcatAllText("\n")
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "Bienvenue")
	assert.Contains(t, text, "Galerie")

	assert.Empty(t, xhtml.NilNode().ConcatAllText("\n"))
}

func TestNodeString(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	s := root.String()
	assert.NotEmpty(t, s)
}

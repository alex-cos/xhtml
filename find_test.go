package xhtml_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/alex-cos/xhtml"
	"github.com/kjk/common/assert"
	"golang.org/x/net/html"
)

func TestFind(t *testing.T) {
	t.Parallel()

	testfile := filepath.Join("testdata", "test1.html")

	data, err := os.ReadFile(testfile)
	assert.NoError(t, err)

	doc, err := html.Parse(bytes.NewReader(data))
	assert.NoError(t, err)

	root := xhtml.NewNode(doc)
	assert.NotNil(t, root)

	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	assert.NotNil(t, body)
	if !testing.Short() {
		fmt.Printf("body = %+v\n", body)
	}
	assert.Equal(t, "body", body.Data)

	node := body.FindFirstNode(xhtml.FilterByRef("mainTitle"))
	assert.NotNil(t, node)
	if !testing.Short() {
		fmt.Printf("node = %+v\n", node)
	}

	nodes := body.FindAllNodes(xhtml.FilterByElement(xhtml.H2))
	assert.Len(t, nodes, 2)
	if !testing.Short() {
		fmt.Printf("nodes = %+v\n", nodes)
	}
	assert.Equal(t, "h2", nodes[0].Data)
	assert.Equal(t, "h2", nodes[1].Data)

	nodes = body.FindAllNodes(xhtml.FilterByClassName("fas"))
	assert.Len(t, nodes, 3)
	if !testing.Short() {
		fmt.Printf("nodes = %+v\n", nodes)
	}

	node = body.FindLastNode(xhtml.FilterByElement(xhtml.A))
	assert.NotNil(t, node)
	if !testing.Short() {
		fmt.Printf("node = %+v\n", node)
	}
	assert.Equal(t, "a", node.Data)

	node = node.NextChildElement()
	assert.NotNil(t, node)
	if !testing.Short() {
		fmt.Printf("node = %+v\n", node)
	}
	assert.Equal(t, "img", node.Data)

	link, _ := node.GetNodeLink()
	assert.Equal(t, "https://via.placeholder.com/150", link)

	keys := node.GetAttributes()
	assert.Equal(t, []string{"src", "alt"}, keys)

	alt, ok := node.GetAttributeValue("Alt")
	assert.True(t, ok)
	assert.Equal(t, "Image 3", alt)
	if !testing.Short() {
		fmt.Printf("alt = %+v\n", alt)
	}

	regex := regexp.MustCompile(`^Galerie.*$`)

	nodes = body.FindAllNodes(xhtml.FilterTextByRegEx(regex))
	assert.True(t, ok)
	if !testing.Short() {
		fmt.Printf("nodes = %+v\n", nodes)
	}

	node = body.FindFirstNode(xhtml.FilterByID("images"))
	assert.NotNil(t, node)
	if !testing.Short() {
		fmt.Printf("node = %+v\n", node)
	}

	node = body.FindFirstNode(xhtml.FilterByClassName("image-list"))
	assert.NotNil(t, node)
	if !testing.Short() {
		fmt.Printf("node = %+v\n", node)
	}

	node = node.NextChildElement().NextElement().NextElement()
	assert.NotNil(t, node)
	if !testing.Short() {
		fmt.Printf("node = %+v\n", node)
	}
	assert.Equal(t, "li", node.Data)

	node = node.PrevElement()
	assert.NotNil(t, node)
	if !testing.Short() {
		fmt.Printf("node = %+v\n", node)
	}
	assert.Equal(t, "li", node.Data)

	node = body.FindNthNode(5, xhtml.FilterByElement(xhtml.LI))
	assert.NotNil(t, node)
	if !testing.Short() {
		fmt.Printf("node = %+v\n", node)
	}
	assert.Equal(t, "li", node.Data)

	alt, ok = node.FindFirstNode(xhtml.FilterByElement(xhtml.IMG)).GetAttributeValue("Alt")
	assert.True(t, ok)
	assert.Equal(t, "Image 2", alt)
	if !testing.Short() {
		fmt.Printf("alt = %+v\n", alt)
	}

	node = body.FindNthNode(-2, xhtml.FilterByElement(xhtml.LI))
	assert.True(t, node.IsNil())

	node = body.FindNthNode(1, xhtml.FilterByClassName("xxxxx"))
	assert.True(t, node.IsNil())

	text := body.ConcatAllText("\n")
	assert.Len(t, text, 141)
	if !testing.Short() {
		fmt.Printf("text = %+v\n", text)
	}
}

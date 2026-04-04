package xhtml_test

import (
	"testing"

	"github.com/alex-cos/xhtml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindFirstNode(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	t.Run("finds existing element", func(t *testing.T) {
		t.Parallel()
		node := body.FindFirstNode(xhtml.FilterByElement(xhtml.H2))
		assert.False(t, node.IsNil())
		assert.Equal(t, xhtml.H2, node.GetData())
	})

	t.Run("finds by ref attribute", func(t *testing.T) {
		t.Parallel()
		node := body.FindFirstNode(xhtml.FilterByRef("mainTitle"))
		assert.False(t, node.IsNil())
		assert.Equal(t, xhtml.H1, node.GetData())
	})

	t.Run("finds by id", func(t *testing.T) {
		t.Parallel()
		node := body.FindFirstNode(xhtml.FilterByID("images"))
		assert.False(t, node.IsNil())
		assert.Equal(t, xhtml.UL, node.GetData())
	})

	t.Run("finds by class name", func(t *testing.T) {
		t.Parallel()
		node := body.FindFirstNode(xhtml.FilterByClassName("image-list"))
		assert.False(t, node.IsNil())
		assert.Equal(t, xhtml.UL, node.GetData())
	})

	t.Run("returns nil for nonexistent", func(t *testing.T) {
		t.Parallel()
		node := body.FindFirstNode(xhtml.FilterByClassName("nonexistent"))
		assert.True(t, node.IsNil())
	})

	t.Run("returns nil on nil node", func(t *testing.T) {
		t.Parallel()
		node := xhtml.NilNode().FindFirstNode(xhtml.FilterAnyElement())
		assert.True(t, node.IsNil())
	})
}

func TestFindLastNode(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	t.Run("finds last link", func(t *testing.T) {
		t.Parallel()
		node := body.FindLastNode(xhtml.FilterByElement(xhtml.A))
		assert.False(t, node.IsNil())
		assert.Equal(t, xhtml.A, node.GetData())
	})

	t.Run("finds last li", func(t *testing.T) {
		t.Parallel()
		node := body.FindLastNode(xhtml.FilterByElement(xhtml.LI))
		assert.False(t, node.IsNil())
		assert.Equal(t, xhtml.LI, node.GetData())
	})

	t.Run("returns nil for nonexistent", func(t *testing.T) {
		t.Parallel()
		node := body.FindLastNode(xhtml.FilterByClassName("nonexistent"))
		assert.True(t, node.IsNil())
	})

	t.Run("returns nil on nil node", func(t *testing.T) {
		t.Parallel()
		node := xhtml.NilNode().FindLastNode(xhtml.FilterAnyElement())
		assert.True(t, node.IsNil())
	})

	t.Run("finds self when no children match", func(t *testing.T) {
		t.Parallel()
		h2 := body.FindFirstNode(xhtml.FilterByElement(xhtml.H2))
		require.False(t, h2.IsNil())

		node := h2.FindLastNode(xhtml.FilterByElement(xhtml.H2))
		assert.False(t, node.IsNil())
		assert.Equal(t, xhtml.H2, node.GetData())
	})
}

func TestFindAllNodes(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	t.Run("finds all h2 elements", func(t *testing.T) {
		t.Parallel()
		nodes := body.FindAllNodes(xhtml.FilterByElement(xhtml.H2))
		assert.Len(t, nodes, 2)
		for _, n := range nodes {
			assert.Equal(t, xhtml.H2, n.GetData())
		}
	})

	t.Run("finds all li elements", func(t *testing.T) {
		t.Parallel()
		nodes := body.FindAllNodes(xhtml.FilterByElement(xhtml.LI))
		assert.Len(t, nodes, 6)
	})

	t.Run("finds all links", func(t *testing.T) {
		t.Parallel()
		nodes := body.FindAllNodes(xhtml.FilterByElement(xhtml.A))
		assert.Len(t, nodes, 6)
	})

	t.Run("finds all images", func(t *testing.T) {
		t.Parallel()
		nodes := body.FindAllNodes(xhtml.FilterByElement(xhtml.IMG))
		assert.Len(t, nodes, 3)
	})

	t.Run("finds by class with multiple spaces", func(t *testing.T) {
		t.Parallel()
		nodes := body.FindAllNodes(xhtml.FilterByClassName("spaced"))
		assert.Len(t, nodes, 1)
	})

	t.Run("finds all elements with class highlight", func(t *testing.T) {
		t.Parallel()
		nodes := body.FindAllNodes(xhtml.FilterByClassName("highlight"))
		assert.Len(t, nodes, 3)
	})

	t.Run("returns empty slice for nonexistent", func(t *testing.T) {
		t.Parallel()
		nodes := body.FindAllNodes(xhtml.FilterByClassName("nonexistent"))
		assert.Empty(t, nodes)
	})

	t.Run("returns empty slice on nil node", func(t *testing.T) {
		t.Parallel()
		nodes := xhtml.NilNode().FindAllNodes(xhtml.FilterAnyElement())
		assert.Empty(t, nodes)
	})

	t.Run("finds with composed filter", func(t *testing.T) {
		t.Parallel()
		filter := xhtml.FilterByElement(xhtml.A).And(xhtml.FilterByClassName("external"))
		nodes := body.FindAllNodes(filter)
		assert.Empty(t, nodes)
	})
}

func TestFindNthNode(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	t.Run("finds first li", func(t *testing.T) {
		t.Parallel()
		node := body.FindNthNode(1, xhtml.FilterByElement(xhtml.LI))
		assert.False(t, node.IsNil())
		assert.Equal(t, xhtml.LI, node.GetData())
	})

	t.Run("finds fifth li", func(t *testing.T) {
		t.Parallel()
		node := body.FindNthNode(5, xhtml.FilterByElement(xhtml.LI))
		assert.False(t, node.IsNil())
		assert.Equal(t, xhtml.LI, node.GetData())
	})

	t.Run("finds last li (sixth)", func(t *testing.T) {
		t.Parallel()
		node := body.FindNthNode(6, xhtml.FilterByElement(xhtml.LI))
		assert.False(t, node.IsNil())
		assert.Equal(t, xhtml.LI, node.GetData())
	})

	t.Run("returns nil for position beyond count", func(t *testing.T) {
		t.Parallel()
		node := body.FindNthNode(100, xhtml.FilterByElement(xhtml.LI))
		assert.True(t, node.IsNil())
	})

	t.Run("returns nil for zero position", func(t *testing.T) {
		t.Parallel()
		node := body.FindNthNode(0, xhtml.FilterByElement(xhtml.LI))
		assert.True(t, node.IsNil())
	})

	t.Run("returns nil for negative position", func(t *testing.T) {
		t.Parallel()
		node := body.FindNthNode(-2, xhtml.FilterByElement(xhtml.LI))
		assert.True(t, node.IsNil())
	})

	t.Run("returns nil for nonexistent class", func(t *testing.T) {
		t.Parallel()
		node := body.FindNthNode(1, xhtml.FilterByClassName("xxxxx"))
		assert.True(t, node.IsNil())
	})

	t.Run("returns nil on nil node", func(t *testing.T) {
		t.Parallel()
		node := xhtml.NilNode().FindNthNode(1, xhtml.FilterAnyElement())
		assert.True(t, node.IsNil())
	})
}

func TestFindIntegration(t *testing.T) {
	t.Parallel()

	root := parseTestFile(t)
	body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
	require.False(t, body.IsNil())

	t.Run("navigation chain from found node", func(t *testing.T) {
		t.Parallel()
		node := body.FindFirstNode(xhtml.FilterByClassName("image-list"))
		require.False(t, node.IsNil())

		child := node.NextChildElement()
		assert.False(t, child.IsNil())
		assert.Equal(t, xhtml.LI, child.GetData())

		next := child.NextElement()
		assert.False(t, next.IsNil())
		assert.Equal(t, xhtml.LI, next.GetData())

		prev := next.PrevElement()
		assert.False(t, prev.IsNil())
		assert.Equal(t, xhtml.LI, prev.GetData())
	})

	t.Run("extract link from nested structure", func(t *testing.T) {
		t.Parallel()
		li := body.FindNthNode(5, xhtml.FilterByElement(xhtml.LI))
		require.False(t, li.IsNil())

		img := li.FindFirstNode(xhtml.FilterByElement(xhtml.IMG))
		require.False(t, img.IsNil())

		alt, ok := img.GetAttributeValue("Alt")
		assert.True(t, ok)
		assert.Equal(t, "Image 2", alt)
	})

	t.Run("concat all text from body", func(t *testing.T) {
		t.Parallel()
		text := body.ConcatAllText("\n")
		assert.Contains(t, text, "Bienvenue")
		assert.Contains(t, text, "Galerie")
		assert.Contains(t, text, "Accueil")
	})

	t.Run("get node link from image", func(t *testing.T) {
		t.Parallel()
		lastLink := body.FindLastNode(xhtml.FilterByElement(xhtml.A))
		require.False(t, lastLink.IsNil())

		img := lastLink.NextChildElement()
		require.False(t, img.IsNil())
		assert.Equal(t, xhtml.IMG, img.GetData())

		link, ok := img.GetNodeLink()
		assert.True(t, ok)
		assert.Equal(t, "https://via.placeholder.com/150", link)
	})

	t.Run("get attributes from image", func(t *testing.T) {
		t.Parallel()
		img := body.FindFirstNode(xhtml.FilterByElement(xhtml.IMG))
		require.False(t, img.IsNil())

		keys := img.GetAttributes()
		assert.Contains(t, keys, "src")
		assert.Contains(t, keys, "alt")
	})
}

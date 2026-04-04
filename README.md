# xhtml

A small Go library for traversing and querying an HTML DOM with ease.

`xhtml` is an ergonomic wrapper around `golang.org/x/net/html` that adds navigation methods (`FindFirstNode`, `FindAllNodes`, `FindNthNode`, `FindLastNode`), composable filters, and utilities for extracting text and attributes.

## Installation

```bash
go get github.com/alex-cos/xhtml
```

## Basic Usage

### Parsing an HTML document

```go
package main

import (
  "bytes"
  "fmt"
  "log"

  "github.com/alex-cos/xhtml"
  "golang.org/x/net/html"
)

func main() {
  raw := `<html><body><h1>Hello</h1></body></html>`

  doc, err := html.Parse(bytes.NewReader([]byte(raw)))
  if err != nil {
    log.Fatal(err)
  }

  root := xhtml.NewNode(doc)
  // root is ready for querying
}
```

### Finding the first matching node

```go
body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))
title := root.FindFirstNode(xhtml.FilterByElement(xhtml.H1))

fmt.Println(title.GetData()) // "h1"
```

### Finding all matching nodes

```go
links := root.FindAllNodes(xhtml.FilterByElement(xhtml.A))

for _, link := range links {
  href, ok := link.GetAttributeValue("href")
  if ok {
    fmt.Println(href)
  }
}
```

### Finding the Nth node

```go
// Find the 3rd <li> element
thirdLi := root.FindNthNode(3, xhtml.FilterByElement(xhtml.LI))

if !thirdLi.IsNil() {
  fmt.Println("Found!")
}
```

### Finding the last matching node

```go
lastLink := root.FindLastNode(xhtml.FilterByElement(xhtml.A))
```

## Available Filters

```go
// By element type
xhtml.FilterByElement(xhtml.DIV)
xhtml.FilterByElement(xhtml.A)
xhtml.FilterAnyElement()
xhtml.FilterAnyText()
xhtml.FilterAnyLeaf()

// By attribute
xhtml.FilterByAttribute("data-type", "active")
xhtml.FilterByID("main-content")
xhtml.FilterByClassName("card")
xhtml.FilterByRef("mainTitle")

// By regular expression (text nodes)
regex := regexp.MustCompile(`^\d+\.?\d*\s*€$`)
xhtml.FilterTextByRegEx(regex)
```

## Composable Filters

Filters can be combined with `And`, `Or`, and `Not`:

```go
// All <a> elements with class "external"
filter := xhtml.FilterByElement(xhtml.A).And(xhtml.FilterByClassName("external"))
links := root.FindAllNodes(filter)

// All elements that are NOT <br>
filter := xhtml.FilterByElement(xhtml.BR).Not()
nodes := root.FindAllNodes(filter)

// Either <h1> or <h2>
filter := xhtml.FilterByElement(xhtml.H1).Or(xhtml.FilterByElement(xhtml.H2))
headings := root.FindAllNodes(filter)
```

## Navigation

```go
node := root.FindFirstNode(xhtml.FilterByID("images"))

// First child element
child := node.NextChildElement()

// Next sibling element
next := child.NextElement()

// Previous sibling element
prev := next.PrevElement()

// All direct child elements
children := node.GetAllChildElements()
```

## Extracting Text and Attributes

```go
// Concatenate all descendant text with a separator
text := body.ConcatAllText("\n")

// Get an attribute (case-insensitive)
alt, ok := imgNode.GetAttributeValue("alt")

// Get the link from an <a> or the src from an <img>
link, ok := node.GetNodeLink()

// List all attribute keys
keys := node.GetAttributes() // ["src", "alt", "class"]
```

## Predefined Tag Constants

The library exports constants for common HTML tags:

```go
xhtml.HTML, xhtml.HEAD, xhtml.BODY, xhtml.DIV, xhtml.SPAN,
xhtml.H1, xhtml.H2, xhtml.H3, xhtml.P, xhtml.A, xhtml.IMG,
xhtml.UL, xhtml.LI, xhtml.TABLE, xhtml.TR, xhtml.TD,
xhtml.FORM, xhtml.INPUT, xhtml.BUTTON, // ... and more
```

## Complete Example

```go
package main

import (
  "bytes"
  "fmt"
  "log"
  "os"

  "github.com/alex-cos/xhtml"
  "golang.org/x/net/html"
)

func main() {
  data, err := os.ReadFile("page.html")
  if err != nil {
    log.Fatal(err)
  }

  doc, err := html.Parse(bytes.NewReader(data))
  if err != nil {
    log.Fatal(err)
  }

  root := xhtml.NewNode(doc)
  body := root.FindFirstNode(xhtml.FilterByElement(xhtml.BODY))

  // Extract all links
  links := body.FindAllNodes(xhtml.FilterByElement(xhtml.A))
  for _, link := range links {
    href, _ := link.GetAttributeValue("href")
    text := link.ConcatAllText(" ")
    fmt.Printf("[%s] %s\n", href, text)
  }

  // Extract all images
  images := body.FindAllNodes(xhtml.FilterByElement(xhtml.IMG))
  for _, img := range images {
    src, _ := img.GetAttributeValue("src")
    alt, _ := img.GetAttributeValue("alt")
    fmt.Printf("Image: %s (alt: %s)\n", src, alt)
  }
}
```

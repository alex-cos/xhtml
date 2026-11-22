# xhtml

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue)](https://go.dev/)
[![Test Status](https://github.com/alex-cos/xhtml/actions/workflows/test.yml/badge.svg)](https://github.com/alex-cos/xhtml/actions/workflows/test.yml)
[![Codecov](https://codecov.io/gh/alex-cos/xhtml/branch/main/graph/badge.svg)](https://codecov.io/gh/alex-cos/xhtml)
[![Lint Status](https://github.com/alex-cos/xhtml/actions/workflows/lint.yml/badge.svg)](https://github.com/alex-cos/xhtml/actions/workflows/lint.yml)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/alex-cos/xhtml)](https://goreportcard.com/report/github.com/alex-cos/xhtml)

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
body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))
title := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagH1))

fmt.Println(title.GetData()) // "h1"
```

### Finding all matching nodes

```go
links := root.FindAllNodes(xhtml.FilterByElement(xhtml.TagA))

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
thirdLi := root.FindNthNode(3, xhtml.FilterByElement(xhtml.TagLi))

if !thirdLi.IsNil() {
  fmt.Println("Found!")
}
```

### Finding the last matching node

```go
lastLink := root.FindLastNode(xhtml.FilterByElement(xhtml.TagA))
```

## Available Filters

```go
// By element type
xhtml.FilterByElement(xhtml.TagDiv)
xhtml.FilterByElement(xhtml.TagA)
xhtml.FilterAnyElement()
xhtml.FilterAnyText()
xhtml.FilterAnyLeaf()

// By attribute
xhtml.FilterByAttribute("data-type", "active")
xhtml.FilterByID("main-content")
xhtml.FilterByClassName("card")
xhtml.FilterByClassNames("card", "active") // all classes required
xhtml.FilterByRef("mainTitle")

// By attribute presence or partial value (values are case-sensitive)
xhtml.FilterByAttributeExists("data-id")
xhtml.FilterByAttributePrefix("href", "https://")
xhtml.FilterByAttributeSuffix("src", ".png")
xhtml.FilterByAttributeContains("href", "example.com")

// By text content (text nodes)
xhtml.FilterByText("Contact")         // exact match, surrounding white space ignored
xhtml.FilterByTextContains("Contact") // substring match

// By regular expression (text nodes)
regex := regexp.MustCompile(`^\d+\.?\d*\s*€$`)
xhtml.FilterTextByRegEx(regex)
```

## Composable Filters

Filters can be combined with `And`, `Or`, and `Not`:

```go
// All <a> elements with class "external"
filter := xhtml.FilterByElement(xhtml.TagA).And(xhtml.FilterByClassName("external"))
links := root.FindAllNodes(filter)

// All elements that are NOT <br>
filter := xhtml.FilterByElement(xhtml.TagBr).Not()
nodes := root.FindAllNodes(filter)

// Either <h1> or <h2>
filter := xhtml.FilterByElement(xhtml.TagH1).Or(xhtml.FilterByElement(xhtml.TagH2))
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
xhtml.TagHTML, xhtml.TagHead, xhtml.TagBody, xhtml.TagDiv, xhtml.TagSpan,
xhtml.TagH1, xhtml.TagH2, xhtml.TagH3, xhtml.TagP, xhtml.TagA, xhtml.TagImg,
xhtml.TagUl, xhtml.TagLi, xhtml.TagTable, xhtml.TagTr, xhtml.TagTd,
xhtml.TagForm, xhtml.TagInput, xhtml.TagButton, // ... and more
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
  body := root.FindFirstNode(xhtml.FilterByElement(xhtml.TagBody))

  // Extract all links
  links := body.FindAllNodes(xhtml.FilterByElement(xhtml.TagA))
  for _, link := range links {
    href, _ := link.GetAttributeValue("href")
    text := link.ConcatAllText(" ")
    fmt.Printf("[%s] %s\n", href, text)
  }

  // Extract all images
  images := body.FindAllNodes(xhtml.FilterByElement(xhtml.TagImg))
  for _, img := range images {
    src, _ := img.GetAttributeValue("src")
    alt, _ := img.GetAttributeValue("alt")
    fmt.Printf("Image: %s (alt: %s)\n", src, alt)
  }
}
```

package xhtml

// FindFirstNode returns the first node, in document order, that matches
// filter. The search starts with node itself and goes through all its
// descendants. It returns a nil Node if nothing matches.
func (node Node) FindFirstNode(filter FilterFunc) Node {
	if node.IsNil() {
		return NilNode()
	}
	if filter(node) {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		found := NewNode(child).FindFirstNode(filter)
		if !found.IsNil() {
			return found
		}
	}

	return NilNode()
}

// FindLastNode returns the last node, in document order, that matches filter.
// The search covers node itself and all its descendants. It returns a nil
// Node if nothing matches.
func (node Node) FindLastNode(filter FilterFunc) Node {
	var found Node

	if node.IsNil() {
		return NilNode()
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		f := NewNode(child).FindLastNode(filter)
		if !f.IsNil() {
			found = f
		}
	}
	if found.IsNil() && filter(node) {
		found = node
	}

	return found
}

// FindAllNodes returns all the nodes, in document order, that match filter.
// The search covers node itself and all its descendants. The result is empty
// if nothing matches.
func (node Node) FindAllNodes(filter FilterFunc) []Node {
	if node.IsNil() {
		return []Node{}
	}
	founded := []Node{}

	if filter(node) {
		founded = append(founded, node)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		nodes := NewNode(child).FindAllNodes(filter)
		founded = append(founded, nodes...)
	}

	return founded
}

// FindNthNode returns the pos-th node (starting at 1), in document order, that
// matches filter. It returns a nil Node if pos is lower than 1 or if there are
// fewer than pos matches.
func (node Node) FindNthNode(pos int, filter FilterFunc) Node {
	if pos <= 0 {
		return NilNode()
	}

	nodes := node.FindAllNodes(filter)
	if pos-1 < len(nodes) {
		return nodes[pos-1]
	}

	return NilNode()
}

package xhtml

func (node Node) FindFirstNode(filter FilterFunc) Node {
	if node.IsNil() {
		return NilNode()
	}
	if filter(&node) {
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
	if found.IsNil() && filter(&node) {
		found = node
	}

	return found
}

func (node Node) FindAllNodes(filter FilterFunc) []Node {
	if node.IsNil() {
		return []Node{}
	}
	founded := []Node{}

	if filter(&node) {
		founded = append(founded, node)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		nodes := NewNode(child).FindAllNodes(filter)
		founded = append(founded, nodes...)
	}

	return founded
}

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

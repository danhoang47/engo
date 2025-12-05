package engine

import "engo/pkg/scene"

// Get parent that node should placed in
// should not return nil (and never be)
func (e *Engine) GetInsertTarget() *scene.Node {
	if len(e.selection) == 0 {
		return e.rootNode
	}

	target := e.selection[0]

	// Group or Frame (page view)
	if target.IsContainer() {
		return target
	}

	if target.Parent != nil {
		return target.Parent
	}

	return e.rootNode
}

func (e *Engine) SetSelection(n *scene.Node, replace bool) {
	if replace {
		e.selection = append(e.selection, n)
		return
	}

	e.selection = []*scene.Node{n}
}

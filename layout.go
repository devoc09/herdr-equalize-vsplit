package main

import "sort"

type LayoutNode struct {
	Type      string      `json:"type"`
	Direction string      `json:"direction,omitempty"`
	Ratio     float64     `json:"ratio,omitempty"`
	PaneID    string      `json:"pane_id,omitempty"`
	First     *LayoutNode `json:"first,omitempty"`
	Second    *LayoutNode `json:"second,omitempty"`
}

type SplitTarget struct {
	Path      []bool
	Ratio     float64
	Direction string
}

func axisSpan(node *LayoutNode, axis string) int {
	if node.Type == "pane" {
		return 1
	}
	if node.Direction == axis {
		return axisSpan(node.First, axis) + axisSpan(node.Second, axis)
	}
	return 1
}

func collectTargets(node *LayoutNode, path []bool, out *[]SplitTarget) {
	if node.Type != "split" {
		return
	}
	first := axisSpan(node.First, node.Direction)
	second := axisSpan(node.Second, node.Direction)
	ratio := float64(first) / float64(first+second)

	p := make([]bool, len(path))
	copy(p, path)
	*out = append(*out, SplitTarget{Path: p, Ratio: ratio, Direction: node.Direction})

	pFirst := make([]bool, len(path)+1)
	copy(pFirst, path)
	pFirst[len(path)] = false
	collectTargets(node.First, pFirst, out)

	pSecond := make([]bool, len(path)+1)
	copy(pSecond, path)
	pSecond[len(path)] = true
	collectTargets(node.Second, pSecond, out)
}

func EqualizeColumns(root *LayoutNode) []SplitTarget {
	var targets []SplitTarget
	collectTargets(root, nil, &targets)

	sort.Slice(targets, func(i, j int) bool {
		return len(targets[i].Path) < len(targets[j].Path)
	})
	return targets
}

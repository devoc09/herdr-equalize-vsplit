package main

import (
	"math"
	"testing"
)

func pane(id string) *LayoutNode {
	return &LayoutNode{Type: "pane", PaneID: id}
}

func rightSplit(ratio float64, first, second *LayoutNode) *LayoutNode {
	return &LayoutNode{Type: "split", Direction: "right", Ratio: ratio, First: first, Second: second}
}

func downSplit(ratio float64, first, second *LayoutNode) *LayoutNode {
	return &LayoutNode{Type: "split", Direction: "down", Ratio: ratio, First: first, Second: second}
}

func TestAxisSpan(t *testing.T) {
	tests := []struct {
		name string
		node *LayoutNode
		want int
	}{
		{"pane", pane("x"), 1},
		{"single right split", rightSplit(0.5, pane("a"), pane("b")), 2},
		{"three columns: p1 | (p2 | p3)", rightSplit(0.5, pane("a"), rightSplit(0.5, pane("b"), pane("c"))), 3},
		{"cross-axis down", downSplit(0.5, pane("a"), pane("b")), 1},
		{"right split containing down: p1 | (p2 / p3)", rightSplit(0.5, pane("a"), downSplit(0.5, pane("b"), pane("c"))), 2},
		{"four columns", rightSplit(0.5, rightSplit(0.5, pane("a"), pane("b")), rightSplit(0.5, pane("c"), pane("d"))), 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := axisSpan(tt.node, "right")
			if got != tt.want {
				t.Errorf("axisSpan() = %d, want %d", got, tt.want)
			}
		})
	}
}

func within(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

func TestEqualizeColumns(t *testing.T) {
	tests := []struct {
		name  string
		root  *LayoutNode
		count int
		check func(*testing.T, []SplitTarget)
	}{
		{
			name:  "single pane, no splits",
			root:  pane("a"),
			count: 0,
		},
		{
			name:  "two columns 50/50",
			root:  rightSplit(0.3, pane("a"), pane("b")),
			count: 1,
			check: func(t *testing.T, targets []SplitTarget) {
				if !within(targets[0].Ratio, 0.5, 0.001) {
					t.Errorf("root ratio = %f, want 0.5", targets[0].Ratio)
				}
			},
		},
		{
			name:  "three columns 33/33/33",
			root:  rightSplit(0.5, pane("a"), rightSplit(0.5, pane("b"), pane("c"))),
			count: 2,
			check: func(t *testing.T, targets []SplitTarget) {
				root := targets[0]
				inner := targets[1]
				if !within(root.Ratio, 1.0/3.0, 0.001) {
					t.Errorf("root ratio = %f, want 0.333", root.Ratio)
				}
				if !within(inner.Ratio, 0.5, 0.001) {
					t.Errorf("inner ratio = %f, want 0.5", inner.Ratio)
				}
			},
		},
		{
			name:  "four columns 25/25/25/25",
			root:  rightSplit(0.5, rightSplit(0.5, pane("a"), pane("b")), rightSplit(0.5, pane("c"), pane("d"))),
			count: 3,
			check: func(t *testing.T, targets []SplitTarget) {
				for _, tr := range targets {
					if !within(tr.Ratio, 0.5, 0.001) {
						t.Errorf("split ratio = %f, want 0.5", tr.Ratio)
					}
				}
			},
		},
		{
			name:  "left-weighted three columns: (p1 | p2) | p3",
			root:  rightSplit(0.5, rightSplit(0.5, pane("a"), pane("b")), pane("c")),
			count: 2,
			check: func(t *testing.T, targets []SplitTarget) {
				root := targets[0]
				inner := targets[1]
				if !within(root.Ratio, 2.0/3.0, 0.001) {
					t.Errorf("root ratio = %f, want 0.667", root.Ratio)
				}
				if !within(inner.Ratio, 0.5, 0.001) {
					t.Errorf("inner ratio = %f, want 0.5", inner.Ratio)
				}
			},
		},
		{
			name:  "down split target computed, caller filters",
			root:  rightSplit(0.5, pane("a"), downSplit(0.7, pane("b"), pane("c"))),
			count: 2,
			check: func(t *testing.T, targets []SplitTarget) {
				for _, tr := range targets {
					if tr.Direction == "right" {
						if !within(tr.Ratio, 0.5, 0.001) {
							t.Errorf("right split ratio = %f, want 0.5", tr.Ratio)
						}
					}
				}
			},
		},
		{
			name:  "five columns 20% each",
			root:  rightSplit(0.5, pane("a"), rightSplit(0.5, pane("b"), rightSplit(0.5, pane("c"), rightSplit(0.5, pane("d"), pane("e"))))),
			count: 4,
			check: func(t *testing.T, targets []SplitTarget) {
				root := targets[0]
				if !within(root.Ratio, 1.0/5.0, 0.001) {
					t.Errorf("root ratio = %f, want 0.2", root.Ratio)
				}
			},
		},
		{
			name:  "equal quad (two rows, two columns) unchanged",
			root:  rightSplit(0.5, downSplit(0.5, pane("a"), pane("c")), downSplit(0.5, pane("b"), pane("d"))),
			count: 3,
			check: func(t *testing.T, targets []SplitTarget) {
				for _, tr := range targets {
					if !within(tr.Ratio, 0.5, 0.001) {
						t.Errorf("split ratio = %f, want 0.5", tr.Ratio)
					}
				}
			},
		},
		{
			name:  "sort order: shorter paths first",
			root:  rightSplit(0.5, pane("a"), rightSplit(0.5, pane("b"), pane("c"))),
			count: 2,
			check: func(t *testing.T, targets []SplitTarget) {
				if len(targets[0].Path) != 0 {
					t.Errorf("first target should be root (len 0), got len %d", len(targets[0].Path))
				}
				if len(targets[1].Path) < 1 {
					t.Errorf("second target should be inner (len >= 1), got len %d", len(targets[1].Path))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EqualizeColumns(tt.root)
			if len(got) != tt.count {
				t.Errorf("EqualizeColumns() returned %d targets, want %d", len(got), tt.count)
			}
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

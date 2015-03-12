package louvain

import (
	"testing"
)

func TestIsolated(t *testing.T) {
	g := MakeNewGraph(100, MergeFn(func(a []*Node) interface{} { return true }))
	for i := 0; i < 100; i++ {
		g.AddLink(i, i, 10)
	}
	g2 := g.NextLevel(-1, 0)
	if g2.Total != g.Total {
		t.Fatalf("Total not matched: %v != %v", g.Total, g2.Total)
	}
	if len(g2.Nodes) != len(g.Nodes) {
		t.Fatalf("Node size not matched: %v != %v", len(g.Nodes), len(g2.Nodes))
	}
	for i := range g2.Nodes {
		node := &g2.Nodes[i]
		if len(node.Children) != 1 {
			t.Fatalf("Too many Childrenren: %v", len(node.Children))
		}
		if node.Children[0].Parent != node {
			t.Fatalf("Children and parent not matched: %v != %v", node.Children[0].Parent, node)
		}
	}
}
func TestConnected(t *testing.T) {
	g := MakeNewGraph(100, MergeFn(func(a []*Node) interface{} { return true }))
	for i := range g.Nodes {
		mod := i % 10
		g.AddLink(i, i-mod, 10)
	}
	g.AddLink(10, 20, 1)
	g2 := g.NextLevel(-1, 0)
	if g2.Total != g.Total {
		t.Fatalf("Total not matched: %v != %v", g.Total, g2.Total)
	}
	if len(g2.Nodes) != 10 {
		t.Fatalf("Node size not matched: %v != 10", len(g2.Nodes))
	}
	for i := range g2.Nodes {
		node := &g2.Nodes[i]
		if len(node.Children) != 10 {
			t.Fatalf("Childrenren size: %v", len(node.Children))
		}
		for ci := range node.Children {
			if node.Children[ci].Parent != node {
				t.Fatalf("Children and parent not matched: %v != %v", node.Children[ci].Parent, node)
			}
		}
	}
}

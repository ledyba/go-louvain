package louvain

import (
	"testing"
)

func TestIsolated(t *testing.T) {
	g := MakeNewGraph(100, MergeFn(func(a []*Node) interface{}{return true}))
	for i:=0;i<100;i++{
		g.Connect(i,i,1)
	}
	g2 := g.NextLevel(-1)
	if g2.Total != g.Total {
		t.Fatalf("Total not matched: %v != %v", g.Total, g2.Total)
	}
	if len(g2.Nodes) != len(g.Nodes) {
		t.Fatalf("Node size not matched: %v != %v", len(g.Nodes), len(g2.Nodes))
	}
	for i := range g2.Nodes{
		node := &g2.Nodes[i]
		if len(node.Child) != 1 {
			t.Fatalf("Too many children: %v", len(node.Child))
		}
		if node.Child[0].Parent != node {
			t.Fatalf("Child and parent not matched: %v != %v", node.Child[0].Parent, node)
		}
	}
}
func TestConnected(t *testing.T) {
	g := MakeNewGraph(100, MergeFn(func(a []*Node) interface{}{return true}))
	for i := range g.Nodes {
		mod := i%10
		g.Connect(i,i-mod, 10)
	}
	g2 := g.NextLevel(-1)
	if g2.Total != g.Total {
		t.Fatalf("Total not matched: %v != %v", g.Total, g2.Total)
	}
	if len(g2.Nodes) != 10 {
		t.Fatalf("Node size not matched: %v != 10", len(g.Nodes))
	}
	for i := range g2.Nodes{
		node := &g2.Nodes[i]
		if len(node.Child) != 10 {
			t.Fatalf("Children size: %v", len(node.Child))
		}
		for ci := range node.Child {
			if node.Child[ci].Parent != node {
				t.Fatalf("Child and parent not matched: %v != %v", node.Child[ci].Parent, node)
			}
		}
	}
}

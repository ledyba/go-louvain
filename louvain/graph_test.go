package louvain

import (
	"testing"
)

func TestEmpty(t *testing.T) {
	nodes := make([]Node, 10)
	tot := 0
	for i := range nodes {
		node := &nodes[i]
		node.Degree = 0
		node.SelfLoop = 1
		tot++
	}
	g := MakeNewGraph(tot, nodes)
	g2 := g.NextLevel()
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

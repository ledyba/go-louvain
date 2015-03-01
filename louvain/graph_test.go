package louvain

import (
	"testing"
)

func TestIsolated(t *testing.T) {
	nodes := make([]Node, 10)
	tot := 0
	for i := range nodes {
		node := &nodes[i]
		node.Degree = 1
		node.SelfLoop = 1
		tot++
	}
	g := MakeNewGraph(tot, nodes, func(a,b interface{})bool{return true})
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
	nodes := make([]Node, 100)
	tot := 0
	for i := range nodes {
		node := &nodes[i]
		node.SelfLoop = 0
		node.Links = []Link{}
		mod := i%10
		if mod > 0{
			p := &nodes[i-mod]
			p.Degree+=10
			tot+=10
			node.Degree+=10
			tot+=10
			node.Links = append(node.Links, Link{10,p})
			p.Links = append(p.Links, Link{10,node})
		}
	}
	g := MakeNewGraph(tot, nodes, func(a,b interface{})bool{return true})
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

package louvain

import (
	"log"
	"math/rand"
)

type Link struct {
	To     int
	Weight int
}
type MergeFn func(cs []*Node) interface{}
type Node struct {
	Data     interface{}
	Parent   *Node
	Children []*Node
	Links    []Link
	Degree   int
	SelfLoop int
}

type Graph struct {
	Total   int
	Nodes   []Node
	mergeFn MergeFn
}

func MakeNewGraph(size int, mergeFn MergeFn) *Graph {
	nodes := make([]Node, size)
	return &Graph{0, nodes, mergeFn}
}
func MakeNewGraphFromNodes(nodes []Node, totalLinks int, mergeFn MergeFn) *Graph {
	return &Graph{totalLinks, nodes, mergeFn}
}

func (node *Node) Print() {
	log.Printf("  Children: %v nodes", len(node.Children))
	log.Printf("  Self loop: %d/%d", node.SelfLoop, node.Degree)
	for to, weight := range node.Links {
		log.Printf("  -> %d: %d/%d", to, weight, node.Degree)
	}
	log.Printf("  Payload: %v", node.Data)
}

func (g *Graph) Print() {
	log.Printf("TotalLinks: %v", g.Total)
	for i := range g.Nodes {
		node := &g.Nodes[i]
		log.Printf("<<Node %d>>", i)
		node.Print()
	}
}

func (g *Graph) NextLevel(limit int, precision float32) *Graph {
	gTotal := float32(g.Total)
	commTotal := make([]int, len(g.Nodes))
	commIn := make([]int, len(g.Nodes))
	tmpComm := make([]int, len(g.Nodes))
	for i := range g.Nodes {
		node := &g.Nodes[i]
		tmpComm[i] = i
		commTotal[i] = node.Degree
		commIn[i] = node.SelfLoop
	}
	order := rand.Perm(len(g.Nodes))
	neighLinks := make([]int, len(g.Nodes))
	neighComm := make([]int, 0, len(g.Nodes))
	changed := len(g.Nodes)
	cnt := 0
	for changed > len(g.Nodes)/100 {
		if limit > 0 && cnt >= limit {
			log.Printf("Exceed limit")
			break
		}
		changed = 0
		cnt++
		for _, rpos := range order {
			node := &g.Nodes[rpos]
			nodeTmpComm := tmpComm[rpos]
			nodeDegree := node.Degree
			nodeSelfLoop := node.SelfLoop
			/* Calculating Neighbor Communities */
			for _, comm := range neighComm {
				neighLinks[comm] = 0
			}
			neighComm = neighComm[0:0]
			for _, link := range node.Links {
				to := tmpComm[link.To]
				if neighLinks[to] <= 0 {
					neighComm = append(neighComm, to)
					neighLinks[to] = link.Weight
				} else {
					neighLinks[to] += link.Weight
				}
			}
			/* Calculating the BEST community */
			bestComm := nodeTmpComm
			bestGain := precision
			for _, comm := range neighComm {
				var gain float32
				if comm == nodeTmpComm {
					gain = float32(neighLinks[comm]) - float32((commTotal[comm]-nodeDegree))*float32(nodeDegree)/gTotal
				} else {
					gain = float32(neighLinks[comm]) - float32(commTotal[comm])*float32(nodeDegree)/gTotal
				}
				if gain > bestGain {
					bestGain = gain
					bestComm = comm
				}
			}
			/* Insert to the best community */
			if nodeTmpComm != bestComm {
				changed++
				tmpComm[rpos] = bestComm
				/* Remove from the original community */
				commTotal[nodeTmpComm] -= nodeDegree
				commIn[nodeTmpComm] -= 2*neighLinks[nodeTmpComm] + nodeSelfLoop
				/* insert */
				commTotal[bestComm] += nodeDegree
				commIn[bestComm] += 2*neighLinks[bestComm] + nodeSelfLoop
			}
		}
	}

	//Calc Next nodes:
	communities := make([]Node, 0, len(g.Nodes)/2)
	oldCommIdx := make([]int, 0, len(g.Nodes)/2)
	c2i := make([]int, len(g.Nodes))
	links := make([]map[int]int, 0)
	for i := range g.Nodes {
		node := &g.Nodes[i]
		nodeTmpComm := tmpComm[i]
		c := c2i[nodeTmpComm]
		var comm *Node
		if c <= 0 {
			c2i[nodeTmpComm] = len(communities) + 1
			communities = append(communities, Node{
				Children: make([]*Node, 0),
			})
			oldCommIdx = append(oldCommIdx, nodeTmpComm)
			links = append(links, nil)
			comm = &communities[len(communities)-1]
		} else {
			comm = &communities[c-1]
		}
		comm.Children = append(comm.Children, node)
	}
	// Merging edges
	for i := range communities {
		comm := &communities[i]
		oldComm := oldCommIdx[i]
		comm.Data = g.mergeFn(comm.Children)
		link := links[i]
		if link == nil {
			link = make(map[int]int)
			links[i] = link
		}
		for _, child := range comm.Children {
			child.Parent = comm
			comm.SelfLoop += child.SelfLoop
			comm.Degree += child.SelfLoop
			for _, edge := range child.Links {
				cLinkToCommNow := tmpComm[edge.To]
				comm.Degree += edge.Weight
				if cLinkToCommNow == oldComm {
					comm.SelfLoop += edge.Weight
				} else {
					link[c2i[cLinkToCommNow]-1] += edge.Weight
				}
			}
		}
	}
	for i := range links {
		comm := &communities[i]
		comm.Links = make([]Link, 0, len(links[i]))
		link := links[i]
		if link == nil {
			continue
		}
		for to, weight := range link {
			comm.Links = append(comm.Links, Link{to, weight})
		}
	}
	return &Graph{g.Total, communities, g.mergeFn}
}

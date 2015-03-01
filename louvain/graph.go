package louvain

import (
	"math/rand"
)

type Link struct {
	Weight int
	To     *Node
}
type Node struct {
	Data     interface{}
	Parent   *Node
	Child    []*Node
	Links    []Link
	Degree   int
	SelfLoop int
	tmpComm  int
}

type Graph struct {
	Total int
	Nodes []Node
}

func MakeNewGraph(total int, nodes []Node) *Graph {
	return &Graph{total, nodes}
}

func shuffleOrder(size int) []int {
	array := make([]int, size)
	for i := range array {
		array[i] = i
	}
	for i := 0; i < size-1; i++ {
		rpos := rand.Intn(size-i) + 1
		array[i], array[rpos] = array[rpos], array[i]
	}
	return array
}

func (g *Graph) NextLevel() *Graph {
	commTotal := make([]int, len(g.Nodes))
	commIn := make([]int, len(g.Nodes))
	for i := range g.Nodes {
		node := &g.Nodes[i]
		node.tmpComm = i
		commTotal[i] = node.Degree
		commIn[i] = node.SelfLoop
	}
	order := shuffleOrder(len(g.Nodes))
	neighLinks := make([]int, len(g.Nodes))
	neighComm := make([]int, 0, len(g.Nodes))
	changed := 1
	for changed > len(g.Nodes)/100 {
		changed = 0
		for _, rpos := range order {
			node := g.Nodes[rpos]
			/* Calculating Neighbor Communities */
			for _, comm := range neighComm {
				neighLinks[comm] = -1
			}
			neighComm = neighComm[0:0]
			for _, link := range node.Links {
				to := link.To.tmpComm
				if neighLinks[to] < 0 {
					neighComm = append(neighComm, to)
					neighLinks[to] = link.Weight
				} else {
					neighLinks[to] += link.Weight
				}
			}
			/* Remove from the original community */
			commTotal[node.tmpComm] -= node.Degree
			commIn[node.tmpComm] -= 2*neighLinks[node.tmpComm] + node.SelfLoop
			/* Calculating the BEST community */
			best_comm := node.tmpComm
			best_gain := float32(0.0)
			for _, comm := range neighComm {
				gain := float32(neighLinks[node.tmpComm]) - float32(commTotal[node.tmpComm]*node.Degree)/float32(g.Total)
				if gain > best_gain {
					best_comm = comm
					best_gain = gain
				}
			}
			/* Insert to the best community */
			if node.tmpComm != best_comm{
				changed++
			}
			node.tmpComm = best_comm
			commTotal[node.tmpComm] += node.Degree
			commIn[node.tmpComm] += 2*neighLinks[node.tmpComm] + node.SelfLoop
		}
	}
	
	//Calc Next nodes:
	communities := make([]Node, 0, len(g.Nodes)/2)
	c2i := make([]int, len(g.Nodes))
	for i := range g.Nodes {
		node := &g.Nodes[i]
		c := c2i[node.tmpComm]
		var comm *Node
		if c <= 0 {
			c2i[node.tmpComm] = len(communities)+1
			communities = append(communities, Node{})
			comm = &communities[len(communities)-1]
			comm.Child = make([]*Node, 0)
			comm.Links = make([]Link, 0)
			comm.tmpComm = node.tmpComm
		}else{
			comm = &communities[c-1]
		}
		comm.Child = append(comm.Child, node)
	}
	for i := range communities {
		comm := &communities[i]
		for _,child := range comm.Child {
			child.Parent = comm
			comm.SelfLoop += child.SelfLoop
			for _,link := range child.Links {
				if link.To.tmpComm == comm.tmpComm {
					comm.SelfLoop += link.Weight
				} else {
					found := false
					lcomm := link.To.tmpComm
					for j := range comm.Links {
						if comm.Links[j].To.tmpComm == lcomm {
							comm.Links[j].Weight += link.Weight
							found = true
						}
					}
					if !found {
						linked := &communities[c2i[lcomm]-1]
						comLink := Link {link.Weight, linked}
						comm.Links = append(comm.Links, comLink)
					}
				}
			}
		}
	}

	return &Graph{g.Total, communities}
}

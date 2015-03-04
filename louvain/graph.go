package louvain

import (
	"math/rand"
)

type Link struct {
	Weight int
	To     *Node
}
type MergeFn func(cs []*Node) interface{}
type Node struct {
	Data     interface{}
	Parent   *Node
	Child    []*Node
	Links    map[int]int
	Degree   int
	SelfLoop int
}

type Graph struct {
	Total int
	Nodes []Node
	mergeFn  MergeFn
}

func MakeNewGraph(size int, mergeFn MergeFn) *Graph {
	nodes := make([]Node, size)
	return &Graph{0, nodes, mergeFn}
}

func (g *Graph) Connect(i,j,w int) {
	if i == j {
		g.Nodes[i].SelfLoop += w*2
	}else{
		if g.Nodes[i].Links == nil{
			g.Nodes[i].Links = make(map[int]int)
		}
		if g.Nodes[j].Links == nil{
			g.Nodes[j].Links = make(map[int]int)
		}
		g.Nodes[i].Links[j]+=w
		g.Nodes[j].Links[i]+=w
	}
	g.Total += w*2
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

func (g *Graph) NextLevel(limit int) *Graph {
	commTotal := make([]int, len(g.Nodes))
	commIn := make([]int, len(g.Nodes))
	tmpComm := make([]int, len(g.Nodes))
	for i := range g.Nodes {
		node := &g.Nodes[i]
		tmpComm[i] = i
		commTotal[i] = node.Degree
		commIn[i] = node.SelfLoop
	}
	order := shuffleOrder(len(g.Nodes))
	neighLinks := make([]int, len(g.Nodes))
	neighComm := make([]int, 0, len(g.Nodes))
	changed := len(g.Nodes)
	cnt := 0
	for changed > len(g.Nodes)/100 && (limit <= 0 || cnt < limit){
		changed = 0
		cnt++
		for _, rpos := range order {
			node := &g.Nodes[rpos]
			nodeTmpComm := tmpComm[rpos]
			/* Calculating Neighbor Communities */
			for _, comm := range neighComm {
				neighLinks[comm] = 0
			}
			neighComm = neighComm[0:0]
			for linkToIdx,weight := range node.Links {
				to := tmpComm[linkToIdx]
				if neighLinks[to] <= 0 {
					neighComm = append(neighComm, to)
					neighLinks[to] = weight
				} else {
					neighLinks[to] += weight
				}
			}
			/* Remove from the original community */
			commTotal[nodeTmpComm] -= node.Degree
			commIn[nodeTmpComm] -= 2*neighLinks[nodeTmpComm] + node.SelfLoop
			/* Calculating the BEST community */
			best_comm := nodeTmpComm
			best_gain := float32(0.0)
			for _, comm := range neighComm {
				gain := float32(neighLinks[comm]) - float32(commTotal[comm]*node.Degree)/float32(g.Total)
				if gain > best_gain {
					best_comm = comm
					best_gain = gain
				}
			}
			/* Insert to the best community */
			if nodeTmpComm != best_comm {
				changed++
			}
			tmpComm[rpos] = best_comm
			commTotal[best_comm] += node.Degree
			commIn[best_comm] += 2*neighLinks[best_comm] + node.SelfLoop
		}
	}

	//Calc Next nodes:
	communities := make([]Node, 0, len(g.Nodes)/2)
	oldCommIdx := make([]int, 0, len(g.Nodes)/2)
	c2i := make([]int, len(g.Nodes))
	for i := range g.Nodes {
		node := &g.Nodes[i]
		nodeTmpComm := tmpComm[i]
		c := c2i[nodeTmpComm]
		var comm *Node
		if c <= 0 {
			c2i[nodeTmpComm] = len(communities) + 1
			communities = append(communities, Node{})
			oldCommIdx = append(oldCommIdx, nodeTmpComm)
			comm = &communities[len(communities)-1]
			comm.Child = make([]*Node, 0)
			comm.Links = make(map[int]int)
		} else {
			comm = &communities[c-1]
		}
		comm.Child = append(comm.Child, node)
	}
	// Merging edges
	for i := range communities {
		comm := &communities[i]
		oldComm := oldCommIdx[i]
		comm.Data = g.mergeFn(comm.Child)
		for _, child := range comm.Child {
			child.Parent = comm
			comm.SelfLoop += child.SelfLoop
			for linkToIdx,weight := range child.Links {
				cLinkToCommNow := tmpComm[linkToIdx]
				if cLinkToCommNow == oldComm {
					comm.SelfLoop += weight
				} else {
					comm.Links[c2i[cLinkToCommNow]-1]+=weight
				}
			}
		}
	}

	return &Graph{g.Total, communities, g.mergeFn}
}

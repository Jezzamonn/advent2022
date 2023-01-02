package main

import (
	"fmt"
)

const startNodeName = "AA"
const startTime = 26

type SearchState struct {
	CurrNode string
	// All the nodes we've stopped at and opened a valve at. This doesn't include nodes we passed through to get to other nodes.
	NodesVisitedSet map[string]struct{}

	// Whether we've restarted the search from the state node, to represent the elephant doing it's own search.
	Restarted bool
	Flow      int
	TimeLeft  int
}

func (s SearchState) Value() int {
	return s.Flow
}

func (s SearchState) UpperBound(g *Graph) int {
	value := UpperBoundFromNode(s.CurrNode, g, s.NodesVisitedSet, s.TimeLeft)
	if !s.Restarted {
		// The upper bound also includes the value of a restarted search from the start node.
		value += UpperBoundFromNode(startNodeName, g, s.NodesVisitedSet, startTime)
	}
	return s.Flow + value
}

// The flow from visiting all the nodes that are reachable within the time left from this node
func UpperBoundFromNode(nodeName string, g *Graph, alreadyVisited map[string]struct{}, timeLeft int) int {
	// For the moment, just add up the flow of each reachable node.

	// We could lower this bound by simulating visiting them all in order of
	// distance, pretending that we can get from the first node to the second
	// node in the difference in distance between the two. This seems good
	// enough for now.
	extraFlow := 0
	for _, node := range g.Nodes {
		if _, ok := alreadyVisited[node.Name]; ok {
			continue
		}
		dist := g.Distances[NamePair{nodeName, node.Name}]
		if dist > timeLeft {
			continue
		}
		extraFlow += node.Value * (timeLeft - dist)
	}
	return extraFlow
}

func (s SearchState) GetSubStates(g *Graph) []SearchState {
	subStates := make([]SearchState, 0)
	currNodeName := s.CurrNode

	// Create a new state for visiting every other node.
	for _, node := range g.Nodes {
		// No point revisiting a node.
		if _, ok := s.NodesVisitedSet[node.Name]; ok {
			continue
		}
		// Add 1 to account for the time to open the node.
		dist := g.Distances[NamePair{currNodeName, node.Name}] + 1
		if (dist > s.TimeLeft) || (dist == 0) {
			continue
		}

		// Duplicate the nodes visited
		newNodesVisitedSet := make(map[string]struct{})
		for k, v := range s.NodesVisitedSet {
			newNodesVisitedSet[k] = v
		}
		newNodesVisitedSet[node.Name] = struct{}{}

		newTimeLeft := s.TimeLeft - dist
		newFlow := s.Flow + node.Value*newTimeLeft

		newState := SearchState{
			CurrNode:        node.Name,
			NodesVisitedSet: newNodesVisitedSet,
			Restarted:       s.Restarted,
			Flow:            newFlow,
			TimeLeft:        newTimeLeft,
		}
		subStates = append(subStates, newState)
	}

	// And we might also be able to stop here and restart the search.
	if !s.Restarted {
		// Copy the nodes visited set because no need to revisit these nodes.
		newNodesVisitedSet := make(map[string]struct{})
		for k, v := range s.NodesVisitedSet {
			newNodesVisitedSet[k] = v
		}

		newState := SearchState{
			CurrNode:        startNodeName,
			NodesVisitedSet: newNodesVisitedSet,
			Restarted:       true,
			Flow:            s.Flow,
			TimeLeft:        startTime,
		}
		subStates = append(subStates, newState)
	}

	return subStates
}

func parseFile(filename string) *Graph {
	return GraphFromFile(filename)
}

func solve(filename string) {
	fmt.Println("Solving", filename)

	graph := parseFile(filename)

	// Do branch and bound
	s := SearchState{
		CurrNode:        startNodeName,
		NodesVisitedSet: map[string]struct{}{startNodeName: {}},
		Restarted:       false,
		Flow:            0,
		TimeLeft:        startTime,
	}
	best := s

	// Just visit all states with a DFS to start
	toVisit := make([]SearchState, 0)
	toVisit = append(toVisit, s)

	searched := 0
	skipped := 0
	for len(toVisit) > 0 {
		// Pop the first node
		currentState := toVisit[0]
		toVisit = toVisit[1:]

		if currentState.Flow > best.Flow {
			fmt.Println("Found new best", currentState.Flow, "\tsearched", searched, "skipped", skipped)
			best = currentState
		}

		// Bound: Don't explore substates if the upper bound is less than the best
		if currentState.UpperBound(graph) < best.Flow {
			skipped++
			continue
		}

		// Add the substates to the list of states to visit
		subStates := currentState.GetSubStates(graph)
		toVisit = append(toVisit, subStates...)
		searched++
	}
	// 2063 is too low (?)

	fmt.Println("Best flow", best.Flow, "\tsearched", searched, "skipped", skipped)
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	solve("16/input.txt")
}

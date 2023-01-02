package main

import (
	"fmt"
)

// Search state can be summarised as the nodes visited so far and the time left
type SearchState struct {
	// All the nodes we've stopped at and opened a valve at. This doesn't include nodes we passed through to get to other nodes.
	NodesVisited    []string
	NodesVisitedSet map[string]struct{}

	Flow     int
	TimeLeft int
}

func (s SearchState) CurrentNodeName() string {
	return s.NodesVisited[len(s.NodesVisited)-1]
}

func (s SearchState) Value() int {
	return s.Flow
}

func (s SearchState) UpperBound(g *Graph) int {
	// Upper bound is the visiting all the nodes that are reachable within the time left from this node
	curr := s.CurrentNodeName()

	// For the moment, just add up the flow.

	// We could lower this bound by simulating visiting them all in order of
	// distance, pretending that we can get from the first node to the second
	// node in the difference in distance between the two. This seems good
	// enough for now.
	extraFlow := 0
	for _, node := range g.Nodes {
		if _, ok := s.NodesVisitedSet[node.Name]; ok {
			continue
		}
		dist := g.Distances[NamePair{curr, node.Name}]
		if dist > s.TimeLeft {
			continue
		}
		extraFlow += node.Value * (s.TimeLeft - dist)
	}
	return s.Flow + extraFlow
}

func (s SearchState) GetSubStates(g *Graph) []SearchState {
	subStates := make([]SearchState, 0)
	currNodeName := s.CurrentNodeName()

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

		// Duplicate the nodes visited data structures
		newNodesVisited := make([]string, len(s.NodesVisited)+1)
		copy(newNodesVisited, s.NodesVisited)
		newNodesVisited[len(s.NodesVisited)] = node.Name

		newNodesVisitedSet := make(map[string]struct{})
		for k, v := range s.NodesVisitedSet {
			newNodesVisitedSet[k] = v
		}
		newNodesVisitedSet[node.Name] = struct{}{}

		newTimeLeft := s.TimeLeft - dist
		newFlow := s.Flow + node.Value*newTimeLeft

		newState := SearchState{
			NodesVisited:    newNodesVisited,
			NodesVisitedSet: newNodesVisitedSet,
			Flow:            newFlow,
			TimeLeft:        newTimeLeft,
		}
		subStates = append(subStates, newState)
	}
	return subStates
}

const startNodeName = "AA"

func parseFile(filename string) *Graph {
	return GraphFromFile(filename)
}

func solve(filename string) {
	fmt.Println("Solving", filename)

	graph := parseFile(filename)

	// Do branch and bound
	s := SearchState{
		NodesVisited:    []string{startNodeName},
		NodesVisitedSet: map[string]struct{}{startNodeName: {}},
		Flow:            0,
		TimeLeft:        30,
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

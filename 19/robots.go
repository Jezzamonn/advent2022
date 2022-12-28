package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var resources = [4]string{"ore", "clay", "obsidian", "geode"}

var resourceIndices = map[string]int{
	"ore":      0,
	"clay":     1,
	"obsidian": 2,
	"geode":    3,
}

// A blueprint is the cost of each robot type.
type Blueprint [4][4]int

type SearchState struct {
	TimeLeft           int
	Resources          [4]int
	ResourcesPerMinute [4]int
}

// Parses the cost of each robot type from a line of the input using regex.
//
// Example line:
// `Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 17 clay. Each geode robot costs 4 ore and 20 obsidian.`
func parseBlueprint(line string) Blueprint {
	// Regex to match one sentence describing one robot.
	// Example: `Each ore robot costs 4 ore.`
	robotRegex := regexp.MustCompile(`Each (\w+) robot costs (.+?)\.`)
	// Find all the sentences describing the robots.
	robotMatches := robotRegex.FindAllStringSubmatch(line, -1)
	// Create a new blueprint
	bp := Blueprint{}
	// For each robot, parse the cost
	for _, robotMatch := range robotMatches {
		robotType := robotMatch[1]
		robotTypeIndex := resourceIndices[robotType]
		// Parse the cost of the robot
		cost := parseCost(robotMatch[2])
		// Set the cost in the blueprint
		bp[robotTypeIndex] = cost
	}
	return bp
}

// Parses a cost string into an array of ints.
//
// Example cost string: `4 ore and 17 clay`
func parseCost(cost string) [4]int {
	// Regex to match one cost component.
	// Example: `4 ore`
	costRegex := regexp.MustCompile(`(\d+) (\w+)`)
	// Find all the cost components
	costMatches := costRegex.FindAllStringSubmatch(cost, -1)
	// Create a new cost array
	c := [4]int{}
	// For each cost component, set the cost in the array
	for _, costMatch := range costMatches {
		// Parse the cost
		value, err := strconv.Atoi(costMatch[1])
		if err != nil {
			panic(err)
		}
		// Set the cost in the array
		resource := costMatch[2]
		resourceIndex := resourceIndices[resource]
		c[resourceIndex] = value
	}
	return c
}

func getSubStates(state SearchState, bp Blueprint) []SearchState {
	subStates := make([]SearchState, 0)
	// For each type of robot, try waiting until we have enough resources and then building it.
outer:
	// r = robot
	for r := 0; r < 4; r++ {
		timeToWait := 0
		// For each resource, figure out how long we'd need to wait until we have enough.
		// If there's no way we can wait long enough, give up.

		// rr = robot resource
		for rr := 0; rr < 4; rr++ {
			amountLeft := bp[r][rr] - state.Resources[rr]
			if amountLeft <= 0 {
				continue
			}
			if amountLeft > 0 && state.ResourcesPerMinute[rr] == 0 {
				// We can't wait long enough to get enough of this resource.
				continue outer
			}
			// Ceiling division
			timeToWaitForResource := (amountLeft + state.ResourcesPerMinute[rr] - 1) / state.ResourcesPerMinute[rr]
			if timeToWaitForResource > timeToWait {
				timeToWait = timeToWaitForResource
			}
		}
		// We have to wait 1 extra minute to build the robot.
		timeToWait += 1

		// If we don't have enough time to wait, don't try
		if timeToWait > state.TimeLeft {
			continue
		}

		// Create a new state with the robot built
		newState := SearchState{
			TimeLeft:           state.TimeLeft - timeToWait,
			Resources:          state.Resources,
			ResourcesPerMinute: state.ResourcesPerMinute,
		}
		// For each resource, add the amount we'd get from waiting
		for rr := 0; rr < 4; rr++ {
			newState.Resources[rr] += timeToWait * state.ResourcesPerMinute[rr]
		}
		// Subtract the cost of the robot
		for rr := 0; rr < 4; rr++ {
			newState.Resources[rr] -= bp[r][rr]
		}
		// Add the resources we get from the robot
		newState.ResourcesPerMinute[r] += 1
		subStates = append(subStates, newState)
	}
	return subStates
}

// The score of this state if we did nothing else.
func scoreState(state SearchState) int {
	// We just care about how many geodes we'd have at the end.
	return state.Resources[3] + state.TimeLeft*state.ResourcesPerMinute[3]
}

// An upper bound of how good this state could possibly be.
func upperBound(state SearchState, bp Blueprint) int {
	// Perhaps not to far off the actual value, but lets assume that from here we can build a geode robot every minute.
	// The amount of geodes we'd have grows with a triangle pattern.
	// Given it takes 1 minute before the robot produces anything, we use state.TimeLeft - 1 in the equation.
	return scoreState(state) + state.TimeLeft*(state.TimeLeft-1)/2
}

func findMaxGeodes(bp Blueprint, time int) int {
	mostGeodes := 0
	// The initial state is we have no resources, but 1 ore robot.
	initialState := SearchState{
		TimeLeft:           time,
		Resources:          [4]int{},
		ResourcesPerMinute: [4]int{1, 0, 0, 0},
	}

	// For debugging
	searched := 0
	skipped := 0

	// Do a depth-first search of the search space.
	// Depth first is good because we can prune branches early.
	// We use a stack to do the search.
	stack := []SearchState{initialState}
	for len(stack) > 0 {
		// Pop the top state off the stack
		state := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// If this state is better than the best we've seen, update the best
		if score := scoreState(state); score > mostGeodes {
			fmt.Println("New best:", score, "searched:", searched, "skipped:", skipped)
			mostGeodes = score
		}

		// If this state could not be better than the best we've seen, don't bother exploring it.
		if upperBound(state, bp) <= mostGeodes {
			skipped++
			continue
		}

		// Get all the states we can get to from this state
		subStates := getSubStates(state, bp)
		// Push them onto the stack
		stack = append(stack, subStates...)
		searched++
	}

	fmt.Println("searched:", searched, "skipped:", skipped)

	return mostGeodes
}

func solvePt1(filename string) int {
	// Read the file, line by line
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	qualitySum := 0

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		// Parse the blueprint
		bp := parseBlueprint(scanner.Text())
		// Find the maximum number of geodes we can get
		mostGeodes := findMaxGeodes(bp, 24)
		quality := mostGeodes * (i + 1)
		qualitySum += quality
		fmt.Println("Geodes:", mostGeodes, "Quality:", quality, "Sum:", qualitySum)
	}

	return qualitySum
}

func solvePt2(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	mostGeodesMultiplied := 1

	scanner := bufio.NewScanner(file)
	// Just read the first 3 lines.
	for i := 0; i < 3 && scanner.Scan(); i++ {
		// Parse the blueprint
		bp := parseBlueprint(scanner.Text())
		// Find the maximum number of geodes we can get
		mostGeodes := findMaxGeodes(bp, 32)
		mostGeodesMultiplied *= mostGeodes
		fmt.Println("Geodes:", mostGeodes)
	}

	return mostGeodesMultiplied
}

func main() {
	filename := "19/input.txt"
	pt1 := solvePt1(filename)
	fmt.Println("Pt1:", pt1)

	pt2 := solvePt2(filename)
	fmt.Println("Pt2:", pt2)
}

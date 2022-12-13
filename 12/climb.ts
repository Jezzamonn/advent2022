const fs = require('fs');
const path = require('path');

interface Point {
    x: number;
    y: number;
}

interface QueueElem {
    position: Point;
    steps: number;
    parent?: QueueElem;
}

const moveDirs = [
    { x: 0, y: -1 },
    { x: 0, y: 1 },
    { x: -1, y: 0 },
    { x: 1, y: 0 },
]

function countStepsToTop(inputText: string, startingLetters: string[]): number {
    const map = inputText.split('\n').map(line => line.split(''));
    const visited = map.map(row => row.map(() => false));

    // Find the start
    const starts = findLetters(map, startingLetters);
    const end = findLetters(map, ['E'])[0];

    let printedOnce = false;

    const queue: QueueElem[] = [];
    // Start a search from each start position. The searches need to happen
    // independently as far as which nodes they've visited, but we can make it
    // more efficient by putting all the queue elements together.

    // ... or not? Maybe lets try sharing the visited array.
    for (const start of starts) {
        queue.push({
            position: start,
            steps: 0,
        });
    }

    // for{ position: start, steps: 0 }];
    function sortQueue() {
        // Make it A* by sorting the list based on steps and distance to end.
        queue.sort((a, b) => heuristic(a, end) - heuristic(b, end));
    }

    sortQueue();

    while (queue.length > 0) {
        const current = queue.shift()!;

        // Skip if visited
        if (visited[current.position.y][current.position.x]) {
            continue;
        }

        // Check if we are at the end
        if (current.position.x === end.x && current.position.y === end.y) {
            return current.steps;
        }

        // Mark as visited
        visited[current.position.y][current.position.x] = true;

        // Add all possible moves to the queue
        const currentElevation = getElevation(map[current.position.y][current.position.x]);
        for (const {x: dx, y: dy} of moveDirs) {
            const newPos = {
                x: current.position.x + dx,
                y: current.position.y + dy,
            };

            if (!isInBounds(map, newPos)) {
                continue;
            }

            if (visited[newPos.y][newPos.x]) {
                continue;
            }

            const newElevation = getElevation(map[newPos.y][newPos.x]);
            const delta = newElevation - currentElevation;

            // We can only move up one step, but we can jump down any number of steps.
            if (delta > 1) {
                continue;
            }

            queue.push({
                position: newPos,
                steps: current.steps + 1,
                parent: current,
            });
        }

        if (printedOnce) {
            // Move back up to where the drew the map last time.
            process.stdout.write('\x1b[0A'.repeat(map.length));
        }
        console.log(printMap(map, visited, queue, current));
        printedOnce = true;

        sortQueue();
    }
    return 0;
}

/**
 * Prints the map, using terminal colors to show visited and unvisited nodes.
 *
 * Unvistied squares are printed in dim grey.
 * Visited squares are printed with green.
 * Nodes in the queue are printed in blue.
 * The currently visited node is printed in bright red, and the path to it is printed in yellow.
 */
function printMap(map: string[][], visited: boolean[][], queue: QueueElem[], current: QueueElem): string {
    // Declare the colors we're going to use.
    const colors = {
        reset: '\x1b[0m',
        dim: '\x1b[2m',
        green: '\x1b[92m',
        blue: '\x1b[94m',
        brightRed: '\x1b[91;1m',
        yellow: '\x1b[93;1m',
    }
    // First create the map with colors based on visited or not
    const mapWithColors = map.map((row, y) =>
        row.map((value, x) => {
            const color = visited[y][x] ? colors.green : colors.dim;
            return color + value + colors.reset;
        }));

    // Then color the queue elements
    for (const { position } of queue) {
        // Use the original map value, not the colored one to remove the dim coloring.
        mapWithColors[position.y][position.x] =
            colors.blue +
            map[position.y][position.x] +
            colors.reset;
    }

    // Color the current node
    mapWithColors[current.position.y][current.position.x] =
        colors.brightRed +
        map[current.position.y][current.position.x] +
        colors.reset;

        // Then draw the path to the current node
    for (let parent = current.parent; parent; parent = parent.parent) {
        mapWithColors[parent.position.y][parent.position.x] =
            colors.yellow +
            map[parent.position.y][parent.position.x] +
            colors.reset;
    }

    // Return as a single string.
    return mapWithColors.map(row => row.join('')).join('\n');
}

function isInBounds(map: string[][], point: Point): boolean {
    return point.x >= 0 && point.x < map[0].length && point.y >= 0 && point.y < map.length;
}

function heuristic(queueElem: QueueElem, end: Point): number {
    const xDist = Math.abs(queueElem.position.x - end.x);
    const yDist = Math.abs(queueElem.position.y - end.y);
    return queueElem.steps + xDist + yDist;
}

function findLetters(map: string[][], letters: string[]): Point[] {
    const lettersSet = new Set(letters);

    const result: Point[] = [];
    for (let y = 0; y < map.length; y++) {
        for (let x = 0; x < map[0].length; x++) {
            if (lettersSet.has(map[y][x])) {
                result.push({ x, y });
            }
        }
    }
    return result;
}

function getElevation(value: string): number {
    if (value === 'S') {
        return getElevation('a');
    }
    if (value === 'E') {
        return getElevation('z');
    }

    return value.charCodeAt(0) - 'a'.charCodeAt(0);
}

function solve(filename: string) {
    const filepath = path.join(__dirname, filename);
    const inputText = fs.readFileSync(filepath, 'utf-8').trimEnd();

    const pt1 = countStepsToTop(inputText, ['S']);
    console.log(`Part 1: ${pt1}`);

    const pt2 = countStepsToTop(inputText, ['S', 'a']);
    console.log(`Part 1: ${pt2}`);
}

solve('demo.txt');
solve('input.txt');
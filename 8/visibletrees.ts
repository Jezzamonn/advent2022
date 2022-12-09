const fs = require('fs');
const path = require('path');

function countVisibleTrees(trees: number[][]) {
    const height = trees.length;
    const width = trees[0].length;

    const isVisible: boolean[][] = trees.map(row => row.map(() => false));

    for (let y = 0; y < height; y++) {
        // Looking from the left
        {
            let maxHeight = -1;
            for (let x = 0; x < width; x++) {
                if (trees[y][x] > maxHeight) {
                    maxHeight = trees[y][x];
                    isVisible[y][x] = true;
                }
            }
        }
        // Looking from the right
        {
            let maxHeight = -1;
            for (let x = width - 1; x >= 0; x--) {
                if (trees[y][x] > maxHeight) {
                    maxHeight = trees[y][x];
                    isVisible[y][x] = true;
                }
            }
        }
    }

    for (let x = 0; x < width; x++) {
        // Looking from the top
        {
            let maxHeight = -1;
            for (let y = 0; y < height; y++) {
                if (trees[y][x] > maxHeight) {
                    maxHeight = trees[y][x];
                    isVisible[y][x] = true;
                }
            }
        }
        // Looking from the bottom
        {
            let maxHeight = -1;
            for (let y = height - 1; y >= 0; y--) {
                if (trees[y][x] > maxHeight) {
                    maxHeight = trees[y][x];
                    isVisible[y][x] = true;
                }
            }
        }
    }

    // printTrees(trees, isVisible);

    // Count visible trees
    return {
        numVisible: isVisible.map(row => row.filter(visible => visible).length).reduce((a, b) => a + b),
        isVisible,
    }
}

function parseTrees(inputText: string): number[][] {
    return inputText.split('\n').filter(line => line.length > 0).map(line => line.split('').map(c => parseInt(c)));
}

function findMostScenicScore(trees: number[][]) {
    const height = trees.length;
    const width = trees[0].length;

    let bestScore = -1;
    let bestCoord = { x: -1, y: -1 };
    for (let y = 0; y < height; y++) {
        for (let x = 0; x < width; x++) {
            const score = calculateScenicScore(trees, x, y);
            if (score > bestScore) {
                bestScore = score;
                bestCoord = { x, y };
            }
        }
    }
    return { bestScore, bestCoord };
}

function calculateScenicScore(trees: number[][], x: number, y: number) {
    // This may be slow to calculate for each tree. We'll see.
    let startHeight = trees[y][x];

    const dirs = [
        { x: 1, y: 0 },
        { x: -1, y: 0 },
        { x: 0, y: 1 },
        { x: 0, y: -1 },
    ];

    let score = 1;
    for (let dir of dirs) {
        let scoreForThisDirection = 0;
        for (let i = 1; true; i++) {
            let xx = x + dir.x * i;
            let yy = y + dir.y * i;
            if (xx < 0 || xx >= trees[0].length || yy < 0 || yy >= trees.length) {
                scoreForThisDirection = i - 1;
                break;
            }
            if (trees[yy][xx] >= startHeight) {
                scoreForThisDirection = i;
                break;
            }
        }
        score *= scoreForThisDirection;
    }
    return score;
}

/**
 * Prints an ascii representation of the trees and their visibility.
 *
 * Uses the regular console color for visible trees and the dim color for
 * invisible trees. The tree with the best score is printed in green.
 */
function printTrees(trees: number[][], isVisible: boolean[][], bestCoord: { x: number, y: number }) {
    const height = trees.length;
    const width = trees[0].length;

    const dim = '\x1b[2m';
    const reset = '\x1b[0m';
    const green = '\x1b[32m';

    for (let y = 0; y < height; y++) {
        let line = '';
        for (let x = 0; x < width; x++) {
            if (x === bestCoord.x && y === bestCoord.y) {
                line += green;
            } else if (isVisible[y][x]) {
                line += '';
            } else {
                line += dim;
            }
            line += trees[y][x];
            line += reset;
        }
        console.log(line);
    }
}

function solve() {
    const filename = path.join(__dirname, 'input.txt');
    const text = fs.readFileSync(filename, 'utf-8');

    const trees = parseTrees(text);

    const {numVisible, isVisible} = countVisibleTrees(trees);

    const {bestScore, bestCoord} = findMostScenicScore(trees);

    printTrees(trees, isVisible, bestCoord);

    console.log(`Part 1: ${numVisible}`);
    console.log(`Part 2: ${bestScore}`);
}

solve();
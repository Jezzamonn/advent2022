const fs = require('fs');
const path = require('path');

/**
 * Simulates the movement of the rope and returns the number of visited squares of the end of the rope.
 */
async function simulateMovement(inputText, numKnots) {
    let rope = Array(numKnots).fill(null).map(_ => ({ x: 0, y: 0 }));
    rope.last = () => rope[rope.length - 1];
    const visitedSquares = new Set(['0,0']);
    let bounds = {
        minX: 0,
        maxX: 0,
        minY: 0,
        maxY: 0,
    }

    for (const movement of inputText.split('\n')) {
        const [direction, amountStr] = movement.split(' ');
        const amount = parseInt(amountStr);

        // Move 'amount' steps, individually.
        for (let i = 0; i < amount; i++) {
            await moveStep(rope, direction, bounds, visitedSquares);
            visitedSquares.add(`${rope.last().x},${rope.last().y}`);

            if (animate && !animateSteps) {
                await renderFrame(bounds, rope, visitedSquares);
            }
        }
    }
    return visitedSquares.size;
}

function updateBounds(rope, bounds) {
    bounds.minX = Math.min(bounds.minX, ...rope.map(r => r.x));
    bounds.maxX = Math.max(bounds.maxX, ...rope.map(r => r.x));
    bounds.minY = Math.min(bounds.minY, ...rope.map(r => r.y));
    bounds.maxY = Math.max(bounds.maxY, ...rope.map(r => r.y));
}

async function moveStep(rope, direction, bounds, visitedSquares) {
    // Move the head of the rope
    switch (direction) {
        case 'U':
            rope[0].y--;
            break;
        case 'D':
            rope[0].y++;
            break;
        case 'L':
            rope[0].x--;
            break;
        case 'R':
            rope[0].x++;
            break;
    }
    updateBounds([rope[0]], bounds);

    if (animateSteps) {
        await renderFrame(bounds, rope, visitedSquares);
    }

    // Update each know in the rope.
    for (let i = 1; i < rope.length; i++) {
        const forward = rope[i - 1];
        const current = rope[i];
        if (Math.abs(forward.x - current.x) >= 2 || Math.abs(forward.y - current.y) >= 2) {
            // Move one step, potentially diagonally, towards the forward knot.
            if (current.x < forward.x) {
                current.x++;
            }
            else if (current.x > forward.x) {
                current.x--;
            }
            if (current.y < forward.y) {
                current.y++;
            }
            else if (current.y > forward.y) {
                current.y--;
            }
        }

        if (animateSteps) {
            await renderFrame(bounds, rope, visitedSquares);
        }
    }
}

async function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// Draws one animation frame to the console, waiting a bit so that the animation isn't too fast.
async function renderFrame(bounds, rope, visitedSquares) {
    // Instead of clearing the console, move the cursor to the top left corner.
    // This way, the animation doesn't flicker.
    console.log('\x1b[0;0H');
    console.log(renderRope(bounds, rope, visitedSquares));
    await sleep(30);
}

/**
 * Draw an ascii art representation of the rope, along with the visited squares.
 *
 * Uses console colors to make it easier to see the rope and the visited squares
 */
function renderRope(bounds, rope, visitedSquares) {
    // Round up the bounds to the nearest power of 2 so that the size of the
    // ascii art doesn't change too often.
    // Also, use a minimum size of 8 for each dimension.
    const roundedBounds = {};
    for (const [key, value] of Object.entries(bounds)) {
        const log = Math.log2(Math.abs(value));
        const rounded = Math.pow(2, Math.ceil(log));
        const roundedWithMin = Math.max(rounded, 8);
        roundedBounds[key] = key.startsWith('min') ? -roundedWithMin : roundedWithMin;
    }

    const width = roundedBounds.maxX - roundedBounds.minX + 1;
    const height = roundedBounds.maxY - roundedBounds.minY + 1;
    const origin = { x: -roundedBounds.minX, y: -roundedBounds.minY };

    const output2dArray = Array(height).fill(null).map(() => Array(width).fill(' '));

    // Define the escape sequences for the colors we plan to use.
    // The visited squares are drawn in dim grey
    // The origin is drawn in bright green.
    // The head of the rope is drawn in bright red, and the rest of the rope is drawn in yellow.
    const colors = {
        reset: '\x1b[0m',
        visited: '\x1b[2m\x1b[90m',
        origin: '\x1b[92m',
        ropeHead: '\x1b[91m',
        rope: '\x1b[93m',
    };

    // Draw the visited squares
    for (const square of visitedSquares) {
        const [x, y] = square.split(',').map(n => parseInt(n));
        output2dArray[y + origin.y][x + origin.x] = `${colors.visited}#${colors.reset}`;
    }
    // Draw the origin square with an 's'
    output2dArray[origin.y][origin.x] = `${colors.origin}s${colors.reset}`;
    // Draw the rope, backwards so that the head is drawn last. Using 'H' for the head, and numbers for the rest of the rope.
    for (let i = rope.length - 1; i >= 0; i--) {
        const { x, y } = rope[i];
        const color = i === 0 ? colors.ropeHead : colors.rope;
        output2dArray[y + origin.y][x + origin.x] = `${color}${i === 0 ? 'H' : i}${colors.reset}`;
    }

    // Return as a string
    return output2dArray.map(row => row.join('')).join('\n');
}

async function solve(filename) {
    const filepath = path.join(__dirname, filename);
    const text = fs.readFileSync(filepath, 'utf-8').trimEnd();

    console.log(`${filename}:`);
    // const pt1 = await simulateMovement(text, 2);
    // console.log('Part 1:', pt1);

    const pt2 = await simulateMovement(text, 10);
    console.log('Part 2:', pt2);
}

async function solveAll() {
    // await solve('demo.txt');
    await solve('input.txt');
}

let animateSteps = false;
let animate = false;

solveAll();
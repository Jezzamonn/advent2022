const fs = require('fs');
const path = require('path');

/**
 * Simulates the movement of the rope and returns the number of visited squares of the end of the rope.
 */
async function simulateMovement(inputText) {
    let head = { x: 0, y: 0 };
    let tail = { x: 0, y: 0 };
    const visitedSquares = new Set([`${tail.x},${tail.y}`]);
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
            moveStep(head, tail, direction);
            visitedSquares.add(`${tail.x},${tail.y}`);

            bounds.minX = Math.min(bounds.minX, head.x, tail.x);
            bounds.maxX = Math.max(bounds.maxX, head.x, tail.x);
            bounds.minY = Math.min(bounds.minY, head.y, tail.y);
            bounds.maxY = Math.max(bounds.maxY, head.y, tail.y);

            await renderFrame(bounds, head, tail, visitedSquares);
        }
    }
    return visitedSquares.size;
}

function moveStep(head, tail, direction) {
    // Move the head of the rope
    switch (direction) {
        case 'U':
            head.y--;
            break;
        case 'D':
            head.y++;
            break;
        case 'L':
            head.x--;
            break;
        case 'R':
            head.x++;
            break;
    }
    // Update the tail
    if (Math.abs(head.x - tail.x) >= 2) {
        tail.x += (tail.x < head.x) ? 1 : -1;
        tail.y = head.y;
    }
    else if (Math.abs(head.y - tail.y) >= 2) {
        tail.y += (tail.y < head.y) ? 1 : -1;
        tail.x = head.x;
    }
}

async function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// Draws one animation frame to the console, waiting a bit so that the animation isn't too fast.
async function renderFrame(bounds, head, tail, visitedSquares) {
    if (!animate) {
        return;
    }
    // Instead of clearing the console, move the cursor to the top left corner.
    // This way, the animation doesn't flicker.
    console.log('\x1b[0;0H');
    console.log(renderRope(bounds, head, tail, visitedSquares));
    await sleep(30);
}

/**
 * Draw an ascii art representation of the rope, along with the visited squares.
 *
 * Uses console colors to make it easier to see the rope and the visited squares.
 */
function renderRope(bounds, head, tail, visitedSquares) {
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

    // Define the escape sequences for the colors we plan to use
    const colors = {
        head: '\x1b[31m\x1b[1m',
        tail: '\x1b[34m\x1b[1m',
        visited: '\x1b[32m\x1b[1m',
        origin: '\x1b[33m\x1b[1m',
        reset: '\x1b[0m',
    };

    // Draw the visited squares
    for (const square of visitedSquares) {
        const [x, y] = square.split(',').map(n => parseInt(n));
        output2dArray[y + origin.y][x + origin.x] = `${colors.visited}#${colors.reset}`;
    }
    // Draw the origin square with an 's'
    output2dArray[origin.y][origin.x] = `${colors.origin}s${colors.reset}`;
    // Draw the tail, and then the head
    output2dArray[tail.y + origin.y][tail.x + origin.x] = `${colors.tail}T${colors.reset}`;
    output2dArray[head.y + origin.y][head.x + origin.x] = `${colors.head}H${colors.reset}`;

    // Return as a string
    return output2dArray.map(row => row.join('')).join('\n');
}

async function solve(filename) {
    const filepath = path.join(__dirname, filename);
    const text = fs.readFileSync(filepath, 'utf-8').trimEnd();

    console.log(`${filename}:`);
    const pt1 = await simulateMovement(text);
    console.log('Part 1:', pt1);
}

let animate = false;

solve('demo.txt').then(
    () => solve('input.txt')
)
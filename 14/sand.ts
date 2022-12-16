import { promises as fs } from 'fs';
import * as path from 'path';

const EMPTY = 0;
const ROCK = 1;
const SAND = 2;

// Console colors:
const colors = {
    reset: '\x1b[0m',
    dim: '\x1b[2m',
    yellow: '\x1b[33m',
    brightYellow: '\x1b[93;1m',
};

/**
 * Map values to symbols for printing. Uses console colors to make things clearer.
 *
 * Rock is printed in regular white. Sand is printed in yellow.
 */
const SYMBOLS: { [key: number]: string } = {
    [EMPTY]: ' ',
    [ROCK]: '#',
    [SAND]: colors.brightYellow + 'o' + colors.reset,
}

let debug = false;

async function sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

class Grid {
    private grid: Uint8Array;
    readonly minX: number;
    readonly maxX: number;
    readonly minY: number;
    readonly maxY: number;
    readonly width: number;
    readonly height: number;
    private printed = false;

    constructor(minX: number, maxX: number, minY: number, maxY: number) {
        this.minX = minX;
        this.maxX = maxX;
        this.minY = minY;
        this.maxY = maxY;
        this.width = maxX - minX + 1;
        this.height = maxY - minY + 1;
        this.grid = new Uint8Array(this.width * this.height);
    }

    get(x: number, y: number): number {
        const index = this.getIndex(x, y);
        return this.grid[index];
    }

    set(x: number, y: number, value: number) {
        const index = this.getIndex(x, y);
        this.grid[index] = value;
    }

    private getIndex(x: number, y: number): number {
        if (x < this.minX || x > this.maxX || y < this.minY || y > this.maxY) {
            throw new Error(`Out of bounds: (${x}, ${y})`);
        }
        const row = y - this.minY;
        const col = x - this.minX;
        const index = row * this.width + col;
        return index;
    }

    print() {
        for (let y = this.minY; y <= this.maxY; y++) {
            let line = '';
            for (let x = this.minX; x <= this.maxX; x++) {
                const value = this.get(x, y);
                line += SYMBOLS[value];
            }
            console.log(line);
        }
    }

    async animate() {
        if (this.printed) {
            this.resetConsolePosition();
        }
        this.print();
        this.printed = true;
        // await sleep(10);
    }

    resetConsolePosition() {
        // Move the cursor back to the top of where we printed.
        const lines = this.maxY - this.minY + 1;
        const up = '\x1b[' + lines + 'A';
        process.stdout.write(up);
    }

    static createFromPaths(rockPaths: string, {withFloor = false}: {withFloor?: boolean} = {}): Grid {
        const paths = parsePaths(rockPaths);
        const minY = Math.min(0, ...paths.map(path => Math.min(...path.map(p => p.y))));
        let maxY = Math.max(...paths.map(path => Math.max(...path.map(p => p.y))));

        if (withFloor) {
            // Add an extra path at the bottom for the floor. Make it wide enough to catch all the sand.
            const floorHeight = maxY + 2;
            const centerX = 500;
            paths.push([
                { x: centerX - floorHeight, y: floorHeight },
                { x: centerX + floorHeight, y: floorHeight }]);
            maxY = floorHeight;
        }

        const minX = Math.min(...paths.map(path => Math.min(...path.map(p => p.x))));
        const maxX = Math.max(...paths.map(path => Math.max(...path.map(p => p.x))));
        // Add extra spaces on each side to account for the sand falling off the
        // rocks, and some extra at the bottom.
        const grid = new Grid(minX - 1, maxX + 1, minY, maxY + 1);

        for (const path of paths) {
            for (let p = 0; p < path.length - 1; p++) {
                const p1 = path[p];
                const p2 = path[p + 1];
                const dx = Math.sign(p2.x - p1.x);
                const dy = Math.sign(p2.y - p1.y);
                for (let x = p1.x, y = p1.y;
                    x !== p2.x || y !== p2.y;
                    x += dx, y += dy) {
                    grid.set(x, y, ROCK);
                }
                grid.set(p2.x, p2.y, ROCK);
            }
        }
        return grid;
    }
}

function parsePaths(rockPaths: string): { x: number; y: number; }[][] {
    return rockPaths.split('\n').map(
        path => path.split(' -> ').map(point => {
            const [x, y] = point.split(',').map(s => parseInt(s));
            return { x, y };
        })
    );
}

async function simulateSand(grid: Grid): Promise<number> {
    // Simulate sand falling!
    let sandCount = 0;
    sandLoop: while (true) {
        // Add a new sand particle at the top.
        const sandPos = { x: 500, y: 0};
        sandCount++;

        if (debug) {
            await grid.animate();
        }

        while (true) {
            // Set the position in the grid, mostly so we can print it.
            grid.set(sandPos.x, sandPos.y, SAND);
            // if (debug) {
            //     await grid.animate();
            // }

            // Check if we've reached the bottom.
            if (sandPos.y === grid.maxY) {
                // We're done!
                sandCount--;
                break sandLoop;
            }


            // Try fall directly down.
            if (grid.get(sandPos.x, sandPos.y + 1) === EMPTY) {
                grid.set(sandPos.x, sandPos.y, EMPTY);
                sandPos.y++;
                continue;
            }
            // Try fall down and to the left.
            if (grid.get(sandPos.x - 1, sandPos.y + 1) === EMPTY) {
                grid.set(sandPos.x, sandPos.y, EMPTY);
                sandPos.x--;
                sandPos.y++;
                continue;
            }
            // Try fall down and to the right.
            if (grid.get(sandPos.x + 1, sandPos.y + 1) === EMPTY) {
                grid.set(sandPos.x, sandPos.y, EMPTY);
                sandPos.x++;
                sandPos.y++;
                continue;
            }
            // Nowhere left to fall! We're stuck. If this grain is at the start,
            // we're done, otherwise move to the next grain.
            if (sandPos.y === 0) {
                break sandLoop;
            }
            break;
        }
    }

    return sandCount;
}

async function solve(filename: string) {
    const filepath = path.join(__dirname, filename);
    const inputText = (await fs.readFile(filepath, 'utf-8')).trimEnd();

    const grid1 = Grid.createFromPaths(inputText);
    const pt1 = await simulateSand(grid1);
    console.log(`Part 1: ${pt1}`);

    const grid2 = Grid.createFromPaths(inputText, { withFloor: true });
    const pt2 = await simulateSand(grid2);
    console.log(`Part 2: ${pt2}`);
}

solve('input.txt');
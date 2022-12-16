import { promises as fs } from 'fs';
import * as path from 'path';

class Range {
    start: number;
    end: number;

    constructor(start: number, end: number) {
        this.start = start;
        this.end = end;
    }

    get length(): number {
        return this.end - this.start + 1;
    }

    isOverlapping(other: Range): boolean {
        return this.start <= other.end && this.end >= other.start;
    }

    merge(other: Range): Range {
        return new Range(Math.min(this.start, other.start), Math.max(this.end, other.end));
    }

    /** Returns a new range if this range can be clamped. If this range is totally outside min and max, returns undefined */
    clamp(min: number, max: number): Range | undefined {
        if (this.start > max || this.end < min) {
            return undefined;
        }
        return new Range(Math.max(this.start, min), Math.min(this.end, max));
    }
}

class Sensor {
    x = 0
    y = 0
    beacon = { x: 0, y: 0 }

    // Example line:
    // Sensor at x=2288642, y=2282562: closest beacon is at x=1581951, y=2271709
    static parse(line: string): Sensor {
        const sensor = new Sensor();
        const parts = line.split(' ');
        sensor.x = parseInt(parts[2].split('=')[1]);
        sensor.y = parseInt(parts[3].split('=')[1]);
        sensor.beacon.x = parseInt(parts[8].split('=')[1]);
        sensor.beacon.y = parseInt(parts[9].split('=')[1]);
        return sensor;
    }

    distToBeacon(): number {
        return Math.abs(this.x - this.beacon.x) + Math.abs(this.y - this.beacon.y);
    }

    knownRangeAtRow(row: number): Range | undefined {
        const sensorDist = this.distToBeacon();
        const distToRow = Math.abs(this.y - row);
        const rangeAtRow = sensorDist - distToRow;
        if (rangeAtRow < 0) {
            return undefined;
        }
        return new Range(this.x - rangeAtRow, this.x + rangeAtRow);
    }

}

function getKnowPositionsRangesAtRow(sensors: Sensor[], { row, minX = -Infinity, maxX = Infinity}: { row: number; minX?: number; maxX?: number }): Range[] {
    const ranges = sensors
        .map((sensor) => sensor.knownRangeAtRow(row))
        .map((range) => range?.clamp(minX, maxX))
        .filter((range) => range !== undefined) as Range[];

    // Merge ranges
    ranges.sort((a, b) => a.start - b.start);
    let mergedRanges: Range[] = [];
    for (const range of ranges) {
        if (mergedRanges.length === 0) {
            mergedRanges.push(range);
            continue;
        }
        const lastRange = mergedRanges[mergedRanges.length - 1];
        if (lastRange.isOverlapping(range)) {
            mergedRanges[mergedRanges.length - 1] = lastRange.merge(range);
        } else {
            mergedRanges.push(range);
        }
    }
    return mergedRanges;
}

function countKnownPositionsAtRow(
    inputText: string,
    {
        row,
        includeBeacons = false,
        minX = -Infinity,
        maxX = Infinity,
    }: { row: number; includeBeacons?: boolean; minX?: number; maxX?: number }
): number {
    const sensors = inputText.split('\n').map((line) => Sensor.parse(line));
    const mergedRanges: Range[] = getKnowPositionsRangesAtRow(sensors, { row, minX, maxX });

    let beaconsCount = 0;
    if (!includeBeacons) {
        // Need to subtract the beacons in this row too.
        const beaconXCoords = new Set(
            sensors
                .filter((sensor) => sensor.beacon.y === row)
                .map((sensor) => sensor.beacon.x)
        );
        beaconsCount = beaconXCoords.size;
    }

    // Count lengths
    return (
        mergedRanges.reduce((sum, range) => sum + range.length, 0) -
        beaconsCount
    );
}

/** Returns the beacon's position as a single number, multiplying its x coordinate by 4000000 and then adding its y coordinate. */
function findHiddenBeacon(inputText: string): number {
    const sensors = inputText.split('\n').map((line) => Sensor.parse(line));

    const size = 4000000;
    for (let row = 0; row < size; row++) {
        const ranges = getKnowPositionsRangesAtRow(sensors, {
            row: row,
            minX: 0,
            maxX: size,
        });
        if (ranges.length !== 1) {
            // We found our beacon.
            const beaconX = ranges[0].end + 1;

            return beaconX * size + row;
        }
    }
    return -1;
}


async function solve(filename: string, row: number) {
    const filepath = path.join(__dirname, filename);
    const inputText = (await fs.readFile(filepath, 'utf-8')).trimEnd();

    console.log(`${filename}:`);

    const pt1 = countKnownPositionsAtRow(inputText, {row});
    console.log(`Part 1: ${pt1}`);

    const p2 = findHiddenBeacon(inputText);
    console.log(`Part 2: ${p2}`);
}

Promise.resolve()
    // .then(() => solve('test.txt', 18))
    // .then(() => solve('demo.txt', 10))
    .then(() => solve('input.txt', 2000000));
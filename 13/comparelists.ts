import { promises as fs } from 'fs';
import * as path from 'path';

type Element = number | Array<Element>;

const log = (msg: string = '') => null;
// const log = console.log;

function countOrderedPairs(inputText: string): number {
    let count = 0;
    const pairs = inputText.split('\n\n');
    for (let p = 0; p < pairs.length; p++) {
        log(`== Pair ${p + 1} ==`)
        const [list1, list2] = pairs[p].split('\n').map(list => JSON.parse(list));
        if (compareElements(list1, list2) < 0) {
            count += (p + 1);
        }
        log();
    }
    return count;
}

function sortAllAndReturnMarkerIndices(inputText: string): number {
    const lists = inputText.split('\n')
        .filter(line => line.length > 0)
        .map(list => JSON.parse(list));

    // Add our two markers
    const marker1 = [[2]];
    const marker2 = [[6]];
    lists.push(marker1);
    lists.push(marker2);

    lists.sort(compareElements);

    return (lists.indexOf(marker1) + 1) * (lists.indexOf(marker2) + 1);
}

/**
 * Compare two elements. If elem1 comes first, return -1. If elem2 comes first,
 * return 1. If they are equal, return 0.
 *
 * Depth is just used for logging.
 */
function compareElements(elem1: Element, elem2: Element, depth=0): number {
    const indentation = ' '.repeat(depth);
    log(`${indentation}- Comparing ${JSON.stringify(elem1)} vs ${JSON.stringify(elem2)}`);
    if (typeof elem1 === 'number' && typeof elem2 === 'number') {
        const diff = elem1 - elem2;
        if (diff < 0) {
            log(`${indentation}  - elem1 is smaller, right order`);
        }
        if (diff > 0) {
            log(`${indentation}  - elem2 is larger, wrong order`);
        }
        return elem1 - elem2;
    }
    if (Array.isArray(elem1) && Array.isArray(elem2)) {
        // Loop through and compare each element one by one.
        for (let i = 0; i < Math.max(elem1.length, elem2.length); i++) {
            if (i >= elem1.length) {
                log(`${indentation}  - elem1 is shorter, right order`)
                return -1;
            }
            if (i >= elem2.length) {
                log(`${indentation}  - elem2 is shorter, wrong order`)
                return 1;
            }
            const diff = compareElements(elem1[i], elem2[i], depth + 1);
            // If the elements are not equal, return the comparison of those
            // elements.
            if (diff !== 0) {
                return diff;
            }
        }
    }
    if (typeof elem1 === 'number') {
        return compareElements([elem1], elem2, depth + 1);
    }
    if (typeof elem2 === 'number') {
        return compareElements(elem1, [elem2], depth + 1);
    }
    // Not really possible to get here, but TypeScript doesn't know that.
    return 0;
}

async function solve(filename: string) {
    const filepath = path.join(__dirname, filename);
    const inputText = (await fs.readFile(filepath, 'utf-8')).trimEnd();

    console.log(`${filename}:`)

    const pt1 = countOrderedPairs(inputText);
    console.log(`Part 1: ${pt1}`);

    const pt2 = sortAllAndReturnMarkerIndices(inputText);
    console.log(`Part 2: ${pt2}`);
}

solve('demo.txt')
    .then(() => solve('input.txt'));

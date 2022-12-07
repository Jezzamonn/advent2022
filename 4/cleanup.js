const fs = require('fs');
const path = require('path');

function countTotallyContainedSchedules(inputText) {
    let numOverlapping = 0;

    for (const line of inputText.split('\n')) {
        const [first, second] = line.split(',').map(parseInterval);
        const overlapping = intervalTotallyContained(first, second);
        if (overlapping) {
            numOverlapping++;
        }
    }
    return numOverlapping;
}

function countPartiallyOverlappingSchedules(inputText) {
    let numOverlapping = 0;

    for (const line of inputText.split('\n')) {
        const [first, second] = line.split(',').map(parseInterval);
        const overlapping = intervalsPartiallyOverlap(first, second);
        if (overlapping) {
            numOverlapping++;
        }
    }
    return numOverlapping;
}

function parseInterval(intervalStr) {
    const [startStr, endStr] = intervalStr.split('-');
    return {
        start: parseInt(startStr),
        end: parseInt(endStr),
    };
}

function intervalsPartiallyOverlap(first, second) {
    return first.start <= second.end && second.start <= first.end;
}

function intervalTotallyContained(longest, shortest) {
    if (longest.end - longest.start < shortest.end - shortest.start) {
        // Shortest was actually the longest.
        return intervalTotallyContained(shortest, longest);
    }

    return longest.start <= shortest.start && longest.end >= shortest.end;
}

function visualizeInterval(interval) {
    const max = 100;
    const arr = new Array(max).fill(' ');
    for (let i = interval.start; i <= interval.end; i++) {
        arr[i] = '-';
    }
    arr[interval.start] = '|';
    arr[interval.end] = '|';
    return arr.join('');
}

function solve() {
    const filename = path.join(__dirname, 'input.txt');
    const text = fs.readFileSync(filename, 'utf-8').trim();

    const pt1 = countTotallyContainedSchedules(text);
    console.log('Pt 1: ', pt1);

    const pt2 = countPartiallyOverlappingSchedules(text);
    console.log('Pt 2: ', pt2);
}

solve();
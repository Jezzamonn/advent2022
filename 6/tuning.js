const fs = require('fs');
const path = require('path');

function findStartOfMessage(message, markerLength) {
    for (let i = 0; i < message.length - (markerLength - 1); i++) {
        const subStr = message.slice(i, i + markerLength);
        if (!containsDuplicateCharacters(subStr)) {
            return i + (markerLength - 1) + 1; // Convert to 1-based indexing
        }
    }
}

function containsDuplicateCharacters(message) {
    const sorted = Array.from(message).sort();
    for (let i = 0; i < sorted.length - 1; i++) {
        if (sorted[i] === sorted[i + 1]) {
            return true;
        }
    }
    return false;
}

function solve() {
    const filename = path.join(__dirname, 'input.txt');
    const text = fs.readFileSync(filename, 'utf-8');

    const pt1 = findStartOfMessage(text, 4);
    console.log('Part 1:', pt1);

    const pt2 = findStartOfMessage(text, 14);
    console.log('Part 2:', pt2);
}

solve();